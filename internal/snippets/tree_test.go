package snippets_test

import (
	"os"
	"testing"

	"github.com/lucasepe/snipr/internal/snippets"
)

func TestTree(t *testing.T) {

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
	_ = md
	//reader := strings.NewReader(strings.Join(md, "\n"))

	reader, err := os.Open("../../testdata/basic.md")
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	err = snippets.Tree(reader, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}

}
