#!/bin/bash
set -euo pipefail

# Launch the enclave using the newly-named .eif
nitro-cli run-enclave \
  --eif-path /opt/enclave_main.eif \
  --cpu-count 2 \
  --memory 512 \
  --enclave-cid 16 \
  --debug-mode

# Wait for vsock to come up
sleep 2

# Grab the EnclaveID
ENCLAVE_ID=$(nitro-cli describe-enclaves | jq -r '.[0].EnclaveID')

# Exec your enclave_main binary inside
nitro-cli exec --enclave-id "$ENCLAVE_ID" -- /opt/enclave_main