package syncer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rpanchyk/fsync/internal/checksum"
)

type Syncer struct {
	VerboseFlag      bool
	Source           string
	Destination      string
	ChecksumVerifier *checksum.Verifier

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
		return fmt.Errorf("cannot synchronize because destination %s is sub-folder of source %s", dst, src)
	}

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("source doesn't exist: %s", src)
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
			return fmt.Errorf("cannot create destination folder: %s", dst)
		} else {
			if s.VerboseFlag {
				fmt.Println("Created destination folder:", dst)
			}
		}
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("cannot analyze source %s error: %s", src, err.Error())
	}
	if srcInfo.IsDir() {
		dirEntries, err := os.ReadDir(src)
		if err != nil {
			return fmt.Errorf("cannot get entries of source folder: %s", src)
		}

		for _, dirEntry := range dirEntries {
			entryInfo, err := dirEntry.Info()
			if err != nil {
				return fmt.Errorf("cannot get entry info: %s", dirEntry)
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
			return fmt.Errorf("cannot copy file %s to %s", src, dstPath)
		}
		if s.VerboseFlag {
			path, _ := strings.CutPrefix(src, s.absoluteSourcePath)
			pathRunes := []rune(path)
			fmt.Println("Copied file", string(pathRunes[1:]), "-->", nBytes, "bytes")
		}
	}

	return nil
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

	ok, err := s.ChecksumVerifier.Same(src, dst)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("source %s and destination %s files differ", src, dst)
	}

	return nBytes, err
}
