# CLAUDE.md — Guia de contribucion para AAVP

## Que es este proyecto

AAVP (Anonymous Age Verification Protocol) es una propuesta de **protocolo abierto y descentralizado** para la verificacion anonima de edad en plataformas digitales. El objetivo final es publicar una especificacion formal en formato RFC y someterla a estandarizacion (W3C / IETF).

Este repositorio es el punto de partida: contiene el white paper (README.md) y la especificacion tecnica (PROTOCOL.md). Con el tiempo contendra la especificacion formal, implementaciones de referencia, test vectors y documentacion de gobernanza.

---

## Estructura del repositorio

```
README.md                          Documento divulgativo (audiencia general)
PROTOCOL.md                        Especificacion tecnica (audiencia tecnica)
CLAUDE.md                          Esta guia (normas de estilo y contribucion)
CHANGELOG.md                       Registro de cambios (Keep a Changelog 1.1.0)
VERSION                            Fuente de verdad para la version (semver)
.github/workflows/bump-version.yml Workflow de bump automatico de version
```

### Regla de separacion de audiencias

El proyecto mantiene **dos niveles de documentacion** con audiencias distintas. Nunca mezclar jerga tecnica en el documento divulgativo ni simplificaciones imprecisas en el tecnico.

| Documento | Audiencia | Tono | Contenido |
|-----------|-----------|------|-----------|
| `README.md` | Reguladores, padres, prensa, plataformas, publico general | Accesible, directo, con analogias | Que problema resuelve, como funciona a alto nivel, por que es diferente, comparativas, regulacion |
| `PROTOCOL.md` | Desarrolladores, criptografos, auditores, implementadores | Riguroso, preciso, sin ambiguedad | Arquitectura, token, criptografia, modelo de amenazas, flujos, glosario tecnico |

**Cuando anadir contenido, preguntarse siempre:** "Esta informacion, va dirigida a alguien que necesita entender el *que* y el *por que* (README), o a alguien que necesita implementar el *como* (PROTOCOL)?"

---

## Terminologia obligatoria

### Distincion critica: Device Agent vs. Control Parental

Esta es la convencion mas importante del proyecto. **Device Agent no es sinonimo de control parental.**

| Termino | Que es | Uso correcto |
|---------|--------|-------------|
| **Device Agent (DA)** | Rol abstracto del protocolo AAVP. Componente de software que genera y gestiona tokens de edad. | "El Device Agent genera un token efimero" |
| **Control parental** | Un tipo de producto que *puede* implementar el rol de DA, pero no es el unico. | "El control parental es uno de los vehiculos posibles para implementar el DA" |
| **Vehiculo de implementacion** | Cualquier software que actue como DA: control parental, componente del SO, extension de navegador, etc. | "Los vehiculos de implementacion del DA incluyen..." |

**Nunca escribir:**
- ~~"El control parental (Device Agent)..."~~
- ~~"El DA, es decir, el sistema de control parental..."~~
- ~~"Software de control parental o componente del SO" como definicion del DA~~

**Si escribir:**
- "El Device Agent — que puede ser implementado por un sistema de control parental, un componente del SO u otro software conforme — genera..."
- "El software que actua como Device Agent..."

### Terminologia del protocolo

Usar siempre los terminos canonicos para los tres roles:

| Rol | Abreviatura | Nunca llamarlo |
|-----|-------------|----------------|
| Device Agent | DA | "el agente", "el cliente", "la app" |
| Verification Gate | VG | "el servidor", "el endpoint", "la API" |
| Implementador | IM | "el proveedor", "la empresa", "el emisor" |

### Franjas de edad

Usar siempre los codigos canonicos en contexto tecnico:

```
UNDER_13   AGE_13_15   AGE_16_17   OVER_18
```

En contexto divulgativo usar las etiquetas descriptivas: "Infantil", "Adolescente temprano", "Adolescente tardio", "Adulto".

### Otros terminos

| Preferir | En lugar de |
|----------|-------------|
| credencial anonima / token | sello, certificado, ticket |
| franja de edad | rango, grupo, categoria de edad |
| firma ciega | firma enmascarada, firma oculta |
| vehiculo de implementacion | plataforma de implementacion, medio |
| senal de edad | dato de edad, informacion de edad |

---

## Principios inviolables del protocolo

Estos cuatro principios son los pilares de AAVP. **Ninguna propuesta, modificacion o extension puede comprometer ninguno de ellos.** Si una idea entra en conflicto con alguno, la idea se descarta, no el principio.

1. **Privacidad por Diseno.** Ningun dato personal identificable abandona el dispositivo. Garantia matematica, no politica.
2. **Descentralizacion.** Sin autoridad central. Cada plataforma decide en quien confiar.
3. **Estandar Abierto.** Sin licencias, tasas ni permisos. Cualquiera puede implementar.
4. **Minimalismo de Datos.** Solo franja de edad. Cada campo adicional requiere justificacion rigurosa.

Al revisar o redactar contenido, verificar siempre que no se introduce lenguaje que contradiga o debilite estos principios. Ejemplos de violaciones sutiles:

- "La autoridad AAVP podria..." — No existe tal autoridad.
- "El token incluye un identificador de sesion..." — Viola minimalismo de datos.
- "El IM conoce la franja del token firmado..." — Viola firmas ciegas / privacidad.
- "Se requiere registro previo en..." — Viola estandar abierto.

---

## Estilo de escritura

### Idioma

- El proyecto se redacta en **espanol**.
- Los terminos tecnicos del protocolo se mantienen en **ingles**: Device Agent, Verification Gate, Implementador (este ultimo ya es espanol), blind signature, token, nonce, TTL, ZKP.
- Los nombres de campos del token van siempre en ingles y monospace: `age_bracket`, `issued_at`, `expires_at`, `nonce`, `implementer_sig`.

### Tono general

- **Serio y profesional**, pero no academico ni burocratico.
- **Directo.** Frases cortas. Ir al grano. Evitar rodeos y muletillas.
- **Honesto.** Reconocer limitaciones abiertamente. No prometer mas de lo que el protocolo puede entregar.
- **Preciso.** Cada afirmacion debe ser defendible. No hacer afirmaciones criptograficas vagas.
- **Sin sensacionalismo.** No usar superlativos innecesarios ni lenguaje de marketing.

### Tono por documento

**README.md (divulgativo):**
- Usar analogias cuando ayuden (portero de discoteca, sobre con papel carbon, SMTP para correo).
- Tutear o usar formas impersonales. Nunca "usted".
- Explicar conceptos como si el lector fuera inteligente pero no tecnico.
- Evitar acronimos no explicados. La primera vez que aparece un acronimo, explicarlo.
- Los diagramas Mermaid deben ser simples y autoexplicativos.

**PROTOCOL.md (tecnico):**
- Ser preciso y completo. No simplificar en exceso.
- Usar terminologia tecnica sin disculparse, pero definir todo en el glosario.
- Los diagramas Mermaid pueden ser detallados.
- Las afirmaciones criptograficas deben ser verificables o estar marcadas como candidatas/pendientes de evaluacion.
- Incluir esquemas candidatos con sus nombres formales (RFC, nombre del algoritmo).

### Formato Markdown

- Encabezados: `##` para secciones principales, `###` para subsecciones, `####` para sub-subsecciones. Nunca usar `#` excepto para el titulo del documento.
- Tablas: usar para comparativas, listas de campos y cualquier informacion tabular. Alinear las columnas visualmente.
- Listas: usar `-` para listas no ordenadas, `1.` para secuencias con orden.
- Enfasis: **negrita** para terminos clave y conceptos criticos. *Cursiva* para terminos en ingles no canonicos y citas conceptuales. No abusar de ninguno de los dos.
- Codigo inline: usar backticks para campos del token, endpoints, valores de enumeracion y fragmentos de codigo.
- Callouts de GitHub: usar `> [!NOTE]`, `> [!TIP]`, `> [!IMPORTANT]` cuando corresponda. No abusar.
- Separadores: usar `---` entre secciones principales.
- Sin emojis.

### Diagramas Mermaid

- Usar Mermaid para todos los diagramas. No imagenes externas.
- Tipos preferidos: `sequenceDiagram` para flujos, `graph` para arquitectura, `stateDiagram-v2` para ciclos de vida, `classDiagram` para estructuras de datos, `gantt` para hojas de ruta, `flowchart` para decisiones.
- Etiquetas en los diagramas: usar los nombres canonicos del protocolo (DA, VG, IM) con su nombre completo entre corchetes. Ejemplo: `DA[Device Agent]`.
- En README.md los diagramas deben ser comprensibles sin contexto tecnico.
- En PROTOCOL.md los diagramas pueden incluir detalle tecnico (funciones criptograficas, parametros).

### Tablas de acentos

El proyecto evita tildes y caracteres especiales en el cuerpo del texto Markdown para maximizar la portabilidad y la legibilidad en editores sin soporte Unicode completo. Esto aplica al texto en prosa, no a citas literales o nombres propios que lo requieran.

---

## Trayectoria hacia RFC

El proyecto tiene como objetivo final una especificacion formal en formato RFC. Esto implica las siguientes consideraciones al escribir:

### Estructura progresiva

1. **Fase actual (white paper):** README.md + PROTOCOL.md. Lenguaje natural con rigor tecnico.
2. **Fase siguiente (Internet-Draft):** Se creara un directorio `spec/` con la especificacion en formato I-D (Internet-Draft) siguiendo las convenciones de la IETF.
3. **Fase final (RFC):** Sometimiento formal al proceso de estandarizacion.

### Preparacion para el I-D

Al redactar PROTOCOL.md, tener en cuenta que su contenido sera la base del Internet-Draft:

- Definir terminos con precision. El glosario de PROTOCOL.md evolucionara hacia la seccion "Terminology" del I-D.
- Usar lenguaje que se pueda mapear a los requisitos RFC 2119 (MUST, SHOULD, MAY). Aunque no usemos esas palabras todavia, las afirmaciones deben ser lo suficientemente precisas para traducirse a ellas.
- Documentar cada decision de diseno con su justificacion. Los RFCs exigen explicar por que se eligio una alternativa sobre las demas.
- Incluir consideraciones de seguridad exhaustivas. El modelo de amenazas de PROTOCOL.md sera la base de la seccion "Security Considerations".
- Los esquemas criptograficos que aun no estan decididos deben marcarse como "candidatos" y documentar los criterios de evaluacion.

---

## Versionado y changelog

### Semantic Versioning

El proyecto sigue [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html) adaptado a un proyecto de especificacion de protocolo (no es una libreria de software, pero los principios aplican):

| Bump | Cuando | Ejemplo |
|------|--------|---------|
| **MAJOR** | Cambios incompatibles con versiones anteriores de la especificacion. Rompe implementaciones existentes. | Cambiar la estructura del token, eliminar un rol del protocolo, modificar el modelo de confianza |
| **MINOR** | Adiciones compatibles hacia atras. No rompe implementaciones existentes. | Nueva seccion en el white paper, extension opcional del token, nuevo vehiculo de implementacion documentado |
| **PATCH** | Correcciones que no alteran el significado de la especificacion. | Typos, clarificaciones de redaccion, mejoras de diagramas, correcciones de formato |

**Reglas clave:**
- Mientras la version sea `0.x.y` (desarrollo inicial), el proyecto esta en fase de borrador y todo puede cambiar.
- `v1.0.0` se reserva para la primera version estable de la especificacion formal (Internet-Draft).
- La version actual se almacena en el archivo `VERSION` (fuente unica de verdad).

### Changelog

El registro de cambios sigue [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/):

- Archivo: `CHANGELOG.md` en la raiz del repositorio.
- Secciones por tipo de cambio: `Added`, `Changed`, `Deprecated`, `Removed`, `Fixed`, `Security`.
- Los cambios pendientes de release se acumulan bajo `[Unreleased]`.
- Al hacer un bump de version, la seccion `[Unreleased]` se estampa automaticamente con la nueva version y fecha.

**Regla para contribuciones:** todo cambio notable debe anadirse manualmente bajo `[Unreleased]` en CHANGELOG.md como parte del commit que introduce el cambio. El workflow de bump se encarga del resto.

### Release automatico via tag (GitHub Actions)

El release se dispara **automaticamente al pushear un tag semver** a main. No hay que recordar ejecutar ningun workflow manualmente.

**Como hacer un release:**

```bash
git tag v0.3.0 && git push origin v0.3.0
```

Eso es todo. El workflow `.github/workflows/bump-version.yml` se encarga del resto.

**Que hace el workflow al detectar el tag:**

1. Extrae la version del nombre del tag (ej: `v0.3.0` -> `0.3.0`).
2. Valida que la version es mayor que la actual en `VERSION` (rechaza tags iguales o menores).
3. Actualiza `VERSION` con la nueva version.
4. Actualiza las cabeceras de `README.md` y `PROTOCOL.md`.
5. Estampa la seccion `[Unreleased]` de `CHANGELOG.md` con la version y fecha, y crea una nueva seccion `[Unreleased]` vacia.
6. Actualiza los links de comparacion al final de `CHANGELOG.md`.
7. Hace commit `Release vX.Y.Z` en main.
8. Mueve el tag para que apunte al commit final (con cabeceras ya actualizadas).

**Flujo de trabajo completo:**

```
1. Hacer cambios en README.md, PROTOCOL.md, etc.
2. Anadir entrada bajo [Unreleased] en CHANGELOG.md
3. Commit y push a main
4. ... repetir 1-3 cuantas veces sea necesario ...
5. Cuando se quiera hacer un release:
   git tag v0.3.0 && git push origin v0.3.0
6. El workflow commitea, actualiza cabeceras y recoloca el tag automaticamente
```

**Validaciones del workflow:**

- El tag debe seguir el formato `vMAJOR.MINOR.PATCH` (ej: `v0.3.0`, `v1.0.0`).
- La version del tag debe ser estrictamente mayor que la version actual en `VERSION`.
- Si la validacion falla, el workflow aborta y muestra un error en la pestaña Actions.

**Patrones de version que el workflow busca y reemplaza:**

- `README.md`: `White Paper vX.Y.Z` y `Anonymous Age Verification Protocol · vX.Y.Z`
- `PROTOCOL.md`: `vX.Y.Z — Borrador` y `Especificacion Tecnica · vX.Y.Z`

Si se anade un nuevo documento con version en la cabecera, hay que anadir el patron `sed` correspondiente al workflow.

---

## Guia de contribucion al contenido

### Antes de anadir un campo al token

Cualquier propuesta de anadir un campo al token AAVP debe superar este test:

1. **Necesidad:** Es estrictamente necesario para el funcionamiento del protocolo?
2. **Minimalismo:** Puede lograrse el mismo objetivo sin este campo?
3. **Fingerprinting:** Puede este campo, solo o combinado con otros, usarse para identificar o rastrear al usuario?
4. **Unlinkability:** Compromete la imposibilidad de correlacionar dos tokens del mismo usuario?

Si la respuesta a 3 o 4 es "si" o "posiblemente", el campo se rechaza.

### Antes de proponer un cambio arquitectonico

Verificar que no viola ninguno de los cuatro principios inviolables. Si el cambio introduce una dependencia en una entidad central, un dato personal o una barrera de acceso, no es compatible con AAVP.

### Antes de anadir una seccion a un documento

Preguntarse:
- Esta informacion existe ya en el otro documento? (Evitar duplicacion.)
- A que audiencia va dirigida? (Colocarla en el documento correcto.)
- Anade valor o es relleno? (Solo contenido que aporte.)

### Verificacion de coherencia

Al modificar un concepto que aparece en ambos documentos, actualizar los dos. El README.md y el PROTOCOL.md deben ser coherentes entre si en todo momento, aunque a distinto nivel de detalle.

---

## Estructura futura prevista del repositorio

```
README.md                              White paper divulgativo
PROTOCOL.md                            Especificacion tecnica
CLAUDE.md                              Guia de contribucion (este archivo)
CHANGELOG.md                           Registro de cambios
VERSION                                Version actual (fuente de verdad)
LICENSE                                Licencia del proyecto
.github/
  workflows/
    bump-version.yml                   Workflow de bump automatico
spec/
  draft-aavp-protocol.md              Internet-Draft (futuro)
reference/
  da/                                  Implementacion de referencia del DA (futuro)
  vg/                                  Implementacion de referencia del VG (futuro)
test-vectors/
  tokens/                              Vectores de test para validacion (futuro)
docs/
  threat-model.md                      Modelo de amenazas extendido (futuro)
  governance.md                        Modelo de gobernanza (futuro)
```

---

## Reglas para commits

- Mensajes de commit en ingles, concisos, en imperativo: "Add threat model for token replay", "Fix DA terminology in README".
- Un commit por cambio logico. No mezclar cambios en README.md y PROTOCOL.md en el mismo commit salvo que sean la misma correccion de coherencia.
- Nunca commitear archivos generados, binarios ni secretos.
- **No modificar VERSION manualmente.** El workflow de GitHub Actions se encarga de ello.
- **No estampar [Unreleased] manualmente en CHANGELOG.md.** Solo anadir entradas bajo esa seccion. El workflow se encarga de estampar la version y la fecha.
- Los commits de release (`Release vX.Y.Z`) son creados exclusivamente por el workflow.
- Para hacer un release: `git tag vX.Y.Z && git push origin vX.Y.Z`. El tag dispara el workflow automaticamente.
- **Sin firmas ni trailers en los commits.** No incluir `Co-Authored-By`, `Signed-off-by`, `Generated by` ni ninguna otra formula de autoria o co-autoria. Los mensajes de commit deben contener exclusivamente la descripcion del cambio.

---

## Resumen de verificacion rapida

Antes de dar por bueno cualquier cambio, verificar:

- [ ] Device Agent no se equipara con control parental
- [ ] Los cuatro principios inviolables no se contradicen
- [ ] La terminologia usa los nombres canonicos del protocolo
- [ ] El contenido esta en el documento correcto segun su audiencia
- [ ] Los diagramas Mermaid son coherentes con el texto
- [ ] No se introducen campos en el token sin justificacion
- [ ] README.md y PROTOCOL.md son coherentes entre si
- [ ] Sin emojis, sin superlativos, sin lenguaje de marketing
- [ ] Sin tildes ni caracteres especiales en el cuerpo del texto
