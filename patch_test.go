package librsync

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var patchTestCases = []string{
	"000-blake2-11-23",
	"000-blake2-512-32",
	"000-md4-256-7",
	"001-blake2-512-32",
	"001-blake2-776-31",
	"001-md4-777-15",
	"002-blake2-512-32",
	"002-blake2-431-19",
	"002-md4-128-16",
	"003-blake2-512-32",
	"003-blake2-1024-13",
	"003-md4-1024-13",
	"004-blake2-1024-28",
	"004-blake2-2222-31",
	"004-blake2-512-32",
	"005-blake2-512-32",
	// "005-blake2-1000-18",
	// "005-md4-999-14",
	"006-blake2-2-32",
	// "007-blake2-5-32",
	"007-blake2-4-32",
	"007-blake2-3-32",
	// "008-blake2-222-30",
	// "008-blake2-512-32",
	// "008-md4-111-11",
	"009-blake2-2048-26",
	"009-blake2-512-32",
	"009-md4-2033-15",
	// "010-blake2-512-32",
	// "010-blake2-7-6",
	// "010-md4-4096-8",
}

// TestPatch tests patching only, using a reference delta.
func TestPatch(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	for _, tt := range patchTestCases {
		t.Run(tt, func(t *testing.T) {
			file, _, _, _, err := argsFromTestName(tt)
			r.NoError(err)

			baseFile, err := os.Open("testdata/" + file + ".old")
			r.NoError(err)

			deltaFile, err := os.Open("testdata/" + tt + ".delta")
			r.NoError(err)

			output := &bytes.Buffer{}
			err = Patch(baseFile, deltaFile, output)
			r.NoError(err)

			wantNewFile, err := ioutil.ReadFile("testdata/" + file + ".new")
			r.NoError(err)

			gotNewFile, err := ioutil.ReadAll(output)
			r.NoError(err)

			a.Equal(wantNewFile, gotNewFile)
		})
	}
}

// TestDeltaAndPatch tests both delta and patching, generating a delta (using a
// reference signature) and then applying this delta.
func TestDeltaAndPatch(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	for _, tt := range patchTestCases {
		t.Run(tt, func(t *testing.T) {
			file, _, _, _, err := argsFromTestName(tt)
			r.NoError(err)

			// Generate delta
			sig, err := readSignatureFile("testdata/" + tt + ".signature")
			r.NoError(err)

			newFile, err := os.Open("testdata/" + file + ".new")
			r.NoError(err)
			deltaBuffer := &bytes.Buffer{}

			err = Delta(sig, newFile, deltaBuffer)
			r.NoError(err)

			// Apply delta
			baseFile, err := os.Open("testdata/" + file + ".old")
			r.NoError(err)

			output := &bytes.Buffer{}
			err = Patch(baseFile, deltaBuffer, output)
			r.NoError(err)

			// Compare
			wantNewFile, err := ioutil.ReadFile("testdata/" + file + ".new")
			r.NoError(err)

			gotNewFile, err := ioutil.ReadAll(output)
			r.NoError(err)

			a.Equal(wantNewFile, gotNewFile)
		})
	}
}
