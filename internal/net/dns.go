package net

import (
	"math/rand"
	"time"
)

type Resolver func(hosts []string) []string

type HostList struct {
	hostList []string
	i        uint32
}

func NewHostList(hosts []string) *HostList {
	rand.Seed(time.Now().UnixNano())
	newHosts := hosts
	rand.Shuffle(len(newHosts), func(i, j int) {
		newHosts[i], newHosts[j] = newHosts[j], newHosts[i]
	})
	return &HostList{
		hostList: newHosts,
		i:        0,
	}
}
func (hosts *HostList) Next() string {
	hosts.i++
	return hosts.hostList[hosts.i%uint32(len(hosts.hostList))]
}
