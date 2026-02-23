// Command generate computes the cryptographic test vector values (TO_BE_COMPUTED)
// in issuance-protocol.json using the PBRSA reference implementation.
//
// It generates a new RSA-2048 key with safe primes, computes all cryptographic
// values for the 4 test vectors, and writes the updated JSON file.
//
// Usage:
//
//	go run ./vectors/generate/
package main

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/aavp-protocol/aavp-go/pbrsa"
	"github.com/aavp-protocol/aavp-go/token"
)

func main() {
	vectorsPath := "../../test-vectors/issuance-protocol.json"

	// Read the existing vectors
	data, err := os.ReadFile(vectorsPath)
	if err != nil {
		fatalf("read %s: %v", vectorsPath, err)
	}

	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		fatalf("parse JSON: %v", err)
	}

	fmt.Println("Generating RSA-2048 safe-prime key...")
	sk, err := pbrsa.GenerateSafePrimeKey(2048)
	if err != nil {
		fatalf("generate key: %v", err)
	}
	fmt.Printf("  p is safe prime: %v\n", pbrsa.IsSafePrime(sk.P))
	fmt.Printf("  q is safe prime: %v\n", pbrsa.IsSafePrime(sk.Q))
	fmt.Printf("  n bits: %d\n", sk.N.BitLen())

	// Marshal SPKI DER
	spkiDER, err := x509.MarshalPKIXPublicKey(&rsa.PublicKey{N: sk.N, E: int(sk.E.Int64())})
	if err != nil {
		fatalf("marshal SPKI: %v", err)
	}

	tokenKeyID := sha256.Sum256(spkiDER)
	tokenKeyIDHex := hex.EncodeToString(tokenKeyID[:])
	tokenKeyIDBase64 := base64.RawURLEncoding.EncodeToString(tokenKeyID[:])
	pkBase64 := base64.RawURLEncoding.EncodeToString(spkiDER)

	fmt.Printf("  token_key_id: %s\n", tokenKeyIDHex)

	// Update test_im_key
	testIMKey := map[string]any{
		"description":          "Clave RSA-2048 de test del Implementador con safe primes. Generada por la implementacion de referencia Go. No usar en produccion.",
		"algorithm":            "RSA",
		"key_size_bits":        2048,
		"safe_primes":          true,
		"n":                    hex.EncodeToString(sk.N.Bytes()),
		"e":                    hex.EncodeToString(sk.E.Bytes()),
		"d":                    hex.EncodeToString(sk.D.Bytes()),
		"p":                    hex.EncodeToString(sk.P.Bytes()),
		"q":                    hex.EncodeToString(sk.Q.Bytes()),
		"spki_der_hex":         hex.EncodeToString(spkiDER),
		"token_key_id_hex":     tokenKeyIDHex,
		"token_key_id_base64url": tokenKeyIDBase64,
		"public_key_base64url": pkBase64,
		"well_known_aavp_issuer_example": map[string]any{
			"note":             "Ejemplo de como apareceria esta clave en el endpoint .well-known/aavp-issuer del IM (PROTOCOL.md seccion 5.2.3).",
			"issuer":           "test-im.example",
			"aavp_version":     "0.11",
			"signing_endpoint": "https://test-im.example/aavp/v1/sign",
			"keys": []map[string]any{
				{
					"token_key_id": tokenKeyIDBase64,
					"token_type":   1,
					"public_key":   pkBase64,
					"not_before":   "2026-01-01T00:00:00Z",
					"not_after":    "2026-06-30T00:00:00Z",
				},
			},
		},
	}
	doc["test_im_key"] = testIMKey

	// Process each vector
	vectors, ok := doc["vectors"].([]any)
	if !ok {
		fatalf("vectors field is not an array")
	}

	// Deterministic blinding factors (one per vector, derived from SHA-256 for reproducibility)
	rSeeds := []string{
		"aavp-test-blinding-factor-over18",
		"aavp-test-blinding-factor-under13",
		"aavp-test-blinding-factor-age1315",
		"aavp-test-blinding-factor-age1617",
	}

	for i, vAny := range vectors {
		v := vAny.(map[string]any)
		fmt.Printf("\nProcessing vector %d: %s\n", i+1, v["name"])

		step1 := v["step_1_prepare"].(map[string]any)

		// Build token from step 1 fields
		nonceHex := step1["nonce"].(string)
		nonce := hexDec(nonceHex)
		var nonceArr [32]byte
		copy(nonceArr[:], nonce)

		ageBracket := uint8(step1["age_bracket_value"].(float64))
		expiresAt := uint64(step1["expires_at"].(float64))

		tok := &token.Token{
			TokenType:  token.TokenTypeRSAPBSSASHA384,
			Nonce:      nonceArr,
			TokenKeyID: tokenKeyID,
			AgeBracket: ageBracket,
			ExpiresAt:  expiresAt,
		}

		msg := tok.MessageToSign()
		metadata := tok.PublicMetadata()

		// Update step_1 with new token_key_id and recomputed message
		step1["token_key_id"] = tokenKeyIDHex
		step1["message_to_sign"] = hex.EncodeToString(msg)
		step1["public_metadata"] = hex.EncodeToString(metadata)
		metaFields := step1["public_metadata_fields"].(map[string]any)
		metaFields["age_bracket"] = hex.EncodeToString([]byte{ageBracket})
		metaFields["expires_at"] = hex.EncodeToString(pbrsa.I2OSP(new(big.Int).SetUint64(expiresAt), 8))

		// Generate deterministic blinding factor
		rSeed := sha256.Sum256([]byte(rSeeds[i]))
		// Extend to 256 bytes for a full-size r
		rBytes := make([]byte, pbrsa.ModulusLen)
		for j := 0; j < pbrsa.ModulusLen; j += 32 {
			seed := sha256.Sum256(append(rSeed[:], byte(j/32)))
			copy(rBytes[j:], seed[:])
		}
		r := new(big.Int).SetBytes(rBytes)
		r.Mod(r, sk.N) // Ensure r < n
		if r.Sign() == 0 {
			r.SetInt64(1)
		}

		// Step 2: Blind
		blindedMsg, state, err := pbrsa.Blind(&sk.PublicKey, msg, metadata, r)
		if err != nil {
			fatalf("Blind vector %d: %v", i+1, err)
		}

		// Step 3: DeriveKey (for reporting)
		skDerived, pkDerived, err := pbrsa.DeriveKeyPair(sk, metadata)
		if err != nil {
			fatalf("DeriveKeyPair vector %d: %v", i+1, err)
		}

		// Step 4: BlindSign
		blindSig, err := pbrsa.BlindSign(sk, blindedMsg, metadata)
		if err != nil {
			fatalf("BlindSign vector %d: %v", i+1, err)
		}

		// Step 5: Finalize
		sig, err := pbrsa.Finalize(&sk.PublicKey, msg, metadata, blindSig, state.Inv)
		if err != nil {
			fatalf("Finalize vector %d: %v", i+1, err)
		}

		// Step 6: Verify
		if err := pbrsa.Verify(&sk.PublicKey, msg, metadata, sig); err != nil {
			fatalf("Verify vector %d: %v", i+1, err)
		}

		// Build complete token
		copy(tok.Authenticator[:], sig)
		fullToken := token.Encode(tok)

		fmt.Printf("  authenticator: %s...\n", hex.EncodeToString(sig[:16]))
		fmt.Printf("  token size: %d bytes\n", len(fullToken))

		// Update step_2_blind
		step2 := v["step_2_blind"].(map[string]any)
		step2Inputs := step2["inputs"].(map[string]any)
		step2Inputs["public_key"] = "pk' (derivada en step_3)"
		step2Inputs["prepared_msg"] = hex.EncodeToString(msg)
		step2Inputs["metadata"] = hex.EncodeToString(metadata)
		step2Outputs := step2["outputs"].(map[string]any)
		step2Outputs["blinded_msg"] = hex.EncodeToString(blindedMsg)
		step2Outputs["blinding_inverse_inv"] = hex.EncodeToString(pbrsa.I2OSP(state.Inv, pbrsa.ModulusLen))
		step2Rand := step2["randomness"].(map[string]any)
		step2Rand["blinding_factor_r"] = hex.EncodeToString(pbrsa.I2OSP(r, pbrsa.ModulusLen))
		step2Rand["note"] = "r derivado deterministicamente de SHA-256('" + rSeeds[i] + "') para reproducibilidad."

		// Update step_3_derive_key
		step3 := v["step_3_derive_key"].(map[string]any)
		step3Inputs := step3["inputs"].(map[string]any)
		step3Inputs["master_secret_key"] = "sk (del test_im_key)"
		step3Inputs["metadata"] = hex.EncodeToString(metadata)
		step3Inputs["hkdf_info"] = "AAVP-RSAPBSSA-SHA384-metadata-" + hex.EncodeToString(metadata)
		step3Outputs := step3["outputs"].(map[string]any)
		step3Outputs["derived_private_key_sk_prime"] = hex.EncodeToString(pbrsa.I2OSP(skDerived.D, pbrsa.ModulusLen))
		step3Outputs["derived_public_key_pk_prime"] = hex.EncodeToString(pbrsa.I2OSP(pkDerived.E, pbrsa.LambdaLen))

		// Update step_4_blind_sign
		step4 := v["step_4_blind_sign"].(map[string]any)
		step4Outputs := step4["outputs"].(map[string]any)
		step4Outputs["blind_sig"] = hex.EncodeToString(blindSig)

		// Update step_5_finalize
		step5 := v["step_5_finalize"].(map[string]any)
		step5Inputs := step5["inputs"].(map[string]any)
		step5Inputs["original_message"] = hex.EncodeToString(msg)
		step5Inputs["metadata"] = hex.EncodeToString(metadata)
		step5Outputs := step5["outputs"].(map[string]any)
		step5Outputs["authenticator"] = hex.EncodeToString(sig)

		// Update expected_token
		expectedToken := v["expected_token"].(map[string]any)
		expectedToken["token_hex"] = hex.EncodeToString(fullToken[:])

		// Update status
		v["status"] = "computed"
		v["status_note"] = "Valores criptograficos generados por la implementacion de referencia Go (github.com/aavp-protocol/aavp-go). Clave RSA-2048 con safe primes."
	}

	// Update generation_status
	doc["generation_status"] = map[string]any{
		"structural_fields":    "complete",
		"cryptographic_values": "complete",
		"note":                 "Todos los valores criptograficos han sido computados por la implementacion de referencia Go. La clave de test usa safe primes como requiere el draft PBRSA.",
		"generated_by":         "github.com/aavp-protocol/aavp-go vectors/generate",
	}

	// Write output
	output, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		fatalf("marshal: %v", err)
	}
	output = append(output, '\n')

	if err := os.WriteFile(vectorsPath, output, 0644); err != nil {
		fatalf("write: %v", err)
	}

	fmt.Printf("\nDone. Updated %s with computed cryptographic values.\n", vectorsPath)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	os.Exit(1)
}

func hexDec(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		fatalf("bad hex %q: %v", s, err)
	}
	return b
}

