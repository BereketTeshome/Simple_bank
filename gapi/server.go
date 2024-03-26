package gapi

import (
	"fmt"

	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/pb"
	"tutorial.sqlc.dev/app/token"
	"tutorial.sqlc.dev/app/util"
	"tutorial.sqlc.dev/app/worker"
)

type Server struct {
	pb.UnimplementedSimplebankServer
	config util.Config
	store db.Store
	tokenMaker token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}
	return server, nil
}
