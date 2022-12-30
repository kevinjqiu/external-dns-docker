package main

import (
	"github.com/docker/docker/client"
	"github.com/kevinjqiu/external-dns-docker/controller"
	"github.com/kevinjqiu/external-dns-docker/dns/cloudflare"
	"github.com/spf13/cobra"
	"log"
)

type cmdFlags struct {
	zoneName     string
	recordSuffix string
}

var flags cmdFlags

func newCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "external-dns-docker",
		Short: "Manage DNS records for docker containers",
		Run: func(cmd *cobra.Command, args []string) {
			if flags.zoneName == "" {
				log.Fatal("--zone-name must be provided")
			}

			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				panic(err)
			}

			defer cli.Close()

			cfProvider, err := cloudflare.NewCloudflareProvider(flags.zoneName, flags.recordSuffix)
			if err != nil {
				panic(err)
			}

			service := controller.NewController(cli, cfProvider)
			service.Run()
		},
	}

	command.PersistentFlags().StringVarP(&flags.zoneName, "zone-name", "z", "", "dns zone name")
	command.PersistentFlags().StringVarP(&flags.recordSuffix, "record-suffix", "s", "", "dns record suffix for records managed by external-dns-docker")
	return command
}

func main() {
	command := newCommand()

	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
