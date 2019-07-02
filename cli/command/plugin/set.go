package plugin

import (
	"golang.org/x/net/context"

	"github.com/alcideio/moby/cli"
	"github.com/alcideio/moby/cli/command"
	"github.com/spf13/cobra"
)

func newSetCommand(dockerCli *command.DockerCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set PLUGIN KEY=VALUE [KEY=VALUE...]",
		Short: "Change settings for a plugin",
		Args:  cli.RequiresMinArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dockerCli.Client().PluginSet(context.Background(), args[0], args[1:])
		},
	}

	return cmd
}
