package librsync

import (
	"bytes"
	"io"
	"math/rand"
	"testing"
	"time"
)

// Benchmarks generating a signature for a file totalBytes long.
func benchmarkSignature(b *testing.B, totalBytes int64) {
	b.SetBytes(totalBytes)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		src := io.LimitReader(rand.New(rand.NewSource(time.Now().UnixNano())), totalBytes)
		signature(b, src, int(totalBytes))
	}
}

func BenchmarkSignature1MB(b *testing.B) {
	benchmarkSignature(b, 1_000_000)
}

func BenchmarkSignature1GB(b *testing.B) {
	benchmarkSignature(b, 1_000_000_000)
}

// Changes the final 10% of the input data.
func benchmarkDeltaChangeTail(b *testing.B, totalBytes int64) {
	newBytes := totalBytes / 10
	oldBytes := totalBytes - newBytes
	oldSeed := time.Now().UnixNano()
	oldData := io.LimitReader(rand.New(rand.NewSource(oldSeed)), totalBytes)
	s := signature(b, oldData, int(totalBytes))

	b.SetBytes(totalBytes)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		newSeed := time.Now().UnixNano()
		newData := io.MultiReader(
			io.LimitReader(rand.New(rand.NewSource(oldSeed)), oldBytes),
			io.LimitReader(rand.New(rand.NewSource(newSeed)), newBytes),
		)

		var buf bytes.Buffer
		if err := Delta(s, newData, &buf); err != nil {
			b.Error(err)
		}

		b.Logf("raw   size:    %v bytes", totalBytes)
		b.Logf("delta size:    %v bytes (%.2f%%)", len(buf.Bytes()), (float64(len(buf.Bytes()))/float64(totalBytes))*100)
	}
}

func BenchmarkDeltaChangeTail1GB(b *testing.B) {
	benchmarkDeltaChangeTail(b, 1_000_000_000)
}

func BenchmarkDeltaChangeTail1MB(b *testing.B) {
	benchmarkDeltaChangeTail(b, 1_000_000)
}

//
// Uninteresting variations of the Delta benchmark above. They all seem to give
// very similar results.
//

// Adds 10% of new data to the end of the input data.
func benchmarkDeltaAppend(b *testing.B, totalBytes int64) {
	newBytes := totalBytes / 10
	oldSeed := time.Now().UnixNano()
	oldData := io.LimitReader(rand.New(rand.NewSource(oldSeed)), totalBytes)
	s := signature(b, oldData, int(totalBytes))

	b.SetBytes(totalBytes)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		newSeed := time.Now().UnixNano()
		newData := io.MultiReader(
			io.LimitReader(rand.New(rand.NewSource(oldSeed)), totalBytes),
			io.LimitReader(rand.New(rand.NewSource(newSeed)), newBytes),
		)

		var buf bytes.Buffer
		if err := Delta(s, newData, &buf); err != nil {
			b.Error(err)
		}

		b.Logf("raw   size:    %v bytes", totalBytes)
		b.Logf("delta size:    %v bytes (%.2f%%)", len(buf.Bytes()), (float64(len(buf.Bytes()))/float64(totalBytes))*100)
	}
}

func BenchmarkDeltaAppend1GB(b *testing.B) {
	benchmarkDeltaAppend(b, 1_000_000_000)
}

func BenchmarkDeltaAppend1MB(b *testing.B) {
	benchmarkDeltaAppend(b, 1_000_000)
}

// Adds 10% of new data to the beginning of the input data.
func benchmarkDeltaPrepend(b *testing.B, totalBytes int64) {
	newBytes := totalBytes / 10
	oldSeed := time.Now().UnixNano()
	oldData := io.LimitReader(rand.New(rand.NewSource(oldSeed)), totalBytes)
	s := signature(b, oldData, int(totalBytes))

	b.SetBytes(totalBytes)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		newSeed := time.Now().UnixNano()
		newData := io.MultiReader(
			io.LimitReader(rand.New(rand.NewSource(newSeed)), newBytes),
			io.LimitReader(rand.New(rand.NewSource(oldSeed)), totalBytes),
		)

		var buf bytes.Buffer
		if err := Delta(s, newData, &buf); err != nil {
			b.Error(err)
		}

		b.Logf("raw   size:    %v bytes", totalBytes)
		b.Logf("delta size:    %v bytes (%.2f%%)", len(buf.Bytes()), (float64(len(buf.Bytes()))/float64(totalBytes))*100)
	}
}

func BenchmarkDeltaPrepend1GB(b *testing.B) {
	benchmarkDeltaPrepend(b, 1_000_000_000)
}

func BenchmarkDeltaPrepend1MB(b *testing.B) {
	benchmarkDeltaPrepend(b, 1_000_000)
}

// Adds 10% of new data to the middle of the input data.
func benchmarkDeltaInpend(b *testing.B, totalBytes int64) {
	newBytes := totalBytes / 10
	oldSeed := time.Now().UnixNano()
	oldData := io.LimitReader(rand.New(rand.NewSource(oldSeed)), totalBytes)

	firstBytes := totalBytes / 3
	lastBytes := totalBytes - firstBytes

	s := signature(b, oldData, int(totalBytes))

	b.SetBytes(totalBytes)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		newSeed := time.Now().UnixNano()
		oldData = io.LimitReader(rand.New(rand.NewSource(oldSeed)), totalBytes)

		newData := io.MultiReader(
			io.LimitReader(oldData, firstBytes),
			io.LimitReader(rand.New(rand.NewSource(newSeed)), newBytes),
			io.LimitReader(oldData, lastBytes),
		)

		var buf bytes.Buffer
		if err := Delta(s, newData, &buf); err != nil {
			b.Error(err)
		}

		b.Logf("raw   size:    %v bytes", totalBytes)
		b.Logf("delta size:    %v bytes (%.2f%%)", len(buf.Bytes()), (float64(len(buf.Bytes()))/float64(totalBytes))*100)
	}
}

func BenchmarkDeltaInpend1GB(b *testing.B) {
	benchmarkDeltaInpend(b, 1_000_000_000)
}

func BenchmarkDeltaInpend1MB(b *testing.B) {
	benchmarkDeltaInpend(b, 1_000_000)
}

//
// Variations of the Delta benchmark that remove 10% of the original data.
//

// Removes 10% of data from the end of the input data.
func benchmarkDeltaCutTail(b *testing.B, totalBytes int64) {
	newBytes := totalBytes - totalBytes/10
	oldSeed := time.Now().UnixNano()
	oldData := io.LimitReader(rand.New(rand.NewSource(oldSeed)), totalBytes)
	s := signature(b, oldData, int(totalBytes))

	b.SetBytes(totalBytes)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		newData := io.MultiReader(
			io.LimitReader(rand.New(rand.NewSource(oldSeed)), newBytes),
		)

		var buf bytes.Buffer
		if err := Delta(s, newData, &buf); err != nil {
			b.Error(err)
		}

		b.Logf("raw   size:    %v bytes", totalBytes)
		b.Logf("delta size:    %v bytes (%.2f%%)", len(buf.Bytes()), (float64(len(buf.Bytes()))/float64(totalBytes))*100)
	}
}

func BenchmarkDeltaCutTail1GB(b *testing.B) {
	benchmarkDeltaCutTail(b, 1_000_000_000)
}

func BenchmarkDeltaCutTail1MB(b *testing.B) {
	benchmarkDeltaCutTail(b, 1_000_000)
}

// Removes 10% of data from the beginning of the input data.
func benchmarkDeltaCutHead(b *testing.B, totalBytes int64) {
	oldSeed := time.Now().UnixNano()
	oldData := io.LimitReader(rand.New(rand.NewSource(oldSeed)), totalBytes)
	s := signature(b, oldData, int(totalBytes))

	b.SetBytes(totalBytes)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		discardBytes := totalBytes / 10
		newBytes := totalBytes - discardBytes

		newData := io.LimitReader(rand.New(rand.NewSource(oldSeed)), newBytes)
		buff := make([]byte, 1024)
		for discardBytes > 0 {
			n, err := newData.Read(buff)
			if err != nil {
				b.Fatal(err)
			}
			discardBytes -= int64(n)
		}

		var buf bytes.Buffer
		if err := Delta(s, newData, &buf); err != nil {
			b.Error(err)
		}

		b.Logf("raw   size:    %v bytes", totalBytes)
		b.Logf("delta size:    %v bytes (%.2f%%)", len(buf.Bytes()), (float64(len(buf.Bytes()))/float64(totalBytes))*100)
	}
}

func BenchmarkDeltaCutHead1GB(b *testing.B) {
	benchmarkDeltaCutHead(b, 1_000_000_000)
}

func BenchmarkDeltaCutHead1MB(b *testing.B) {
	benchmarkDeltaCutHead(b, 1_000_000)
}
