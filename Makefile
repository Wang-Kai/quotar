include VERSION

idl:
	@rm -fr pb/*.go
	@protoc -I=. pb/*.proto --go_out=plugins=grpc:.

bin/quotar:
	@rm -fr quotar_linux
	@GOOS=linux go build -v -o ./bin/quotar_linux ./quotar/main.go

bin/nfs-provisioner:
	rm -fr bin/nfs-provisioner
	GOOS=linux go build -v -o bin/nfs-provisioner ./nfs-provisioner/*.go

publish: bin/linux
	rsync -azP ./bin/quotar_linux root@106.75.223.111:/root/xfs/

image/nfs-provisioner: bin/nfs-provisioner
	docker build -t uhub.service.ucloud.cn/sf_open/nfs-provisioner:${VERSION} .
	docker push uhub.service.ucloud.cn/sf_open/nfs-provisioner:${VERSION}
	