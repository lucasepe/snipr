package snippets

import (
	"fmt"
	"io"
	"strings"
)

func Tree(in io.Reader, wri io.Writer) (err error) {
	all, err := extractCodeBlocks(in)
	if err != nil {
		return err
	}

	printTree(all, wri)

	return nil
}

/*
func printTree(blocks []*codeBlockInfo, wri io.Writer) {
	// mappa nome -> blocco
	blockMap := make(map[string]*codeBlockInfo)
	for _, b := range blocks {
		blockMap[b.name] = b
	}

	// mappa nome -> dipendenze
	depMap := make(map[string][]string)
	for _, b := range blocks {
		for _, dep := range strings.FieldsFunc(b.params["depends"], func(r rune) bool { return r == ',' || r == ' ' }) {
			dep = strings.TrimSpace(dep)
			if dep != "" {
				depMap[b.name] = append(depMap[b.name], dep)
			}
		}
	}

	visited := make(map[string]bool)

	var walk func(name string, prefix string, last bool, isRoot bool)
	walk = func(name, prefix string, last, isRoot bool) {
		if visited[name] {
			fmt.Fprintln(wri, prefix+name+" (cycle?)")
			return
		}
		visited[name] = true

		blk, ok := blockMap[name]
		lang := "unknown"
		if ok {
			lang = blk.lang
		}
		display := fmt.Sprintf("%s (%s)", name, lang)

		// stampa root senza branch
		if isRoot {
			fmt.Fprintln(wri, display)
		} else {
			branch := "├── "
			if last {
				branch = "└── "
			}
			fmt.Fprintln(wri, prefix+branch+display)
		}

		// stampa params come elenco puntato
		if blk != nil && len(blk.params) > 0 {
			for key, value := range blk.params {
				fmt.Fprintf(wri, "%s    •  %s: %s\n", prefix, key, value)
			}
		}

		deps := depMap[name]
		for i, dep := range deps {
			isLast := i == len(deps)-1
			newPrefix := prefix
			if !isRoot {
				if last {
					newPrefix += "    "
				} else {
					newPrefix += "│   "
				}
			}
			// se la dipendenza non esiste
			if _, exists := blockMap[dep]; !exists {
				fmt.Fprintf(wri, "%s%s%s (missing!)\n", newPrefix, "└── ", dep)
				continue
			}
			walk(dep, newPrefix, isLast, false)
		}
	}

	// calcola root blocks (non dipendenza di altri)
	rootList := []*codeBlockInfo{}
	isDep := make(map[string]bool)
	for _, deps := range depMap {
		for _, d := range deps {
			isDep[d] = true
		}
	}
	for _, b := range blocks {
		if !isDep[b.name] {
			rootList = append(rootList, b)
		}
	}

	// stampa i root nell'ordine di apparizione
	for _, root := range rootList {
		walk(root.name, "", true, true)
	}
}
*/

func printTree(blocks []*codeBlockInfo, wri io.Writer) {
	// mappa nome -> blocco
	blockMap := make(map[string]*codeBlockInfo)
	for _, b := range blocks {
		blockMap[b.name] = b
	}

	// mappa nome -> dipendenze
	depMap := make(map[string][]string)
	for _, b := range blocks {
		for _, dep := range strings.FieldsFunc(b.params["depends"], func(r rune) bool { return r == ',' || r == ' ' }) {
			dep = strings.TrimSpace(dep)
			if dep != "" {
				depMap[b.name] = append(depMap[b.name], dep)
			}
		}
	}

	visited := make(map[string]bool)

	var walk func(name string, prefix string, last bool, isRoot bool)
	walk = func(name, prefix string, last bool, isRoot bool) {
		if visited[name] {
			fmt.Fprintln(wri, prefix+name+" (cycle?)")
			return
		}
		visited[name] = true

		blk, ok := blockMap[name]
		lang := "unknown"
		if ok {
			lang = blk.lang
		}
		display := fmt.Sprintf("%s (%s)", name, lang)

		if isRoot {
			fmt.Fprintln(wri, display)
		} else {
			branch := "├── "
			if last {
				branch = "└── "
			}
			fmt.Fprintln(wri, prefix+branch+display)
		}

		// Stampa i parametri come elenco puntato
		if blk != nil && len(blk.params) > 0 {
			for key, value := range blk.params {
				//fmt.Fprintf(wri, "%s    ├── %s: %s\n", prefix, key, value)
				fmt.Fprintf(wri, "%s    •  %s: %s\n", prefix, key, value)
			}
		}

		deps := depMap[name]
		for i, dep := range deps {
			isLast := i == len(deps)-1
			newPrefix := prefix
			if !isRoot {
				if last {
					newPrefix += "    "
				} else {
					newPrefix += "│   "
				}
			}
			walk(dep, newPrefix, isLast, false)
		}
	}

	// calcola root blocks (quelli che non sono dipendenza di altri)
	roots := make(map[string]bool)
	for _, b := range blocks {
		roots[b.name] = true
	}
	for _, deps := range depMap {
		for _, d := range deps {
			delete(roots, d)
		}
	}

	// stampa a partire dai root
	for root := range roots {
		walk(root, "", true, true)
	}
}
