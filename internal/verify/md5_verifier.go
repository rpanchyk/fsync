package verify

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"sync"
)

type MD5Verifier struct {
}

func (v MD5Verifier) Same(file1, file2 string) (bool, error) {
	var wg sync.WaitGroup
	outChan := make(chan string, 2)
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		v.getSum(file1, outChan, errChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		v.getSum(file2, outChan, errChan)
	}()

	go func() {
		wg.Wait()
		close(outChan)
		close(errChan)
	}()

	for err := range errChan {
		return false, err
	}

	sumA := <-outChan
	sumB := <-outChan
	return sumA == sumB, nil
}

func (v MD5Verifier) getSum(filePath string, outChan chan string, errChan chan error) {
	file, err := os.Open(filePath)
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		errChan <- err
		return
	}
	outChan <- fmt.Sprintf("%x", hash.Sum(nil))
}
