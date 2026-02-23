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
can be stored, leaked, or misused. These approaches create a false
dilemma between child protection and privacy.

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
   data minimalism test defined in {{minimalism-test}}.

## Conventions and Definitions

{::boilerplate bcp14-tagged}

# Terminology

This section defines the key terms used throughout this document.

Device Agent (DA):
: A software component residing on the user's device that generates,
  manages, and rotates AAVP age tokens. The Device Agent is an abstract
  protocol role, not synonymous with "parental control software". It
  MAY be implemented by parental control applications, operating system
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

AAVP defines three roles with distinct responsibilities. The design
ensures that none needs to blindly trust the others: cryptographic
verifiability replaces institutional trust.

~~~ ascii-art
+--------+            +---------+           +---------+
|        | blind sign |         |   token   |         |
| Device +----------->| Imple-  |           | Verif.  |
| Agent  |<-----------+ menter  |  +------->+  Gate   |
|  (DA)  | signature  |  (IM)   |  |        |  (VG)   |
+---+----+            +---------+  |        +----+----+
    |                              |             |
    +------------------------------+        session +
           token presentation             age_bracket
                                                |
                                          +-----v-----+
                                          | Platform   |
                                          +-----------+
~~~

## Roles {#roles}

### Device Agent (DA)

The Device Agent is an abstract protocol role: a software component
residing on the minor's device that is responsible for generating,
managing, and rotating age tokens.

The Device Agent is NOT synonymous with "parental control software".
It is a protocol role that MAY be implemented by various vehicles:

| Implementation Vehicle | Example |
|------------------------|---------|
| Parental control system | Qustodio, Bark, carrier software |
| Native OS component | Built-in module in iOS, Android, Windows |
| Browser extension | Extension conforming to the specification |
| Device firmware | Routers with integrated parental controls |

The separation between the role (Device Agent) and its implementation
vehicle is deliberate: it allows the ecosystem to evolve without
modifying the protocol.

Responsibilities of the DA:

- Generate local key pairs in secure storage (Secure Enclave, TPM,
  StrongBox).
- Generate ephemeral tokens with the configured age bracket.
- Obtain a partially blind signature from the Implementer: the public
  metadata (age_bracket, expires_at) is visible to the IM, but the
  nonce remains blinded.
- Present signed tokens to the Verification Gate.
- Rotate tokens before expiration.
- Protect the age bracket configuration via parental PIN or equivalent
  OS-level mechanism.

### Verification Gate (VG)

A dedicated endpoint of the digital platform acting as a gateway to
the service. It validates the AAVP token and establishes an internal
session with the age bracket mark.

Responsibilities of the VG:

- Expose the discovery endpoint .well-known/aavp and optionally the
  DNS record _aavp as specified in {{service-discovery}}.
- Validate the cryptographic signature of the token against the public
  keys of accepted Implementers.
- Verify the token TTL.
- Extract the age bracket and establish an internal session.
- Reject expired, malformed, or untrusted tokens.

### Implementer (IM)

An organization that develops software acting as a Device Agent,
conforming to the AAVP standard.

Responsibilities of the IM:

- Publish its public key on its own domain via the
  .well-known/aavp-issuer endpoint as specified in {{im-keys}}.
- Maintain auditable code (preferably open source).
- Provide a partially blind signature service to the Device Agent.
- Comply with the open specification.

## Verification Gate Model {#vg-model}

A naive approach would send the age credential in every HTTP request,
continuously exposing it to potential interception. AAVP adopts a
different model: the gateway.

The age token travels only once per session, during a dedicated initial
handshake. After that, the platform operates with its own session
system.

The age token never coexists with regular application traffic. It
is a separate channel, a one-time handshake. After the handshake, the
information "this user is a minor" becomes an internal platform flag,
completely decoupled from the original token.

~~~ ascii-art
  User       DA        VG       Platform
   |          |         |          |
   |--open--->|         |          |
   |          |<-AAVP-->|          |
   |          | supported          |
   |          |         |          |
   |          |--token->|          |
   |          |  (once) |          |
   |          |<--OK----|          |
   |          |         |          |
   |          |         |--session->|
   |          |         |+age_bracket
   |<---------+---------+--content-|
   |          |         |          |
   |   ... normal session ...      |
   |          |         |          |
   |          |--new--->|          |
   |          | token   |          |
   |          |<--OK----|          |
   |          |         |--renew-->|
~~~

Advantages of the gateway model:

- **Reduced attack surface:** the age token travels only once per
  session, not in every request.
- **Context separation:** age information never coexists with
  application data traffic.
- **Compatibility:** platforms already manage sessions; AAVP only adds
  a prior step.
- **Minimal MITM window:** intercepting the initial handshake requires
  compromising TLS within a very brief window. All protocol channels
  (DA-VG and DA-IM) MUST use TLS 1.3 or higher.

## Security Assumptions {#security-assumptions}

Every cryptographic protocol rests on explicit and implicit
assumptions. AAVP makes them explicit so that implementers, auditors,
and regulators can evaluate the guarantees and their limits. For a
detailed analysis of each assumption, see the companion document
SECURITY-ANALYSIS.md.

### Category A -- Assumptions Resolved in the Specification

These assumptions have direct coverage in the technical specification.

| ID | Assumption | Level | Reference |
|----|-----------|-------|-----------|
| S1 | TLS 1.3 + Certificate Transparency protects DA-VG and DA-IM channels | MUST | {{fingerprinting-prevention}}, {{im-keys}} |
| S3 | Partially blind signatures prevent linking token to user | MUST | {{partially-blind-signatures}} |
| S4 | Token rotation prevents longitudinal tracking | MUST | {{token-rotation}} |
| S5 | Post-handshake sessions are secure | MUST | {{session-credential}} |
| S9 | The DA-IM channel is confidential and has integrity | MUST | {{fingerprinting-prevention}} |
| S10 | Clocks are reasonably synchronized (defined tolerance) | MUST | {{token-rotation}} |
| S11 | Each IM publishes its keys on its own domain | MUST | {{im-keys}} |
| S14 | IM revocation occurs through natural expiration and bilateral decision | MUST | {{trust-mechanisms}} |

Each assumption is backed by cryptographic or protocol mechanisms
defined in the referenced sections. The guarantees are verifiable by
any implementer.

### Category B -- Partially Resolved Assumptions

These assumptions have partial mitigations in the specification but
retain residual risk.

| ID | Assumption | Level | Status |
|----|-----------|-------|--------|
| S2 | Secure hardware protects DA keys | SHOULD (MUST when available) | Optional key attestation ({{device-attestation}}) |
| S6 | Open source auditing prevents malicious IMs | SHOULD | No runtime verification |
| S7 | Parental PIN or OS protection prevents deactivation by the minor | Depends on implementation vehicle | Partially mitigated by account-level persistence ({{additive-model}}) |
| S8 | The device is not compromised (root/jailbreak) | Explicit; not guaranteeable by the protocol | Optional device attestation ({{device-attestation}}) |
| S12 | Platforms correctly implement segmentation | SHOULD (with public verification) | SAF ({{saf}}) mitigates with SPD + PTL + OVP |

**S2 -- Secure hardware.** DA keys SHOULD be generated in secure
hardware when the device supports it ({{da-keys-hardware}}). Key
attestation ({{key-attestation}}) allows the IM to differentiate
between hardware-backed and software-only keys. Residual risk: devices
without TEE and attacks on specific TEE implementations.

**S6 -- Open source auditing.** The standard recommends auditable code
({{auditable-code}}), but no runtime verification exists that the
published code is what actually executes. The mitigation depends on
reproducible builds and periodic audits, which are outside the
protocol's scope.

**S7 -- Deactivation protection.** Effectiveness depends on the DA
implementation vehicle. Account-level persistence ({{additive-model}})
partially mitigates: even if the minor uninstalls the DA, account
restrictions persist. Residual risk: devices where the DA has no
OS-level protection.

**S8 -- Uncompromised device.** On a rooted or jailbroken device, all
DA guarantees can be defeated. {{device-attestation}} defines optional
device attestation (key attestation + integrity signals) as partial
mitigation. Explicitly documented as a protocol limitation
({{attestation-limitations}}).

**S12 -- Correct segmentation.** The SAF ({{saf}}) defines
accountability infrastructure: signed SPD, transparency logs (PTL),
and open verification (OVP). Residual risk: dynamic content and UGC
make exhaustive verification difficult.

### Category C -- Recognized Protocol Limitations

| ID | Assumption | Justification |
|----|-----------|---------------|
| S13 | The minor does not have access to a second device without DA | Inherent limitation. AAVP protects devices where it is present. The additive model ({{additive-model}}) partially mitigates: the account retains restrictions even when accessed from another device with DA. A device without DA generates no AAVP signal. |

These limitations are inherent to the model and cannot be resolved
without compromising the protocol's principles.

# Token Structure {#token-structure}

The token is a fixed-size cryptographic structure of 331 bytes designed
to be minimal. Each field has a specific justification and passes the
data minimalism test of the protocol.

| Field | Content | Purpose |
|-------|---------|---------|
| token_type | uint16, identifies the cryptographic scheme | Enables cryptographic agility and post-quantum migration. |
| nonce | 32 cryptographically secure random bytes | Prevents reuse and ensures uniqueness of each token. Blinded during issuance. |
| token_key_id | SHA-256 of the IM's public key (32 bytes) | Allows the VG to identify which key to use for signature verification. |
| age_bracket | Enumeration: UNDER_13 (0x00), AGE_13_15 (0x01), AGE_16_17 (0x02), OVER_18 (0x03) | Age bracket signal. Public metadata of the partially blind signature. |
| expires_at | uint64 big-endian, Unix timestamp with 1-hour precision | Validity window. Public metadata. Coarse precision groups tokens temporally. |
| authenticator | Partially blind signature RSAPBSSA-SHA384 (256 bytes) | Proves the token originates from a legitimate IM without linking to the user. |

## Binary Format {#binary-format}

The token has a fixed binary format of 331 bytes, with no separators
or encoding metadata. Canonicalization is implicit in the format:
fields are concatenated in the specified order with deterministic
offsets.

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

All conformant implementations MUST produce tokens of exactly 331
bytes. A token of a different size is invalid.

## Public Metadata vs. Blinded Content {#metadata-vs-blinded}

The AAVP token uses partially blind signatures. This implies a
distinction between two types of content within the token:

- **Public metadata** (age_bracket, expires_at): visible to the IM
  during the signing process. The IM uses them to derive a specific
  signing key via HKDF. They are part of the visible contract between
  the DA and the IM.
- **Blinded content** (nonce): hidden from the IM during issuance. Only
  the DA and the VG know its value. Cryptographic blinding ensures the
  IM cannot read it.

The IM knows the age bracket of the token it signs, but it cannot link
that information to the identity of the user who requested it. Within
the same bracket, all tokens are indistinguishable to the IM. This
preserves unlinkability: the age bracket is not personal data, it is
the minimal signal that the protocol needs to transmit.

This architecture is acceptable because:

1. The age bracket is precisely the signal the protocol transmits. It
   is not additional information.
2. The IM does not obtain anything that the VG does not also obtain
   when verifying the token.
3. The IM can act as a second barrier against age_bracket spoofing by
   verifying coherence with the DA's configuration.

## Data Minimalism Test {#minimalism-test}

The fields token_type and token_key_id are additions relative to
earlier versions of the specification. Both pass the data minimalism
test:

- **token_type**: necessary for cryptographic agility (post-quantum
  migration). Its value is identical for all tokens of the same scheme.
  It does not enable individual fingerprinting.
- **token_key_id**: necessary for the VG to identify the verification
  key without trying all known keys. Derived from the IM's public key
  (not the user's). Identical for all tokens from the same IM.

Any proposal to add a field to the AAVP token MUST pass this test:

1. **Necessity:** Is it strictly necessary for protocol operation?
2. **Minimalism:** Can the same goal be achieved without this field?
3. **Fingerprinting:** Can this field, alone or combined with others,
   be used to identify or track the user?
4. **Unlinkability:** Does it compromise the impossibility of
   correlating two tokens from the same user?

If the answer to 3 or 4 is "yes" or "possibly", the field MUST be
rejected.

## Explicitly Excluded Fields {#excluded-fields}

The token MUST NOT contain:

- User identity
- Device identifier
- IP address
- Geographic location
- Software version
- Operating system
- Issuance timestamp (issued_at): removed because freshness is managed
  via coarse expires_at, and an issuance timestamp with jitter is an
  unnecessary fingerprinting surface
- Any other data that enables correlation or tracking

Each additional piece of data is a potential fingerprinting vector and
MUST be rigorously justified before inclusion in future versions.

## Nonce Generation {#nonce-generation}

The 32-byte nonce MUST be generated using a Cryptographically Secure
Pseudo-Random Number Generator (CSPRNG) provided by the operating
system. The use of weak randomness sources compromises the token's
uniqueness guarantees and can facilitate prediction attacks.

Required APIs by platform:

| Platform | Required API | Alternative |
|----------|-------------|-------------|
| iOS / macOS | SecRandomCopyBytes (Security.framework) | CryptoKit (SymmetricKey(size:)) |
| Android | java.security.SecureRandom (backed by /dev/urandom) | -- |
| Linux | getrandom(2) or getentropy() | /dev/urandom |
| Windows | BCryptGenRandom (BCrypt) | -- |
| Web (browser) | crypto.getRandomValues() (Web Crypto API) | -- |

The following randomness sources MUST NOT be used, as they produce
non-conformant tokens:

- Math.random() (JavaScript) or non-cryptographic equivalents in other
  languages.
- rand(), random() (standard C/C++).
- Derivation from timestamps, MAC addresses, or device identifiers.
- PRNGs initialized with predictable or low-entropy seeds.

The 256-bit nonce space is sufficient to prevent collisions (negligible
probability even after 10^18 tokens). However, a weak PRNG with 32
bits of effective entropy reduces the space to 2^32 possible values,
making collisions probable after approximately 65,000 tokens.

Conformance testing for DA implementations MUST verify:

1. That the source code exclusively uses the OS CSPRNG APIs listed
   above.
2. That a sample of 10,000 generated nonces passes the NIST SP 800-22
   statistical test suite or equivalent.
3. That consecutive nonces share no detectable prefixes, suffixes, or
   patterns.

# Token Rotation {#token-rotation}

Even without personal data, a static token could become a persistent
pseudo-identifier if reused. Therefore, AAVP implements mandatory
rotation:

- **Maximum time to live (TTL):** Each token has a validity defined by
  expires_at, with a RECOMMENDED range between 1 and 4 hours. The VG
  validates expires_at against its own clock.

- **Coarse precision of expires_at:** The expires_at value MUST be
  rounded to the nearest complete hour. This means all tokens issued
  within the same hour share the same expiration value, which increases
  the anonymity set and hinders temporal correlation.

- **Clock skew tolerance:** The VG MUST apply an asymmetric tolerance
  when validating expires_at:
  - **Expired tokens:** The VG MUST accept tokens whose expires_at
    has passed by no more than 300 seconds (5 minutes). This
    accommodates imperfect clock synchronization on mobile devices,
    consistent with the de facto tolerance of Kerberos ({{RFC4120}})
    and JWT recommendations ({{RFC7519}}).
  - **Future tokens:** The VG MUST reject tokens whose expires_at
    exceeds the VG's current time by more than the maximum allowed TTL
    (4 hours) plus 60 seconds. An excessively future expires_at
    indicates a manipulated clock or a fabricated token.
  - **Canonical values:** Conformant implementations MUST use
    CLOCK_SKEW_TOLERANCE_PAST = 300 and CLOCK_SKEW_TOLERANCE_FUTURE =
    60 as default values. VGs MAY adjust these values for their
    environment but SHOULD NOT exceed the recommended maximums.

- **Proactive rotation:** The Device Agent MAY generate a new token
  before expiration to maintain session continuity.

- **Unlinkability:** Two consecutive tokens from the same device MUST
  NOT be correlatable. Each token is cryptographically independent of
  the previous one.

# Cryptographic Foundations {#crypto}

## Partially Blind Signatures {#partially-blind-signatures}

The central mechanism of AAVP for decoupling user identity from the
age signal is the use of partially blind signatures, an evolution of
the blind signatures proposed by David Chaum in 1983.

The selected scheme is RSAPBSSA-SHA384 (RSA Partially Blind Signature
Scheme with Appendix), based on {{RFC9474}} and
{{I-D.irtf-cfrg-partially-blind-rsa}}. This scheme allows the IM to
see the public metadata (age_bracket, expires_at) while the nonce
remains blinded.

~~~ ascii-art
  DA                              IM
  |                                |
  |-- generate nonce, build token  |
  |-- define public metadata:      |
  |   age_bracket, expires_at      |
  |-- blind message with factor r  |
  |                                |
  |--- blinded msg + metadata ---->|
  |                                |-- sees age_bracket, expires_at
  |                                |-- does NOT see nonce
  |                                |-- verify metadata coherence
  |                                |-- derive key:
  |                                |   sk' = DeriveKeyPair(sk, metadata)
  |                                |-- sign:
  |                                |   blind_sig = BlindSign(sk', blinded_msg)
  |<-------- blind_sig ------------|
  |                                |
  |-- unblind:                     |
  |   auth = Finalize(pk, token,   |
  |           metadata, blind_sig, |
  |           r)                   |
  |   auth is valid signature      |
  |   over complete token          |
~~~

**Key derivation per metadata:** The IM has a single master key pair
(sk, pk). For each combination of public metadata (age_bracket,
expires_at), a derived key pair (sk', pk') is automatically computed
via HKDF. The VG, which knows the master public key and the token's
metadata, performs the same derivation to verify. This
cryptographically binds the metadata to the signature without revealing
the blinded content.

**Result:** The Implementer knows the age bracket but cannot link a
signed token to the DA that requested it. Within the same bracket, all
tokens are indistinguishable to the IM. The bracket is not personal
data: it is the signal that the protocol transmits.

### DA-IM Channel Security

The partially blind signature guarantees that the IM cannot link the
finalized token to the signing request (blindness). However, the
transport channel can leak metadata that compromises this property:

- **Minimum requirement:** The DA-IM channel MUST use TLS 1.3 or
  higher. The DA MUST verify the IM's certificate chain against the
  operating system's root certificates. Certificate integrity MUST be
  backed by Certificate Transparency ({{RFC9162}}).
- **Network metadata:** Even with TLS, the IM observes the DA's IP
  address, its TLS fingerprint (JA3/JA4), and the temporal patterns of
  requests. These leaks are inherent to TCP/IP transport and are
  mitigated by the traffic analysis resistance measures defined in
  {{traffic-analysis}}.
- **Traffic analysis resistance:** Implementations SHOULD apply the
  measures defined in {{traffic-analysis}}: pre-signing with temporal
  decoupling ({{pre-signing}}), message padding ({{message-padding}}),
  presentation jitter ({{presentation-jitter}}), and optionally
  Oblivious HTTP ({{RFC9458}}) for the DA-IM channel ({{ohttp}}).

The DA-VG channel has the same TLS requirements. The difference is that
on the DA-VG channel the token is already finalized and contains no
data linkable to the DA, so network metadata leakage has lower privacy
impact.

### Candidate Schemes

- **RSAPBSSA-SHA384** ({{RFC9474}} +
  {{I-D.irtf-cfrg-partially-blind-rsa}}) -- Selected scheme for AAVP.
- Blind BLS Signatures -- Future alternative due to reduced signature
  size (48 bytes). No published RFC.
- ZKP (Bulletproofs) -- Complement for initial age verification against
  an official document.

## Zero-Knowledge Proofs {#zkp}

As an alternative or complement to blind signatures, AAVP contemplates
the use of Zero-Knowledge Proofs (ZKP).

A ZKP allows proving a statement -- "my age is within bracket X" --
without revealing any additional data. This is particularly useful in
scenarios where the initial age verification is performed against an
official document: the ZKP would prove that the date of birth meets
the bracket criterion without exposing the date, name, or any other
field of the document.

Candidate schemes:

- zk-SNARKs (Groth16, PLONK)
- zk-STARKs (no trusted setup)
- Bulletproofs (for range proofs over age)

## Fingerprinting Prevention {#fingerprinting-prevention}

Each field of the token is designed to minimize information that could
be used to identify or track the user:

| Measure | Affected field | Purpose |
|---------|---------------|---------|
| Coarse precision | expires_at | Rounding to the hour eliminates temporal correlation. All tokens issued in the same hour share the same value. |
| Cryptographic nonce | nonce | Generated without derivation from device identifiers |
| Minimal metadata | age_bracket, expires_at | Only two public metadata fields. age_bracket partitions the anonymity set into 4 groups (inherent to the protocol's purpose). Hourly precision of expires_at groups all tokens from the same hour. |
| Frequent rotation | expires_at | Short-lived tokens prevent longitudinal tracking |
| Fixed size | (entire token) | All tokens are exactly 331 bytes |

## Device Attestation {#device-attestation}

The cryptographic guarantees of AAVP depend on the Device Agent
operating in an intact environment. If the device is compromised
(root/jailbreak), an attacker can bypass the DA's protections. This
section defines device integrity mechanisms as optional trust signals
that modulate -- but never condition -- access to the protocol.

### DA Keys in Secure Hardware {#da-keys-hardware}

The Device Agent's cryptographic keys SHOULD be generated and stored in
secure hardware when the device supports it. Private keys are
non-exportable: all cryptographic operations occur within the secure
environment.

| Platform | Secure Hardware | Key Generation API |
|----------|----------------|-------------------|
| iOS / macOS | Secure Enclave | SecKey with kSecAttrTokenIDSecureEnclave |
| Android | StrongBox / TEE | Android Keystore with setIsStrongBoxBacked() |
| Windows | TPM 2.0 | CNG with NCRYPT_PROVIDER_HANDLE (TPM) |
| Linux | TPM 2.0 | tpm2-tss / PKCS#11 |

When secure hardware is not available, keys MUST be generated in
software-based secure storage (Keychain, Android Keystore without
StrongBox). The resulting trust level is lower, but the DA remains
functional.

### Key Attestation {#key-attestation}

Key attestation is the central device integrity mechanism in AAVP. It
allows the DA to prove to the IM that its cryptographic keys reside in
genuine secure hardware, not in an emulated environment.

Key attestation flow:

1. The DA generates a key pair inside the TEE.
2. The TEE produces an attestation certificate chain whose root is the
   hardware manufacturer's CA.
3. The DA presents the chain to the IM along with the registration
   request.
4. The IM verifies the chain against known attestation roots.
5. If the chain is valid: the key is marked as hardware-backed (high
   trust). If not: the key is accepted as software-only (base trust).

Key attestation mechanisms by platform:

| Platform | Mechanism | Trust Chain |
|----------|-----------|-------------|
| Android (API 24+) | Key Attestation (android.security.keystore) | Google Hardware Attestation Root |
| iOS / macOS | App Attest (DCAppAttestService) | Apple Attestation CA |
| Windows | TPM 2.0 Key Attestation | TPM manufacturer CA |

Privacy restrictions on the attestation chain:

| Data | IM MAY inspect | IM MUST NOT use |
|------|:--------------:|:---------------:|
| Hardware security level (TEE, StrongBox, software) | Yes | -- |
| Key properties (non-exportable, restricted usage) | Yes | -- |
| Device model | -- | No (fingerprinting) |
| Device identifier | -- | No (linkability) |
| OS version | -- | No (fingerprinting) |

The attestation chain does not reveal the user's identity. The IM MUST
only inspect the hardware security level and key properties.

Key attestation is a trust signal, not an access gate. An IM that
categorically rejects requests without attestation is creating a
barrier that violates the Open Standard principle. The IM MUST accept
software-only keys at base trust level.

### Device Integrity Signals

In addition to key attestation, the DA MAY verify the integrity of its
execution environment locally:

- **Root/jailbreak detection:** The DA MAY refuse to operate if it
  detects the device is compromised. Root detection is an arms race
  between detectors and evaders, so it cannot be considered an
  absolute guarantee.
- **Device integrity APIs** (Play Integrity on Android, App Attest on
  iOS): the DA MAY use them as local self-checks. The results of these
  checks MUST NOT be transmitted to the IM or VG.

Three design decisions resolve the tension with decentralization:

1. Key attestation operates at the DA-IM level (bilateral
   relationship). It introduces no new central authority: each IM
   independently decides which attestation roots to accept.
2. Device integrity signals are local DA self-checks. They MUST NOT be
   transmitted to third parties.
3. There is no attestation gate. Attestation modulates trust level,
   not protocol access.

### DA Key Rotation

DA keys have a limited lifetime. Periodic rotation forces
re-attestation and bounds the exploitation window in case of
compromise.

| Parameter | Recommended Value | Justification |
|-----------|-------------------|---------------|
| Rotation period | 7 days | Limits the exploitation window for compromised keys |
| Minimum overlap | >= maximum token TTL (4 hours) | Tokens signed with the previous key remain verifiable |
| Maximum period without rotation | 30 days | Forces re-attestation even on intermittently used devices |

During rotation, the DA generates a new key pair, obtains a new
attestation chain, and registers it with the IM. The previous key
remains valid during the overlap period so that already-issued tokens
are not prematurely invalidated.

### Attestation Limitations {#attestation-limitations}

AAVP assumes the device's operating system is intact. On a rooted or
jailbroken device, all DA guarantees can be defeated: an attacker can
intercept keys, modify the age bracket in memory, or inject fabricated
tokens.

Key attestation ({{key-attestation}}) and device integrity signals
offer partial detection. Weekly key rotation limits the exploitation
window. However, a fully compromised device can evade these
mitigations.

This is a recognized protocol limitation, not a flaw. AAVP explicitly
documents this assumption ({{security-assumptions}}, S8) and offers
the viable technical mitigations without compromising the
decentralization principle.

## Traffic Analysis Resistance {#traffic-analysis}

AAVP operates over the public network and is subject to traffic
analysis by observers with privileged position (ISP, state entities).
This section defines the mitigations incorporated into the protocol to
hinder session correlation.

The AAVP architecture implements privacy partitioning ({{RFC9614}}):
the IM knows the age bracket but not the destination platform; the VG
knows the platform but not the user's identity. No single entity has
access to both pieces of data simultaneously.

### Pre-signing and Temporal Decoupling {#pre-signing}

The DA SHOULD temporally decouple the acquisition of signed tokens from
their presentation to the VG. The pattern "request to IM followed by
request to VG" within a short interval (~100-500 ms) is a correlatable
signal.

Requirements:

- The DA SHOULD obtain signed tokens at times independent of platform
  access (background refresh). Tokens are stored locally in the
  device's keychain.
- The DA MAY request multiple tokens in a single interaction with the
  IM (batch issuance), reducing contact frequency.
- Between obtaining a token and presenting it to the VG, a minimum
  interval SHOULD exist. The DA SHOULD NOT contact the IM and the VG
  within the same 5-minute temporal window.

This pattern is consistent with the batch issuance of Privacy Pass
(draft-ietf-privacypass-batched-tokens) and Apple Private Access Tokens
implementation.

### Message Padding {#message-padding}

AAVP protocol requests and responses SHOULD be indistinguishable in
size from standard HTTP/API calls.

Requirements:

- The HTTP request and response bodies of the AAVP handshake SHOULD be
  padded to a multiple of 2048 bytes (2 KiB).
- The padding consists of random bytes added in a "padding" JSON field
  or as additional bytes in the HTTP body, ignored by the receiver.
- The receiver MUST ignore the "padding" field without error.

The AAVP token size is fixed (331 bytes), which facilitates making
padded exchanges indistinguishable from typical API responses.

### Presentation Jitter {#presentation-jitter}

The DA SHOULD introduce a random delay before presenting a token to
the VG.

Requirements:

- Uniform jitter between 0 and 300 seconds before the first
  presentation of a token to a new VG.
- The jitter applies only when the DA does not have a valid session
  credential for the platform. Subsequent renewals do not require
  additional jitter (the temporal decoupling from {{pre-signing}} is
  sufficient).

### Oblivious HTTP for Maximum Privacy {#ohttp}

Implementations seeking to minimize network metadata leakage SHOULD use
Oblivious HTTP ({{RFC9458}}) for the DA-IM channel.

The OHTTP architecture interposes a relay between the DA and the IM:

- The relay observes the DA's IP but cannot read the request content
  (encrypted with the IM's public key).
- The IM reads the request but only observes the relay's IP.
- No single entity simultaneously observes the user's identity (IP)
  and the request content (blind signature).

OHTTP relay operators are available in production (Cloudflare Privacy
Gateway, Fastly). The privacy properties of OHTTP have been formally
verified.

Note that OHTTP protects against the IM and against network observers
between the DA and the relay. It does not protect against a relay that
colludes with the IM. Non-collusion between relay and IM is an explicit
trust requirement ({{RFC9614}}, Section 4).

# Decentralized Trust Model {#trust-model}

## Trust Without Central Authority {#no-central-authority}

AAVP explicitly rejects the centralized Certificate Authority model.
Centralization of certification creates:

- **Perverse incentives:** the central entity acquires veto power.
- **Priority target:** for political pressure and attacks.
- **Single point of failure:** whose compromise invalidates the entire
  system.

AAVP adopts a distributed trust model, inspired by DMARC/DKIM for
email authentication.

~~~ ascii-art
Centralized model (rejected):

  IM1 --\                 /--> Platform 1
  IM2 ---+-> Central CA --+--> Platform 2
  IM3 --/                 \--> Platform 3

AAVP model (adopted) -- each platform decides whom to trust:

  IM1 ---+--> Platform 1
  IM1 ---+--> Platform 2
  IM2 ---+--> Platform 1
  IM2 ---+--> Platform 2
  IM3 ---+--> Platform 1
~~~

## Trust Mechanisms {#trust-mechanisms}

### Open and Verifiable Standard {#auditable-code}

Any organization can implement AAVP. Its tokens are cryptographically
verifiable by any platform that also implements the standard. No
permission from any third party is needed. Trust comes from
mathematical verifiability, not institutional authorization.

The standard strongly recommends -- and regulation may require -- that
Device Agent implementations be open source or, at minimum, auditable
by independent third parties. This is analogous to Certificate
Transparency logs: the community can verify that the software complies
with the specification.

### IM Key Publication {#im-keys}

Each Implementer publishes its cryptographic material on its own
domain. No centralized registry exists: the IM is the authoritative
source for its own keys.

Primary endpoint: https://\[IM-domain\]/.well-known/aavp-issuer

The endpoint MUST be served over TLS 1.3, with certificate integrity
backed by Certificate Transparency ({{RFC9162}}).

JSON response (application/json):

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| issuer | string (hostname) | Yes | FQDN of the IM. MUST match the domain serving the endpoint. |
| aavp_version | string | Yes | Protocol version supported. Format MAJOR.MINOR. |
| signing_endpoint | string (URI) | Yes | HTTPS URI of the partially blind signature service. Same domain or subdomain of issuer. |
| keys | array | Yes | Active signing keys (includes keys in rotation). |
| keys[].token_key_id | string | Yes | SHA-256 of the public key, encoded in base64url without padding (43 characters). Matches the token_key_id field in the token. |
| keys[].token_type | uint16 | Yes | Cryptographic scheme identifier. Matches token_type in the token (see registry in {{token-type-registry}}). |
| keys[].public_key | string | Yes | Public key in SPKI DER format (SubjectPublicKeyInfo, {{RFC5280}}), encoded in base64url without padding. |
| keys[].not_before | string (ISO 8601) | Yes | Start of the validity period. |
| keys[].not_after | string (ISO 8601) | Yes | End of the validity period. not_after - not_before MUST be <= 180 days. |

HTTP requirements:

- HTTPS mandatory (TLS 1.3, verifiable CT).
- Cache-Control: public, max-age=86400 (24 hours).
- Access-Control-Allow-Origin: * (for browser-based DAs).
- The client MUST verify that issuer matches the domain from which the
  document was obtained.

Complementary DNS: TXT record _aavp-keys.\[IM-domain\]:

    v=aavp1; url=https://im.example/.well-known/aavp-issuer

This endpoint is NOT an approval authority. Any organization can
publish keys on its domain. Trust does not come from being published,
but from the independent decision of each VG to accept that
Implementer.

### IM Key Lifecycle

IM signing keys have a limited lifetime. This reduces the exposure
window if a key is compromised and eliminates the need for a
centralized revocation mechanism.

- **Maximum recommended lifetime:** 6 months (180 days).
  Implementations MUST NOT accept keys with a validity period exceeding
  this.
- **Rotation:** When the IM generates a new key, it MUST publish both
  (old and new) simultaneously. The overlap period MUST be at least
  equal to the maximum token TTL (4 hours) so that tokens signed with
  the previous key remain verifiable until expiration. An overlap of
  at least 24 hours is RECOMMENDED to give VGs time to refresh their
  cache.
- **Natural expiration:** A key past its expiration date is no longer
  valid for signature verification. VGs MUST reject tokens whose
  token_key_id corresponds to an expired key.
- **No centralized revocation:** No central mechanism exists to revoke
  an IM key. If an IM detects that its key has been compromised, it
  removes the key from its endpoint. Effective revocation is bilateral:
  each VG manages its own trust store and can remove an IM at any time.

### VG Trust Store Management

Each VG maintains a local trust store: a list of accepted Implementers
along with their public keys. The decision to trust an IM is
independent for each VG, with no mediation by any central authority.

- **Key acquisition:** The VG obtains public keys directly from the
  IM's domain over TLS 1.3. The VG MUST verify the TLS certificate
  chain and the presence in Certificate Transparency logs before
  accepting the cryptographic material.
- **Cache and refresh:** The VG SHOULD cache the keys of accepted IMs.
  Refresh SHOULD be periodic (RECOMMENDED: at least every 24 hours) to
  detect key rotations and possible withdrawals.
- **Trust revocation:** The VG MAY remove an IM from its trust store
  at any time, without coordination with other VGs. This is analogous
  to how a browser can unilaterally stop trusting a Certificate
  Authority.
- **Discovery of new IMs:** VGs MUST NOT automatically trust unknown
  IMs. Incorporating a new IM into the trust store is a deliberate
  decision following a reputational process.

### Reputation-Based Trust

Platforms individually decide which Implementers to trust, just as
browsers decide which CAs to trust for TLS. There is no centralized
decision, but multiple independent decisions that tend to converge.

## Service Discovery {#service-discovery}

Platforms supporting AAVP announce it via a discovery endpoint. The
Device Agent queries this endpoint to determine whether the platform
accepts AAVP tokens, which Implementers it recognizes, and to which
URL to send the handshake.

### VG Discovery Endpoint

URI: https://\[platform-domain\]/.well-known/aavp

JSON response (application/json):

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| aavp_version | string | Yes | Protocol version supported. Format MAJOR.MINOR. |
| vg_endpoint | string (URI) | Yes | HTTPS URI of the handshake endpoint where the DA presents tokens. Same domain or subdomain. |
| accepted_ims | array of objects | Yes | Implementers accepted by this VG. |
| accepted_ims[].domain | string (hostname) | Yes | FQDN of the IM. The DA uses this domain to locate .well-known/aavp-issuer. |
| accepted_ims[].token_key_ids | array of strings | No | Currently accepted token_key_ids (base64url). If omitted, all active keys from the IM are accepted. |
| accepted_token_types | array of uint16 | Yes | Accepted token_type values (see registry in {{token-type-registry}}). |
| age_policy | string (URI) | No | URI of the Segmentation Policy Declaration (SPD). If omitted, the platform does not publish a verifiable segmentation policy. See {{saf}}. |

HTTP requirements:

- HTTPS mandatory.
- Cache-Control: public, max-age=3600 (1 hour).
- Access-Control-Allow-Origin: *.
- The DA MUST validate that vg_endpoint is on the same domain or
  subdomain as the .well-known host.

Response codes:

| Code | Meaning | DA Behavior |
|------|---------|-------------|
| 200 | AAVP supported | Parse and proceed |
| 404 | No AAVP support | Fallback to DNS; if DNS also fails, no support |
| 429 | Rate limiting | Retry with exponential backoff |
| 5xx | Server error | Retry once; use cache if available; fallback to DNS |

### DNS as Complementary Mechanism

TXT record _aavp.\[platform-domain\]:

    v=aavp1; e=https://platform.example/aavp/verify; im=im1.example,im2.example

| Key | Required | Description |
|-----|:--------:|-------------|
| v | Yes | Version tag. Fixed value aavp1. |
| e | Yes | Handshake endpoint URI (equivalent to vg_endpoint). |
| im | No | Comma-separated list of accepted IM domains. |

DNS is informational: if the DA can reach .well-known/aavp, the JSON
document takes precedence.

### Discovery Priority

The DA follows this priority chain to detect AAVP support:

1. **Local cache** of known platforms (if entry has not expired).
2. **.well-known/aavp** over HTTPS (primary mechanism).
3. **DNS _aavp TXT** as fallback.

If no mechanism responds, the DA concludes that the platform does not
support AAVP. A negative result SHOULD be cached for 1 hour.

### Endpoint Connection Flow

1. The DA obtains .well-known/aavp from the platform.
2. Reads accepted_ims and verifies whether its IM is in the list.
3. If it has a valid pre-signed token (obtained previously via the
   temporal decoupling strategy of {{pre-signing}}) with a key accepted
   by the VG (token_key_ids), it uses it directly without contacting
   the IM.
4. If it needs a new token, queries .well-known/aavp-issuer from the
   IM to obtain active keys.
5. Requests a partially blind signature from the IM's signing_endpoint.
6. Presents the token to the platform's vg_endpoint.

## Token Type Registry {#token-type-registry}

The token_type field of the token and the accepted_token_types and
keys[].token_type fields of the discovery endpoints share the same
value space. This registry is the protocol's central mechanism for
cryptographic agility: each value uniquely identifies a partially blind
signature scheme with all its parameters fixed.

Design principle: Each token_type completely defines the cryptographic
scheme, hash function, key size, signature size, and the metadata-based
key derivation procedure. There is no optionality within a single
token_type. This avoids the class of algorithm confusion attacks
documented in JWT/JWS (RFC 7518), where parameter flexibility within
a single algorithm identifier has produced vulnerabilities such as
acceptance of the "none" algorithm or confusion between symmetric and
asymmetric keys.

Registry values:

| Value | Scheme | Hash | Key Size | Sig Size | Reference | Status |
|-------|--------|------|:--------:|:--------:|-----------|--------|
| 0x0000 | Reserved | -- | -- | -- | -- | Do not use |
| 0x0001 | RSAPBSSA-SHA384 | SHA-384 | 2048 bits | 256 bytes | {{RFC9474}}, {{I-D.irtf-cfrg-partially-blind-rsa}} | Active |
| 0x0002-0x00FF | Unassigned | -- | -- | -- | -- | Reserved for RSA-based schemes |
| 0x0100-0x01FF | Unassigned | -- | -- | -- | -- | Reserved for elliptic curve schemes |
| 0x0200-0x02FF | Unassigned | -- | -- | -- | -- | Reserved for post-quantum schemes |
| 0x0300-0xFFFE | Unassigned | -- | -- | -- | -- | Reserved for future schemes |
| 0xFFFF | Reserved | -- | -- | -- | -- | Do not use |

Registry notes:

- **Immutability:** Once published, a token_type value MUST NOT be
  reassigned to a different scheme or have its parameters modified.
  If a scheme needs different parameters, a new value MUST be assigned.
- **Range segmentation:** The reserved ranges facilitate logical
  organization of future schemes without imposing technical
  restrictions.
- **Token size:** The total token size MAY vary between schemes (since
  signature sizes differ). The VG determines the expected size from
  the token_type before parsing the rest of the token.
- **IANA registration:** When AAVP is formalized as an Internet-Draft,
  token_type values will be registered in a dedicated IANA registry
  with "Specification Required" policy ({{RFC8126}}, Section 4.6),
  analogous to the Token Types registry of Privacy Pass ({{RFC9578}},
  Section 8.1).

## Cryptographic Agility and Algorithm Migration {#crypto-agility}

Cryptographic agility is a protocol's ability to transition between
cryptographic schemes without disrupting service or compromising
security during the transition. {{RFC7696}} (BCP 201) establishes the
general principles for designing this capability in Internet protocols.
This section defines how AAVP applies them.

### Design Principles

The AAVP cryptographic agility model is grounded in three principles
derived from accumulated experience with prior protocols:

**1. Out-of-band selection, not in-line negotiation.**

Unlike TLS ({{RFC8446}}), which negotiates the cipher suite during the
handshake, AAVP follows the Privacy Pass model ({{RFC9578}}): the
token_type is agreed out of band via the discovery endpoints. The DA
queries accepted_token_types from the VG and keys[].token_type from
the IM, selects a compatible value, and generates the token. There is
no negotiation during token presentation.

This model eliminates the attack surface of negotiation protocols
(downgrade attacks such as DROWN in TLS or the reinstallation of weak
cipher suites in IKEv1).

**2. Integral and unalterable type.**

The token_type field travels within the token and is covered by the
authenticator signature. The VG reads it before attempting any
cryptographic operation. This prevents type confusion attacks where an
adversary modifies the scheme identifier to force verification with a
weaker algorithm.

**3. Scheme fully specified by type.**

Each token_type fixes all cryptographic parameters: signature
algorithm, hash function, key size, salt length, metadata-based key
derivation method. There is no optionality within a type. This design
follows the philosophy of TLS 1.3 cipher suites ({{RFC8446}}, Section
4.2) and COSE code points ({{RFC9053}}), where each identifier is
self-contained.

### Token Type Selection Procedure (DA)

When the Device Agent needs to generate a token for a platform, it
MUST execute the following procedure:

1. Obtain the accepted_token_types list from the platform's VG (from
   local cache or .well-known/aavp).
2. Obtain the keys[].token_type list from the configured IM (from
   local cache or .well-known/aavp-issuer).
3. Compute the intersection of both lists.
4. If the intersection is empty, the DA cannot generate a token for
   this VG-IM combination. Log the error and, if the DA supports
   multiple IMs, try another IM.
5. If the intersection contains a single value, use it.
6. If the intersection contains multiple values, select the one with
   the highest numerical value whose status in the registry is
   "Active". The numerical order reflects the chronological progression
   of registered schemes: higher values correspond to more recent
   schemes. This rule is deterministic and does not require complex
   prioritization logic in the DA.

### Migration Between Cryptographic Schemes

Migration from an old scheme to a new one follows a five-phase
procedure inspired by the documented practices of the most widely
deployed Internet protocols:

- **DNSSEC** ({{RFC6781}}): defines a key rollover procedure with
  overlap periods.
- **WebPKI** (CA/Browser Forum): SHA-1 deprecation followed a
  four-phase schedule between 2014 and 2017.
- **TLS** ({{RFC8446}}): the transition from TLS 1.2 to 1.3 maintained
  compatibility via the supported_versions extension mechanism.

Migration phases:

**Phase 1 -- Publication.** The IM generates a key pair for the new
scheme and adds it to the keys[] array of .well-known/aavp-issuer with
a future not_before. VGs do not yet include the new type in
accepted_token_types. No impact on service.

**Phase 2 -- Adoption.** VGs add the new token_type to their
accepted_token_types list. DAs that support the new scheme and find a
non-empty intersection begin generating tokens with the new type. DAs
that do not yet support the new scheme continue generating tokens with
the old scheme, which is still accepted.

**Phase 3 -- Overlap.** Both schemes coexist. This phase MUST last long
enough for the vast majority of DAs to update. The duration depends on
the ecosystem's update velocity, but MUST consider: the update cycle of
major DA implementation vehicles (OS components: weeks to months;
parental control applications: days to weeks); discovery endpoint cache
TTL (24 hours); natural IM key rotation (maximum not_after - not_before
<= 180 days).

**Phase 4 -- Deprecation.** VGs remove the old type from
accepted_token_types. DAs still using the old scheme will no longer
find an intersection for updated VGs and MUST migrate. IMs mark old
keys as expired (not_after in the past) but keep them published during
an additional period so VGs can verify tokens issued before deprecation
that have not yet expired.

**Phase 5 -- Retirement.** IMs remove old keys from their endpoint. The
type is marked as "Deprecated" in the registry. No protocol actor
generates or accepts tokens with that type.

### Per-Actor Responsibilities During Migration

| Actor | Responsibility |
|-------|---------------|
| **IM** | Publish new scheme keys in advance (Phase 1). Maintain keys for both schemes during overlap (Phases 2-3). Retire old keys only after their not_after has expired and tokens issued with them have expired (Phases 4-5). |
| **VG** | Add the new type to accepted_token_types (Phase 2). Maintain both types during overlap (Phase 3). Retire the old type (Phase 4). Verify each token with the scheme indicated by its token_type; MUST NOT attempt multiple schemes. |
| **DA** | Update implementation to support the new scheme. Select type per the intersection rule ({{crypto-agility}}). MUST NOT generate tokens with types the VG does not accept. |

### Degradation Attack Protection

A degradation (downgrade) attack forces participants to use a weaker
scheme than what both support. AAVP mitigates this vector through:

- **token_type covered by signature.** The field is inside the signed
  message; modifying it invalidates the authenticator. An adversary
  cannot change a token's type in transit.
- **No negotiation.** There is no round-trip protocol where an
  intermediary can manipulate each party's announced capabilities.
- **VG cache over HTTPS.** Discovery endpoints are served over TLS
  with 24-hour Cache-Control. A network adversary cannot modify the
  accepted_token_types list without compromising the TLS connection.
- **Degradation detection at VG.** If a VG observes an anomalous
  increase in tokens with an old type after having published a new type
  in accepted_token_types, it may indicate an attack or distribution
  problem. The VG MAY log this anomaly for operational investigation.

### Post-Quantum Migration Considerations

RSA-based partially blind signature schemes are threatened long-term
by advances in quantum computing (Shor's algorithm). AAVP's design
contemplates this migration:

**Current research status.** The NIST PQC published standards (FIPS 203
ML-KEM, FIPS 204 ML-DSA, FIPS 205 SLH-DSA) cover key encapsulation
and standard digital signatures, but do NOT cover blind or partially
blind signatures. Lattice-based blind signatures are an active research
area; no standardized scheme exists yet.

**Transition mechanism.** When a post-quantum partially blind signature
scheme is standardized:

1. A new token_type value will be assigned in the 0x0200-0x02FF range.
2. The token size will increase (post-quantum signatures are
   significantly larger; ML-DSA-65 produces 3,309-byte signatures vs.
   256 bytes for RSA-2048). The VG will determine the expected size
   from the token_type.
3. Migration will follow the five phases defined above.

**Hybrid schemes.** During a transition phase, it MAY be desirable to
issue tokens offering both classical and post-quantum security.
{{RFC9794}} documents algorithm combination strategies. A hybrid scheme
would be registered as an independent token_type with its own signature
size defined as the concatenation of both signatures. The VG would
verify both signatures to accept the token.

## Analogy with DMARC/DKIM {#dmarc-analogy}

| Aspect | DMARC/DKIM | AAVP |
|--------|-----------|------|
| Central authority | None | None |
| Who can issue | Any mail server | Any Implementer |
| Who decides to trust | Each receiver (Gmail, Outlook...) | Each digital platform |
| Basis of trust | Standard compliance + track record | Standard compliance + auditing |
| Consequence of fraud | Emails rejected / spam | Tokens rejected by platforms |

# Operational Flow {#operational-flow}

## Initial Configuration {#initial-config}

This step is performed by parents or guardians. It is the only moment
that requires conscious human intervention.

1. Parents activate the AAVP functionality on the minor's device. The
   vehicle may be a parental control system, a native OS
   configuration, or other conformant software.
2. The software acting as Device Agent generates a local key pair in
   the device's secure storage (Secure Enclave on iOS, StrongBox/TEE
   on Android, TPM on Windows/Linux).
3. The DA establishes a one-time connection with the Implementer's
   signing service to obtain the partially blind signature capability.
4. The corresponding age bracket is configured for the minor.

## Platform Access {#platform-access}

This process is completely transparent to the user:

1. The user opens the application or accesses the website.
2. The Device Agent detects that the platform supports AAVP following
   the discovery chain defined in {{service-discovery}}: local cache,
   .well-known/aavp over HTTPS, and DNS _aavp TXT as fallback.
3. The DA generates an ephemeral token, blinds the nonce, sends the
   blinded message along with the public metadata (age_bracket,
   expires_at) to the Implementer for partially blind signing, unblinds
   the signature, and presents the token to the Verification Gate.
4. The VG validates the signature against the public keys of accepted
   Implementers, verifies the TTL, and extracts the age bracket.
5. The platform establishes a session according to the session
   credential model described in {{session-credential}}.
6. Content is filtered according to the platform's policy for that
   bracket.
7. When the token expires, the DA generates a new one and the VG
   renews the session. The process is transparent.

## Deactivation {#deactivation}

If the software acting as Device Agent is deactivated during an active
session, it stops issuing tokens. On the next revalidation, the session
cannot be renewed and transitions to an "unverified" state.

The policy for sessions where the DA disappears is each platform's
decision. However, the protocol establishes guidelines in
{{additive-model}}: if the platform has previously registered a minor
bracket for that account, restrictions SHOULD be maintained until a
valid OVER_18 credential is presented.

# Session Credential {#session-credential}

Once the VG validates an AAVP token, it needs a mechanism to maintain
the age bracket signal during the user's session without retaining the
original token or storing server-side state. This section defines the
recommended model: a self-contained, ephemeral, and unlinkable session
credential.

## Design Principles {#session-design-principles}

| AAVP Principle | Application to Session |
|---------------|----------------------|
| **Privacy by Design** | The credential contains exclusively age_bracket. The complete AAVP token is discarded after validation. |
| **Decentralization** | Each VG generates and validates its own session credentials. No centralized session service exists. |
| **Open Standard** | The model is part of the open specification. Each platform chooses its concrete format. |
| **Data Minimalism** | Only three fields: age_bracket, session_expires_at, and VG signature. No additional data. |

## Mandatory Token Discard {#token-discard}

After validating an AAVP token, the VG:

1. MUST extract exclusively age_bracket from the token.
2. MUST discard the complete token. The VG MUST NOT store, log, or
   retransmit any field of the AAVP token after validation. This
   includes nonce, authenticator, token_key_id, and expires_at.
3. MUST NOT generate or store token derivatives (hashes, digests) that
   could act as pseudo-identifiers.

Token discard is the most important privacy property of session
management. A VG that stores complete tokens is inadvertently creating
a repository of pseudo-identifiers that, in case of a security breach,
could compromise user privacy.

## Session Credential Structure {#session-credential-structure}

The VG issues a self-contained session credential containing
exclusively:

| Field | Type | Purpose |
|-------|------|---------|
| age_bracket | Enumeration (UNDER_13, AGE_13_15, AGE_16_17, OVER_18) | Age bracket extracted from the validated token |
| session_expires_at | Timestamp | Expiration moment of the session credential |
| vg_signature | VG signature | Guarantees integrity and authenticity of the credential |

The credential is self-contained: it includes all information necessary
for its validation. The VG does not need to maintain server-side state
(session store). Verification is performed by checking vg_signature and
that session_expires_at has not passed.

The concrete credential format (binary structure, signature algorithm,
transport mechanism to client) is an implementation decision of each
platform. The specification defines the mandatory fields and the
properties the credential MUST satisfy, not the exact format.

## Lifetime and Renewal {#session-lifetime}

The TTL of the session credential MUST be strictly less than or equal
to the TTL of the AAVP token that originated it.

| Property | Recommended Value | Justification |
|----------|-------------------|---------------|
| Credential TTL | 15-30 minutes | Limits the exploitation window for session hijacking |
| Relationship with token TTL | session_expires_at <= expires_at of the token | The session MUST NOT outlive the token that generated it |
| Renewal frequency | On credential expiration | The DA presents a new AAVP token |

Credential renewal implies a complete cycle:

1. The session credential expires.
2. The DA generates a new AAVP token, cryptographically independent of
   the previous one.
3. The DA presents the new token to the VG.
4. The VG validates the new token, discards the token, and issues a new
   session credential.
5. The new credential is not linkable to the previous one.

This model guarantees that each renewal produces a cryptographically
independent session: the VG cannot correlate two consecutive
credentials from the same user.

## Lifecycle {#session-lifecycle}

~~~ ascii-art
                +-------------+
                | Unverified  |<-----------+
                +------+------+            |
                       |                   |
            DA presents token        no new token
                       |             (DA deactivated)
                       v                   |
                +-----------+              |
            +-->| Validating|--invalid---->+
            |   +-----+-----+
            |         | valid
            |         v
            |   +----------+
            |   |  Active  |---user logs out--->+
            |   +----+-----+                    |
            |        |                          |
            |  session_expires_at               |
            |        |                          |
            |        v                          v
            |   +---------+            +--------+-----+
            +---| Expired |            |  Unverified  |
                +---------+            +--------------+
~~~

## Security Properties {#session-security}

The session credential satisfies the following properties:

- **Self-contained.** Does not require server-side state. Validatable
  with the VG's verification key. Compatible with CDN/edge
  architectures where validation can occur at the network edge.
- **Ephemeral.** Short TTL (15-30 minutes) that limits the exploitation
  window for credential theft (session hijacking).
- **Unlinkable.** Two consecutive credentials from the same user are
  independent. Each comes from a different AAVP token, and the
  credential contains no identifiers enabling correlation.
- **Purely additive.** AAVP only restricts when there is an active
  minor bracket signal. The absence of AAVP signal (user without DA)
  implies no restriction: the experience is identical to what would
  exist without AAVP. Restrictions apply exclusively to accounts that
  have received a minor bracket.
- **Account-level persistent.** When a platform receives a minor
  bracket for an account, that restriction SHOULD persist at the
  account level. Removing restrictions requires a valid OVER_18
  credential; the mere absence of AAVP signal is not sufficient.
- **Minimalist.** Contains exclusively age_bracket, session_expires_at,
  and vg_signature. Any additional data violates the data minimalism
  principle and MUST NOT be included.

## Additive Model and Account-Level Persistence {#additive-model}

AAVP is a purely additive protocol: it only adds restrictions when a
DA actively sends a minor bracket signal. A user without a DA is not
affected at all.

### Fundamental Principle

| Situation | User Experience |
|-----------|----------------|
| No DA (no AAVP handshake ever occurred) | No restrictions. Experience identical to a platform without AAVP. |
| DA present, minor bracket | Restrictions according to the received bracket. |
| DA present, OVER_18 bracket | No restrictions. Cryptographic proof of age available. |

### Account-Level Persistence

When a platform receives a minor bracket signal for a user account, it
SHOULD persist that signal as an internal account flag. This flag
survives session credential expiration and DA deactivation:

| Situation | Recommended Behavior |
|-----------|---------------------|
| Account marked as minor, active session credential | Restrictions per credential bracket |
| Account marked as minor, expired credential without renewal | Restrictions persist (account still marked) |
| Account marked as minor, DA uninstalled | Restrictions persist (account still marked) |
| Restriction removal | Only via valid OVER_18 credential |

Account-level persistence is the primary defense against evasion by DA
deactivation. If a minor uninstalls the DA software, platforms where
their account was already marked as minor continue applying
restrictions. To remove them, a valid OVER_18 credential is required --
which needs a DA configured with that bracket and an Implementer's
signature.

### Signaling

The platform SHOULD signal to the user when restrictions are active and
offer a mechanism to present an OVER_18 credential if the user
considers the restrictions inapplicable.

## Edge Cases {#session-edge-cases}

- **User clears cookies during an active session.** The platform loses
  the session credential, but the minor bracket flag at the account
  level persists. On the next interaction, if the DA is present, a new
  AAVP handshake starts. If the DA is not present, account
  restrictions persist.
- **Multiple tabs or windows.** Each tab MAY have its own session
  credential. The DA MUST be able to manage multiple concurrent
  handshakes without reusing tokens.
- **Session expires without DA availability.** The credential expires
  naturally. Without a DA to renew, account restrictions persist per
  {{additive-model}}.
- **Minor uninstalls DA and creates new account on a device without
  DA.** The new account has no AAVP signal history. Without a DA on the
  device, there is no handshake and the platform does not apply
  restrictions. This is the same vector that exists without AAVP: a
  minor with access to an uncontrolled device.
- **Minor turns 18.** Parents update the bracket in the DA or remove
  it. The young person presents an OVER_18 credential to platforms.
  Account restrictions are removed.

## CDN and Edge Architecture Compatibility {#cdn-compatibility}

The self-contained credential is compatible with architectures where
validation occurs at edge nodes:

- The edge node can validate vg_signature without querying the origin
  server.
- The VG's verification key can be distributed to all edge nodes of the
  platform.
- The AAVP handshake endpoint MUST NOT be cached:
  Cache-Control: no-store.
- Segmented content responses SHOULD include Vary with a bracket
  identifier so the CDN distinguishes variants by age_bracket.

# Segmentation Accountability Framework {#saf}

## Motivation and Scope {#saf-motivation}

AAVP delivers a reliable age signal (age_bracket) with cryptographic
guarantees. But the system's effectiveness depends on platforms using
that signal to effectively segment content. Without a verification
mechanism, the age signal could become a rubber stamp with no real
effect.

AAVP controls the token issuance and validation phases with
cryptographic guarantees. The content segmentation phase is outside its
direct control. The Segmentation Accountability Framework (SAF)
addresses this gap with an accountability (detection) approach, not an
enforcement (imposition) approach:

- **What AAVP can do:** Define infrastructure for platforms to publicly
  declare their segmentation policies, register them transparently and
  immutably, and allow any party to verify compliance.
- **What AAVP cannot do:** Force a platform to segment correctly.
  Enforcement belongs to regulatory frameworks (DSA, AADC, OSA, COPPA).

The approach is analogous to Certificate Transparency ({{RFC6962}}):
CT does not prevent a CA from issuing a fraudulent certificate, but
guarantees that any issuance is recorded and detectable. Similarly, SAF
does not prevent a platform from ignoring the age signal, but
guarantees that its declared policy is public, auditable, and
verifiable.

## Segmentation Policy Declaration {#spd}

The Segmentation Policy Declaration (SPD) is a signed JSON document
declaring a platform's content segmentation policy by age bracket. It
is a public, machine-readable commitment.

### SPD Endpoint

URI: https://\[platform-domain\]/.well-known/aavp-age-policy.json

HTTP requirements:

- HTTPS mandatory (TLS 1.3).
- Cache-Control: public, max-age=86400 (24 hours).
- Access-Control-Allow-Origin: *.
- Content-Type: application/json.

### SPD JSON Schema

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| spd_version | string | Yes | SPD schema version. Current value: "1.0". |
| platform | string (hostname) | Yes | FQDN of the platform. MUST match the domain serving the endpoint. |
| published | string (ISO 8601) | Yes | Publication date of this policy version. |
| taxonomy_version | string | Yes | Identifier of the content taxonomy used. See {{content-taxonomy}}. |
| segmentation | object | Yes | Map of age brackets to segmentation rules. Each key is a bracket code (UNDER_13, AGE_13_15, AGE_16_17, OVER_18). |
| segmentation[bracket].restricted | array of strings | Yes | Content categories completely blocked for this bracket. |
| segmentation[bracket].adapted | array of strings | Yes | Content categories modified or reduced for this bracket. |
| segmentation[bracket].unrestricted | array of strings | Yes | Content categories with no restrictions for this bracket. Value "*" indicates all categories. |
| policy_url | string (URI) | Yes | URI of the human-readable segmentation policy. |
| ugc_handling | object | No | Declaration of the platform's approach to user-generated content moderation. |
| ugc_handling.moderation | string | Yes (if ugc_handling present) | Moderation approach: "automated", "human", "hybrid". |
| ugc_handling.response_target | string | Yes (if ugc_handling present) | Target time to act on reported content. ISO 8601 duration format (e.g., "PT4H" = 4 hours). |
| ugc_handling.description | string | No | Free-text description of the moderation approach, human-readable. |
| spts | array of objects | No | Signed Policy Timestamps obtained from transparency logs. See {{ptl}}. |
| signature | string | Yes | RSASSA-PKCS1-v1_5 with SHA-256 signature over the canonical JSON (without the signature field itself), encoded in base64url. The signing key is the VG's key. |

### Content Taxonomy {#content-taxonomy}

SAF defines a minimal content category taxonomy. Platforms MUST map
their content to these categories as a minimum.

| Category | Code | Description |
|----------|------|-------------|
| Explicit sexual | explicit-sexual | Content with explicit sexual activity |
| Graphic violence | violence-graphic | Graphic representations of violence or gore |
| Real-money gambling | gambling | Games of chance with real-money wagering |
| Substances | substances | Promotion or depiction of alcohol, tobacco, or drugs |
| Self-harm | self-harm | Content promoting self-harm or suicide |
| Explicit language | profanity | Vulgar language or intense profanity |

Extensibility: Platforms MAY add additional categories with the x-
prefix (e.g., x-de-jugendschutz for German JMStV, x-uk-vsc for BBFC
categories). Extended categories do not affect interoperability:
verifiers that do not recognize them ignore them.

Three action levels per category and bracket:

| Level | Meaning |
|-------|---------|
| restricted | Content completely blocked for this bracket |
| adapted | Content modified or reduced (e.g., filters, edited versions, warnings) |
| unrestricted | No restrictions for this bracket |

### SPD Signing

The SPD is signed with the VG's key to guarantee integrity and
authenticity.

- **Algorithm:** RSASSA-PKCS1-v1_5 with SHA-256 over canonical JSON.
- **Canonical JSON:** The complete SPD document without the signature
  field, serialized with keys sorted alphabetically, no whitespace
  ({{RFC8785}}, JSON Canonicalization Scheme).
- **signature field:** The signing result encoded in base64url without
  padding.
- **Verification:** The DA or any monitor can verify the signature
  using the VG's public key obtained from .well-known/aavp or from the
  IM's trust store.

## Policy Transparency Log {#ptl}

### PTL Architecture

The Policy Transparency Log is an append-only registry, inspired by
Certificate Transparency ({{RFC6962}}), where platform SPDs are
recorded. Multiple independent logs operate in parallel, guaranteeing
decentralization.

~~~ ascii-art
  Platform 1 ---SPD---> Log Operator A ---SPT---> Platform 1
  Platform 1 ---SPD---> Log Operator B ---SPT---> Platform 1
  Platform 2 ---SPD---> Log Operator A
  Platform 2 ---SPD---> Log Operator B
                              |                      |
                         (read)                 (read)
                              v                      v
                         Monitor 1              Device Agent
                              |
                         (alert)
                              v
                           Public
~~~

Principles:

- **Decentralization:** Any organization can operate a log. No official
  or mandatory log exists.
- **Append-only:** Entries MUST NOT be modified or deleted once
  recorded.
- **Verifiability:** Any party can verify the inclusion of an SPD in
  the log.

### Log Structure

Each log maintains a Merkle append-only tree where each leaf contains:

- The complete SPD.
- Registration timestamp in the log.
- Log operator's signature over the entry.

**Signed Policy Timestamp (SPT):**

When registering an SPD, the log issues an SPT as cryptographic proof
of registration:

| Field | Type | Description |
|-------|------|-------------|
| log_id | string | SHA-256 of the log's public key, encoded in base64url. |
| timestamp | string (ISO 8601) | Moment when the log recorded the SPD. |
| signature | string | Log operator's signature over the SPD + timestamp, encoded in base64url. |

The platform includes the obtained SPTs in its SPD (spts field),
demonstrating that the policy has been registered in transparency logs.

### Log Operator Requirements

| Requirement | Description |
|-------------|-------------|
| **Availability** | The log MUST be publicly accessible over HTTPS. |
| **Retention** | Entries MUST be retained for at least 2 years. |
| **Audit** | The log MUST expose a query API allowing verification of any SPD inclusion (Merkle inclusion proof). |
| **Consistency** | The log MUST provide consistency proofs (Merkle consistency proof) demonstrating it has not deleted or modified entries. |

### Monitor Role

Monitors are independent entities that observe transparency logs to
detect anomalies:

- **Policy changes:** Detect when a platform modifies its SPD
  (tightening or relaxing restrictions).
- **Inconsistencies:** Detect if a platform presents different SPDs to
  different observers (split-view attack).
- **Absence of registration:** Alert when a platform declares an
  age_policy in .well-known/aavp but does not register its SPD in any
  log.

Any organization or individual can operate a monitor. No authorization
or registration is required.

## Open Verification Protocol {#ovp}

### OVP Methodology

The Open Verification Protocol defines a standardized methodology for
verifying that a platform complies with its declared SPD:

1. The verifier obtains the platform's SPD from
   .well-known/aavp-age-policy.json.
2. Accesses the platform with valid AAVP tokens for each age bracket
   (UNDER_13, AGE_13_15, AGE_16_17, OVER_18).
3. For each bracket, evaluates the accessibility of content in each
   taxonomy category.
4. Compares observed results against the policy declared in the SPD.

### Compliance Metrics

| Metric | Description | Target |
|--------|-------------|--------|
| **Consistency ratio** | % of categories where observed behavior matches the declared SPD | > 95% |
| **Inter-bracket delta** | Content restriction difference between adjacent brackets | > 0 (each bracket has fewer restrictions than the previous) |
| **False negatives** | Content declared as restricted that is accessible | < 1% |
| **False positives** | Content declared as unrestricted that is blocked | < 5% |

### Verifiers

- **Decentralized:** Any organization or individual can execute OVP
  verifications. No authorization is required.
- **Reference implementation:** An open-source OVP verification tool
  will be published as future work.
- **Crawling limitations:** Dynamic content (algorithmic feeds,
  personalized recommendations) and user-generated content (UGC) make
  exhaustive verification difficult. The sampling methodology
  ({{sampling}}) addresses this limitation with statistical rigor.

### Sampling Methodology {#sampling}

Exhaustive verification of all platform content is infeasible. OVP
adopts a stratified sampling approach with statistical rigor, consistent
with established practices in quality control (ISO 2859) and content
auditing.

Sampling strata:

| Stratum | Description | Methodology |
|---------|-------------|-------------|
| **Curated content** | Catalog, editorial content, fixed sections | Random sampling by SPD category |
| **Algorithmic content** | Feeds, recommendations, trends, "Explore" | Sampling with multiple profiles, times, and contexts |
| **UGC** | User-generated content | Post-publication sampling; measure response time for non-compliant content |

Statistical requirements for OVP reports:

| Requirement | Description |
|-------------|-------------|
| **Sample size** | Document sample size per stratum and category |
| **Confidence interval** | Report results with 95% confidence interval |
| **Margin of error** | Declare margin of error per metric |
| **Period** | Minimum audit duration: not a single snapshot but sampling at multiple moments |

Differentiated metrics by content type:

| Content Type | Primary Metric | Target |
|-------------|---------------|--------|
| Curated | Consistency with SPD | > 99% |
| Algorithmic | Consistency with SPD (with 95% CI) | > 95% |
| UGC | Response time for non-compliant content | Within the SPD's response_target |

## Compliance Signal {#compliance-signal}

### Handshake Extension

The VG MAY optionally include in its handshake response information
about its segmentation policy:

- **SPD hash:** SHA-256 of the canonical JSON of the current SPD,
  encoded in base64url.
- **SPTs:** List of Signed Policy Timestamps proving registration in
  transparency logs.

The DA can verify the SPD hash against the SPD obtained from
.well-known/aavp-age-policy.json and verify SPTs against known log
keys.

### User Indicators

The DA MAY present a compliance indicator to the user with three
states:

| State | Meaning | Condition |
|-------|---------|-----------|
| **Verified policy** | The platform has a signed SPD registered in at least one transparency log | Valid SPD + at least 1 verifiable SPT |
| **Unlogged policy** | The platform has a signed SPD but not registered in transparency logs | Valid SPD + 0 SPTs |
| **No policy** | The platform does not publish an SPD | age_policy field absent in .well-known/aavp or SPD not available |

The DA MUST NOT classify content or judge the quality of the declared
policy. It only informs the user about the platform's transparency
regarding its segmentation commitment.

## SAF Conformance Levels {#saf-conformance-levels}

SAF defines three conformance levels for platforms implementing AAVP:

| Level | Name | Requirements |
|-------|------|-------------|
| **Level 1** | Basic | The platform implements a conformant VG and accepts valid AAVP tokens. |
| **Level 2** | Intermediate | Level 1 + the platform publishes a signed SPD at .well-known/aavp-age-policy.json with documented segmentation policy. |
| **Level 3** | Advanced | Level 2 + the SPD is registered in at least one PTL and the platform undergoes periodic OVP verification (public results). |

No central certification authority exists. SAF conformance levels are
verifiable by any party using the mechanisms defined in this framework.
A platform meets a level if the conditions are verifiably true.

## Limitations and Residual Risk {#saf-limitations}

SAF defines the accountability infrastructure but does not eliminate
all risks:

- **Dynamic content:** Algorithmic feeds and recommendation systems
  generate personalized content that is difficult to audit. The OVP
  sampling methodology ({{sampling}}) addresses this risk with
  stratified sampling of algorithmic content, requiring multiple
  profiles and observation moments.
- **User-generated content (UGC):** Exhaustive real-time UGC
  classification is infeasible. OVP measures UGC compliance as
  response time against non-compliant content, consistent with the
  ugc_handling.response_target field of the SPD. Moderation systems
  have inherent error rates that OVP documents statistically.
- **Segmentation is not censorship:** Segmentation adapts content to
  the age bracket, it does not eliminate it. Platforms SHOULD allow
  documented exceptions (e.g., educational health content for
  AGE_16_17).
- **Declare and not comply:** A platform can publish a restrictive SPD
  and not implement it. This risk is detectable via OVP but not
  preventable by the protocol.

AAVP defines the accountability infrastructure. Effective enforcement
belongs to regulatory frameworks (DSA, AADC, OSA, COPPA).

# Conformance Requirements {#conformance}

For AAVP to be credible as an open standard, implementations of its
three roles MUST be verifiable without depending on a central
certification authority. This section defines a conformance framework
that allows any party to evaluate whether an implementation complies
with the specification, following the protocol's principles of
decentralization and verifiability.

The framework design is grounded in the established practices of major
Internet protocols: FIDO Alliance (tiered certification), OpenID
Connect (role-based conformance suite), Certificate Transparency
({{RFC6962}}, continuous compliance monitoring), NIST CAVP/ACVP
(structured test vectors), Privacy Pass ({{RFC9578}}, specification
test vectors), and PCI DSS (tiered assessment model).

## Conformance Framework Principles {#conformance-principles}

| Protocol Principle | Application to Conformance |
|-------------------|---------------------------|
| **Privacy by Design** | Conformance verification does not require access to user data. Tests operate on synthetic tokens and public configurations. |
| **Decentralization** | No authority "grants" conformance. Any party can execute verifications. Platforms (VGs) individually decide what conformance evidence they require from accepted IMs. |
| **Open Standard** | Verification tools, test vectors, and conformance criteria are public and freely usable. |
| **Data Minimalism** | Verifications evaluate exclusively specification compliance. No usage metrics or operational data beyond what is needed for evaluation are collected. |

## Per-Role Conformance Requirements {#per-role-requirements}

Each protocol role has specific conformance requirements, organized in
three categories per {{RFC2119}} ({{RFC8174}}) terminology:

- **MUST:** Requirements whose non-compliance prevents
  interoperability or compromises security. An implementation that
  fails a MUST requirement is non-conformant.
- **SHOULD:** Requirements whose compliance improves security or
  privacy. Non-compliance MUST be justified and documented.
- **MAY:** Additional capabilities the implementation can offer.

### Device Agent (DA) Requirements

| ID | Requirement | Category | Verification |
|----|------------|----------|-------------|
| DA-01 | Generate tokens of exactly the size specified by the token_type (331 bytes for 0x0001) | MUST | Test vector: token-encoding.json |
| DA-02 | Concatenate token fields in the order and format defined in {{token-structure}} | MUST | Test vector: token-encoding.json |
| DA-03 | Generate 32-byte nonce via OS CSPRNG (SecRandomCopyBytes, SecureRandom, getrandom(2), BCryptGenRandom, or crypto.getRandomValues()) | MUST | Code analysis + statistical test (NIST SP 800-22) |
| DA-04 | Correctly execute the partially blind signature protocol: Prepare, Blind, Finalize | MUST | Test vector: issuance-protocol.json |
| DA-05 | Select token_type per the intersection rule defined in {{crypto-agility}} | MUST | Interoperability test |
| DA-06 | Not reuse nonce between tokens | MUST | Statistical test: generate 10,000 tokens and verify uniqueness |
| DA-07 | Not include hidden metadata in the token beyond the 6 specified fields | MUST | Statistical test: chi-squared on nonce and authenticator bytes (10,000 tokens) |
| DA-08 | Produce unlinkable tokens: two consecutive tokens from the same DA MUST NOT be distinguishable from tokens from different DAs | MUST | Statistical test: binary classifier with epsilon < 0.01 |
| DA-09 | Store keys in secure hardware when the device supports it (Secure Enclave, StrongBox, TPM) | SHOULD | Key attestation ({{device-attestation}}) |
| DA-10 | Implement pre-signing with temporal decoupling ({{pre-signing}}) | SHOULD | Traffic analysis: verify token acquisition does not coincide temporally with presentation |
| DA-11 | Apply 2 KiB message padding in the handshake ({{message-padding}}) | SHOULD | Traffic capture |
| DA-12 | Apply uniform 0-300s jitter before first presentation to a new VG ({{presentation-jitter}}) | SHOULD | Timing analysis |

### Verification Gate (VG) Requirements

| ID | Requirement | Category | Verification |
|----|------------|----------|-------------|
| VG-01 | Correctly parse the token binary format per {{token-structure}} | MUST | Test vector: token-encoding.json |
| VG-02 | Verify the authenticator signature against accepted IM public keys | MUST | Test vector: issuance-protocol.json |
| VG-03 | Reject tokens with expired expires_at, applying asymmetric tolerance: 300s past, 60s future above maximum TTL | MUST | Test vector: token-validation.json |
| VG-04 | Reject tokens with age_bracket outside range (values other than 0x00-0x03) | MUST | Test vector: token-validation.json |
| VG-05 | Reject tokens with unsupported or reserved (0x0000) token_type | MUST | Test vector: token-validation.json |
| VG-06 | Reject tokens with incorrect size for the indicated token_type | MUST | Test vector: token-validation.json |
| VG-07 | Reject tokens whose authenticator does not verify against the IM's public key | MUST | Test vector: token-validation.json |
| VG-08 | Discard the complete token after validation; MUST NOT store, log, or retransmit any token field ({{token-discard}}) | MUST | Code audit + storage inspection |
| VG-09 | Extract exclusively age_bracket from the token for the session | MUST | Functional test: tokens with different nonces and same bracket MUST produce identical sessions |
| VG-10 | Expose .well-known/aavp over HTTPS (TLS 1.3) with the mandatory fields defined in {{service-discovery}} | MUST | Endpoint test |
| VG-11 | Determine the cryptographic scheme from token_type before attempting verification; MUST NOT try multiple schemes | MUST | Code analysis |
| VG-12 | Verify tokens in approximately constant time to avoid timing side-channels | SHOULD | Timing test: difference < 5% between valid and invalid tokens (10,000 samples) |
| VG-13 | Implement the self-contained session credential per {{session-credential}} | SHOULD | Functional test |
| VG-14 | Publish SPD at .well-known/aavp-age-policy.json ({{spd}}) | MAY (SAF Level 2+) | Endpoint test |

### Implementer (IM) Requirements

| ID | Requirement | Category | Verification |
|----|------------|----------|-------------|
| IM-01 | Publish .well-known/aavp-issuer over HTTPS (TLS 1.3) with mandatory fields defined in {{im-keys}} | MUST | Endpoint test |
| IM-02 | Correctly execute BlindSign and metadata-based key derivation (HKDF) per {{partially-blind-signatures}} | MUST | Test vector: issuance-protocol.json |
| IM-03 | Publish keys with validity period not_after - not_before <= 180 days | MUST | Endpoint test: verify not_before and not_after fields |
| IM-04 | Generate token_key_id as SHA-256 of the public key in SPKI DER format | MUST | Test: derive token_key_id from published public key and compare |
| IM-05 | Partially blind signatures MUST be truly blind: no correlation between signing request and resulting token | MUST | Blindness test ({{blindness-test}}) |
| IM-06 | Not retain signing request logs that enable posterior correlation | MUST | Code audit + retention configuration verification |
| IM-07 | Serve keys with Cache-Control: public, max-age=86400 | MUST | HTTP test |
| IM-08 | Not introduce metadata in signatures enabling correlation between tokens | MUST | Cryptographic analysis of signature sample |
| IM-09 | Publish auditable code (preferably open source) | SHOULD | Public repository verification |
| IM-10 | Run the signing service in an ephemeral environment (container without persistent volumes) with memory cleanup after each operation | SHOULD | Infrastructure audit |
| IM-11 | Support OHTTP ({{RFC9458}}) for the DA-IM channel | MAY | Connectivity test |

## Verification Methodology {#verification-methodology}

Conformance verification combines three complementary approaches, from
lowest to highest cost:

### Automated Verification via Test Vectors

The test vectors published in the test-vectors/ directory constitute
the first line of verification. A conformant implementation MUST
produce identical results for all vectors applicable to its role:

| File | Role Verified | What It Validates |
|------|--------------|-------------------|
| token-encoding.json | DA (encoding), VG (decoding) | 331-byte binary format |
| token-validation.json | VG | Validation logic: expiration, clock skew, invalid fields |
| issuance-protocol.json | DA (blinding, finalize), IM (blind sign, key derivation) | Complete partially blind signature flow |

The vector structure follows CFRG conventions and the NIST ACVP
hierarchy: each file contains vectors organized by test case, with
input values, intermediate values, and expected output. Values are
hex-encoded without prefix, consistent with {{RFC9474}} Appendix A.

Verification procedure:

1. For each vector in the applicable file, feed the implementation
   with input values.
2. Compare the output with the expected value (expected_token_hex,
   expected_valid, intermediate values).
3. Record the result: PASS (exact match), FAIL (discrepancy), SKIP
   (non-applicable vector).
4. An implementation is conformant at the test vector level if all
   applicable vectors produce PASS.

### IM Blindness Test {#blindness-test}

Verifying that signatures are truly blind requires a specific
black-box test:

1. The verifier generates N tokens (T1, T2, ..., Tn).
2. Blinds each token with a random factor.
3. Randomly permutes the order.
4. Sends blinded messages in the permuted order to the IM.
5. The IM returns blind signatures.
6. The verifier unblinds signatures and verifies them.
7. The verifier challenges the IM: match signature Si with token Tj.
8. Success is only expected with probability 1/N!.

Test parameters:

| Parameter | Value |
|-----------|-------|
| N (tokens per round) | >= 10 |
| Rounds | >= 100 |
| Acceptance threshold | The IM does not match correctly with probability > 1/N! + 3 sigma |

This test is analogous to the blindness verification documented in
{{RFC9474}} (Section 5.2, Blindness).

### Interoperability Testing

Interoperability verification validates that independent
implementations of the three roles work correctly together. It follows
the OpenID Connect Conformance model and FIDO Alliance
interoperability events.

Interoperability scenarios:

| Scenario | Actors | What It Validates |
|----------|--------|-------------------|
| Complete issuance | DA1 + IM1 | DA obtains a signed token from IM |
| Cross-verification | DA1 + IM1 + VG1 | VG accepts token issued by DA1 via IM1 |
| Multi-IM | DA1 + IM1 + IM2 + VG1 | VG accepts tokens from both IMs |
| Scheme migration | DA1 + IM1 + VG1 (with two token_types) | DA selects correct type; VG accepts both |
| Untrusted IM rejection | DA1 + IM3 (not accepted) + VG1 | VG rejects token signed by IM3 |

Modalities:

- **In-person or virtual event:** Multiple implementers connect their
  implementations and execute scenarios. Each implementation pair MUST
  complete all applicable scenarios.
- **On-demand verification:** An implementation connects to a test
  environment simulating the other roles. This follows the FIDO
  Alliance On-Demand Testing pattern.
- **Continuous integration:** Reference implementations (when
  available) run automatically against each new version. This follows
  the W3C Web Platform Tests (wpt.fyi) model.

## Implementation Conformance Levels {#impl-conformance-levels}

Implementations of each role can achieve three conformance levels. The
levels are cumulative: each includes the requirements of the previous
one.

| Level | Name | Requirements | Verification Method |
|-------|------|-------------|-------------------|
| **Level 1** | Functional | All MUST requirements for the role. All applicable test vectors produce PASS. | Self-assessment with automated test suite. |
| **Level 2** | Verified | Level 1 + all SHOULD requirements for the role. Interoperability testing with at least two independent implementations of the other roles. Operational verification (accessible endpoints, valid keys, response times). | Self-assessment + documented interoperability testing. |
| **Level 3** | Audited | Level 2 + external audit by an independent third party. Audit report published. For IMs: includes blindness test ({{blindness-test}}) and no-log-retention verification. | Third-party audit + report publication. |

Decentralized trust decision: The levels are not "granted" by any
authority. An implementation declares its conformance level and
publishes the evidence. Platforms (VGs) individually decide what
conformance level they require from accepted IMs, analogous to how
each browser decides which CAs to include in its trust store.

## Continuous Operational Verification {#continuous-verification}

Beyond point-in-time conformance verification, the protocol defines
continuous monitoring mechanisms that any party can execute without
authorization or privileged access.

### Discovery Endpoint Monitoring

Any party MAY periodically verify the operational health of IMs and VGs
by querying their public endpoints:

| Verification | Endpoint | Recommended Frequency | Failure Condition |
|-------------|----------|----------------------|------------------|
| IM: active and valid keys | .well-known/aavp-issuer | Daily | Expired keys, token_key_id inconsistent with public key, invalid HTTPS |
| IM: token_type coherence | .well-known/aavp-issuer keys[].token_type | Daily | Unregistered or deprecated value |
| VG: endpoint accessible | .well-known/aavp | Daily | HTTP 4xx/5xx, missing mandatory fields |
| VG: accepted IMs operational | .well-known/aavp accepted_ims[].domain | Daily | Referenced IM has no operational endpoint |

### Token Linting

A token linter is a tool that validates the structure of an AAVP token
without verifying the cryptographic signature. Any party can run a
linter on test tokens to verify a DA's structural conformance:

Linter checks:

1. Total size = expected size per token_type (331 bytes for 0x0001).
2. token_type is a registered and active value ({{token-type-registry}}).
3. age_bracket is within valid range (0x00-0x03).
4. expires_at is a reasonable Unix timestamp (not zero, not excessively
   future).
5. Nonce and authenticator bytes are not all zeros or identical to each
   other.

The linter does not replace cryptographic verification; it complements
automated verification to detect obvious format errors.

### Aggregate Verification Reports

VGs MAY generate aggregate reports on token verification results,
following the DMARC aggregate reporting model ({{RFC7489}}). These
reports provide a feedback mechanism that does not compromise privacy.

Report structure:

| Field | Type | Description |
|-------|------|-------------|
| report_id | string | Unique report identifier |
| vg_domain | string | Domain of the VG generating the report |
| period_start | ISO 8601 | Start of the covered period |
| period_end | ISO 8601 | End of the covered period |
| im_domain | string | Domain of the evaluated IM |
| total_tokens | uint64 | Total tokens presented from the IM in the period |
| valid_tokens | uint64 | Tokens that passed all verifications |
| failures | object | Breakdown by failure type |
| failures.expired | uint64 | Tokens rejected for expiration |
| failures.bad_signature | uint64 | Tokens with invalid signature |
| failures.malformed | uint64 | Tokens with incorrect format |
| failures.unknown_key | uint64 | Tokens with unknown token_key_id |

Privacy properties of the report:

- **No per-user data.** The report contains only aggregate counters. It
  MUST NOT include individual tokens, nonces, timestamps, or any data
  enabling user identification or tracking.
- **Minimum temporal granularity.** The report period MUST be >= 24
  hours to prevent temporal correlation.
- **Optional publication.** The VG MAY send the report to the IM,
  publish it, or both. There is no publication obligation.

## Framework Sources and References {#framework-references}

| Source | Adopted Pattern | Application in AAVP |
|--------|----------------|---------------------|
| FIDO Alliance Authenticator Certification | Tiered levels (L1-L3+) | Three conformance levels per role |
| OpenID Connect Conformance Suite | Role-based profiles, self-certification | Test suite with DA, VG, IM profiles |
| {{RFC7489}} (DMARC) | Anonymous aggregate reports | VG-to-IM verification reports |
| {{RFC6962}} (Certificate Transparency) | Public auditable logs, continuous monitoring | Endpoint monitoring, transparency |
| NIST ACVP | Structured JSON test vectors | test-vectors/ in CFRG format |
| {{RFC9578}} (Privacy Pass) | Specification test vectors, interoperability by independent implementations | Test vectors + interoperability tests |
| PCI DSS | Tiered self-assessment (SAQ) vs. external audit (QSA) | Level 1 = self-assessment, Level 3 = audit |
| ISO/IEC 17065 | Certification lifecycle: assessment - decision - surveillance | Continuous operational verification |
| OAuth Security BCP | MUST/SHOULD/MAY requirements per role | Numbered per-role requirements by category |
| CA/Browser Forum | Platforms as decentralized trust gatekeepers | VGs decide what conformance evidence to require |

# Security Considerations {#security}

Every security protocol must honestly analyze its attack vectors:

| Threat | Mitigation | Residual Risk |
|--------|-----------|---------------|
| **Bypass via device without DA** | Platform policy for sessions without token | Medium |
| **Fraudulent Implementer** | Open source auditing, reputation, platform exclusion | Low |
| **IM domain compromise** | Limited-life keys (<= 6 months), TLS 1.3 + CT for key acquisition, key pinning by VGs, bilateral revocation | Low |
| **MITM on handshake** | TLS 1.3, Certificate Transparency, minimal temporal window | Very Low |
| **Token correlation** | Rotation, nonces, coarse expires_at, partially blind signatures, fixed size (331 bytes) | Very Low |
| **Minor deactivates DA** | OS-level protection, parental PIN, MDM policies. Account-level persistence: restrictions not lifted when DA uninstalled ({{additive-model}}) | Low |
| **Token fabrication** | RSAPBSSA-SHA384 cryptographic signature computationally infeasible to forge | Very Low |
| **Implementer colludes with platform** | The IM knows age_bracket (public metadata) but cannot link the token to a specific DA. The VG also knows age_bracket (it is the received signal). The IM gains no additional information useful for correlation. | Very Low |
| **Token replay** | Unique nonce + expires_at validated by VG against its own clock | Very Low |
| **Clock manipulation** | Asymmetric expires_at tolerance: 300s past, 60s future; rejection of tokens with excessively future expires_at | Low |
| **age_bracket spoofing** | With Partially Blind RSA, the IM can verify age_bracket coherence with the DA's configuration, acting as a second validation barrier | Low |
| **Rooted device** | Optional device attestation ({{device-attestation}}): key attestation detects emulated TEE; weekly rotation forces re-attestation. Recognized limitation ({{security-assumptions}}, S8) | Medium-High |

## Recognized Limitations {#recognized-limitations}

- **Uncontrolled devices:** If a minor accesses from a device without
  DA software, AAVP cannot protect them. The protocol protects the
  doors, not the windows.
- **Compromised devices (root/jailbreak):** On a rooted device, DA
  guarantees can be defeated. Key attestation ({{device-attestation}})
  offers partial detection. See {{security-assumptions}}, assumption
  S8.
- **Implementation quality:** A deficient DA or VG implementation can
  void the protocol's theoretical guarantees.
- **Complement, not substitute:** AAVP is a technical tool that
  complements digital education and family supervision.

# Privacy Considerations {#privacy}

AAVP is designed with privacy as a fundamental property, not an
afterthought. This section consolidates the privacy characteristics
across all protocol components.

**Data minimalism.** The token carries exactly six fields, each with a
specific justification. Any additional field MUST pass the data
minimalism test ({{minimalism-test}}). Fields that enable correlation
or tracking are explicitly excluded ({{excluded-fields}}).

**Unlinkability.** Partially blind signatures ({{partially-blind-signatures}})
ensure the IM cannot link a signed token to the DA that requested it.
Token rotation ({{token-rotation}}) ensures two consecutive tokens from
the same device are cryptographically independent. The VG cannot
correlate two session credentials from the same user
({{session-security}}).

**Token discard.** The VG MUST discard the complete token after
validation ({{token-discard}}). Only the age_bracket is retained for
session management. No token derivatives that could act as
pseudo-identifiers are permitted.

**Privacy partitioning.** The architecture implements privacy
partitioning ({{RFC9614}}): the IM knows the age bracket but not the
destination platform; the VG knows the platform but not the user's
identity. No single entity has access to both pieces of data
simultaneously.

**Fingerprinting prevention.** Coarse timestamps (1-hour precision),
fixed token size (331 bytes), and CSPRNG-generated nonces minimize
fingerprinting vectors ({{fingerprinting-prevention}}).

**Traffic analysis resistance.** Pre-signing ({{pre-signing}}), message
padding ({{message-padding}}), presentation jitter
({{presentation-jitter}}), and optional OHTTP ({{ohttp}}) mitigate
network-level correlation.

**Account-level persistence.** While designed for privacy, the
account-level persistence model ({{additive-model}}) creates a
one-directional privacy trade-off: once an account is marked as minor,
removing the mark requires active proof (OVER_18 credential). This is
a deliberate design choice that prioritizes child protection over the
minor's ability to remove restrictions unilaterally.

# IANA Considerations {#iana}

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

The AAVP specification includes three sets of test vectors with real
cryptographic values computed by the Go reference implementation:

- **Token Encoding** (token-encoding.json): Binary encoding test cases
  verifying that token construction produces the expected hexadecimal
  output for each field combination. Four vectors, one per age bracket.

- **Token Validation** (token-validation.json): Validation logic test
  cases covering expiration, clock skew, malformed tokens, invalid
  signatures, and boundary conditions.

- **Issuance Protocol** (issuance-protocol.json): End-to-end partially
  blind signature flow test cases including key generation (RSA-2048
  with safe primes), key derivation per metadata (HKDF), Prepare,
  Blind, BlindSign, Finalize, and Verify steps, with all intermediate
  values recorded.

The vector structure follows CFRG conventions and the NIST ACVP
hierarchy. Values are hex-encoded without prefix, consistent with
{{RFC9474}} Appendix A.

# Formal Verification {#formal-verification}

The security properties of AAVP have been formally verified using the
Tamarin Prover. Three models are maintained:

- **aavp.spthy**: Core protocol model proving unforgeability (an
  adversary cannot produce a valid token without the IM's signing key),
  nonce uniqueness, metadata binding (the IM's derived key
  cryptographically binds age_bracket and expires_at to the signature),
  and executability (a valid protocol trace exists).

- **aavp-unlinkability.spthy**: Unlinkability model using observational
  equivalence via Tamarin's --diff mode. Proves that an adversary
  observing the protocol cannot distinguish whether two tokens
  originate from the same or different Device Agents.

- **aavp-saf.spthy**: Segmentation Accountability Framework model
  proving seven properties of the SPD/PTL/OVP system, including policy
  immutability, log append-only semantics, and detection of non-
  compliant platforms.

# Acknowledgments
{:numbered="false"}

The author thanks the open-source community for the tools and libraries
that made the reference implementation and formal verification possible.
