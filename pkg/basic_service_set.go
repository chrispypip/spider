package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	basicServiceSetInterface       = "net.connman.iwd.BasicServiceSet"
	basicServiceSetPropertyAddress = basicServiceSetInterface + ".Address"
)

var (
	bssLogger *log.Entry
)

type BasicServiceSetter interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	GetAddress() string
}

type BasicServiceSet struct {
	conn    *dbus.Conn
	obj     dbus.BusObject
	path    dbus.ObjectPath
	address string
}

func NewBasicServiceSet(conn *dbus.Conn, path dbus.ObjectPath) (*BasicServiceSet, error) {
	log.SetReportCaller(true)
	bssLogger = log.WithFields(log.Fields{
		"type": "BasicServiceSet",
		"path": path,
	})
	obj := conn.Object(IwdService, path)
	bss := &BasicServiceSet{
		conn: conn,
		obj:  obj,
		path: path,
	}
	if variant, err := bss.obj.GetProperty(basicServiceSetPropertyAddress); err != nil {
		bssLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'address'")
		return nil, err
	} else {
		if err2 := variant.Store(&bss.address); err2 != nil {
			bssLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'address'")
			return nil, err2
		}
		bssLogger.Debugf("Address = %s", bss.address)
	}
	return bss, nil
}

func (bss *BasicServiceSet) GetPath() dbus.ObjectPath {
	return bss.path
}

func (bss *BasicServiceSet) GetInterface() string {
	return basicServiceSetInterface
}

func (bss *BasicServiceSet) GetAddress() string {
	return bss.address
}

func (bss *BasicServiceSet) String() string {
	return fmt.Sprintf("{Path: %s, Interface: %s, Address: %s}", bss.path, bss.GetInterface(), bss.address)
}
