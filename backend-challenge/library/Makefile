.PHONY: dep test generate-mocks precommit

dep:
	go mod tidy

test:
	go test -failfast -v ./... -coverprofile=coverage/coverage.out && go tool cover -html=coverage/coverage.out -o coverage/coverage.html && go tool cover -func coverage/coverage.out

generate-mocks:
	mockgen -source=internal/repository/interface.go -destination=internal/repository/mocks/interface.go -package mocks	

precommit: update-version generate-mocks dep test