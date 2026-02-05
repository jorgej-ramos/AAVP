# AAVP — Anonymous Age Verification Protocol

> **White Paper v0.1 — Borrador Inicial — Febrero 2026**
>
> Un protocolo abierto y descentralizado para la verificación anónima de edad en plataformas digitales.

---

> [!NOTE]
> **Principio fundamental:** Es posible transmitir una señal fiable de franja de edad a las plataformas digitales sin recopilar datos personales, sin identificación del usuario y sin posibilidad de rastreo inverso.

---

## Índice

- [1. Resumen Ejecutivo](#1-resumen-ejecutivo)
- [2. Definición del Problema](#2-definición-del-problema)
  - [2.1 Situación Actual](#21-situación-actual)
  - [2.2 El Dilema Privacidad vs. Protección](#22-el-dilema-privacidad-vs-protección)
- [3. Visión y Principios de Diseño](#3-visión-y-principios-de-diseño)
  - [3.1 Privacidad por Diseño](#31-privacidad-por-diseño-privacy-by-design)
  - [3.2 Descentralización](#32-descentralización)
  - [3.3 Estándar Abierto](#33-estándar-abierto)
  - [3.4 Minimalismo de Datos](#34-minimalismo-de-datos)
- [4. Arquitectura del Protocolo](#4-arquitectura-del-protocolo)
  - [4.1 Actores del Sistema](#41-actores-del-sistema)
  - [4.2 Modelo de Puerta de Entrada](#42-modelo-de-puerta-de-entrada-verification-gate)
  - [4.3 Estructura del Token](#43-estructura-del-token-aavp)
  - [4.4 Rotación de Tokens](#44-rotación-de-tokens)
- [5. Modelo de Confianza Descentralizado](#5-modelo-de-confianza-descentralizado)
  - [5.1 Confianza sin Autoridad Central](#51-confianza-sin-autoridad-central)
  - [5.2 Mecanismos de Confianza](#52-mecanismos-de-confianza)
  - [5.3 Analogía con DMARC/DKIM](#53-analogía-con-dmarcdkim)
- [6. Fundamentos Criptográficos](#6-fundamentos-criptográficos)
  - [6.1 Firmas Ciegas](#61-firmas-ciegas-blind-signatures)
  - [6.2 Pruebas de Conocimiento Cero](#62-pruebas-de-conocimiento-cero-zkp)
  - [6.3 Prevención de Fingerprinting](#63-prevención-de-fingerprinting)
- [7. Flujo Operativo Detallado](#7-flujo-operativo-detallado)
  - [7.1 Configuración Inicial](#71-configuración-inicial-una-sola-vez)
  - [7.2 Acceso a una Plataforma](#72-acceso-a-una-plataforma-cada-sesión)
  - [7.3 Desactivación del Control Parental](#73-desactivación-del-control-parental)
- [8. Modelo de Amenazas](#8-modelo-de-amenazas)
- [9. Estrategia de Adopción](#9-estrategia-de-adopción)
- [10. Comparativa con Soluciones Existentes](#10-comparativa-con-soluciones-existentes)
- [11. Trabajo Futuro](#11-trabajo-futuro-y-líneas-abiertas)
- [12. Conclusión](#12-conclusión)
- [Glosario](#glosario)

---

## 1. Resumen Ejecutivo

La protección de menores en el entorno digital se enfrenta a lo que parece un dilema irresoluble. Los sistemas de verificación de edad actuales —documentos de identidad, biometría facial, tarjetas de crédito— exigen la recolección de datos personales sensibles, creando nuevos riesgos de privacidad y vigilancia que afectan a todos los usuarios, incluidos los propios menores a los que se pretende proteger.

Este documento presenta el **Anonymous Age Verification Protocol (AAVP)**, un protocolo abierto, descentralizado y respetuoso con la privacidad que permite a las plataformas digitales adaptar su contenido y funcionalidades según la franja de edad del usuario, sin recopilar datos personales identificables y sin posibilidad de rastreo inverso.

AAVP se apoya en tres pilares:
```mermaid
mindmap
  root((AAVP))
    Privacidad por Diseno
      Ningun dato personal abandona el dispositivo
      Garantia matematica no politica de privacidad
    Descentralizacion
      Sin autoridad central
      Sin punto unico de fallo
      Sin captura regulatoria
    Estandar Abierto
      Implementacion libre
      Sin licencias ni permisos
      Barrera exclusivamente tecnica
```

Los bloques de construcción criptográficos necesarios —firmas ciegas, pruebas de conocimiento cero, almacenamiento seguro en hardware— existen y están probados. Lo que falta es la voluntad de articularlos en un estándar común. Este white paper es un primer paso hacia esa articulación.

---

## 2. Definición del Problema

### 2.1 Situación Actual

Las redes sociales y plataformas digitales se enfrentan a una presión regulatoria creciente para verificar la edad de sus usuarios. Las soluciones desplegadas hasta la fecha presentan deficiencias significativas:
```mermaid
quadrantChart
    title Soluciones actuales de verificacion de edad
    x-axis Baja Privacidad --> Alta Privacidad
    y-axis Baja Fiabilidad --> Alta Fiabilidad
    quadrant-1 Objetivo ideal
    quadrant-2 Fiable pero invasivo
    quadrant-3 Ni fiable ni privado
    quadrant-4 Privado pero inutil
    DNI-Pasaporte: [0.15, 0.85]
    Biometria-facial: [0.12, 0.55]
    Tarjeta-de-credito: [0.30, 0.50]
    Autodeclaracion: [0.90, 0.08]
    AAVP: [0.88, 0.82]
```

- **Verificación por documento de identidad.** El usuario sube una copia de su DNI o pasaporte. Esto crea bases de datos de documentos sensibles que se convierten en objetivos de alto valor para atacantes. Las filtraciones de este tipo de datos —ya documentadas en múltiples incidentes— tienen consecuencias especialmente graves porque un documento de identidad no puede «revocarse» como una contraseña.

- **Verificación biométrica.** El análisis facial para estimar la edad implica la recolección de datos biométricos, que el RGPD clasifica como categoría especial con la protección más estricta. Además, estos sistemas presentan sesgos documentados por género, etnia y condiciones de iluminación.

- **Verificación por tarjeta de crédito.** Asume que solo los adultos poseen tarjetas de crédito, lo cual es incorrecto (las tarjetas prepago y las cuentas juveniles son habituales). Además, vincula la identidad financiera del usuario con su actividad en plataformas digitales.

- **Autodeclaración.** El sistema más extendido («tengo más de 18 años») es trivialmente eludible. Ningún menor que desee acceder a un contenido se detiene ante una casilla de verificación.

### 2.2 El Dilema Privacidad vs. Protección

Existe una tensión aparente entre verificar la edad de un usuario y proteger su privacidad: cuanto más fiable es la verificación, más datos personales parece exigir.

> [!IMPORTANT]
> Este white paper sostiene que **dicha tensión es un artefacto del diseño actual de los sistemas, no una limitación fundamental**. Es posible —y este documento describe cómo— transmitir una señal fiable de franja de edad sin transmitir ningún dato que identifique al usuario.

---

## 3. Visión y Principios de Diseño

AAVP se construye sobre cuatro principios no negociables. Cualquier implementación que comprometa alguno de ellos no es conforme con el protocolo.

### 3.1 Privacidad por Diseño (Privacy by Design)

Ningún dato personal identificable abandona el dispositivo del usuario en ningún punto del protocolo. La señal de edad es una afirmación criptográfica anónima, no un dato personal. Esto no es una política de privacidad: es una **garantía matemática**. No es posible, ni siquiera con recursos computacionales ilimitados, vincular un token AAVP con la identidad de un usuario específico.

### 3.2 Descentralización

No existe una autoridad central que certifique, autorice o controle el acceso al protocolo. La confianza emerge de la adopción del estándar abierto, la auditoría pública del código y la reputación de los implementadores. Esto elimina tres riesgos críticos: **incentivos perversos** (nadie tiene poder de veto), **puntos únicos de fallo** (no hay «la autoridad» que comprometer) y **captura regulatoria** (un gobierno no puede presionar a una entidad que no existe).

### 3.3 Estándar Abierto

La especificación del protocolo es pública y libre. Cualquier empresa de control parental puede emitir tokens AAVP y cualquier plataforma digital puede verificarlos, sin licencias, tasas ni permisos. La única barrera de entrada es técnica: implementar correctamente la especificación. Esto es análogo a cómo cualquier servidor puede implementar SMTP para enviar correo electrónico.

### 3.4 Minimalismo de Datos

El token transmite la mínima información necesaria: una **franja de edad**, no una edad exacta. Las franjas propuestas inicialmente son:

| Franja | Codigo | Rango de edad |
|--------|--------|---------------|
| Infantil | `UNDER_13` | Menor de 13 |
| Adolescente temprano | `AGE_13_15` | Entre 13 y 15 |
| Adolescente tardio | `AGE_16_17` | Entre 16 y 17 |
| Adulto | `OVER_18` | Mayor de 18 |

Cada dato adicional es un vector potencial de fingerprinting y debe justificarse rigurosamente antes de incluirse.

---

## 4. Arquitectura del Protocolo

### 4.1 Actores del Sistema

AAVP define tres actores con roles diferenciados. El diseño garantiza que ninguno de ellos necesita confiar ciegamente en los otros: la verificabilidad criptográfica sustituye a la confianza institucional.
```mermaid
graph TB
    subgraph Dispositivo
        DA[Device Agent - Control parental]
    end
    subgraph Implementador_
        IM[Implementador - Empresa de control parental]
    end
    subgraph Plataforma
        VG[Verification Gate - Puerta de entrada]
        APP[Aplicacion Web - Contenido filtrado]
    end

    DA -- "1. Solicita firma ciega" --> IM
    IM -- "2. Devuelve firma ciega" --> DA
    DA -- "3. Presenta token firmado" --> VG
    VG -. "4. Valida firma vs clave publica" .-> IM
    VG -- "5. Establece sesion con flag de edad" --> APP
```

| Actor | Descripción | Responsabilidad |
|-------|-------------|-----------------|
| **Device Agent (DA)** | Sistema de control parental o componente del SO instalado en el dispositivo del menor. | Genera, custodia y rota los tokens de edad. Es el único componente que conoce la configuración real. |
| **Verification Gate (VG)** | Endpoint dedicado de la plataforma digital que actúa como puerta de entrada al servicio. | Valida el token AAVP y establece una sesión interna con la marca de franja de edad. |
| **Implementador (IM)** | Empresa u organización que desarrolla software de control parental conforme al estándar. | Publica su clave pública, mantiene código auditable y cumple con la especificación abierta. |

### 4.2 Modelo de Puerta de Entrada (Verification Gate)

Un enfoque ingenuo enviaría la credencial de edad en cada petición HTTP, exponiéndola continuamente a posibles interceptaciones. AAVP adopta un modelo diferente: la **puerta de entrada**. El token de edad solo viaja una vez por sesión, durante un handshake inicial dedicado. Después, la plataforma trabaja con su propio sistema de sesiones.

> [!TIP]
> **Idea clave:** El token de edad nunca convive con el tráfico regular de la aplicación. Es un canal separado, un handshake puntual. Después, la información «este usuario es menor» es un flag interno de la plataforma, completamente desacoplado del token original.
```mermaid
sequenceDiagram
    participant U as Usuario
    participant DA as Device Agent
    participant VG as Verification Gate
    participant APP as Plataforma

    U->>APP: Abre la aplicacion
    APP-->>DA: Senal AAVP soportado

    rect rgb(219, 234, 254)
        Note over DA,VG: Handshake de verificacion - canal TLS dedicado
        DA->>DA: Genera token efimero y firma ciega
        DA->>VG: Presenta token firmado
        VG->>VG: Valida firma y TTL
        VG-->>DA: Verificacion OK
    end

    VG->>APP: Establece sesion con flag age_bracket
    APP-->>U: Contenido filtrado segun franja de edad

    Note over U,APP: Sesion normal - el token AAVP ya no existe

    rect rgb(254, 243, 199)
        Note over DA,VG: Revalidacion transparente al caducar sesion
        DA->>DA: Genera NUEVO token con rotacion
        DA->>VG: Presenta nuevo token
        VG-->>APP: Renueva sesion
    end
```

**Ventajas del modelo de puerta de entrada:**

- **Superficie de ataque reducida:** el token de edad solo viaja una vez por sesión, no en cada request.
- **Separación de contextos:** la información de edad nunca convive con el tráfico de datos de la aplicación.
- **Compatibilidad:** las plataformas ya gestionan sesiones; AAVP solo añade un paso previo.
- **Ventana temporal mínima para MITM:** interceptar el handshake inicial requiere comprometer TLS con certificate pinning en una ventana muy breve.

### 4.3 Estructura del Token AAVP

El token es una estructura criptográfica diseñada para ser mínima. Cada campo tiene una justificación específica:
```mermaid
classDiagram
    class AAVPToken {
        +AgeBracket age_bracket
        +uint64 issued_at
        +uint64 expires_at
        +bytes32 nonce
        +bytes implementer_sig
    }

    class AgeBracket {
        UNDER_13
        AGE_13_15
        AGE_16_17
        OVER_18
    }

    class DatosExcluidos {
        identidad_usuario
        id_dispositivo
        direccion_IP
        localizacion
        version_software
        sistema_operativo
    }

    AAVPToken --> AgeBracket
    AAVPToken ..> DatosExcluidos : No contiene
```

| Campo | Contenido | Propósito |
|-------|-----------|-----------|
| `age_bracket` | Enumeración: `UNDER_13`, `AGE_13_15`, `AGE_16_17`, `OVER_18` | Señal de franja de edad para filtrado de contenido. |
| `issued_at` | Timestamp Unix con ruido aleatorio | Verificar frescura sin revelar el momento exacto de emisión. |
| `expires_at` | Timestamp Unix | Ventana de validez. Fuerza la rotación. |
| `nonce` | Valor aleatorio criptográficamente seguro | Previene reutilización y asegura unicidad de cada token. |
| `implementer_sig` | Firma ciega (blind signature) | Demuestra que el token proviene de un IM legítimo sin vincular al usuario. |

Es igualmente importante lo que el token **no contiene**: identidad del usuario, identificador del dispositivo, dirección IP, localización geográfica, versión del software, ni ningún otro dato que permita correlación o rastreo.

### 4.4 Rotación de Tokens

Incluso sin datos personales, un token estático podría convertirse en un pseudoidentificador persistente si se reutiliza. Por ello, AAVP implementa rotación obligatoria:
```mermaid
stateDiagram-v2
    [*] --> TokenActivo : DA genera token
    TokenActivo --> Validado : VG verifica firma
    Validado --> SesionActiva : Plataforma crea sesion
    SesionActiva --> Caducado : TTL expira entre 1 y 4h
    Caducado --> TokenActivo : DA genera NUEVO token
    SesionActiva --> NoVerificado : Control parental desactivado

    state TokenActivo {
        [*] --> Firmado : Firma ciega del IM
        Firmado --> Listo : Nonce y ruido temporal
    }

    note right of Caducado
        Cada nuevo token es
        criptograficamente
        independiente del anterior.
        No hay correlacion posible.
    end note
```

- **Tiempo de vida máximo (TTL):** Cada token tiene una validez configurable, recomendándose entre 1 y 4 horas.
- **Rotación proactiva:** El Device Agent puede generar un nuevo token antes de la expiración para mantener la continuidad.
- **No vinculabilidad (*unlinkability*):** Dos tokens consecutivos del mismo dispositivo no son correlacionables entre sí. Cada token es criptográficamente independiente del anterior.

---

## 5. Modelo de Confianza Descentralizado

### 5.1 Confianza sin Autoridad Central

AAVP rechaza explícitamente el modelo de Autoridad de Certificación centralizada. ¿Por qué? Porque la centralización de la certificación crea incentivos perversos: la entidad central adquiere poder de veto sobre quién participa en el ecosistema, se convierte en objetivo prioritario de presión política, y genera un punto único de fallo cuya compromisión invalida todo el sistema.

AAVP adopta un **modelo de confianza distribuida**, inspirado en protocolos como DMARC/DKIM para autenticación de correo electrónico.
```mermaid
graph LR
    subgraph Centralizado - RECHAZADO
        CA[Autoridad Central]
        IM1a[IM 1] --> CA
        IM2a[IM 2] --> CA
        IM3a[IM 3] --> CA
        CA --> P1a[Plataforma 1]
        CA --> P2a[Plataforma 2]
    end

    subgraph AAVP - ADOPTADO
        IM1b[IM 1] --> P1b[Plataforma 1]
        IM1b --> P2b[Plataforma 2]
        IM2b[IM 2] --> P1b
        IM2b --> P2b
        IM3b[IM 3] --> P1b
        IM3b --> P2b
    end
```

### 5.2 Mecanismos de Confianza

#### 5.2.1 Estándar abierto y verificable

Cualquier empresa puede implementar AAVP. Al hacerlo, sus tokens son verificables criptográficamente por cualquier plataforma que también implemente el estándar. No se necesita permiso de ningún tercero. La confianza proviene de la verificabilidad matemática, no de una autorización institucional.

#### 5.2.2 Código auditable

El estándar recomienda firmemente —y la regulación podría exigir— que las implementaciones de Device Agent sean de código abierto o, como mínimo, auditables por terceros independientes. Esto es análogo a los logs de Certificate Transparency: la comunidad puede verificar que el software cumple con la especificación.

#### 5.2.3 Registro público de Implementadores

Se propone un registro público descentralizado (potencialmente basado en un log transparente) donde los Implementadores publican sus claves públicas y declaran conformidad con el estándar. Es importante subrayar que esto **no es una autoridad de aprobación**: es un directorio público auditable, abierto a cualquiera.

#### 5.2.4 Confianza por reputación

Las plataformas digitales deciden individualmente en qué Implementadores confían, del mismo modo que los navegadores web deciden en qué CAs confían para TLS. No hay una decisión centralizada, sino múltiples decisiones independientes que tienden a converger.
```mermaid
flowchart TD
    IM[Nuevo Implementador] -->|Publica clave publica| REG[Registro Publico Descentralizado]
    IM -->|Publica codigo fuente| GH[Codigo Auditable Open Source]

    GH -->|Comunidad audita| AUDIT{Cumple el estandar}
    AUDIT -->|Si| TRUST[Plataformas incluyen clave en trust store]
    AUDIT -->|No| REJECT[Plataformas rechazan tokens]

    TRUST -->|Emite tokens fraudulentos| REVOKE[Plataformas retiran confianza]
    REVOKE --> REJECT
```

### 5.3 Analogía con DMARC/DKIM

Para entender intuitivamente cómo funciona este modelo, resulta útil compararlo con DMARC/DKIM, el sistema de autenticación de correo electrónico:

| Aspecto | DMARC/DKIM | AAVP |
|---------|-----------|------|
| Autoridad central | No existe | No existe |
| Quién puede emitir | Cualquier servidor de correo | Cualquier Implementador |
| Quién decide confiar | Cada receptor (Gmail, Outlook...) | Cada plataforma digital |
| Base de la confianza | Cumplimiento del estandar + historial | Cumplimiento del estandar + auditoria |
| Consecuencia del fraude | Correos rechazados / spam | Tokens rechazados por plataformas |

---

## 6. Fundamentos Criptográficos

### 6.1 Firmas Ciegas (Blind Signatures)

El mecanismo central de AAVP para desacoplar la identidad del usuario de la señal de edad es el uso de **firmas ciegas**, una técnica propuesta por David Chaum en 1983. La analogía clásica: un sobre con papel carbón. El firmante estampa su firma sobre el sobre cerrado, y la firma se transfiere al documento interior sin que el firmante lo vea.
```mermaid
sequenceDiagram
    participant DA as Device Agent
    participant IM as Implementador

    DA->>DA: 1. Genera token con franja de edad
    DA->>DA: 2. Enmascara token via blinding

    rect rgb(243, 232, 255)
        DA->>IM: 3. Envia token ENMASCARADO
        Note over IM: El IM NO puede ver el contenido del token
        IM->>IM: 4. Firma el token enmascarado
        IM-->>DA: 5. Devuelve firma
    end

    DA->>DA: 6. Desenmascara la firma
    Note over DA: Resultado: firma valida sobre el token original

    DA->>DA: 7. Token listo para presentar al VG

    Note over DA,IM: El IM nunca supo que token firmo. No puede vincular firma con usuario.
```

El resultado práctico: el Implementador puede certificar que un token es legítimo (proviene de una instalación real de control parental) sin saber qué token ha firmado. **Ni siquiera el Implementador puede vincular un token con un usuario.**

### 6.2 Pruebas de Conocimiento Cero (ZKP)

Como alternativa o complemento a las firmas ciegas, AAVP contempla el uso de **pruebas de conocimiento cero** (Zero-Knowledge Proofs). Un ZKP permite demostrar una afirmación —por ejemplo, «mi edad está dentro de la franja X»— sin revelar ningún dato adicional. Esto sería especialmente útil en escenarios donde la verificación inicial de la edad se realiza contra un documento oficial: el ZKP demostraría que la fecha de nacimiento cumple el criterio de franja sin exponer la fecha.
```mermaid
flowchart LR
    DOC[Documento oficial] --> ZKP[Motor ZKP]
    ZKP --> PROOF[Prueba verificable: Edad en AGE_13_15]
    ZKP -.->|NO revela| HIDDEN[Fecha exacta / Nombre / Num documento]
```

### 6.3 Prevención de Fingerprinting

Cada campo del token está diseñado para minimizar la información que podría usarse para identificar o rastrear al usuario:

- **Ruido temporal:** El campo `issued_at` incluye un jitter aleatorio para evitar correlación por momento de emisión.
- **Nonce único:** Generado criptográficamente sin derivación de ningún identificador del dispositivo.
- **Sin metadatos:** No se incluye versión del software, sistema operativo ni ningún dato del entorno.
- **Rotación frecuente:** Tokens de corta vida que impiden el seguimiento longitudinal.

---

## 7. Flujo Operativo Detallado

### 7.1 Configuración Inicial (una sola vez)

Este paso lo realizan los padres o tutores y es el único momento en que se requiere intervención humana consciente:
```mermaid
flowchart TD
    P[Padres o Tutores] -->|1. Activan control parental| CP[Sistema de Control Parental]
    CP -->|2. Genera par de claves| SE[Secure Enclave o TPM]
    CP -->|3. Conexion unica| IM[Servicio de firma del Implementador]
    IM -->|4. Capacidad de firma ciega| CP
    P -->|5. Configura franja de edad| CP
```

1. Los padres activan el control parental en el dispositivo del menor.
2. El Device Agent genera un par de claves locales en el almacenamiento seguro del dispositivo (Secure Enclave, TPM o equivalente).
3. El DA establece una conexión única con el servicio de firma del Implementador para obtener la capacidad de firma ciega.
4. Se configura la franja de edad correspondiente al menor.

### 7.2 Acceso a una Plataforma (cada sesión)

Este proceso es **completamente transparente** para el usuario:

1. El usuario abre la aplicación o accede al sitio web.
2. El Device Agent detecta que la plataforma soporta AAVP (vía un registro DNS bien conocido o un endpoint estándar tipo `.well-known/aavp`).
3. El DA genera un token efímero, lo firma ciegamente y lo presenta al Verification Gate.
4. El VG valida la firma, verifica el TTL y extrae la franja de edad.
5. La plataforma establece una sesión con un flag interno de franja.
6. El contenido se filtra según la política de la plataforma para esa franja.
7. Al caducar la sesión, el proceso se repite automáticamente con un nuevo token.

### 7.3 Desactivación del Control Parental

Si el control parental se desactiva durante una sesión activa, el Device Agent deja de emitir tokens. En la siguiente revalidación, la sesión no puede renovarse y transiciona a un estado «no verificado». La política de qué hacer con sesiones no verificadas es decisión exclusiva de cada plataforma —el protocolo es deliberadamente agnóstico en este punto.

---

## 8. Modelo de Amenazas

Todo protocolo de seguridad debe analizar honestamente sus vectores de ataque:

| Amenaza | Descripción | Mitigación | Riesgo residual |
|---------|-------------|------------|-----------------|
| **Bypass por dispositivo alternativo** | El menor accede desde un dispositivo sin DA instalado. | Política de plataforma para sesiones sin token. | Medio |
| **Implementador fraudulento** | Un IM emite tokens de adulto a menores. | Auditoría open source, pérdida de reputación, exclusión. | Bajo |
| **MITM en el handshake** | Interceptar el token durante la validación inicial. | TLS con certificate pinning, ventana temporal mínima. | Muy bajo |
| **Correlación de tokens** | Vincular tokens sucesivos del mismo usuario. | Rotación, nonces aleatorios, ruido temporal, firmas ciegas. | Muy bajo |
| **Desactivación del control parental** | El menor desactiva el DA por su cuenta. | Protección a nivel de SO, PIN parental. | Medio |
| **Fabricación manual de tokens** | Crear tokens válidos sin un DA legítimo. | Firma criptográfica del IM lo hace computacionalmente inviable. | Muy bajo |
```mermaid
pie title Distribucion de riesgo residual
    "Muy bajo" : 3
    "Bajo" : 1
    "Medio" : 2
```

### Limitaciones Reconocidas

AAVP no pretende ser una solución completa. Es importante ser transparentes:

- **Dispositivos no controlados:** Si un menor accede desde un dispositivo sin control parental, AAVP no puede protegerle. El protocolo protege las puertas, no las ventanas.
- **Calidad de la implementación:** Como cualquier protocolo criptográfico, una implementación deficiente puede anular las garantías teóricas.
- **Complemento, no sustituto:** AAVP es una herramienta técnica que complementa la educación digital y la supervisión familiar, no las reemplaza.

---

## 9. Estrategia de Adopción

### El Problema del Bootstrapping

Todo protocolo de dos lados enfrenta el clásico dilema del huevo y la gallina. Para superarlo se propone una estrategia en tres fases:
```mermaid
gantt
    title Hoja de ruta de adopcion de AAVP
    dateFormat YYYY
    axisFormat %Y

    section Fase 1 Especificacion
    Publicacion del estandar abierto     :a1, 2026, 1y
    Implementaciones de referencia OSS   :a2, 2026, 1y
    Pruebas controladas cripto y UX      :a3, 2026, 1y

    section Fase 2 Adopcion temprana
    Integracion en SO Apple y Google     :b1, 2027, 2y
    Controles parentales emiten tokens   :b2, 2027, 1y
    Plataformas pioneras con VG          :b3, 2027, 2y

    section Fase 3 Masa critica
    Impulso regulatorio DSA COPPA AADC   :c1, 2028, 2y
    Ciclo virtuoso de adopcion           :c2, 2029, 2y
    Gobernanza comunitaria W3C o IETF    :c3, 2029, 2y
```

### Compatibilidad Regulatoria

AAVP está diseñado para encajar en los marcos regulatorios existentes y emergentes:

| Regulación | Compatibilidad con AAVP |
|------------|------------------------|
| **RGPD / GDPR** | Al no procesar datos personales, AAVP minimiza las obligaciones regulatorias. No se requiere consentimiento específico. |
| **Digital Services Act (DSA)** | La DSA exige medidas de protección de menores. AAVP proporciona la señal técnica sin crear sistemas de vigilancia. |
| **COPPA (EE.UU.)** | Facilita el cumplimiento al identificar menores de 13 sin recopilar datos personales de menores. |
| **Age Appropriate Design Code (UK)** | Compatible con el enfoque de «diseño apropiado para la edad» al proporcionar la señal para adaptar la experiencia. |

---

## 10. Comparativa con Soluciones Existentes

| Criterio | AAVP | DNI / Pasaporte | Biometria facial | Tarjeta credito | Autodeclaracion |
|----------|------|-----------------|------------------|-----------------|-----------------|
| **Privacidad** | Alta | Muy baja | Muy baja | Baja | Alta |
| **Fiabilidad** | Alta | Alta | Media | Media | Nula |
| **Descentralizado** | Si | No | No | No | Si |
| **Riesgo de filtracion** | Minimo | Critico | Critico | Alto | Ninguno |
| **Coste de implementacion** | Medio | Alto | Muy alto | Medio | Bajo |
| **RGPD nativo** | Si | No | No | No | Si |

---

## 11. Trabajo Futuro y Líneas Abiertas

- **Especificación formal:** Desarrollar una especificación técnica completa en formato RFC, incluyendo formatos de mensaje, algoritmos específicos y procedimientos de prueba de conformidad.
- **Implementación de referencia:** Crear bibliotecas open source en múltiples lenguajes para Device Agent y Verification Gate.
- **Análisis formal de seguridad:** Verificación formal de las propiedades de privacidad y seguridad mediante herramientas como ProVerif o Tamarin.
- **Pruebas de usabilidad:** Evaluar la experiencia de usuario completa, especialmente la transparencia de la revalidación y la configuración inicial.
- **Evaluación de rendimiento:** Medir el impacto en latencia del handshake y optimizar para conexiones móviles de baja calidad.
- **Extensibilidad del token:** Explorar señales adicionales (preferencias de privacidad parentales, por ejemplo) manteniendo las garantías de anonimato.
- **Gobernanza del estándar:** Definir un modelo de gobernanza comunitaria para la evolución del protocolo, potencialmente bajo el W3C o el IETF.

---

## 12. Conclusión

La protección de menores en el entorno digital no tiene por qué venir a costa de la privacidad de todos los usuarios. AAVP demuestra que es técnicamente viable construir un sistema de verificación de edad que sea simultáneamente fiable, anónimo, descentralizado y compatible con los marcos regulatorios existentes.

Los bloques de construcción criptográficos necesarios existen y están probados. Lo que falta es la voluntad de articularlos en un estándar abierto y la presión regulatoria y social para impulsar su adopción.

> [!IMPORTANT]
> Invitamos a la comunidad técnica, a los reguladores, a las empresas de control parental y a las plataformas digitales a contribuir a la evolución de AAVP hacia un estándar robusto, auditable y verdaderamente protector tanto de los menores como de la privacidad de todos.

---

## Glosario

| Término | Definición |
|---------|-----------|
| **AAVP** | Anonymous Age Verification Protocol. El protocolo propuesto en este documento. |
| **Blind Signature** | Técnica criptográfica donde un firmante puede firmar un mensaje sin conocer su contenido. |
| **Certificate Pinning** | Técnica de seguridad que asocia un servicio con su certificado específico, previniendo ataques MITM. |
| **Device Agent (DA)** | Software de control parental o componente del SO que genera y gestiona los tokens AAVP. |
| **Fingerprinting** | Técnica de rastreo que identifica usuarios por características únicas de su dispositivo o comportamiento. |
| **Implementador (IM)** | Empresa que desarrolla software conforme al estándar AAVP. |
| **TTL (Time To Live)** | Tiempo de vida máximo de un token antes de que expire y deba ser reemplazado. |
| **Verification Gate (VG)** | Endpoint dedicado de una plataforma que valida tokens AAVP y establece sesiones. |
| **ZKP** | Zero-Knowledge Proof. Prueba criptográfica que demuestra una afirmación sin revelar información adicional. |

---

<div align="center">

**AAVP** · Anonymous Age Verification Protocol · v0.1

*Documento de trabajo — Sujeto a revisión*

</div>
