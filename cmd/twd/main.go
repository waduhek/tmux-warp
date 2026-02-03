package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"

	clih "github.com/waduhek/tmux-warp/internal/cli"
	"github.com/waduhek/tmux-warp/internal/parser"
	"github.com/waduhek/tmux-warp/internal/tmux"
	"github.com/waduhek/tmux-warp/internal/warp"
)

func main() {
	tmuxCommands := tmux.NewTmuxCommands()
	tmuxManager := tmux.NewTmuxManager(tmuxCommands)
	parse := parser.NewParser()
	warp := warp.NewWarper(parse)
	cliHandler := clih.NewCLIHandler(tmuxManager, warp)

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "start a new session and warp to it",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:      "name",
						UsageText: "<name>",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					nameArg := c.StringArg("name")
					if nameArg == "" {
						return cli.Exit("name argument cannot be empty", 1)
					}

					if err := cliHandler.StartSession(nameArg); err != nil {
						return cli.Exit(err.Error(), 1)
					}
					return nil
				},
			},
		},
	}

	cmd.Run(context.Background(), os.Args)
}
