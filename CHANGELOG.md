# Changelog

Todos los cambios notables de este proyecto se documentan en este archivo.

El formato se basa en [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/)
y este proyecto se adhiere a [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Sección 8 "Segmentation Accountability Framework (SAF)" en PROTOCOL.md: declaración de política de segmentación firmada (SPD) en `.well-known/aavp-age-policy.json`, logs de transparencia (PTL) inspirados en Certificate Transparency, protocolo de verificación abierto (OVP) y señal de cumplimiento en el handshake. Taxonomía de contenido mínima con 6 categorías extensibles. Tres niveles de conformidad.
- Campo opcional `age_policy` en endpoint `.well-known/aavp` de PROTOCOL.md.
- Términos SAF, SPD, PTL, SPT y OVP añadidos al glosario de PROTOCOL.md.

### Changed

- Área "Segmentación de contenido" de SECURITY-ANALYSIS.md actualizada de rojo a amarillo: SAF mitiga la brecha de verificación con riesgo residual en contenido dinámico.
- Distribución del semáforo actualizada: 0 rojo, 6 amarillo, 1 verde.
- Vulnerabilidad V11 (segmentación no verificable) y supuesto S12 de SECURITY-ANALYSIS.md marcados como mitigados.
- Sección 7 de SECURITY-ANALYSIS.md reescrita para referenciar el SAF formalizado.
- Secciones 8 y 9 de PROTOCOL.md renumeradas a 9 y 10 respectivamente.

## [0.7.0] - 2026-02-09

### Added

- Sección 7 "Credencial de Sesión del Verification Gate" en PROTOCOL.md: credencial de sesión autocontenida emitida por el VG tras validar un token AAVP con descarte obligatorio del token, TTL de 15-30 minutos (siempre ≤ TTL del token), renovación con token AAVP independiente y no vinculable, modelo aditivo (sin DA = sin restricciones) con persistencia a nivel de cuenta (la franja menor persiste aunque el DA desaparezca; solo una credencial `OVER_18` la retira), y compatibilidad con CDN/edge.
- Términos "Credencial de sesión", "Fail-closed" y "Self-contained" añadidos al glosario de PROTOCOL.md.
- Endpoint `.well-known/aavp-issuer` especificado en PROTOCOL.md sección 5.2.3: esquema JSON con campos `issuer`, `signing_endpoint`, `keys[]` (con `token_key_id`, `token_type`, `public_key` en SPKI DER base64url, `not_before`, `not_after`). Cache de 24 horas.
- Endpoint `.well-known/aavp` especificado en PROTOCOL.md sección 5.3: esquema JSON con campos `aavp_version`, `vg_endpoint`, `accepted_ims[]` (con `domain` y `token_key_ids` opcional), `accepted_token_types`. Cache de 1 hora.
- Registros DNS complementarios `_aavp` y `_aavp-keys` con formato TXT definido.
- Cadena de prioridad de descubrimiento: caché local → `.well-known` HTTPS → DNS TXT.
- Registro informal de valores de `token_type` en PROTOCOL.md sección 5.4.
- Términos `.well-known/aavp` y `.well-known/aavp-issuer` añadidos al glosario de PROTOCOL.md.

### Changed

- Nueva subsección "Generación del nonce" en PROTOCOL.md sección 2: APIs de CSPRNG obligatorias por plataforma, fuentes prohibidas y tests de conformidad con NIST SP 800-22.
- Vulnerabilidad T-4.6 (calidad de fuente de aleatoriedad del nonce) de SECURITY-ANALYSIS.md marcada como resuelta.
- Área "Estructura del token" de SECURITY-ANALYSIS.md actualizada de amarillo a verde: 6 vulnerabilidades resueltas.
- Distribución del semáforo actualizada: 1 rojo, 5 amarillo, 1 verde.
- Vulnerabilidades I-5.2 (gestión de sesiones post-handshake) e I-5.3 (política de contenido no verificado) de SECURITY-ANALYSIS.md marcadas como resueltas.
- Recomendación R6 (política de sesiones no verificadas) de SECURITY-ANALYSIS.md marcada como resuelta.
- Área "Gestión de sesiones (VG)" de SECURITY-ANALYSIS.md actualizada de rojo a amarillo.
- Supuesto S5 (sesiones post-handshake seguras) fortalecido con referencia a la nueva especificación.
- Amenaza "Menor desactiva DA" en PROTOCOL.md reclasificada de riesgo Medio a Bajo gracias a la persistencia a nivel de cuenta.
- Secciones 7 y 8 de PROTOCOL.md renumeradas a 8 y 9 respectivamente.
- Vulnerabilidad I-5.1 (descubrimiento de servicio) de SECURITY-ANALYSIS.md marcada como resuelta.
- Especificación E4 (formato de `.well-known/aavp`) de SECURITY-ANALYSIS.md marcada como resuelta.
- Referencia "Pendiente: E4" eliminada de las áreas "Modelo de confianza" y "Gestión de sesiones" del semáforo.
- Sección 5.2.3 de PROTOCOL.md actualizada: reemplazado placeholder ("se definirá en el Internet-Draft") con especificación completa.
- Sección 5.3 "Analogía con DMARC/DKIM" renumerada a 5.5.

## [0.6.0] - 2026-02-07

### Changed

- Modelo de confianza del registro de Implementadores especificado en PROTOCOL.md: auto-publicación de claves por cada IM en su dominio (TLS 1.3 + CT). Sin registro centralizado. Vulnerabilidad S11 y recomendación R8 de SECURITY-ANALYSIS.md marcadas como resueltas.
- Ciclo de vida de claves del IM definido en PROTOCOL.md: vida máxima de 6 meses, rotación con solapamiento, expiración natural.
- Revocación bilateral definida: cada VG gestiona su trust store de forma independiente. Vulnerabilidad S14 y especificación E3 de SECURITY-ANALYSIS.md marcadas como resueltas.
- Vector V-2.4 (ataque al registro de IMs) reclasificado en SECURITY-ANALYSIS.md: de Crítico (registro central) a Medio (compromiso de dominio individual).
- Área "Modelo de confianza (registro de IMs)" de SECURITY-ANALYSIS.md actualizada de rojo a amarillo.
- Seguridad del canal DA-IM especificada en PROTOCOL.md: TLS 1.3 con Certificate Transparency (RFC 9162). OHTTP (RFC 9458) recomendado como medida opcional de máxima privacidad. Vulnerabilidad S9 y recomendación R9 de SECURITY-ANALYSIS.md marcadas como resueltas.
- Tolerancia asimétrica de *clock skew* definida para validación de `expires_at`: 300 segundos para tokens recién expirados, 60 segundos para tokens del futuro. Coherente con Kerberos (RFC 4120) y JWT (RFC 7519). Vulnerabilidad S10 y recomendación R4 de SECURITY-ANALYSIS.md marcadas como resueltas.
- Mención de *certificate pinning* reemplazada por TLS 1.3 + Certificate Transparency en toda la especificación. Tabla de amenazas de PROTOCOL.md actualizada.
- Entrada de glosario "Certificate Pinning" reemplazada por "Certificate Transparency (CT)" y "Clock skew".

## [0.5.0] - 2026-02-06

### Added

- Modelo formal del protocolo en Tamarin Prover (`formal/aavp.spthy`, `formal/aavp-unlinkability.spthy`). Verifica matemáticamente: unforgeability (un token válido requiere participación del IM), unlinkability (dos tokens del mismo DA son indistinguibles), unicidad de nonce, y vinculación de metadatos.
- Workflow de GitHub Actions para verificación formal (`formal-verification.yml`) con Docker (`infsec/tamarin-prover`).

## [0.4.0] - 2026-02-06

### Added

- `SECURITY-ANALYSIS.md`: estudio exhaustivo de vulnerabilidades y análisis de seguridad del protocolo, con supuestos de seguridad, vectores de ataque no documentados, análisis de esquemas criptográficos, vulnerabilidades del token, modelo de implementación para plataformas, protocolo de auditoría, verificación de segmentación de contenido, escenarios de ataque compuestos y recomendaciones priorizadas.

### Changed

- Estructura del token AAVP: adopción de Partially Blind RSA (RSAPBSSA-SHA384) como esquema criptográfico. `age_bracket` y `expires_at` pasan a ser metadatos públicos de la firma parcialmente ciega. Añadidos campos `token_type` (2 bytes, agilidad criptográfica) y `token_key_id` (32 bytes, identificación de clave del IM). Eliminado `issued_at`. Tamaño fijo del token: 331 bytes.
- Vulnerabilidades T-4.1 a T-4.5 de SECURITY-ANALYSIS.md marcadas como resueltas. Estado del área "Estructura del token" actualizado de rojo a amarillo.
- Frase sobre firmas ciegas en README.md ajustada para reflejar firmas parcialmente ciegas.
- Página `/seguridad/` en el sitio web que renderiza SECURITY-ANALYSIS.md con soporte de diagramas Mermaid.
- Enlace a análisis de seguridad en la navegación del sitio y en la landing page.
- `SECURITY-ANALYSIS.md` añadido como trigger del workflow de deploy del sitio.

## [0.3.0] - 2026-02-06

### Added

- Sitio web público con Astro desplegado en GitHub Pages (`site/`).
- Landing page con hero, sección del problema, cómo funciona, principios y CTA.
- Páginas que renderizan README.md, PROTOCOL.md y CHANGELOG.md con soporte de diagramas Mermaid (SVG en build time).
- Modo claro/oscuro con detección de preferencia del sistema y persistencia en localStorage.
- Workflow de GitHub Actions para deploy automático del sitio (`deploy-site.yml`).
- Sección de documentación del sitio web en CLAUDE.md.

### Changed

- Ortografía castellana corregida en todos los documentos: tildes, eñes y signos de puntuación.
- Regla de ortografía actualizada en CLAUDE.md: el proyecto usa ortografía castellana correcta.

## [0.2.0] - 2026-02-06

### Added

- `PROTOCOL.md`: nuevo documento de especificación técnica separado del white paper divulgativo.
- `CLAUDE.md`: guía de contribución con normas de estilo, terminología obligatoria y principios inviolables.
- `CHANGELOG.md`: registro de cambios con Semantic Versioning.
- `VERSION`: archivo fuente de verdad para el versionado del proyecto.
- `.github/workflows/bump-version.yml`: workflow de GitHub Actions para bump automático de versión, actualización de cabeceras y creación de tags.
- Distinción explícita entre Device Agent (rol del protocolo) y sus vehículos de implementación (control parental, SO, extensión, etc.).
- Tabla de vehículos de implementación del DA en PROTOCOL.md.
- Dos amenazas nuevas en el modelo de amenazas: colusión IM+plataforma y replay de tokens.
- Esquemas criptográficos candidatos documentados (RSA Blind Signatures RFC 9474, BLS, zk-SNARKs, zk-STARKs, Bulletproofs).
- Sección de prevención de fingerprinting con tabla de medidas.
- Estructura futura del repositorio documentada en CLAUDE.md.

### Changed

- **Separación de audiencias:** el contenido se divide en README.md (divulgativo, audiencia general) y PROTOCOL.md (técnico, audiencia implementadora).
- README.md reescrito con lenguaje accesible, analogías y diagramas simplificados.
- PROTOCOL.md contiene toda la arquitectura, criptografía, modelo de amenazas y flujo operativo detallado.
- Terminología corregida: Device Agent ya no se equipara con "control parental" en ningún documento.
- Definición del DA actualizada a "rol abstracto del protocolo" en lugar de "sistema de control parental o componente del SO".
- Glosario ampliado con nuevos términos: Unlinkability, vehículo de implementación.
- Diagramas Mermaid reorganizados según la audiencia de cada documento.

### Removed

- Mezcla de contenido técnico y divulgativo en un solo documento.
- Equivalencia directa DA = control parental en definiciones, glosario y flujos.

## [0.1.0] - 2026-02-06

### Added

- White paper inicial de AAVP (Anonymous Age Verification Protocol) v0.1.
- Definición del problema de verificación de edad en plataformas digitales.
- Arquitectura del protocolo con tres actores: Device Agent, Verification Gate, Implementador.
- Estructura del token AAVP con cinco campos: age_bracket, issued_at, expires_at, nonce, implementer_sig.
- Cuatro franjas de edad: UNDER_13, AGE_13_15, AGE_16_17, OVER_18.
- Fundamentos criptográficos: firmas ciegas y pruebas de conocimiento cero.
- Modelo de confianza descentralizado inspirado en DMARC/DKIM.
- Modelo de amenazas con seis vectores analizados.
- Comparativa con soluciones existentes (DNI, biometría, tarjeta de crédito, autodeclaración).
- Compatibilidad regulatoria (RGPD, DSA, COPPA, Age Appropriate Design Code UK).
- Hoja de ruta en tres fases (2026-2029+).
- Diagramas Mermaid de arquitectura, flujos, ciclo de vida del token y hoja de ruta.

[Unreleased]: https://github.com/jorgej-ramos/AAVP/compare/v0.7.0...HEAD
[0.7.0]: https://github.com/jorgej-ramos/AAVP/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/jorgej-ramos/AAVP/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/jorgej-ramos/AAVP/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/jorgej-ramos/AAVP/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/jorgej-ramos/AAVP/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/jorgej-ramos/AAVP/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/jorgej-ramos/AAVP/releases/tag/v0.1.0
