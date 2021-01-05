idl:
	protoc -I=. pb/*.proto --go_out=plugins=grpc:.