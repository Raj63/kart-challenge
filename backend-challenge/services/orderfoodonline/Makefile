.PHONY: dep build test generate-mocks generate-docs generate precommit update-version

dep:
	go mod tidy

build:
	docker build -t orderfoodonline .

test:
	go test -failfast -v ./... -coverprofile=coverage/coverage.out && go tool cover -html=coverage/coverage.out -o coverage/coverage.html && go tool cover -func coverage/coverage.out

generate-mocks:
	mockgen -source=internal/http/handlers/interface.go -destination=internal/http/handlers/mocks/interface.go -package mocks	
	mockgen -source=internal/http/middlewares/interface.go -destination=internal/http/middlewares/mocks/interface.go -package mocks	
	mockgen -source=internal/repository/interface.go -destination=internal/repository/mocks/interface.go -package mocks	
	mockgen -source=internal/service/interface.go -destination=internal/service/mocks/interface.go -package mocks	
	
generate-docs:
	swag init --generalInfo cmd/rest/main.go --output cmd/rest/docs

update-version:
	@echo "Checking for route changes to determine version bump..."
	@version_file=internal/constants/version.go; \
	commit_hash=$$(git rev-parse --short HEAD); \
	old_version=$$(grep 'Version =' $$version_file | awk '{print $$3}' | tr -d '"'); \
	IFS='.' read -r major minor patch <<< "$$old_version"; \
	if git diff --name-only HEAD | grep -E 'internal/http/routes/|internal/http/handlers/'; then \
		minor=$$((minor+1)); patch=0; \
		new_version="$$major.$$minor.$$patch"; \
		bump_type="minor"; \
	else \
		patch=$$((patch+1)); \
		new_version="$$major.$$minor.$$patch"; \
		bump_type="patch"; \
	fi; \
	if sed --version >/dev/null 2>&1; then \
		sed_cmd="sed -i -E"; \
	else \
		sed_cmd="sed -i '' -E"; \
	fi; \
	$$sed_cmd "s/Version = \\\"[0-9]+\\.[0-9]+\\.[0-9]+\\\"/Version = \\\"$$new_version\\\"/" $$version_file; \
	$$sed_cmd "s/CommitHash = \\\"[a-f0-9]+\\\"/CommitHash = \\\"$$commit_hash\\\"/" $$version_file; \
	echo "Bumped $$bump_type version to $$new_version with commit hash $$commit_hash"

precommit: update-version generate-docs generate-mocks dep test