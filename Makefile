.PHONY: generate-proto build deploy

generate-proto:
	protoc \
	--go_out=. \
	--go_opt=Mproto/telegram.proto \
	--go-grpc_out=. \
	--go-grpc_opt=require_unimplemented_servers=false \
	--go-grpc_opt=Mproto/telegram.proto \
	proto/telegram.proto

build:
	mkdir -p build
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o build/telegram cmd/main.go

deploy: build
	ssh kiber "systemctl stop telegram.service"
	scp build/telegram kiber:/var/local/bin/telegram
	ssh kiber "systemctl start telegram.service"
