package librsync

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// TestSignatureBinaryMarshal verifies that the marshalled `StructureType` raw
// bytes match the that is written to `output`.
func TestSignatureBinaryMarshal(t *testing.T) {
	var (
		totalBytes int64  = 10000000 // 10 MB
		magic             = BLAKE2_SIG_MAGIC
		blockLen   uint32 = 512
		strongLen  uint32 = 32
		bufSize           = 65536
	)

	for i := 1; i <= 5; i++ {
		t.Run(fmt.Sprintf("%s-%d", t.Name(), i), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			src := io.LimitReader(rand.New(rand.NewSource(time.Now().UnixNano())), totalBytes)

			s, err := Signature(
				bufio.NewReaderSize(src, bufSize),
				&buf,
				blockLen, strongLen, magic)
			if err != nil {
				t.Fatal(err)
			}

			data, err := s.MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}

			if reflect.DeepEqual(buf.Bytes(), data) {
				t.Fatal(err)
			}
		})
	}
}

// TestSignatureBinaryUnmarshal verifies that the marshalled `StructureType` raw
// bytes can be unmarshalled correctly
func TestSignatureBinaryUnmarshal(t *testing.T) {
	var (
		totalBytes int64  = 10000000 // 10 MB
		magic             = BLAKE2_SIG_MAGIC
		blockLen   uint32 = 512
		strongLen  uint32 = 32
		bufSize           = 65536
	)

	for i := 1; i <= 5; i++ {
		t.Run(fmt.Sprintf("%s-%d", t.Name(), i), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			src := io.LimitReader(rand.New(rand.NewSource(time.Now().UnixNano())), totalBytes)

			s, err := Signature(
				bufio.NewReaderSize(src, bufSize),
				&buf,
				blockLen, strongLen, magic)
			if err != nil {
				t.Fatal(err)
			}

			var data SignatureType
			if err := data.UnmarshalBinary(buf.Bytes()); err != nil {
				t.Fatal(err)
			}

			if reflect.DeepEqual(s, data) {
				t.Fatal(err)
			}
		})
	}
}
