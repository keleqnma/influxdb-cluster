default:
	arun sh -c 'go generate doc.go && APP_ENV=local go run cmd/main.go'

arun:
	arun sh -c 'APP_ENV=local go run cmd/main.go'
dev:
	arun sh -c 'APP_ENV=dev go run cmd/main.go'
swagger: mod-tidied
	go generate doc.go

run: mod-tidied  swagger
	go run cmd/main.go

mod-tidied:
	GO111MODULE=on go mod tidy

mod-vendor:
	GO111MODULE=on go mod vendor

