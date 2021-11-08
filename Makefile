run:
	go run ./src

build:
	go build -o backend-challenge ./src

test:
	go test -v ./src

protoc:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./src/pb/discount.proto

image:
	docker build -t backend-challenge_ecommerce .

compose:
	docker-compose up