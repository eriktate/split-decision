package api

type Hasher interface {
	Hash(input []byte) ([]byte, error)
	HashString(input string) ([]byte, error)
	Compare(hashed, input []byte) bool
	CompareString(hashed []byte, input string) bool
}
