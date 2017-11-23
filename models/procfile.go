// https://github.com/DarthSim/hivemind/blob/master/procfile.go

package models

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type procfileEntry struct {
	Name    string
	Command string
}

func parseProcfile(path string) (entries []procfileEntry) {
	re, _ := regexp.Compile(`^([\w-]+):\s+(.+)$`)

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
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

		entries = append(entries, procfileEntry{name, cmd})

	}

	if len(entries) == 0 {
		fmt.Println("No entries found")
	}

	return
}
