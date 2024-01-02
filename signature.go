package librsync

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	BLAKE2_SUM_LENGTH = 32
	MD4_SUM_LENGTH    = 16
)

type SignatureType struct {
	sigType   MagicNumber
	blockLen  uint32
	strongLen uint32

	// strongSigs is interpreted as a slice of strongLen-byte values. Previously
	// it was declared as a [][]byte, which was more clear, but didn't allow us
	// to pre-allocate the whole memory block at once and avoid allocations
	// during the hot loop of the signature calculation.
	strongSigs []byte
	weak2block map[uint32]int
}

func Signature(input io.Reader, output io.Writer, blockLen, strongLen uint32, sigType MagicNumber) (*SignatureType, error) {
	return SignatureWithBlockCount(input, output, blockLen, strongLen, sigType, 0)
}

// SignatureWithBlockCount is a version of Signature that allows the caller to
// pass in the expected number of blocks in the Signature. This is used to
// pre-allocate the internal data structures.
func SignatureWithBlockCount(input io.Reader, output io.Writer, blockLen, strongLen uint32, sigType MagicNumber, blockCount int) (*SignatureType, error) {
	var maxStrongLen uint32

	switch sigType {
	case BLAKE2_SIG_MAGIC:
		maxStrongLen = BLAKE2_SUM_LENGTH
	case MD4_SIG_MAGIC:
		maxStrongLen = MD4_SUM_LENGTH
	default:
		return nil, fmt.Errorf("invalid sigType %#x", sigType)
	}

	if strongLen > maxStrongLen {
		return nil, fmt.Errorf("invalid strongLen %d for sigType %#x", strongLen, sigType)
	}

	err := binary.Write(output, binary.BigEndian, sigType)
	if err != nil {
		return nil, err
	}
	err = binary.Write(output, binary.BigEndian, blockLen)
	if err != nil {
		return nil, err
	}
	err = binary.Write(output, binary.BigEndian, strongLen)
	if err != nil {
		return nil, err
	}

	block := make([]byte, blockLen)

	var ret SignatureType
	ret.weak2block = make(map[uint32]int, blockCount)
	ret.strongSigs = make([]byte, 0, blockCount*int(strongLen))
	ret.sigType = sigType
	ret.strongLen = strongLen
	ret.blockLen = blockLen

	summer, err := NewSummer(sigType, strongLen)
	if err != nil {
		return nil, err
	}

	for {
		n, err := io.ReadAtLeast(input, block, int(blockLen))
		if err == io.EOF {
			// We reached the end of the input, we are done with the signature
			break
		} else if err == nil || err == io.ErrUnexpectedEOF {
			if n == 0 {
				// No real error and no new data either: that also signals the
				// end the input; we are done with the signature
				break
			}
			// No real error, got data. Leave this `if` and checksum this block
		} else if err != nil {
			// Got a real error, report it back to the caller
			return nil, err
		}

		data := block[:n]

		weak := WeakChecksum(data)
		err = binary.Write(output, binary.BigEndian, weak)
		if err != nil {
			return nil, err
		}

		strong := summer.Sum(data)
		output.Write(strong)

		ret.weak2block[weak] = len(ret.strongSigs) / int(strongLen)
		ret.strongSigs = append(ret.strongSigs, strong...)
	}

	return &ret, nil
}

// ReadSignature reads a signature from an io.Reader.
func ReadSignature(r io.Reader) (*SignatureType, error) {
	var magic MagicNumber
	err := binary.Read(r, binary.BigEndian, &magic)
	if err != nil {
		return nil, err
	}

	var blockLen uint32
	err = binary.Read(r, binary.BigEndian, &blockLen)
	if err != nil {
		return nil, err
	}

	var strongLen uint32
	err = binary.Read(r, binary.BigEndian, &strongLen)
	if err != nil {
		return nil, err
	}

	strongSigs := []byte{}
	weak2block := map[uint32]int{}

	for {
		var weakSum uint32
		err = binary.Read(r, binary.BigEndian, &weakSum)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		strongSum := make([]byte, strongLen)
		_, err := io.ReadFull(r, strongSum)
		if err != nil {
			return nil, err
		}

		weak2block[weakSum] = len(strongSigs) / int(strongLen)
		strongSigs = append(strongSigs, strongSum...)
	}

	return &SignatureType{
		sigType:    magic,
		blockLen:   blockLen,
		strongLen:  strongLen,
		strongSigs: strongSigs,
		weak2block: weak2block,
	}, nil
}

// ReadSignatureFile reads a signature from the file at path.
func ReadSignatureFile(path string) (*SignatureType, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadSignature(f)
}
