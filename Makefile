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

# Run Postman collection tests
test-api:
	@echo "🚀 Running Postman Collection Tests..."
	@echo "📋 Checking if Newman is installed..."
	@which newman > /dev/null || (echo "❌ Newman not found. Installing..." && npm install -g newman)
	@echo "🔧 Setting up environment variables..."
	@echo "🌐 Starting services if not running..."
	@docker compose -f docker-compose.local.yml up -d
	@echo "⏳ Waiting for services to be ready..."
	@sleep 10
	@echo "🧪 Running API tests with Postman collection..."
	newman run "Order Food Online.postman_collection.json" \
		--env-var "host=localhost" \
		--env-var "port=8080" \
		--env-var "api_key=your-api-key-here" \
		--reporters cli,json \
		--reporter-json-export postman-results.json
	@echo "📊 Test results saved to postman-results.json"
	@echo "✅ API tests completed!"

# Run Postman collection tests with environment file
test-api-env:
	@echo "🚀 Running Postman Collection Tests with environment file..."
	@echo "📋 Checking if Newman is installed..."
	@which newman > /dev/null || (echo "❌ Newman not found. Installing..." && npm install -g newman --unsafe-perm=true || echo "⚠️  Try: sudo npm install -g newman or make install-newman")
	@echo "🔧 Setting up environment variables..."
	@echo "🌐 Starting services if not running..."
	@docker compose -f docker-compose.local.yml up -d
	@echo "⏳ Waiting for services to be ready..."
	@sleep 10
	@echo "🧪 Running API tests with Postman collection and environment..."
	newman run "Order Food Online.postman_collection.json" \
		--environment "Order Food Online.postman_environment.json" \
		--reporters cli,json \
		--reporter-json-export postman-results.json
	@echo "📊 Test results saved to postman-results.json"
	@echo "✅ API tests completed!"

# Run Postman collection tests with custom API key
test-api-with-key:
	@echo "🚀 Running Postman Collection Tests with custom API key..."
	@echo "📋 Checking if Newman is installed..."
	@which newman > /dev/null || (echo "❌ Newman not found. Installing..." && npm install -g newman --unsafe-perm=true || echo "⚠️  Try: sudo npm install -g newman or make install-newman")
	@echo "🔧 Setting up environment variables..."
	@echo "🌐 Starting services if not running..."
	@docker compose -f docker-compose.local.yml up -d
	@echo "⏳ Waiting for services to be ready..."
	@sleep 10
	@echo "🧪 Running API tests with Postman collection..."
	@read -p "Enter your API key: " api_key; \
	newman run "Order Food Online.postman_collection.json" \
		--env-var "host=localhost" \
		--env-var "port=8080" \
		--env-var "api_key=$$api_key" \
		--reporters cli,json \
		--reporter-json-export postman-results.json
	@echo "📊 Test results saved to postman-results.json"
	@echo "✅ API tests completed!"

# Run Postman collection tests with npx (no global install required)
test-api-npx:
	@echo "🚀 Running Postman Collection Tests with npx..."
	@echo "🔧 Setting up environment variables..."
	@echo "🌐 Starting services if not running..."
	@docker compose -f docker-compose.local.yml up -d
	@echo "⏳ Waiting for services to be ready..."
	@sleep 10
	@echo "🧪 Running API tests with Postman collection using npx..."
	npx newman run "Order Food Online.postman_collection.json" \
		--environment "Order Food Online.postman_environment.json" \
		--reporters cli,json \
		--reporter-json-export postman-results.json
	@echo "📊 Test results saved to postman-results.json"
	@echo "✅ API tests completed!"

# Install Newman CLI tool
install-newman:
	@echo "📦 Installing Newman CLI tool..."
	@echo "🔧 Trying different installation methods..."
	@npm install -g newman --unsafe-perm=true || \
	(echo "⚠️  Permission denied. Trying with sudo..." && sudo npm install -g newman) || \
	(echo "⚠️  Alternative: Install via Homebrew or use npx" && echo "💡 Try: brew install newman or npx newman run ...")
	@echo "✅ Newman installation completed!"

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
