package librsync

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPatch tests patching only, using a reference delta.
func TestPatch(t *testing.T) {
	r := require.New(t)
	a := assert.New(t)

	for _, tt := range allTestCases {
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

			wantNewFile, err := os.ReadFile("testdata/" + file + ".new")
			r.NoError(err)

			gotNewFile, err := io.ReadAll(output)
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

	for _, tt := range allTestCases {
		t.Run(tt, func(t *testing.T) {
			file, _, _, _, err := argsFromTestName(tt)
			r.NoError(err)

			// Generate delta
			sig, err := ReadSignatureFile("testdata/" + tt + ".signature")
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
			wantNewFile, err := os.ReadFile("testdata/" + file + ".new")
			r.NoError(err)

			gotNewFile, err := io.ReadAll(output)
			r.NoError(err)

			a.Equal(wantNewFile, gotNewFile)
		})
	}
}
