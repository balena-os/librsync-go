package librsync

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		io.Discard,
		blockLen, strongLen, magic)
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

			inputData, err := os.ReadFile("testdata/" + file + ".old")
			r.NoError(err)
			input := bytes.NewReader(inputData)

			output := &bytes.Buffer{}
			gotSig, err := Signature(input, output, blockLen, strongLen, magic)
			r.NoError(err)

			wantSig, err := ReadSignatureFile("testdata/" + tt + ".signature")
			r.NoError(err)
			a.Equal(wantSig.blockLen, gotSig.blockLen)
			a.Equal(wantSig.sigType, gotSig.sigType)
			a.Equal(wantSig.strongLen, gotSig.strongLen)

			outputData, err := io.ReadAll(output)
			r.NoError(err)
			expectedData, err := os.ReadFile("testdata/" + tt + ".signature")
			r.NoError(err)
			a.Equal(expectedData, outputData)
		})
	}
}
