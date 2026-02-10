# AAVP Formal Verification — Tamarin Prover
#
# Prerequisites:
#   Option A: Tamarin installed locally (https://tamarin-prover.com)
#   Option B: Docker (uses infsec/tamarin-prover image)
#
# Usage:
#   make verify              Run all proofs (security + unlinkability + SAF)
#   make security            Run security properties only (unforgeability, etc.)
#   make unlinkability       Run unlinkability proof (observational equivalence)
#   make saf                 Run SAF properties only
#   make docker-verify       Run all proofs via Docker
#   make docker-security     Run security properties via Docker
#   make docker-unlinkability Run unlinkability via Docker
#   make docker-saf          Run SAF properties via Docker

TAMARIN     ?= tamarin-prover
DOCKER_IMG  := infsec/tamarin-prover
FORMAL_DIR  := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
PROOFS_DIR  := $(FORMAL_DIR)proofs

.PHONY: verify security unlinkability saf
.PHONY: docker-verify docker-security docker-unlinkability docker-saf
.PHONY: clean

# --- Local execution ---

verify: security unlinkability saf

security:
	@echo "=== AAVP: Security properties (unforgeability, executability) ==="
	$(TAMARIN) --prove $(FORMAL_DIR)aavp.spthy

unlinkability:
	@echo "=== AAVP: Unlinkability (observational equivalence) ==="
	$(TAMARIN) --diff --prove $(FORMAL_DIR)aavp-unlinkability.spthy

saf:
	@echo "=== AAVP: SAF properties (SPD/SPT unforgeability, transparency) ==="
	$(TAMARIN) --prove $(FORMAL_DIR)aavp-saf.spthy

# --- Docker execution ---

docker-verify: docker-security docker-unlinkability docker-saf

docker-security:
	@echo "=== AAVP: Security properties via Docker ==="
	docker run --rm -v "$(FORMAL_DIR):/data" $(DOCKER_IMG) \
		--prove /data/aavp.spthy

docker-unlinkability:
	@echo "=== AAVP: Unlinkability via Docker ==="
	docker run --rm -v "$(FORMAL_DIR):/data" $(DOCKER_IMG) \
		--diff --prove /data/aavp-unlinkability.spthy

docker-saf:
	@echo "=== AAVP: SAF properties via Docker ==="
	docker run --rm -v "$(FORMAL_DIR):/data" $(DOCKER_IMG) \
		--prove /data/aavp-saf.spthy

# --- Cleanup ---

clean:
	rm -rf $(PROOFS_DIR)
