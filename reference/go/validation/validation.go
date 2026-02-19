// Package validation implements the Verification Gate token validation logic
// as specified in PROTOCOL.md section 3.
package validation

import (
	"errors"
	"time"

	"github.com/aavp-protocol/aavp-go/token"
)

// Protocol constants.
const (
	ClockSkewTolerancePast   = 300   // seconds (5 minutes)
	ClockSkewToleranceFuture = 60    // seconds (1 minute)
	MaxTTLHours              = 4     // hours
	MaxTTLSeconds            = MaxTTLHours * 3600
)

// Validation errors.
var (
	ErrInvalidTokenSize           = errors.New("invalid_token_size")
	ErrUnsupportedTokenType       = errors.New("unsupported_token_type")
	ErrInvalidAgeBracket          = errors.New("invalid_age_bracket")
	ErrTokenExpired               = errors.New("token_expired")
	ErrExpiresAtTooFarFuture      = errors.New("expires_at_too_far_future")
	ErrSignatureVerificationFailed = errors.New("signature_verification_failed")
)

// ValidationResult contains the validated token information.
type ValidationResult struct {
	AgeBracket uint8
	ExpiresAt  time.Time
	TokenType  uint16
}

// AcceptedTokenTypes is the default set of accepted token types.
var AcceptedTokenTypes = []uint16{token.TokenTypeRSAPBSSASHA384}

// Validate checks a raw token according to VG validation rules.
// verifySignature is a callback that verifies the token's cryptographic signature.
// If verifySignature is nil, signature verification is skipped.
func Validate(tokenBytes []byte, now time.Time, verifySignature func([]byte) error) (*ValidationResult, error) {
	// 1. Size check.
	if len(tokenBytes) != token.TokenSize {
		return nil, ErrInvalidTokenSize
	}

	// 2. Decode fields.
	tok, err := token.Decode(tokenBytes)
	if err != nil {
		return nil, ErrInvalidTokenSize
	}

	// 3. token_type check.
	if !isAcceptedTokenType(tok.TokenType) {
		return nil, ErrUnsupportedTokenType
	}

	// 4. age_bracket check.
	if !token.ValidAgeBracket(tok.AgeBracket) {
		return nil, ErrInvalidAgeBracket
	}

	// 5. Expiration (past) check.
	nowUnix := uint64(now.Unix())
	expiresAt := tok.ExpiresAt
	if nowUnix > expiresAt && (nowUnix-expiresAt) > ClockSkewTolerancePast {
		return nil, ErrTokenExpired
	}

	// 6. Expiration (future) check.
	maxFuture := uint64(MaxTTLSeconds + ClockSkewToleranceFuture)
	if expiresAt > nowUnix && (expiresAt-nowUnix) > maxFuture {
		return nil, ErrExpiresAtTooFarFuture
	}

	// 7. Signature verification (if callback provided).
	if verifySignature != nil {
		if err := verifySignature(tokenBytes); err != nil {
			return nil, ErrSignatureVerificationFailed
		}
	}

	return &ValidationResult{
		AgeBracket: tok.AgeBracket,
		ExpiresAt:  time.Unix(int64(expiresAt), 0).UTC(),
		TokenType:  tok.TokenType,
	}, nil
}

func isAcceptedTokenType(tt uint16) bool {
	for _, accepted := range AcceptedTokenTypes {
		if tt == accepted {
			return true
		}
	}
	return false
}
