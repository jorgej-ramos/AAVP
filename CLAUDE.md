# CLAUDE.md — Guía de contribución para AAVP

## Qué es este proyecto

AAVP (Anonymous Age Verification Protocol) es una propuesta de **protocolo abierto y descentralizado** para la verificación anónima de edad en plataformas digitales. El objetivo final es publicar una especificación formal en formato RFC y someterla a estandarización (W3C / IETF).

Este repositorio es el punto de partida: contiene el white paper (README.md) y la especificación técnica (PROTOCOL.md). Con el tiempo contendrá la especificación formal, implementaciones de referencia, test vectors y documentación de gobernanza.

---

## Estructura del repositorio

```
README.md                          Documento divulgativo (audiencia general)
PROTOCOL.md                        Especificacion tecnica (audiencia tecnica)
SECURITY-ANALYSIS.md               Estudio de vulnerabilidades y analisis de seguridad
CLAUDE.md                          Esta guia (normas de estilo y contribucion)
CHANGELOG.md                       Registro de cambios (Keep a Changelog 1.1.0)
VERSION                            Fuente de verdad para la version (semver)
.github/workflows/bump-version.yml Workflow de bump automatico de version
.github/workflows/deploy-site.yml  Workflow de deploy del sitio web a GitHub Pages
.github/workflows/spec-consistency.yml Workflow de coherencia de la especificacion
scripts/check-spec-consistency.sh  Script de verificacion de coherencia entre documentos
site/                              Sitio web publico (Astro, desplegado en GitHub Pages)
```

### Regla de separación de audiencias

El proyecto mantiene **dos niveles de documentación** con audiencias distintas. Nunca mezclar jerga técnica en el documento divulgativo ni simplificaciones imprecisas en el técnico.

| Documento | Audiencia | Tono | Contenido |
|-----------|-----------|------|-----------|
| `README.md` | Reguladores, padres, prensa, plataformas, público general | Accesible, directo, con analogías | Qué problema resuelve, cómo funciona a alto nivel, por qué es diferente, comparativas, regulación |
| `PROTOCOL.md` | Desarrolladores, criptógrafos, auditores, implementadores | Riguroso, preciso, sin ambigüedad | Arquitectura, token, criptografía, modelo de amenazas, flujos, glosario técnico |

**Cuando añadir contenido, preguntarse siempre:** "¿Esta información, va dirigida a alguien que necesita entender el *qué* y el *por qué* (README), o a alguien que necesita implementar el *cómo* (PROTOCOL)?"

---

## Terminología obligatoria

### Distinción crítica: Device Agent vs. Control Parental

Esta es la convención más importante del proyecto. **Device Agent no es sinónimo de control parental.**

| Término | Qué es | Uso correcto |
|---------|--------|-------------|
| **Device Agent (DA)** | Rol abstracto del protocolo AAVP. Componente de software que genera y gestiona tokens de edad. | "El Device Agent genera un token efímero" |
| **Control parental** | Un tipo de producto que *puede* implementar el rol de DA, pero no es el único. | "El control parental es uno de los vehículos posibles para implementar el DA" |
| **Vehículo de implementación** | Cualquier software que actúe como DA: control parental, componente del SO, extensión de navegador, etc. | "Los vehículos de implementación del DA incluyen..." |

**Nunca escribir:**
- ~~"El control parental (Device Agent)..."~~
- ~~"El DA, es decir, el sistema de control parental..."~~
- ~~"Software de control parental o componente del SO" como definición del DA~~

**Sí escribir:**
- "El Device Agent — que puede ser implementado por un sistema de control parental, un componente del SO u otro software conforme — genera..."
- "El software que actúa como Device Agent..."

### Terminología del protocolo

Usar siempre los términos canónicos para los tres roles:

| Rol | Abreviatura | Nunca llamarlo |
|-----|-------------|----------------|
| Device Agent | DA | "el agente", "el cliente", "la app" |
| Verification Gate | VG | "el servidor", "el endpoint", "la API" |
| Implementador | IM | "el proveedor", "la empresa", "el emisor" |

### Franjas de edad

Usar siempre los códigos canónicos en contexto técnico:

```
UNDER_13   AGE_13_15   AGE_16_17   OVER_18
```

En contexto divulgativo usar las etiquetas descriptivas: "Infantil", "Adolescente temprano", "Adolescente tardío", "Adulto".

### Otros términos

| Preferir | En lugar de |
|----------|-------------|
| credencial anónima / token | sello, certificado, ticket |
| franja de edad | rango, grupo, categoría de edad |
| firma parcialmente ciega | firma enmascarada, firma oculta |
| vehículo de implementación | plataforma de implementación, medio |
| señal de edad | dato de edad, información de edad |

---

## Principios inviolables del protocolo

Estos cuatro principios son los pilares de AAVP. **Ninguna propuesta, modificación o extensión puede comprometer ninguno de ellos.** Si una idea entra en conflicto con alguno, la idea se descarta, no el principio.

1. **Privacidad por Diseño.** Ningún dato personal identificable abandona el dispositivo. Garantía matemática, no política.
2. **Descentralización.** Sin autoridad central. Cada plataforma decide en quién confiar.
3. **Estándar Abierto.** Sin licencias, tasas ni permisos. Cualquiera puede implementar.
4. **Minimalismo de Datos.** Solo franja de edad. Cada campo adicional requiere justificación rigurosa.

Al revisar o redactar contenido, verificar siempre que no se introduce lenguaje que contradiga o debilite estos principios. Ejemplos de violaciones sutiles:

- "La autoridad AAVP podría..." — No existe tal autoridad.
- "El token incluye un identificador de sesión..." — Viola minimalismo de datos.
- "El IM conoce la identidad del usuario que solicita la firma..." — Viola firmas parcialmente ciegas / privacidad. (Nota: el IM sí conoce la franja de edad como metadato público, pero no la identidad del usuario.)
- "Se requiere registro previo en..." — Viola estándar abierto.

---

## Estilo de escritura

### Idioma

- El proyecto se redacta en **español**.
- Los términos técnicos del protocolo se mantienen en **inglés**: Device Agent, Verification Gate, Implementador (este último ya es español), blind signature, token, nonce, TTL, ZKP.
- Los nombres de campos del token van siempre en inglés y monospace: `token_type`, `nonce`, `token_key_id`, `age_bracket`, `expires_at`, `authenticator`.

### Tono general

- **Serio y profesional**, pero no académico ni burocrático.
- **Directo.** Frases cortas. Ir al grano. Evitar rodeos y muletillas.
- **Honesto.** Reconocer limitaciones abiertamente. No prometer más de lo que el protocolo puede entregar.
- **Preciso.** Cada afirmación debe ser defendible. No hacer afirmaciones criptográficas vagas.
- **Sin sensacionalismo.** No usar superlativos innecesarios ni lenguaje de marketing.

### Tono por documento

**README.md (divulgativo):**
- Usar analogías cuando ayuden (portero de discoteca, sobre con papel carbón, SMTP para correo).
- Tutear o usar formas impersonales. Nunca "usted".
- Explicar conceptos como si el lector fuera inteligente pero no técnico.
- Evitar acrónimos no explicados. La primera vez que aparece un acrónimo, explicarlo.
- Los diagramas Mermaid deben ser simples y autoexplicativos.

**PROTOCOL.md (técnico):**
- Ser preciso y completo. No simplificar en exceso.
- Usar terminología técnica sin disculparse, pero definir todo en el glosario.
- Los diagramas Mermaid pueden ser detallados.
- Las afirmaciones criptográficas deben ser verificables o estar marcadas como candidatas/pendientes de evaluación.
- Incluir esquemas candidatos con sus nombres formales (RFC, nombre del algoritmo).

### Formato Markdown

- Encabezados: `##` para secciones principales, `###` para subsecciones, `####` para sub-subsecciones. Nunca usar `#` excepto para el título del documento.
- Tablas: usar para comparativas, listas de campos y cualquier información tabular. Alinear las columnas visualmente.
- Listas: usar `-` para listas no ordenadas, `1.` para secuencias con orden.
- Énfasis: **negrita** para términos clave y conceptos críticos. *Cursiva* para términos en inglés no canónicos y citas conceptuales. No abusar de ninguno de los dos.
- Código inline: usar backticks para campos del token, endpoints, valores de enumeración y fragmentos de código.
- Callouts de GitHub: usar `> [!NOTE]`, `> [!TIP]`, `> [!IMPORTANT]` cuando corresponda. No abusar.
- Separadores: usar `---` entre secciones principales.
- Sin emojis.

### Diagramas Mermaid

- Usar Mermaid para todos los diagramas. No imágenes externas.
- Tipos preferidos: `sequenceDiagram` para flujos, `graph` para arquitectura, `stateDiagram-v2` para ciclos de vida, `classDiagram` para estructuras de datos, `gantt` para hojas de ruta, `flowchart` para decisiones.
- Etiquetas en los diagramas: usar los nombres canónicos del protocolo (DA, VG, IM) con su nombre completo entre corchetes. Ejemplo: `DA[Device Agent]`.
- En README.md los diagramas deben ser comprensibles sin contexto técnico.
- En PROTOCOL.md los diagramas pueden incluir detalle técnico (funciones criptográficas, parámetros).

### Ortografía

El proyecto se redacta con ortografía castellana correcta: tildes, eñes, signos de apertura de interrogación y exclamación. No se omiten caracteres especiales del español.

---

## Trayectoria hacia RFC

El proyecto tiene como objetivo final una especificación formal en formato RFC. Esto implica las siguientes consideraciones al escribir:

### Estructura progresiva

1. **Fase actual (white paper):** README.md + PROTOCOL.md. Lenguaje natural con rigor técnico.
2. **Fase siguiente (Internet-Draft):** Se creará un directorio `spec/` con la especificación en formato I-D (Internet-Draft) siguiendo las convenciones de la IETF.
3. **Fase final (RFC):** Sometimiento formal al proceso de estandarización.

### Preparación para el I-D

Al redactar PROTOCOL.md, tener en cuenta que su contenido será la base del Internet-Draft:

- Definir términos con precisión. El glosario de PROTOCOL.md evolucionará hacia la sección "Terminology" del I-D.
- Usar lenguaje que se pueda mapear a los requisitos RFC 2119 (MUST, SHOULD, MAY). Aunque no usemos esas palabras todavía, las afirmaciones deben ser lo suficientemente precisas para traducirse a ellas.
- Documentar cada decisión de diseño con su justificación. Los RFCs exigen explicar por qué se eligió una alternativa sobre las demás.
- Incluir consideraciones de seguridad exhaustivas. El modelo de amenazas de PROTOCOL.md será la base de la sección "Security Considerations".
- Los esquemas criptográficos que aún no están decididos deben marcarse como "candidatos" y documentar los criterios de evaluación.

---

## Versionado y changelog

### Semantic Versioning

El proyecto sigue [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html) adaptado a un proyecto de especificación de protocolo (no es una librería de software, pero los principios aplican):

| Bump | Cuándo | Ejemplo |
|------|--------|---------|
| **MAJOR** | Cambios incompatibles con versiones anteriores de la especificación. Rompe implementaciones existentes. | Cambiar la estructura del token, eliminar un rol del protocolo, modificar el modelo de confianza |
| **MINOR** | Adiciones compatibles hacia atrás. No rompe implementaciones existentes. | Nueva sección en el white paper, extensión opcional del token, nuevo vehículo de implementación documentado |
| **PATCH** | Correcciones que no alteran el significado de la especificación. | Typos, clarificaciones de redacción, mejoras de diagramas, correcciones de formato |

**Reglas clave:**
- Mientras la versión sea `0.x.y` (desarrollo inicial), el proyecto está en fase de borrador y todo puede cambiar.
- `v1.0.0` se reserva para la primera versión estable de la especificación formal (Internet-Draft).
- La versión actual se almacena en el archivo `VERSION` (fuente única de verdad).

### Changelog

El registro de cambios sigue [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/):

- Archivo: `CHANGELOG.md` en la raíz del repositorio.
- Secciones por tipo de cambio: `Added`, `Changed`, `Deprecated`, `Removed`, `Fixed`, `Security`.
- Los cambios pendientes de release se acumulan bajo `[Unreleased]`.
- Al hacer un bump de versión, la sección `[Unreleased]` se estampa automáticamente con la nueva versión y fecha.

**Regla para contribuciones:** todo cambio notable debe añadirse manualmente bajo `[Unreleased]` en CHANGELOG.md como parte del commit que introduce el cambio. El workflow de bump se encarga del resto.

### Release automático vía tag (GitHub Actions)

El release se dispara **automáticamente al pushear un tag semver** a main. No hay que recordar ejecutar ningún workflow manualmente.

**Cómo hacer un release:**

```bash
git tag v0.3.0 && git push origin v0.3.0
```

Eso es todo. El workflow `.github/workflows/bump-version.yml` se encarga del resto.

**Qué hace el workflow al detectar el tag:**

1. Extrae la versión del nombre del tag (ej: `v0.3.0` -> `0.3.0`).
2. Valida que la versión es mayor que la actual en `VERSION` (rechaza tags iguales o menores).
3. Actualiza `VERSION` con la nueva versión.
4. Actualiza las cabeceras de `README.md` y `PROTOCOL.md`.
5. Estampa la sección `[Unreleased]` de `CHANGELOG.md` con la versión y fecha, y crea una nueva sección `[Unreleased]` vacía.
6. Actualiza los links de comparación al final de `CHANGELOG.md`.
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
- La versión del tag debe ser estrictamente mayor que la versión actual en `VERSION`.
- Si la validación falla, el workflow aborta y muestra un error en la pestaña Actions.

**Patrones de versión que el workflow busca y reemplaza:**

- `README.md`: `White Paper vX.Y.Z` y `Anonymous Age Verification Protocol · vX.Y.Z`
- `PROTOCOL.md`: `vX.Y.Z — Borrador` y `Especificacion Tecnica · vX.Y.Z`

Si se añade un nuevo documento con versión en la cabecera, hay que añadir el patrón `sed` correspondiente al workflow.

---

## Guía de contribución al contenido

### Antes de añadir un campo al token

Cualquier propuesta de añadir un campo al token AAVP debe superar este test:

1. **Necesidad:** ¿Es estrictamente necesario para el funcionamiento del protocolo?
2. **Minimalismo:** ¿Puede lograrse el mismo objetivo sin este campo?
3. **Fingerprinting:** ¿Puede este campo, solo o combinado con otros, usarse para identificar o rastrear al usuario?
4. **Unlinkability:** ¿Compromete la imposibilidad de correlacionar dos tokens del mismo usuario?

Si la respuesta a 3 o 4 es "sí" o "posiblemente", el campo se rechaza.

### Antes de proponer un cambio arquitectónico

Verificar que no viola ninguno de los cuatro principios inviolables. Si el cambio introduce una dependencia en una entidad central, un dato personal o una barrera de acceso, no es compatible con AAVP.

### Antes de añadir una sección a un documento

Preguntarse:
- ¿Esta información existe ya en el otro documento? (Evitar duplicación.)
- ¿A qué audiencia va dirigida? (Colocarla en el documento correcto.)
- ¿Añade valor o es relleno? (Solo contenido que aporte.)

### Al modificar PROTOCOL.md: re-evaluación contra SECURITY-ANALYSIS.md

Cada cambio en la especificación técnica puede resolver, agravar o introducir vulnerabilidades. **Antes de dar por bueno un cambio en PROTOCOL.md**, consultar el resumen ejecutivo de `SECURITY-ANALYSIS.md` y verificar:

1. **¿El cambio mitiga alguna vulnerabilidad documentada?** Si es así, actualizar la severidad o el riesgo residual de la entrada correspondiente en SECURITY-ANALYSIS.md. Si la vulnerabilidad queda resuelta, indicarlo explícitamente.
2. **¿El cambio introduce una nueva superficie de ataque?** Evaluar si crea nuevos vectores no contemplados. Si es así, documentarlos en la sección correspondiente de SECURITY-ANALYSIS.md.
3. **¿El cambio afecta a algún supuesto de seguridad (S1-S14)?** Verificar si fortalece o debilita alguno de los supuestos explícitos o implícitos listados en la sección 1.
4. **¿El cambio modifica la estructura del token?** Revisar las vulnerabilidades T-4.1 a T-4.6 para verificar que no se introducen nuevas carencias de especificación.

**Ejemplo:** Si se define el formato binario del token (331 bytes fijos), las vulnerabilidades T-4.1 (formato no definido) y T-4.2 (tamaño no especificado) pasan a estado resuelto y deben actualizarse.

**Regla práctica:** Todo commit que modifique PROTOCOL.md debe ir acompañado de una revisión del resumen ejecutivo de SECURITY-ANALYSIS.md. Si el cambio afecta a alguna entrada, el commit debe incluir también la actualización correspondiente en SECURITY-ANALYSIS.md.

### Verificación de coherencia

Al modificar un concepto que aparece en varios documentos, actualizar todos los afectados. README.md, PROTOCOL.md y SECURITY-ANALYSIS.md deben ser coherentes entre sí en todo momento, aunque a distinto nivel de detalle.

- **README.md ↔ PROTOCOL.md:** Coherencia de conceptos a diferente nivel de audiencia.
- **PROTOCOL.md ↔ SECURITY-ANALYSIS.md:** Coherencia entre la especificación y el análisis de sus vulnerabilidades. Si PROTOCOL.md cambia, las vulnerabilidades documentadas en SECURITY-ANALYSIS.md pueden cambiar de severidad o quedar resueltas.

### Al modificar cualquier documento: actualizar la verificación automatizada

El proyecto mantiene un script de verificación automatizada (`scripts/check-spec-consistency.sh`) que se ejecuta en CI con cada push y pull request. **Si un cambio en la especificación invalida o requiere ampliar algún check, el commit debe incluir la actualización correspondiente del script.**

Situaciones que requieren actualizar el script:

1. **Cambio en la estructura del token.** Si se añade, elimina o modifica un campo:
   - Actualizar la lista `CANONICAL_FIELDS` con los campos activos.
   - Actualizar la lista `REMOVED_FIELDS` si se elimina un campo.
   - La verificación de offsets y tamaño total se recalcula automáticamente desde la tabla de PROTOCOL.md.

2. **Cambio en el tamaño del token.** Si el tamaño total cambia:
   - Actualizar la variable `TOKEN_SIZE` en el script.

3. **Cambio en el esquema criptográfico.** Si se adopta un esquema diferente:
   - Actualizar los checks de la sección 5 (actualmente buscan `RSAPBSSA-SHA384` y `RFC 9474`).

4. **Cambio en las franjas de edad.** Si se añade o modifica una franja:
   - Actualizar la lista `AGE_BRACKETS`.

5. **Cambio en el semáforo de SECURITY-ANALYSIS.md.** Si se cambia el estado de un área (rojo/amarillo/verde):
   - El script verifica automáticamente que el resumen textual coincide con los iconos de la tabla. Solo hay que asegurarse de que ambos se actualizan en el mismo commit.

6. **Nuevo patrón de versión en un documento.** Si se añade un nuevo documento con versión en cabecera o pie:
   - Añadir el check de versión correspondiente en la sección 1 del script.
   - Añadir el patrón `sed` correspondiente en el workflow de release (`bump-version.yml`).

**Regla práctica:** Después de hacer cambios, ejecutar `bash scripts/check-spec-consistency.sh` localmente. Si falla, corregir la inconsistencia o actualizar el script según corresponda. No pushear con checks rotos.

**Al añadir un nuevo check al script**, actualizar también la página de verificación del sitio web (`site/src/pages/verificacion.astro`) para que la lista de comprobaciones publicada sea coherente con lo que realmente se ejecuta.

---

## Sitio web (GitHub Pages)

El proyecto incluye un sitio web estático en el directorio `site/`, construido con **Astro** y desplegado automáticamente en GitHub Pages.

### Estructura del sitio

```
site/
  astro.config.mjs          Config de Astro (base path, markdown, mermaid)
  package.json               Dependencias
  tsconfig.json              Config TypeScript
  src/
    layouts/
      Base.astro             Layout base (head, nav, footer)
      Doc.astro              Layout para paginas de documentacion
    pages/
      index.astro            Landing page
      white-paper.astro      Renderiza README.md
      protocolo.astro        Renderiza PROTOCOL.md
      seguridad.astro        Renderiza SECURITY-ANALYSIS.md
      verificacion.astro     Pagina de verificacion automatizada
      changelog.astro        Renderiza CHANGELOG.md
    styles/
      global.css             Estilos globales
  public/
    favicon.svg              Favicon del proyecto
```

### Páginas

| Ruta | Contenido | Fuente |
|------|-----------|--------|
| `/` | Landing page con hero, principios y CTA | `site/src/pages/index.astro` |
| `/white-paper/` | White paper divulgativo | `README.md` (renderizado) |
| `/protocolo/` | Especificación técnica | `PROTOCOL.md` (renderizado) |
| `/seguridad/` | Análisis de seguridad | `SECURITY-ANALYSIS.md` (renderizado) |
| `/verificacion/` | Verificación automatizada de la especificación | `site/src/pages/verificacion.astro` |
| `/changelog/` | Registro de cambios | `CHANGELOG.md` (renderizado) |

### Cómo funciona

- Las páginas de documentación (`white-paper`, `protocolo`, `changelog`) leen los archivos Markdown de la raíz del repositorio en build time y los renderizan a través del pipeline de Astro con `createMarkdownProcessor`.
- Los diagramas Mermaid se renderizan a SVG en build time vía `rehype-mermaid` (Playwright + Chromium). No se envía JavaScript al cliente.
- El deploy se dispara automáticamente al pushear a main si cambian archivos relevantes (`site/**`, `README.md`, `PROTOCOL.md`, `CHANGELOG.md`, `VERSION`).

### Desarrollo local

```bash
cd site
npm install
npx playwright install --with-deps chromium
npm run dev      # Servidor de desarrollo
npm run build    # Build de produccion en site/dist/
npm run preview  # Preview del build
```

---

## Estructura futura prevista del repositorio

```
README.md                              White paper divulgativo
PROTOCOL.md                            Especificacion tecnica
SECURITY-ANALYSIS.md                   Estudio de vulnerabilidades y analisis de seguridad
CLAUDE.md                              Guia de contribucion (este archivo)
CHANGELOG.md                           Registro de cambios
VERSION                                Version actual (fuente de verdad)
LICENSE                                Licencia del proyecto
site/                                  Sitio web publico (Astro + GitHub Pages)
.github/
  workflows/
    bump-version.yml                   Workflow de bump automatico
    deploy-site.yml                    Workflow de deploy del sitio web
    spec-consistency.yml               Workflow de coherencia de la especificacion
scripts/
  check-spec-consistency.sh            Verificacion de coherencia entre documentos
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

- Mensajes de commit en inglés, concisos, en imperativo: "Add threat model for token replay", "Fix DA terminology in README".
- Un commit por cambio lógico. No mezclar cambios en README.md y PROTOCOL.md en el mismo commit salvo que sean la misma corrección de coherencia.
- Nunca commitear archivos generados, binarios ni secretos.
- **No modificar VERSION manualmente.** El workflow de GitHub Actions se encarga de ello.
- **No estampar [Unreleased] manualmente en CHANGELOG.md.** Solo añadir entradas bajo esa sección. El workflow se encarga de estampar la versión y la fecha.
- Los commits de release (`Release vX.Y.Z`) son creados exclusivamente por el workflow.
- Para hacer un release: `git tag vX.Y.Z && git push origin vX.Y.Z`. El tag dispara el workflow automáticamente.
- **Sin firmas ni trailers en los commits.** No incluir `Co-Authored-By`, `Signed-off-by`, `Generated by` ni ninguna otra fórmula de autoría o co-autoría. Los mensajes de commit deben contener exclusivamente la descripción del cambio.

---

## Resumen de verificación rápida

Antes de dar por bueno cualquier cambio, verificar:

- [ ] Device Agent no se equipara con control parental
- [ ] Los cuatro principios inviolables no se contradicen
- [ ] La terminología usa los nombres canónicos del protocolo
- [ ] El contenido está en el documento correcto según su audiencia
- [ ] Los diagramas Mermaid son coherentes con el texto
- [ ] No se introducen campos en el token sin justificación
- [ ] README.md, PROTOCOL.md y SECURITY-ANALYSIS.md son coherentes entre sí
- [ ] Si se modificó PROTOCOL.md: se revisó el resumen ejecutivo de SECURITY-ANALYSIS.md y se actualizaron las vulnerabilidades afectadas
- [ ] Sin emojis, sin superlativos, sin lenguaje de marketing
- [ ] Ortografía castellana correcta (tildes, eñes, ¿? ¡!)
- [ ] `bash scripts/check-spec-consistency.sh` pasa sin errores
- [ ] Si se modificó el script de verificación: se actualizó `site/src/pages/verificacion.astro`
