generate-proto:
	protoc \
	--go_out=. \
	--go_opt=Mproto/telegram.proto \
	--go-grpc_out=. \
	--go-grpc_opt=require_unimplemented_servers=false \
	--go-grpc_opt=Mproto/telegram.proto \
	proto/telegram.proto
