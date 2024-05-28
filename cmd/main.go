package main

import (
	"log"
	"net"
	"os"

	"github.com/akshaybt001/DatingApp_MatchMaking_Service/db"
	"github.com/akshaybt001/DatingApp_MatchMaking_Service/initializer"
	"github.com/akshaybt001/DatingApp_proto_files/pb"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err:=godotenv.Load("../.env");err!=nil{
		log.Fatalf(err.Error())
	}
	addr:=os.Getenv("DB_KEY")
	DB,err:=db.InitDB(addr)
	if err!=nil{
		log.Fatal(err.Error())
	}
	services:=initializer.Initializer(DB)
	server:=grpc.NewServer()
	pb.RegisterMatchServiceServer(server,services)
	listener,err:=net.Listen("tcp",":8084")
	if err!=nil{
		log.Fatalf("failed to listen on port 8084 %v",err)
	}
	log.Printf("MatchMaking sevice listening on port 8084")
	if err=server.Serve(listener);err!=nil{
		log.Fatalf("failed to listen on port 8084 %v",err)
	}
} 