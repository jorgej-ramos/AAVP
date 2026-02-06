# Changelog

Todos los cambios notables de este proyecto se documentan en este archivo.

El formato se basa en [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/)
y este proyecto se adhiere a [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-02-06

## [0.2.0] - 2026-02-06

### Added

- `PROTOCOL.md`: nuevo documento de especificacion tecnica separado del white paper divulgativo.
- `CLAUDE.md`: guia de contribucion con normas de estilo, terminologia obligatoria y principios inviolables.
- `CHANGELOG.md`: registro de cambios con Semantic Versioning.
- `VERSION`: archivo fuente de verdad para el versionado del proyecto.
- `.github/workflows/bump-version.yml`: workflow de GitHub Actions para bump automatico de version, actualizacion de cabeceras y creacion de tags.
- Distincion explicita entre Device Agent (rol del protocolo) y sus vehiculos de implementacion (control parental, SO, extension, etc.).
- Tabla de vehiculos de implementacion del DA en PROTOCOL.md.
- Dos amenazas nuevas en el modelo de amenazas: colusion IM+plataforma y replay de tokens.
- Esquemas criptograficos candidatos documentados (RSA Blind Signatures RFC 9474, BLS, zk-SNARKs, zk-STARKs, Bulletproofs).
- Seccion de prevencion de fingerprinting con tabla de medidas.
- Estructura futura del repositorio documentada en CLAUDE.md.

### Changed

- **Separacion de audiencias:** el contenido se divide en README.md (divulgativo, audiencia general) y PROTOCOL.md (tecnico, audiencia implementadora).
- README.md reescrito con lenguaje accesible, analogias y diagramas simplificados.
- PROTOCOL.md contiene toda la arquitectura, criptografia, modelo de amenazas y flujo operativo detallado.
- Terminologia corregida: Device Agent ya no se equipara con "control parental" en ningun documento.
- Definicion del DA actualizada a "rol abstracto del protocolo" en lugar de "sistema de control parental o componente del SO".
- Glosario ampliado con nuevos terminos: Unlinkability, vehiculo de implementacion.
- Diagramas Mermaid reorganizados segun la audiencia de cada documento.

### Removed

- Mezcla de contenido tecnico y divulgativo en un solo documento.
- Equivalencia directa DA = control parental en definiciones, glosario y flujos.

## [0.1.0] - 2026-02-06

### Added

- White paper inicial de AAVP (Anonymous Age Verification Protocol) v0.1.
- Definicion del problema de verificacion de edad en plataformas digitales.
- Arquitectura del protocolo con tres actores: Device Agent, Verification Gate, Implementador.
- Estructura del token AAVP con cinco campos: age_bracket, issued_at, expires_at, nonce, implementer_sig.
- Cuatro franjas de edad: UNDER_13, AGE_13_15, AGE_16_17, OVER_18.
- Fundamentos criptograficos: firmas ciegas y pruebas de conocimiento cero.
- Modelo de confianza descentralizado inspirado en DMARC/DKIM.
- Modelo de amenazas con seis vectores analizados.
- Comparativa con soluciones existentes (DNI, biometria, tarjeta de credito, autodeclaracion).
- Compatibilidad regulatoria (RGPD, DSA, COPPA, Age Appropriate Design Code UK).
- Hoja de ruta en tres fases (2026-2029+).
- Diagramas Mermaid de arquitectura, flujos, ciclo de vida del token y hoja de ruta.

[Unreleased]: https://github.com/USER/AAVP/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/jorgej-ramos/AAVP/compare/v0.1.0...v0.2.0
[0.2.0]: https://github.com/USER/AAVP/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/USER/AAVP/releases/tag/v0.1.0
