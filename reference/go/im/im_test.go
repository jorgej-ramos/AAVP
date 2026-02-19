package im

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/aavp-protocol/aavp-go/pbrsa"
)

func testKey() *pbrsa.PrivateKey {
	n, _ := new(big.Int).SetString("deeb38cfa184281eee07142527e4bc39f7d5aba3b9b4acfc5202a1f3e2f02c5f41677ba9fa8c2345cb3ebf7b80c48626b3a6f191909cbcfafd64b5ca92ad5c5c2b32b57660abc2783f5d9093a9a7d1057330fe09a992324dffca3a6469bae4817573e6ad419f926c5074cb1fa6818c95ce376b339056a6a0800bd5b4f7ff2be0981df5c558842d054a61c9c180bfd02abe5cec6e07e05160739fd64734186b8f01b6489f4d33177100bb5c12c629f3f4f4bcb3a39abd776834491101ff554aabddc9abea1820b6039c4aca86141a9d2295e113d8203a976ee0cc210b7912a917b4ab4366b6c4c787672c5e9d6f29a97ea184111b3893b0403e8a0fbf155b4545", 16)
	e := big.NewInt(65537)
	d, _ := new(big.Int).SetString("24196980ce422d911cb0cec559998415cb19b20af886d6c0a1b34570ce5e608128814e986f37847ac7f8286022b1309c51d98623318d005990f15f3327dfa52653e489585b3d5567cdb324379570d4bb9234ebdebab42f2b4c71fe54c67e7a84b0758d749f3ced24573f22a9c47814412a3cf5424b6c8cdd4eff1ba38bc9a9dc0eb0303ac02074ebc70d88c1ad0df5ae6e82ccc2c2e461d8b9281835fbbe411e77d41878cb02cdca7494bde3eb5ba51f6e751f464ad2671f501703520b9c7cf7f6f971303b4d29903213f58290146e1eeeea49944b44b265a9d622ec09f0824d7f03d8ecaa4374ca90776b2bf6e4d72b6c915bfbf82286ce62014f639aa833bd", 16)
	p, _ := new(big.Int).SetString("faaa76e6b37407e625d4c023fa90e23b4dc35a285adff2ff31181ad47b94acd4fe74f63042dff2cf590f0097284fabff96daf8562c23ace0a239f73f2fcf967e1059fb672084654221d23004f6e6925c8f75ae618adfa5017c271328165c2eb17293e01c25283bfdb9eadb316a6c1d669323855700f7a4e0816ea384c479fe8b", 16)
	q, _ := new(big.Int).SetString("e3a99a0e52f91ba2e4ef751227ce23045932534d484a9280ee600b3103e98a5da73e217856eb79696b4af1a58269be4bc79176545a3fea91d127113ea6be8b3bea34088ffb1f7143428dd71bd69fd5ea34ccb0a2ae5232f179e059b28b0dfc0d37f3bce9f54d5a7349336f6a857b17637116435275407500ec142d45eacd956f", 16)
	return pbrsa.NewPrivateKey(n, e, d, p, q)
}

func TestTokenKeyID(t *testing.T) {
	sk := testKey()
	spkiDER, _ := hex.DecodeString("30820122300d06092a864886f70d01010105000382010f003082010a0282010100deeb38cfa184281eee07142527e4bc39f7d5aba3b9b4acfc5202a1f3e2f02c5f41677ba9fa8c2345cb3ebf7b80c48626b3a6f191909cbcfafd64b5ca92ad5c5c2b32b57660abc2783f5d9093a9a7d1057330fe09a992324dffca3a6469bae4817573e6ad419f926c5074cb1fa6818c95ce376b339056a6a0800bd5b4f7ff2be0981df5c558842d054a61c9c180bfd02abe5cec6e07e05160739fd64734186b8f01b6489f4d33177100bb5c12c629f3f4f4bcb3a39abd776834491101ff554aabddc9abea1820b6039c4aca86141a9d2295e113d8203a976ee0cc210b7912a917b4ab4366b6c4c787672c5e9d6f29a97ea184111b3893b0403e8a0fbf155b45450203010001")

	im := NewImplementor(sk, spkiDER, "test-im.example")
	keyID := im.TokenKeyID()

	expectedKeyID := sha256.Sum256(spkiDER)
	if keyID != expectedKeyID {
		t.Error("TokenKeyID mismatch")
	}

	expectedHex := "fffea9ba9efa735080cf1af734625994ed056c1c4f94a8d82f4676a017ab2c7c"
	if hex.EncodeToString(keyID[:]) != expectedHex {
		t.Errorf("TokenKeyID hex: got %s, want %s", hex.EncodeToString(keyID[:]), expectedHex)
	}
}

func TestWellKnownResponse(t *testing.T) {
	sk := testKey()
	spkiDER, _ := hex.DecodeString("30820122300d06092a864886f70d01010105000382010f003082010a0282010100deeb38cfa184281eee07142527e4bc39f7d5aba3b9b4acfc5202a1f3e2f02c5f41677ba9fa8c2345cb3ebf7b80c48626b3a6f191909cbcfafd64b5ca92ad5c5c2b32b57660abc2783f5d9093a9a7d1057330fe09a992324dffca3a6469bae4817573e6ad419f926c5074cb1fa6818c95ce376b339056a6a0800bd5b4f7ff2be0981df5c558842d054a61c9c180bfd02abe5cec6e07e05160739fd64734186b8f01b6489f4d33177100bb5c12c629f3f4f4bcb3a39abd776834491101ff554aabddc9abea1820b6039c4aca86141a9d2295e113d8203a976ee0cc210b7912a917b4ab4366b6c4c787672c5e9d6f29a97ea184111b3893b0403e8a0fbf155b45450203010001")

	imInst := NewImplementor(sk, spkiDER, "test-im.example")

	notBefore, _ := time.Parse(time.RFC3339, "2026-01-01T00:00:00Z")
	notAfter, _ := time.Parse(time.RFC3339, "2026-06-30T00:00:00Z")

	resp := imInst.WellKnownResponse(notBefore, notAfter)

	if resp.Issuer != "test-im.example" {
		t.Errorf("issuer: got %q", resp.Issuer)
	}
	if len(resp.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(resp.Keys))
	}
	if resp.Keys[0].TokenKeyID != "__6pup76c1CAzxr3NGJZlO0FbBxPlKjYL0Z2oBerLHw" {
		t.Errorf("token_key_id: got %q", resp.Keys[0].TokenKeyID)
	}
	if resp.Keys[0].TokenType != 1 {
		t.Errorf("token_type: got %d", resp.Keys[0].TokenType)
	}
}

func TestSign(t *testing.T) {
	sk := testKey()
	spkiDER, _ := hex.DecodeString("30820122300d06092a864886f70d01010105000382010f003082010a0282010100deeb38cfa184281eee07142527e4bc39f7d5aba3b9b4acfc5202a1f3e2f02c5f41677ba9fa8c2345cb3ebf7b80c48626b3a6f191909cbcfafd64b5ca92ad5c5c2b32b57660abc2783f5d9093a9a7d1057330fe09a992324dffca3a6469bae4817573e6ad419f926c5074cb1fa6818c95ce376b339056a6a0800bd5b4f7ff2be0981df5c558842d054a61c9c180bfd02abe5cec6e07e05160739fd64734186b8f01b6489f4d33177100bb5c12c629f3f4f4bcb3a39abd776834491101ff554aabddc9abea1820b6039c4aca86141a9d2295e113d8203a976ee0cc210b7912a917b4ab4366b6c4c787672c5e9d6f29a97ea184111b3893b0403e8a0fbf155b45450203010001")

	imInst := NewImplementor(sk, spkiDER, "test-im.example")

	// Blind a test message
	pk := &sk.PublicKey
	msg := []byte("test message for IM signing")
	info, _ := hex.DecodeString("030000000069a39da0")

	blindedMsg, state, err := pbrsa.Blind(pk, msg, info, nil)
	if err != nil {
		t.Fatalf("Blind: %v", err)
	}

	// IM signs
	blindSig, err := imInst.Sign(blindedMsg, info)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	// Finalize and verify
	sig, err := pbrsa.Finalize(pk, msg, info, blindSig, state.Inv)
	if err != nil {
		t.Fatalf("Finalize: %v", err)
	}

	if err := pbrsa.Verify(pk, msg, info, sig); err != nil {
		t.Fatalf("Verify: %v", err)
	}
}
