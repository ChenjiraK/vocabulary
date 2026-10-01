GO=go

.PHONY: migrate-up migrate-down

migrate-up:
	$(GO) run . migrate

migrate-down:
	@test -n "$(VERSION)" || (echo "VERSION is required. Usage: make migrate-down VERSION=2" && exit 1)
	$(GO) run . migrate down $(VERSION)
