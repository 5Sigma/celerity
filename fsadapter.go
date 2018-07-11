package celerity

import "github.com/spf13/afero"

// FSAdapter is the adapter used for gaining access to the file system
type FSAdapter interface {
	RootPath(string) afero.Fs
}

// OSAdapter gives access to the file system
type OSAdapter struct{}

// RootPath reutrns a filesystem with OS access
func (o *OSAdapter) RootPath(path string) afero.Fs {
	return afero.NewBasePathFs(afero.NewOsFs(), path)
}

// MEMAdapter give access to an in memory file system for testing
type MEMAdapter struct {
	MEMFS afero.Fs
}

// NewMEMAdapter creates a new in memory FS adapter for testing
func NewMEMAdapter() *MEMAdapter {
	mm := afero.NewMemMapFs()
	return &MEMAdapter{mm}
}

// RootPath reutrns a filesystem with in memory access
func (m *MEMAdapter) RootPath(path string) afero.Fs {
	m.MEMFS.MkdirAll(path, 0755)
	return afero.NewBasePathFs(m.MEMFS, path)
}
