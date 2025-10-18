package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lucasepe/snipr/internal/cmd/list"
	"github.com/lucasepe/snipr/internal/cmd/run"
	"github.com/lucasepe/x/cl"
	xtext "github.com/lucasepe/x/text"
)

const (
	appName = "snipr"
)

var (
	Version = "v0.0.0"
)

func main() {
	ctx := context.Background()

	top := flag.NewFlagSet(appName, flag.ExitOnError)

	tool := cl.NewTool(top, appName)
	tool.Output = os.Stdout
	tool.Error = os.Stderr
	tool.Header = func(w io.Writer) {
		fmt.Fprint(w, "      ▗       \n")
		fmt.Fprintf(w, "▞▀▘▛▀▖▄ ▛▀▖▙▀▖ (%s)\n", Version)
		fmt.Fprint(w, "▝▀▖▌ ▌▐ ▙▄▘▌  \n")
		fmt.Fprint(w, "▀▀ ▘ ▘▀▘▌  ▘ \n")
		fmt.Fprint(w, " by Luca Sepe (https://github.com/lucasepe)\n\n")
		fmt.Fprint(w, "Turn your Markdown into a runnable script.\n")
		fmt.Fprintln(w)
	}

	tool.Footer = func(w io.Writer) {
		fmt.Fprintf(w, "Use \"%s <command> -h\" for more information about a command.\n\n", appName)
		fmt.Fprintln(w)

		fmt.Fprint(w, "SUPPORT:\n\n")
		fmt.Fprint(w, xtext.Indent(
			strings.Join([]string{
				"If you find this tool helpful consider supporting with a donation.",
				"Every bit helps cover development time and fuels future improvements.\n",
				"Your support truly makes a difference — thank you!\n",
				"  * https://www.paypal.com/donate/?hosted_button_id=FV575PVWGXZBY\n",
			}, "\n"), "  "))
		fmt.Fprint(w, "\n")
	}

	// Register tasks
	tool.Register(list.Task(appName), "")
	tool.Register(run.Task(appName), "")

	top.Parse(os.Args[1:])

	os.Exit(int(tool.Execute(ctx)))
}
