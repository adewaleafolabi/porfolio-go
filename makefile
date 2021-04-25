
install:
	@go mod download
	@cd ui && npm i

build-ui:
	@cd ui && npm run build && rm -fr ../cmd/dist && mv dist ../cmd/

build-backend:
	@go fmt ./...
	@go build -o portfolio cmd/main.go

build:
	@make install
	@make build-ui
	@make build-backend

deploy:
	@make build && mv portfolio ~/Portfoilio