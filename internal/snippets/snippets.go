package snippets

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type executorFunc func(nfo *codeBlockInfo, env map[string]string) error

func execWithOptions(cmd *exec.Cmd, dir string, timeoutSec int, env map[string]string) (string, error) {
	var out bytes.Buffer

	var envSlice []string
	for k, v := range env {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Env = envSlice
	cmd.Stdout = &out
	cmd.Stderr = &out
	if dir != "" {
		cmd.Dir = dir
	}

	runWithTimeout := func() error {
		if timeoutSec <= 0 {
			return cmd.Run()
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
		defer cancel()

		name := cmd.Path
		args := cmd.Args[1:]
		cmdCtx := exec.CommandContext(ctx, name, args...)
		cmdCtx.Stdout = &out
		cmdCtx.Stderr = &out
		cmdCtx.Dir = dir
		cmdCtx.Env = envSlice

		err := cmdCtx.Run()
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout of %ds exceeded", timeoutSec)
		}
		return err
	}

	err := runWithTimeout()

	if err != nil {
		trimmed := strings.TrimSpace(out.String())
		if trimmed != "" {
			return trimmed, fmt.Errorf("%v\n%s", err, trimmed)
		}
		return trimmed, err
	}

	return strings.TrimSpace(out.String()), nil
}

func printEnvSlice(env []string) {
	sort.Strings(env)

	for _, el := range env {
		fmt.Fprintln(os.Stdout, el)
	}

	fmt.Fprintln(os.Stdout)

}
