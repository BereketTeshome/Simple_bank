package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"tutorial.sqlc.dev/app/api"
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/gapi"
	"tutorial.sqlc.dev/app/pb"
	"tutorial.sqlc.dev/app/util"
)


func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to db:", err)
	}

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store){
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Can not create a server: ", err)
	}
	
	grpcServer := grpc.NewServer()
	pb.RegisterSimplebankServer(grpcServer, server)
	reflection.Register(grpcServer)
	

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("can not create listener")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("can not start gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store){
	c, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Can not create a server: ", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})


	grpcMux := runtime.NewServeMux(jsonOption)
	err = pb.RegisterSimplebankHandlerServer(c, grpcMux, server)
	if err != nil {
		log.Fatal("can not register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("can not create listener")
	}


	log.Printf("start Http server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("can not start http server")
	}
}

func runGinServer(config util.Config, store db.Store){
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Can not create a server: ", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("Server couldn't start: ", err)
	}
}
