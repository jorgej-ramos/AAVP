// Package da implements the Device Agent role of the AAVP protocol.
package da

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math/big"
	"time"

	"github.com/aavp-protocol/aavp-go/pbrsa"
	"github.com/aavp-protocol/aavp-go/token"
)

// SignerFunc represents the IM's blind signing service. It receives a blinded message
// and public metadata, and returns the blind signature.
type SignerFunc func(blindedMsg, metadata []byte) ([]byte, error)

// DeviceAgent holds the configuration for a Device Agent.
type DeviceAgent struct {
	IMPublicKey *pbrsa.PublicKey
	TokenKeyID  [32]byte
}

// NewDeviceAgent creates a new DeviceAgent from the IM's master public key.
// The TokenKeyID is computed as SHA-256 of the SPKI DER encoding.
func NewDeviceAgent(imPK *pbrsa.PublicKey, spkiDER []byte) *DeviceAgent {
	keyID := sha256.Sum256(spkiDER)
	return &DeviceAgent{
		IMPublicKey: imPK,
		TokenKeyID:  keyID,
	}
}

// NewDeviceAgentWithKeyID creates a DeviceAgent with an explicit token_key_id.
func NewDeviceAgentWithKeyID(imPK *pbrsa.PublicKey, keyID [32]byte) *DeviceAgent {
	return &DeviceAgent{
		IMPublicKey: imPK,
		TokenKeyID:  keyID,
	}
}

// PrepareResult contains the output of the Prepare step.
type PrepareResult struct {
	Token    *token.Token
	Metadata []byte // 9 bytes: age_bracket(1) || expires_at(8)
}

// Prepare builds a token with its fields and extracts the public metadata.
// The nonce is generated using CSPRNG. The expires_at is rounded to the nearest hour.
func (da *DeviceAgent) Prepare(ageBracket uint8, ttl time.Duration) (*PrepareResult, error) {
	if !token.ValidAgeBracket(ageBracket) {
		return nil, errors.New("da: invalid age bracket")
	}

	// Generate random nonce
	var nonce [32]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}

	// Compute expires_at: round to nearest hour
	expiresAt := time.Now().Add(ttl).Truncate(time.Hour).Unix()
	if expiresAt <= 0 {
		return nil, errors.New("da: expires_at must be positive")
	}

	tok := &token.Token{
		TokenType:  token.TokenTypeRSAPBSSASHA384,
		Nonce:      nonce,
		TokenKeyID: da.TokenKeyID,
		AgeBracket: ageBracket,
		ExpiresAt:  uint64(expiresAt),
	}

	return &PrepareResult{
		Token:    tok,
		Metadata: tok.PublicMetadata(),
	}, nil
}

// PrepareWithValues builds a token with explicit nonce and expires_at (for test vectors).
func (da *DeviceAgent) PrepareWithValues(ageBracket uint8, nonce [32]byte, expiresAt uint64) (*PrepareResult, error) {
	if !token.ValidAgeBracket(ageBracket) {
		return nil, errors.New("da: invalid age bracket")
	}

	tok := &token.Token{
		TokenType:  token.TokenTypeRSAPBSSASHA384,
		Nonce:      nonce,
		TokenKeyID: da.TokenKeyID,
		AgeBracket: ageBracket,
		ExpiresAt:  expiresAt,
	}

	return &PrepareResult{
		Token:    tok,
		Metadata: tok.PublicMetadata(),
	}, nil
}

// BlindResult contains the output of the Blind step.
type BlindResult struct {
	BlindedMsg []byte
	State      *pbrsa.BlindingState
}

// Blind blinds the message_to_sign for sending to the IM.
// If rFixed is not nil, uses that as the blinding factor (for deterministic tests).
func (da *DeviceAgent) Blind(tok *token.Token, metadata []byte, rFixed *big.Int) (*BlindResult, error) {
	msg := tok.MessageToSign()
	blindedMsg, state, err := pbrsa.Blind(da.IMPublicKey, msg, metadata, rFixed)
	if err != nil {
		return nil, err
	}
	return &BlindResult{BlindedMsg: blindedMsg, State: state}, nil
}

// Finalize unblinds the signature and sets the token's authenticator.
func (da *DeviceAgent) Finalize(tok *token.Token, blindSig []byte, state *pbrsa.BlindingState, metadata []byte) error {
	msg := tok.MessageToSign()
	authenticator, err := pbrsa.Finalize(da.IMPublicKey, msg, metadata, blindSig, state.Inv)
	if err != nil {
		return err
	}
	copy(tok.Authenticator[:], authenticator)
	return nil
}

// IssueToken performs the full issuance flow: Prepare → Blind → Sign (via callback) → Finalize.
func (da *DeviceAgent) IssueToken(ageBracket uint8, ttl time.Duration, signer SignerFunc) (*token.Token, error) {
	result, err := da.Prepare(ageBracket, ttl)
	if err != nil {
		return nil, err
	}

	blindResult, err := da.Blind(result.Token, result.Metadata, nil)
	if err != nil {
		return nil, err
	}

	blindSig, err := signer(blindResult.BlindedMsg, result.Metadata)
	if err != nil {
		return nil, err
	}

	if err := da.Finalize(result.Token, blindSig, blindResult.State, result.Metadata); err != nil {
		return nil, err
	}

	return result.Token, nil
}
