start:
	docker compose -f docker-compose.local.yml up --build

stop:
	docker compose -f docker-compose.local.yml down

precommit-coupons:

precommit-orderfoodonline:
	@echo "Run Generate mocks, docs and run build/tests"
	cd backend-challenge/services/orderfoodonline && make precommit

kill:
	sudo lsof -t -i -P -n | xargs sudo kill -9
