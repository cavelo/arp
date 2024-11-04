package arp

import (
	"sync"
	"time"
)

type cache struct {
	sync.RWMutex
	table ArpTable

	Updated      time.Time
	UpdatedCount int
}

func (c *cache) Refresh() {
	c.Lock()
	defer c.Unlock()

	c.table = Table()
	c.Updated = time.Now()
	c.UpdatedCount += 1
}

func (c *cache) Search(ip string) []string {
	c.RLock()
	defer c.RUnlock()

	entries, ok := c.table[ip]

	if !ok {
		c.RUnlock()
		c.Refresh()
		c.RLock()
		entries = c.table[ip]
	}

	macs := []string{}
	for _, entry := range entries {
		macs = append(macs, entry.MAC)
	}

	return macs
}

func (c *cache) SearchEntries(ip string) []ArpTableEntry {
	c.RLock()
	defer c.RUnlock()

	entries, ok := c.table[ip]

	if !ok {
		c.RUnlock()
		c.Refresh()
		c.RLock()
		entries = c.table[ip]
	}

	return entries
}
