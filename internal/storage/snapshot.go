package storage

import (
	"time"
)

type PathMetadata struct {
	Path       string    `json:"path"`
	FileSize   int64     `json:"fileSize"`
	FilePerm   string    `json:"filePermission"`
	ModifiedAt time.Time `json:"modifiedAt"`

	Chunks []string `json:"chunks"`
}

type ChunkLoc struct {
	Chunk string `json:"chunk"`
	Pack  string `json:"pack"`

	ChunkStart int64 `json:"chunkStartOffset"`
	ChunkEnd   int64 `json:"chunkEndOffset"`
}

type Snapshot struct {
	BackupVolume string         `json:"backupVolume"`
	Paths        []PathMetadata `json:"paths"`

	ChunkLocs []ChunkLoc `json:"chunkLocations"`

	CreatedAt time.Time `json:"createdAt"`
	Hostname  string    `json:"hostname"`
	Username  string    `json:"username"`
}
