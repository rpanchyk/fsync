package verify

type Verifier interface {
	Same(file1, file2 string) (bool, error)
}
