.PHONY: lint
lint: 
	golangci-lint run --config=.golangci.yml ./...

.PHONY: build-proxy-server
build-proxy-server: 
	go build -o /bin/proxy-server cmd/proxy/main.go

.PHONY: run-proxy-server
run-proxy-server: 
	go run cmd/proxy/main.go