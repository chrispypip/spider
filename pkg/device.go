package spider

import (
	"fmt"
	"slices"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	deviceInterface       = "net.connman.iwd.Device"
	devicePropertyName    = deviceInterface + ".Name"
	devicePropertyAddress = deviceInterface + ".Address"
	devicePropertyPowered = deviceInterface + ".Powered"
	devicePropertyAdapter = deviceInterface + ".Adapter"
	devicePropertyMode    = deviceInterface + ".Mode"
)

var (
	deviceLogger *log.Entry
)

type Devicer interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	GetName() string
	GetAddress() string
	GetPowered() bool
	GetAdapter() Adapter
	GetMode() string
	SetPowered(powered bool) error
	SetMode(mode string) error
	GetSupportedModes() []string
}

type Device struct {
	conn    *dbus.Conn
	path    dbus.ObjectPath
	obj     dbus.BusObject
	name    string
	address string
	powered bool
	adapter Adapter
	mode    string
}

func NewDevice(conn *dbus.Conn, path dbus.ObjectPath) (*Device, error) {
	log.SetReportCaller(true)
	deviceLogger = log.WithFields(log.Fields{
		"type": "Device",
		"path": path,
	})
	obj := conn.Object(IwdService, path)
	d := &Device{
		conn: conn,
		path: path,
		obj:  obj,
	}
	if variant, err := d.obj.GetProperty(devicePropertyName); err != nil {
		deviceLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'name'")
		return nil, err
	} else {
		if err2 := variant.Store(&d.name); err2 != nil {
			deviceLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'name'")
			return nil, err2
		}
		deviceLogger.Debugf("Name = %s", d.name)
	}
	if variant, err := d.obj.GetProperty(devicePropertyAddress); err != nil {
		deviceLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'address'")
		return nil, err
	} else {
		if err2 := variant.Store(&d.address); err2 != nil {
			deviceLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'address'")
			return nil, err2
		}
		deviceLogger.Debugf("Address = %s", d.address)
	}
	if variant, err := d.obj.GetProperty(devicePropertyPowered); err != nil {
		deviceLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'powered'")
		return nil, err
	} else {
		if err2 := variant.Store(&d.powered); err2 != nil {
			deviceLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'powered'")
			return nil, err2
		}
		deviceLogger.Debugf("Powered = %t", d.powered)
	}
	if variant, err := d.obj.GetProperty(devicePropertyAdapter); err != nil {
		deviceLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'adapter'")
		return nil, err
	} else {
		var path dbus.ObjectPath
		if err2 := variant.Store(&path); err2 != nil {
			deviceLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'adapter'")
			return nil, err2
		}
		if adapter, err2 := NewAdapter(d.conn, path); err2 != nil {
			deviceLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to created to Adapter from property 'adapter'")
			return nil, err2
		} else {
			deviceLogger.Debugf("Created new Adapter from property 'adapter': %s", adapter)
			d.adapter = *adapter
		}
	}
	if variant, err := d.obj.GetProperty(devicePropertyMode); err != nil {
		deviceLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'mode'")
		return nil, err
	} else {
		if err2 := variant.Store(&d.mode); err2 != nil {
			deviceLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'mode'")
			return nil, err2
		}
		deviceLogger.Debugf("Mode = %s", d.mode)
	}
	return d, nil
}

func (d Device) String() string {
	return fmt.Sprintf("{Path: %s, Interface: %s, Name: %s, Address: %s, Powered: %t, Adapter: %s, Mode: %s}", d.path, d.GetInterface(), d.name, d.address, d.powered, d.adapter, d.mode)
}

func (d *Device) GetPath() dbus.ObjectPath {
	return d.path
}

func (d *Device) GetInterface() string {
	return deviceInterface
}

func (d *Device) GetName() string {
	return d.name
}

func (d *Device) GetAddress() string {
	return d.address
}

func (d *Device) GetPowered() bool {
	return d.powered
}

func (d *Device) GetAdapter() Adapter {
	return d.adapter
}

func (d *Device) GetMode() string {
	return d.mode
}

func (d *Device) SetPowered(powered bool) error {
	obj := d.conn.Object(IwdService, d.path)
	if err := obj.SetProperty(devicePropertyPowered, dbus.MakeVariant(powered)); err != nil {
		deviceLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to set property 'powered' to %t", powered)
		return err
	} else {
		deviceLogger.Debugf("SetPowered %t", powered)
		d.powered = powered
		return nil
	}
}

func (d *Device) SetMode(mode string) error {
	supportedModes := d.GetSupportedModes()
	if !slices.Contains(supportedModes, mode) {
		err := fmt.Errorf("invalid mode %s. Valid modes are %v", mode, supportedModes)
		deviceLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("invalid mode %s; valid modes are %v", mode, supportedModes)
		return err
	}
	if err := d.obj.SetProperty(devicePropertyMode, dbus.MakeVariant(mode)); err != nil {
		deviceLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to set property 'mode' to %s", mode)
		return err
	} else {
		deviceLogger.Debugf("SetMode %s", mode)
		d.mode = mode
		return nil
	}
}

func (d *Device) GetSupportedModes() []string {
	return []string{"ad-hoc", "station", "ap"}
}
