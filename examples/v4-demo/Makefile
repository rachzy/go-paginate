.PHONY: help run local remote tidy-local tidy-remote

help:
	@echo "Targets:"
	@echo "  make local   - run against ../../v4 (your fork)"
	@echo "  make remote  - run against the published module (go get)"
	@echo "  make run     - alias for local (default for fork testing)"
	@echo "  make tidy-local / tidy-remote - refresh go.sum for each mode"

run: local

local:
	go run -modfile=go.local.mod .

remote:
	go run .

tidy-local:
	go mod tidy -modfile=go.local.mod

tidy-remote:
	go mod tidy
