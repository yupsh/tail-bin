package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	yup "github.com/gloo-foo/framework"
	. "github.com/yupsh/tail"
)

const (
	flagLines   = "lines"
	flagBytes   = "bytes"
	flagFollow  = "follow"
	flagQuiet   = "quiet"
	flagVerbose = "verbose"
)

func main() {
	app := &cli.App{
		Name:  "tail",
		Usage: "output the last part of files",
		UsageText: `tail [OPTIONS] [FILE...]

   Print the last 10 lines of each FILE to standard output.
   With more than one FILE, precede each with a header giving the file name.
   With no FILE, or when FILE is -, read standard input.`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    flagLines,
				Aliases: []string{"n"},
				Usage:   "output the last NUM lines, instead of the last 10",
				Value:   10,
			},
			&cli.IntFlag{
				Name:    flagBytes,
				Aliases: []string{"c"},
				Usage:   "output the last NUM bytes",
			},
			&cli.BoolFlag{
				Name:    flagFollow,
				Aliases: []string{"f"},
				Usage:   "output appended data as the file grows",
			},
			&cli.BoolFlag{
				Name:    flagQuiet,
				Aliases: []string{"q", "silent"},
				Usage:   "never output headers giving file names",
			},
			&cli.BoolFlag{
				Name:    flagVerbose,
				Aliases: []string{"v"},
				Usage:   "always output headers giving file names",
			},
		},
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "tail: %v\n", err)
		os.Exit(1)
	}
}

func action(c *cli.Context) error {
	var params []any

	// Add file arguments (or none for stdin)
	for i := 0; i < c.NArg(); i++ {
		params = append(params, yup.File(c.Args().Get(i)))
	}

	// Add flags based on CLI options
	if c.IsSet(flagLines) {
		params = append(params, LineCount(c.Int(flagLines)))
	}
	if c.IsSet(flagBytes) {
		params = append(params, ByteCount(c.Int(flagBytes)))
	}
	if c.Bool(flagFollow) {
		params = append(params, Follow)
	}
	if c.Bool(flagQuiet) {
		params = append(params, Quiet)
	}
	if c.Bool(flagVerbose) {
		params = append(params, Verbose)
	}

	// Create and execute the tail command
	cmd := Tail(params...)
	return yup.Run(cmd)
}
