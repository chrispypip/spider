package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
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
	groupCipher *string }

func NewAccessPoint(conn *dbus.Conn, path dbus.ObjectPath) (*AccessPoint, error) {
	obj := conn.Object(IwdService, path)
	ap := &AccessPoint{
		conn: conn,
		obj: obj,
		path: path,
	}

	if variant, err := ap.obj.GetProperty(accessPointPropertyStarted); err != nil {
		fmt.Printf("Failed to get started property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&ap.started); err2 != nil {
			fmt.Printf("Failed to store started property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyName); err != nil {
		fmt.Printf("Failed to get name property: %s\n", err)
		// log err
		ap.name = nil
	} else {
		if err2 := variant.Store(&ap.name); err2 != nil {
			fmt.Printf("Failed to store name property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyFrequency); err != nil {
		fmt.Printf("Failed to get frequency property: %s\n", err)
		// log err
		ap.frequency = nil
	} else {
		if err2 := variant.Store(&ap.frequency); err2 != nil {
			fmt.Printf("Failed to store frequency property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyPairwiseCiphers); err != nil {
		fmt.Printf("Failed to get pairwise ciphers property: %s\n", err)
		// log err
		ap.pairwiseCiphers = nil
	} else {
		if err2 := variant.Store(&ap.pairwiseCiphers); err2 != nil {
			fmt.Printf("Failed to store pairwise ciphers property: %s\n", err2)
			// log err
			return nil, err2
		}
	}
	if variant, err := ap.obj.GetProperty(accessPointPropertyGroupCipher); err != nil {
		fmt.Printf("Failed to get group cipher property: %s\n", err)
		// log err
		ap.groupCipher = nil
	} else {
		if err2 := variant.Store(&ap.groupCipher); err2 != nil {
			fmt.Printf("Failed to store group cipher property: %s\n", err2)
			// log err
			return nil, err2
		}
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
		return err
	}
	
	ap.started = true
	return nil
}

func (ap *AccessPoint) Stop() error {
	if err := ap.obj.Call(accessPointMethodStop, 0).Err; err != nil {
		return err
	}

	ap.started = false
	return nil
}

func (ap *AccessPoint) StartProfile(ssid string) error {
	if err := ap.obj.Call(accessPointMethodStartProfile, 0).Err; err != nil {
		return err
	}

	ap.started = true
	return nil
}

func (ap *AccessPoint) Scan() error {
	return ap.obj.Call(accessPointMethodScan, 0).Err
}

func (ap * AccessPoint) GetOrderedNetworks() ([]*AccessPointOrderedNetwork, error) {
	var objects [][]dbus.Variant
	if err := ap.obj.Call(accessPointMethodGetOrderedNetworks, 0).Store(&objects); err != nil {
		return nil, err
	}

	orderedNetworks := make([]*AccessPointOrderedNetwork, 0, len(objects))
	for _, obj := range objects {
		name := obj[0].Value().(string)
		signalStrength := obj[1].Value().(int16)
		security := obj[2].Value().(string)
		orderedNetworks = append(orderedNetworks, &AccessPointOrderedNetwork{
			Name: name,
			SignalStrength: signalStrength,
			Security: security,
		})
	}
	
	return orderedNetworks, nil
}

func (ap *AccessPoint) GetScanning() (bool, error) {
	var scanning bool
	if variant, err := ap.obj.GetProperty(accessPointPropertyScanning); err != nil {
		fmt.Printf("Failed to get scanning property: %s\n", err)
		// log err
		return false, err
	} else {
		if err2 := variant.Store(&scanning); err2 != nil {
			fmt.Printf("Failed to store scanning property: %s\n", err)
			// log err
			return false, err2
		}
	}
	return scanning, nil
}
