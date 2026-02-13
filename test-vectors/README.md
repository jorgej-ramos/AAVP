# AAVP Test Vectors

> **v0.10.0 — Vectores de test para validación de implementaciones**
>
> Conjunto de entradas, valores intermedios y salidas esperadas para verificar la conformidad de implementaciones de Device Agent, Verification Gate e Implementador con la especificación AAVP (PROTOCOL.md).

---

## Convenciones de formato

Los vectores de test de AAVP siguen las convenciones establecidas por el CFRG (Crypto Forum Research Group) y los RFCs de referencia del protocolo:

- **Formato de fichero:** JSON, legible por máquina.
- **Codificación de bytes:** Cadenas hexadecimales en minúsculas, sin prefijo `0x`. Ejemplo: `"0001"` representa los bytes `[0x00, 0x01]`.
- **Enteros:** Los campos `token_type` y `age_bracket` se representan como entero en el JSON y como hexadecimal en el campo `token_hex` del token codificado.
- **Timestamps:** Enteros Unix (segundos desde epoch), representados como `uint64` big-endian en el formato binario del token.
- **Orden de campos en el token binario:** Tal como define PROTOCOL.md sección 2: `token_type` (2) || `nonce` (32) || `token_key_id` (32) || `age_bracket` (1) || `expires_at` (8) || `authenticator` (256) = 331 bytes.
- **Nonces de test:** Generados como SHA-256 de cadenas descriptivas determinísticas. Esto permite reproducibilidad sin comprometer la documentación de que en producción deben usarse CSPRNG del SO.

Estas convenciones son coherentes con:

- RFC 9474, Appendix A — RSA Blind Signatures test vectors (formato hexadecimal, valores intermedios incluidos).
- RFC 9578, Appendix B — Privacy Pass Issuance Protocol test vectors (estructura JSON en implementaciones de referencia).
- draft-irtf-cfrg-cryptography-specification — CFRG guidelines for test vectors (formato legible por máquina, determinismo, cobertura de caminos lógicos).

---

## Ficheros

| Fichero | Descripción | Verificable sin criptografía |
|---------|-------------|:----------------------------:|
| `token-encoding.json` | Codificación/decodificación del formato binario de 331 bytes | Sí |
| `token-validation.json` | Lógica de validación del VG: expiración, clock skew, campos inválidos | Sí |
| `issuance-protocol.json` | Flujo completo de firma parcialmente ciega RSAPBSSA-SHA384 | No (requiere implementación RSAPBSSA) |

---

## Descripción de cada fichero

### token-encoding.json

Verifica que una implementación construye y parsea correctamente el formato binario del token AAVP de 331 bytes. Cada vector incluye:

- Valores individuales de cada campo del token.
- La representación hexadecimal esperada del token completo.
- El tamaño esperado (siempre 331 bytes).

**Qué valida:** Que la implementación concatena los campos en el orden y formato correctos (big-endian para enteros, offsets determinísticos).

**Casos cubiertos:**
- Un token por cada franja de edad (`UNDER_13`, `AGE_13_15`, `AGE_16_17`, `OVER_18`).
- Diferentes timestamps de expiración.
- Nonces distintos para cada vector.

**No requiere operaciones criptográficas.** El campo `authenticator` usa un valor placeholder determinístico; la validez de la firma se verifica en `issuance-protocol.json`.

### token-validation.json

Verifica que la lógica de validación del Verification Gate aplica correctamente las reglas de PROTOCOL.md:

- Tolerancia de clock skew asimétrica (300s pasado, 60s futuro sobre el TTL máximo).
- Rechazo de tokens expirados fuera de tolerancia.
- Rechazo de tokens con `expires_at` excesivamente futuro.
- Rechazo de valores de `age_bracket` fuera de rango.
- Rechazo de `token_type` reservado (0x0000).
- Rechazo de tokens con tamaño incorrecto (!=  331 bytes).
- Detección de `authenticator` manipulado.

Cada vector incluye el token codificado, el tiempo de referencia del VG y el resultado esperado (`valid` o `invalid` con código de error).

**No requiere verificación de firma.** Los vectores verifican la lógica de validación estructural y temporal. La verificación criptográfica del `authenticator` se cubre en `issuance-protocol.json`.

### issuance-protocol.json

Verifica el flujo completo de emisión de un token AAVP mediante firma parcialmente ciega RSAPBSSA-SHA384:

1. **Preparación (DA):** Generación de nonce, construcción de metadatos públicos.
2. **Blinding (DA):** Cegamiento del mensaje con factor aleatorio `r`.
3. **Derivación de clave (IM):** HKDF desde la clave maestra con los metadatos públicos.
4. **Firma ciega (IM):** `BlindSign(sk', blinded_msg)`.
5. **Finalize (DA):** Descegamiento de la firma: `Finalize(pk, msg, metadata, blind_sig, r)`.
6. **Verificación (VG):** Verificación de la firma con la clave pública derivada.

Cada vector incluye:

- Material criptográfico del IM (clave RSA de test con factores primos).
- Todos los valores intermedios del protocolo de firma parcialmente ciega.
- El token final con `authenticator` criptográficamente válido.

**Requiere una implementación de RSAPBSSA-SHA384.** Los valores intermedios permiten verificar cada paso del protocolo de forma independiente. La estructura sigue el formato de RFC 9474 Appendix A (variante parcialmente ciega).

> [!IMPORTANT]
> Los valores criptográficos de `issuance-protocol.json` deben generarse con una implementación conforme de RSAPBSSA-SHA384 (RFC 9474 + draft-irtf-cfrg-partially-blind-rsa). La sección "Generación de los vectores criptográficos" documenta el procedimiento exacto.

---

## Clave RSA de test

Todos los vectores que requieren material criptográfico del Implementador utilizan una **única clave RSA-2048 de test** incluida en `issuance-protocol.json`. Esta clave es exclusivamente para test; no debe usarse en producción.

La clave de test incluye todos los componentes necesarios para verificar cada paso del protocolo:

| Campo | Descripción |
|-------|-------------|
| `n` | Módulo RSA (256 bytes) |
| `e` | Exponente público (típicamente 65537) |
| `d` | Exponente privado |
| `p`, `q` | Factores primos |

El `token_key_id` de test se deriva como SHA-256 de la clave pública en formato SPKI DER, codificado como 32 bytes hexadecimales.

---

## Generación de los vectores criptográficos

### Vectores estructurales (token-encoding, token-validation)

Generados sin dependencias criptográficas. Los nonces son SHA-256 de cadenas descriptivas (ej: `SHA-256("aavp-test-vector-nonce-over18")`). Los authenticators son valores placeholder determinísticos. Cualquier implementación puede verificar estos vectores con operaciones de concatenación y comparación de bytes.

### Vectores de emisión (issuance-protocol)

Requieren una implementación conforme de:

- **RSAPBSSA-SHA384** para la firma parcialmente ciega (RFC 9474 + draft-irtf-cfrg-partially-blind-rsa).
- **HKDF-SHA384** para la derivación de clave por metadato.
- **EMSA-PSS** para el encoding del mensaje (SHA-384, salt length = 0 para variante determinística).

**Procedimiento de generación:**

1. Generar o seleccionar una clave RSA-2048 de test con factores primos conocidos.
2. Para cada franja de edad, definir los metadatos públicos (`age_bracket`, `expires_at`).
3. Derivar la clave de firma: `(sk', pk') = DeriveKeyPair(sk, metadata)` via HKDF.
4. Construir el mensaje a firmar: concatenación binaria de los campos del token sin el `authenticator`.
5. Ejecutar el protocolo de firma parcialmente ciega:
   a. `prepared_msg = Prepare(msg)`
   b. `(blinded_msg, inv) = Blind(pk', prepared_msg, metadata)`
   c. `blind_sig = BlindSign(sk', blinded_msg, metadata)`
   d. `authenticator = Finalize(pk', blind_sig, inv, msg, metadata)`
6. Verificar: `Verify(pk', msg, metadata, authenticator)` devuelve éxito.
7. Registrar todos los valores intermedios en el JSON.

**Implementaciones de referencia para RSAPBSSA:**

- [blindrsa-ts](https://github.com/cloudflare/blindrsa-ts) — TypeScript (Cloudflare). Conforme a RFC 9474.
- [blindrsa-go](https://github.com/cloudflare/pat-go) — Go (Cloudflare). Incluye variante parcialmente ciega.
- [rust-blind-rsa-signatures](https://github.com/nicbh/rust-blind-rsa-signatures) — Rust.

> [!NOTE]
> Hasta que se disponga de una implementación de referencia de AAVP, los vectores criptográficos de `issuance-protocol.json` se marcan con `"status": "awaiting_reference_implementation"`. La estructura del JSON es definitiva; los valores hexadecimales de los campos criptográficos (`blinded_msg`, `blind_sig`, `authenticator`) se completarán cuando se genere la implementación de referencia del DA.

---

## Cómo utilizar estos vectores

### Para implementadores de Device Agent

1. Verificar que la construcción del token produce la salida hexadecimal de `token-encoding.json` para cada vector.
2. Verificar que el flujo de firma ciega produce los valores intermedios de `issuance-protocol.json`.

### Para implementadores de Verification Gate

1. Verificar que el parsing del token extrae correctamente los campos de cada vector de `token-encoding.json`.
2. Verificar que la lógica de validación produce el resultado esperado de `token-validation.json` para cada caso.
3. Verificar que la verificación de firma acepta los tokens válidos de `issuance-protocol.json`.

### Para implementadores de Implementador (IM)

1. Verificar que la derivación de clave por metadato produce las claves derivadas de `issuance-protocol.json`.
2. Verificar que `BlindSign` produce las firmas ciegas esperadas para cada vector.

---

## Fuentes y referencias

- RFC 9474 — RSA Blind Signatures. Appendix A: test vectors para las cuatro variantes de RSABSSA. Formato hexadecimal con valores intermedios. https://www.rfc-editor.org/rfc/rfc9474.html
- RFC 9578 — Privacy Pass Issuance Protocol. Appendix B: test vectors para emisión con VOPRF y Blind RSA. https://www.rfc-editor.org/rfc/rfc9578.html
- RFC 9497 — Oblivious Pseudorandom Functions. Appendix A: test vectors organizados por ciphersuite y modo. https://www.rfc-editor.org/rfc/rfc9497.html
- RFC 8032 — Edwards-Curve Digital Signature Algorithm. Section 7: test vectors con claves, mensajes y firmas. https://www.rfc-editor.org/rfc/rfc8032.html
- draft-irtf-cfrg-cryptography-specification — CFRG guidelines for cryptographic specifications. Recomienda formato JSON legible por máquina, determinismo y cobertura de todos los caminos lógicos. https://www.ietf.org/archive/id/draft-irtf-cfrg-cryptography-specification-02.html
- NIST ACVP — Automated Cryptographic Validation Protocol. Esquema JSON con test groups y test cases. https://pages.nist.gov/ACVP/

---

<div align="center">

**AAVP** · Anonymous Age Verification Protocol · Test Vectors

</div>
