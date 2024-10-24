package internal

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const alphabets = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func RandomID() string {
	return gonanoid.MustGenerate(alphabets, 24)
}
