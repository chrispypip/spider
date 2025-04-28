package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	networkInterface = "net.connman.iwd.Network"
	networkPropertyName = networkInterface + ".Name"
	networkPropertyConnected = networkInterface + ".Connected"
	networkPropertyDevice = networkInterface + ".Device"
	networkPropertyType = networkInterface + ".Type"
	networkPropertyKnownNetwork = networkInterface + ".KnownNetwork"
	networkPropertyExtendedServiceSet = networkInterface + ".ExtendedServiceSet"
	networkMethodConnect = networkInterface + ".Connect"
)

type Networker interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	GetName() string
	GetConnected() bool
	GetDevice() Device
	GetType() string
	GetKnownNetwork() *KnownNetwork
	GetExtendedServiceSet() *[]BasicServiceSet
	Connect() error
}

type Network struct {
	conn *dbus.Conn
	obj dbus.BusObject
	path dbus.ObjectPath
	name string
	connected bool
	device Device
	netType string
	knownNetwork *KnownNetwork
	extendedServiceSet *[]BasicServiceSet
}

func NewNetwork(conn *dbus.Conn, path dbus.ObjectPath) (*Network, error) {
	obj := conn.Object(IwdService, path)
	n := &Network{
		conn: conn,
		obj: obj,
		path: path,
	}

	if variant, err := n.obj.GetProperty(networkPropertyName); err != nil {
		fmt.Printf("Failed to get name property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&n.name); err2 != nil {
			fmt.Printf("Failed to store name property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := n.obj.GetProperty(networkPropertyConnected); err != nil {
		fmt.Printf("Failed to get connected property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&n.connected); err2 != nil {
			fmt.Printf("Failed to store connected property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := n.obj.GetProperty(networkPropertyDevice); err != nil {
		fmt.Printf("Failed to get device property: %s\n", err)
		// log err
		return nil, err
	} else {
		var path dbus.ObjectPath
		if err2 := variant.Store(&path); err2 != nil {
			fmt.Printf("Failed to store device path property: %s\n", err2)
			// log err
			return nil, err2
		}
		if device, err2 := NewDevice(n.conn, path); err2 != nil {
			fmt.Printf("Failed to create new device for network: %s\n", err2)
			return nil, err2
		} else {
			n.device = *device
		}
	}
	if variant, err := n.obj.GetProperty(networkPropertyType); err != nil {
		fmt.Printf("Failed to get type property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&n.netType); err2 != nil {
			fmt.Printf("Failed to store type property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := n.obj.GetProperty(networkPropertyKnownNetwork); err != nil {
		fmt.Printf("Failed to get known network property: %s\n", err)
		// log err
		n.knownNetwork = nil
	} else {
		var path dbus.ObjectPath
		if err2 := variant.Store(&path); err2 != nil {
			fmt.Printf("Failed to store known network path property: %s\n", err2)
			// log err
			return nil, err2
		}
		if knownNetwork, err2 := NewKnownNetwork(conn, path); err2 != nil {
			fmt.Printf("Failed to create new known network for network: %s\n", err2)
			return nil, err2
		} else {
			n.knownNetwork = knownNetwork
		}
	}
	if variant, err := n.obj.GetProperty(networkPropertyExtendedServiceSet); err != nil {
		fmt.Printf("Failed to get extended service set property: %s\n", err)
		// log err
		n.extendedServiceSet = nil
	} else {
		var paths []dbus.ObjectPath
		var ess = make([]BasicServiceSet, 0)
		if err2 := variant.Store(&paths); err2 != nil {
			fmt.Printf("Failed to store extended service set paths property: %s\n", err2)
			// log err
			return nil, err2
		}
		for _, path := range(paths) {
			if bss, err2 := NewBasicServiceSet(n.conn, path); err2 != nil {
				fmt.Printf("Failed to create new basic service set for network: %s\n", err2)
				return nil, err2
			} else {
				ess = append(ess, *bss)
			}
		}
		n.extendedServiceSet = &ess
	}

	return n, nil
}

func (n *Network) GetPath() dbus.ObjectPath {
	return n.path
}

func (n *Network) GetInterface() string {
	return networkInterface
}

func (n *Network) GetName() string {
	return n.name
}

func (n *Network) GetConnected() bool {
	return n.connected
}

func (n *Network) GetDevice() Device {
	return n.device
}

func (n *Network) GetType() string {
	return n.netType
}

func (n *Network) GetKnownNetwork() *KnownNetwork {
	return n.knownNetwork
}

func (n *Network) String() string {
	var kn string
	var ess string
	if n.knownNetwork == nil {
		kn = "<nil>"
	} else {
		kn = n.knownNetwork.String()
	}
	if n.extendedServiceSet == nil {
		ess = "<nil>"
	} else {
		ess = fmt.Sprintf("%v", n.extendedServiceSet)
	}
	return fmt.Sprintf("{Path: %s, Interface: %s, Name: %s, Type: %s, Connected: %t, Device: %s, KnownNetwork: %s, ExtendedServiceSet: %v}", n.path, n.GetInterface(), n.name, n.netType, n.connected, &n.device, kn, ess)
}

func (n *Network) Connect() error {
	if err := n.obj.Call(networkMethodConnect, 0).Err; err != nil {
		return err
	}
	n.connected = true
	return nil
}
