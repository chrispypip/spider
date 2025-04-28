package spider

import (
	"github.com/godbus/dbus/v5"
)

const (
	AgentInterface = "net.connman.iwd.Agent"
)

type AgentClient interface {
	GetConnection() *dbus.Conn
	GetPath() dbus.ObjectPath
	Release() *dbus.Error
	RequestPassphrase(network dbus.ObjectPath) (string, *dbus.Error)
	RequestPrivateKeyPassphrase(network dbus.ObjectPath) (string, *dbus.Error)
	RequestUserNameAndPassword(network dbus.ObjectPath) (string, string, *dbus.Error)
	RequestUserPassword(network dbus.ObjectPath, user string) (string, *dbus.Error)
	Cancel(reason string) *dbus.Error
}
