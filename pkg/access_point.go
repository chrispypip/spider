package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	accessPointInterface                = "net.connman.iwd.AccessPoint"
	accessPointPropertyStarted          = accessPointInterface + ".Started"
	accessPointPropertyName             = accessPointInterface + ".Name"
	accessPointPropertyScanning         = accessPointInterface + ".Scanning"
	accessPointPropertyFrequency        = accessPointInterface + ".Frequency"
	accessPointPropertyPairwiseCiphers  = accessPointInterface + ".PairwiseCiphers"
	accessPointPropertyGroupCipher      = accessPointInterface + ".GroupCipher"
	accessPointMethodStart              = accessPointInterface + ".Start"
	accessPointMethodStop               = accessPointInterface + ".Stop"
	accessPointMethodStartProfile       = accessPointInterface + ".StartProfile"
	accessPointMethodScan               = accessPointInterface + ".Scan"
	accessPointMethodGetOrderedNetworks = accessPointInterface + ".GetOrderedNetworks"
)

var (
	apLogger *log.Entry
)

type AccessPointer interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	GetStarted() bool
	GetName() *string
	GetFrequency() *uint32
	GetPairwiseCiphers() *[]string
	GetGroupCipher() *string
	Start(ssid, psk string) error
	Stop() error
	StartProfile(ssid string) error
	Scan() error
	GetOrderedNetworks() ([]*AccessPointOrderedNetwork, error)
	GetScanning() (bool, error)
}

type AccessPointOrderedNetwork struct {
	Name           string
	SignalStrength int16
	Security       string
}

type AccessPoint struct {
	conn            *dbus.Conn
	obj             dbus.BusObject
	path            dbus.ObjectPath
	started         bool
	name            *string
	frequency       *uint32
	pairwiseCiphers *[]string
	groupCipher     *string
}

func NewAccessPoint(conn *dbus.Conn, path dbus.ObjectPath) (*AccessPoint, error) {
	log.SetReportCaller(true)
	apLogger = log.WithFields(log.Fields{
		"type": "AccessPoint",
		"path": path,
	})
	obj := conn.Object(IwdService, path)
	ap := &AccessPoint{
		conn: conn,
		obj:  obj,
		path: path,
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyStarted); err != nil {
		apLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'started'")
		return nil, err
	} else {
		if err2 := variant.Store(&ap.started); err2 != nil {
			apLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("failed to store property 'started'")
			return nil, err2
		}
		apLogger.Debugf("Started = %t", ap.started)
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyName); err != nil {
		apLogger.WithFields(log.Fields{
			"err": err,
		}).Info("Failed to get optional property 'name'")
		ap.name = nil
	} else {
		if err2 := variant.Store(&ap.name); err2 != nil {
			apLogger.WithFields(log.Fields{
				"err": err2,
			}).Errorf("failed to store property 'name'")
			return nil, err2
		}
		apLogger.Debugf("Name = %s", *ap.name)
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyFrequency); err != nil {
		apLogger.WithFields(log.Fields{
			"err": err,
		}).Info("Failed to get optional property 'frequency'")
		ap.frequency = nil
	} else {
		if err2 := variant.Store(&ap.frequency); err2 != nil {
			apLogger.WithFields(log.Fields{
				"err": err,
			}).Error("failed to store property 'frequency'")
			return nil, err2
		}
		apLogger.Debugf("Frequency = %d", *ap.frequency)
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyPairwiseCiphers); err != nil {
		apLogger.WithFields(log.Fields{
			"err": err,
		}).Info("Failed to get optional property 'pairwise ciphers'")
		ap.pairwiseCiphers = nil
	} else {
		if err2 := variant.Store(&ap.pairwiseCiphers); err2 != nil {
			apLogger.WithFields(log.Fields{
				"err": err,
			}).Error("failed to store property 'pairwise cipher'")
			return nil, err2
		}
		apLogger.Debugf("Pairwise Ciphers = %v", *ap.pairwiseCiphers)
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyGroupCipher); err != nil {
		apLogger.WithFields(log.Fields{
			"err": err,
		}).Info("Failed to get optional property 'group cipher'")
		ap.groupCipher = nil
	} else {
		if err2 := variant.Store(&ap.groupCipher); err2 != nil {
			apLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'group cipher'")
			return nil, err2
		}
		apLogger.Debugf("Group Cipher = %s", *ap.groupCipher)
	}
	return ap, nil
}

func (ap *AccessPoint) GetPath() dbus.ObjectPath {
	return ap.path
}

func (ap *AccessPoint) GetInterface() string {
	return accessPointInterface
}

func (ap *AccessPoint) GetStarted() bool {
	return ap.started
}

func (ap *AccessPoint) GetName() *string {
	return ap.name
}

func (ap *AccessPoint) GetFrequency() *uint32 {
	return ap.frequency
}

func (ap *AccessPoint) GetPairwiseCiphers() *[]string {
	return ap.pairwiseCiphers
}

func (ap *AccessPoint) GetGroupCipher() *string {
	return ap.groupCipher
}

func (ap *AccessPoint) String() string {
	var name string
	var frequency string
	var pairwiseCiphers string
	var groupCipher string
	if ap.name == nil {
		name = "<nil>"
	} else {
		name = *ap.name
	}
	if ap.frequency == nil {
		frequency = "<nil>"
	} else {
		frequency = fmt.Sprintf("%d", *ap.frequency)
	}
	if ap.pairwiseCiphers == nil {
		pairwiseCiphers = "<nil>"
	} else {
		pairwiseCiphers = fmt.Sprintf("%v", ap.pairwiseCiphers)
	}
	if ap.groupCipher == nil {
		groupCipher = "<nil>"
	} else {
		groupCipher = *ap.groupCipher
	}
	scanning, _ := ap.GetScanning()
	return fmt.Sprintf("{Path: %s, Interface: %s, Name: %s, Started: %t, Frequency: %s, PairwiseCiphers: %v, GroupCipher: %s, Scanning: %t}", ap.path, ap.GetInterface(), name, ap.started, frequency, pairwiseCiphers, groupCipher, scanning)
}

func (ap *AccessPoint) Start(ssid, psk string) error {
	if err := ap.obj.Call(accessPointMethodStart, 0, ssid, psk).Err; err != nil {
		apLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to start Access Point")
		return err
	}
	apLogger.Debugf("Started SSID = %s", ssid)
	ap.started = true
	return nil
}

func (ap *AccessPoint) Stop() error {
	if err := ap.obj.Call(accessPointMethodStop, 0).Err; err != nil {
		apLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to stop Access Point")
		return err
	}
	apLogger.Debug("Stopped")
	ap.started = false
	return nil
}

func (ap *AccessPoint) StartProfile(ssid string) error {
	if err := ap.obj.Call(accessPointMethodStartProfile, 0).Err; err != nil {
		apLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("Failed to start Access Point profile: SSID %s", ssid)
		return err
	}
	apLogger.Debugf("Profile Started SSID = %s", ssid)
	ap.started = true
	return nil
}

func (ap *AccessPoint) Scan() error {
	if err := ap.obj.Call(accessPointMethodScan, 0).Err; err != nil {
		apLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to start scan")
		return err
	}
	apLogger.Debug("Scanning")
	return nil
}

func (ap *AccessPoint) GetOrderedNetworks() ([]*AccessPointOrderedNetwork, error) {
	var objects [][]dbus.Variant
	if err := ap.obj.Call(accessPointMethodGetOrderedNetworks, 0).Store(&objects); err != nil {
		apLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get ordered networks")
		return nil, err
	}
	objLen := len(objects)
	orderedNetworks := make([]*AccessPointOrderedNetwork, 0, objLen)
	for i, obj := range objects {
		name := obj[0].Value().(string)
		signalStrength := obj[1].Value().(int16)
		security := obj[2].Value().(string)
		orderedNetworks = append(orderedNetworks, &AccessPointOrderedNetwork{
			Name:           name,
			SignalStrength: signalStrength,
			Security:       security,
		})
		apLogger.Debugf("Ordered Network %d = %v", i, *orderedNetworks[i])
	}
	apLogger.Debugf("Found %d ordered networks", objLen)
	return orderedNetworks, nil
}

func (ap *AccessPoint) GetScanning() (bool, error) {
	var scanning bool
	if variant, err := ap.obj.GetProperty(accessPointPropertyScanning); err != nil {
		apLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'scanning'")
		return false, err
	} else {
		if err2 := variant.Store(&scanning); err2 != nil {
			apLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'scanning'")
			return false, err2
		}
	}
	apLogger.Debugf("Scanning = %t", scanning)
	return scanning, nil
}
