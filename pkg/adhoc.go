package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	adhocInterface              = "net.connman.iwd.AdHoc"
	adhocPropertyConnectedPeers = adhocInterface + ".ConnectedPeers"
	adhocPropertyStarted        = adhocInterface + ".Started"
	adhocMethodStart            = adhocInterface + ".Start"
	adhocMethodStartOpen        = adhocInterface + ".StartOpen"
	adhocMethodStop             = adhocInterface + ".Stop"
)

var (
	adhocLogger *log.Entry
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
	conn    *dbus.Conn
	obj     dbus.BusObject
	path    dbus.ObjectPath
	started bool
}

func NewAdHoc(conn *dbus.Conn, path dbus.ObjectPath) (*AdHoc, error) {
	log.SetReportCaller(true)
	adhocLogger = log.WithFields(log.Fields{
		"type": "AdHoc",
		"path": path,
	})
	obj := conn.Object(IwdService, path)
	ah := &AdHoc{
		conn: conn,
		obj:  obj,
		path: path,
	}
	if variant, err := ah.obj.GetProperty(adhocPropertyStarted); err != nil {
		adhocLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'started'")
		return nil, err
	} else {
		if err2 := variant.Store(&ah.started); err2 != nil {
			adhocLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'started'")
			return nil, err2
		}
		adhocLogger.Debugf("Started = %t", ah.started)
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
		adhocLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to start AdHoc")
		return err
	}
	adhocLogger.Debugf("Started SSID %s", ssid)
	ah.started = true
	return nil
}

func (ah *AdHoc) StartOpen(ssid string) error {
	if err := ah.obj.Call(adhocMethodStartOpen, 0).Err; err != nil {
		adhocLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("Failed to start open AdHoc: SSID %s", ssid)
		return err
	}
	adhocLogger.Debugf("Start Open SSID %s", ssid)
	ah.started = true
	return nil
}

func (ah *AdHoc) Stop() error {
	if err := ah.obj.Call(adhocMethodStop, 0).Err; err != nil {
		adhocLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to stop AdHoc")
		return err
	}
	adhocLogger.Debug("Stopped")
	ah.started = false
	return nil
}

func (ah *AdHoc) GetConnectedPeers() ([]string, error) {
	var objects [][]dbus.Variant
	if variant, err := ah.obj.GetProperty(adhocPropertyConnectedPeers); err != nil {
		adhocLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'connected peers'")
		return nil, err
	} else {
		if err2 := variant.Store(&objects); err2 != nil {
			adhocLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'connected peers'")
			return nil, err2
		}
	}
	objLen := len(objects)
	connectedPeers := make([]string, 0, objLen)
	for i, obj := range objects {
		connectedPeers = append(connectedPeers, obj[0].Value().(string))
		adhocLogger.Debugf("Connected Peers %d = %s", i, connectedPeers[i])
	}
	return connectedPeers, nil
}
