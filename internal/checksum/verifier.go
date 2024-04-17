package checksum

import (
	"sync"
)

type Verifier struct {
	checkSum CheckSum
}

func NewVerifier(checkSumType CheckSumType) *Verifier {
	verifier := &Verifier{}
	switch checkSumType {
	case CRC32:
		verifier.checkSum = &CRC32CheckSum{}
	case MD5:
		fallthrough
	default:
		verifier.checkSum = &MD5CheckSum{}
	}
	return verifier
}

func (v Verifier) Same(file1, file2 string) (bool, error) {
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

func (v Verifier) getSum(filePath string, outChan chan string, errChan chan error) {
	if sum, err := v.checkSum.GetCheckSum(filePath); err != nil {
		errChan <- err
	} else {
		outChan <- sum
	}
}
