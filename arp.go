package arp

import (
	"time"
)

type ArpTableEntry struct {
	MAC  string
	Line string
}

type ArpTable map[string][]ArpTableEntry

var (
	stop     = make(chan struct{})
	arpCache = &cache{
		table: make(ArpTable),
	}
)

func AutoRefresh(t time.Duration) {
	go func() {
		for {
			select {
			case <-time.After(t):
				arpCache.Refresh()
			case <-stop:
				return
			}
		}
	}()
}

func StopAutoRefresh() {
	stop <- struct{}{}
}

func CacheUpdate() {
	arpCache.Refresh()
}

func CacheLastUpdate() time.Time {
	return arpCache.Updated
}

func CacheUpdateCount() int {
	return arpCache.UpdatedCount
}

// Search looks up the MAC address for an IP address
// in the arp table
func Search(ip string) string {
	macs := arpCache.Search(ip)
	if len(macs) == 0 {
		return ""
	}

	//For reverse compatability, return the last entry
	return macs[len(macs)-1]

}

func SearchEntries(ip string) []ArpTableEntry {
	return arpCache.SearchEntries(ip)
}
