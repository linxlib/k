package utils

import (
	"errors"
	"io"
	"os"
)

var DefaultPermCopy = os.FileMode(0777)

func Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

// IsFile checks whether given <path> a file, which means it's not a directory.
// Note that it returns false if the <path> does not exist.
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func CopyFile(src, dst string) (err error) {
	if src == "" {
		return errors.New("source file cannot be empty")
	}
	if dst == "" {
		return errors.New("destination file cannot be empty")
	}
	// If src and dst are the same path, it does nothing.
	if src == dst {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer func() {
		if e := in.Close(); e != nil {
			err = e
		}
	}()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	_, err = io.Copy(out, in)
	if err != nil {
		return
	}
	err = out.Sync()
	if err != nil {
		return
	}
	err = os.Chmod(dst, DefaultPermCopy)
	if err != nil {
		return
	}
	return
}

func Mkdir(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}
