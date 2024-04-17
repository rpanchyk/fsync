package verify

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

type MD5Verifier struct {
}

func (v MD5Verifier) Same(file1, file2 string) (bool, error) {
	sum1, err := v.getSum(file1)
	if err != nil {
		return false, err
	}

	sum2, err := v.getSum(file2)
	if err != nil {
		return false, err
	}

	return sum1 == sum2, nil
}

func (v MD5Verifier) getSum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
