package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// This is an excersise on how to implement etcd and use it on the project
func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
		// handle error!
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	resp, err := cli.Put(ctx, "sample_key", "sample_value")
	if err != nil {
		fmt.Printf("some err")
		panic(err)
	}
	defer cancel()
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}

	fmt.Printf("response: %+v\n", resp)

	defer cli.Close()
}
