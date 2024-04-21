package bcrypt

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
	cost int
}

func New(cost int) Hasher {
	return Hasher{
		cost,
	}
}

func (h Hasher) Hash(input []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(input, h.cost)
}

func (h Hasher) HashString(input string) ([]byte, error) {
	return h.Hash([]byte(input))
}

func (h Hasher) Compare(hash, input []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, input) == nil
}

func (h Hasher) CompareString(hash []byte, input string) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(input)) == nil
}
