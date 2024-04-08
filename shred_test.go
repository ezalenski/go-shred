package shred

import (
	"errors"
	"os"
	"testing"
)

const test_data = "some data"

func TestShredInvalidPath(t *testing.T) {
	path := "badfilename"
	s := &Shred{
		os: &OsShim{},
	}
	err := s.Shred(path)
	if err == nil || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf(`Shred(%q) = %v, want %v`, path, err, os.ErrNotExist)
	}
}

func TestShredInvalidPerm(t *testing.T) {
	d, err := os.MkdirTemp(".", "test-*")
	if err != nil {
		t.Fatalf(`Failed to make tmp dir: %v`, err)
	}
	path := d + string(os.PathSeparator) + "badfileperm"
	_, err = os.OpenFile(path, os.O_CREATE|os.O_TRUNC, 0000)
	if err != nil {
		t.Fatalf(`Failed to make file: %q, %v`, path, err)
	}
	s := &Shred{
		os: &OsShim{},
	}
	err = s.Shred(path)
	if err == nil || !errors.Is(err, os.ErrPermission) {
		t.Fatalf(`Shred(%q) = %v, want %v`, path, err, os.ErrPermission)
	}
	err = os.RemoveAll(d)
	if err != nil {
		t.Fatalf(`Failed to remove %q: %v`, d, err)
	}
}

func TestShredDelete(t *testing.T) {
	path := "newfile"
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf(`Failed to open filpackage shred
		e %q: %v`, path, err)
	}
	i, err := f.WriteString(test_data)
	if err != nil || i != len(test_data) {
		t.Fatalf(`Failed to write file: %d, %v want %d, nil`, i, err, len(test_data))
	}
	s := &Shred{
		os: &OsShim{},
	}
	err = s.Shred(path)
	if err != nil {
		t.Fatalf(`Shred(%q) = %v, want nil`, path, err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf(`Failed to delete file`)
	}
}

func TestShredAmount(t *testing.T) {
	path := "newfile"
	writeAmount := int64(0)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf(`Failed to open file %q: %v`, path, err)
	}
	i, err := f.WriteString(test_data)
	if err != nil || i != len(test_data) {
		t.Fatalf(`Failed to write file: %d, %v want %d, nil`, i, err, len(test_data))
	}
	fShim := &FileShim{
		Delegate: f,
	}
	mockFile := &MockFile{
		WriteAtFunc: func(b []byte, off int64) (n int, err error) {
			writeAmount += int64(len(b))
			return fShim.WriteAt(b, off)
		},
	}
	osShim := OsShim{}
	mockOs := &MockOS{
		StatFunc: osShim.Stat,
		OpenFileFunc: func(name string, flag int, perm os.FileMode) (File, error) {
			_, err := osShim.OpenFile(name, flag, perm)
			if err != nil {
				t.Fatalf(`Failed to open file %q: %v`, name, err)
			}
			return mockFile, nil
		},
		RemoveFunc: osShim.Remove,
	}
	s := &Shred{
		os: mockOs,
	}
	s.Shred(path)
	if err != nil {
		t.Fatalf(`Shred(%q) = %v, want nil`, path, err)
	}
	if writeAmount != int64(len(test_data)*3) {
		t.Fatalf(`Bytes written to file: %d want %d`, writeAmount, len(test_data)*3)
	}
}
