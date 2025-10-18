package snippets

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lucasepe/x/text"
)

func Run(in io.Reader, names ...string) error {
	ws, err := newGoWorkspace()
	if err != nil {
		return err
	}
	defer ws.cleanup()

	executors := map[string]executorFunc{
		"shell": runShell,
		"bash":  runShell,
		"sh":    runShell,
		"go":    runGo(ws),
	}

	cbr, err := newCodeBlockRunner(executors)
	if err != nil {
		return err
	}

	return cbr.parseAndRun(in, names...)
}

func newCodeBlockRunner(executors map[string]executorFunc) (cbr *codeBlockRunner, err error) {
	cbr = &codeBlockRunner{
		env:       make(map[string]string),
		executors: executors,
	}

	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			cbr.env[parts[0]] = parts[1]
		}
	}

	return cbr, err
}

type codeBlockRunner struct {
	env       map[string]string
	executors map[string]executorFunc
}

func (cbr *codeBlockRunner) parseAndRun(in io.Reader, names ...string) error {
	blocks, err := extractCodeBlocks(in)
	if err != nil {
		return err
	}

	sortedBlocks, err := topologicalSort(blocks)
	if err != nil {
		return err
	}

	runAll := len(names) == 0

	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}

	blockMap := make(map[string]*codeBlockInfo)
	for _, b := range blocks {
		blockMap[b.name] = b
	}

	if !runAll {
		required := make(map[string]bool)

		var addDeps func(string)
		addDeps = func(name string) {
			if required[name] {
				return
			}
			required[name] = true

			b := blockMap[name]
			if b == nil {
				return
			}

			for _, dep := range strings.FieldsFunc(b.params["depends"],
				func(r rune) bool {
					return r == ',' || r == ' '
				}) {
				dep = strings.TrimSpace(dep)
				if dep != "" {
					addDeps(dep)
				}
			}
		}

		for n := range nameSet {
			addDeps(n)
		}

		nameSet = required
	}

	results := make(map[string]*codeBlockInfo)
	errs := []error{}

	for _, b := range sortedBlocks {
		if !runAll && !nameSet[b.name] {
			continue
		}

		if shouldSkip(b) {
			results[b.name] = b
			continue
		}

		if skip, depName := dependsOnSkipped(b, results); skip {
			results[b.name] = b
			errs = append(errs,
				fmt.Errorf("block %q skipped because dependency %q is marked as skipped",
					b.name, depName))
			continue
		}

		err = cbr.runCodeBlock(b)
		results[b.name] = b

		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	return errors.Join(errs...)
}

func (cbr *codeBlockRunner) runCodeBlock(in *codeBlockInfo) (err error) {
	execFn, ok := cbr.executors[in.lang]
	if !ok {
		return fmt.Errorf("block %q language %q not supported", in.name, in.lang)
	}

	fmt.Fprintf(os.Stdout, "running: %s (%s)\n", in.name, in.lang)

	err = execFn(in, cbr.env)
	if err != nil {
		return fmt.Errorf("code block %q (%s) execution failed: %w", in.name, in.lang, err)
	}

	lines := strings.Split(in.output, "\n")
	for i := range lines {
		ln := strings.TrimSpace(lines[i])
		fmt.Fprintln(os.Stdout, text.Indent(text.Wrap(ln, 60), "  » "))
	}
	fmt.Fprintln(os.Stdout)

	if exp := in.params["export"]; exp != "" {
		for key := range strings.SplitSeq(exp, ",") {
			key = strings.TrimSpace(key)
			if key != "" {
				cbr.env[key] = in.output
			}
		}
	}

	// Scrivi su file se richiesto
	if filePath := in.params["file"]; filePath != "" {
		err := os.WriteFile(filePath, []byte(in.output), 0644)
		if err != nil {
			return fmt.Errorf("unable to write output to file %q: %w", filePath, err)
		}
	}

	return nil
}
