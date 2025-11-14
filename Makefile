.PHONY: ogen
ogen:
	ogen --target generated --clean api/openapi.yml

.PHONY: lint-fix
lint-fix:
	@echo "running golangci-lint"
	golangci-lint run --fix ./...
	@echo "lint complete"
