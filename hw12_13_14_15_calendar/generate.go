package proto

//go:generate 	protoc --go_out=internal/server/pb --go-grpc_out=internal/server/pb --grpc-gateway_out=internal/server/pb --grpc-gateway_opt generate_unbound_methods=true --openapiv2_out . api/EventService.proto
