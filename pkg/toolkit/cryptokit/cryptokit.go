// Package cryptokit provides cryptographic hashing utilities.
//
// Two main use cases:
//
//  1. Short IDs: Blake2b-based Hash() generates compact identifiers
//     with custom base encoding (b32, b62, b64, hex)
//
//  2. Cryptographic hashing: SHA family functions for content addressing,
//     integrity verification, and security (MD5Hex, SHA1Hex, SHA256Hex, SHA512Hex)
//
// Example - Generate short ID:
//
//	id, _ := cryptokit.Hash("user@example.com")           // b62 encoding (default)
//	id, _ := cryptokit.Hash([]byte("data"), "hex")        // hex encoding
//
// Example - Content addressing:
//
//	hash := cryptokit.SHA256Hex(fileData)                  // for blob deduplication
//	hash := cryptokit.MD5Hex([]byte("password"))           // for checksums
package cryptokit

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/hex"
	"hash"
	"math/big"
	"strings"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
)


const (
	DEFAULT_DIGEST_SIZE = 12  // Consider moving to 16 depending on usage
	DEFAULT_BASE        = "b62"
)


// HashFunc is a function that returns a hash.Hash implementation.
type HashFunc func() hash.Hash

var (
	// It's an older code, sir, but it checks out.
	MD5  HashFunc = md5.New
	SHA1 HashFunc = sha1.New

	// SHA-2
	SHA224 HashFunc = sha256.New224
	SHA256 HashFunc = sha256.New
	SHA384 HashFunc = sha512.New384
	SHA512 HashFunc = sha512.New

	// SHA-3
	SHA3_224 HashFunc = sha3.New224
	SHA3_256 HashFunc = sha3.New256
	SHA3_384 HashFunc = sha3.New384
	SHA3_512 HashFunc = sha3.New512
)


// BaseEncoding defines a custom base encoding scheme with character set and base.
type BaseEncoding struct {
	CharSet string
	Length  int
}

// CustomBaseMap contains predefined base encoding schemes.
var CustomBaseMap = map[string]BaseEncoding{
	"b32": {CharSet: "0123456789ABCDEFGHJKMNPQRSTVWXYZ", Length: 32},
	"b62": {CharSet: "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", Length: 62},
	"b64": {CharSet: "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_", Length: 64},
	"hex": {CharSet: "0123456789abcdef", Length: 16},
}


// Hash generates a Blake2b hash with custom base encoding.
// Accepts string or []byte data. Default encoding is b62.
func Hash(data any, base ...string) (string, error) {
	var bytes []byte
	switch v := data.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		bytes = []byte(v.(string))
	}

	return HashBytes(bytes, base...)
}


// HashBytes generates a Blake2b hash from bytes with custom base encoding.
func HashBytes(data []byte, base ...string) (string, error) {
	encoding := DEFAULT_BASE
	if len(base) > 0 {
		encoding = base[0]
	}

	h, err := blake2b.New(DEFAULT_DIGEST_SIZE, nil)
	if err != nil {
		return "", err
	}
	h.Write(data)
	return Encode(h, encoding), nil
}


// Encode converts a hash to a string using the specified base encoding.
func Encode(rawHash hash.Hash, base string) string {
	if base == "hex" {
		return hex.EncodeToString(rawHash.Sum(nil))
	}

	digest := rawHash.Sum(nil)

	if base == "b32" {
		return strings.ToLower(strings.TrimRight(base32.StdEncoding.EncodeToString(digest), "="))
	}

	if encoding, ok := CustomBaseMap[base]; ok {
		intVal := new(big.Int).SetBytes(digest)
		return customBaseEncoding(intVal, encoding.CharSet, encoding.Length)
	}

	return ""
}


func customBaseEncoding(integer *big.Int, charSet string, charSetLength int) string {
	if integer.Sign() == 0 {
		return string(charSet[0])
	}

	digits := make([]byte, 0, 32)
	base   := big.NewInt(int64(charSetLength))
	mod    := new(big.Int)
	val    := new(big.Int).Set(integer)

	for val.Sign() > 0 {
		val.DivMod(val, base, mod)
		digits = append(digits, charSet[mod.Int64()])
	}

	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}

	return string(digits)
}


// SHA family convenience functions - all accept []byte and return hex string

// MD5Hex calculates MD5 hash and returns hex string.
func MD5Hex(data []byte) string {
	h := md5.Sum(data)
	return hex.EncodeToString(h[:])
}

// SHA1Hex calculates SHA1 hash and returns hex string.
func SHA1Hex(data []byte) string {
	h := sha1.Sum(data)
	return hex.EncodeToString(h[:])
}

// SHA256Hex calculates SHA256 hash and returns hex string.
func SHA256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// SHA512Hex calculates SHA512 hash and returns hex string.
func SHA512Hex(data []byte) string {
	h := sha512.Sum512(data)
	return hex.EncodeToString(h[:])
}
