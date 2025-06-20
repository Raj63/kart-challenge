start:
	docker compose -f docker-compose.local.yml up --build

stop:
	docker compose -f docker-compose.local.yml down

precommit-coupons:

precommit-orderfoodonline:
	@echo "Run Generate mocks, docs and run build/tests"
	cd ./backend-challenge/services/orderfoodonline && make precommit
	
hot-reload:
	@echo "Setting up local environment..."
	@if [ ! -f ./backend-challenge/services/orderfoodonline/.env ]; then \
		echo "Error: ./backend-challenge/services/orderfoodonline/.env file not found. Please create it first."; \
		exit 1; \
	fi
	
	cd ./backend-challenge/services/orderfoodonline && air -c .air.toml --build.cmd "go build -o ./tmp/orderfoodonline ./cmd/rest/main.go" --build.bin "./tmp/orderfoodonline"

kill:
	sudo lsof -t -i -P -n | xargs sudo kill -9

