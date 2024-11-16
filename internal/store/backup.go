package store

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/julianstephens/warden/internal/chunker"
	"github.com/julianstephens/warden/internal/crypto"
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

	fmt.Printf("%-64s  %s\n", "HASH", "CHUNK SIZE")
	for _, p := range pathsToBackup {
		chunkAndHash(store, p)
	}

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

func chunkAndHash(store *Store, filepath string) (err error) {
	warden.Log.Debug().Msgf("checking file %s exists...", filepath)

	if _, err = os.Stat(filepath); err == nil {
		var file *os.File
		file, err = os.Open(filepath)
		if err != nil {
			return
		}
		defer file.Close()

		warden.Log.Debug().Msg("chunking and hashing file...")
		cKr := chunker.NewChunker(file)

		for {
			var chunk chunker.Chunk
			chunk, err = cKr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}

			hashedChunk := crypto.SecureHash(chunk.Data, store.master.user.Data)
			fmt.Println(hashedChunk)
			// TODO: check if chunk has already been backed up
		}
	} else {
		warden.Log.Info().Msgf("file %s does not exist. skipping...", filepath)
	}

	return
}
