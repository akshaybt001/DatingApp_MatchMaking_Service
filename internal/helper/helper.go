package helper

import "google.golang.org/grpc"

func DialGrpc(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	return conn, err
}