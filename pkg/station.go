package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	stationInterface = "net.connman.iwd.Station"
	stationPropertyState = stationInterface + ".State"
	stationPropertyConnectedNetwork = stationInterface + ".ConnectedNetwork"
	stationPropertyScanning = stationInterface + ".Scanning"
	stationPropertyConnectedAccessPoint = stationInterface + ".ConnectedAccessPoint"
	stationMethodScan = stationInterface + ".Scan"
	stationMethodDisconnect = stationInterface + ".Disconnect"
	stationMethodGetOrderedNetworks = stationInterface + ".GetOrderedNetworks"
	stationMethodGetHiddenAccessPoints = stationInterface + ".GetHiddenAccessPoints"
	stationMethodConnectHiddenNetwork = stationInterface + ".ConnectHiddenNetwork"
	stationMethodRegisterSignalLevelAgent = stationInterface + ".RegisterSignalLevelAgent"
	stationMethodUnregisterSignalLevelAgent = stationInterface + ".UnregisterSignalLevelAgent"
)

type Stationer interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	GetState() (string, error)
	GetConnectedNetwork() (*Network, error)
	GetScanning() (bool, error)
	GetConnectedAccessPoint() (*BasicServiceSet, error)
	Scan() error
	Disconnect() error
 	GetOrderedNetworks() ([]StationOrderedNetwork, error)
	GetHiddenAccessPoints() ([]HiddenAccessPoint, error)
	ConnectHiddenNetwork(ssid string) error
	RegisterSignalLevelAgent(client SignalLevelAgentClient, levels []int16) error
	UnregisterSignalLevelAgent(client SignalLevelAgentClient) error
}

type StationOrderedNetwork struct {
	Network *Network
	SignalStrength int16
}

func (on *StationOrderedNetwork) String() string {
	return fmt.Sprintf("{Network: %s, SignalStength: %d}", on.Network, on.SignalStrength)
}

type HiddenAccessPoint struct {
	Address string
	SignalStrength int16
	Type string
}

type Station struct {
	conn *dbus.Conn
	obj dbus.BusObject
	path dbus.ObjectPath
	state string
}

func NewStation(conn *dbus.Conn, path dbus.ObjectPath) (*Station, error) {
	obj := conn.Object(IwdService, path)
	s := &Station{
		conn: conn,
		obj: obj,
		path: path,
	}

	if variant, err := s.obj.GetProperty(stationPropertyState); err != nil {
		fmt.Printf("Failed to get state property: %s\n", err)
		// log err
		return nil, err
	} else {
		if err2 := variant.Store(&s.state); err2 != nil {
			fmt.Printf("Failed to store state property: %s\n", err2)
			// log err
			return nil, err2
		}
	}

	return s, nil
}

func (s *Station) GetPath() dbus.ObjectPath {
	return s.path
}

func (s *Station) GetInterface() string {
	return stationInterface
}

func (s *Station) GetState() (string, error) {
	var state string
	if variant, err := s.obj.GetProperty(stationPropertyState); err != nil {
		fmt.Printf("Failed to get state property: %s\n", err)
		// log err
		return "", err
	} else {
		if err2 := variant.Store(&state); err2 != nil {
			fmt.Printf("Failed to store state property: %s\n", err)
			// log err
			return "", err2
		}
	}
	return state, nil
}

func (s *Station) GetConnectedNetwork() (*Network, error) {
	if variant, err := s.obj.GetProperty(stationPropertyConnectedNetwork); err != nil {
		fmt.Printf("Failed to get connected network property: %s\n", err)
		// log err
		return nil, nil
	} else {
		var path dbus.ObjectPath
		if err2 := variant.Store(&path); err2 != nil {
			fmt.Printf("Failed to store connected network path property: %s\n", err2)
			// log err
			return nil, err2
		}
		if network, err2 := NewNetwork(s.conn, path); err2 != nil {
			fmt.Printf("Failed to create new network for station: %s\n", err2)
			return nil, err2
		} else {
			return network, nil
		}
	}
}

func (s *Station) GetScanning() (bool, error) {
	var scanning bool
	if variant, err := s.obj.GetProperty(stationPropertyScanning); err != nil {
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

func (s *Station) GetConnectedAccessPoint() (*BasicServiceSet, error) {
	if variant, err := s.obj.GetProperty(stationPropertyConnectedAccessPoint); err != nil {
		fmt.Printf("Failed to get connected access point property: %s\n", err)
		// log err
		return nil, nil
	} else {
		var path dbus.ObjectPath
		if err2 := variant.Store(&path); err2 != nil {
			fmt.Printf("Failed to store connected access point path property: %s\n", err2)
			// log err
			return nil, err2
		}
		if bss, err2 := NewBasicServiceSet(s.conn, path); err2 != nil {
			fmt.Printf("Failed to create new basic service set for station: %s\n", err2)
			return nil, err2
		} else {
			return bss, nil
		}
	}
}

func (s *Station) String() string {
	var cn string
	var ap string
	state, _ := s.GetState()
	scanning, _ := s.GetScanning()
	connectedNetwork, _ := s.GetConnectedNetwork()
	connectedAccessPoint, _ := s.GetConnectedAccessPoint()
	if connectedNetwork == nil {
		cn = "<nil>"
	} else {
		cn = connectedNetwork.String()
	}
	if connectedAccessPoint == nil {
		ap = "<nil>"
	} else {
		ap = connectedAccessPoint.String()
	}
	return fmt.Sprintf("{Path: %s, Interface: %s, State: %s, ConnectedNetwork: %s, ConnectedAccessPoint: %s, Scanning: %t}", s.path, s.GetInterface(), state, cn, ap, scanning)
}

func (s *Station) Scan() error {
	if err := s.obj.Call(stationMethodScan, 0).Err; err != nil {
		return err
	}
	
	return nil
}

func (s *Station) Disconnect() error {
	if err := s.obj.Call(stationMethodDisconnect, 0).Err; err != nil {
		return err
	}

	return nil
}

func (s *Station) GetOrderedNetworks() ([]*StationOrderedNetwork, error) {
	var objects [][]dbus.Variant
	if err := s.obj.Call(stationMethodGetOrderedNetworks, 0).Store(&objects); err != nil {
		return nil, err
	}

	orderedNetworks := make([]*StationOrderedNetwork, 0, len(objects))
	for _, obj := range objects {
		networkPath := obj[0].Value().(dbus.ObjectPath)
		signalStrength := obj[1].Value().(int16)
		if network, err := NewNetwork(s.conn, networkPath); err != nil {
			fmt.Printf("Failed to created network for station ordered networks: %s\n", err)
			return orderedNetworks, nil
		} else {
			orderedNetworks = append(orderedNetworks, &StationOrderedNetwork{
				Network: network,
				SignalStrength: signalStrength,
			})
		}
	}
	
	return orderedNetworks, nil
}

func (s *Station) GetHiddenAccessPoints() ([]*HiddenAccessPoint, error) {
	var objects [][]dbus.Variant
	if err := s.obj.Call(stationMethodGetHiddenAccessPoints, 0).Store(&objects); err != nil {
		return nil, err
	}

	accessPoints := make([]*HiddenAccessPoint, 0, len(objects))
	for _, obj := range objects {
		address := obj[0].Value().(string)
		signalStrength := obj[1].Value().(int16)
		accessPoints = append(accessPoints, &HiddenAccessPoint{
			Address: address,
			SignalStrength: signalStrength,
		})
	}

	return accessPoints, nil
}

func (s *Station) ConnectHiddenNetwork(ssid string) error {
if err := s.obj.Call(stationMethodConnectHiddenNetwork, 0, ssid).Err; err != nil {
		return err
	}

	return nil
}

func (s *Station) RegisterSignalLevelAgent(client SignalLevelAgentClient, levels []int16) error {
	if err := s.obj.Call(stationMethodRegisterSignalLevelAgent, 0, client.GetPath(), levels).Err; err != nil {
		return err
	}

	return nil
}

func (s *Station) UnregisterSignalLevelAgent(client SignalLevelAgentClient) error {
	if err := s.obj.Call(stationMethodUnregisterSignalLevelAgent, 0, client.GetPath()).Err; err != nil {
		return err
	}

	return nil
}
