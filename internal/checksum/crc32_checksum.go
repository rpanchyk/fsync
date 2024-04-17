package checksum

import "errors"

type CRC32CheckSum struct {
}

func (c CRC32CheckSum) GetCheckSum(filePath string) (string, error) {
	return "", errors.New("not implemented")
}
