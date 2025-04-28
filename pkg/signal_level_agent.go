package spider

import (
	"github.com/godbus/dbus/v5"
)

const (
	SignalLevelAgentInterface = "net.connman.iwd.SignalLevelAgent"
)

type SignalLevelAgentClient interface {
	GetConnection() *dbus.Conn
	GetPath() dbus.ObjectPath
	Release(device dbus.ObjectPath) *dbus.Error
	Changed(device dbus.ObjectPath, level uint8) *dbus.Error
}
