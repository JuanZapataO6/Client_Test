package main

import (
	"context"
	"fmt"
	"log"

	dgo "github.com/dgraph-io/dgo/v2"
	api "github.com/dgraph-io/dgo/v2/protos/api"
	grpc "google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	op := &api.Operation{}
	op.Schema = `
		transaction_id: string @index(exact) .
		age: int .
		married: bool .
		loc: geo .
		dob: datetime .
	`

	ctx := context.Background()
	err = dg.Alter(ctx, op)
	if err != nil {
		log.Fatal(err)
	}

	// Ask for the type of name and age.
	resp, err := dg.NewTxn().Query(ctx, `schema(pred: [name, age]) {type}`)
	if err != nil {
		log.Fatal(err)
	}

	// resp.Json contains the schema query response.
	fmt.Println(string(resp.Json))
}
