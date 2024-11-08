//go:build windows
// +build windows

package arp

// Windows arp table reader added by Claudio Matsuoka.
// Tested only in Windows 8.1, hopefully the arp command output format
// is the same in other Windows versions.

import (
	"os/exec"
	"strings"
)

func Table() ArpTable {
	data, err := exec.Command("arp", "-a").Output()
	if err != nil {
		return nil
	}

	var table = make(ArpTable)
	skipNext := false
	for _, line := range strings.Split(string(data), "\n") {
		// skip empty lines
		if len(line) <= 0 {
			continue
		}
		// skip Interface: lines
		if line[0] != ' ' {
			skipNext = true
			continue
		}
		// skip column headers
		if skipNext {
			skipNext = false
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		ip := fields[0]

		if _, ok := table[ip]; !ok {
			table[ip] = []ArpTableEntry{}
		}

		table[ip] = append(table[ip], ArpTableEntry{
			MAC:  strings.Replace(fields[1], "-", ":", -1), // Normalize MAC address to colon-separated format
			Line: line,
		})

	}

	return table
}
