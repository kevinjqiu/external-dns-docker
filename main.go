package main

import (
	"github.com/docker/docker/client"
	"github.com/kevinjqiu/external-dns-docker/controller"
	"github.com/kevinjqiu/external-dns-docker/dns/cloudflare"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
				logrus.Fatal("--zone-name must be provided")
			}

			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				logrus.Fatal("cannot create docker client: %v", err)
			}

			defer cli.Close()

			cfProvider, err := cloudflare.NewCloudflareProvider(flags.zoneName, flags.recordSuffix)
			if err != nil {
				logrus.Fatal("cannot create dns provider: %v", err)
			}

			service := controller.NewController(cli, cfProvider)
			if err := service.Run(); err != nil {
				logrus.Fatal("encountered error: %v", err)
			}
		},
	}

	command.PersistentFlags().StringVarP(&flags.zoneName, "zone-name", "z", "", "dns zone name")
	command.PersistentFlags().StringVarP(&flags.recordSuffix, "record-suffix", "s", "", "dns record suffix for records managed by external-dns-docker")
	return command
}

func main() {
	command := newCommand()

	if err := command.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
