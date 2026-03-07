package dto

type SnapshotSelectItem struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

type SnapshotSelectResponse struct {
	Snapshots []SnapshotSelectItem `json:"snapshots"`
}
