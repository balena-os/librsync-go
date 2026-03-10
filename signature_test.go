package librsync

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errorI interface {
	// testing#T.Error || testing#B.Error
	Error(args ...interface{})
}

func signature(t errorI, src io.Reader, inputSize int64) *SignatureType {
	var (
		magic            = BLAKE2_SIG_MAGIC
		blockLen  uint32 = 512
		strongLen uint32 = 32
		bufSize          = 65536
	)

	s, err := Signature(
		bufio.NewReaderSize(src, bufSize),
		ioutil.Discard,
		blockLen, strongLen, magic, inputSize)
	if err != nil {
		t.Error(err)
	}

	return s
}

func TestSignature(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	for _, tt := range allTestCases {
		t.Run(tt, func(t *testing.T) {
			file, magic, blockLen, strongLen, err := argsFromTestName(tt)
			r.NoError(err)

			inputData, err := ioutil.ReadFile("testdata/" + file + ".old")
			r.NoError(err)
			input := bytes.NewReader(inputData)

			output := &bytes.Buffer{}
			gotSig, err := Signature(input, output, blockLen, strongLen, magic, int64(len(inputData)))
			r.NoError(err)

			wantSig, err := ReadSignatureFile("testdata/" + tt + ".signature")
			r.NoError(err)
			a.Equal(wantSig.blockLen, gotSig.blockLen)
			a.Equal(wantSig.sigType, gotSig.sigType)
			a.Equal(wantSig.strongLen, gotSig.strongLen)

			outputData, err := ioutil.ReadAll(output)
			r.NoError(err)
			expectedData, err := ioutil.ReadFile("testdata/" + tt + ".signature")
			r.NoError(err)
			a.Equal(expectedData, outputData)
		})
	}
}

func TestRecommendBlockLen(t *testing.T) {
	a := assert.New(t)

	tt := []struct {
		InputSize int64
		BlockLen  uint32
	}{
		{0, DEFAULT_BLOCK_LEN},
		{1, MIN_BLOCK_LEN},
		{1000000, 896},
	}

	for _, test := range tt {
		blockLen := recommendBlockLen(test.InputSize)

		a.Equal(test.BlockLen, blockLen)
	}
}
