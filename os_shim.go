package shred

import (
	"io/fs"
	"os"
)

type FileShim struct {
	Delegate *os.File
}

type File interface {
	WriteAt(b []byte, off int64) (n int, err error)
}

func (f *FileShim) WriteAt(b []byte, off int64) (n int, err error) {
	return f.Delegate.WriteAt(b, off)
}

type FileInfoShim struct {
	Delegate fs.FileInfo
}

func (f *FileInfoShim) Size() (n int64) {
	return f.Delegate.Size()
}

type FileInfo interface {
	Size() int64
}

type OsShim struct{}

type OS interface {
	Stat(name string) (FileInfo, error)
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	Remove(name string) error
}

func (sh *OsShim) Stat(name string) (FileInfo, error) {
	o, err := os.Stat(name)
	return &FileInfoShim{Delegate: o}, err
}

func (sh *OsShim) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	o, err := os.OpenFile(name, flag, perm)
	return &FileShim{Delegate: o}, err
}

func (sh *OsShim) Remove(name string) error {
	return os.Remove(name)
}
