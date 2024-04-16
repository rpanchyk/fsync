package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Syncer struct {
	VerboseFlag bool
}

func (s *Syncer) Copy(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return errors.New(fmt.Sprint("Source doesn't exist:", src))
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		if err = os.MkdirAll(dst, os.ModeDir); err != nil {
			return errors.New(fmt.Sprint("Cannot create destination folder:", dst))
		} else {
			if s.VerboseFlag {
				fmt.Println("Created destination folder:", dst)
			}
		}
	}

	srcIsDir, err := s.isDirectory(src)
	if err != nil {
		return errors.New(fmt.Sprint("Cannot analyze source:", src))
	}

	if srcIsDir {
		entries, err := os.ReadDir(src)
		if err != nil {
			return errors.New(fmt.Sprint("Cannot get entries of source folder:", src))
		}

		for _, entry := range entries {
			entryInfo, err := entry.Info()
			if err != nil {
				return errors.New(fmt.Sprint("Cannot get entry info:", entry))
			}

			srcPath := filepath.Join(src, entryInfo.Name())

			dstPath := dst
			if entryInfo.IsDir() {
				dstPath = filepath.Join(dst, entryInfo.Name())
			}

			err = s.Copy(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	} else {
		dstPath := filepath.Join(dst, filepath.Base(src))
		nBytes, err := s.copyFile(src, dstPath)
		if err != nil {
			return errors.New(fmt.Sprint("Cannot copy file:", src, "to", dstPath))
		}
		if s.VerboseFlag {
			fmt.Println("Copied file", src, "of", nBytes, "bytes")
		}
	}

	return nil
}

func (s *Syncer) isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func (s *Syncer) copyFile(src, dst string) (int64, error) {
	sourceFileInfo, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileInfo.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
