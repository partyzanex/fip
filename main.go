package main

import (
	"github.com/partyzanex/fip/env"
	"github.com/partyzanex/fip/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	prefix = pflag.String("prefix", "fip", "environment prefix")
)

func main() {
	pflag.Parse()

	config, err := env.Read(*prefix)
	if err != nil {
		logrus.Fatal(err)
	}

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.SetLevel(level)

	logrus.Fatal(server.Run(config))
}
