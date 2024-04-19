package storage

type StorageFeature struct {
	Feature   string `db:"feature"`
	IsEnabled bool   `db:"is_enabled"`
}
