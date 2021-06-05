package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	dgo "github.com/dgraph-io/dgo/v2"
	api "github.com/dgraph-io/dgo/v2/protos/api"
	grpc "google.golang.org/grpc"
)

type buyer struct {
	buyer_id   string     `json:"buyer_id,omitempty"`
	buyer_name string     `json:"buyer_name,omitempty"`
	dob        *time.Time `json:"dob,omitempty"`
	age        int        `json:"age,omitempty"`
}
type product struct {
	product_id   string  `json:"product_id,omitempty"`
	product_name string  `json:"product_name,omitempty"`
	price        float32 `json:"price,omitempty"`
}
type Transaction struct {
	transaction_id string    `json:"transaction_id ,omitempty"`
	buyer          []buyer   `json:"buyer,omitempty"`
	product        []product `json:"product,omitempty"`
	direction_ip   string    `json:"direction_ip,omitempty"`
	device         string    `json:"device,omitempty"`
}

func main() {
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	dob := time.Date(1980, 01, 01, 23, 0, 0, 0, time.UTC)
	T := Transaction{
		transaction_id: "#000060b6ca00",
		buyer: []buyer{
			{
				buyer_id:   "490d6704",
				buyer_name: "beaumont",
				dob:        &dob,
				age:        34,
			},
		},
		product: []product{
			{
				product_id:   "9e160ac0",
				product_name: "Cream  of mushroom condensed soup",
				price:        5020,
			},
			{
				product_id:   "8bb1b853",
				product_name: "Spanish style rice",
				price:        546,
			},
			{
				product_id:   "efef0fea",
				product_name: "Shells + vegan cheddar plant-based mac with chickpea pasta",
				price:        2283,
			},
			{
				product_id:   "d343222d",
				product_name: "Uncured pepperoni",
				price:        2607,
			},
			{
				product_id:   "57296c39",
				product_name: "Pepperoni & mozzarella cheese pizza",
				price:        2949,
			},
		},
		direction_ip: "157.62.23.254",
		device:       "mac",
	}
	txn := dgraphClient.NewTxn()
	ctx := context.Background()
	defer txn.Discard(ctx)
	pb, err := json.Marshal(T)
	if err != nil {
		log.Fatal(err)
	}
	mu := &api.Mutation{
		SetJson: pb,
	}
	res, err := txn.Mutate(ctx, mu)
	if err != nil {
		fmt.Println("Aqui toy")
		log.Fatal(err)
	} else {
		fmt.Println(res)
	}
}
