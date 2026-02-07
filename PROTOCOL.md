# AAVP — Especificación Técnica del Protocolo

> **v0.6.0 — Borrador Inicial — Febrero 2026**
>
> Este documento describe la arquitectura, los fundamentos criptográficos y el modelo de seguridad del Anonymous Age Verification Protocol. Para una introducción accesible, consultar [README.md](README.md).

---

## Índice

- [1. Arquitectura del Protocolo](#1-arquitectura-del-protocolo)
- [2. Estructura del Token AAVP](#2-estructura-del-token-aavp)
- [3. Rotación de Tokens](#3-rotación-de-tokens)
- [4. Fundamentos Criptográficos](#4-fundamentos-criptográficos)
- [5. Modelo de Confianza Descentralizado](#5-modelo-de-confianza-descentralizado)
- [6. Flujo Operativo Detallado](#6-flujo-operativo-detallado)
- [7. Modelo de Amenazas](#7-modelo-de-amenazas)
- [8. Trabajo Futuro y Líneas Abiertas](#8-trabajo-futuro-y-líneas-abiertas)
- [Glosario](#glosario)

---

## 1. Arquitectura del Protocolo

### 1.1 Roles del Protocolo

AAVP define tres roles con responsabilidades diferenciadas. El diseño garantiza que ninguno necesita confiar ciegamente en los otros: la verificabilidad criptográfica sustituye a la confianza institucional.

```mermaid
graph LR
    DA[Device Agent] -->|firma parcialmente ciega| IM[Implementador]
    IM -->|firma| DA
    DA -->|token| VG[Verification Gate]
    VG -.->|valida clave publica| IM
    VG -->|sesion| APP[Plataforma]
```

#### Device Agent (DA)

El Device Agent es un **rol abstracto del protocolo**: un componente de software que reside en el dispositivo del menor y es responsable de generar, custodiar y rotar los tokens de edad.

**Qué es:** Una pieza de software que implementa la especificación AAVP para la generación y gestión de tokens. Es el único componente del sistema que conoce la configuración real de franja de edad.

**Qué NO es:** El Device Agent no es sinónimo de "control parental". Es un rol del protocolo que puede ser implementado por distintos vehículos:

| Vehículo de implementación | Ejemplo |
|---------------------------|---------|
| Sistema de control parental | Qustodio, Bark, software de operador |
| Componente nativo del SO | Módulo integrado en iOS, Android, Windows |
| Extensión de navegador | Extensión conforme a la especificación |
| Firmware del dispositivo | Routers con control parental integrado |

La separación entre el rol (Device Agent) y su vehículo de implementación es deliberada: permite que el ecosistema evolucione sin modificar el protocolo. Hoy, el vehículo más probable es el software de control parental existente; mañana, podría ser un componente nativo del sistema operativo.

**Responsabilidades del DA:**
- Generar pares de claves locales en almacenamiento seguro (Secure Enclave, TPM, StrongBox).
- Generar tokens efímeros con la franja de edad configurada.
- Obtener firma parcialmente ciega del Implementador: los metadatos públicos (`age_bracket`, `expires_at`) son visibles al IM, pero el `nonce` permanece cegado.
- Presentar tokens firmados al Verification Gate.
- Rotar tokens antes de su expiración.
- Proteger la configuración de franja mediante PIN parental o mecanismo equivalente a nivel de SO.

#### Verification Gate (VG)

Endpoint dedicado de la plataforma digital que actúa como puerta de entrada al servicio. Valida el token AAVP y establece una sesión interna con la marca de franja de edad.

**Responsabilidades del VG:**
- Exponer un endpoint estándar (`.well-known/aavp`) o registro DNS para anunciar soporte AAVP.
- Validar la firma criptográfica del token contra las claves públicas de Implementadores aceptados.
- Verificar el TTL del token.
- Extraer la franja de edad y establecer una sesión interna.
- Rechazar tokens expirados, malformados o firmados por Implementadores no confiables.

#### Implementador (IM)

Empresa u organización que desarrolla software que actúa como Device Agent, conforme al estándar AAVP.

**Responsabilidades del IM:**
- Publicar su clave pública en su propio dominio mediante un endpoint estándar.
- Mantener código auditable (preferentemente open source).
- Proveer servicio de firma parcialmente ciega al Device Agent.
- Cumplir con la especificación abierta.

### 1.2 Modelo de Puerta de Entrada (Verification Gate)

Un enfoque ingenuo enviaría la credencial de edad en cada petición HTTP, exponiéndola continuamente a posibles interceptaciones. AAVP adopta un modelo diferente: la **puerta de entrada**.

El token de edad solo viaja una vez por sesión, durante un handshake inicial dedicado. Después, la plataforma trabaja con su propio sistema de sesiones.

> **Idea clave:** El token de edad nunca convive con el tráfico regular de la aplicación. Es un canal separado, un handshake puntual. Después, la información "este usuario es menor" es un flag interno de la plataforma, completamente desacoplado del token original.

```mermaid
sequenceDiagram
    participant U as Usuario
    participant DA as Device Agent
    participant VG as Verif. Gate
    participant APP as Plataforma

    U->>APP: Abre app
    APP-->>DA: AAVP soportado

    Note over DA,VG: Handshake TLS
    DA->>DA: Genera token
    DA->>VG: Token firmado
    VG->>VG: Valida firma
    VG-->>DA: OK

    VG->>APP: Sesion + age_bracket
    APP-->>U: Contenido filtrado

    Note over U,APP: Sesion normal

    Note over DA,VG: Revalidacion
    DA->>DA: Nuevo token
    DA->>VG: Token rotado
    VG-->>APP: Renueva sesion
```

**Ventajas del modelo de puerta de entrada:**

- **Superficie de ataque reducida:** el token de edad solo viaja una vez por sesión, no en cada request.
- **Separación de contextos:** la información de edad nunca convive con el tráfico de datos de la aplicación.
- **Compatibilidad:** las plataformas ya gestionan sesiones; AAVP solo añade un paso previo.
- **Ventana temporal mínima para MITM:** interceptar el handshake inicial requiere comprometer TLS en una ventana muy breve. Todos los canales del protocolo (DA-VG y DA-IM) requieren TLS 1.3 o superior.

---

## 2. Estructura del Token AAVP

El token es una estructura criptográfica de tamaño fijo (331 bytes) diseñada para ser mínima. Cada campo tiene una justificación específica y supera el test de minimalismo de datos del protocolo.

```mermaid
classDiagram
    class AAVPToken {
        +uint16 token_type
        +bytes32 nonce
        +bytes32 token_key_id
        +AgeBracket age_bracket
        +uint64 expires_at
        +bytes authenticator
    }
    class AgeBracket {
        UNDER_13
        AGE_13_15
        AGE_16_17
        OVER_18
    }
    AAVPToken --> AgeBracket
```

| Campo | Contenido | Propósito |
|-------|-----------|-----------|
| `token_type` | uint16, identifica el esquema criptográfico | Permite agilidad criptográfica y migración post-cuántica. |
| `nonce` | 32 bytes aleatorios criptográficamente seguros | Previene reutilización y asegura unicidad de cada token. Cegado durante la emisión. |
| `token_key_id` | SHA-256 de la clave pública del IM (32 bytes) | Permite al VG identificar qué clave usar para verificar la firma. |
| `age_bracket` | Enumeración: `UNDER_13` (0x00), `AGE_13_15` (0x01), `AGE_16_17` (0x02), `OVER_18` (0x03) | Señal de franja de edad. Metadato público de la firma parcialmente ciega. |
| `expires_at` | uint64 big-endian, timestamp Unix con precisión de 1 hora | Ventana de validez. Metadato público. La precisión gruesa agrupa tokens temporalmente. |
| `authenticator` | Firma parcialmente ciega RSAPBSSA-SHA384 (256 bytes) | Demuestra que el token proviene de un IM legítimo sin vincular al usuario. |

### Formato binario

El token tiene un formato binario fijo de 331 bytes, sin separadores ni metadatos de codificación. La canonicalización está implícita en el formato: los campos se concatenan en el orden especificado con offsets determinísticos.

```
Offset  Tamaño  Campo                Visibilidad
0       2       token_type           Público
2       32      nonce                Cegado (oculto al IM durante emisión)
34      32      token_key_id         Público
66      1       age_bracket          Metadato público (0x00-0x03)
67      8       expires_at           Metadato público (uint64 BE, precisión 1h)
75      256     authenticator        Firma parcialmente ciega (RSAPBSSA-SHA384)
---
Total: 331 bytes (fijo)
```

Todas las implementaciones conformes deben producir tokens de exactamente 331 bytes. Un token de tamaño diferente es inválido.

### Metadatos públicos vs. contenido cegado

El token AAVP utiliza **firmas parcialmente ciegas** (Partially Blind Signatures). Esto implica una distinción entre dos tipos de contenido dentro del token:

- **Metadatos públicos** (`age_bracket`, `expires_at`): visibles al IM durante el proceso de firma. El IM los utiliza para derivar una clave de firma específica via HKDF. Son parte del contrato visible entre el DA y el IM.
- **Contenido cegado** (`nonce`): oculto al IM durante la emisión. Solo el DA y el VG conocen su valor. El cegamiento criptográfico garantiza que el IM no puede leerlo.

El IM conoce la franja de edad del token que firma, pero **no puede vincular esa información con la identidad del usuario** que la solicita. Dentro de una misma franja, todos los tokens son indistinguibles para el IM. Esto preserva la *unlinkability*: la franja de edad no es un dato personal, es la señal mínima que el protocolo necesita transmitir.

Esta arquitectura es aceptable porque:
1. La franja de edad es precisamente la señal que el protocolo transmite. No es información adicional.
2. El IM no obtiene nada que el VG no obtenga también al verificar el token.
3. El IM puede actuar como segunda barrera contra la suplantación de `age_bracket`, verificando coherencia con la configuración del DA.

### Test de minimalismo de los campos nuevos

Los campos `token_type` y `token_key_id` son adiciones respecto a versiones anteriores de la especificación. Ambos superan el test de minimalismo de datos:

- **`token_type`**: necesario para la agilidad criptográfica (migración post-cuántica). Su valor es idéntico para todos los tokens del mismo esquema. No permite *fingerprinting* individual.
- **`token_key_id`**: necesario para que el VG identifique la clave de verificación sin probar todas las claves conocidas. Derivado de la clave pública del IM (no del usuario). Idéntico para todos los tokens del mismo IM.

### Campos explícitamente excluidos

El token **no contiene** y **no puede contener**:

- Identidad del usuario
- Identificador del dispositivo
- Dirección IP
- Localización geográfica
- Versión del software
- Sistema operativo
- Timestamp de emisión (`issued_at`): eliminado porque la frescura se gestiona con `expires_at` grueso, y un timestamp de emisión con jitter es una superficie innecesaria de *fingerprinting*
- Ningún otro dato que permita correlación o rastreo

Cada dato adicional es un vector potencial de *fingerprinting* y debe justificarse rigurosamente antes de incluirse en futuras versiones del protocolo.

---

## 3. Rotación de Tokens

Incluso sin datos personales, un token estático podría convertirse en un pseudoidentificador persistente si se reutiliza. Por ello, AAVP implementa rotación obligatoria:

```mermaid
stateDiagram-v2
    [*] --> Generado
    Generado --> Validado
    Validado --> Activo
    Activo --> Caducado : TTL 1-4h
    Caducado --> Generado : Nuevo token
    Activo --> Revocado : DA desactivado
```

- **Tiempo de vida máximo (TTL):** Cada token tiene una validez definida por `expires_at`, recomendándose entre 1 y 4 horas. El VG valida `expires_at` contra su propio reloj.
- **Precisión gruesa de `expires_at`:** El valor de `expires_at` se redondea a la hora completa más cercana. Esto implica que todos los tokens emitidos en la misma hora comparten el mismo valor de expiración, lo que incrementa el *anonymity set* y dificulta la correlación temporal.
- **Tolerancia de reloj (*clock skew*):** El VG aplica una tolerancia asimétrica al validar `expires_at`:
  - **Tokens expirados:** El VG acepta tokens cuyo `expires_at` haya pasado hace no más de **5 minutos** (300 segundos). Esto acomoda la sincronización imperfecta de relojes en dispositivos móviles, coherente con la tolerancia de facto de Kerberos (RFC 4120) y las recomendaciones de JWT (RFC 7519).
  - **Tokens del futuro:** El VG rechaza tokens cuyo `expires_at` supere la hora actual del VG en más del TTL máximo permitido (4 horas) **más 1 minuto** (60 segundos). Un `expires_at` excesivamente futuro indica un reloj manipulado o un token fabricado.
  - **Valor canónico:** Las implementaciones conformes usan `CLOCK_SKEW_TOLERANCE_PAST = 300` y `CLOCK_SKEW_TOLERANCE_FUTURE = 60` como valores por defecto. Los VGs pueden ajustar estos valores según su entorno, pero no deben exceder los máximos recomendados.
- **Rotación proactiva:** El Device Agent puede generar un nuevo token antes de la expiración para mantener la continuidad de la sesión.
- **No vinculabilidad (*unlinkability*):** Dos tokens consecutivos del mismo dispositivo no son correlacionables entre sí. Cada token es criptográficamente independiente del anterior.

---

## 4. Fundamentos Criptográficos

### 4.1 Firmas Parcialmente Ciegas (Partially Blind Signatures)

El mecanismo central de AAVP para desacoplar la identidad del usuario de la señal de edad es el uso de **firmas parcialmente ciegas**, una evolución de las firmas ciegas propuestas por David Chaum en 1983.

**Analogía:** Un sobre con papel carbón y una ventanilla transparente. El firmante estampa su firma sobre el sobre cerrado: ve a través de la ventanilla la franja de edad (metadato público), pero el resto del contenido permanece oculto.

**Esquema elegido:** RSAPBSSA-SHA384 (*RSA Partially Blind Signature Scheme with Appendix*), basado en RFC 9474 y draft-irtf-cfrg-partially-blind-rsa. Este esquema permite que el IM vea los metadatos públicos (`age_bracket`, `expires_at`) mientras el `nonce` permanece cegado.

```mermaid
sequenceDiagram
    participant DA as Device Agent
    participant IM as Implementador

    DA->>DA: Genera nonce, construye token
    DA->>DA: Define metadatos publicos: age_bracket, expires_at
    DA->>DA: Ciega el mensaje con factor r

    DA->>IM: Mensaje cegado + metadatos publicos
    Note over IM: Ve age_bracket y expires_at
    Note over IM: NO ve el nonce
    IM->>IM: Verifica coherencia de metadatos
    IM->>IM: Deriva clave: sk' = DeriveKeyPair(sk, metadatos)
    IM->>IM: Firma: blind_sig = BlindSign(sk', msg_cegado)
    IM-->>DA: Devuelve blind_sig

    DA->>DA: Desciega: authenticator = Finalize(pk, token, metadatos, blind_sig, r)
    Note over DA: authenticator es firma valida sobre el token completo
```

**Derivación de clave por metadato:** El IM tiene una sola clave maestra (sk, pk). Para cada combinación de metadatos públicos (age_bracket, expires_at), se deriva automáticamente un par de claves (sk', pk') via HKDF. El VG, que conoce la clave pública maestra y los metadatos del token, realiza la misma derivación para verificar. Esto vincula criptográficamente los metadatos a la firma sin revelar el contenido cegado.

**Resultado:** El Implementador conoce la franja de edad pero **no puede vincular un token firmado con el DA que lo solicitó**. Dentro de la misma franja, todos los tokens son indistinguibles para el IM. La franja no es un dato personal: es la señal que el protocolo transmite.

#### Seguridad del canal DA-IM

La firma parcialmente ciega garantiza que el IM no puede vincular el token finalizado con la petición de firma (*blindness*). Sin embargo, el canal de transporte puede filtrar metadatos que comprometan esta propiedad:

- **Requisito mínimo:** El canal DA-IM requiere TLS 1.3 o superior. El DA verifica la cadena de certificados del IM contra los certificados raíz del sistema operativo. La integridad de los certificados se respalda con Certificate Transparency (RFC 9162).
- **Metadatos de red:** Incluso con TLS, el IM observa la dirección IP del DA, su *fingerprint* TLS (JA3/JA4) y los patrones temporales de las peticiones. Estas fugas son inherentes al transporte TCP/IP y se mitigan parcialmente con la rotación de tokens y la precisión gruesa de `expires_at`.
- **Recomendación para máxima privacidad:** Las implementaciones que busquen minimizar la fuga de metadatos de red pueden utilizar Oblivious HTTP (RFC 9458) para el canal DA-IM, interponiendo un relay que oculte la IP del DA al IM. Esta medida es opcional y queda fuera del alcance mínimo del protocolo.

> [!NOTE]
> El canal DA-VG tiene los mismos requisitos de TLS. La diferencia es que en DA-VG el token ya está finalizado y no contiene datos vinculables con el DA, por lo que la fuga de metadatos de red tiene menor impacto en la privacidad.

**Esquema principal y alternativas:**
- **RSAPBSSA-SHA384** (RFC 9474 + draft-irtf-cfrg-partially-blind-rsa) — Esquema elegido por AAVP.
- Blind BLS Signatures — Alternativa futura por tamaño de firma reducido (48 bytes). Sin RFC publicado.
- ZKP (Bulletproofs) — Complemento para la verificación inicial de edad contra un documento oficial.

### 4.2 Pruebas de Conocimiento Cero (ZKP)

Como alternativa o complemento a las firmas ciegas, AAVP contempla el uso de **pruebas de conocimiento cero** (Zero-Knowledge Proofs).

Un ZKP permite demostrar una afirmación — "mi edad está dentro de la franja X" — sin revelar ningún dato adicional. Esto es especialmente útil en escenarios donde la verificación inicial de la edad se realiza contra un documento oficial: el ZKP demostraría que la fecha de nacimiento cumple el criterio de franja sin exponer la fecha, el nombre ni ningún otro campo del documento.

```mermaid
flowchart LR
    A[Documento fuente] --> B[Motor ZKP]
    B --> C{Prueba: edad en franja}
    C --> D[age_bracket = AGE_13_15]
    B -.->|NO revela| E[Fecha / Nombre / DNI]
```

**Esquemas candidatos:**
- zk-SNARKs (Groth16, PLONK)
- zk-STARKs (sin trusted setup)
- Bulletproofs (para range proofs sobre edad)

### 4.3 Prevención de Fingerprinting

Cada campo del token está diseñado para minimizar la información que podría usarse para identificar o rastrear al usuario:

| Medida | Campo afectado | Propósito |
|--------|---------------|-----------|
| Precisión gruesa | `expires_at` | Redondeo a la hora elimina correlación temporal. Todos los tokens emitidos en la misma hora comparten el mismo valor. |
| Nonce criptográfico | `nonce` | Generado sin derivación de identificadores del dispositivo |
| Metadatos mínimos | `age_bracket`, `expires_at` | Solo dos metadatos públicos. `age_bracket` particiona el *anonymity set* en 4 grupos (inherente al propósito del protocolo). La precisión horaria de `expires_at` agrupa todos los tokens de la misma hora. |
| Rotación frecuente | `expires_at` | Tokens de corta vida impiden seguimiento longitudinal |
| Tamaño fijo | (todo el token) | Todos los tokens tienen exactamente 331 bytes |

---

## 5. Modelo de Confianza Descentralizado

### 5.1 Confianza sin Autoridad Central

AAVP rechaza explícitamente el modelo de Autoridad de Certificación centralizada. La centralización de la certificación crea:

- **Incentivos perversos:** la entidad central adquiere poder de veto.
- **Objetivo prioritario:** de presión política y ataques.
- **Punto único de fallo:** cuya compromisión invalida todo el sistema.

AAVP adopta un **modelo de confianza distribuida**, inspirado en DMARC/DKIM para autenticación de correo electrónico.

**Modelo centralizado (rechazado):**

```mermaid
graph LR
    IM1[IM 1] --> CA[Autoridad Central]
    IM2[IM 2] --> CA
    IM3[IM 3] --> CA
    CA --> P1[Plataforma 1]
    CA --> P2[Plataforma 2]
```

**Modelo AAVP (adoptado) — cada plataforma decide en quién confiar:**

```mermaid
graph LR
    IM1[IM 1] --> P1[Plataforma 1]
    IM1 --> P2[Plataforma 2]
    IM2[IM 2] --> P1
    IM2 --> P2
    IM3[IM 3] --> P1
    IM3 --> P2
```

### 5.2 Mecanismos de Confianza

#### 5.2.1 Estándar abierto y verificable

Cualquier organización puede implementar AAVP. Sus tokens son verificables criptográficamente por cualquier plataforma que también implemente el estándar. No se necesita permiso de ningún tercero. La confianza proviene de la verificabilidad matemática, no de una autorización institucional.

#### 5.2.2 Código auditable

El estándar recomienda firmemente — y la regulación podría exigir — que las implementaciones del Device Agent sean de código abierto o, como mínimo, auditables por terceros independientes. Esto es análogo a los logs de Certificate Transparency: la comunidad puede verificar que el software cumple con la especificación.

#### 5.2.3 Publicación de claves del Implementador

Cada Implementador publica su material criptográfico en su propio dominio. No existe un registro centralizado: el IM es la fuente autoritativa de sus propias claves.

**Canales de publicación:**

- **Primario:** Endpoint HTTPS en el dominio del IM (la ruta y formato concretos se definirán en el Internet-Draft). El endpoint se sirve sobre TLS 1.3, con integridad de certificados respaldada por Certificate Transparency (RFC 9162).
- **Complementario:** Registro DNS TXT bajo `_aavp-keys.[dominio-IM]` con la información mínima para que un VG pueda localizar el endpoint primario.

**Información publicada por el IM:**

Para cada clave de firma activa, el IM publica como mínimo:

- Clave pública (formato según el esquema criptográfico, actualmente RSAPBSSA-SHA384).
- `token_key_id`: SHA-256 de la clave pública (coincide con el campo del token).
- Identificador del esquema criptográfico (coincide con `token_type`).
- Periodo de validez: fecha de inicio y fecha de expiración de la clave.

El `token_key_id` ya presente en cada token permite al VG identificar qué clave usar para verificar la firma sin probar todas las claves conocidas.

> [!IMPORTANT]
> Este endpoint **no es una autoridad de aprobación**. Cualquier organización puede publicar claves en su dominio. La confianza no proviene de estar publicado, sino de la decisión independiente de cada VG de aceptar a ese Implementador.

#### 5.2.4 Ciclo de vida de las claves del Implementador

Las claves de firma del IM tienen una vida limitada. Esto reduce la ventana de exposición si una clave es comprometida y elimina la necesidad de un mecanismo de revocación centralizado.

- **Vida máxima recomendada:** 6 meses (180 días). Las implementaciones no deben aceptar claves con un periodo de validez superior.
- **Rotación:** Cuando el IM genera una nueva clave, publica ambas (antigua y nueva) simultáneamente. El periodo de solapamiento debe ser al menos igual al TTL máximo de los tokens (4 horas), para que los tokens firmados con la clave anterior sigan siendo verificables hasta su expiración. Se recomienda un solapamiento de al menos 24 horas para dar margen a los VGs a refrescar su caché.
- **Expiración natural:** Una clave que supera su fecha de expiración deja de ser válida para la verificación de firmas. Los VGs rechazan tokens cuyo `token_key_id` corresponda a una clave expirada.
- **Sin revocación centralizada:** No existe un mecanismo central para revocar una clave de IM. Si un IM detecta que su clave ha sido comprometida, retira la clave de su endpoint. La revocación efectiva es bilateral: cada VG gestiona su propio trust store y puede retirar a un IM en cualquier momento.

```mermaid
stateDiagram-v2
    [*] --> Publicada : IM genera y publica
    Publicada --> Activa : Fecha de inicio alcanzada
    Activa --> Solapamiento : IM publica nueva clave
    Solapamiento --> Expirada : Fecha de expiración
    Activa --> Retirada : IM detecta compromiso
    Expirada --> [*]
    Retirada --> [*]
```

#### 5.2.5 Gestión de confianza por el Verification Gate

Cada VG mantiene un **trust store local**: una lista de Implementadores aceptados junto con sus claves públicas. La decisión de confiar en un IM es independiente para cada VG, sin mediación de ninguna autoridad central.

- **Obtención de claves:** El VG obtiene las claves públicas directamente del dominio del IM, sobre TLS 1.3. El VG verifica la cadena de certificados TLS y la presencia en logs de Certificate Transparency antes de aceptar el material criptográfico.
- **Caché y refresco:** El VG cachea las claves de los IMs aceptados. El refresco debe ser periódico (recomendado: al menos cada 24 horas) para detectar rotaciones de clave y posibles retiros.
- **Revocación de confianza:** El VG puede retirar a un IM de su trust store en cualquier momento, sin coordinación con otros VGs. Esto es análogo a cómo un navegador puede dejar de confiar en una Autoridad de Certificación unilateralmente.
- **Descubrimiento de nuevos IMs:** Los VGs no confían automáticamente en IMs desconocidos. La incorporación de un nuevo IM al trust store es una decisión deliberada que sigue el proceso reputacional descrito en la sección 5.2.6.

#### 5.2.6 Confianza por reputación

Las plataformas deciden individualmente en qué Implementadores confían, del mismo modo que los navegadores deciden en qué CAs confían para TLS. No hay una decisión centralizada, sino múltiples decisiones independientes que tienden a converger.

```mermaid
flowchart TD
    A[Nuevo Implementador] --> B[Publica clave publica]
    A --> C[Publica codigo fuente]
    C --> D{Auditoria independiente}
    D -->|Cumple especificacion| E[Aceptado por plataformas]
    D -->|No cumple| F[Rechazado]
    E -->|Fraude detectado| G[Confianza revocada]
    G --> F
```

### 5.3 Analogía con DMARC/DKIM

| Aspecto | DMARC/DKIM | AAVP |
|---------|-----------|------|
| Autoridad central | No existe | No existe |
| Quién puede emitir | Cualquier servidor de correo | Cualquier Implementador |
| Quién decide confiar | Cada receptor (Gmail, Outlook...) | Cada plataforma digital |
| Base de la confianza | Cumplimiento del estándar + historial | Cumplimiento del estándar + auditoría |
| Consecuencia del fraude | Correos rechazados / spam | Tokens rechazados por plataformas |

---

## 6. Flujo Operativo Detallado

### 6.1 Configuración Inicial (una sola vez)

Este paso lo realizan los padres o tutores. Es el único momento que requiere intervención humana consciente.

```mermaid
flowchart TD
    A[Padres / Tutores] -->|1. Activan| B[Software con rol de DA]
    B -->|2. Genera claves| C[Secure Enclave / TPM]
    B -->|3. Conecta con| D[Implementador]
    D -->|4. Establece firma parcialmente ciega| B
    A -->|5. Configura franja| B
```

1. Los padres activan la funcionalidad AAVP en el dispositivo del menor. El vehículo puede ser un sistema de control parental, una configuración nativa del SO, u otro software conforme.
2. El software que actúa como Device Agent genera un par de claves locales en el almacenamiento seguro del dispositivo (Secure Enclave en iOS, StrongBox/TEE en Android, TPM en Windows/Linux).
3. El DA establece una conexión única con el servicio de firma del Implementador para obtener la capacidad de firma parcialmente ciega.
4. Se configura la franja de edad correspondiente al menor.

### 6.2 Acceso a una Plataforma (cada sesión)

Este proceso es **completamente transparente** para el usuario:

1. El usuario abre la aplicación o accede al sitio web.
2. El Device Agent detecta que la plataforma soporta AAVP (vía registro DNS `_aavp` o endpoint `.well-known/aavp`).
3. El DA genera un token efímero, ciega el `nonce`, envía el mensaje cegado junto con los metadatos públicos (`age_bracket`, `expires_at`) al Implementador para firma parcialmente ciega, desciega la firma y presenta el token al Verification Gate.
4. El VG valida la firma contra las claves públicas de los Implementadores aceptados, verifica el TTL y extrae la franja de edad.
5. La plataforma establece una sesión con un flag interno de franja de edad.
6. El contenido se filtra según la política de la plataforma para esa franja.
7. Al caducar el token, el DA genera uno nuevo y el VG renueva la sesión. El proceso es transparente.

### 6.3 Desactivación

Si el software que actúa como Device Agent se desactiva durante una sesión activa, deja de emitir tokens. En la siguiente revalidación, la sesión no puede renovarse y transiciona a un estado "no verificado".

La política de qué hacer con sesiones no verificadas es decisión exclusiva de cada plataforma. El protocolo es deliberadamente agnóstico en este punto.

---

## 7. Modelo de Amenazas

Todo protocolo de seguridad debe analizar honestamente sus vectores de ataque:

| Amenaza | Mitigación | Riesgo residual |
|---------|------------|-----------------|
| **Bypass por dispositivo sin DA** | Política de plataforma para sesiones sin token | Medio |
| **Implementador fraudulento** | Auditoría open source, reputación, exclusión por plataformas | Bajo |
| **Compromiso del dominio del IM** | Claves de vida limitada (≤ 6 meses), TLS 1.3 + CT para la obtención de claves, key pinning por VGs, revocación bilateral | Bajo |
| **MITM en handshake** | TLS 1.3, Certificate Transparency, ventana temporal mínima | Muy bajo |
| **Correlación de tokens** | Rotación, nonces, `expires_at` con precisión gruesa, firmas parcialmente ciegas, tamaño fijo (331 bytes) | Muy bajo |
| **Menor desactiva DA** | Protección a nivel de SO, PIN parental, políticas MDM | Medio |
| **Fabricación de tokens** | Firma criptográfica RSAPBSSA-SHA384 computacionalmente inviable de falsificar | Muy bajo |
| **Implementador colude con plataforma** | El IM conoce `age_bracket` (metadato público) pero no puede vincular el token con un DA concreto. El VG también conoce `age_bracket` (es la señal que recibe). El IM no gana información adicional útil para correlación. | Muy bajo |
| **Replay de tokens** | Nonce único + `expires_at` validado por el VG contra su propio reloj | Muy bajo |
| **Manipulación de reloj** | Tolerancia asimétrica de `expires_at`: 300s pasado, 60s futuro; rechazo de tokens con `expires_at` excesivamente futuro | Bajo |
| **Suplantación de `age_bracket`** | Con Partially Blind RSA, el IM puede verificar coherencia de `age_bracket` con la configuración del DA, actuando como segunda barrera de validación | Bajo |

```mermaid
pie title Distribucion de riesgo residual
    "Muy bajo" : 5
    "Bajo" : 4
    "Medio" : 2
```

### Limitaciones reconocidas

- **Dispositivos no controlados:** Si un menor accede desde un dispositivo sin software que actúe como DA, AAVP no puede protegerle. El protocolo protege las puertas, no las ventanas.
- **Calidad de la implementación:** Una implementación deficiente del DA o del VG puede anular las garantías teóricas del protocolo.
- **Complemento, no sustituto:** AAVP es una herramienta técnica que complementa la educación digital y la supervisión familiar.

---

## 8. Trabajo Futuro y Líneas Abiertas

- **Especificación formal:** Desarrollar una especificación técnica completa en formato RFC, incluyendo formatos de mensaje, procedimientos de prueba de conformidad y test vectors.
- **Implementación de referencia:** Crear bibliotecas open source en múltiples lenguajes para Device Agent y Verification Gate.
- **Análisis formal de seguridad:** Verificación formal de las propiedades de privacidad y seguridad mediante herramientas como ProVerif o Tamarin.
- **Pruebas de usabilidad:** Evaluar la experiencia de usuario completa, especialmente la transparencia de la revalidación y la configuración inicial.
- **Evaluación de rendimiento:** Medir el impacto en latencia del handshake y optimizar para conexiones móviles de baja calidad.
- **Extensibilidad del token:** Explorar señales adicionales (preferencias de privacidad parentales, por ejemplo) manteniendo las garantías de anonimato.
- **Gobernanza del estándar:** Definir un modelo de gobernanza comunitaria para la evolución del protocolo, potencialmente bajo el W3C o el IETF.
- **Migración post-cuántica:** El campo `token_type` permite adoptar esquemas de firma ciega post-cuánticos cuando estén estandarizados. Las firmas basadas en retículos (*lattice-based blind signatures*) son la línea de investigación más prometedora. Los estándares NIST PQC actuales (ML-KEM, ML-DSA, SLH-DSA) no cubren firmas ciegas.
- **Registro IANA:** Registrar los valores de `token_type` y la extensión `age_bracket` en los registros IANA correspondientes cuando se formalice el Internet-Draft.

---

## Glosario

| Término | Definición |
|---------|-----------|
| **AAVP** | Anonymous Age Verification Protocol. El protocolo propuesto en este documento. |
| **Blind Signature** | Técnica criptográfica donde un firmante puede firmar un mensaje sin conocer su contenido. |
| **Certificate Transparency (CT)** | Estándar abierto (RFC 9162) que exige a las Autoridades de Certificación registrar todos los certificados emitidos en logs públicos auditables. Reemplaza al *certificate pinning* (deprecado) como mecanismo principal de detección de certificados fraudulentos. |
| **Clock skew** | Diferencia de sincronización entre los relojes de dos sistemas. En AAVP, la tolerancia asimétrica (`CLOCK_SKEW_TOLERANCE_PAST = 300`, `CLOCK_SKEW_TOLERANCE_FUTURE = 60`) acomoda las divergencias entre el reloj del DA y el del VG. |
| **Device Agent (DA)** | Rol del protocolo AAVP: componente de software en el dispositivo del menor que genera y gestiona los tokens de edad. No es sinónimo de "control parental"; puede ser implementado por distintos tipos de software. |
| **Fingerprinting** | Técnica de rastreo que identifica usuarios por características únicas de su dispositivo o comportamiento. |
| **Implementador (IM)** | Organización que desarrolla software conforme al estándar AAVP, actuando como proveedor de la funcionalidad de Device Agent. |
| **Metadato público** | Parte del token visible al IM durante la firma parcialmente ciega, vinculada criptográficamente via derivación de clave. En AAVP: `age_bracket` y `expires_at`. |
| **Partially Blind Signature** | Esquema de firma donde el firmante ve una parte del mensaje (metadato público) pero no el resto. En AAVP, el IM ve `age_bracket` y `expires_at` pero no el `nonce`. |
| **RSAPBSSA** | RSA Partially Blind Signature Scheme with Appendix. Esquema concreto elegido por AAVP, basado en RFC 9474 y draft-irtf-cfrg-partially-blind-rsa. Utiliza SHA-384 como función hash. |
| **`token_key_id`** | SHA-256 de la clave pública del IM. Permite al VG identificar qué clave usar para verificar la firma sin probar todas las claves conocidas. |
| **`token_type`** | Campo del token que identifica el esquema criptográfico utilizado. Permite agilidad criptográfica y migración futura a esquemas post-cuánticos. |
| **Trust store** | Lista de Implementadores aceptados por un Verification Gate, junto con sus claves públicas. Cada VG mantiene su propio trust store de forma independiente. Análogo al trust store de certificados raíz de un navegador. |
| **TTL (Time To Live)** | Tiempo de vida máximo de un token antes de que expire y deba ser reemplazado. |
| **Unlinkability** | Propiedad criptográfica que impide correlacionar dos tokens como pertenecientes al mismo usuario o dispositivo. |
| **Verification Gate (VG)** | Endpoint dedicado de una plataforma que valida tokens AAVP y establece sesiones con franja de edad. |
| **ZKP** | Zero-Knowledge Proof. Prueba criptográfica que demuestra una afirmación sin revelar información adicional. |

---

<div align="center">

**AAVP** · Anonymous Age Verification Protocol · Especificación Técnica · v0.6.0

*Documento de trabajo — Sujeto a revisión*

</div>
