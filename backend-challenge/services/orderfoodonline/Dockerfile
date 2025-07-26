FROM --platform=linux/amd64 golang:1.23-alpine AS development

WORKDIR /app/backend-challenge/services/orderfoodonline

# Copy go.mod and go.sum
COPY backend-challenge/services/orderfoodonline/go.mod .
COPY backend-challenge/services/orderfoodonline/go.sum .

# Copy shared library using full relative path
COPY backend-challenge/library /app/backend-challenge/library

# Download dependencies
RUN go mod download

# Install tools
# RUN go install github.com/air-verse/air@v1.52.3
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the rest of the code
COPY backend-challenge/services/orderfoodonline/. .
# COPY backend-challenge/services/orderfoodonline/.air.toml ./

# Download dependencies
RUN go mod download

# Generate Swagger docs
RUN swag init --generalInfo ./cmd/rest/main.go --output ./cmd/rest/docs

# 👇 Build static binary compatible with alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./rest-app ./cmd/rest/main.go

# Create a minimal image to run the microservice
FROM --platform=linux/amd64 alpine:latest

COPY --from=development /app/backend-challenge/services/orderfoodonline/rest-app ./
COPY --from=development /app/backend-challenge/services/orderfoodonline/config.json ./
COPY --from=development /app/backend-challenge/services/orderfoodonline/cmd/rest/docs /cmd/rest/docs
COPY --from=development /app/backend-challenge/services/orderfoodonline/migrations /migrations
RUN chmod +x /rest-app
RUN mkdir -p /logs

EXPOSE 8080

CMD [ "./rest-app" ]