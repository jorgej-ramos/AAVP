// Package im implements the Implementor (IM) role of the AAVP protocol.
package im

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"time"

	"github.com/aavp-protocol/aavp-go/pbrsa"
)

// Implementor holds the IM's master private key and configuration.
type Implementor struct {
	PrivateKey *pbrsa.PrivateKey
	SPKIDER    []byte // SPKI DER encoding of the public key
	Domain     string
}

// NewImplementor creates a new Implementor from a private key.
// The SPKI DER is computed from the key components.
func NewImplementor(sk *pbrsa.PrivateKey, spkiDER []byte, domain string) *Implementor {
	return &Implementor{
		PrivateKey: sk,
		SPKIDER:    spkiDER,
		Domain:     domain,
	}
}

// Sign performs BlindSign on a blinded message with the given metadata.
func (im *Implementor) Sign(blindedMsg, metadata []byte) ([]byte, error) {
	return pbrsa.BlindSign(im.PrivateKey, blindedMsg, metadata)
}

// TokenKeyID returns the SHA-256 of the IM's public key in SPKI DER format.
func (im *Implementor) TokenKeyID() [32]byte {
	return sha256.Sum256(im.SPKIDER)
}

// WellKnownIssuer represents the .well-known/aavp-issuer response.
type WellKnownIssuer struct {
	Issuer          string          `json:"issuer"`
	AAVPVersion     string          `json:"aavp_version"`
	SigningEndpoint string          `json:"signing_endpoint"`
	Keys            []WellKnownKey  `json:"keys"`
}

// WellKnownKey represents a key entry in the .well-known/aavp-issuer response.
type WellKnownKey struct {
	TokenKeyID string `json:"token_key_id"`
	TokenType  uint16 `json:"token_type"`
	PublicKey  string `json:"public_key"`
	NotBefore  string `json:"not_before"`
	NotAfter   string `json:"not_after"`
}

// WellKnownResponse generates the .well-known/aavp-issuer JSON structure.
func (im *Implementor) WellKnownResponse(notBefore, notAfter time.Time) *WellKnownIssuer {
	keyID := im.TokenKeyID()
	keyIDBase64 := base64.RawURLEncoding.EncodeToString(keyID[:])
	pkBase64 := base64.RawURLEncoding.EncodeToString(im.SPKIDER)

	return &WellKnownIssuer{
		Issuer:          im.Domain,
		AAVPVersion:     "0.10",
		SigningEndpoint: "https://" + im.Domain + "/aavp/v1/sign",
		Keys: []WellKnownKey{
			{
				TokenKeyID: keyIDBase64,
				TokenType:  1,
				PublicKey:  pkBase64,
				NotBefore:  notBefore.UTC().Format(time.RFC3339),
				NotAfter:   notAfter.UTC().Format(time.RFC3339),
			},
		},
	}
}

// MarshalSPKIDER creates a SPKI DER encoding from the master public key components.
// Only works for master keys where E fits in int (e.g., e=65537).
func MarshalSPKIDER(pk *pbrsa.PublicKey) ([]byte, error) {
	stdPK := &rsa.PublicKey{
		N: pk.N,
		E: int(pk.E.Int64()),
	}
	return x509.MarshalPKIXPublicKey(stdPK)
}
