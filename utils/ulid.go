package utils

import (
	"crypto/rand"
	"io"
	"time"

	"github.com/oklog/ulid/v2"
)

var entropy io.Reader

func init() {
	entropy = ulid.DefaultEntropy()
}

// GenerateULID ULIDを生成する
func GenerateULID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

// GenerateULIDWithEntropy カスタムエントロピーを使用してULIDを生成する
func GenerateULIDWithEntropy() string {
	customEntropy := rand.Reader
	return ulid.MustNew(ulid.Timestamp(time.Now()), customEntropy).String()
}
