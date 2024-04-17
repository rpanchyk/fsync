package checksum

type CheckSum interface {
	GetCheckSum(filePath string) (string, error)
}

type CheckSumType uint8

const (
	MD5 CheckSumType = iota
	CRC32
)
