.PHONY: lint
lint: 
	golangci-lint run --config=.golangci.yml ./...

.PHONY: build-proxy-server
build-proxy-server: 
	go build -o /bin/proxy-server cmd/proxy/main.go

.PHONY: run-proxy-server
run-proxy-server: 
	go run cmd/proxy/main.go

.PHONY: up-proxy-server
up-proxy-server: 
	docker compose -f deploy/docker-compose.yaml up -d
	
.PHONY: down-proxy-server
down-proxy-server: 
	docker compose -f deploy/docker-compose.yaml down

.PHONY: gen-ca
gen-ca:
	mkdir certs && ./scripts/gen_ca.sh operelygin
	
