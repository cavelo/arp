//go:build darwin

// only tested on OSX
// decided to go with exec.Command after I couldn't figure
// out how to extract the arp cache out of the kernel with
// golang's syscall or Sysctl()
//
// ... Help appreciated :)

package arp

import (
	"os/exec"
	"strings"
)

func Table() ArpTable {
	data, err := exec.Command("arp", "-an").Output()
	if err != nil {
		return nil
	}

	var table = make(ArpTable)
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// strip brackets around IP
		ip := strings.Replace(fields[1], "(", "", -1)
		ip = strings.Replace(ip, ")", "", -1)

		if _, ok := table[ip]; !ok {
			table[ip] = []ArpTableEntry{}
		}

		table[ip] = append(table[ip], ArpTableEntry{
			MAC:  fields[3],
			Line: line,
		})

	}

	return table
}
