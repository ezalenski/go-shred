package shred

import (
	"fmt"
	"math/rand"
	"os"
)

const numShreds = 3

type Shred struct {
	os OS
}

func (s *Shred) ShredAt(file File, off int64, n int64) error {
	b := make([]byte, n)
	rand.Read(b)
	bytesWritten, err := file.WriteAt(b, off)
	if err != nil {
		return err
	}
	if int64(bytesWritten) != n {
		return fmt.Errorf(`failed to overwrite byte at %x`, off)
	}
	return nil
}

func (s *Shred) Shred(file_path string) error {
	stat, err := s.os.Stat(file_path)
	if err != nil {
		return err
	}
	file, err := s.os.OpenFile(file_path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	for i := 0; i < numShreds; i++ {
		if err = s.ShredAt(file, 0, stat.Size()); err != nil {
			return err
		}
	}
	if err := s.os.Remove(file_path); err != nil {
		return err
	}
	return nil
}
