// Package vectors provides test vector verification against the JSON files
// in test-vectors/.
package vectors

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/aavp-protocol/aavp-go/token"
	"github.com/aavp-protocol/aavp-go/validation"
)

// --- Token Encoding Vectors ---

type encodingVectorFile struct {
	Vectors []encodingVector `json:"vectors"`
}

type encodingVector struct {
	Name         string         `json:"name"`
	Fields       encodingFields `json:"fields"`
	ExpectedHex  string         `json:"expected_token_hex"`
	ExpectedSize int            `json:"expected_size"`
}

type encodingFields struct {
	TokenType     uint16 `json:"token_type"`
	Nonce         string `json:"nonce"`
	TokenKeyID    string `json:"token_key_id"`
	AgeBracketVal uint8  `json:"age_bracket_value"`
	ExpiresAt     uint64 `json:"expires_at"`
	Authenticator string `json:"authenticator"`
}

func TestTokenEncodingVectors(t *testing.T) {
	data, err := os.ReadFile("../../../test-vectors/token-encoding.json")
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var f encodingVectorFile
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatalf("parse: %v", err)
	}

	for _, v := range f.Vectors {
		t.Run(v.Name, func(t *testing.T) {
			nonce := hexTo32(t, v.Fields.Nonce)
			keyID := hexTo32(t, v.Fields.TokenKeyID)
			auth := hexTo256(t, v.Fields.Authenticator)

			tok := &token.Token{
				TokenType:     v.Fields.TokenType,
				Nonce:         nonce,
				TokenKeyID:    keyID,
				AgeBracket:    v.Fields.AgeBracketVal,
				ExpiresAt:     v.Fields.ExpiresAt,
				Authenticator: auth,
			}

			encoded := token.Encode(tok)
			expected := hexToBytes(t, v.ExpectedHex)

			if len(encoded) != v.ExpectedSize {
				t.Fatalf("size: got %d, want %d", len(encoded), v.ExpectedSize)
			}
			for i := range encoded {
				if encoded[i] != expected[i] {
					t.Fatalf("mismatch at byte %d: got 0x%02x, want 0x%02x", i, encoded[i], expected[i])
				}
			}

			// Decode and re-encode
			decoded, err := token.Decode(expected)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			reencoded := token.Encode(decoded)
			if reencoded != encoded {
				t.Error("re-encoded mismatch")
			}
		})
	}
}

// --- Token Validation Vectors ---

type validationVectorFile struct {
	Vectors []validationVector `json:"vectors"`
}

type validationVector struct {
	Name               string `json:"name"`
	TokenHex           string `json:"token_hex"`
	TokenSize          int    `json:"token_size"`
	VGCurrentTime      int64  `json:"vg_current_time"`
	ExpectedResult     string `json:"expected_result"`
	ExpectedError      string `json:"expected_error,omitempty"`
	ExpectedAgeBracket string `json:"expected_age_bracket,omitempty"`
}

func TestTokenValidationVectors(t *testing.T) {
	data, err := os.ReadFile("../../../test-vectors/token-validation.json")
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var f validationVectorFile
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatalf("parse: %v", err)
	}

	for _, v := range f.Vectors {
		t.Run(v.Name, func(t *testing.T) {
			tokenBytes := hexToBytes(t, v.TokenHex)
			now := time.Unix(v.VGCurrentTime, 0).UTC()

			var sigVerifier func([]byte) error
			if v.ExpectedError == "signature_verification_failed" {
				sigVerifier = func([]byte) error { return errors.New("bad sig") }
			}

			result, err := validation.Validate(tokenBytes, now, sigVerifier)

			if v.ExpectedResult == "valid" {
				if err != nil {
					t.Fatalf("expected valid, got: %v", err)
				}
				if v.ExpectedAgeBracket != "" {
					got := token.AgeBracketName(result.AgeBracket)
					if got != v.ExpectedAgeBracket {
						t.Errorf("bracket: got %s, want %s", got, v.ExpectedAgeBracket)
					}
				}
			} else {
				if err == nil {
					t.Fatal("expected error, got valid")
				}
				if err.Error() != v.ExpectedError {
					t.Errorf("error: got %q, want %q", err.Error(), v.ExpectedError)
				}
			}
		})
	}
}

// --- Issuance Protocol Vectors (structural checks) ---

type issuanceVectorFile struct {
	TestIMKey struct {
		N           string `json:"n"`
		E           string `json:"e"`
		D           string `json:"d"`
		P           string `json:"p"`
		Q           string `json:"q"`
		TokenKeyID  string `json:"token_key_id_hex"`
		SPKIDERHex  string `json:"spki_der_hex"`
	} `json:"test_im_key"`
	Vectors []issuanceVector `json:"vectors"`
}

type issuanceVector struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Step1       struct {
		TokenType    uint16 `json:"token_type"`
		Nonce        string `json:"nonce"`
		TokenKeyID   string `json:"token_key_id"`
		AgeBracket   uint8  `json:"age_bracket_value"`
		ExpiresAt    uint64 `json:"expires_at"`
		MsgToSign    string `json:"message_to_sign"`
		MsgSize      int    `json:"message_to_sign_size"`
		Metadata     string `json:"public_metadata"`
		MetadataSize int    `json:"public_metadata_size"`
	} `json:"step_1_prepare"`
}

func TestIssuanceVectorsStructural(t *testing.T) {
	data, err := os.ReadFile("../../../test-vectors/issuance-protocol.json")
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var f issuanceVectorFile
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatalf("parse: %v", err)
	}

	for _, v := range f.Vectors {
		t.Run(v.Name, func(t *testing.T) {
			// Verify Step 1 structural correctness
			nonce := hexTo32(t, v.Step1.Nonce)
			keyID := hexTo32(t, v.Step1.TokenKeyID)

			tok := &token.Token{
				TokenType:  v.Step1.TokenType,
				Nonce:      nonce,
				TokenKeyID: keyID,
				AgeBracket: v.Step1.AgeBracket,
				ExpiresAt:  v.Step1.ExpiresAt,
			}

			msg := tok.MessageToSign()
			if len(msg) != v.Step1.MsgSize {
				t.Errorf("message size: got %d, want %d", len(msg), v.Step1.MsgSize)
			}
			if hex.EncodeToString(msg) != v.Step1.MsgToSign {
				t.Errorf("message_to_sign mismatch")
			}

			meta := tok.PublicMetadata()
			if len(meta) != v.Step1.MetadataSize {
				t.Errorf("metadata size: got %d, want %d", len(meta), v.Step1.MetadataSize)
			}
			if hex.EncodeToString(meta) != v.Step1.Metadata {
				t.Errorf("metadata mismatch:\ngot:  %s\nwant: %s", hex.EncodeToString(meta), v.Step1.Metadata)
			}
		})
	}
}

// Helpers

func hexToBytes(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("bad hex: %v", err)
	}
	return b
}

func hexTo32(t *testing.T, s string) [32]byte {
	t.Helper()
	b := hexToBytes(t, s)
	var arr [32]byte
	copy(arr[:], b)
	return arr
}

func hexTo256(t *testing.T, s string) [256]byte {
	t.Helper()
	b := hexToBytes(t, s)
	var arr [256]byte
	copy(arr[:], b)
	return arr
}
