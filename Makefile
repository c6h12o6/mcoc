

pb: proto/service.proto proto/service.pb.go
	PATH=${PATH}:/home/jbf/go/bin protoc --go-grpc_out=. --go_out=. proto/service.proto --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative

run:
	 SQL_HOST=34.138.249.209 go run mcoc.go mcoc_grpc.go
