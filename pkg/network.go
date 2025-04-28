package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	networkInterface                  = "net.connman.iwd.Network"
	networkPropertyName               = networkInterface + ".Name"
	networkPropertyConnected          = networkInterface + ".Connected"
	networkPropertyDevice             = networkInterface + ".Device"
	networkPropertyType               = networkInterface + ".Type"
	networkPropertyKnownNetwork       = networkInterface + ".KnownNetwork"
	networkPropertyExtendedServiceSet = networkInterface + ".ExtendedServiceSet"
	networkMethodConnect              = networkInterface + ".Connect"
)

var (
	networkLogger *log.Entry
)

type Networker interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	GetName() (string, error)
	GetConnected() bool
	GetDevice() Device
	GetType() (string, error)
	GetKnownNetwork() *KnownNetwork
	GetExtendedServiceSet() *[]BasicServiceSet
	Connect() error
}

type Network struct {
	conn               *dbus.Conn
	obj                dbus.BusObject
	path               dbus.ObjectPath
	name               string
	connected          bool
	device             Device
	netType            string
	knownNetwork       *KnownNetwork
	extendedServiceSet *[]BasicServiceSet
}

func NewNetwork(conn *dbus.Conn, path dbus.ObjectPath) (*Network, error) {
	log.SetReportCaller(true)
	networkLogger = log.WithFields(log.Fields{
		"type": "Network",
		"path": path,
	})
	obj := conn.Object(IwdService, path)
	n := &Network{
		conn: conn,
		obj:  obj,
		path: path,
	}
	if variant, err := n.obj.GetProperty(networkPropertyName); err != nil {
		networkLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'name'")
		return nil, err
	} else {
		if err2 := variant.Store(&n.name); err2 != nil {
			networkLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'name'")
			return nil, err2
		}
		networkLogger.Debugf("Name = %s", n.name)
	}
	if variant, err := n.obj.GetProperty(networkPropertyConnected); err != nil {
		networkLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'connected'")
		return nil, err
	} else {
		if err2 := variant.Store(&n.connected); err2 != nil {
			networkLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'connected'")
			return nil, err2
		}
		networkLogger.Debugf("Connected = %t", n.connected)
	}
	if variant, err := n.obj.GetProperty(networkPropertyDevice); err != nil {
		networkLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'device'")
		return nil, err
	} else {
		var path dbus.ObjectPath
		if err2 := variant.Store(&path); err2 != nil {
			networkLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'device'")
			return nil, err2
		}
		if device, err2 := NewDevice(n.conn, path); err2 != nil {
			networkLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to create new Device from property 'device'")
			return nil, err2
		} else {
			networkLogger.Debugf("Created new Device from property 'device': %s", device)
			n.device = *device
		}
	}
	if variant, err := n.obj.GetProperty(networkPropertyType); err != nil {
		networkLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'type'")
		return nil, err
	} else {
		if err2 := variant.Store(&n.netType); err2 != nil {
			networkLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'type'")
			return nil, err2
		}
		networkLogger.Debugf("Type = %s", n.netType)
	}
	if variant, err := n.obj.GetProperty(networkPropertyKnownNetwork); err != nil {
		networkLogger.WithFields(log.Fields{
			"err": err,
		}).Info("Failed to get optional property 'known network'")
		n.knownNetwork = nil
	} else {
		var path dbus.ObjectPath
		if err2 := variant.Store(&path); err2 != nil {
			networkLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'known network'")
			return nil, err2
		}
		if knownNetwork, err2 := NewKnownNetwork(conn, path); err2 != nil {
			networkLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to create new Known Network from property 'known network'")
			return nil, err2
		} else {
			networkLogger.Debugf("Created new Known Network from property 'known network': %s", knownNetwork)
			n.knownNetwork = knownNetwork
		}
	}
	if variant, err := n.obj.GetProperty(networkPropertyExtendedServiceSet); err != nil {
		networkLogger.WithFields(log.Fields{
			"err": err,
		}).Info("Failed to get optional property 'extended service set'")
		n.extendedServiceSet = nil
	} else {
		var paths []dbus.ObjectPath
		var ess = make([]BasicServiceSet, 0)
		if err2 := variant.Store(&paths); err2 != nil {
			networkLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'extended service set'")
			return nil, err2
		}
		for i, path := range paths {
			if bss, err2 := NewBasicServiceSet(n.conn, path); err2 != nil {
				networkLogger.WithFields(log.Fields{
					"err": err2,
				}).Error("Failed to create new Basic Service Set from property 'extended service set'")
				return nil, err2
			} else {
				networkLogger.Debugf("Created new Basic Service Set %d from property 'extended service set': %s", i, bss)
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

func (n *Network) GetName() (string, error) {
	return n.name, nil
}

func (n *Network) GetConnected() bool {
	return n.connected
}

func (n *Network) GetDevice() Device {
	return n.device
}

func (n *Network) GetType() (string, error) {
	return n.netType, nil
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
		networkLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to call Connect")
		return err
	}
	networkLogger.Debugf("Connected")
	n.connected = true
	return nil
}
