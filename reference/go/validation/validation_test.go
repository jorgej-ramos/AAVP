package validation

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/aavp-protocol/aavp-go/token"
)

type validationVectorFile struct {
	Constants struct {
		ClockSkewPast   int      `json:"CLOCK_SKEW_TOLERANCE_PAST"`
		ClockSkewFuture int      `json:"CLOCK_SKEW_TOLERANCE_FUTURE"`
		MaxTTLHours     int      `json:"MAX_TTL_HOURS"`
		ValidTypes      []uint16 `json:"VALID_TOKEN_TYPES"`
		ValidBrackets   []uint8  `json:"VALID_AGE_BRACKETS"`
		TokenSizeBytes  int      `json:"TOKEN_SIZE_BYTES"`
	} `json:"constants"`
	Vectors []validationVector `json:"vectors"`
}

type validationVector struct {
	Name          string `json:"name"`
	TokenHex      string `json:"token_hex"`
	TokenSize     int    `json:"token_size"`
	VGCurrentTime int64  `json:"vg_current_time"`
	ParsedFields  *struct {
		TokenType      uint16 `json:"token_type"`
		AgeBracket     string `json:"age_bracket"`
		AgeBracketVal  uint8  `json:"age_bracket_value"`
		ExpiresAt      uint64 `json:"expires_at"`
	} `json:"parsed_fields,omitempty"`
	ExpectedResult     string `json:"expected_result"`
	ExpectedError      string `json:"expected_error,omitempty"`
	ExpectedAgeBracket string `json:"expected_age_bracket,omitempty"`
}

func loadValidationVectors(t *testing.T) *validationVectorFile {
	t.Helper()
	data, err := os.ReadFile("../../../test-vectors/token-validation.json")
	if err != nil {
		t.Fatalf("failed to read test vectors: %v", err)
	}
	var f validationVectorFile
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatalf("failed to parse test vectors: %v", err)
	}
	return &f
}

func TestValidationConstants(t *testing.T) {
	f := loadValidationVectors(t)
	if f.Constants.ClockSkewPast != ClockSkewTolerancePast {
		t.Errorf("ClockSkewPast: got %d, want %d", ClockSkewTolerancePast, f.Constants.ClockSkewPast)
	}
	if f.Constants.ClockSkewFuture != ClockSkewToleranceFuture {
		t.Errorf("ClockSkewFuture: got %d, want %d", ClockSkewToleranceFuture, f.Constants.ClockSkewFuture)
	}
	if f.Constants.MaxTTLHours != MaxTTLHours {
		t.Errorf("MaxTTLHours: got %d, want %d", MaxTTLHours, f.Constants.MaxTTLHours)
	}
	if f.Constants.TokenSizeBytes != token.TokenSize {
		t.Errorf("TokenSize: got %d, want %d", token.TokenSize, f.Constants.TokenSizeBytes)
	}
}

func TestValidationVectors(t *testing.T) {
	f := loadValidationVectors(t)
	for _, v := range f.Vectors {
		t.Run(v.Name, func(t *testing.T) {
			tokenBytes, err := hex.DecodeString(v.TokenHex)
			if err != nil {
				t.Fatalf("invalid hex: %v", err)
			}
			if len(tokenBytes) != v.TokenSize {
				t.Fatalf("token_size mismatch: got %d, want %d", len(tokenBytes), v.TokenSize)
			}

			now := time.Unix(v.VGCurrentTime, 0).UTC()

			// For the signature_verification_failed test, use a callback that always fails.
			// For all other tests, skip signature verification (placeholder authenticators).
			var sigVerifier func([]byte) error
			if v.ExpectedError == "signature_verification_failed" {
				sigVerifier = func([]byte) error {
					return errors.New("bad signature")
				}
			}

			result, err := Validate(tokenBytes, now, sigVerifier)

			if v.ExpectedResult == "valid" {
				if err != nil {
					t.Fatalf("expected valid, got error: %v", err)
				}
				if v.ExpectedAgeBracket != "" {
					expectedName := v.ExpectedAgeBracket
					gotName := token.AgeBracketName(result.AgeBracket)
					if gotName != expectedName {
						t.Errorf("age_bracket: got %s, want %s", gotName, expectedName)
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
