package storage_test

import (
	"os"
	"testing"
	"time"

	mp "github.com/agiledragon/gomonkey/v2"

	"github.com/julianstephens/warden/internal/storage"
	"github.com/julianstephens/warden/internal/warden"
)

const testVolume = "/tmp/notes"

func initTestVolume(t *testing.T) {
	if err := os.Mkdir(testVolume, os.ModePerm); err != nil {
		t.Fatalf("unable to create test volume: %+v", err)
	}
}

func TestSnapshot(t *testing.T) {
	initTestVolume(t)

	testTime := time.Now().Add(-(10 * time.Minute))
	patches := mp.NewPatches()
	patches.ApplyFuncReturn(time.Now, testTime)
	defer patches.Reset()

	_, err := storage.NewSnapshot("")
	if err == nil {
		t.Fatal("should error on empty backup :volume")
	}

	_, err = storage.NewSnapshot("/path/to/nowhere")
	if err == nil {
		t.Fatal("should error on nonexistent backup volume")
	}

	snapshot, err := storage.NewSnapshot(testVolume)
	if err != nil || snapshot == nil {
		t.Fatal("should create a new snapshot instead nil or error")
	}

	if snapshot.BackupVolume != testVolume {
		t.Fatalf("expected backup volume %s, got %s", testVolume, snapshot.BackupVolume)
	}
	if !snapshot.CreatedAt.Equal(testTime) {
		t.Fatalf("expected created at %s, got %s", testTime.String(), snapshot.CreatedAt.String())
	}
	if len(snapshot.Paths) != 0 {
		t.Fatalf("expected empty paths slice, got len %d", len(snapshot.Paths))
	}
	if len(snapshot.PackedChunks) != 0 {
		t.Fatalf("expected empty packed chunk slice, got len %d", len(snapshot.PackedChunks))
	}

	os.RemoveAll(testVolume)
}

func assertChunksEqual(t *testing.T, chunk1, chunk2 storage.PackedChunk) {
	if chunk1.Chunk != chunk2.Chunk {
		t.Fatalf("chunk %s does not equal %s", chunk1.Chunk, chunk2.Chunk)
	}
	if chunk1.Pack != chunk2.Pack {
		t.Fatalf("pack %s does not equal %s", chunk1.Pack, chunk2.Pack)
	}
	if chunk1.ChunkStart != chunk2.ChunkStart {
		t.Fatalf("chunk start %d does not equal %d", chunk1.ChunkStart, chunk2.ChunkStart)
	}
	if chunk1.ChunkEnd != chunk2.ChunkEnd {
		t.Fatalf("chunk end %d does not equal %d", chunk1.ChunkEnd, chunk2.ChunkEnd)
	}
}

func TestGetPack(t *testing.T) {
	initTestVolume(t)

	snapshot, err := storage.NewSnapshot(testVolume)
	if err != nil || snapshot == nil {
		t.Fatal("should create a new snapshot instead nil or error")
	}

	testChunkHash := warden.NewID().String()
	testPackHash := warden.NewID().String()
	testChunkStart := 0
	testChunkEnd := 1024
	testChunk := storage.PackedChunk{Chunk: testChunkHash, Pack: testPackHash, ChunkStart: int64(testChunkStart), ChunkEnd: int64(testChunkEnd)}

	snapshot.PackedChunks = append(snapshot.PackedChunks,
		testChunk,
		storage.PackedChunk{Chunk: "xxx", Pack: "xxx", ChunkStart: int64(testChunkStart + 1024), ChunkEnd: int64(testChunkEnd + 1024)},
	)

	res := snapshot.GetPack(testChunkHash)
	assertChunksEqual(t, testChunk, *res)

	res = snapshot.GetPack("nonexistent")
	if res != nil {
		t.Fatalf("expected nil, got %+v", res)
	}

	os.RemoveAll(testVolume)
}
