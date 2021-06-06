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

/* type School struct {
	Name string `json:"name,omitempty"`
}

type loc struct {
	Type   string    `json:"type,omitempty"`
	Coords []float64 `json:"coordinates,omitempty"`
} */

// If omitempty is not set, then edges with empty values (0 for int/float, "" for string, false
// for bool) would be created for values not specified explicitly.

/* type Person struct {
	Uid      string     `json:"uid,omitempty"`
	Name     string     `json:"name,omitempty"`
	Age      int        `json:"age,omitempty"`
	Dob      *time.Time `json:"dob,omitempty"`
	Married  bool       `json:"married,omitempty"`
	Raw      []byte     `json:"raw_bytes,omitempty"`
	Friends  []Person   `json:"friend,omitempty"`
	Location loc        `json:"loc,omitempty"`
	School   []School   `json:"school,omitempty"`
} */
type Buyer struct {
	Id   string     `json:"id,omitempty"`
	Name string     `json:"name,omitempty"`
	Dob  *time.Time `json:"dob,omitempty"`
	Age  int        `json:"age,omitempty"`
}
type Product struct {
	Id    string  `json:"id,omitempty"`
	Name  string  `json:"name,omitempty"`
	Price float32 `json:"price,omitempty"`
}
type Transaction struct {
	Uid         string    `json:"uid,omitempty"`
	Id          string    `json:"id,omitempty"`
	BuyerId     []Buyer   `json:"buyerId,omitempty"`
	ProductId   []Product `json:"productId,omitempty"`
	DirectionIp string    `json:"direction_ip,omitempty"`
	Device      string    `json:"device,omitempty"`
}

func main() {
	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}
	defer conn.Close()

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)

	dob := time.Date(1980, 01, 01, 23, 0, 0, 0, time.UTC)
	// While setting an object if a struct has a Uid then its properties in the graph are updated
	// else a new node is created.
	// In the example below new nodes for Alice, Bob and Charlie and school are created (since they
	// dont have a Uid).
	T := Transaction{
		Uid: "_:#000060b6ca00",
		Id:  "#000060b6ca00",
		BuyerId: []Buyer{
			{
				Id:   "490d6704",
				Name: "beaumont",
				Dob:  &dob,
				Age:  34,
			},
		},
		ProductId: []Product{
			{

				Id:    "9e160ac0",
				Name:  "Cream  of mushroom condensed soup",
				Price: 5020,
			},
			{
				Id:    "8bb1b853",
				Name:  "Spanish style rice",
				Price: 546,
			},
			{
				Id:    "efef0fea",
				Name:  "Shells + vegan cheddar plant-based mac with chickpea pasta",
				Price: 2283,
			},
			{
				Id:    "d343222d",
				Name:  "Uncured pepperoni",
				Price: 2607,
			},
			{
				Id:    "57296c39",
				Name:  "Pepperoni & mozzarella cheese pizza",
				Price: 2949,
			},
		},
		DirectionIp: "157.62.23.254",
		Device:      "Mac",
	}
	/* p := Person{
		Uid:     "_:alice",
		Name:    "Alice",
		Age:     26,
		Married: true,
		Location: loc{
			Type:   "Point",
			Coords: []float64{1.1, 2},
		},
		Dob: &dob,
		Raw: []byte("raw_bytes"),
		Friends: []Person{{
			Name: "Bob",
			Age:  24,
		}, {
			Name: "Charlie",
			Age:  29,
		}},
		School: []School{{
			Name: "Crown Public School",
		}},
	}
	*/
	op := &api.Operation{}
	op.Schema = `
		id:string .
		buyerId: [uid] .
		productId: [uid] .
		Id: string .
		name: string .
		price: int .
		age: int . 
		type product{
			id: string
			name: string
			price: int
		}
		type buyer {
			id: string  
			name: string
			age: int 
		}
		directionIp: string .
		device: string .
	`

	ctx := context.Background()
	err = dg.Alter(ctx, op)
	if err != nil {
		log.Fatal(err)
	}

	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(T)
	if err != nil {
		log.Fatal(err)
	}

	mu.SetJson = pb
	assigned, err := dg.NewTxn().Mutate(ctx, mu)
	if err != nil {
		log.Fatal(err)
	}

	// Assigned uids for nodes which were created would be returned in the assigned.Uids map.
	variables := map[string]string{"$id1": assigned.Uids["#000060b6ca00"]}
	q := `query Me($id1: string){
		me(func:  uid($id1)) {
		direction_ip
		device
		buyerId  @filter(eq(name, "beaumont")){
			name
			id
			age
		}
		productId{
			name
			id
			price
		}
		}
	}`
	resp, err := dg.NewTxn().QueryWithVars(ctx, q, variables)
	if err != nil {
		log.Fatal(err)
	}

	type Root struct {
		Me []Transaction `json:"me"`
	}

	var r Root
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Me: %+v\n", r.Me)
	fmt.Printf("Product:", r.productId)
	// R.Me would be same as the person that we set above
}
