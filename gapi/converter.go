package gapi

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/pb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username: 			user.Username,
		FullName: 			user.FullName,
		Email: 				user.Email,
		PasswordChangedAt: 	timestamppb.New(user.PasswordChangedAt),
		CreatedAt: 			timestamppb.New(user.CreatedAt),
	} 
}






// rm -f pb/*.go
// 	rm -f doc/swagger/*.swagger.json
// 	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
// 	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
// 	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
// 	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
// 	proto/*.proto
// 	statik -src=./doc/swagger -dest=./doc -f