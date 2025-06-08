package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	accessPointInterface = "net.connman.iwd.AccessPoint"
	accessPointPropertyStarted = accessPointInterface + ".Started"
	accessPointPropertyName = accessPointInterface + ".Name"
	accessPointPropertyScanning = accessPointInterface + ".Scanning"
	accessPointPropertyFrequency = accessPointInterface + ".Frequency"
	accessPointPropertyPairwiseCiphers = accessPointInterface + ".PairwiseCiphers"
	accessPointPropertyGroupCipher = accessPointInterface + ".GroupCipher"
	accessPointMethodStart = accessPointInterface + ".Start"
	accessPointMethodStop = accessPointInterface + ".Stop"
	accessPointMethodStartProfile = accessPointInterface + ".StartProfile"
	accessPointMethodScan = accessPointInterface + ".Scan"
	accessPointMethodGetOrderedNetworks = accessPointInterface + ".GetOrderedNetworks"
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
	Name string
	SignalStrength int16
	Security string
}

type AccessPoint struct {
	conn *dbus.Conn
	obj dbus.BusObject
	path dbus.ObjectPath
	started bool
	name *string
	frequency *uint32
	pairwiseCiphers *[]string
	groupCipher *string
}

func NewAccessPoint(conn *dbus.Conn, path dbus.ObjectPath) (*AccessPoint, error) {
	log.SetReportCaller(true)
	obj := conn.Object(IwdService, path)
	ap := &AccessPoint{
		conn: conn,
		obj: obj,
		path: path,
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyStarted); err != nil {
		log.WithFields(log.Fields{
			"type": "AccessPoint",
			"path": path,
			"err": err,
		}).Error("Failed to get property 'started'")
		return nil, err
	} else {
		if err2 := variant.Store(&ap.started); err2 != nil {
			log.WithFields(log.Fields{
				"type": "AccessPoint",
				"path": path,
				"err": err2,
			}).Error("Failed to store proptery 'started'")
			return nil, err2
		}
		log.Debugf("AccessPoint: Started = %b", ap.started)
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyName); err != nil {
		log.WithFields(log.Fields{
			"type": "AccessPoint",
			"path": path,
			"err": err,
		}).Info("Failed to get optional property 'name'")
		ap.name = nil
	} else {
		if err2 := variant.Store(&ap.name); err2 != nil {
			log.WithFields(log.Fields{
				"type": "AccessPoint",
				"path": path,
				"err": err2,
			}).Errorf("Failed to store property 'name'")
			return nil, err2
		}
		log.Debugf("AccessPoint: Name = %s", *ap.name)
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyFrequency); err != nil {
		log.WithFields(log.Fields{
			"type": "AccessPoint",
			"path": path,
			"err": err,
		}).Info("Failed to get optional property 'frequency'")
		ap.frequency = nil
	} else {
		if err2 := variant.Store(&ap.frequency); err2 != nil {
			log.WithFields(log.Fields{
				"type": "AccessPoint",
				"path": path,
				"err": err,
			}).Error("Failed to store property 'frequency'")
			return nil, err2
		}
		log.Debugf("AccessPoint: Frequency = %d", *ap.frequency)
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyPairwiseCiphers); err != nil {
		log.WithFields(log.Fields{
			"type": "AccessPoint",
			"path": path,
			"err": err,
		}).Info("Failed to get optional property 'pairwise ciphers'")
		ap.pairwiseCiphers = nil
	} else {
		if err2 := variant.Store(&ap.pairwiseCiphers); err2 != nil {
			log.WithFields(log.Fields{
				"type": "AccessPoint",
				"path": path,
				"err": err,
			}).Error("Failed to store property 'pairwise cipher'")
			return nil, err2
		}
		log.Debugf("AccessPoint: Pairwise Ciphers = %v", *ap.pairwiseCiphers)
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyGroupCipher); err != nil {
		log.WithFields(log.Fields{
			"type": "AccessPoint",
			"path": path,
			"err": err,
		}).Info("Failed to get optional property 'group cipher'")
		ap.groupCipher = nil
	} else {
		if err2 := variant.Store(&ap.groupCipher); err2 != nil {
			log.WithFields(log.Fields{
				"type": "AccessPoint",
				"path": path,
				"err": err2,
			}).Error("Failed to store property 'group cipher'")
			return nil, err2
		}
		log.Debugf("AccessPoint: Group Cipher = %s", *ap.groupCipher)
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
	return fmt.Sprintf("{Path: %s, Interface: %s, Name: %s, Started: %t, Frequency: %d, PairwiseCiphers: %v, GroupCipher: %s, Scanning: %t}", ap.path, ap.GetInterface(), name, ap.started, frequency, pairwiseCiphers, groupCipher, scanning)
}

func (ap *AccessPoint) Start(ssid, psk string) error {
	if err := ap.obj.Call(accessPointMethodStart, 0, ssid, psk).Err; err != nil {
		log.WithFields(log.Fields{
			"type": "AccessPoint",
			"path": ap.GetPath(),
			"err": err,
		}).Error("Failed to start Access Point")
		return err
	}
	log.Debugf("AccessPoint: Started")
	ap.started = true
	return nil
}

func (ap *AccessPoint) Stop() error {
	if err := ap.obj.Call(accessPointMethodStop, 0).Err; err != nil {
		log.WithFields(log.Fields{
			"type": "AccessPoint",
			"path": ap.GetPath(),
			"err": err,
		}).Error("Failed to stop Access Point")
		return err
	}
	log.Debugf("AccessPoint: Stopped")
	ap.started = false
	return nil
}

func (ap *AccessPoint) StartProfile(ssid string) error {
	if err := ap.obj.Call(accessPointMethodStartProfile, 0).Err; err != nil {
		log.WithFields(log.Fields{
			"type": "AccessPoint",
			"path": ap.GetPath(),
			"err": err,
		}).Errorf("Failed to start Access Point profile: SSID %s", ssid)
		return err
	}
	log.Debugf("AccessPoint: Profile Started")
	ap.started = true
	return nil
}

func (ap *AccessPoint) Scan() error {
	if err := ap.obj.Call(accessPointMethodScan, 0).Err; err != nil {
		log.WithFields(log.Fields{
			"type": "AccessPoint",
			"path": ap.GetPath(),
			"err": err,
		}).Error("Failed to start scan")
		return err
	}
	log.Debugf("AccessPoint: Scanning")
	return nil
}

func (ap * AccessPoint) GetOrderedNetworks() ([]*AccessPointOrderedNetwork, error) {
	var objects [][]dbus.Variant
	if err := ap.obj.Call(accessPointMethodGetOrderedNetworks, 0).Store(&objects); err != nil {
		log.WithFields(log.Fields{
			"type": "AccessPoint",
			"path": ap.GetPath(),
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
			Name: name,
			SignalStrength: signalStrength,
			Security: security,
		})
		log.Debugf("AccessPoint: Ordered Network %d = %v", i, *orderedNetworks[i])
	}
	log.Debugf("AccessPoint: Found %d ordered networks", objLen)
	return orderedNetworks, nil
}

func (ap *AccessPoint) GetScanning() (bool, error) {
	var scanning bool
	apPath := ap.GetPath()
	if variant, err := ap.obj.GetProperty(accessPointPropertyScanning); err != nil {
		log.WithFields(log.Fields{
			"type": "AccessPoint",
			"path": apPath,
			"err": err,
		}).Error("Failed to get property 'scanning'")
		return false, err
	} else {
		if err2 := variant.Store(&scanning); err2 != nil {
			log.WithFields(log.Fields{
				"type": "AccessPoint",
				"path": apPath,
				"err": err2,
			}).Error("Failed to store property 'scanning'")
			return false, err2
		}
	}
	log.Debugf("AccessPoint: Scanning = %b", scanning)
	return scanning, nil
}
