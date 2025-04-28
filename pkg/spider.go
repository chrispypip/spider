package spider

import (
	"fmt"
	"sort"
	"time"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

type Spider struct {
	AccessPoints       []*AccessPoint
	Adapters           []*Adapter
	AdHocs             []*AdHoc
	AgentManager       *AgentManager
	BasicServiceSets   []*BasicServiceSet
	Devices            []*Device
	KnownNetworks      []*KnownNetwork
	Networks           []*Network
	Stations           []*Station
	StationDiagnostics []*StationDiagnostic
}

func GetAdapters(conn *dbus.Conn) ([]*Adapter, error) {
	adapters := make([]*Adapter, 0)
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	objectManager := conn.Object(IwdService, "/")
	if err := objectManager.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects); err != nil {
		log.Errorf("Failed to get managed objects: %s", err)
		return nil, fmt.Errorf("failed to get managed objects: %s", err)
	}
	for k, v := range objects {
		for resource := range v {
			switch resource {
			case "net.connman.iwd.Adapter":
				if adapter, err := NewAdapter(conn, k); err != nil {
					log.Errorf("failed to create Adapter from %s: %s", k, err)
					return nil, fmt.Errorf("failed to create Adapter from %s: %s", k, err)
				} else {
					adapters = append(adapters, adapter)
				}
			}
		}
	}
	log.Debugf("Adapters: %v", adapters)
	return adapters, nil
}

func GetAdapterByName(conn *dbus.Conn, name string) (*Adapter, error) {
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	objectManager := conn.Object(IwdService, "/")
	if err := objectManager.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects); err != nil {
		log.Errorf("Failed to get managed objects: %s", err)
		return nil, fmt.Errorf("failed to get managed objects: %s", err)
	}
	for k, v := range objects {
		for resource := range v {
			switch resource {
			case "net.connman.iwd.Adapter":
				if adapter, err := NewAdapter(conn, k); err != nil {
					log.Errorf("Failed to create Adapter from %s: %s", k, err)
					return nil, fmt.Errorf("failed to create Adapter from %s: %s", k, err)
				} else if adapter.GetName() == name {
					log.Debugf("Found Adapter with name %s: %s", name, *adapter)
					return adapter, nil
				}
			}
		}
	}
	log.Warnf("Could not find Device with name %s", name)
	return nil, fmt.Errorf("could not find Adapter with name %s", name)
}

func GetDevices(conn *dbus.Conn) ([]*Device, error) {
	devices := make([]*Device, 0)
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	objectManager := conn.Object(IwdService, "/")
	if err := objectManager.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects); err != nil {
		log.Errorf("Failed to get managed objects: %s", err)
		return nil, fmt.Errorf("failed to get managed objects: %s", err)
	}
	for k, v := range objects {
		for resource := range v {
			switch resource {
			case "net.connman.iwd.Device":
				if device, err := NewDevice(conn, k); err != nil {
					log.Errorf("Failed to create Device from %s: %s", k, err)
					return nil, fmt.Errorf("failed to create Device from %s: %s", k, err)
				} else {
					devices = append(devices, device)
				}
			}
		}
	}
	log.Debugf("Devices: %v", devices)
	return devices, nil
}

func GetDeviceByName(conn *dbus.Conn, name string) (*Device, error) {
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	objectManager := conn.Object(IwdService, "/")
	if err := objectManager.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects); err != nil {
		log.Errorf("Failed to get managed objects: %s", err)
		return nil, fmt.Errorf("failed to get managed objects: %s", err)
	}
	for k, v := range objects {
		for resource := range v {
			switch resource {
			case "net.connman.iwd.Device":
				if device, err := NewDevice(conn, k); err != nil {
					log.Errorf("Failed to create Device from %s: %s", k, err)
					return nil, fmt.Errorf("failed to create Device from %s: %s", k, err)
				} else if device.GetName() == name {
					log.Debugf("Found Device with name %s: %s", name, *device)
					return device, nil
				}
			}
		}
	}
	log.Warnf("Could not find Device with name %s", name)
	return nil, fmt.Errorf("could not find Device with name %s", name)
}

func ScanForNetworks(conn *dbus.Conn, scanTime uint8) ([]*StationOrderedNetwork, error) {
	stations := make([]*Station, 0)
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	objectManager := conn.Object(IwdService, "/")
	if err := objectManager.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects); err != nil {
		log.Errorf("Failed to get managed objects: %s", err)
		return nil, fmt.Errorf("failed to get managed objects: %s", err)
	}
	for k, v := range objects {
		for resource := range v {
			switch resource {
			case "net.connman.iwd.Station":
				if station, err := NewStation(conn, k); err != nil {
					log.Errorf("Failed to create Station from %s: %s", k, err)
					return nil, fmt.Errorf("failed to create Station from %s: %s", k, err)
				} else {
					stations = append(stations, station)
				}
			}
		}
	}
	orderedNetworks := make([]*StationOrderedNetwork, 0)
	if len(stations) == 0 {
		log.Error("No stations found")
		return nil, fmt.Errorf("no stations found; Adapters must be powered on and Devices must be powered on in 'station' mode")
	}
	for _, st := range stations {
		if err := st.Scan(); err != nil {
			log.Errorf("failed to scan Station %s", st.GetPath())
			return nil, fmt.Errorf("failed to scan Station %s", st.GetPath())
		} else {
			if scanTime > 0 {
				time.Sleep(time.Duration(scanTime) * time.Second)
			}
			if stOns, err2 := st.GetOrderedNetworks(); err2 != nil {
				log.Errorf("failed to get ordered networks for Station %s: %s", st.GetPath(), err2)
				return nil, fmt.Errorf("failed to get ordered networks for Station %s: %s", st.GetPath(), err2)
			} else {
				orderedNetworks = append(orderedNetworks, stOns...)
			}
		}
	}
	sort.SliceStable(orderedNetworks, func(i, j int) bool {
		name1, _ := orderedNetworks[i].Network.GetName()
		name2, _ := orderedNetworks[i].Network.GetName()
		return name1 < name2
	})
	sort.SliceStable(orderedNetworks, func(i, j int) bool { return orderedNetworks[i].SignalStrength > orderedNetworks[j].SignalStrength })
	sort.SliceStable(orderedNetworks, func(i, j int) bool { return orderedNetworks[i].Network.GetKnownNetwork() != nil })
	sort.SliceStable(orderedNetworks, func(i, j int) bool { return orderedNetworks[i].Network.GetConnected() })
	log.Debugf("OrderedNetworks: %v", orderedNetworks)
	return orderedNetworks, nil
}

func GetKnownNetworks(conn *dbus.Conn) ([]*KnownNetwork, error) {
	knownNetworks := make([]*KnownNetwork, 0)
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	objectManager := conn.Object(IwdService, "/")
	if err := objectManager.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects); err != nil {
		log.Errorf("failed to get managed objects: %s", err)
		return nil, fmt.Errorf("failed to get managed objects: %s", err)
	}
	for k, v := range objects {
		for resource := range v {
			switch resource {
			case "net.connman.iwd.KnownNetwork":
				if knownNetwork, err := NewKnownNetwork(conn, k); err != nil {
					log.Errorf("failed to create KnownNetwork from %s: %s", k, err)
					return nil, fmt.Errorf("failed to create KnownNetwork from %s: %s", k, err)
				} else {
					knownNetworks = append(knownNetworks, knownNetwork)
				}
			}
		}
	}
	log.Debugf("KnownNetworks: %v", knownNetworks)
	return knownNetworks, nil
}

func GetKnownNetworkByName(conn *dbus.Conn, name string) (*KnownNetwork, error) {
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	objectManager := conn.Object(IwdService, "/")
	if err := objectManager.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects); err != nil {
		log.Errorf("failed to get managed objects: %s", err)
		return nil, fmt.Errorf("failed to get managed objects: %s", err)
	}
	for k, v := range objects {
		for resource := range v {
			switch resource {
			case "net.connman.iwd.KnownNetwork":
				if knownNetwork, err := NewKnownNetwork(conn, k); err != nil {
					log.Errorf("failed to create KnownNetwork from %s: %s", k, err)
					return nil, fmt.Errorf("failed to create KnownNetwork from %s: %s", k, err)
				} else {
					knName, err2 := knownNetwork.GetName()
					if err2 != nil {
						log.Errorf("failed to get name of KnownNetwork %s: %s", k, err2)
						continue
					}
					if knName == name {
						log.Debugf("found KnownNetwork with name %s: %s", name, knownNetwork)
						return knownNetwork, nil
					}
				}
			}
		}
	}
	log.Warnf("Could not find KnownNetwork with name %s", name)
	return nil, fmt.Errorf("could not find KnownNetwork with name %s", name)
}

func GetAgentManager(conn *dbus.Conn) (*AgentManager, error) {
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	objectManager := conn.Object(IwdService, "/")
	if err := objectManager.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects); err != nil {
		log.Errorf("failed to get managed objects: %s", err)
		return nil, fmt.Errorf("failed to get managed objects: %s", err)
	}
	for k, v := range objects {
		for resource := range v {
			switch resource {
			case "net.connman.iwd.AgentManager":
				return NewAgentManager(conn, k), nil
			}
		}
	}
	return nil, fmt.Errorf("could not find an AgentManager")
}
