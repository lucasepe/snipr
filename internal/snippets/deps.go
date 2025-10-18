package snippets

import (
	"fmt"
	"strings"
)

func topologicalSort(blocks []*codeBlockInfo) ([]*codeBlockInfo, error) {
	// Mappa nome → blocco
	blockMap := make(map[string]*codeBlockInfo)
	for _, b := range blocks {
		blockMap[b.name] = b
	}

	visited := make(map[string]bool)
	tempMark := make(map[string]bool)
	result := []*codeBlockInfo{}

	var visit func(b *codeBlockInfo) error
	visit = func(b *codeBlockInfo) error {
		if tempMark[b.name] {
			return fmt.Errorf("dependency cycle detected at block %q", b.name)
		}
		if visited[b.name] {
			return nil
		}

		tempMark[b.name] = true
		for _, dep := range strings.FieldsFunc(b.params["depends"], func(r rune) bool { return r == ',' || r == ' ' }) {
			dep = strings.TrimSpace(dep)
			if dep == "" {
				continue
			}
			depBlock, ok := blockMap[dep]
			if !ok {
				return fmt.Errorf("block %q depends on unknown block %q", b.name, dep)
			}
			if err := visit(depBlock); err != nil {
				return err
			}
		}

		tempMark[b.name] = false
		visited[b.name] = true

		result = append(result, b)
		return nil
	}

	for _, b := range blocks {
		if !visited[b.name] {
			if err := visit(b); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}
