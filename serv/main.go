package main

import (
	pb "../getData"
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	//"log"
	log "github.com/sirupsen/logrus"
	"net"
)

type server struct {
}

func (s *server) GetDataById(ctx context.Context, in *pb.RequestId) (*pb.Data, error) {
	log.Print("Requested id:", in.Id)
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
	)
	if err != nil {
		log.Fatal("something went wrong", err)
	}
	data := searchDataInES(client, in.Id)
	return &data, nil
}

func searchDataInES(client *elastic.Client, id int32) pb.Data {

	q := elastic.NewMultiMatchQuery(id, "_id").Type("phrase")
	res, err := client.Search().
		Index("test_index").
		Pretty(true).
		Query(q).
		Do(context.Background())
	if err != nil {
		log.Fatal("something went wrong", err)
	}

	var data pb.Data
	data.Id = id
	for _, hit := range res.Hits.Hits {
		err := json.Unmarshal(hit.Source, &data)
		if err != nil {
			log.Fatal("something went wrong", err)
		}
		break
	}
	return data
}
func main() {

	listener, err := net.Listen("tcp", ":20100")
	if err != nil {
		log.Fatal("failed to listen", err)
	}
	log.Printf("start listening for id's at port %s", ":20100")

	rpcserv := grpc.NewServer()

	pb.RegisterGetDataServer(rpcserv, &server{})
	reflection.Register(rpcserv)

	err = rpcserv.Serve(listener)
	if err != nil {
		log.Fatal("failed to serve", err)
	}
}
