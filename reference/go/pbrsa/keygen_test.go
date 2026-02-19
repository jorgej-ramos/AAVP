package pbrsa

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestGenerateSafePrimeKey(t *testing.T) {
	if testing.Short() {
		t.Skip("safe prime generation is slow")
	}
	sk, err := GenerateSafePrimeKey(2048)
	if err != nil {
		t.Fatalf("GenerateSafePrimeKey: %v", err)
	}
	if !IsSafePrime(sk.P) {
		t.Error("p is not a safe prime")
	}
	if !IsSafePrime(sk.Q) {
		t.Error("q is not a safe prime")
	}
	if sk.N.BitLen() != 2048 {
		t.Errorf("key size: got %d bits, want 2048", sk.N.BitLen())
	}

	// Verify d*e mod phi == 1
	pMinus1 := new(big.Int).Sub(sk.P, big.NewInt(1))
	qMinus1 := new(big.Int).Sub(sk.Q, big.NewInt(1))
	phi := new(big.Int).Mul(pMinus1, qMinus1)
	product := new(big.Int).Mul(sk.D, sk.E)
	product.Mod(product, phi)
	if product.Cmp(big.NewInt(1)) != 0 {
		t.Error("d*e mod phi != 1")
	}

	// Print key for use in test vectors (only on verbose)
	t.Logf("n = %s", hex.EncodeToString(sk.N.Bytes()))
	t.Logf("e = %x", sk.E)
	t.Logf("d = %s", hex.EncodeToString(sk.D.Bytes()))
	t.Logf("p = %s", hex.EncodeToString(sk.P.Bytes()))
	t.Logf("q = %s", hex.EncodeToString(sk.Q.Bytes()))

	// Verify DeriveKeyPair works with this key for all bracket metadata
	brackets := []string{
		"030000000069a39da0",
		"000000000069cc6000",
		"01000000006a2fe940",
		"020000000069b66700",
	}
	for _, infoHex := range brackets {
		info, _ := hex.DecodeString(infoHex)
		_, _, err := DeriveKeyPair(sk, info)
		if err != nil {
			t.Errorf("DeriveKeyPair failed for info %s: %v", infoHex, err)
		}
	}

	// Also test with random metadata
	for i := 0; i < 100; i++ {
		info := []byte(fmt.Sprintf("random-metadata-%d", i))
		_, _, err := DeriveKeyPair(sk, info)
		if err != nil {
			t.Errorf("DeriveKeyPair failed for random info %d: %v", i, err)
		}
	}
}

func TestIsSafePrime(t *testing.T) {
	// 7 = 2*3 + 1, and 3 is prime => 7 is a safe prime
	if !IsSafePrime(big.NewInt(7)) {
		t.Error("7 should be a safe prime")
	}
	// 11 = 2*5 + 1, and 5 is prime => 11 is a safe prime
	if !IsSafePrime(big.NewInt(11)) {
		t.Error("11 should be a safe prime")
	}
	// 9 is not prime
	if IsSafePrime(big.NewInt(9)) {
		t.Error("9 should not be a safe prime")
	}
	// 13 = 2*6 + 1, 6 is not prime => 13 is not a safe prime
	if IsSafePrime(big.NewInt(13)) {
		t.Error("13 should not be a safe prime (13 is prime but (13-1)/2=6 is not)")
	}
}
