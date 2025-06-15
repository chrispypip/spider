package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	agentManagerInterface                                 = "net.connman.iwd.AgentManager"
	agentManagerMethodRegisterAgent                       = agentManagerInterface + ".RegisterAgent"
	agentManagerMethodUnregisterAgent                     = agentManagerInterface + ".UnregisterAgent"
	agentManagerMethodRegisterNetworkConfigurationAgent   = agentManagerInterface + ".RegisterNetworkConfigurationAgent"
	agentManagerMethodUnregisterNetworkConfigurationAgent = agentManagerInterface + ".UnregisterNetworkConfigurationAgent"
)

var (
	amLogger *log.Entry
)

type AgentManagerer interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	RegisterAgent(agent AgentClient) error
	UnregisterAgent(agent AgentClient) error
	RegisterNetworkConfigurationAgent(agent NetworkConfigurationAgentClient) error
	UnregisterNetworkConfigurationAgent(agent NetworkConfigurationAgentClient) error
}

type AgentManager struct {
	conn *dbus.Conn
	obj  dbus.BusObject
	path dbus.ObjectPath
}

func NewAgentManager(conn *dbus.Conn, path dbus.ObjectPath) *AgentManager {
	log.SetReportCaller(true)
	amLogger = log.WithFields(log.Fields{
		"type": "AgentManager",
		"path": path,
	})
	obj := conn.Object(IwdService, path)
	am := &AgentManager{
		conn: conn,
		obj:  obj,
		path: path,
	}
	return am
}

func (am *AgentManager) GetPath() dbus.ObjectPath {
	return am.path
}

func (am *AgentManager) GetInterface() string {
	return agentManagerInterface
}

func (am *AgentManager) String() string {
	return fmt.Sprintf("{Path: %s, Interface: %s}", am.path, am.GetInterface())
}

func (am *AgentManager) RegisterAgent(client AgentClient) error {
	clientPath := client.GetPath()
	if err := am.obj.Call(agentManagerMethodRegisterAgent, 0, clientPath).Err; err != nil {
		amLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to register Agent %s", clientPath)
		return err
	}
	amLogger.Debugf("registered Agent %s", clientPath)
	return nil
}

func (am *AgentManager) UnregisterAgent(client AgentClient) error {
	clientPath := client.GetPath()
	if err := am.obj.Call(agentManagerMethodUnregisterAgent, 0, clientPath).Err; err != nil {
		amLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to unregister Agent %s", clientPath)
		return err
	}
	amLogger.Debugf("Unregistered Agent %s", clientPath)
	return nil
}

func (am *AgentManager) RegisterNetworkConfigurationAgent(client NetworkConfigurationAgentClient) error {
	clientPath := client.GetPath()
	if err := am.obj.Call(agentManagerMethodRegisterNetworkConfigurationAgent, 0, clientPath).Err; err != nil {
		amLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to register NetworkConfigurationAgent %s", clientPath)
		return err
	}
	amLogger.Debugf("Registered NetworkConfigurationAgent %s", clientPath)
	return nil
}

func (am *AgentManager) UnregisterNetworkConfigurationAgent(client NetworkConfigurationAgentClient) error {
	clientPath := client.GetPath()
	if err := am.obj.Call(agentManagerMethodUnregisterNetworkConfigurationAgent, 0, clientPath).Err; err != nil {
		amLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to unregister NetworkConfigurationAgent %s", clientPath)
		return err
	}
	amLogger.Debugf("unregistered NetworkConfigurationAgent %s", clientPath)
	return nil
}
