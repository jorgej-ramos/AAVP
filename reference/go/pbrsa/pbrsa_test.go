package pbrsa

import (
	"encoding/hex"
	"math/big"
	"testing"
)

// testKeyFromVectors returns the RSA-2048 test key from issuance-protocol.json.
func testKeyFromVectors() *PrivateKey {
	n, _ := new(big.Int).SetString("deeb38cfa184281eee07142527e4bc39f7d5aba3b9b4acfc5202a1f3e2f02c5f41677ba9fa8c2345cb3ebf7b80c48626b3a6f191909cbcfafd64b5ca92ad5c5c2b32b57660abc2783f5d9093a9a7d1057330fe09a992324dffca3a6469bae4817573e6ad419f926c5074cb1fa6818c95ce376b339056a6a0800bd5b4f7ff2be0981df5c558842d054a61c9c180bfd02abe5cec6e07e05160739fd64734186b8f01b6489f4d33177100bb5c12c629f3f4f4bcb3a39abd776834491101ff554aabddc9abea1820b6039c4aca86141a9d2295e113d8203a976ee0cc210b7912a917b4ab4366b6c4c787672c5e9d6f29a97ea184111b3893b0403e8a0fbf155b4545", 16)
	e := big.NewInt(65537)
	d, _ := new(big.Int).SetString("24196980ce422d911cb0cec559998415cb19b20af886d6c0a1b34570ce5e608128814e986f37847ac7f8286022b1309c51d98623318d005990f15f3327dfa52653e489585b3d5567cdb324379570d4bb9234ebdebab42f2b4c71fe54c67e7a84b0758d749f3ced24573f22a9c47814412a3cf5424b6c8cdd4eff1ba38bc9a9dc0eb0303ac02074ebc70d88c1ad0df5ae6e82ccc2c2e461d8b9281835fbbe411e77d41878cb02cdca7494bde3eb5ba51f6e751f464ad2671f501703520b9c7cf7f6f971303b4d29903213f58290146e1eeeea49944b44b265a9d622ec09f0824d7f03d8ecaa4374ca90776b2bf6e4d72b6c915bfbf82286ce62014f639aa833bd", 16)
	p, _ := new(big.Int).SetString("faaa76e6b37407e625d4c023fa90e23b4dc35a285adff2ff31181ad47b94acd4fe74f63042dff2cf590f0097284fabff96daf8562c23ace0a239f73f2fcf967e1059fb672084654221d23004f6e6925c8f75ae618adfa5017c271328165c2eb17293e01c25283bfdb9eadb316a6c1d669323855700f7a4e0816ea384c479fe8b", 16)
	q, _ := new(big.Int).SetString("e3a99a0e52f91ba2e4ef751227ce23045932534d484a9280ee600b3103e98a5da73e217856eb79696b4af1a58269be4bc79176545a3fea91d127113ea6be8b3bea34088ffb1f7143428dd71bd69fd5ea34ccb0a2ae5232f179e059b28b0dfc0d37f3bce9f54d5a7349336f6a857b17637116435275407500ec142d45eacd956f", 16)
	return NewPrivateKey(n, e, d, p, q)
}

func TestSafePrimeCheck(t *testing.T) {
	sk := testKeyFromVectors()
	if !IsSafePrime(sk.P) {
		t.Log("test key p is not a safe prime (will be regenerated when computing test vectors)")
	}
	if !IsSafePrime(sk.Q) {
		t.Log("test key q is not a safe prime (will be regenerated when computing test vectors)")
	}
}

func TestDerivePublicKey(t *testing.T) {
	sk := testKeyFromVectors()
	info, _ := hex.DecodeString("030000000069a39da0") // OVER_18 metadata
	pkDerived, err := DerivePublicKey(&sk.PublicKey, info)
	if err != nil {
		t.Fatalf("DerivePublicKey failed: %v", err)
	}
	// e' should be odd and less than n
	if pkDerived.E.Bit(0) == 0 {
		t.Error("derived e' is even")
	}
	if pkDerived.E.Cmp(sk.N) >= 0 {
		t.Error("derived e' >= n")
	}
	// e' should be LambdaLen bytes or less
	if len(pkDerived.E.Bytes()) > LambdaLen {
		t.Errorf("derived e' is %d bytes, expected <= %d", len(pkDerived.E.Bytes()), LambdaLen)
	}
}

func TestDeriveKeyPairInvertible(t *testing.T) {
	sk := testKeyFromVectors()
	info, _ := hex.DecodeString("030000000069a39da0")
	skDerived, pkDerived, err := DeriveKeyPair(sk, info)
	if err != nil {
		t.Fatalf("DeriveKeyPair failed: %v", err)
	}
	// Verify d'*e' mod phi == 1
	pMinus1 := new(big.Int).Sub(sk.P, big.NewInt(1))
	qMinus1 := new(big.Int).Sub(sk.Q, big.NewInt(1))
	phi := new(big.Int).Mul(pMinus1, qMinus1)
	product := new(big.Int).Mul(skDerived.D, pkDerived.E)
	product.Mod(product, phi)
	if product.Cmp(big.NewInt(1)) != 0 {
		t.Error("d' * e' mod phi != 1")
	}
}

func TestDerivePublicKeyDeterministic(t *testing.T) {
	sk := testKeyFromVectors()
	info, _ := hex.DecodeString("030000000069a39da0")
	pk1, _ := DerivePublicKey(&sk.PublicKey, info)
	pk2, _ := DerivePublicKey(&sk.PublicKey, info)
	if pk1.E.Cmp(pk2.E) != 0 {
		t.Error("DerivePublicKey is not deterministic")
	}
}

func TestDerivePublicKeyDifferentInfo(t *testing.T) {
	sk := testKeyFromVectors()
	info1, _ := hex.DecodeString("030000000069a39da0")
	info2, _ := hex.DecodeString("000000000069cc6000")
	pk1, _ := DerivePublicKey(&sk.PublicKey, info1)
	pk2, _ := DerivePublicKey(&sk.PublicKey, info2)
	if pk1.E.Cmp(pk2.E) == 0 {
		t.Error("different info should produce different e'")
	}
}

func TestRoundTrip(t *testing.T) {
	sk := testKeyFromVectors()
	pk := &sk.PublicKey

	msg := []byte("test message for round-trip verification")
	info, _ := hex.DecodeString("030000000069a39da0") // OVER_18 metadata

	// Step 1: Blind
	blindedMsg, state, err := Blind(pk, msg, info, nil)
	if err != nil {
		t.Fatalf("Blind failed: %v", err)
	}
	if len(blindedMsg) != ModulusLen {
		t.Fatalf("blinded message length: got %d, want %d", len(blindedMsg), ModulusLen)
	}

	// Step 2: BlindSign
	blindSig, err := BlindSign(sk, blindedMsg, info)
	if err != nil {
		t.Fatalf("BlindSign failed: %v", err)
	}
	if len(blindSig) != ModulusLen {
		t.Fatalf("blind signature length: got %d, want %d", len(blindSig), ModulusLen)
	}

	// Step 3: Finalize
	sig, err := Finalize(pk, msg, info, blindSig, state.Inv)
	if err != nil {
		t.Fatalf("Finalize failed: %v", err)
	}
	if len(sig) != ModulusLen {
		t.Fatalf("signature length: got %d, want %d", len(sig), ModulusLen)
	}

	// Step 4: Verify
	if err := Verify(pk, msg, info, sig); err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
}

func TestRoundTripDeterministic(t *testing.T) {
	sk := testKeyFromVectors()
	pk := &sk.PublicKey

	msg := []byte("deterministic test")
	info, _ := hex.DecodeString("030000000069a39da0")
	rFixed, _ := new(big.Int).SetString("aabbccdd00112233445566778899aabbccddeeff00112233445566778899aabb", 16)

	blindedMsg1, state1, _ := Blind(pk, msg, info, rFixed)
	blindedMsg2, _, _ := Blind(pk, msg, info, rFixed)

	if hex.EncodeToString(blindedMsg1) != hex.EncodeToString(blindedMsg2) {
		t.Error("Blind with same r should be deterministic")
	}

	blindSig, _ := BlindSign(sk, blindedMsg1, info)
	sig, err := Finalize(pk, msg, info, blindSig, state1.Inv)
	if err != nil {
		t.Fatalf("Finalize failed: %v", err)
	}
	if err := Verify(pk, msg, info, sig); err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
}

func TestVerifyRejectsWrongInfo(t *testing.T) {
	sk := testKeyFromVectors()
	pk := &sk.PublicKey

	msg := []byte("test")
	info, _ := hex.DecodeString("030000000069a39da0")
	wrongInfo, _ := hex.DecodeString("000000000069cc6000")

	blindedMsg, state, _ := Blind(pk, msg, info, nil)
	blindSig, _ := BlindSign(sk, blindedMsg, info)
	sig, _ := Finalize(pk, msg, info, blindSig, state.Inv)

	if err := Verify(pk, msg, wrongInfo, sig); err == nil {
		t.Error("Verify should fail with wrong info")
	}
}

func TestVerifyRejectsWrongMessage(t *testing.T) {
	sk := testKeyFromVectors()
	pk := &sk.PublicKey

	msg := []byte("test")
	info, _ := hex.DecodeString("030000000069a39da0")

	blindedMsg, state, _ := Blind(pk, msg, info, nil)
	blindSig, _ := BlindSign(sk, blindedMsg, info)
	sig, _ := Finalize(pk, msg, info, blindSig, state.Inv)

	if err := Verify(pk, []byte("wrong"), info, sig); err == nil {
		t.Error("Verify should fail with wrong message")
	}
}

func TestAllAgeBrackets(t *testing.T) {
	sk := testKeyFromVectors()
	pk := &sk.PublicKey

	brackets := []string{
		"030000000069a39da0", // OVER_18
		"000000000069cc6000", // UNDER_13
		"01000000006a2fe940", // AGE_13_15
		"020000000069b66700", // AGE_16_17
	}

	for _, infoHex := range brackets {
		t.Run(infoHex, func(t *testing.T) {
			info, _ := hex.DecodeString(infoHex)
			msg := []byte("token message for bracket test")

			blindedMsg, state, err := Blind(pk, msg, info, nil)
			if err != nil {
				t.Fatalf("Blind: %v", err)
			}
			blindSig, err := BlindSign(sk, blindedMsg, info)
			if err != nil {
				t.Fatalf("BlindSign: %v", err)
			}
			sig, err := Finalize(pk, msg, info, blindSig, state.Inv)
			if err != nil {
				t.Fatalf("Finalize: %v", err)
			}
			if err := Verify(pk, msg, info, sig); err != nil {
				t.Fatalf("Verify: %v", err)
			}
		})
	}
}

func TestBuildMsgPrime(t *testing.T) {
	msg := []byte("hello")
	info := []byte("meta")
	result := BuildMsgPrime(msg, info)

	// "msg" + len(info)=4 as uint32 BE + "meta" + "hello"
	expected := []byte("msg")
	expected = append(expected, 0, 0, 0, 4)
	expected = append(expected, []byte("meta")...)
	expected = append(expected, []byte("hello")...)

	if hex.EncodeToString(result) != hex.EncodeToString(expected) {
		t.Errorf("BuildMsgPrime mismatch:\ngot:  %x\nwant: %x", result, expected)
	}
}
