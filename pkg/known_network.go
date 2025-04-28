package spider

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	knownNetworkInterface                 = "net.connman.iwd.KnownNetwork"
	knownNetworkPropertyName              = knownNetworkInterface + ".Name"
	knownNetworkPropertyType              = knownNetworkInterface + ".Type"
	knownNetworkPropertyHidden            = knownNetworkInterface + ".Hidden"
	knownNetworkPropertyLastConnectedTime = knownNetworkInterface + ".LastConnectedTime"
	knownNetworkPropertyAutoConnect       = knownNetworkInterface + ".AutoConnect"
	knownNetworkMethodForget              = knownNetworkInterface + ".Forget"
)

var (
	knLogger            *log.Entry
	ErrNetworkForgotten = errors.New("network has been forgotten")
)

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
	conn              *dbus.Conn
	obj               dbus.BusObject
	path              dbus.ObjectPath
	name              string
	netType           string
	hidden            bool
	lastConnectedTime *string
	autoConnect       bool
	forgotten         bool
}

func NewKnownNetwork(conn *dbus.Conn, path dbus.ObjectPath) (*KnownNetwork, error) {
	log.SetReportCaller(true)
	knLogger = log.WithFields(log.Fields{
		"type": "KnownNetwork",
		"path": path,
	})
	obj := conn.Object(IwdService, path)
	kn := &KnownNetwork{
		conn:      conn,
		obj:       obj,
		path:      path,
		forgotten: false,
	}
	if variant, err := kn.obj.GetProperty(knownNetworkPropertyName); err != nil {
		knLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'name'")
		return nil, err
	} else {
		if err2 := variant.Store(&kn.name); err2 != nil {
			knLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'name'")
			return nil, err2
		}
		knLogger.Debugf("Name = %s", kn.name)
	}
	if variant, err := kn.obj.GetProperty(knownNetworkPropertyType); err != nil {
		knLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'type'")
		return nil, err
	} else {
		if err2 := variant.Store(&kn.netType); err2 != nil {
			knLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'type'")
			return nil, err2
		}
		knLogger.Debugf("Type = %s", kn.netType)
	}
	if variant, err := kn.obj.GetProperty(knownNetworkPropertyHidden); err != nil {
		knLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'hidden'")
		return nil, err
	} else {
		if err2 := variant.Store(&kn.hidden); err2 != nil {
			knLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'hidden'")
			return nil, err2
		}
		knLogger.Debugf("Hidden = %t", kn.hidden)
	}
	if variant, err := kn.obj.GetProperty(knownNetworkPropertyLastConnectedTime); err != nil {
		knLogger.WithFields(log.Fields{
			"err": err,
		}).Info("Failed to get optional property 'last connected time'")
		kn.lastConnectedTime = nil
	} else {
		if err2 := variant.Store(&kn.lastConnectedTime); err2 != nil {
			knLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'last connected time'")
			return nil, err2
		}
		knLogger.Debugf("Last Connected Time = %s", *kn.lastConnectedTime)
	}
	if variant, err := kn.obj.GetProperty(knownNetworkPropertyAutoConnect); err != nil {
		knLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'auto connect'")
		return nil, err
	} else {
		if err2 := variant.Store(&kn.autoConnect); err2 != nil {
			knLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'auto connect'")
			return nil, err2
		}
		knLogger.Debugf("Auto Connect = %t", kn.autoConnect)
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
		knLogger.WithFields(log.Fields{
			"err": ErrNetworkForgotten,
		}).Errorf("Network %s has been forgotten", kn.name)
		return "", ErrNetworkForgotten
	}
	knLogger.Debugf("GetName %s", kn.name)
	return kn.name, nil
}

func (kn *KnownNetwork) GetType() (string, error) {
	if kn.forgotten {
		return "", ErrNetworkForgotten
	}
	knLogger.Debugf("GetType %s", kn.netType)
	return kn.netType, nil
}

func (kn *KnownNetwork) GetHidden() (bool, error) {
	if kn.forgotten {
		knLogger.WithFields(log.Fields{
			"err": ErrNetworkForgotten,
		}).Errorf("Network %s has been forgotten", kn.name)
		return false, ErrNetworkForgotten
	}
	knLogger.Debugf("GetHidden %t", kn.hidden)
	return kn.hidden, nil
}

func (kn *KnownNetwork) GetLastConnectedTime() (*string, error) {
	if kn.forgotten {
		knLogger.WithFields(log.Fields{
			"err": ErrNetworkForgotten,
		}).Errorf("Network %s has been forgotten", kn.name)
		return nil, ErrNetworkForgotten
	}
	if kn.lastConnectedTime != nil {
		knLogger.Debugf("GetLastConnectedTime %s", *kn.lastConnectedTime)
	} else {
		knLogger.Debugf("No LastConnectedTime")
	}
	return kn.lastConnectedTime, nil
}

func (kn *KnownNetwork) GetAutoConnect() (bool, error) {
	if kn.forgotten {
		knLogger.WithFields(log.Fields{
			"err": ErrNetworkForgotten,
		}).Errorf("Network %s has been forgotten", kn.name)
		return false, ErrNetworkForgotten
	}
	knLogger.Debugf("GetAutoConnect %t", kn.autoConnect)
	return kn.autoConnect, nil
}

func (kn *KnownNetwork) SetAutoConnect(autoConnect bool) error {
	if kn.forgotten {
		knLogger.WithFields(log.Fields{
			"err": ErrNetworkForgotten,
		}).Errorf("Network %s has been forgotten", kn.name)
		return ErrNetworkForgotten
	}
	obj := kn.conn.Object(IwdService, kn.path)
	if err := obj.SetProperty(knownNetworkPropertyAutoConnect, dbus.MakeVariant(autoConnect)); err != nil {
		knLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to set property 'auto connect' to %t", autoConnect)
		return err
	} else {
		knLogger.Debugf("SetAutoConnect %t", autoConnect)
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
		knLogger.WithFields(log.Fields{
			"err": ErrNetworkForgotten,
		}).Errorf("Network %s has been forgotten", kn.name)
		return ErrNetworkForgotten
	}
	if err := kn.obj.Call(knownNetworkMethodForget, 0).Err; err != nil {
		knLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to call Forget")
		return err
	}
	knLogger.Debugf("Forgotten")
	kn.forgotten = true
	return nil
}
