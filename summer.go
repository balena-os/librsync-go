package librsync

import (
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/md4"
)

// A Summer computes the strong checksum of a block of data. It uses a
// pre-allocated buffer to avoid allocations during the computation.
type Summer struct {
	hasher    hash.Hash
	strongLen uint32
	buffer    []byte
}

// NewSummer returns a new Summer that will use a hashing algorithm specified
// by sigType. The strongLen parameter specifies the length of the strong
// checksum that will be returned by Sum().
func NewSummer(sigType MagicNumber, strongLen uint32) (*Summer, error) {
	result := &Summer{
		hasher:    nil,
		strongLen: strongLen,
		buffer:    nil,
	}

	switch sigType {
	case BLAKE2_SIG_MAGIC:
		var err error
		result.hasher, err = blake2b.New256(nil)
		if err != nil {
			return nil, fmt.Errorf("constructing Summer: %v", err)
		}
		result.buffer = make([]byte, 0, BLAKE2_SUM_LENGTH)

	case MD4_SIG_MAGIC:
		result.hasher = md4.New()
		result.buffer = make([]byte, 0, MD4_SUM_LENGTH)

	default:
		return nil, fmt.Errorf("constructing Summer: invalid sigType %#x", sigType)
	}

	return result, nil
}

// Sum returns the strong checksum data. The returned slice is a reference to
// the internal buffer, and is only valid until the next call to Sum().
func (s *Summer) Sum(data []byte) []byte {
	s.hasher.Reset()
	s.hasher.Write(data)
	return s.hasher.Sum(s.buffer[:0])[:s.strongLen]
}
