package main

import (
	"github.com/docker/docker/client"
	"github.com/kevinjqiu/external-dns-docker/controller"
	"github.com/kevinjqiu/external-dns-docker/dns"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	defer cli.Close()

	service := controller.NewController(cli, []dns.Provider{})
	service.Run()
	// messageChan, errChan := cli.Events(context.Background(), types.EventsOptions{})

	// for {
	// 	select {
	// 	case message := <-messageChan:
	// 		fmt.Println(message)

	// 	case err := <-errChan:
	// 		panic(err)
	// 	}
	// }
}
