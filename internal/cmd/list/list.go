package list

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/lucasepe/snipr/internal/snippets"
	ioutil "github.com/lucasepe/snipr/internal/util/io"
	"github.com/lucasepe/x/cl"
	"github.com/lucasepe/x/log"
)

func Task(appName string) cl.Task {
	return &listCmd{
		ctx:     context.Background(),
		appName: appName,
	}
}

var (
	_ cl.Task = (*listCmd)(nil)
)

type listCmd struct {
	appName  string
	ctx      context.Context
	filename string
}

func (impl *listCmd) Name() string {
	return "list"
}

func (impl *listCmd) Synopsis() string {
	return "List all executable code blocks in a Markdown document"
}

func (impl *listCmd) Usage() string {
	wri := bytes.Buffer{}
	fmt.Fprintf(&wri, "%s\n\n", impl.Synopsis())

	fmt.Fprint(&wri, "USAGE:\n\n")
	fmt.Fprintf(&wri, "  %s %s [MARKDOWN_FILE]\n\n", impl.appName, impl.Name())

	fmt.Fprint(&wri, "DESCRIPTION:\n\n")
	fmt.Fprintln(&wri, "  Scans a Markdown document and lists all executable code blocks found within.")
	fmt.Fprintln(&wri, "  Each block shows its language, name, and any declared parameters such as")
	fmt.Fprint(&wri, "  dependencies or exports.\n\n")
	fmt.Fprint(&wri, "  The input can be provided either as a file path or via standard input (stdin).\n\n")

	fmt.Fprint(&wri, "EXAMPLES:\n\n")
	fmt.Fprint(&wri, "  List all code blocks in a Markdown file:\n\n")
	fmt.Fprintf(&wri, "    %s list example.md\n\n", impl.appName)

	fmt.Fprintf(&wri, "  Pipe a Markdown file into %s:\n\n", impl.appName)
	fmt.Fprintf(&wri, "    cat example.md | %s list\n\n", impl.appName)

	return wri.String()
}

func (impl *listCmd) Ctx() context.Context {
	return impl.ctx
}

func (impl *listCmd) SetFlags(_ *flag.FlagSet) {}

func (impl *listCmd) Execute(ctx context.Context, fs *flag.FlagSet, args ...any) cl.ExitStatus {
	if err := impl.configure(fs); err != nil {
		log.E("unable to configure task",
			log.String("task", impl.Name()),
			log.Err("err", err),
		)

		return cl.ExitFailure
	}

	sin, close, err := ioutil.FileOrStdin(impl.filename)
	if err != nil {
		log.E("unable to read markdown content",
			log.String("task", impl.Name()),
			log.Err("err", err),
		)

		return cl.ExitFailure
	}
	defer close()

	err = snippets.List(sin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to list all executable markdown codeblocks: %v\n", err)
		return cl.ExitFailure
	}

	return cl.ExitSuccess
}

func (impl *listCmd) configure(fs *flag.FlagSet) (err error) {
	if fs.NArg() > 0 {
		impl.filename = fs.Arg(0)
	}

	return nil
}
