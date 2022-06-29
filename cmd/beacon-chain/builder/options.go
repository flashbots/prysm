package buildercmd

import (
	"github.com/prysmaticlabs/prysm/beacon-chain/builder"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var log = logrus.WithField("prefix", "cmd-powchain")

// FlagOptions for powchain service flag configurations.
func FlagOptions(c *cli.Context) ([]builder.Option, error) {
	opts := []builder.Option{}
	return opts, nil
}
