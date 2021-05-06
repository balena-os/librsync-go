package librsync

import (
	"bufio"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

type errorI interface {
	// testing#T.Error || testing#B.Error
	Error(args ...interface{})
}

func signature(t errorI, src io.Reader) *SignatureType {
	var (
		magic            = BLAKE2_SIG_MAGIC
		blockLen  uint32 = 512
		strongLen uint32 = 32
		bufSize          = 65536
	)

	s, err := Signature(
		bufio.NewReaderSize(src, bufSize),
		ioutil.Discard,
		blockLen, strongLen, magic)
	if err != nil {
		t.Error(err)
	}

	return s
}

func TestSignature(t *testing.T) {
	var totalBytes int64 = 10000000 // 1 MB
	src := io.LimitReader(rand.New(rand.NewSource(time.Now().UnixNano())), totalBytes)
	signature(t, src)
}
