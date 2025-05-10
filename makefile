run:
	go run cmd/server/main.go

swagger:
	swag init -g ./cmd/server/main.go -o ./docs

migrate:
	rm -f data.sqlite
	touch data.sqlite
	go run cmd/migration/main.go
