// Package pbrsa implements the Partially Blind RSA Signature scheme
// (RSAPBSSA-SHA384-PSSZERO-Deterministic) as specified in
// draft-amjad-cfrg-partially-blind-rsa, using SHA-384 with salt_length=0.
package pbrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"io"
	"math/big"

	"golang.org/x/crypto/hkdf"
)

const (
	// ModulusLen is the RSA modulus size in bytes (2048 bits).
	ModulusLen = 256
	// LambdaLen is half the modulus length.
	LambdaLen = ModulusLen / 2
	// HKDFExpandLen is the HKDF output length (lambda_len + 16 bytes for bias reduction).
	HKDFExpandLen = LambdaLen + 16
	// HashLen is the output size of SHA-384.
	HashLen = 48
	// SaltLen is the PSS salt length (0 for PSSZERO variant).
	SaltLen = 0
)

var (
	hkdfInputPrefix = []byte("key")
	hkdfLabel       = []byte("PBRSA")
	msgPrimePrefix  = []byte("msg")
)

// PublicKey represents an RSA public key with a big.Int exponent,
// necessary because derived exponents are 128 bytes (too large for int).
type PublicKey struct {
	N *big.Int
	E *big.Int
}

// PrivateKey represents an RSA private key with big.Int exponents.
type PrivateKey struct {
	PublicKey
	D *big.Int
	P *big.Int
	Q *big.Int
}

// BlindingState holds the blinding factor inverse and derived public key,
// needed by the client to finalize the signature.
type BlindingState struct {
	Inv       *big.Int
	R         *big.Int
	PKDerived *PublicKey
}

// FromStdPublicKey converts a standard rsa.PublicKey to our PublicKey type.
func FromStdPublicKey(pk *rsa.PublicKey) *PublicKey {
	return &PublicKey{
		N: new(big.Int).Set(pk.N),
		E: big.NewInt(int64(pk.E)),
	}
}

// FromStdPrivateKey converts a standard rsa.PrivateKey to our PrivateKey type.
func FromStdPrivateKey(sk *rsa.PrivateKey) *PrivateKey {
	return &PrivateKey{
		PublicKey: PublicKey{
			N: new(big.Int).Set(sk.N),
			E: big.NewInt(int64(sk.E)),
		},
		D: new(big.Int).Set(sk.D),
		P: new(big.Int).Set(sk.Primes[0]),
		Q: new(big.Int).Set(sk.Primes[1]),
	}
}

// NewPrivateKey creates a PrivateKey from raw big.Int components.
func NewPrivateKey(n, e, d, p, q *big.Int) *PrivateKey {
	return &PrivateKey{
		PublicKey: PublicKey{N: new(big.Int).Set(n), E: new(big.Int).Set(e)},
		D:         new(big.Int).Set(d),
		P:         new(big.Int).Set(p),
		Q:         new(big.Int).Set(q),
	}
}

// DerivePublicKey derives a public key (n, e') from the master public key
// and the public metadata (info) using HKDF-SHA384.
func DerivePublicKey(pk *PublicKey, info []byte) (*PublicKey, error) {
	// IKM = "key" || info || 0x00
	ikm := make([]byte, 0, len(hkdfInputPrefix)+len(info)+1)
	ikm = append(ikm, hkdfInputPrefix...)
	ikm = append(ikm, info...)
	ikm = append(ikm, 0x00)

	// salt = I2OSP(n, modulus_len)
	salt := I2OSP(pk.N, ModulusLen)

	// HKDF-SHA384(IKM, salt, "PBRSA", L=HKDFExpandLen)
	hkdfReader := hkdf.New(sha512.New384, ikm, salt, hkdfLabel)
	expanded := make([]byte, HKDFExpandLen)
	if _, err := io.ReadFull(hkdfReader, expanded); err != nil {
		return nil, err
	}

	// Take first lambda_len bytes, manipulate bits
	ePrimeBytes := make([]byte, LambdaLen)
	copy(ePrimeBytes, expanded[:LambdaLen])
	ePrimeBytes[0] &= 0x3F           // Clear top 2 bits
	ePrimeBytes[LambdaLen-1] |= 0x01 // Ensure odd

	ePrime := new(big.Int).SetBytes(ePrimeBytes)
	return &PublicKey{N: new(big.Int).Set(pk.N), E: ePrime}, nil
}

// DeriveKeyPair derives a full key pair (sk', pk') from the master private key
// and the public metadata (info). Uses phi = (p-1)*(q-1).
func DeriveKeyPair(sk *PrivateKey, info []byte) (*PrivateKey, *PublicKey, error) {
	pkDerived, err := DerivePublicKey(&sk.PublicKey, info)
	if err != nil {
		return nil, nil, err
	}

	// phi = (p-1) * (q-1)
	pMinus1 := new(big.Int).Sub(sk.P, big.NewInt(1))
	qMinus1 := new(big.Int).Sub(sk.Q, big.NewInt(1))
	phi := new(big.Int).Mul(pMinus1, qMinus1)

	// d' = inverse_mod(e', phi)
	dPrime := new(big.Int).ModInverse(pkDerived.E, phi)
	if dPrime == nil {
		return nil, nil, errors.New("pbrsa: e' is not invertible mod phi; key may not use safe primes")
	}

	skDerived := &PrivateKey{
		PublicKey: *pkDerived,
		D:         dPrime,
		P:         new(big.Int).Set(sk.P),
		Q:         new(big.Int).Set(sk.Q),
	}
	return skDerived, pkDerived, nil
}

// BuildMsgPrime constructs msg_prime = "msg" || len(info) as uint32 BE || info || msg.
func BuildMsgPrime(msg, info []byte) []byte {
	result := make([]byte, 0, 3+4+len(info)+len(msg))
	result = append(result, msgPrimePrefix...)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(info)))
	result = append(result, lenBuf...)
	result = append(result, info...)
	result = append(result, msg...)
	return result
}

// Blind blinds a message for partially blind signing. If rFixed is not nil, it is used
// as the blinding factor (for deterministic test vectors). Otherwise, a random r is generated.
func Blind(pk *PublicKey, msg, info []byte, rFixed *big.Int) (blindedMsg []byte, state *BlindingState, err error) {
	pkDerived, err := DerivePublicKey(pk, info)
	if err != nil {
		return nil, nil, err
	}

	// Construct msg_prime
	msgPrime := BuildMsgPrime(msg, info)

	// EMSA-PSS-ENCODE(msg_prime, bit_len(n) - 1) with salt_len=0
	emBits := pk.N.BitLen() - 1
	em, err := emsaPSSEncode(msgPrime, emBits)
	if err != nil {
		return nil, nil, err
	}

	// m = OS2IP(em)
	m := new(big.Int).SetBytes(em)

	// Generate or use fixed blinding factor r
	var r *big.Int
	if rFixed != nil {
		r = new(big.Int).Set(rFixed)
	} else {
		r, err = randInt(rand.Reader, pk.N)
		if err != nil {
			return nil, nil, err
		}
	}

	// inv = inverse_mod(r, n)
	inv := new(big.Int).ModInverse(r, pk.N)
	if inv == nil {
		return nil, nil, errors.New("pbrsa: blinding factor not invertible")
	}

	// x = r^e' mod n
	x := new(big.Int).Exp(r, pkDerived.E, pk.N)

	// z = m * x mod n
	z := new(big.Int).Mul(m, x)
	z.Mod(z, pk.N)

	// blind_msg = I2OSP(z, modulus_len)
	blindMsg := I2OSP(z, ModulusLen)

	return blindMsg, &BlindingState{Inv: inv, R: r, PKDerived: pkDerived}, nil
}

// BlindSign signs a blinded message using the IM's private key and public metadata.
func BlindSign(sk *PrivateKey, blindedMsg, info []byte) ([]byte, error) {
	if len(blindedMsg) != ModulusLen {
		return nil, errors.New("pbrsa: invalid blinded message length")
	}

	// m = OS2IP(blind_msg)
	m := new(big.Int).SetBytes(blindedMsg)

	// Derive key pair
	skDerived, pkDerived, err := DeriveKeyPair(sk, info)
	if err != nil {
		return nil, err
	}

	// s = m^d' mod n (using CRT for efficiency)
	s := crtExp(m, skDerived)

	// Verification check: s^e' mod n == m
	mCheck := new(big.Int).Exp(s, pkDerived.E, skDerived.N)
	if mCheck.Cmp(m) != 0 {
		return nil, errors.New("pbrsa: signing verification failed")
	}

	return I2OSP(s, ModulusLen), nil
}

// Finalize unblinds the blind signature and verifies it.
func Finalize(pk *PublicKey, msg, info, blindSig []byte, inv *big.Int) ([]byte, error) {
	if len(blindSig) != ModulusLen {
		return nil, errors.New("pbrsa: invalid blind signature length")
	}

	// z = OS2IP(blind_sig)
	z := new(big.Int).SetBytes(blindSig)

	// s = z * inv mod n
	s := new(big.Int).Mul(z, inv)
	s.Mod(s, pk.N)

	// sig = I2OSP(s, modulus_len)
	sig := I2OSP(s, ModulusLen)

	// Verify the unblinded signature
	if err := Verify(pk, msg, info, sig); err != nil {
		return nil, errors.New("pbrsa: finalize verification failed: " + err.Error())
	}

	return sig, nil
}

// Verify verifies a partially blind RSA signature.
func Verify(pk *PublicKey, msg, info, sig []byte) error {
	if len(sig) != ModulusLen {
		return errors.New("pbrsa: invalid signature length")
	}

	pkDerived, err := DerivePublicKey(pk, info)
	if err != nil {
		return err
	}

	msgPrime := BuildMsgPrime(msg, info)

	// s = OS2IP(sig)
	s := new(big.Int).SetBytes(sig)

	// m = s^e' mod n (RSAVP1)
	m := new(big.Int).Exp(s, pkDerived.E, pk.N)

	// em = I2OSP(m, emLen)
	emBits := pk.N.BitLen() - 1
	emLen := (emBits + 7) / 8
	em := I2OSP(m, emLen)

	// EMSA-PSS-VERIFY(msg_prime, em, emBits)
	return emsaPSSVerify(msgPrime, em, emBits)
}

// crtExp computes m^d mod n using CRT.
func crtExp(m *big.Int, sk *PrivateKey) *big.Int {
	pMinus1 := new(big.Int).Sub(sk.P, big.NewInt(1))
	dp := new(big.Int).Mod(sk.D, pMinus1)

	qMinus1 := new(big.Int).Sub(sk.Q, big.NewInt(1))
	dq := new(big.Int).Mod(sk.D, qMinus1)

	qInv := new(big.Int).ModInverse(sk.Q, sk.P)

	sp := new(big.Int).Exp(m, dp, sk.P)
	sq := new(big.Int).Exp(m, dq, sk.Q)

	h := new(big.Int).Sub(sp, sq)
	if h.Sign() < 0 {
		h.Add(h, sk.P)
	}
	h.Mul(h, qInv)
	h.Mod(h, sk.P)

	s := new(big.Int).Mul(h, sk.Q)
	s.Add(s, sq)
	return s
}

// I2OSP converts a non-negative big.Int to a byte string of the given length (I2OSP from RFC 8017).
func I2OSP(n *big.Int, length int) []byte {
	b := n.Bytes()
	if len(b) > length {
		return b[len(b)-length:]
	}
	result := make([]byte, length)
	copy(result[length-len(b):], b)
	return result
}

// randInt returns a uniform random big.Int in [1, max).
func randInt(random io.Reader, max *big.Int) (*big.Int, error) {
	for {
		r, err := rand.Int(random, max)
		if err != nil {
			return nil, err
		}
		if r.Sign() > 0 {
			return r, nil
		}
	}
}

// emsaPSSEncode implements EMSA-PSS-ENCODE from RFC 8017 Section 9.1.1
// with salt_len=0 (deterministic) and SHA-384.
func emsaPSSEncode(msg []byte, emBits int) ([]byte, error) {
	hash := sha512.New384()
	emLen := (emBits + 7) / 8

	hash.Write(msg)
	mHash := hash.Sum(nil)
	hash.Reset()

	if emLen < HashLen+SaltLen+2 {
		return nil, errors.New("encoding error: emLen too short")
	}

	// M' = 0x00{8} || mHash (salt is empty)
	mPrime := make([]byte, 8+HashLen)
	copy(mPrime[8:], mHash)
	hash.Write(mPrime)
	h := hash.Sum(nil)
	hash.Reset()

	// DB = PS || 0x01 (salt is empty)
	dbLen := emLen - HashLen - 1
	db := make([]byte, dbLen)
	db[dbLen-1] = 0x01

	dbMask := mgf1SHA384(h, dbLen)

	maskedDB := make([]byte, dbLen)
	for i := range maskedDB {
		maskedDB[i] = db[i] ^ dbMask[i]
	}

	topBits := uint(8*emLen - emBits)
	if topBits > 0 {
		maskedDB[0] &= 0xFF >> topBits
	}

	em := make([]byte, emLen)
	copy(em, maskedDB)
	copy(em[dbLen:], h)
	em[emLen-1] = 0xBC

	return em, nil
}

// emsaPSSVerify implements EMSA-PSS-VERIFY from RFC 8017 Section 9.1.2
// with salt_len=0 and SHA-384.
func emsaPSSVerify(msg, em []byte, emBits int) error {
	hash := sha512.New384()
	emLen := (emBits + 7) / 8

	hash.Write(msg)
	mHash := hash.Sum(nil)
	hash.Reset()

	if emLen < HashLen+SaltLen+2 {
		return errors.New("inconsistent")
	}

	if em[emLen-1] != 0xBC {
		return errors.New("inconsistent: no 0xBC trailer")
	}

	dbLen := emLen - HashLen - 1
	maskedDB := em[:dbLen]
	h := em[dbLen : dbLen+HashLen]

	topBits := uint(8*emLen - emBits)
	if topBits > 0 {
		mask := byte(0xFF >> topBits)
		if maskedDB[0]&^mask != 0 {
			return errors.New("inconsistent: top bits not zero")
		}
	}

	dbMask := mgf1SHA384(h, dbLen)

	db := make([]byte, dbLen)
	for i := range db {
		db[i] = maskedDB[i] ^ dbMask[i]
	}

	if topBits > 0 {
		db[0] &= 0xFF >> topBits
	}

	// Check DB: zeros || 0x01 (no salt)
	psLen := dbLen - 1
	for i := 0; i < psLen; i++ {
		if db[i] != 0x00 {
			return errors.New("inconsistent: padding not zero")
		}
	}
	if db[psLen] != 0x01 {
		return errors.New("inconsistent: no 0x01 separator")
	}

	// M' = 0x00{8} || mHash (salt is empty)
	mPrime := make([]byte, 8+HashLen)
	copy(mPrime[8:], mHash)
	hash.Write(mPrime)
	hPrime := hash.Sum(nil)

	if !constantTimeEqual(h, hPrime) {
		return errors.New("inconsistent: H != H'")
	}

	return nil
}

// mgf1SHA384 implements MGF1 from RFC 8017 Appendix B.2.1 using SHA-384.
func mgf1SHA384(seed []byte, length int) []byte {
	var result []byte
	counter := make([]byte, 4)
	for i := 0; len(result) < length; i++ {
		binary.BigEndian.PutUint32(counter, uint32(i))
		h := sha512.New384()
		h.Write(seed)
		h.Write(counter)
		result = append(result, h.Sum(nil)...)
	}
	return result[:length]
}

func constantTimeEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := range a {
		v |= a[i] ^ b[i]
	}
	return v == 0
}
