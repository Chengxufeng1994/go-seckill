SHELL = /bin/bash

.PHONY: proto

proto:
	rm -rf pb/*.go
	protoc --proto_path=pb --go_out=pb --go_opt=paths=source_relative \
			--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
			pb/*.proto