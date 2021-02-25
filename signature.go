package librsync

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/md4"
)

const (
	BLAKE2_SUM_LENGTH = 32
	MD4_SUM_LENGTH    = 16
)

type SignatureType struct {
	sigType    MagicNumber
	blockLen   uint32
	strongLen  uint32
	strongSigs [][]byte
	weak2block map[uint32]int
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

func Signature(input io.Reader, output io.Writer, blockLen, strongLen uint32, sigType MagicNumber) (*SignatureType, error) {
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
	ret.weak2block = make(map[uint32]int)
	ret.sigType = sigType
	ret.strongLen = strongLen
	ret.blockLen = blockLen

	for {
		n, err := io.ReadFull(input, block)
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			break
		} else if err != nil {
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

func (s *SignatureType) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, s.sigType)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.BigEndian, s.blockLen)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.BigEndian, s.strongLen)
	if err != nil {
		return nil, err
	}
	idx := 0
	for weak := range s.weak2block {
		err := binary.Write(&buf, binary.BigEndian, weak)
		if err != nil {
			return nil, err
		}
		buf.Write(s.strongSigs[idx])
		idx++
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (s *SignatureType) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.BigEndian, &s.sigType); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &s.blockLen); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &s.strongLen); err != nil {
		return err
	}
	idx := 0
	s.weak2block = make(map[uint32]int)
	for {
		var weak uint32
		err := binary.Read(buf, binary.BigEndian, &weak)
		if err != nil {
			// we're done here
			if err == io.EOF {
				break
			}
			return err
		}
		s.weak2block[weak] = idx
		idx++

		var strong []byte
		if err := binary.Read(buf, binary.BigEndian, &strong); err != nil {
			return err
		}
		s.strongSigs = append(s.strongSigs, strong)
	}
	return nil
}
