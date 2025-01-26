build_proto:
	protoc --go_out=./teashop_proto --go_opt=paths=source_relative \
		--go-grpc_out=./teashop_proto --go-grpc_opt=paths=source_relative \
		teashop.proto
