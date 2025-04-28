package spider

import (
	"github.com/godbus/dbus/v5"
)

const (
	NetworkConfigurationAgentInterface = "net.connman.iwd.NetworkConfigurationAgent"
)

type NetworkConfigurationRouteDestination struct {
	IP     string
	Length byte
}

type NetworkConfigurationAddress struct {
	Address           string
	PrefixLength      *byte
	Broadcast         *string
	ValidLifetime     *uint32
	PreferredLifetime *uint32
}

type NetworkConfigurationRoute struct {
	Destination     *NetworkConfigurationRouteDestination
	Router          *string
	PreferredSource *string
	Lifetime        *uint32
	Priority        uint32
	Preference      *byte
	MTU             *uint32
}

type NetworkConfigurationConfig struct {
	Method            string
	Addresses         []NetworkConfigurationAddress
	Routes            []NetworkConfigurationRoute
	DomainNameServers *[]string
	DomainNames       *[]string
	MDNS              *string
}

type NetworkConfigurationAgentClient interface {
	GetConnection() *dbus.Conn
	GetPath() dbus.ObjectPath
	Release() *dbus.Error
	//ConfigureIPv4(device dbus.ObjectPath, config map[string]dbus.Variant) *dbus.Error
	//ConfigureIPv6(device dbus.ObjectPath, config map[string]dbus.Variant) *dbus.Error
	CancelIPv4(device dbus.ObjectPath, reason string) *dbus.Error
	CancelIPv6(device dbus.ObjectPath, reason string) *dbus.Error
}
