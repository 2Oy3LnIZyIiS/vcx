// Package compressionkit provides zstd compression utilities.
//
// Automatically determines if compression is beneficial based on:
//   - Minimum size: 512 bytes
//   - Compression ratio: compressed must be <95% of original
//
// Uses zstd compression with default speed settings for balanced performance.
// Provides both simple and buffer-reuse APIs for memory efficiency.
package compressionkit

import (
	"github.com/klauspost/compress/zstd"
)

const (
	MIN_COMPRESS_SIZE     = 512
	MIN_COMPRESSION_RATIO = 0.95 // Accept if compressed is < 95% of original
)

var (
	encoder *zstd.Encoder
	decoder *zstd.Decoder
)


func init() {
	var err error
	encoder, err = zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.SpeedDefault))
	if err != nil {
		panic(err)
	}
	decoder, err = zstd.NewReader(nil)
	if err != nil {
		panic(err)
	}
}


// Compress compresses data if beneficial (>512 bytes, <95% ratio).
// Returns compressed data and whether compression was applied.
func Compress(data []byte) ([]byte, bool) {
	return CompressBuffer(data, make([]byte, 0, len(data)))
}


// CompressBuffer compresses data using a provided buffer for output.
func CompressBuffer(data []byte, dst []byte) ([]byte, bool) {
	if len(data) < MIN_COMPRESS_SIZE {
		return data, false
	}

	compressed := encoder.EncodeAll(data, dst[:0])

	if float64(len(compressed)) < float64(len(data))*MIN_COMPRESSION_RATIO {
		return compressed, true
	}

	return data, false
}


// Decompress decompresses zstd-compressed data.
func Decompress(data []byte) ([]byte, error) {
	return DecompressBuffer(data, nil)
}


// DecompressBuffer decompresses data using a provided buffer for output.
func DecompressBuffer(data []byte, dest []byte) ([]byte, error) {
	if dest == nil {
		dest = make([]byte, 0, len(data)*2)
	}

	decompressed, err := decoder.DecodeAll(data, dest[:0])
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}
