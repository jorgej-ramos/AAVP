---
title: "Anonymous Age Verification Protocol (AAVP)"
abbrev: "AAVP"
docname: draft-ramos-aavp-protocol-00
category: std
ipr: trust200902
area: Security
workgroup: Independent Submission
keyword:
  - age verification
  - privacy
  - blind signatures
  - anonymous credentials
stand_alone: yes
pi:
  toc: yes
  sortrefs: yes
  symrefs: yes
  compact: yes
  subcompact: no

author:
  - name: Jorge Juan Ramos Garnero
    organization: Independent
    email: jramos@example.com

normative:
  RFC2119:
  RFC5280:
  RFC8174:
  RFC8785:
  RFC9162:
  RFC9474:

informative:
  RFC4120:
  RFC6781:
  RFC6962:
  RFC7489:
  RFC7519:
  RFC7696:
  RFC8126:
  RFC8446:
  RFC9053:
  RFC9458:
  RFC9578:
  RFC9614:
  RFC9794:
  I-D.irtf-cfrg-partially-blind-rsa:
    title: "Partially Blind RSA Signatures"
    author:
      - name: Frank Denis
      - name: Frederic Jacobs
      - name: Christopher A. Wood
    date: 2024
    target: https://datatracker.ietf.org/doc/draft-irtf-cfrg-partially-blind-rsa/

--- abstract

This document specifies the Anonymous Age Verification Protocol (AAVP),
a decentralized, privacy-preserving protocol that enables digital
platforms to verify the age bracket of their users without collecting
or transmitting any personally identifiable information. AAVP uses
partially blind signatures to produce ephemeral, unlinkable tokens that
carry only an age bracket signal (one of four predefined ranges) and an
expiration timestamp. The protocol defines three roles -- Device Agent,
Verification Gate, and Implementer -- and a decentralized trust model
where each platform independently decides which Implementers to trust.
No central authority, registry, or coordination point exists. AAVP is
designed as an open standard: anyone can implement it without licenses,
fees, or permissions.

--- middle

# Introduction

Digital platforms increasingly face regulatory requirements to verify
the age of their users, particularly minors. Existing approaches --
such as uploading identity documents, facial recognition, or credit
card verification -- inherently require disclosing personal data that
can be stored, leaked, or misused. These approaches create a
false dilemma between child protection and privacy.

AAVP resolves this tension by introducing a protocol where the age
verification signal (an age bracket, not a precise age) travels from
the user's device to the platform without any entity learning who the
user is. The cryptographic core of the protocol is a partially blind
signature scheme: the Implementer (the entity that signs tokens) sees
the age bracket but cannot link the signature to the user who requested
it. The Verification Gate (the platform endpoint) sees the signed age
bracket but cannot learn who signed it or trace the token back to a
specific device.

## Design Principles

AAVP is built on four inviolable principles. No extension, modification,
or implementation of this protocol MAY compromise any of them:

1. **Privacy by Design.** No personally identifiable information leaves
   the device. This is a mathematical guarantee, not a policy promise.

2. **Decentralization.** No central authority exists. Each platform
   decides independently which Implementers to trust.

3. **Open Standard.** No licenses, fees, or permissions are required.
   Anyone can implement this protocol.

4. **Data Minimalism.** The token carries only an age bracket. Every
   additional field requires rigorous justification and MUST pass the
   data minimalism test defined in {{token-structure}}.

## Conventions and Definitions

{::boilerplate bcp14-tagged}

# Terminology

This section defines the key terms used throughout this document.

Device Agent (DA):
: A software component residing on the user's device that generates,
  manages, and rotates AAVP age tokens. The Device Agent is an abstract
  protocol role, not synonymous with "parental control software". It MAY
  be implemented by parental control applications, operating system
  components, browser extensions, or other conformant software.

Verification Gate (VG):
: A dedicated endpoint operated by a digital platform that validates
  AAVP tokens and establishes age-bracketed sessions. The VG acts as
  a gateway: the age token is presented once during an initial
  handshake, after which the platform operates with its own internal
  session mechanism.

Implementer (IM):
: An organization that develops software conforming to the AAVP
  standard, acting as a provider of Device Agent functionality. The IM
  operates a partially blind signature service and publishes its public
  keys at a well-known endpoint on its own domain.

Age Bracket:
: One of four predefined age ranges used as the minimal age signal in
  AAVP tokens. The canonical values are: UNDER_13 (0x00), AGE_13_15
  (0x01), AGE_16_17 (0x02), and OVER_18 (0x03).

Partially Blind Signature:
: A cryptographic signature scheme where the signer sees a designated
  portion of the message (the public metadata) but cannot see the
  remaining content. In AAVP, the IM sees the age_bracket and
  expires_at fields but the nonce remains blinded.

Token:
: A fixed-size (331-byte) cryptographic structure containing an age
  bracket, an expiration timestamp, a nonce, a token key identifier,
  a token type, and a partially blind signature (authenticator).

Public Metadata:
: The portion of the token visible to the IM during the partially blind
  signing process. In AAVP, this comprises the age_bracket and
  expires_at fields. Public metadata is cryptographically bound to the
  signature via key derivation (HKDF).

Blinded Content:
: The portion of the token hidden from the IM during the signing
  process. In AAVP, this is the nonce field. Cryptographic blinding
  ensures the IM cannot read this value.

Trust Store:
: The list of Implementers accepted by a Verification Gate, along with
  their public keys. Each VG maintains its own trust store
  independently, analogous to a browser's root certificate store.

Token Key ID:
: The SHA-256 hash of the IM's public key. Allows the VG to identify
  which key to use for signature verification without trying all known
  keys.

Token Type:
: A 16-bit unsigned integer identifying the cryptographic scheme used
  to produce the token. Enables cryptographic agility and future
  migration to post-quantum schemes. See {{token-type-registry}}.

Session Credential:
: A self-contained structure issued by the VG after validating an AAVP
  token. Contains exclusively the age_bracket, a session expiration
  timestamp, and the VG's own signature. Does not require server-side
  state.

Clock Skew:
: The synchronization difference between two systems' clocks. AAVP
  defines asymmetric tolerances: CLOCK_SKEW_TOLERANCE_PAST = 300
  seconds, CLOCK_SKEW_TOLERANCE_FUTURE = 60 seconds.

Unlinkability:
: The cryptographic property ensuring that no party can correlate two
  tokens as belonging to the same user or device.

Certificate Transparency (CT):
: An open standard ({{RFC9162}}) requiring Certificate Authorities to
  log all issued certificates in public, auditable logs. Used in AAVP
  to ensure the integrity of TLS certificates on the DA-IM and DA-VG
  channels.

CSPRNG:
: Cryptographically Secure Pseudo-Random Number Generator. The
  operating system's secure randomness source, required for generating
  the 32-byte nonce in each token.

Device Attestation:
: The process by which a device proves the integrity of its execution
  environment to a third party. Includes key attestation and device
  integrity signals. In AAVP, this is an optional mechanism that
  modulates trust, not an access gate. See {{device-attestation}}.

Fail-closed:
: A security policy where loss of verification signal maintains active
  restrictions. In AAVP, applied at the account level: an account
  flagged as belonging to a minor retains restrictions even if the DA
  becomes unavailable. Only a valid OVER_18 credential removes them.

Oblivious HTTP (OHTTP):
: A protocol ({{RFC9458}}) that interposes a relay between client and
  server to conceal the client's identity from the server. In AAVP, an
  optional measure for the DA-IM channel. See {{traffic-analysis}}.

Privacy Partitioning:
: An architectural principle ({{RFC9614}}) of distributing data across
  multiple parties such that no single entity possesses both the user's
  identity and their activity content.

RSAPBSSA:
: RSA Partially Blind Signature Scheme with Appendix. The concrete
  scheme selected by AAVP, based on {{RFC9474}} and
  {{I-D.irtf-cfrg-partially-blind-rsa}}. Uses SHA-384 as the hash
  function.

SPD (Segmentation Policy Declaration):
: A signed JSON document declaring a platform's content segmentation
  policy by age bracket. Served at
  .well-known/aavp-age-policy.json. See {{spd}}.

PTL (Policy Transparency Log):
: An append-only log, inspired by Certificate Transparency
  ({{RFC6962}}), where platform SPDs are recorded. Multiple
  independent logs operate in parallel. See {{ptl}}.

OVP (Open Verification Protocol):
: An open, standardized methodology for verifying that a platform
  complies with its declared SPD. Any party can execute OVP
  verifications. See {{ovp}}.

SAF (Segmentation Accountability Framework):
: The framework comprising SPD, PTL, OVP, and a compliance signal that
  together provide accountability for platform content segmentation
  policies. See {{saf}}.

Conformance:
: The degree to which an implementation meets the requirements of the
  AAVP specification. Three levels are defined: Functional
  (self-assessment), Verified (interoperability testing), and Audited
  (independent third-party audit). See {{conformance}}.

# Protocol Architecture {#architecture}

TODO: Translate PROTOCOL.md Section 1 (Roles, Verification Gate model,
Security Assumptions).

~~~ ascii-art
+--------+           +--------+           +--------+
|        | partially |        |   token   |        |
| Device |   blind   | Imple- |           | Verif. |
| Agent  |   sign    | menter |           |  Gate  |
|  (DA)  +---------->+  (IM)  |  +------->+  (VG)  |
|        |<----------+        |  |        |        |
|        |  signature|        |  |        |        |
+--------+           +--------+  |        +---+----+
     |                            |            |
     +----------------------------+            |
              token presentation          session +
                                        age_bracket
                                               |
                                          +----v----+
                                          |Platform |
                                          +---------+
~~~

## Roles {#roles}

TODO: Translate roles (DA, VG, IM) from PROTOCOL.md Section 1.1.

## Verification Gate Model {#vg-model}

TODO: Translate the gate model from PROTOCOL.md Section 1.2.

## Security Assumptions {#security-assumptions}

TODO: Translate security assumptions (S1-S14) from PROTOCOL.md
Section 1.3, organized in three categories:

- Category A: Assumptions resolved in the specification
- Category B: Partially resolved assumptions
- Category C: Recognized protocol limitations

# Token Structure {#token-structure}

TODO: Translate PROTOCOL.md Section 2 (token fields, binary format,
public metadata vs. blinded content, minimalism test, excluded fields,
nonce generation).

The AAVP token is a fixed-size cryptographic structure of 331 bytes.

~~~ ascii-art
Offset  Size    Field            Visibility
0       2       token_type       Public
2       32      nonce            Blinded (hidden from IM during issuance)
34      32      token_key_id     Public
66      1       age_bracket      Public metadata (0x00-0x03)
67      8       expires_at       Public metadata (uint64 BE, 1h precision)
75      256     authenticator    Partially blind signature (RSAPBSSA-SHA384)
---
Total: 331 bytes (fixed)
~~~

## Binary Format {#binary-format}

TODO: Detailed binary format specification.

## Public Metadata vs. Blinded Content {#metadata-vs-blinded}

TODO: Translate from PROTOCOL.md.

## Data Minimalism Test {#minimalism-test}

TODO: Translate the minimalism test for new fields.

## Explicitly Excluded Fields {#excluded-fields}

TODO: Translate the list of explicitly excluded fields.

## Nonce Generation {#nonce-generation}

TODO: Translate nonce generation requirements (CSPRNG).

# Token Rotation {#token-rotation}

TODO: Translate PROTOCOL.md Section 3 (TTL, clock skew tolerances,
rotation strategy, unlinkability preservation).

# Cryptographic Foundations {#crypto}

TODO: Translate PROTOCOL.md Section 4.

## Partially Blind Signatures {#partially-blind-signatures}

TODO: Translate Section 4.1 (RSAPBSSA-SHA384 scheme, key derivation
via HKDF for public metadata binding).

## Zero-Knowledge Proofs {#zkp}

TODO: Translate Section 4.2 (ZKP candidates for future extensions).

## Fingerprinting Prevention {#fingerprinting-prevention}

TODO: Translate Section 4.3 (TLS requirements, Certificate
Transparency).

## Device Attestation {#device-attestation}

TODO: Translate Section 4.4 (key attestation, integrity signals,
TEE requirements, trust modulation, limitations).

## Traffic Analysis Resistance {#traffic-analysis}

TODO: Translate Section 4.5 (pre-signing, message padding, jitter,
OHTTP).

# Decentralized Trust Model {#trust-model}

TODO: Translate PROTOCOL.md Section 5.

## Trust Without Central Authority {#no-central-authority}

TODO: Translate Section 5.1.

## Trust Mechanisms {#trust-mechanisms}

TODO: Translate Section 5.2 (IM evaluation criteria, publishing keys,
.well-known/aavp-issuer endpoint).

## Service Discovery {#service-discovery}

TODO: Translate Section 5.3 (.well-known/aavp endpoint, DNS _aavp
record, DA-VG and DA-IM discovery).

## Token Type Registry {#token-type-registry}

TODO: Translate Section 5.4 (registry structure, initial values,
IANA registration policy).

## Cryptographic Agility and Algorithm Migration {#crypto-agility}

TODO: Translate Section 5.5 (five-phase migration, degradation attack
protection, post-quantum considerations).

## Analogy with DMARC/DKIM {#dmarc-analogy}

TODO: Translate Section 5.6.

# Operational Flow {#operational-flow}

TODO: Translate PROTOCOL.md Section 6.

## Initial Configuration {#initial-config}

TODO: Translate Section 6.1 (one-time setup).

## Platform Access {#platform-access}

TODO: Translate Section 6.2 (per-session flow).

## Deactivation {#deactivation}

TODO: Translate Section 6.3.

# Session Credential {#session-credential}

TODO: Translate PROTOCOL.md Section 7.

## Design Principles {#session-design-principles}

TODO: Translate Section 7.1.

## Mandatory Token Discard {#token-discard}

TODO: Translate Section 7.2.

## Session Credential Structure {#session-credential-structure}

TODO: Translate Section 7.3.

## Lifetime and Renewal {#session-lifetime}

TODO: Translate Section 7.4.

## Lifecycle {#session-lifecycle}

TODO: Translate Section 7.5.

## Security Properties {#session-security}

TODO: Translate Section 7.6.

## Additive Model and Account-Level Persistence {#additive-model}

TODO: Translate Section 7.7.

## Edge Cases {#session-edge-cases}

TODO: Translate Section 7.8.

## CDN and Edge Architecture Compatibility {#cdn-compatibility}

TODO: Translate Section 7.9.

# Segmentation Accountability Framework {#saf}

TODO: Translate PROTOCOL.md Section 8.

## Motivation and Scope {#saf-motivation}

TODO: Translate Section 8.1.

## Segmentation Policy Declaration {#spd}

TODO: Translate Section 8.2 (SPD structure, JSON schema, categories,
signing, examples).

## Policy Transparency Log {#ptl}

TODO: Translate Section 8.3 (log structure, SPT, monitoring).

## Open Verification Protocol {#ovp}

TODO: Translate Section 8.4 (sampling methodology, automated
verification, reporting).

## Compliance Signal {#compliance-signal}

TODO: Translate Section 8.5.

## Conformance Levels {#saf-conformance-levels}

TODO: Translate Section 8.6.

## Limitations and Residual Risk {#saf-limitations}

TODO: Translate Section 8.7.

# Conformance Requirements {#conformance}

TODO: Translate PROTOCOL.md Section 9.

## Conformance Framework Principles {#conformance-principles}

TODO: Translate Section 9.1.

## Per-Role Conformance Requirements {#per-role-requirements}

TODO: Translate Section 9.2 (DA-01..DA-12, VG-01..VG-14,
IM-01..IM-11). Map Obligatorio/Recomendado/Opcional to
MUST/SHOULD/MAY.

## Verification Methodology {#verification-methodology}

TODO: Translate Section 9.3 (test vectors, interoperability, blind
test).

## Implementation Conformance Levels {#impl-conformance-levels}

TODO: Translate Section 9.4 (Functional, Verified, Audited).

## Continuous Operational Verification {#continuous-verification}

TODO: Translate Section 9.5.

## Framework Sources and References {#framework-references}

TODO: Translate Section 9.6.

# Security Considerations {#security}

TODO: Translate PROTOCOL.md Section 10 + Security Assumptions from
Section 1.3.

This section consolidates the security analysis of AAVP, including
the threat model, known attack vectors, mitigations, and residual
risks.

## Threat Model {#threat-model}

TODO: Translate threat model from PROTOCOL.md Section 10.

## Recognized Limitations {#recognized-limitations}

TODO: Translate limitations from PROTOCOL.md Section 10.

# Privacy Considerations {#privacy}

TODO: Consolidate privacy considerations from PROTOCOL.md Sections 2
(token minimalism), 4.3 (fingerprinting prevention), 7.2 (token
discard), and 7.6 (session security properties).

This section addresses the privacy properties of AAVP, including
unlinkability, data minimalism, fingerprinting resistance, and the
separation of knowledge between protocol roles (privacy partitioning).

# IANA Considerations {#iana}

TODO: Define IANA registry requests from PROTOCOL.md Sections 5.4
and 11.

## AAVP Token Type Registry {#iana-token-type}

This document requests IANA to create a new registry titled
"AAVP Token Types" with the following initial values:

| Value  | Name            | Hash   | Key Size  | Sig Size  | Reference                               | Status |
|--------|-----------------|--------|-----------|-----------|-----------------------------------------|--------|
| 0x0000 | Reserved        | N/A    | N/A       | N/A       | This document                           | N/A    |
| 0x0001 | RSAPBSSA-SHA384 | SHA-384| 2048 bits | 256 bytes | {{RFC9474}}, {{I-D.irtf-cfrg-partially-blind-rsa}} | Active |

New registrations in this registry require "Specification Required"
policy ({{RFC8126}}, Section 4.6).

## AAVP Age Bracket Registry {#iana-age-bracket}

This document requests IANA to create a new registry titled
"AAVP Age Brackets" with the following initial values:

| Value | Name      | Description           |
|-------|-----------|-----------------------|
| 0x00  | UNDER_13  | Under 13 years old    |
| 0x01  | AGE_13_15 | 13 to 15 years old    |
| 0x02  | AGE_16_17 | 16 to 17 years old    |
| 0x03  | OVER_18   | 18 years old or older |

Values 0x04 through 0xFF are reserved for future use.

--- back

# Test Vectors {#test-vectors}

TODO: Include summary of test vectors from test-vectors/ directory.
Reference the repository for complete vector sets.

The AAVP specification includes three sets of test vectors:

- **Token Encoding** (token-encoding.json): Binary encoding test
  cases verifying that token construction produces the expected
  hexadecimal output for each field combination.

- **Token Validation** (token-validation.json): Validation logic test
  cases covering expiration, clock skew, malformed tokens, and
  invalid signatures.

- **Issuance Protocol** (issuance-protocol.json): End-to-end partially
  blind signature flow test cases with real cryptographic values
  computed by the Go reference implementation.

# Formal Verification {#formal-verification}

The security properties of AAVP have been formally verified using the
Tamarin Prover. Three models are maintained:

- **aavp.spthy**: Core protocol model proving unforgeability, nonce
  uniqueness, metadata binding, and executability.

- **aavp-unlinkability.spthy**: Unlinkability model using
  observational equivalence via Tamarin's --diff mode.

- **aavp-saf.spthy**: Segmentation Accountability Framework model
  proving seven properties of the SPD/PTL/OVP system.

# Acknowledgments
{:numbered="false"}

TODO: Add acknowledgments.
