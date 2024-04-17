package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rpanchyk/fsync/internal/verify"
)

type Syncer struct {
	VerboseFlag bool
	Source      string
	Destination string
	Verifier    verify.Verifier

	// runtime
	absoluteSourcePath      string
	absoluteDestinationPath string
}

func (s *Syncer) Sync() error {
	src, dst, err := s.normalizePaths(s.Source, s.Destination)
	if err != nil {
		return err
	}
	s.absoluteSourcePath = src
	s.absoluteDestinationPath = dst

	if strings.HasPrefix(dst, src) {
		return errors.New(fmt.Sprint("Cannot synchronize because destination", dst, "is sub-folder of source", src))
	}

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return errors.New(fmt.Sprint("Source doesn't exist:", src))
	}

	return s.copy(src, dst)
}

func (s *Syncer) normalizePaths(src, dst string) (string, string, error) {
	srcPath, err := s.absolutePath(src)
	if err != nil {
		return "", "", err
	}

	dstPath, err := s.absolutePath(dst)
	if err != nil {
		return "", "", err
	}

	if s.VerboseFlag {
		fmt.Println("Normalized source:", srcPath)
		fmt.Println("Normalized destination:", dstPath)
	}

	return srcPath, dstPath, nil
}

func (s *Syncer) absolutePath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		currDir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		if s.VerboseFlag {
			fmt.Println("Current folder:", currDir)
		}

		return filepath.Join(currDir, path), nil
	}
	return path, nil
}

func (s *Syncer) copy(src, dst string) error {
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
		dirEntries, err := os.ReadDir(src)
		if err != nil {
			return errors.New(fmt.Sprint("Cannot get entries of source folder:", src))
		}

		for _, dirEntry := range dirEntries {
			entryInfo, err := dirEntry.Info()
			if err != nil {
				return errors.New(fmt.Sprint("Cannot get entry info:", dirEntry))
			}

			srcPath := filepath.Join(src, entryInfo.Name())

			dstPath := dst
			if entryInfo.IsDir() {
				dstPath = filepath.Join(dst, entryInfo.Name())
			}

			err = s.copy(srcPath, dstPath)
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
			path, _ := strings.CutPrefix(src, s.absoluteSourcePath)
			pathRunes := []rune(path)
			fmt.Println("Copied file", string(pathRunes[1:]), "-->", nBytes, "bytes")
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

	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	nBytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		return 0, err
	}

	// https://stackoverflow.com/questions/57850902/is-it-possible-to-calculate-md5-of-the-file-and-write-file-to-the-disk-in-the-sa
	ok, err := s.Verifier.Same(src, dst)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("source %s and destination %s files differ", src, dst)
	}

	return nBytes, err
}
