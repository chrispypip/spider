package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	adapterInterface = "net.connman.iwd.Adapter"
	adapterPropertyPowered = adapterInterface + ".Powered"
	adapterPropertyName = adapterInterface + ".Name"
	adapterPropertyModel = adapterInterface + ".Model"
	adapterPropertyVendor = adapterInterface + ".Vendor"
	adapterPropertySupportedModes = adapterInterface + ".SupportedModes"
)

type Adapterer interface {
 	GetPath() dbus.ObjectPath
 	GetInterface() string
 	GetName() string
 	GetModel() *string
 	GetVendor() *string
 	GetSupportedModes() []string
 	GetPowered() bool
 	SetPowered(powered bool) error
}


type Adapter struct {
	conn *dbus.Conn
	path dbus.ObjectPath
	obj dbus.BusObject
	name string
	model *string
	vendor *string
	supportedModes []string
	powered bool
}

func NewAdapter(conn *dbus.Conn, path dbus.ObjectPath) (*Adapter, error) {
	obj := conn.Object(IwdService, path)
	a := &Adapter{
		conn: conn,
		path: path,
		obj: obj,
	}

	if variant, err := a.obj.GetProperty(adapterPropertyName); err != nil {
		fmt.Printf("Failed to get name property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&a.name); err2 != nil {
			fmt.Printf("Failed to store name property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := a.obj.GetProperty(adapterPropertyModel); err != nil {
		fmt.Printf("Failed to get model property: %s\n", err)
		// log error
		a.model = nil
	} else {
		if err2 := variant.Store(&a.model); err2 != nil {
			fmt.Printf("Failed to store model property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := a.obj.GetProperty(adapterPropertyVendor); err != nil {
		fmt.Printf("Failed to get vendor property: %s\n", err)
		// log error
		a.vendor = nil
	} else {
		if err2 := variant.Store(&a.vendor); err2 != nil {
			fmt.Printf("Failed to store vendor property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := a.obj.GetProperty(adapterPropertySupportedModes); err != nil {
		fmt.Printf("Failed to get supported modes property: %s\n", err)
		// log error
		return nil, err
	} else {
		if err2 := variant.Store(&a.supportedModes); err2 != nil {
			fmt.Printf("Failed to store supported modes property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := a.obj.GetProperty(adapterPropertyPowered); err != nil {
		fmt.Printf("Failed to get powered  property: %s\n", err)
		// log error
		return nil, err
	} else {
		if err2 := variant.Store(&a.powered); err2 != nil {
			fmt.Printf("Failed to store powered property: %s\n", err2)
			// log err
			return nil, err2
		}
	}

	return a, nil
}

func (a Adapter) String() string {
	var model string
	var vendor string
	if a.model == nil {
		model = "<nil>"
	} else {
		model = *a.model
	}
	if a.vendor == nil {
		vendor = "<nil>"
	} else {
		vendor = *a.vendor
	}
	return fmt.Sprintf("{Path: %s, Interface: %s, Name: %s, Model: %s, Vendor: %s, SupportedModes: %v, Powered: %t}", a.path, a.GetInterface(), a.name, model, vendor, a.supportedModes, a.powered)
}

func (a *Adapter) GetPath() dbus.ObjectPath {
	return a.path
}

func (a *Adapter) GetInterface() string {
	return adapterInterface
}

func (a *Adapter) GetName() string {
	return a.name
}

func (a *Adapter) GetModel() *string {
	return a.model
}

func (a *Adapter) GetVendor() *string {
	return a.vendor
}

func (a *Adapter) GetSupportedModes() []string {
	return a.supportedModes
}

func (a *Adapter) GetPowered() bool {
	return a.powered
}

func (a *Adapter) SetPowered(powered bool) error {
	if err := a.obj.SetProperty(adapterPropertyPowered, dbus.MakeVariant(powered)); err != nil {
		// log err
		return err
	} else {
		a.powered = powered
		return nil
	}
}
