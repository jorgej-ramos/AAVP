package pbrsa

import (
	"crypto/rand"
	"math/big"
)

// GenerateSafePrimeKey generates an RSA key pair where both p and q are safe primes
// (p = 2p'+1, q = 2q'+1 where p' and q' are also prime).
// This is required by draft-amjad-cfrg-partially-blind-rsa.
func GenerateSafePrimeKey(bits int) (*PrivateKey, error) {
	halfBits := bits / 2
	p, err := generateSafePrime(halfBits)
	if err != nil {
		return nil, err
	}
	var q *big.Int
	for {
		q, err = generateSafePrime(halfBits)
		if err != nil {
			return nil, err
		}
		if p.Cmp(q) != 0 {
			break
		}
	}

	n := new(big.Int).Mul(p, q)
	e := big.NewInt(65537)

	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
	qMinus1 := new(big.Int).Sub(q, big.NewInt(1))
	phi := new(big.Int).Mul(pMinus1, qMinus1)

	d := new(big.Int).ModInverse(e, phi)
	if d == nil {
		return nil, errInvertible
	}

	return &PrivateKey{
		PublicKey: PublicKey{N: n, E: e},
		D:         d,
		P:         p,
		Q:         q,
	}, nil
}

var errInvertible = errorf("pbrsa: e is not invertible mod phi")

func errorf(s string) error { return &stringError{s} }

type stringError struct{ s string }

func (e *stringError) Error() string { return e.s }

// generateSafePrime generates a safe prime of the given bit size.
// A safe prime p satisfies p = 2p'+1 where p' is also prime.
func generateSafePrime(bits int) (*big.Int, error) {
	for {
		// Generate a random prime p' of (bits-1) bits
		pPrime, err := rand.Prime(rand.Reader, bits-1)
		if err != nil {
			return nil, err
		}
		// p = 2*p' + 1
		p := new(big.Int).Mul(pPrime, big.NewInt(2))
		p.Add(p, big.NewInt(1))
		// Check if p is also prime
		if p.ProbablyPrime(20) && p.BitLen() == bits {
			return p, nil
		}
	}
}

// IsSafePrime checks if p is a safe prime (p = 2p'+1 where p' is prime).
func IsSafePrime(p *big.Int) bool {
	if !p.ProbablyPrime(20) {
		return false
	}
	// p' = (p-1)/2
	pPrime := new(big.Int).Sub(p, big.NewInt(1))
	pPrime.Rsh(pPrime, 1)
	return pPrime.ProbablyPrime(20)
}
