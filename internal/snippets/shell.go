package snippets

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func runShell(nfo *codeBlockInfo, env map[string]string) error {
	timeout := 0
	if t, ok := nfo.params["timeout"]; ok {
		if n, err := strconv.Atoi(t); err == nil {
			timeout = n
		}
	}

	dir := nfo.params["dir"]

	// Step 1: costruiamo un dump iniziale dell'ambiente (prima di eseguire il blocco)
	preCmd := exec.Command("bash", "-c", "env")
	preDump, err := execWithOptions(preCmd, dir, 0, env)
	if err != nil {
		return err
	}
	initialVars := parseEnvDump(preDump)

	// Step 2: costruiamo il comando bash che:
	// - esegue il codice del blocco
	// - stampa il marker
	// - stampa il nuovo ambiente
	cmdStr := fmt.Sprintf(`set -a
%s
echo "__ENV_DUMP_START__"
env
`, nfo.text)

	cmd := exec.Command("bash", "-c", cmdStr)

	output, err := execWithOptions(cmd, dir, timeout, env)
	if err != nil {
		return err
	}

	// Step 3: separiamo output “normale” da quello dell’ambiente
	parts := strings.SplitN(output, "__ENV_DUMP_START__", 2)
	nfo.output = strings.TrimSpace(parts[0])

	if len(parts) == 2 {
		envDump := strings.TrimSpace(parts[1])
		newVars := parseEnvDump(envDump)

		// Step 4: aggiorniamo solo le variabili che sono nuove o cambiate
		for k, v := range newVars {
			if oldVal, exists := initialVars[k]; !exists || oldVal != v {
				env[k] = v
			}
		}
	}

	return nil
}

// parseEnvDump converte un output tipo "KEY=VAL\n..." in una mappa
func parseEnvDump(dump string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(dump, "\n") {
		line = strings.TrimSpace(line)
		if kv := strings.SplitN(line, "=", 2); len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}
