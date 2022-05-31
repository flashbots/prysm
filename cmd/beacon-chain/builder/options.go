package buildercmd

import (
	"github.com/prysmaticlabs/prysm/beacon-chain/builder"
	"github.com/prysmaticlabs/prysm/cmd/beacon-chain/flags"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var log = logrus.WithField("prefix", "cmd-powchain")

// FlagOptions for powchain service flag configurations.
func FlagOptions(c *cli.Context) ([]builder.Option, error) {
	endpoints := parseBuilderEndpoints(c)
	opts := []builder.Option{
		builder.WithHttpEndpoints(endpoints),
	}
	return opts, nil
}

func parseBuilderEndpoints(c *cli.Context) []string {
	if c.String(flags.HTTPBuilderFlag.Name) == "" && len(c.StringSlice(flags.FallbackBuilderFlag.Name)) == 0 {
		log.Error("No builder specified to run with the beacon node.")
	}
	endpoints := []string{c.String(flags.HTTPBuilderFlag.Name)}
	endpoints = append(endpoints, c.StringSlice(flags.FallbackBuilderFlag.Name)...)
	return endpoints
}
