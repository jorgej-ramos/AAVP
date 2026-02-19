package token

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

type testVectorFile struct {
	Vectors []testVector `json:"vectors"`
}

type testVector struct {
	Name         string     `json:"name"`
	Fields       testFields `json:"fields"`
	ExpectedHex  string     `json:"expected_token_hex"`
	ExpectedSize int        `json:"expected_size"`
}

type testFields struct {
	TokenType     uint16 `json:"token_type"`
	Nonce         string `json:"nonce"`
	TokenKeyID    string `json:"token_key_id"`
	AgeBracket    string `json:"age_bracket"`
	AgeBracketVal uint8  `json:"age_bracket_value"`
	ExpiresAt     uint64 `json:"expires_at"`
	Authenticator string `json:"authenticator"`
}

func loadEncodingVectors(t *testing.T) []testVector {
	t.Helper()
	data, err := os.ReadFile("../../../test-vectors/token-encoding.json")
	if err != nil {
		t.Fatalf("failed to read test vectors: %v", err)
	}
	var f testVectorFile
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatalf("failed to parse test vectors: %v", err)
	}
	return f.Vectors
}

func hexToBytes(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("invalid hex string: %v", err)
	}
	return b
}

func hexToArray32(t *testing.T, s string) [32]byte {
	t.Helper()
	b := hexToBytes(t, s)
	if len(b) != 32 {
		t.Fatalf("expected 32 bytes, got %d", len(b))
	}
	var arr [32]byte
	copy(arr[:], b)
	return arr
}

func hexToArray256(t *testing.T, s string) [256]byte {
	t.Helper()
	b := hexToBytes(t, s)
	if len(b) != 256 {
		t.Fatalf("expected 256 bytes, got %d", len(b))
	}
	var arr [256]byte
	copy(arr[:], b)
	return arr
}

func TestEncodeDecodeVectors(t *testing.T) {
	vectors := loadEncodingVectors(t)
	for _, v := range vectors {
		t.Run(v.Name, func(t *testing.T) {
			tok := &Token{
				TokenType:     v.Fields.TokenType,
				Nonce:         hexToArray32(t, v.Fields.Nonce),
				TokenKeyID:    hexToArray32(t, v.Fields.TokenKeyID),
				AgeBracket:    v.Fields.AgeBracketVal,
				ExpiresAt:     v.Fields.ExpiresAt,
				Authenticator: hexToArray256(t, v.Fields.Authenticator),
			}

			// Encode and compare with expected hex.
			encoded := Encode(tok)
			expectedBytes := hexToBytes(t, v.ExpectedHex)
			if len(expectedBytes) != v.ExpectedSize {
				t.Fatalf("expected size %d, got %d", v.ExpectedSize, len(expectedBytes))
			}
			if len(encoded) != TokenSize {
				t.Fatalf("encoded size %d, expected %d", len(encoded), TokenSize)
			}
			for i := range encoded {
				if encoded[i] != expectedBytes[i] {
					t.Fatalf("byte mismatch at offset %d: got 0x%02x, want 0x%02x", i, encoded[i], expectedBytes[i])
				}
			}

			// Decode the expected hex and verify fields match.
			decoded, err := Decode(expectedBytes)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}
			if decoded.TokenType != v.Fields.TokenType {
				t.Errorf("TokenType: got %d, want %d", decoded.TokenType, v.Fields.TokenType)
			}
			if hex.EncodeToString(decoded.Nonce[:]) != v.Fields.Nonce {
				t.Errorf("Nonce mismatch")
			}
			if hex.EncodeToString(decoded.TokenKeyID[:]) != v.Fields.TokenKeyID {
				t.Errorf("TokenKeyID mismatch")
			}
			if decoded.AgeBracket != v.Fields.AgeBracketVal {
				t.Errorf("AgeBracket: got %d, want %d", decoded.AgeBracket, v.Fields.AgeBracketVal)
			}
			if decoded.ExpiresAt != v.Fields.ExpiresAt {
				t.Errorf("ExpiresAt: got %d, want %d", decoded.ExpiresAt, v.Fields.ExpiresAt)
			}
			if hex.EncodeToString(decoded.Authenticator[:]) != v.Fields.Authenticator {
				t.Errorf("Authenticator mismatch")
			}

			// Re-encode the decoded token and verify byte-for-byte.
			reencoded := Encode(decoded)
			if reencoded != encoded {
				t.Error("re-encoded token does not match original encoding")
			}
		})
	}
}

func TestDecodeInvalidSize(t *testing.T) {
	_, err := Decode(make([]byte, 330))
	if err == nil {
		t.Error("expected error for 330-byte input")
	}
	_, err = Decode(make([]byte, 332))
	if err == nil {
		t.Error("expected error for 332-byte input")
	}
	_, err = Decode(nil)
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestMessageToSign(t *testing.T) {
	tok := &Token{
		TokenType:  TokenTypeRSAPBSSASHA384,
		AgeBracket: AgeBracketOver18,
		ExpiresAt:  1772330400,
	}
	msg := tok.MessageToSign()
	if len(msg) != MessageToSignSize {
		t.Fatalf("MessageToSign length: got %d, want %d", len(msg), MessageToSignSize)
	}
}

func TestPublicMetadata(t *testing.T) {
	tok := &Token{
		AgeBracket: AgeBracketOver18,
		ExpiresAt:  1772330400,
	}
	meta := tok.PublicMetadata()
	if len(meta) != PublicMetadataSize {
		t.Fatalf("PublicMetadata length: got %d, want %d", len(meta), PublicMetadataSize)
	}
	if meta[0] != AgeBracketOver18 {
		t.Errorf("metadata[0]: got %d, want %d", meta[0], AgeBracketOver18)
	}
	expectedHex := "030000000069a39da0"
	if hex.EncodeToString(meta) != expectedHex {
		t.Errorf("metadata hex: got %s, want %s", hex.EncodeToString(meta), expectedHex)
	}
}

func TestAgeBracketName(t *testing.T) {
	tests := []struct {
		val  uint8
		name string
	}{
		{0, "UNDER_13"},
		{1, "AGE_13_15"},
		{2, "AGE_16_17"},
		{3, "OVER_18"},
	}
	for _, tt := range tests {
		if got := AgeBracketName(tt.val); got != tt.name {
			t.Errorf("AgeBracketName(%d) = %q, want %q", tt.val, got, tt.name)
		}
	}
}
