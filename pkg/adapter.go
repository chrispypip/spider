package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	adapterInterface              = "net.connman.iwd.Adapter"
	adapterPropertyPowered        = adapterInterface + ".Powered"
	adapterPropertyName           = adapterInterface + ".Name"
	adapterPropertyModel          = adapterInterface + ".Model"
	adapterPropertyVendor         = adapterInterface + ".Vendor"
	adapterPropertySupportedModes = adapterInterface + ".SupportedModes"
)

var (
	adapterLogger *log.Entry
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
	conn           *dbus.Conn
	path           dbus.ObjectPath
	obj            dbus.BusObject
	name           string
	model          *string
	vendor         *string
	supportedModes []string
	powered        bool
}

func NewAdapter(conn *dbus.Conn, path dbus.ObjectPath) (*Adapter, error) {
	log.SetReportCaller(true)
	adapterLogger = log.WithFields(log.Fields{
		"type": "Adapter",
		"path": path,
	})
	obj := conn.Object(IwdService, path)
	a := &Adapter{
		conn: conn,
		path: path,
		obj:  obj,
	}
	if variant, err := a.obj.GetProperty(adapterPropertyName); err != nil {
		adapterLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'name'")
		return nil, err
	} else {
		if err2 := variant.Store(&a.name); err2 != nil {
			adapterLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'name'")
			return nil, err2
		}
		adapterLogger.Debugf("Name = %s", a.name)
	}
	if variant, err := a.obj.GetProperty(adapterPropertyModel); err != nil {
		adapterLogger.WithFields(log.Fields{
			"err": err,
		}).Info("Failed to get optional property 'model'")
		a.model = nil
	} else {
		if err2 := variant.Store(&a.model); err2 != nil {
			adapterLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'model'")
			return nil, err2
		}
		adapterLogger.Debugf("Model = %s", *a.model)
	}
	if variant, err := a.obj.GetProperty(adapterPropertyVendor); err != nil {
		adapterLogger.WithFields(log.Fields{
			"err": err,
		}).Info("Failed to get optional property 'vendor'")
		a.vendor = nil
	} else {
		if err2 := variant.Store(&a.vendor); err2 != nil {
			adapterLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'vendor'")
			return nil, err2
		}
		adapterLogger.Debugf("Vendor = %s", *a.vendor)
	}
	if variant, err := a.obj.GetProperty(adapterPropertySupportedModes); err != nil {
		adapterLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'supported modes'")
		return nil, err
	} else {
		if err2 := variant.Store(&a.supportedModes); err2 != nil {
			adapterLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'supported modes'")
			return nil, err2
		}
		adapterLogger.Debugf("Supported Modes = %v", a.supportedModes)
	}
	if variant, err := a.obj.GetProperty(adapterPropertyPowered); err != nil {
		adapterLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'powered'")
		return nil, err
	} else {
		if err2 := variant.Store(&a.powered); err2 != nil {
			adapterLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'powered'")
			return nil, err2
		}
		adapterLogger.Debugf("Powered = %t", a.powered)
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
		adapterLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to set property 'powered'")
		return err
	}
	adapterLogger.Debugf("SetPowered %t", powered)
	a.powered = powered
	return nil
}
