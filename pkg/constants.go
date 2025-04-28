package spider

import (
	"github.com/godbus/dbus/v5"
)

const (
	IwdService = "net.connman.iwd"
)

type ListedNetworker interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	GetName() (string, error)
	GetType() (string, error)
}
