// Package uuidkit provides UUID generation and short code utilities.
//
// Supports:
//   - UUID v4 (random) and v7 (time-ordered) generation
//   - Short codes: compact hash-based identifiers from UUIDs (8-32 characters)
//   - Useful for user-facing IDs that are shorter than full UUIDs
package uuidkit

import (
	"crypto/sha256"
	"encoding/hex"

	_uuid "github.com/google/uuid"
)


// NewUUID generates a new UUID v4.
func NewUUID() string {
    return _uuid.New().String()
}

// NewUUIDv7 generates a time-ordered UUID v7.
func NewUUIDv7() _uuid.UUID {
    uuid, err := _uuid.NewV7()
    if err != nil {
        return _uuid.New()
    }
    return uuid
}


// NewUUIDv7AsString generates a time-ordered UUID v7 as a string.
func NewUUIDv7AsString() string {
    return NewUUIDv7().String()
}


// NewShortCode generates an 8-character short code from a new UUID.
func NewShortCode() string {
    return ShortCodeWithLength(nil, 8) // 8 is the default length
}

// ShortCode generates an 8-character short code from a UUID.
func ShortCode(value *_uuid.UUID) string {
    return ShortCodeWithLength(value, 8) // 8 is the default length
}


// ShortCodeWithLength generates a short code of specified length from a UUID.
// If value is nil, generates a new UUID. Length is clamped between 1-32.
func ShortCodeWithLength(value *_uuid.UUID, codeLength int) string {
    if value == nil {
        // If no UUID provided, generate a new one
        newUUID := NewUUID()
        uuidObj, _ := _uuid.Parse(newUUID)
        value = &uuidObj
    }
    if codeLength <= 0 {
        codeLength = 8  // Default to 8 if invalid length provided
    } else if codeLength > 32 {
        codeLength = 32  // Max length is 32
    }

    bytes, _ := value.MarshalBinary()
    hashBytes := sha256.Sum256(bytes)  // Always 32 bytes

    h := hex.EncodeToString(hashBytes[32 - ((codeLength + 1)/2):])
    return h[len(h) - codeLength:]  // may or may not need to trim a char
}


/*

    @staticmethod
    def random(castType: Optional[type] = None) -> Any:
        '''
        returns a UUID4 (random)
        '''
        return uuid.uuid4() if castType is None else castType( uuid.uuid4() )


    @staticmethod
    def uuid0(castType: Optional[type] = None) -> Any:
        '''
        returns an all 0 UUID
        '''
        return _UUID(int = 0) if castType is None else castType( _UUID(int = 0) )


    @staticmethod
    def uuid1(castType: Optional[type] = None) -> Any:
        '''
        returns an UUID1
        '''
        return uuid.uuid1() if castType is None else castType( uuid.uuid1() )


    @staticmethod
    def fromString(value: str) -> _UUID:
        '''
        returns an UUID from a string
        '''
        return _UUID(f'{{{value}}}')


    @staticmethod
    def getEpoch(value: _UUID) -> float:
        '''
        returns the epoch time of an UUID
        '''
        return (value.time - 0x01b21dd213814000) / 1e7


    @staticmethod
    def uuid7(castType: Optional[str]) -> Union[_UUID, str, int, bytes]:
        '''returns a UUIDv7
        cast types: 'bytes', 'hex', 'int', 'str', 'uuid', None
        '''
        castType = castType or 'uuid'
        return _uuid7(as_type=castType)

*/
