package gapi

import (
	"fmt"

	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/pb"
	"tutorial.sqlc.dev/app/token"
	"tutorial.sqlc.dev/app/util"
)

type Server struct {
	pb.UnimplementedSimplebankServer
	config util.Config
	store db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	return server, nil
}
