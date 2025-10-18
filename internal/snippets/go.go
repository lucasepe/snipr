package snippets

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func runGo(workspace *goWorkspace) executorFunc {
	return func(nfo *codeBlockInfo, env map[string]string) (err error) {
		timeout := 0
		if t, ok := nfo.params["timeout"]; ok {
			if n, err := strconv.Atoi(t); err == nil {
				timeout = n
			}
		}

		// Write the code in a file
		blockName := nfo.params["name"]
		if blockName == "" {
			blockName = fmt.Sprintf("block-%d", time.Now().UnixNano())
		}
		srcPath := filepath.Join(workspace.dir, blockName+".go")

		err = os.WriteFile(srcPath, []byte(nfo.text), 0644)
		if err != nil {
			return
		}

		// Update deps
		err = workspace.tidy(timeout)
		if err != nil {
			return err
		}

		cmd := exec.Command("go", "run", srcPath)

		nfo.output, err = execWithOptions(cmd, workspace.dir, timeout, env)

		return
	}
}

type goWorkspace struct {
	dir string
}

func newGoWorkspace() (*goWorkspace, error) {
	dir, err := os.MkdirTemp("", "mdgo-*")
	if err != nil {
		return nil, err
	}

	modCmd := exec.Command("go", "mod", "init", "mdrun-temp")
	if out, err := execWithOptions(modCmd, dir, 0, nil); err != nil {
		return nil, fmt.Errorf("go mod init failed: %v\n%s", err, out)
	}

	return &goWorkspace{dir: dir}, nil
}

func (w *goWorkspace) tidy(timeout int) error {
	tidyCmd := exec.Command("go", "mod", "tidy")
	if out, err := execWithOptions(tidyCmd, w.dir, timeout, nil); err != nil {
		return fmt.Errorf("go mod tidy failed: %v\n%s", err, out)
	}
	return nil
}

func (w *goWorkspace) cleanup() {
	os.RemoveAll(w.dir)
}
