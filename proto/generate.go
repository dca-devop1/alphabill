package proto

//go:generate protoc -I. --go_out=paths=source_relative:../internal/rpc/alphabill/. --go-grpc_out=paths=source_relative:../internal/rpc/alphabill/ alphabill.proto
//go:generate protoc -I. --go_out=paths=source_relative:../internal/rpc/transaction/. --go-grpc_out=paths=source_relative:../internal/rpc/transaction/ alphabill.transaction.proto
