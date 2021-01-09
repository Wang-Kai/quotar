idl:
	@rm -fr pb/*.go
	@protoc -I=. pb/*.proto --go_out=plugins=grpc:.

bin/linux:
	@rm -fr quotar_linux
	@GOOS=linux go build -v -o quotar_linux ./cmd/main.go

run:
	@go run ./cmd/main.go

publish: bin/linux
	rsync -azP quotar_linux root@106.75.223.111:/root/