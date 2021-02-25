package librsync

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"
)

func signature(b *testing.B, src io.Reader) *SignatureType {
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
		b.Error(err)
	}

	return s
}

func BenchmarkSignature(b *testing.B) {
	var totalBytes int64 = 1000000000 // 1 GB
	src := io.LimitReader(rand.New(rand.NewSource(time.Now().UnixNano())), totalBytes)

	b.SetBytes(totalBytes)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		signature(b, src)
	}
}

func BenchmarkDelta(b *testing.B) {
	var totalBytes int64 = 1000000000 // 1 GB

	var srcBuf bytes.Buffer
	src := io.TeeReader(
		io.LimitReader(rand.New(rand.NewSource(time.Now().UnixNano())), totalBytes),
		&srcBuf)
	b.SetBytes(totalBytes)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := signature(b, src)

		var buf bytes.Buffer

		// create 10% of difference by appending new random data
		newBytes := totalBytes / 10
		srcBuf.Truncate(int(totalBytes - newBytes))
		_, err := io.CopyN(&srcBuf, rand.New(rand.NewSource(time.Now().UnixNano())), newBytes)
		if err != nil {
			b.Error(err)
		}

		if err := Delta(s, &srcBuf, &buf); err != nil {
			b.Error(err)
		}

		b.Logf("raw   size:    %v bytes", totalBytes)
		b.Logf("delta size:    %v bytes (%.2f%%)", len(buf.Bytes()), (float64(len(buf.Bytes()))/float64(totalBytes))*100)
	}
}

func BenchmarkDeltaWithCache(b *testing.B) {
	var totalBytes int64 = 1000000000 // 1 GB

	var srcBuf bytes.Buffer
	src := io.TeeReader(
		io.LimitReader(rand.New(rand.NewSource(time.Now().UnixNano())), totalBytes),
		&srcBuf)
	s := signature(b, src)

	b.SetBytes(totalBytes)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer

		// create 10% of difference by appending new random data
		newBytes := totalBytes / 10
		srcBuf.Truncate(int(totalBytes - newBytes))
		_, err := io.CopyN(&srcBuf, rand.New(rand.NewSource(time.Now().UnixNano())), newBytes)
		if err != nil {
			b.Error(err)
		}

		if err := Delta(s, &srcBuf, &buf); err != nil {
			b.Error(err)
		}

		b.Logf("raw   size:    %v bytes", totalBytes)
		b.Logf("delta size:    %v bytes (%.2f%%)", len(buf.Bytes()), (float64(len(buf.Bytes()))/float64(totalBytes))*100)
	}
}
