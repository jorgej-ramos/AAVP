package da

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/aavp-protocol/aavp-go/im"
	"github.com/aavp-protocol/aavp-go/internal/testkeys"
	"github.com/aavp-protocol/aavp-go/pbrsa"
	"github.com/aavp-protocol/aavp-go/token"
)

func TestPrepareWithValues(t *testing.T) {
	// Use the vector key for structural tests (no signing needed)
	sk := testkeys.VectorKey()
	keyID, _ := hex.DecodeString("fffea9ba9efa735080cf1af734625994ed056c1c4f94a8d82f4676a017ab2c7c")
	var kid [32]byte
	copy(kid[:], keyID)

	agent := NewDeviceAgentWithKeyID(&sk.PublicKey, kid)

	nonce, _ := hex.DecodeString("a6d9e1762e690e1e01f5a9c5df7b45ffa4850e59cec9ef0942319916cba412d8")
	var nonceArr [32]byte
	copy(nonceArr[:], nonce)

	result, err := agent.PrepareWithValues(token.AgeBracketOver18, nonceArr, 1772330400)
	if err != nil {
		t.Fatalf("PrepareWithValues: %v", err)
	}

	expectedMsg := "0001a6d9e1762e690e1e01f5a9c5df7b45ffa4850e59cec9ef0942319916cba412d8fffea9ba9efa735080cf1af734625994ed056c1c4f94a8d82f4676a017ab2c7c030000000069a39da0"
	gotMsg := hex.EncodeToString(result.Token.MessageToSign())
	if gotMsg != expectedMsg {
		t.Errorf("message_to_sign mismatch:\ngot:  %s\nwant: %s", gotMsg, expectedMsg)
	}

	expectedMeta := "030000000069a39da0"
	gotMeta := hex.EncodeToString(result.Metadata)
	if gotMeta != expectedMeta {
		t.Errorf("metadata mismatch:\ngot:  %s\nwant: %s", gotMeta, expectedMeta)
	}
}

func TestIssueTokenRoundTrip(t *testing.T) {
	// Use safe-prime key for operations involving DeriveKeyPair with random metadata
	sk := testkeys.SafePrimeKey()
	spkiDER, err := im.MarshalSPKIDER(&sk.PublicKey)
	if err != nil {
		t.Fatalf("MarshalSPKIDER: %v", err)
	}
	kid := sha256.Sum256(spkiDER)

	agent := NewDeviceAgentWithKeyID(&sk.PublicKey, kid)

	signer := func(blindedMsg, metadata []byte) ([]byte, error) {
		return pbrsa.BlindSign(sk, blindedMsg, metadata)
	}

	tok, err := agent.IssueToken(token.AgeBracketOver18, 3*time.Hour, signer)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}

	if tok.TokenType != token.TokenTypeRSAPBSSASHA384 {
		t.Errorf("token_type: got %d, want %d", tok.TokenType, token.TokenTypeRSAPBSSASHA384)
	}
	if tok.AgeBracket != token.AgeBracketOver18 {
		t.Errorf("age_bracket: got %d, want %d", tok.AgeBracket, token.AgeBracketOver18)
	}

	// Verify the token's signature
	encoded := token.Encode(tok)
	msg := encoded[:token.MessageToSignSize]
	metadata := tok.PublicMetadata()
	if err := pbrsa.Verify(&sk.PublicKey, msg, metadata, tok.Authenticator[:]); err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}
}

func TestStepByStepIssuance(t *testing.T) {
	sk := testkeys.VectorKey()
	keyID, _ := hex.DecodeString("fffea9ba9efa735080cf1af734625994ed056c1c4f94a8d82f4676a017ab2c7c")
	var kid [32]byte
	copy(kid[:], keyID)

	agent := NewDeviceAgentWithKeyID(&sk.PublicKey, kid)

	nonce, _ := hex.DecodeString("a6d9e1762e690e1e01f5a9c5df7b45ffa4850e59cec9ef0942319916cba412d8")
	var nonceArr [32]byte
	copy(nonceArr[:], nonce)

	// Step 1: Prepare
	prepResult, err := agent.PrepareWithValues(token.AgeBracketOver18, nonceArr, 1772330400)
	if err != nil {
		t.Fatalf("Prepare: %v", err)
	}

	// Step 2: Blind (with deterministic r)
	rFixed, _ := new(big.Int).SetString("aabbccdd00112233445566778899aabb", 16)
	blindResult, err := agent.Blind(prepResult.Token, prepResult.Metadata, rFixed)
	if err != nil {
		t.Fatalf("Blind: %v", err)
	}

	// Step 3-4: IM signs
	blindSig, err := pbrsa.BlindSign(sk, blindResult.BlindedMsg, prepResult.Metadata)
	if err != nil {
		t.Fatalf("BlindSign: %v", err)
	}

	// Step 5: Finalize
	err = agent.Finalize(prepResult.Token, blindSig, blindResult.State, prepResult.Metadata)
	if err != nil {
		t.Fatalf("Finalize: %v", err)
	}

	// Step 6: Verify
	encoded := token.Encode(prepResult.Token)
	msg := encoded[:token.MessageToSignSize]
	metadata := prepResult.Token.PublicMetadata()
	if err := pbrsa.Verify(&sk.PublicKey, msg, metadata, prepResult.Token.Authenticator[:]); err != nil {
		t.Fatalf("Verify: %v", err)
	}
}
