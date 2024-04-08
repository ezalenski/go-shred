package shred

import "os"

type MockFileInfo struct {
	SizeFunc func() int64
}

func (fi *MockFileInfo) Size() int64 {
	return fi.SizeFunc()
}

type MockFile struct {
	WriteAtFunc func(b []byte, off int64) (n int, err error)
}

func (f *MockFile) WriteAt(b []byte, off int64) (n int, err error) {
	return f.WriteAtFunc(b, off)
}

type MockOS struct {
	StatFunc     func(name string) (FileInfo, error)
	OpenFileFunc func(name string, flag int, perm os.FileMode) (File, error)
	RemoveFunc   func(name string) error
}

func (os *MockOS) Stat(name string) (FileInfo, error) {
	return os.StatFunc(name)
}

func (os *MockOS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFileFunc(name, flag, perm)
}

func (os *MockOS) Remove(name string) error {
	return os.RemoveFunc(name)
}
