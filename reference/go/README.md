# AAVP Reference Implementation (Go)

Reference implementation of the Anonymous Age Verification Protocol (AAVP) in Go, covering the three protocol roles (Device Agent, Implementor, Verification Gate) and the RSAPBSSA-SHA384-PSSZERO-Deterministic partially blind signature scheme.

## Structure

```
token/       Token binary format (331 bytes): encode, decode, field access
validation/  VG validation logic: clock skew, TTL, field checks
pbrsa/       Partially Blind RSA signatures (draft-amjad-cfrg-partially-blind-rsa)
da/          Device Agent role: prepare, blind, finalize tokens
im/          Implementor role: blind sign, key management, .well-known
vg/          Verification Gate role: full token verification
vectors/     Test vector verification and generation tooling
```

## Requirements

- Go 1.22 or later
- `golang.org/x/crypto` (for HKDF)

## Running tests

```bash
go test ./...
```

To include the slow safe-prime generation test:

```bash
go test ./pbrsa/ -v -timeout 300s
```

## Generating test vectors

The `vectors/generate` tool computes the cryptographic values for `test-vectors/issuance-protocol.json`:

```bash
go run ./vectors/generate/
```

This generates a new RSA-2048 key with safe primes and computes all `TO_BE_COMPUTED` values in the issuance protocol test vectors.

## Test coverage

- **token-encoding.json**: 4 vectors covering all age brackets (encode/decode round-trip)
- **token-validation.json**: 14 vectors covering valid tokens, expiration, clock skew, invalid fields, size checks, and signature failure
- **issuance-protocol.json**: 4 vectors covering the full 6-step issuance flow for all age brackets
