start:
	docker compose -f docker-compose.local.yml up -d
	docker compose -f docker-compose.local.yml logs -f --tail=10

start-with-build:
	docker compose -f docker-compose.local.yml up --build -d
	docker compose -f docker-compose.local.yml logs -f --tail=10

build-orderfoodonline:
	docker compose -f docker-compose.local.yml build orderfoodonline

build-coupons:
	docker compose -f docker-compose.local.yml build coupons-processor

stop:
	docker compose -f docker-compose.local.yml down

precommit-library:
	@echo "Run Generate mocks, docs and run build/tests"
	cd backend-challenge/library && make precommit

precommit-coupons:
	@echo "Run Generate mocks, docs and run build/tests"
	cd backend-challenge/services/coupons && make precommit

precommit-orderfoodonline:
	@echo "Run Generate mocks, docs and run build/tests"
	cd backend-challenge/services/orderfoodonline && make precommit

kill:
	sudo lsof -t -i -P -n | xargs sudo kill -9
