# Changelog

Todos los cambios notables de este proyecto se documentan en este archivo.

El formato se basa en [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/)
y este proyecto se adhiere a [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.4.0] - 2026-02-06

### Added

- `SECURITY-ANALYSIS.md`: estudio exhaustivo de vulnerabilidades y análisis de seguridad del protocolo, con supuestos de seguridad, vectores de ataque no documentados, análisis de esquemas criptográficos, vulnerabilidades del token, modelo de implementación para plataformas, protocolo de auditoría, verificación de segmentación de contenido, escenarios de ataque compuestos y recomendaciones priorizadas.
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

[Unreleased]: https://github.com/jorgej-ramos/AAVP/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/jorgej-ramos/AAVP/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/jorgej-ramos/AAVP/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/jorgej-ramos/AAVP/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/jorgej-ramos/AAVP/releases/tag/v0.1.0
