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

var signatureTestCases = []string{
	// "000-blake2-11-23",
	// "000-blake2-512-32",
	// "000-md4-256-7",
	// "001-blake2-512-32",
	// "001-blake2-776-31",
	// "001-md4-777-15",
	// "002-blake2-512-32",
	// "002-blake2-431-19",
	// "002-md4-128-16",
	"003-blake2-512-32",
	"003-blake2-1024-13",
	"003-md4-1024-13",
	// "004-blake2-1024-28",
	// "004-blake2-2222-31",
	// "004-blake2-512-32",
	// "005-blake2-512-32",
	// "005-blake2-1000-18",
	// "005-md4-999-14",
	// "006-blake2-2-32",
	"007-blake2-5-32",
	// "007-blake2-4-32",
	// "007-blake2-3-32",
	// "008-blake2-222-30",
	// "008-blake2-512-32",
	// "008-md4-111-11",
	"009-blake2-2048-26",
	"009-blake2-512-32",
	"009-md4-2033-15",
	"010-blake2-512-32",
	"010-blake2-7-6",
	"010-md4-4096-8",
}

func TestSignature(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	for _, tt := range signatureTestCases {
		t.Run(tt, func(t *testing.T) {
			file, magic, blockLen, strongLen, err := argsFromTestName(tt)
			r.NoError(err)

			inputData, err := ioutil.ReadFile("testdata/" + file + ".old")
			r.NoError(err)
			input := bytes.NewReader(inputData)

			output := &bytes.Buffer{}
			gotSig, err := Signature(input, output, blockLen, strongLen, magic)
			r.NoError(err)

			wantSig, err := readSignatureFile("testdata/" + tt + ".signature")
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
