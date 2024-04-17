package checksum

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

type MD5CheckSum struct {
}

func (c MD5CheckSum) GetCheckSum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
