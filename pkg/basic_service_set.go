package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	basicServiceSetInterface = "net.connman.iwd.BasicServiceSet"
	basicServiceSetPropertyAddress = basicServiceSetInterface + ".Address"
)

type BasicServiceSetter interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	GetAddress() string
}

type BasicServiceSet struct {
	conn *dbus.Conn
	obj dbus.BusObject
	path dbus.ObjectPath
	address string
}

func NewBasicServiceSet(conn *dbus.Conn, path dbus.ObjectPath) (*BasicServiceSet, error) {
	obj := conn.Object(IwdService, path)
	bss := &BasicServiceSet{
		conn: conn,
		obj: obj,
		path: path,
	}

	if variant, err := bss.obj.GetProperty(basicServiceSetPropertyAddress); err != nil {
		fmt.Printf("Failed to get address property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&bss.address); err2 != nil {
			fmt.Printf("Failed to store address property: %s\n", err2)
			// log err
			return nil, err2
		}
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
	return fmt.Sprintf("{Path: %s, Interface: %s, Address: %t}", bss.path, bss.GetInterface(), bss.address)
}
