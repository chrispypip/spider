package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	agentManagerInterface = "net.connman.iwd.AgentManager"
	agentManagerMethodRegisterAgent = agentManagerInterface + ".RegisterAgent"
	agentManagerMethodUnregisterAgent = agentManagerInterface + ".UnregisterAgent"
	agentManagerMethodRegisterNetworkConfigurationAgent = agentManagerInterface + ".RegisterNetworkConfigurationAgent"
	agentManagerMethodUnregisterNetworkConfigurationAgent = agentManagerInterface + ".UnregisterNetworkConfigurationAgent"
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
	obj dbus.BusObject
	path dbus.ObjectPath
}

func NewAgentManager(conn *dbus.Conn, path dbus.ObjectPath) *AgentManager {
	obj := conn.Object(IwdService, path)
	am := &AgentManager{
		conn: conn,
		obj: obj,
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

func (am *AgentManager) RegisterAgent(agent AgentClient) error {
	if err := am.obj.Call(agentManagerMethodRegisterAgent, 0, agent.GetPath()).Err; err != nil {
		return err
	}
	return nil
}

func (am *AgentManager) UnregisterAgent(agent AgentClient) error {
	if err := am.obj.Call(agentManagerMethodUnregisterAgent, 0, agent.GetPath()).Err; err != nil {
		return err
	}
	return nil
}

func (am *AgentManager) RegisterNetworkConfigurationAgent(agent NetworkConfigurationAgentClient) error {
	if err := am.obj.Call(agentManagerMethodRegisterNetworkConfigurationAgent, 0, agent.GetPath()).Err; err != nil {
		return err
	}
	return nil
}

func (am *AgentManager) UnregisterNetworkConfigurationAgent(agent NetworkConfigurationAgentClient) error {
	if err := am.obj.Call(agentManagerMethodUnregisterNetworkConfigurationAgent, 0, agent.GetPath()).Err; err != nil {
		return err
	}
	return nil
}
