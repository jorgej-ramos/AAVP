package vg

import (
	"crypto/sha256"
	"testing"
	"time"

	"github.com/aavp-protocol/aavp-go/da"
	"github.com/aavp-protocol/aavp-go/im"
	"github.com/aavp-protocol/aavp-go/internal/testkeys"
	"github.com/aavp-protocol/aavp-go/pbrsa"
	"github.com/aavp-protocol/aavp-go/token"
)

func setupProtocol(t *testing.T) (*da.DeviceAgent, *VerificationGate, *pbrsa.PrivateKey) {
	t.Helper()
	sk := testkeys.SafePrimeKey()
	spkiDER, err := im.MarshalSPKIDER(&sk.PublicKey)
	if err != nil {
		t.Fatalf("MarshalSPKIDER: %v", err)
	}
	kid := sha256.Sum256(spkiDER)

	agent := da.NewDeviceAgentWithKeyID(&sk.PublicKey, kid)
	gate := NewVerificationGate()
	gate.AddTrustedIM(kid, &sk.PublicKey)

	return agent, gate, sk
}

func TestFullProtocolRoundTrip(t *testing.T) {
	agent, gate, sk := setupProtocol(t)

	signer := func(blindedMsg, metadata []byte) ([]byte, error) {
		return pbrsa.BlindSign(sk, blindedMsg, metadata)
	}

	brackets := []uint8{
		token.AgeBracketUnder13,
		token.AgeBracketAge13_15,
		token.AgeBracketAge16_17,
		token.AgeBracketOver18,
	}

	for _, bracket := range brackets {
		t.Run(token.AgeBracketName(bracket), func(t *testing.T) {
			tok, err := agent.IssueToken(bracket, 3*time.Hour, signer)
			if err != nil {
				t.Fatalf("IssueToken: %v", err)
			}

			encoded := token.Encode(tok)
			result, err := gate.Verify(encoded[:], time.Now().UTC())
			if err != nil {
				t.Fatalf("Verify: %v", err)
			}

			if result.AgeBracket != bracket {
				t.Errorf("age_bracket: got %d, want %d", result.AgeBracket, bracket)
			}
		})
	}
}

func TestVerifyRejectsTamperedToken(t *testing.T) {
	agent, gate, sk := setupProtocol(t)

	signer := func(blindedMsg, metadata []byte) ([]byte, error) {
		return pbrsa.BlindSign(sk, blindedMsg, metadata)
	}

	tok, err := agent.IssueToken(token.AgeBracketOver18, 3*time.Hour, signer)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}

	encoded := token.Encode(tok)
	encoded[75] ^= 0x80 // Flip a bit in the authenticator

	_, err = gate.Verify(encoded[:], time.Now().UTC())
	if err == nil {
		t.Error("expected verification to fail for tampered token")
	}
}

func TestVerifyRejectsUnknownIM(t *testing.T) {
	agent, _, sk := setupProtocol(t)
	emptyGate := NewVerificationGate() // No trusted IMs

	signer := func(blindedMsg, metadata []byte) ([]byte, error) {
		return pbrsa.BlindSign(sk, blindedMsg, metadata)
	}

	tok, err := agent.IssueToken(token.AgeBracketOver18, 3*time.Hour, signer)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}

	encoded := token.Encode(tok)
	_, err = emptyGate.Verify(encoded[:], time.Now().UTC())
	if err == nil {
		t.Error("expected verification to fail for unknown IM")
	}
}

func TestVerifyRejectsExpiredToken(t *testing.T) {
	agent, gate, sk := setupProtocol(t)

	signer := func(blindedMsg, metadata []byte) ([]byte, error) {
		return pbrsa.BlindSign(sk, blindedMsg, metadata)
	}

	tok, err := agent.IssueToken(token.AgeBracketOver18, 1*time.Hour, signer)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}

	encoded := token.Encode(tok)
	// Verify at a time far in the future
	futureTime := time.Now().Add(24 * time.Hour).UTC()
	_, err = gate.Verify(encoded[:], futureTime)
	if err == nil {
		t.Error("expected verification to fail for expired token")
	}
}
