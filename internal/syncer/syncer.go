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
	DeleteFlag       bool
	Source           string
	Destination      string
	ChecksumVerifier *checksum.Verifier

	// runtime
	absoluteSourcePath      string
	absoluteDestinationPath string
}

const TEMP_FILE_EXT = ".fsync-tmp"

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
			return fmt.Errorf("cannot create destination folder %s error: %s", dst, err.Error())
		} else {
			if s.VerboseFlag {
				fmt.Println("Created destination folder:", dst)
			}
		}
	}

	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("cannot analyze source %s error: %s", src, err.Error())
	}
	if srcFileInfo.IsDir() {
		srcDirEntries, err := os.ReadDir(src)
		if err != nil {
			return fmt.Errorf("cannot get entries of source folder %s error: %s", src, err.Error())
		}

		srcEntries := make(map[string]struct{})

		for _, dirEntry := range srcDirEntries {
			entryInfo, err := dirEntry.Info()
			if err != nil {
				return fmt.Errorf("cannot get entry info %s error: %s", dirEntry, err.Error())
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

			if s.DeleteFlag {
				relativePath := s.relativePath(s.absoluteSourcePath, srcPath)
				srcEntries[relativePath] = struct{}{}
			}
		}

		if s.DeleteFlag {
			dstDirEntries, err := os.ReadDir(dst)
			if err != nil {
				return fmt.Errorf("cannot get entries of destination folder %s error: %s", dst, err.Error())
			}

			for _, dirEntry := range dstDirEntries {
				entryInfo, err := dirEntry.Info()
				if err != nil {
					return fmt.Errorf("cannot get entry info %s error: %s", dirEntry, err.Error())
				}

				dstPath := filepath.Join(dst, entryInfo.Name())
				relativePath := s.relativePath(s.absoluteDestinationPath, dstPath)
				if _, ok := srcEntries[relativePath]; !ok {
					if err = s.removeExtraneous(dstPath, dirEntry.IsDir()); err != nil {
						return fmt.Errorf("cannot remove extraneous %s error: %s", dstPath, err.Error())
					}
					if s.VerboseFlag {
						fmt.Println("Removed extraneous", relativePath)
					}
				}
			}
		}
	} else {
		dstPath := filepath.Join(dst, filepath.Base(src))

		if dstFileInfo, err := os.Stat(dstPath); err == nil && !dstFileInfo.IsDir() {
			ok, err := s.ChecksumVerifier.Same(src, dstPath)
			if err != nil {
				return fmt.Errorf("cannot get checksum for file %s error: %s", dstPath, err.Error())
			}
			if ok {
				fmt.Println("File", s.relativePath(s.absoluteDestinationPath, dstPath), "is up do date")
				return nil
			}
		}

		tempDstPath := dstPath + TEMP_FILE_EXT
		nBytes, err := s.copyFile(src, tempDstPath)
		if err != nil {
			return fmt.Errorf("cannot copy file %s to %s error: %s", src, tempDstPath, err.Error())
		}
		if err = os.Rename(tempDstPath, dstPath); err != nil {
			return fmt.Errorf("cannot rename file %s to %s error: %s", tempDstPath, dstPath, err.Error())
		}
		if s.VerboseFlag {
			fmt.Println("Copied file", s.relativePath(s.absoluteDestinationPath, dstPath), "-->", nBytes, "bytes")
		}
	}

	return nil
}

func (s *Syncer) copyFile(src, dst string) (int64, error) {
	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !srcFileInfo.Mode().IsRegular() {
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

func (s *Syncer) relativePath(parentPath, childPath string) string {
	if path, ok := strings.CutPrefix(childPath, parentPath); !ok {
		return childPath
	} else {
		pathRunes := []rune(path)
		return string(pathRunes[1:])
	}
}

func (s *Syncer) removeExtraneous(path string, isDir bool) error {
	if isDir {
		return os.RemoveAll(path)
	} else {
		return os.Remove(path)
	}
}
