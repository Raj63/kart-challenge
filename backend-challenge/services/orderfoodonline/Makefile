.PHONY: dep build test generate-mocks generate-docs  generate precommit

dep:
	go mod tidy

build:
	docker build -t orderfoodonline .

test:
	go test -count=1 -failfast -v -race ./... -coverprofile=coverage/coverage.out && go tool cover -html=coverage/coverage.out -o coverage/coverage.html && go tool cover -func coverage/coverage.out

generate-mocks:
	mockgen -source=internal/http/handlers/interface.go -destination=internal/http/handlers/mocks/interface.go -package mocks	
	mockgen -source=internal/http/middlewares/interface.go -destination=internal/http/middlewares/mocks/interface.go -package mocks	
	
#	mockgen -destination=repository/mocks/mocks.go -package=mocks go-backend/repository LevelsRepo,ProductRepo,OrganizationRepo,UserRepo,RolesRepo,PackagingRepo
#	mockgen -source=repository/interfaces.go -destination=repository/mocksrepo/mocks.go -package mocks

generate-docs:
	swag init --generalInfo cmd/rest/main.go --output cmd/rest/docs

precommit: generate-docs generate-mocks dep test