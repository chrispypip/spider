package spider

import (
	"fmt"
	"slices"

	"github.com/godbus/dbus/v5"
)

const (
	deviceInterface = "net.connman.iwd.Device"
	devicePropertyName = deviceInterface + ".Name"
	devicePropertyAddress = deviceInterface + ".Address"
	devicePropertyPowered = deviceInterface + ".Powered"
	devicePropertyAdapter = deviceInterface + ".Adapter"
	devicePropertyMode = deviceInterface + ".Mode"
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
	conn *dbus.Conn
	path dbus.ObjectPath
	obj dbus.BusObject
	name string
	address string
	powered bool
	adapter Adapter
	mode string
}

func NewDevice(conn *dbus.Conn, path dbus.ObjectPath) (*Device, error) {
	obj := conn.Object(IwdService, path)
	d := &Device{
		conn: conn,
		path: path,
		obj: obj,
	}

	if variant, err := d.obj.GetProperty(devicePropertyName); err != nil {
		fmt.Printf("Failed to get name property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&d.name); err2 != nil {
			fmt.Printf("Failed to store name property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := d.obj.GetProperty(devicePropertyAddress); err != nil {
		fmt.Printf("Failed to get address property: %s\n", err)
		// log error
		return nil, err
	} else {
		if err2 := variant.Store(&d.address); err2 != nil {
			fmt.Printf("Failed to store address property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := d.obj.GetProperty(devicePropertyPowered); err != nil {
		fmt.Printf("Failed to get powered property: %s\n", err)
		// log error
		return nil, err
	} else {
		if err2 := variant.Store(&d.powered); err2 != nil {
			fmt.Printf("Failed to store powered property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := d.obj.GetProperty(devicePropertyAdapter); err != nil {
		fmt.Printf("Failed to get adapter property: %s\n", err)
		// log error
		return nil, err
	} else {
		var path dbus.ObjectPath
		if err2 := variant.Store(&path); err2 != nil {
			fmt.Printf("Failed to store adapter path property: %s\n", err2)
			// log err
			return nil, err2
		}
		if adapter, err2 := NewAdapter(d.conn, path); err2 != nil {
			fmt.Printf("Failed to create new adapter for device: %s\n", err2)
			return nil, err2
		} else {
			d.adapter = *adapter
		}
	}
	if variant, err := d.obj.GetProperty(devicePropertyMode); err != nil {
		fmt.Printf("Failed to get mode property: %s\n", err)
		// log error
		return nil, err
	} else {
		if err2 := variant.Store(&d.mode); err2 != nil {
			fmt.Printf("Failed to store mode property: %s\n", err2)
			// log err
			return nil, err2
		}
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
		// log err
		return err
	} else {
		d.powered = powered
		return nil
	}
}

func (d *Device) SetMode(mode string) error {
	supportedModes := d.GetSupportedModes()
	if !slices.Contains(supportedModes, mode) {
		return fmt.Errorf("Invalid mode %s. Valid modes are %v", mode, supportedModes)
	}

	if err := d.obj.SetProperty(devicePropertyMode, dbus.MakeVariant(mode)); err != nil {
		// log err
		return err
	} else {
		d.mode = mode
		return nil
	}
}

func (d *Device) GetSupportedModes() []string {
	return []string{"ad-hoc", "station", "ap"}
}
