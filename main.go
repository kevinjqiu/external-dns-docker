package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	messageChan, errChan := cli.Events(context.Background(), types.EventsOptions{})

	for {
		select {
		case message := <-messageChan:
			fmt.Println(message)

		case err := <-errChan:
			panic(err)
		}
	}
}
