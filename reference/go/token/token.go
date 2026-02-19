// Package token implements encoding and decoding of the AAVP 331-byte binary token format
// as specified in PROTOCOL.md section 2.
package token

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	// TokenSize is the fixed size of an AAVP token in bytes.
	TokenSize = 331

	// Field offsets within the token.
	OffsetTokenType     = 0
	OffsetNonce         = 2
	OffsetTokenKeyID    = 34
	OffsetAgeBracket    = 66
	OffsetExpiresAt     = 67
	OffsetAuthenticator = 75

	// Field sizes.
	SizeTokenType     = 2
	SizeNonce         = 32
	SizeTokenKeyID    = 32
	SizeAgeBracket    = 1
	SizeExpiresAt     = 8
	SizeAuthenticator = 256

	// MessageToSignSize is the size of the portion signed (everything except authenticator).
	MessageToSignSize = 75

	// PublicMetadataSize is the size of the public metadata (age_bracket + expires_at).
	PublicMetadataSize = 9
)

// Token type values.
const (
	TokenTypeReserved       uint16 = 0x0000
	TokenTypeRSAPBSSASHA384 uint16 = 0x0001
)

// Age bracket values.
const (
	AgeBracketUnder13  uint8 = 0
	AgeBracketAge13_15 uint8 = 1
	AgeBracketAge16_17 uint8 = 2
	AgeBracketOver18   uint8 = 3
)

// AgeBracketName returns the canonical name for an age bracket value.
func AgeBracketName(b uint8) string {
	switch b {
	case AgeBracketUnder13:
		return "UNDER_13"
	case AgeBracketAge13_15:
		return "AGE_13_15"
	case AgeBracketAge16_17:
		return "AGE_16_17"
	case AgeBracketOver18:
		return "OVER_18"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", b)
	}
}

// ValidAgeBracket returns true if the given value is a valid age bracket.
func ValidAgeBracket(b uint8) bool {
	return b <= AgeBracketOver18
}

// Token represents an AAVP token with its six fields.
type Token struct {
	TokenType     uint16
	Nonce         [SizeNonce]byte
	TokenKeyID    [SizeTokenKeyID]byte
	AgeBracket    uint8
	ExpiresAt     uint64
	Authenticator [SizeAuthenticator]byte
}

// Encode serializes a Token into the 331-byte binary format.
func Encode(t *Token) [TokenSize]byte {
	var buf [TokenSize]byte
	binary.BigEndian.PutUint16(buf[OffsetTokenType:], t.TokenType)
	copy(buf[OffsetNonce:], t.Nonce[:])
	copy(buf[OffsetTokenKeyID:], t.TokenKeyID[:])
	buf[OffsetAgeBracket] = t.AgeBracket
	binary.BigEndian.PutUint64(buf[OffsetExpiresAt:], t.ExpiresAt)
	copy(buf[OffsetAuthenticator:], t.Authenticator[:])
	return buf
}

// Decode deserializes a byte slice into a Token. Returns an error if the size is not exactly 331 bytes.
func Decode(b []byte) (*Token, error) {
	if len(b) != TokenSize {
		return nil, errors.New("invalid token size: expected 331 bytes")
	}
	t := &Token{}
	t.TokenType = binary.BigEndian.Uint16(b[OffsetTokenType:])
	copy(t.Nonce[:], b[OffsetNonce:OffsetNonce+SizeNonce])
	copy(t.TokenKeyID[:], b[OffsetTokenKeyID:OffsetTokenKeyID+SizeTokenKeyID])
	t.AgeBracket = b[OffsetAgeBracket]
	t.ExpiresAt = binary.BigEndian.Uint64(b[OffsetExpiresAt:])
	copy(t.Authenticator[:], b[OffsetAuthenticator:OffsetAuthenticator+SizeAuthenticator])
	return t, nil
}

// MessageToSign returns the first 75 bytes of the encoded token (everything except the authenticator).
func (t *Token) MessageToSign() []byte {
	buf := Encode(t)
	result := make([]byte, MessageToSignSize)
	copy(result, buf[:MessageToSignSize])
	return result
}

// PublicMetadata returns the 9-byte public metadata: age_bracket (1 byte) || expires_at (8 bytes BE).
func (t *Token) PublicMetadata() []byte {
	meta := make([]byte, PublicMetadataSize)
	meta[0] = t.AgeBracket
	binary.BigEndian.PutUint64(meta[1:], t.ExpiresAt)
	return meta
}
