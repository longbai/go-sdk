package storage

import (
	"github.com/longbai/go-sdk/internal/net"
)

type Config struct {
	PartSizeMb uint32
	UpHosts    []string
	CheckSum   bool
	Resolver   net.Resolver
	Recorder   UploadRecorder
}
