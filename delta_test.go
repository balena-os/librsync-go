package librsync

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
	"time"
)

func TestDelta(t *testing.T) {
	var totalBytes int64 = 10000000 // 1 MB

	var srcBuf bytes.Buffer
	src := io.TeeReader(
		io.LimitReader(rand.New(rand.NewSource(time.Now().UnixNano())), totalBytes),
		&srcBuf)

	s := signature(t, src)

	var buf bytes.Buffer

	// create 10% of difference by appending new random data
	newBytes := totalBytes / 10
	srcBuf.Truncate(int(totalBytes - newBytes))
	_, err := io.CopyN(&srcBuf, rand.New(rand.NewSource(time.Now().UnixNano())), newBytes)
	if err != nil {
		t.Error(err)
	}

	if err := Delta(s, &srcBuf, &buf); err != nil {
		t.Error(err)
	}
}
