package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	adhocInterface = "net.connman.iwd.AdHoc"
	adhocPropertyConnectedPeers = adhocInterface + ".ConnectedPeers"
	adhocPropertyStarted = adhocInterface + ".Started"
	adhocMethodStart = adhocInterface + ".Start"
	adhocMethodStartOpen = adhocInterface + ".StartOpen"
	adhocMethodStop = adhocInterface + ".Stop"
)

type AdHocer interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	GetStarted() bool
	Start(ssid, psk string) error
	StartOpen(ssid string) error
	Stop() error
	GetConnectedPeers() ([]string, error)
}

type AdHoc struct {
	conn *dbus.Conn
	obj dbus.BusObject
	path dbus.ObjectPath
	connectedPeers []string
	started bool
}

func NewAdHoc(conn *dbus.Conn, path dbus.ObjectPath) (*AdHoc, error) {
	obj := conn.Object(IwdService, path)
	ah := &AdHoc{
		conn: conn,
		obj: obj,
		path: path,
	}

	if variant, err := ah.obj.GetProperty(adhocPropertyStarted); err != nil {
		fmt.Printf("Failed to get started property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&ah.started); err2 != nil {
			fmt.Printf("Failed to store started property: %s\n", err2)
			// log err
			return nil, err2
		}
	}

	return ah, nil
}

func (ah *AdHoc) GetPath() dbus.ObjectPath {
	return ah.path
}

func (ah *AdHoc) GetInterface() string {
	return adhocInterface
}

func (ah *AdHoc) GetStarted() bool {
	return ah.started
}

func (ah *AdHoc) String() string {
	connectedPeers, _ := ah.GetConnectedPeers()
	return fmt.Sprintf("{Path: %s, Interface: %s, Started: %t, ConnectedPeers: %v}", ah.path, ah.GetInterface(), ah.started, connectedPeers)
}

func (ah *AdHoc) Start(ssid, psk string) error {
	if err := ah.obj.Call(adhocMethodStart, 0, ssid, psk).Err; err != nil {
		return err
	}
	
	ah.started = true
	return nil
}

func (ah *AdHoc) StartOpen(ssid string) error {
	if err := ah.obj.Call(adhocMethodStartOpen, 0).Err; err != nil {
		return err
	}

	ah.started = true
	return nil
}

func (ah *AdHoc) Stop() error {
	if err := ah.obj.Call(adhocMethodStop, 0).Err; err != nil {
		return err
	}

	ah.started = false
	return nil
}

func (ah *AdHoc) GetConnectedPeers() ([]string, error) {
	var connectedPeers []string
	if variant, err := ah.obj.GetProperty(adhocPropertyConnectedPeers); err != nil {
		fmt.Printf("Failed to get connected peers property: %s\n", err)
		// log err
		connectedPeers = append(connectedPeers, "<nil>")
	} else {
		if err2 := variant.Store(connectedPeers); err2 != nil {
			fmt.Printf("Failed to store connected peers property: %s\n", err)
			// log err
			return nil, err2
		}
	}
	return connectedPeers, nil
}
