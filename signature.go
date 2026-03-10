package librsync

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/md4"
)

const (
	BLAKE2_SUM_LENGTH            = 32
	MD4_SUM_LENGTH               = 16
	INPUT_SIZE_FOR_MIN_BLOCK_LEN = 65536
	DEFAULT_BLOCK_LEN            = 2048
	MIN_BLOCK_LEN                = 256
)

type SignatureType struct {
	sigType    MagicNumber
	blockLen   uint32
	strongLen  uint32
	strongSigs [][]byte
	weak2block map[uint32]int
}

// Based on the original librsync.
// https://github.com/librsync/librsync/blob/271744de428c0177e993f95c045d70b8cf2558f2/src/sumset.c#L141
// Attempts to give a block len that is a reasonable compromise between signature size, delta size,
// and performance.
func recommendBlockLen(size int64) uint32 {
	if size <= 0 {
		// don't know what the input size is, use a balanced block len
		return DEFAULT_BLOCK_LEN
	} else if size <= INPUT_SIZE_FOR_MIN_BLOCK_LEN {
		return MIN_BLOCK_LEN
	} else {
		sizeSqrt := math.Floor(math.Sqrt(float64(size)))
		// round up the result to be a multiple of 128 to match blake2b internal buffer
		size128 := uint32(sizeSqrt) & ^uint32(127)

		return size128
	}
}

func CalcStrongSum(data []byte, sigType MagicNumber, strongLen uint32) ([]byte, error) {
	switch sigType {
	case BLAKE2_SIG_MAGIC:
		d := blake2b.Sum256(data)
		return d[:strongLen], nil
	case MD4_SIG_MAGIC:
		d := md4.New()
		d.Write(data)
		return d.Sum(nil)[:strongLen], nil
	}
	return nil, fmt.Errorf("Invalid sigType %#x", sigType)
}

func Signature(input io.Reader, output io.Writer, blockLen, strongLen uint32, sigType MagicNumber, inputSize int64) (*SignatureType, error) {
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

	if blockLen == 0 {
		blockLen = recommendBlockLen(inputSize)
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

	nbBlocks := int64(0)
	if inputSize > 0 {
		nbBlocks = inputSize / int64(blockLen)
		if inputSize%int64(blockLen) != 0 {
			nbBlocks += 1
		}
	}

	block := make([]byte, blockLen)

	var ret SignatureType
	ret.weak2block = make(map[uint32]int, nbBlocks)
	ret.strongSigs = make([][]byte, 0, nbBlocks)
	ret.sigType = sigType
	ret.strongLen = strongLen
	ret.blockLen = blockLen

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

		strong, _ := CalcStrongSum(data, sigType, strongLen)
		output.Write(strong)

		ret.weak2block[weak] = len(ret.strongSigs)
		ret.strongSigs = append(ret.strongSigs, strong)
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

	strongSigs := [][]byte{}
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

		weak2block[weakSum] = len(strongSigs)
		strongSigs = append(strongSigs, strongSum)
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
