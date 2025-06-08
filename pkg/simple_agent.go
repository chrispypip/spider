package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	SimpleAgentPath = "/spider/simpleagent"
	SimpleAgentUser = "User"
	SimpleAgentPassphrase = "Passphrase"
	SimpleAgentPrivateKeyPassphrase = "Passphrase"
)

type Agent struct {
	conn *dbus.Conn
	path dbus.ObjectPath
	user string
	passphrase string
	privateKeyPassphrase string
}

func NewSimpleAgent(conn *dbus.Conn) *Agent {
	return &Agent{
		conn: conn,
		path: SimpleAgentPath,
		user: SimpleAgentUser,
		passphrase: SimpleAgentPassphrase,
		privateKeyPassphrase: SimpleAgentPrivateKeyPassphrase,
	}
}

func (a *Agent) SetUser(user string) {
	a.user = user
}

func (a *Agent) SetPassphrase(passphrase string) {
	a.passphrase = passphrase
}

func (a *Agent) SetPrivateKeyPassphrase(passphrase string) {
	a.privateKeyPassphrase = passphrase
}

func (a *Agent) GetConnection() *dbus.Conn {
	return a.conn
}

func (a *Agent) GetPath() dbus.ObjectPath {
	return a.path
}

func (a *Agent) Release() *dbus.Error {
	fmt.Printf("Releasing agent\n")
	return nil
}

func (a *Agent) RequestPassphrase(network dbus.ObjectPath) (string, *dbus.Error) {
	fmt.Printf("Requesting passhprase for %s\n", network)
	return a.passphrase, nil
}

func (a *Agent) RequestPrivateKeyPassphrase(network dbus.ObjectPath) (string, *dbus.Error) {
	fmt.Printf("Requesting private key passphrase for %s\n", network)
	return a.privateKeyPassphrase, nil
}

func (a *Agent) RequestUserNameAndPassword(network dbus.ObjectPath) (string, string, *dbus.Error) {
	fmt.Printf("Requesting user name and password for %s\n", network)
	return a.user, a.passphrase, nil
}

func (a *Agent) RequestUserPassword(network dbus.ObjectPath, user string) (string, *dbus.Error) {
	fmt.Printf("Requesting password for user %s for %s\n", user, network)
	return a.passphrase, nil
}

func (a *Agent) Cancel(reason string) *dbus.Error {
	fmt.Printf("Canceling request because: %s\n", reason)
	return nil
}
