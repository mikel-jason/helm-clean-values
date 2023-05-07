#!/usr/bin/env bash

# wanted array expansion -> read as different arguments & flags
# shellcheck disable=SC2068
go run ./cmd/helm-clean-values $@
