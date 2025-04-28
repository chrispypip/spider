package spider

import (
	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	SimpleAgentPath                 = "/spider/simpleagent"
	SimpleAgentUser                 = "User"
	SimpleAgentPassphrase           = "Passphrase"
	SimpleAgentPrivateKeyPassphrase = "Passphrase"
)

var (
	saLogger *log.Entry
)

type Agent struct {
	conn                 *dbus.Conn
	path                 dbus.ObjectPath
	user                 string
	passphrase           string
	privateKeyPassphrase string
}

func NewSimpleAgent(conn *dbus.Conn) *Agent {
	log.SetReportCaller(true)
	saLogger = log.WithFields(log.Fields{
		"type": "SimpleAgent",
		"path": SimpleAgentPath,
	})
	return &Agent{
		conn:                 conn,
		path:                 SimpleAgentPath,
		user:                 SimpleAgentUser,
		passphrase:           SimpleAgentPassphrase,
		privateKeyPassphrase: SimpleAgentPrivateKeyPassphrase,
	}
}

func (a *Agent) SetUser(user string) {
	saLogger.Debugf("SetUser")
	a.user = user
}

func (a *Agent) SetPassphrase(passphrase string) {
	saLogger.Debugf("SetPassphrase")
	a.passphrase = passphrase
}

func (a *Agent) SetPrivateKeyPassphrase(passphrase string) {
	saLogger.Debugf("SetPrivateKeyPassphrase")
	a.privateKeyPassphrase = passphrase
}

func (a *Agent) GetConnection() *dbus.Conn {
	return a.conn
}

func (a *Agent) GetPath() dbus.ObjectPath {
	return a.path
}

func (a *Agent) Release() *dbus.Error {
	saLogger.Infof("Released SimpleAgent")
	return nil
}

func (a *Agent) RequestPassphrase(network dbus.ObjectPath) (string, *dbus.Error) {
	saLogger.Infof("Requesting passphrase for Network %s", network)
	return a.passphrase, nil
}

func (a *Agent) RequestPrivateKeyPassphrase(network dbus.ObjectPath) (string, *dbus.Error) {
	saLogger.Infof("Requesting private key passphrase for Network %s", network)
	return a.privateKeyPassphrase, nil
}

func (a *Agent) RequestUserNameAndPassword(network dbus.ObjectPath) (string, string, *dbus.Error) {
	saLogger.Infof("Requesting username and password for Network %s", network)
	return a.user, a.passphrase, nil
}

func (a *Agent) RequestUserPassword(network dbus.ObjectPath, user string) (string, *dbus.Error) {
	saLogger.Infof("Requesting user password for Network %s", network)
	return a.passphrase, nil
}

func (a *Agent) Cancel(reason string) *dbus.Error {
	saLogger.Infof("Canceling request because: %s", reason)
	return nil
}
