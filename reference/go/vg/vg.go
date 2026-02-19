// Package vg implements the Verification Gate (VG) role of the AAVP protocol.
package vg

import (
	"errors"
	"time"

	"github.com/aavp-protocol/aavp-go/pbrsa"
	"github.com/aavp-protocol/aavp-go/token"
	"github.com/aavp-protocol/aavp-go/validation"
)

// VerificationGate holds the trust store and configuration for a VG.
type VerificationGate struct {
	// TrustStore maps token_key_id (hex) to the IM's master public key.
	TrustStore map[[32]byte]*pbrsa.PublicKey
}

// NewVerificationGate creates a new VG with the given trust store entries.
func NewVerificationGate() *VerificationGate {
	return &VerificationGate{
		TrustStore: make(map[[32]byte]*pbrsa.PublicKey),
	}
}

// AddTrustedIM adds an IM's master public key to the trust store.
func (vg *VerificationGate) AddTrustedIM(tokenKeyID [32]byte, pk *pbrsa.PublicKey) {
	vg.TrustStore[tokenKeyID] = pk
}

// VerificationResult contains the result of a successful token verification.
type VerificationResult struct {
	AgeBracket uint8
	ExpiresAt  time.Time
}

// Verify performs full validation of a token: structural checks, temporal checks,
// and cryptographic signature verification.
func (vg *VerificationGate) Verify(tokenBytes []byte, now time.Time) (*VerificationResult, error) {
	sigVerifier := func(tb []byte) error {
		return vg.verifySignature(tb)
	}

	result, err := validation.Validate(tokenBytes, now, sigVerifier)
	if err != nil {
		return nil, err
	}

	return &VerificationResult{
		AgeBracket: result.AgeBracket,
		ExpiresAt:  result.ExpiresAt,
	}, nil
}

// verifySignature performs the cryptographic signature verification.
func (vg *VerificationGate) verifySignature(tokenBytes []byte) error {
	tok, err := token.Decode(tokenBytes)
	if err != nil {
		return err
	}

	// Look up the IM's master public key
	pk, ok := vg.TrustStore[tok.TokenKeyID]
	if !ok {
		return errors.New("unknown token_key_id")
	}

	// Extract message_to_sign (first 75 bytes) and metadata
	msg := tokenBytes[:token.MessageToSignSize]
	metadata := tok.PublicMetadata()
	sig := tok.Authenticator[:]

	// Verify the partially blind RSA signature
	return pbrsa.Verify(pk, msg, metadata, sig)
}
