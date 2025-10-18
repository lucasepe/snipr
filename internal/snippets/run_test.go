package snippets_test

import (
	"strings"
	"testing"

	"github.com/lucasepe/snipr/internal/snippets"
)

func TestRun(t *testing.T) {
	md := []string{
		"# Test Markdown Script",
		"",
		"```shell {name=setup export=TEST_DIR}",
		"echo \"/tmp/testdir\"",
		"```",
		"",
		"```shell {name=build depends=setup}",
		"cd $TEST_DIR",
		"echo \"Build in $(pwd)\"",
		"```",
		"",
		"```go {name=hello depends=build}",
		"package main",
		"import \"fmt\"",
		"func main() {",
		"    fmt.Println(\"Hello world\")",
		"}",
		"```",
		"",
		"```shell {name=skipped depends=nonexistent}",
		"echo \"This should not run\"",
		"```",
	}

	// Create an io.Reader for the Markdown
	reader := strings.NewReader(strings.Join(md, "\n"))

	// Run the Markdown script
	err := snippets.Run(reader)
	if err == nil {
		t.Fatal("expected error for skipped block due to unsatisfied dependencies, got nil")
	}

	// Optional: check error message contains skipped block name
	if !strings.Contains(err.Error(), "skipped") {
		t.Errorf("expected error message to mention skipped block, got: %v", err)
	}
}
