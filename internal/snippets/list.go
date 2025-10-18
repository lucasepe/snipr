package snippets

import (
	"fmt"
	"io"

	"github.com/lucasepe/x/text/table"
)

func List(in io.Reader, wri io.Writer) (err error) {
	all, err := extractCodeBlocks(in)
	if err != nil {
		return err
	}

	tbl := table.New()
	tbl.Separator = "   "
	for _, el := range all {
		i := 0
		for k, v := range el.params {
			if i == 0 {
				tbl.AddRow(el.name, el.lang, fmt.Sprintf("%s=%s", k, v))
			} else {
				tbl.AddRow("", "", fmt.Sprintf("%s=%s", k, v))
			}
			i += 1
		}
	}

	_, err = fmt.Fprintln(wri, tbl.String())
	return
}
