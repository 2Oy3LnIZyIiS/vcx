// Package filekit provides file content analysis utilities.
//
// Primary use case is binary detection for determining if files should be compressed.
// Uses the same heuristic as Git: checks first 8KB for null bytes.
package filekit

const binarySampleSize = 8192 // Check first 8KB like Git

// IsBinary checks if data contains null bytes (indicating binary content).
func IsBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	sampleSize := min(len(data), binarySampleSize)

	for _, b := range data[:sampleSize] {
		if b == 0 {
			return true
		}
	}

	return false
}
