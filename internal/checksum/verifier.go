package checksum

type Verifier struct {
	checkSum CheckSum
}

type resultOrError struct {
	result string
	err    error
}

func NewVerifier(checkSumType CheckSumType) *Verifier {
	verifier := &Verifier{}

	switch checkSumType {
	case CRC32:
		verifier.checkSum = &CRC32CheckSum{}
	case MD5:
		verifier.checkSum = &MD5CheckSum{}
	default:
		panic("Unknown verification checksum type: " + string(checkSumType))
	}

	return verifier
}

func (v Verifier) Same(file1, file2 string) (bool, error) {
	outChan := make(chan resultOrError, 2)
	defer close(outChan)

	go func() {
		v.getCheckSum(file1, outChan)
	}()
	go func() {
		v.getCheckSum(file2, outChan)
	}()

	first, last := <-outChan, <-outChan
	if first.err != nil {
		return false, first.err
	}
	if last.err != nil {
		return false, last.err
	}

	return first.result == last.result, nil
}

func (v Verifier) getCheckSum(filePath string, outChan chan resultOrError) {
	roe := resultOrError{}
	roe.result, roe.err = v.checkSum.GetCheckSum(filePath)
	outChan <- roe
}
