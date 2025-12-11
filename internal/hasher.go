package internal

import (
	"crypto/sha1"
	"fmt"
)

type Hasher struct{}

func NewHasher() *Hasher { return &Hasher{} }

func (h *Hasher) Hash(content string) string {
	sum := sha1.Sum([]byte(content))
	return fmt.Sprintf("%x", sum)
}