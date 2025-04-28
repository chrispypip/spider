package spider

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	knownNetworkInterface = "net.connman.iwd.KnownNetwork"
	knownNetworkPropertyName = knownNetworkInterface + ".Name"
	knownNetworkPropertyType = knownNetworkInterface + ".Type"
	knownNetworkPropertyHidden = knownNetworkInterface + ".Hidden"
	knownNetworkPropertyLastConnectedTime = knownNetworkInterface + ".LastConnectedTime"
	knownNetworkPropertyAutoConnect = knownNetworkInterface + ".AutoConnect"
	knownNetworkMethodForget = knownNetworkInterface + ".Forget"
)

var KnownNetworkError = errors.New("Network has been forgotten")

type KnownNetworker interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	GetName() (string, error)
	GetType() (string, error)
	GetHidden() (bool, error)
	GetLastConnectedTime() (string, error)
	GetAutoConnect() (bool, error)
	SetAutoConnect(autoConnect bool) error
	Forget() error
}

type KnownNetwork struct {
	conn *dbus.Conn
	obj dbus.BusObject
	path dbus.ObjectPath
	name string
	netType string
	hidden bool
	lastConnectedTime *string
	autoConnect bool
	forgotten bool
}

func NewKnownNetwork(conn *dbus.Conn, path dbus.ObjectPath) (*KnownNetwork, error) {
	obj := conn.Object(IwdService, path)
	kn := &KnownNetwork{
		conn: conn,
		obj: obj,
		path: path,
		forgotten: false,
	}

	if variant, err := kn.obj.GetProperty(knownNetworkPropertyName); err != nil {
		fmt.Printf("Failed to get name property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&kn.name); err2 != nil {
			fmt.Printf("Failed to store name property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := kn.obj.GetProperty(knownNetworkPropertyType); err != nil {
		fmt.Printf("Failed to get type property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&kn.netType); err2 != nil {
			fmt.Printf("Failed to store type property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := kn.obj.GetProperty(knownNetworkPropertyHidden); err != nil {
		fmt.Printf("Failed to get hidden property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&kn.hidden); err2 != nil {
			fmt.Printf("Failed to store hidden property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := kn.obj.GetProperty(knownNetworkPropertyLastConnectedTime); err != nil {
		fmt.Printf("Failed to get last connected time property: %s\n", err)
		// log err
		kn.lastConnectedTime = nil
	} else {
		if err2 := variant.Store(&kn.lastConnectedTime); err2 != nil {
			fmt.Printf("Failed to store last connected time property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := kn.obj.GetProperty(knownNetworkPropertyAutoConnect); err != nil {
		fmt.Printf("Failed to get auto connect property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&kn.autoConnect); err2 != nil {
			fmt.Printf("Failed to store auto connect property: %s\n", err2)
			// log err
			return nil, err2
		}
	}

	return kn, nil
}

func (kn *KnownNetwork) GetPath() dbus.ObjectPath {
	return kn.path
}

func (kn *KnownNetwork) GetInterface() string {
	return knownNetworkInterface
}

func (kn *KnownNetwork) GetName() (string, error) {
	if kn.forgotten {
		return "", KnownNetworkError
	}
	return kn.name, nil
}

func (kn *KnownNetwork) GetType() (string, error) {
	if kn.forgotten {
		return "", KnownNetworkError
	}
	return kn.netType, nil
}

func (kn *KnownNetwork) GetHidden() (bool, error) {
	if kn.forgotten {
		return false, KnownNetworkError
	}
	return kn.hidden, nil
}

func (kn *KnownNetwork) GetLastConnectedTime() (*string, error) {
	if kn.forgotten {
		return nil, KnownNetworkError
	}
	return kn.lastConnectedTime, nil
}

func (kn *KnownNetwork) GetAutoConnect() (bool, error) {
	if kn.forgotten {
		return false, KnownNetworkError
	}
	return kn.autoConnect, nil
}

func (kn *KnownNetwork) SetAutoConnect(autoConnect bool) error {
	if kn.forgotten {
		return KnownNetworkError
	}
	obj := kn.conn.Object(IwdService, kn.path)
	if err := obj.SetProperty(knownNetworkPropertyAutoConnect, dbus.MakeVariant(autoConnect)); err != nil {
		// log err
		return err
	} else {
		kn.autoConnect = autoConnect
		return nil
	}
}

func (kn *KnownNetwork) String() string {
	var lastConnectedTime string
	if kn.lastConnectedTime == nil {
		lastConnectedTime = "<nil>"
	} else {
		lastConnectedTime = *kn.lastConnectedTime
	}
	return fmt.Sprintf("{Path: %s, Interface: %s, Name: %s, Type: %s, Hidden: %t, LastConnectedTime: %s, AutoConnect: %t}", kn.path, kn.GetInterface(), kn.name, kn.netType, kn.hidden, lastConnectedTime, kn.autoConnect)
}

func (kn *KnownNetwork) Forget() error {
	if kn.forgotten {
		return KnownNetworkError
	}
	if err := kn.obj.Call(knownNetworkMethodForget, 0).Err; err != nil {
		return err
	}
	kn.forgotten = true
	return nil
}
