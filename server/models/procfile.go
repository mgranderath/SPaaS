// https://github.com/DarthSim/hivemind/blob/master/procfile.go

package models

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// ProcfileEntry contains the procfile entry
type ProcfileEntry struct {
	Name    string
	Command string
}

// ParseProcfile parses yml files
func ParseProcfile(path string) (entries []ProcfileEntry) {
	re, _ := regexp.Compile(`^([\w-]+):\s+(.+)$`)

	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	names := make(map[string]bool)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}

		params := re.FindStringSubmatch(scanner.Text())
		if len(params) != 3 {
			continue
		}

		name, cmd := params[1], params[2]

		if names[name] {
			fmt.Println("must be unique")
		}
		names[name] = true

		entries = append(entries, ProcfileEntry{name, cmd})

	}

	return
}
