#!/usr/bin/env bash
# check-spec-consistency.sh
#
# Verifica la coherencia entre los documentos de especificación de AAVP.
# Pensado para ejecutarse en CI (GitHub Actions) y localmente.
# Compatible con macOS (BSD grep) y Linux (GNU grep).
#
# Exit code 0 = todo coherente, 1 = hay inconsistencias.
# Genera un fichero results.json con el detalle de cada check.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

PASS=0
FAIL=0
CHECKS=()

# ─── Helpers ────────────────────────────────────────────────────────

check() {
  local name="$1"
  local description="$2"
  local result="$3" # "pass" or "fail"
  local detail="${4:-}"

  if [[ "$result" == "pass" ]]; then
    PASS=$((PASS + 1))
    printf "  \033[32m✓\033[0m %s\n" "$description"
  else
    FAIL=$((FAIL + 1))
    printf "  \033[31m✗\033[0m %s\n" "$description"
    if [[ -n "$detail" ]]; then
      printf "    → %s\n" "$detail"
    fi
  fi

  CHECKS+=("{\"name\":\"$name\",\"description\":\"$description\",\"result\":\"$result\",\"detail\":\"${detail//\"/\\\"}\"}")
}

# Portable version extraction: extracts the version from a grep match.
# Usage: extract_version "line containing vX.Y.Z" "prefix_pattern"
extract_version_after() {
  local text="$1"
  local prefix="$2"
  echo "$text" | sed -n "s/.*${prefix}v\{0,1\}\([0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*\).*/\1/p" | head -1
}

# Count matches safely (always returns a single number)
count_matches() {
  local pattern="$1"
  local file="$2"
  local n
  n=$(grep -c "$pattern" "$file" 2>/dev/null) || true
  printf "%d" "${n:-0}"
}

# ─── 1. Versión ─────────────────────────────────────────────────────

printf "\n\033[1m1. Coherencia de versión\033[0m\n"

VERSION=$(tr -d '[:space:]' < "$REPO_ROOT/VERSION")

# README.md header: "White Paper vX.Y.Z"
README_HEADER_LINE=$(grep "White Paper v" "$REPO_ROOT/README.md" 2>/dev/null | head -1 || echo "")
README_HEADER_VERSION=$(extract_version_after "$README_HEADER_LINE" "White Paper v")
README_HEADER_VERSION="${README_HEADER_VERSION:-NOT_FOUND}"
if [[ "$README_HEADER_VERSION" == "$VERSION" ]]; then
  check "version_readme_header" "README.md header coincide con VERSION ($VERSION)" "pass"
else
  check "version_readme_header" "README.md header coincide con VERSION" "fail" "VERSION=$VERSION, README header=$README_HEADER_VERSION"
fi

# README.md footer: "Protocol · vX.Y.Z"
README_FOOTER_LINE=$(grep "Protocol · v" "$REPO_ROOT/README.md" 2>/dev/null | head -1 || echo "")
README_FOOTER_VERSION=$(extract_version_after "$README_FOOTER_LINE" "Protocol · v")
README_FOOTER_VERSION="${README_FOOTER_VERSION:-NOT_FOUND}"
if [[ "$README_FOOTER_VERSION" == "$VERSION" ]]; then
  check "version_readme_footer" "README.md footer coincide con VERSION ($VERSION)" "pass"
else
  check "version_readme_footer" "README.md footer coincide con VERSION" "fail" "VERSION=$VERSION, README footer=$README_FOOTER_VERSION"
fi

# PROTOCOL.md header: "vX.Y.Z — Borrador"
PROTOCOL_HEADER_LINE=$(grep "Borrador" "$REPO_ROOT/PROTOCOL.md" 2>/dev/null | head -1 || echo "")
PROTOCOL_HEADER_VERSION=$(extract_version_after "$PROTOCOL_HEADER_LINE" "v")
PROTOCOL_HEADER_VERSION="${PROTOCOL_HEADER_VERSION:-NOT_FOUND}"
if [[ "$PROTOCOL_HEADER_VERSION" == "$VERSION" ]]; then
  check "version_protocol_header" "PROTOCOL.md header coincide con VERSION ($VERSION)" "pass"
else
  check "version_protocol_header" "PROTOCOL.md header coincide con VERSION" "fail" "VERSION=$VERSION, PROTOCOL header=$PROTOCOL_HEADER_VERSION"
fi

# PROTOCOL.md footer: "Especificación Técnica · vX.Y.Z"
PROTOCOL_FOOTER_LINE=$(grep "cnica · v" "$REPO_ROOT/PROTOCOL.md" 2>/dev/null | head -1 || echo "")
PROTOCOL_FOOTER_VERSION=$(extract_version_after "$PROTOCOL_FOOTER_LINE" "cnica · v")
PROTOCOL_FOOTER_VERSION="${PROTOCOL_FOOTER_VERSION:-NOT_FOUND}"
if [[ "$PROTOCOL_FOOTER_VERSION" == "$VERSION" ]]; then
  check "version_protocol_footer" "PROTOCOL.md footer coincide con VERSION ($VERSION)" "pass"
else
  check "version_protocol_footer" "PROTOCOL.md footer coincide con VERSION" "fail" "VERSION=$VERSION, PROTOCOL footer=$PROTOCOL_FOOTER_VERSION"
fi

# SECURITY-ANALYSIS.md header: "vX.Y.Z — Documento de trabajo"
SECURITY_HEADER_LINE=$(grep "Documento de trabajo" "$REPO_ROOT/SECURITY-ANALYSIS.md" 2>/dev/null | head -1 || echo "")
SECURITY_HEADER_VERSION=$(extract_version_after "$SECURITY_HEADER_LINE" "v")
SECURITY_HEADER_VERSION="${SECURITY_HEADER_VERSION:-NOT_FOUND}"
if [[ "$SECURITY_HEADER_VERSION" == "$VERSION" ]]; then
  check "version_security_header" "SECURITY-ANALYSIS.md header coincide con VERSION ($VERSION)" "pass"
else
  check "version_security_header" "SECURITY-ANALYSIS.md header coincide con VERSION" "fail" "VERSION=$VERSION, SECURITY header=$SECURITY_HEADER_VERSION"
fi

# SECURITY-ANALYSIS.md footer: "Vulnerabilidades · vX.Y.Z"
SECURITY_FOOTER_LINE=$(grep "Vulnerabilidades · v" "$REPO_ROOT/SECURITY-ANALYSIS.md" 2>/dev/null | head -1 || echo "")
SECURITY_FOOTER_VERSION=$(extract_version_after "$SECURITY_FOOTER_LINE" "Vulnerabilidades · v")
SECURITY_FOOTER_VERSION="${SECURITY_FOOTER_VERSION:-NOT_FOUND}"
if [[ "$SECURITY_FOOTER_VERSION" == "$VERSION" ]]; then
  check "version_security_footer" "SECURITY-ANALYSIS.md footer coincide con VERSION ($VERSION)" "pass"
else
  check "version_security_footer" "SECURITY-ANALYSIS.md footer coincide con VERSION" "fail" "VERSION=$VERSION, SECURITY footer=$SECURITY_FOOTER_VERSION"
fi

# ─── 2. Tamaño del token ────────────────────────────────────────────

printf "\n\033[1m2. Tamaño del token\033[0m\n"

TOKEN_SIZE="331"

for doc in README.md PROTOCOL.md SECURITY-ANALYSIS.md; do
  filepath="$REPO_ROOT/$doc"
  [[ -f "$filepath" ]] || continue

  STALE_COUNT=$(count_matches "305 bytes" "$filepath")
  if [[ "$STALE_COUNT" -gt 0 ]]; then
    check "token_size_stale_${doc}" "$doc no contiene referencias obsoletas a '305 bytes'" "fail" "$STALE_COUNT ocurrencias de '305 bytes'"
  else
    check "token_size_stale_${doc}" "$doc no contiene referencias obsoletas a '305 bytes'" "pass"
  fi
done

CORRECT_COUNT=$(count_matches "${TOKEN_SIZE} bytes" "$REPO_ROOT/PROTOCOL.md")
if [[ "$CORRECT_COUNT" -gt 0 ]]; then
  check "token_size_protocol" "PROTOCOL.md menciona '${TOKEN_SIZE} bytes' ($CORRECT_COUNT veces)" "pass"
else
  check "token_size_protocol" "PROTOCOL.md menciona '${TOKEN_SIZE} bytes'" "fail" "No se encontró '${TOKEN_SIZE} bytes'"
fi

# ─── 3. Campos canónicos del token ──────────────────────────────────

printf "\n\033[1m3. Campos canónicos del token\033[0m\n"

CANONICAL_FIELDS=("token_type" "nonce" "token_key_id" "age_bracket" "expires_at" "authenticator")

for field in "${CANONICAL_FIELDS[@]}"; do
  COUNT=$(count_matches "${field}" "$REPO_ROOT/PROTOCOL.md")
  if [[ "$COUNT" -gt 0 ]]; then
    check "field_${field}" "Campo '${field}' presente en PROTOCOL.md ($COUNT refs)" "pass"
  else
    check "field_${field}" "Campo '${field}' presente en PROTOCOL.md" "fail" "No se encontró"
  fi
done

# Check that removed fields don't appear as active in field definition tables
# We look for table rows (starting with |) that contain the field name in backticks
# but exclude rows that mention exclusion/elimination
REMOVED_FIELDS=("issued_at" "implementer_sig")
for field in "${REMOVED_FIELDS[@]}"; do
  # Get table rows with this field, then exclude rows about exclusion
  ACTIVE_ROWS=$(grep "^|" "$REPO_ROOT/PROTOCOL.md" 2>/dev/null \
    | grep "\`${field}\`" \
    | grep -v -i -E "exclu|elimin|remov|no se incluye|eliminado" \
    || true)
  if [[ -z "$ACTIVE_ROWS" ]]; then
    STRUCT_COUNT=0
  else
    STRUCT_COUNT=$(echo "$ACTIVE_ROWS" | wc -l | tr -d ' ')
  fi
  if [[ "$STRUCT_COUNT" -eq 0 ]]; then
    check "removed_field_${field}" "Campo eliminado '${field}' no aparece en la estructura activa del token" "pass"
  else
    check "removed_field_${field}" "Campo eliminado '${field}' no aparece en la estructura activa del token" "fail" "$STRUCT_COUNT ocurrencias en tablas de estructura"
  fi
done

# ─── 4. Franjas de edad ─────────────────────────────────────────────

printf "\n\033[1m4. Franjas de edad\033[0m\n"

AGE_BRACKETS=("UNDER_13" "AGE_13_15" "AGE_16_17" "OVER_18")

for bracket in "${AGE_BRACKETS[@]}"; do
  if grep -q "$bracket" "$REPO_ROOT/PROTOCOL.md"; then
    check "bracket_${bracket}" "Franja '${bracket}' definida en PROTOCOL.md" "pass"
  else
    check "bracket_${bracket}" "Franja '${bracket}' definida en PROTOCOL.md" "fail" "No se encontró"
  fi
done

# ─── 5. Esquema criptográfico ───────────────────────────────────────

printf "\n\033[1m5. Esquema criptográfico\033[0m\n"

if grep -q "RSAPBSSA-SHA384" "$REPO_ROOT/PROTOCOL.md"; then
  check "crypto_scheme" "RSAPBSSA-SHA384 referenciado en PROTOCOL.md" "pass"
else
  check "crypto_scheme" "RSAPBSSA-SHA384 referenciado en PROTOCOL.md" "fail" "No se encontró"
fi

if grep -q "RFC 9474" "$REPO_ROOT/PROTOCOL.md"; then
  check "crypto_rfc" "RFC 9474 referenciado en PROTOCOL.md" "pass"
else
  check "crypto_rfc" "RFC 9474 referenciado en PROTOCOL.md" "fail" "No se encontró"
fi

if [[ -f "$REPO_ROOT/SECURITY-ANALYSIS.md" ]]; then
  if grep -q "RSAPBSSA-SHA384" "$REPO_ROOT/SECURITY-ANALYSIS.md"; then
    check "crypto_scheme_security" "RSAPBSSA-SHA384 referenciado en SECURITY-ANALYSIS.md" "pass"
  else
    check "crypto_scheme_security" "RSAPBSSA-SHA384 referenciado en SECURITY-ANALYSIS.md" "fail" "No se encontró"
  fi
fi

# ─── 6. Terminología: Device Agent ≠ control parental ───────────────

printf "\n\033[1m6. Terminología\033[0m\n"

FORBIDDEN_PATTERNS=(
  "el control parental (Device Agent)"
  "el DA, es decir, el sistema de control parental"
  "el DA es el sistema de control parental"
  "Device Agent (control parental)"
)

TERM_OK=true
for pattern in "${FORBIDDEN_PATTERNS[@]}"; do
  for doc in README.md PROTOCOL.md; do
    MATCH_COUNT=$(count_matches "$pattern" "$REPO_ROOT/$doc")
    if [[ "$MATCH_COUNT" -gt 0 ]]; then
      check "terminology_da_${doc}" "$doc no equipara DA con control parental" "fail" "Patron encontrado: '$pattern'"
      TERM_OK=false
    fi
  done
done

if [[ "$TERM_OK" == true ]]; then
  check "terminology_da" "Ningun documento equipara DA con control parental" "pass"
fi

for doc in PROTOCOL.md SECURITY-ANALYSIS.md; do
  [[ -f "$REPO_ROOT/$doc" ]] || continue
  if grep -q "Device Agent" "$REPO_ROOT/$doc" && grep -q "Verification Gate" "$REPO_ROOT/$doc"; then
    check "terminology_roles_${doc}" "$doc usa nombres completos de los roles del protocolo" "pass"
  else
    check "terminology_roles_${doc}" "$doc usa nombres completos de los roles del protocolo" "fail" "Falta 'Device Agent' o 'Verification Gate'"
  fi
done

# ─── 7. Coherencia cruzada ──────────────────────────────────────────

printf "\n\033[1m7. Coherencia cruzada\033[0m\n"

PROTOCOL_PARTIAL=$(count_matches "parcialmente ciega" "$REPO_ROOT/PROTOCOL.md")
README_PARTIAL=$(count_matches "parcialmente ciega" "$REPO_ROOT/README.md")

if [[ "$PROTOCOL_PARTIAL" -gt 0 && "$README_PARTIAL" -gt 0 ]]; then
  check "cross_blind_sig" "README.md menciona 'parcialmente ciega' (coherente con PROTOCOL.md)" "pass"
elif [[ "$PROTOCOL_PARTIAL" -gt 0 && "$README_PARTIAL" -eq 0 ]]; then
  check "cross_blind_sig" "README.md menciona 'parcialmente ciega' (coherente con PROTOCOL.md)" "fail" "PROTOCOL usa 'parcialmente ciega' pero README no la menciona"
else
  check "cross_blind_sig" "Coherencia de terminologia de firmas ciegas" "pass"
fi

README_BRACKETS=$(grep -cE 'UNDER_13|AGE_13_15|AGE_16_17|OVER_18|Menor de 13|Entre 13 y 15|Entre 16 y 17|Mayor de 18' "$REPO_ROOT/README.md" 2>/dev/null) || true
README_BRACKETS="${README_BRACKETS:-0}"
if [[ "$README_BRACKETS" -ge 4 ]]; then
  check "cross_brackets" "README.md define las 4 franjas de edad" "pass"
else
  check "cross_brackets" "README.md define las 4 franjas de edad" "fail" "Solo se encontraron $README_BRACKETS referencias a franjas"
fi

# ─── 8. Formato binario del token ────────────────────────────────────

printf "\n\033[1m8. Formato binario del token\033[0m\n"

# Parse the binary format table from PROTOCOL.md and verify offsets add up.
# The table looks like:
#   0       2       token_type           Público
#   2       32      nonce                Cegado
#   ...
# We extract offset and size columns and verify:
#   - Each offset == previous offset + previous size
#   - Last offset + last size == TOKEN_SIZE (331)

OFFSET_ERRORS=""
PREV_OFFSET=-1
PREV_SIZE=-1
EXPECTED_OFFSET=0
LINE_COUNT=0

while IFS= read -r line; do
  # Extract leading numbers: offset and size
  offset=$(echo "$line" | awk '{print $1}')
  size=$(echo "$line" | awk '{print $2}')

  # Skip non-numeric lines
  case "$offset" in
    [0-9]*) ;;
    *) continue ;;
  esac
  case "$size" in
    [0-9]*) ;;
    *) continue ;;
  esac

  LINE_COUNT=$((LINE_COUNT + 1))

  if [[ "$offset" -ne "$EXPECTED_OFFSET" ]]; then
    field=$(echo "$line" | awk '{print $3}')
    OFFSET_ERRORS="${OFFSET_ERRORS}Campo ${field}: offset ${offset} esperado ${EXPECTED_OFFSET}. "
  fi

  EXPECTED_OFFSET=$((offset + size))
done < <(sed -n '/^Offset.*Campo/,/^---/p' "$REPO_ROOT/PROTOCOL.md" | tail -n +2 | sed '/^---/d')

if [[ "$LINE_COUNT" -eq 0 ]]; then
  check "binary_format_parse" "Tabla de formato binario encontrada en PROTOCOL.md" "fail" "No se pudo parsear la tabla de offsets"
elif [[ -n "$OFFSET_ERRORS" ]]; then
  check "binary_format_offsets" "Offsets del formato binario son secuenciales" "fail" "$OFFSET_ERRORS"
else
  check "binary_format_offsets" "Offsets del formato binario son secuenciales" "pass"
fi

if [[ "$LINE_COUNT" -gt 0 && "$EXPECTED_OFFSET" -eq "$TOKEN_SIZE" ]]; then
  check "binary_format_total" "Suma de campos = $TOKEN_SIZE bytes" "pass"
elif [[ "$LINE_COUNT" -gt 0 ]]; then
  check "binary_format_total" "Suma de campos = $TOKEN_SIZE bytes" "fail" "La suma es $EXPECTED_OFFSET, esperado $TOKEN_SIZE"
fi

# ─── 9. Semáforo de SECURITY-ANALYSIS.md ─────────────────────────────

printf "\n\033[1m9. Semaforo de SECURITY-ANALYSIS.md\033[0m\n"

if [[ -f "$REPO_ROOT/SECURITY-ANALYSIS.md" ]]; then
  # Count actual semaphore icons in the status table (lines 33-39 approx)
  # Only count rows in the main status table (lines with | **Name** | icon |)
  TABLE_LINES=$(grep -E '^\| \*\*' "$REPO_ROOT/SECURITY-ANALYSIS.md" | grep -E '🔴|🟡|🟢' || true)

  if [[ -n "$TABLE_LINES" ]]; then
    ACTUAL_RED=$(echo "$TABLE_LINES" | grep -c '🔴') || true
    ACTUAL_RED="${ACTUAL_RED:-0}"
    ACTUAL_YELLOW=$(echo "$TABLE_LINES" | grep -c '🟡') || true
    ACTUAL_YELLOW="${ACTUAL_YELLOW:-0}"
    ACTUAL_GREEN=$(echo "$TABLE_LINES" | grep -c '🟢') || true
    ACTUAL_GREEN="${ACTUAL_GREEN:-0}"

    # Extract claimed counts from the summary line
    SUMMARY_LINE=$(grep -E '[0-9]+ .reas? en rojo' "$REPO_ROOT/SECURITY-ANALYSIS.md" | head -1 || echo "")

    if [[ -n "$SUMMARY_LINE" ]]; then
      CLAIMED_RED=$(echo "$SUMMARY_LINE" | sed -n 's/.*\([0-9][0-9]*\) .reas\{0,1\} en rojo.*/\1/p')
      CLAIMED_YELLOW=$(echo "$SUMMARY_LINE" | sed -n 's/.*\([0-9][0-9]*\) en amarillo.*/\1/p')
      CLAIMED_GREEN=$(echo "$SUMMARY_LINE" | sed -n 's/.*\([0-9][0-9]*\) en verde.*/\1/p')

      if [[ "$ACTUAL_RED" -eq "${CLAIMED_RED:-0}" && "$ACTUAL_YELLOW" -eq "${CLAIMED_YELLOW:-0}" && "$ACTUAL_GREEN" -eq "${CLAIMED_GREEN:-0}" ]]; then
        check "semaphore_counts" "Resumen del semaforo coincide con la tabla (${ACTUAL_RED}R/${ACTUAL_YELLOW}A/${ACTUAL_GREEN}V)" "pass"
      else
        check "semaphore_counts" "Resumen del semaforo coincide con la tabla" "fail" \
          "Tabla: ${ACTUAL_RED}R/${ACTUAL_YELLOW}A/${ACTUAL_GREEN}V. Resumen: ${CLAIMED_RED:-?}R/${CLAIMED_YELLOW:-?}A/${CLAIMED_GREEN:-?}V"
      fi
    else
      check "semaphore_counts" "Resumen del semaforo presente en SECURITY-ANALYSIS.md" "fail" "No se encontro la linea de distribucion"
    fi
  else
    check "semaphore_counts" "Tabla de semaforo encontrada en SECURITY-ANALYSIS.md" "fail" "No se encontraron filas con iconos de semaforo"
  fi
fi

# ─── 10. CHANGELOG tiene [Unreleased] ────────────────────────────────

printf "\n\033[1m10. Integridad del CHANGELOG\033[0m\n"

if grep -q '^\## \[Unreleased\]' "$REPO_ROOT/CHANGELOG.md" 2>/dev/null; then
  check "changelog_unreleased" "CHANGELOG.md contiene seccion [Unreleased]" "pass"
else
  check "changelog_unreleased" "CHANGELOG.md contiene seccion [Unreleased]" "fail" "La seccion [Unreleased] es necesaria para el workflow de release"
fi

# Check that the [Unreleased] comparison link points to the current version
UNRELEASED_LINK=$(grep '^\[Unreleased\]:' "$REPO_ROOT/CHANGELOG.md" 2>/dev/null | head -1 || echo "")
if echo "$UNRELEASED_LINK" | grep -q "v${VERSION}\.\.\.HEAD"; then
  check "changelog_link" "Link de [Unreleased] apunta a v${VERSION}...HEAD" "pass"
else
  check "changelog_link" "Link de [Unreleased] apunta a v${VERSION}...HEAD" "fail" "Link actual: $UNRELEASED_LINK"
fi

# ─── 11. Credencial de sesion del VG ──────────────────────────────────

printf "\n\033[1m11. Credencial de sesion del VG\033[0m\n"

if grep -q "Credencial de Sesión del Verification Gate" "$REPO_ROOT/PROTOCOL.md"; then
  check "session_credential_section" "PROTOCOL.md contiene la seccion de credencial de sesion" "pass"
else
  check "session_credential_section" "PROTOCOL.md contiene la seccion de credencial de sesion" "fail" \
    "Seccion 'Credencial de Sesion del Verification Gate' no encontrada"
fi

if grep -q "session_expires_at" "$REPO_ROOT/PROTOCOL.md"; then
  check "session_expires_field" "Campo session_expires_at documentado en PROTOCOL.md" "pass"
else
  check "session_expires_field" "Campo session_expires_at documentado en PROTOCOL.md" "fail" \
    "No se encontro referencia a session_expires_at"
fi

# ─── 12. Endpoints de descubrimiento ──────────────────────────────────

printf "\n\033[1m12. Endpoints de descubrimiento\033[0m\n"

if grep -q '\.well-known/aavp-issuer' "$REPO_ROOT/PROTOCOL.md"; then
  check "discovery_im_endpoint" "PROTOCOL.md especifica .well-known/aavp-issuer" "pass"
else
  check "discovery_im_endpoint" "PROTOCOL.md especifica .well-known/aavp-issuer" "fail" \
    "No se encontro referencia a .well-known/aavp-issuer"
fi

if grep -q '\.well-known/aavp' "$REPO_ROOT/PROTOCOL.md" && grep -q 'vg_endpoint' "$REPO_ROOT/PROTOCOL.md"; then
  check "discovery_vg_endpoint" "PROTOCOL.md especifica .well-known/aavp con vg_endpoint" "pass"
else
  check "discovery_vg_endpoint" "PROTOCOL.md especifica .well-known/aavp con vg_endpoint" "fail" \
    "No se encontro especificacion completa de .well-known/aavp"
fi

if grep -q 'accepted_token_types' "$REPO_ROOT/PROTOCOL.md"; then
  check "discovery_token_types" "Campo accepted_token_types documentado en PROTOCOL.md" "pass"
else
  check "discovery_token_types" "Campo accepted_token_types documentado en PROTOCOL.md" "fail" \
    "No se encontro referencia a accepted_token_types"
fi

# Verificar que no quedan referencias obsoletas al esquema borrador
STALE_ALGO=$(grep -c '"rsa-blind-2048"' "$REPO_ROOT/SECURITY-ANALYSIS.md" 2>/dev/null) || true
STALE_ALGO="${STALE_ALGO:-0}"
if [[ "$STALE_ALGO" -eq 0 ]]; then
  check "discovery_no_stale_algo" "SECURITY-ANALYSIS.md no contiene referencias a esquema obsoleto rsa-blind-2048" "pass"
else
  check "discovery_no_stale_algo" "SECURITY-ANALYSIS.md no contiene referencias a esquema obsoleto rsa-blind-2048" "fail" \
    "$STALE_ALGO ocurrencias de rsa-blind-2048"
fi

# ─── Resumen ─────────────────────────────────────────────────────────

TOTAL=$((PASS + FAIL))
printf "\n\033[1mResultado: %d/%d checks passed\033[0m\n" "$PASS" "$TOTAL"

# Generate JSON results
RESULTS_FILE="$REPO_ROOT/spec-consistency-results.json"
{
  printf '{\n'
  printf '  "total": %d,\n' "$TOTAL"
  printf '  "passed": %d,\n' "$PASS"
  printf '  "failed": %d,\n' "$FAIL"
  printf '  "status": "%s",\n' "$(if [[ "$FAIL" -eq 0 ]]; then echo "consistent"; else echo "inconsistent"; fi)"
  printf '  "checks": [\n'
  for i in "${!CHECKS[@]}"; do
    if [[ $i -lt $((${#CHECKS[@]} - 1)) ]]; then
      printf '    %s,\n' "${CHECKS[$i]}"
    else
      printf '    %s\n' "${CHECKS[$i]}"
    fi
  done
  printf '  ]\n'
  printf '}\n'
} > "$RESULTS_FILE"

if [[ "$FAIL" -gt 0 ]]; then
  printf "\n\033[31m%d inconsistencias encontradas.\033[0m\n\n" "$FAIL"
  exit 1
else
  printf "\n\033[32mEspecificacion coherente.\033[0m\n\n"
  exit 0
fi
