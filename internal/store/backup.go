package store

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/julianstephens/warden/internal/storage"
	"github.com/julianstephens/warden/internal/warden"
)

func (s *Store) Backup(ctx context.Context, backupDir string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if backupDir == s.Location {
		return fmt.Errorf("cannot backup warden store")
	}
	err := backup(s, ctx, backupDir)
	if err != nil {
		return fmt.Errorf("unable to backup dir %s: %+v", backupDir, err)
	}

	return nil
}

func backup(store *Store, ctx context.Context, backupDir string) (err error) {
	latestSnapshot, err := getLastestSnapshot(store, ctx, backupDir)
	if err != nil {
		err = fmt.Errorf("unable to retrieve latest snapshot for backup dir %s: %+v", backupDir, err)
		return
	}

	// TODO: chunk and hash
	pathsToBackup, pathsToCopy, err := sortBackupPaths(latestSnapshot, backupDir)
	if err != nil {
		return
	}
	fmt.Println(pathsToBackup)
	fmt.Println(pathsToCopy)

	return
}

func getLastestSnapshot(store *Store, ctx context.Context, backupDir string) (snap *storage.Snapshot, err error) {
	snapshots, err := store.backend.ListSnapshots(ctx)
	if err != nil {
		return
	}
	fmt.Print(len(snapshots))

	backupSnaps := warden.Filter(snapshots, func(t storage.Snapshot) bool {
		return t.BackupVolume == backupDir
	})

	if len(backupSnaps) == 0 {
		return
	}

	sort.Slice(backupSnaps, func(i, j int) bool {
		return backupSnaps[i].CreatedAt.After(backupSnaps[j].CreatedAt)
	})
	snap = &backupSnaps[0]

	return
}

func sortBackupPaths(latestSnapshot *storage.Snapshot, backupDir string) (pathsToBackup []string, pathsToCopy []string, err error) {
	err = filepath.WalkDir(backupDir, func(path string, entry fs.DirEntry, err error) error {
		if entry.IsDir() {
			return nil
		}

		if latestSnapshot != nil {
			for _, p := range latestSnapshot.Paths {
				if p.Path == path {
					entryInfo, err := entry.Info()
					if err != nil || p.FileSize != entryInfo.Size() || !p.ModifiedAt.Equal(entryInfo.ModTime()) {
						pathsToBackup = append(pathsToBackup, path)
						return nil
					} else {
						pathsToCopy = append(pathsToCopy, path)
						return nil
					}
				}
			}
		}

		pathsToBackup = append(pathsToBackup, path)
		return nil
	})

	return
}
