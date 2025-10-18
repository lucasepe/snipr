package snippets

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/lucasepe/snipr/internal/markdown"
)

// codeBlockInfo rappresenta un blocco di codice raccolto
type codeBlockInfo struct {
	name   string
	lang   string
	text   string
	params map[string]string
	output string
}

func shouldSkip(nfo *codeBlockInfo) bool {
	val, ok := nfo.params["skip"]
	if !ok {
		return false
	}

	return strings.EqualFold(val, "true")
}

func dependsOnSkipped(b *codeBlockInfo, results map[string]*codeBlockInfo) (bool, string) {
	// split sui separatori , o spazio (come usato altrove)
	for _, dep := range strings.FieldsFunc(b.params["depends"], func(r rune) bool { return r == ',' || r == ' ' }) {
		dep = strings.TrimSpace(dep)
		if dep == "" {
			continue
		}
		if depBlock, ok := results[dep]; ok {
			if shouldSkip(depBlock) {
				return true, dep
			}
		}
	}
	return false, ""
}

// blockWalker è una funzione chiamata su ogni nodo dell’albero Markdown
type blockWalker func(b markdown.Block)

// WalkBlocks attraversa ricorsivamente l’albero Markdown
func walkBlocks(b markdown.Block, fn blockWalker) {
	fn(b)
	switch n := b.(type) {
	case *markdown.Document:
		for _, child := range n.Blocks {
			walkBlocks(child, fn)
		}
	case *markdown.List:
		for _, child := range n.Items {
			walkBlocks(child, fn)
		}
	case *markdown.Quote:
		for _, child := range n.Blocks {
			walkBlocks(child, fn)
		}
	case *markdown.Paragraph, *markdown.Heading:
		// nessun figlio
	}
}

// collectCodeBlocks usa il BlockWalker per trovare tutti i blocchi di codice
func collectCodeBlocks(doc markdown.Block) (all []*codeBlockInfo) {
	walkBlocks(doc, func(b markdown.Block) {
		if cb, ok := b.(*markdown.CodeBlock); ok {
			lang, params := parseInfo(cb.Info)

			name := params["name"]
			if name == "" {
				name = fmt.Sprintf("anon-%d", len(all)+1)
			}
			all = append(all, &codeBlockInfo{
				name:   name,
				lang:   lang,
				text:   strings.Join(cb.Text, "\n"),
				params: params,
			})
		}
	})

	return
}

func extractCodeBlocks(in io.Reader) (all []*codeBlockInfo, err error) {
	buf := bytes.Buffer{}
	_, err = io.Copy(&buf, in)
	if err != nil {
		return nil, err
	}

	p := new(markdown.Parser)
	doc := p.Parse(buf.String())

	all = collectCodeBlocks(doc)

	return
}

// parseInfo analizza la stringa Info, estraendo linguaggio e parametri
func parseInfo(info string) (lang string, params map[string]string) {
	params = make(map[string]string)
	info = strings.TrimSpace(info)
	if info == "" {
		return "", params
	}

	// Esempio info: `go {run=false name=test}`
	re := regexp.MustCompile(`^([a-zA-Z0-9_+-]+)?\s*(?:\{(.*)\})?$`)
	m := re.FindStringSubmatch(info)
	if len(m) < 2 {
		return info, params
	}

	lang = strings.TrimSpace(m[1])
	if len(m) > 2 && m[2] != "" {
		for _, kv := range strings.Fields(m[2]) {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) == 2 {
				params[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}
	return lang, params
}
