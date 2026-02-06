# AAVP Formal Verification — Tamarin Prover
#
# Prerequisites:
#   Option A: Tamarin installed locally (https://tamarin-prover.com)
#   Option B: Docker (uses infsec/tamarin-prover image)
#
# Usage:
#   make verify              Run all proofs (security properties + unlinkability)
#   make security            Run security properties only (unforgeability, etc.)
#   make unlinkability       Run unlinkability proof (observational equivalence)
#   make docker-verify       Run all proofs via Docker
#   make docker-security     Run security properties via Docker
#   make docker-unlinkability Run unlinkability via Docker

TAMARIN     ?= tamarin-prover
DOCKER_IMG  := infsec/tamarin-prover
FORMAL_DIR  := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
PROOFS_DIR  := $(FORMAL_DIR)proofs

.PHONY: verify security unlinkability
.PHONY: docker-verify docker-security docker-unlinkability
.PHONY: clean

# --- Local execution ---

verify: security unlinkability

security:
	@echo "=== AAVP: Security properties (unforgeability, executability) ==="
	$(TAMARIN) --prove $(FORMAL_DIR)aavp.spthy

unlinkability:
	@echo "=== AAVP: Unlinkability (observational equivalence) ==="
	$(TAMARIN) --diff --prove $(FORMAL_DIR)aavp-unlinkability.spthy

# --- Docker execution ---

docker-verify: docker-security docker-unlinkability

docker-security:
	@echo "=== AAVP: Security properties via Docker ==="
	docker run --rm -v "$(FORMAL_DIR):/data" $(DOCKER_IMG) \
		--prove /data/aavp.spthy

docker-unlinkability:
	@echo "=== AAVP: Unlinkability via Docker ==="
	docker run --rm -v "$(FORMAL_DIR):/data" $(DOCKER_IMG) \
		--diff --prove /data/aavp-unlinkability.spthy

# --- Cleanup ---

clean:
	rm -rf $(PROOFS_DIR)
