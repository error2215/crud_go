package main

import (
	"context"
	"fmt"
	//"log"
	"os"

	pb "../getData"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	var id int32

	conn, err := grpc.Dial("localhost:20100", grpc.WithInsecure())
	if err != nil {
		log.Fatal("something went wrong", err)
	}
	defer conn.Close()

	c := pb.NewGetDataClient(conn)

	log.Print("pls write some id: ")

	_, err = fmt.Fscan(os.Stdin, &id)
	if err != nil {
		log.Fatal("something went wrong", err)
	}
	rply, err := c.GetDataById(context.Background(), &pb.RequestId{Id: id})
	if err != nil {
		log.Fatal("something went wrong", err)
	}
	log.Println(rply)

}
