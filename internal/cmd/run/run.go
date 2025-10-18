package run

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/lucasepe/snipr/internal/snippets"
	errorsutil "github.com/lucasepe/snipr/internal/util/errors"
	ioutil "github.com/lucasepe/snipr/internal/util/io"
	"github.com/lucasepe/x/cl"
	"github.com/lucasepe/x/log"
)

func Task(appName string) cl.Task {
	return &runCmd{
		ctx:     context.Background(),
		appName: appName,
	}
}

var (
	_ cl.Task = (*runCmd)(nil)
)

type runCmd struct {
	appName    string
	ctx        context.Context
	filename   string
	blockNames []string
}

func (impl *runCmd) Name() string {
	return "run"
}

func (impl *runCmd) Synopsis() string {
	return "Execute all or selected runnable code blocks in a Markdown document"
}

func (impl *runCmd) Usage() string {
	wri := bytes.Buffer{}
	fmt.Fprintf(&wri, "%s\n\n", impl.Synopsis())

	fmt.Fprint(&wri, "USAGE:\n\n")
	fmt.Fprintf(&wri, "  %s %s [MARKDOWN_FILE] [BLOCK_NAMES...]\n\n", impl.appName, impl.Name())

	fmt.Fprint(&wri, "DESCRIPTION:\n\n")
	fmt.Fprintln(&wri, "  Parses a Markdown document and executes runnable code blocks in dependency order.")
	fmt.Fprint(&wri, "  If no file is provided, Markdown is read from stdin.\n\n")
	fmt.Fprint(&wri, "  You can optionally specify one or more block names to run.\n\n")

	fmt.Fprint(&wri, "  Each code block may define:\n\n")
	fmt.Fprint(&wri, "    • name        — a unique identifier for referencing dependencies\n")
	fmt.Fprint(&wri, "    • depends     — comma-separated list of required blocks\n")
	fmt.Fprint(&wri, "    • export      — comma-separated list of environment variables to export\n")
	fmt.Fprint(&wri, "    • file        — path to a file where the block's output will be written\n")
	fmt.Fprint(&wri, "    • dir         — working directory for the command\n")
	fmt.Fprint(&wri, "    • skip        — if \"true\", the block is skipped (useful for temporary disable)\n")
	fmt.Fprint(&wri, "    • timeout     — execution timeout in seconds\n\n")

	fmt.Fprint(&wri, "EXAMPLES:\n\n")
	fmt.Fprint(&wri, "  Run all code blocks in a file:\n\n")
	fmt.Fprintf(&wri, "    %s %s example.md\n\n", impl.Name(), impl.appName)

	fmt.Fprintf(&wri, "  Run only specific blocks (and their dependencies) from a file:\n\n")
	fmt.Fprintf(&wri, "    %s %s example.md build test\n\n", impl.Name(), impl.appName)

	fmt.Fprint(&wri, "  Run all code blocks from stdin:\n\n")
	fmt.Fprintf(&wri, "    cat example.md | %s %s\n\n", impl.Name(), impl.appName)

	fmt.Fprint(&wri, "  Run only specific blocks from stdin:\n\n")
	fmt.Fprintf(&wri, "    cat example.md build test | %s %s\n\n", impl.Name(), impl.appName)

	fmt.Fprint(&wri, "NOTES:\n\n")
	fmt.Fprintf(&wri, "  • Blocks are executed respecting dependency order.\n")
	fmt.Fprintf(&wri, "  • Skipped blocks (unsatisfied dependencies) are reported at the end.\n")
	fmt.Fprintf(&wri, "  • Environment variables from the host system are inherited and can be\n")
	fmt.Fprintf(&wri, "    overridden by variables exported from previous blocks.\n")
	fmt.Fprintf(&wri, "  • Supported languages: shell, sh, bash, go\n\n")
	//fmt.Fprintf(&wri, "  • Execution stops on first critical error, unless --continue is specified.\n\n")

	return wri.String()
}

func (impl *runCmd) Ctx() context.Context {
	return impl.ctx
}

func (impl *runCmd) SetFlags(_ *flag.FlagSet) {}

func (impl *runCmd) Execute(ctx context.Context, fs *flag.FlagSet, args ...any) cl.ExitStatus {
	if err := impl.configure(fs); err != nil {
		log.E("unable to configure task",
			log.String("task", impl.Name()),
			log.Err("err", err),
		)

		return cl.ExitFailure
	}

	sin, close, err := ioutil.FileOrStdin(impl.filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read markdown content: %v", err)
		return cl.ExitFailure
	}
	defer close()

	err = snippets.Run(sin, impl.blockNames...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to run markdown codeblocks:")
		errorsutil.FormatJoinedErrors(os.Stderr, "  • ", err)
		return cl.ExitFailure
	}

	//fmt.Fprintln(os.Stdout, renderer.String())

	return cl.ExitSuccess
}

func (impl *runCmd) configure(fs *flag.FlagSet) error {
	args := fs.Args()

	if len(args) == 0 {
		// nessun argomento: leggi da stdin
		impl.filename = ""
		return nil
	}

	// se il primo argomento è un file esistente, usalo
	if _, err := os.Stat(args[0]); err == nil {
		impl.filename = args[0]
		if len(args) > 1 {
			impl.blockNames = args[1:]
		}
		return nil
	}

	// altrimenti, leggi da stdin e tratta tutti gli argomenti come nomi di blocchi
	impl.filename = ""
	impl.blockNames = args

	return nil
}
