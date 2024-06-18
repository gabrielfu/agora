package internal

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func RandomID() string {
	return gonanoid.Must(16)
}
