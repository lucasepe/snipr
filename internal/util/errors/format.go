package errors

import (
	"fmt"
	"io"
	"strings"
)

// FormatJoinedErrors stampa tutti gli errori combinati (anche con errors.Join)
// in modo leggibile, uno per riga, con indentazione per le linee multilinea.
func FormatJoinedErrors(w io.Writer, prefix string, err error) {
	if err == nil {
		return
	}

	// Tipo per controllare se è un errors.Join
	type joiner interface {
		Unwrap() []error
	}

	if je, ok := err.(joiner); ok {
		for _, e := range je.Unwrap() {
			FormatJoinedErrors(w, prefix, e)
		}
		return
	}

	// errore singolo → splitta le linee
	lines := strings.Split(err.Error(), "\n")
	for i, line := range lines {
		linePrefix := prefix
		if i > 0 {
			// per le linee successive, allinea sotto il prefisso
			linePrefix = strings.Repeat(" ", len(prefix))
		}
		fmt.Fprintf(w, "%s%s\n", linePrefix, line)
	}
}
