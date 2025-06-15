package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	stationInterface                        = "net.connman.iwd.Station"
	stationPropertyState                    = stationInterface + ".State"
	stationPropertyConnectedNetwork         = stationInterface + ".ConnectedNetwork"
	stationPropertyScanning                 = stationInterface + ".Scanning"
	stationPropertyConnectedAccessPoint     = stationInterface + ".ConnectedAccessPoint"
	stationMethodScan                       = stationInterface + ".Scan"
	stationMethodDisconnect                 = stationInterface + ".Disconnect"
	stationMethodGetOrderedNetworks         = stationInterface + ".GetOrderedNetworks"
	stationMethodGetHiddenAccessPoints      = stationInterface + ".GetHiddenAccessPoints"
	stationMethodConnectHiddenNetwork       = stationInterface + ".ConnectHiddenNetwork"
	stationMethodRegisterSignalLevelAgent   = stationInterface + ".RegisterSignalLevelAgent"
	stationMethodUnregisterSignalLevelAgent = stationInterface + ".UnregisterSignalLevelAgent"
)

var (
	stationLogger *log.Entry
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
	Network        *Network
	SignalStrength int16
}

func (on *StationOrderedNetwork) String() string {
	return fmt.Sprintf("{Network: %s, SignalStength: %d}", on.Network, on.SignalStrength)
}

type HiddenAccessPoint struct {
	Address        string
	SignalStrength int16
	Type           string
}

type Station struct {
	conn *dbus.Conn
	obj  dbus.BusObject
	path dbus.ObjectPath
}

func NewStation(conn *dbus.Conn, path dbus.ObjectPath) (*Station, error) {
	log.SetReportCaller(true)
	stationLogger = log.WithFields(log.Fields{
		"type": "Station",
		"path": path,
	})
	obj := conn.Object(IwdService, path)
	s := &Station{
		conn: conn,
		obj:  obj,
		path: path,
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
		stationLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'state'")
		return "", err
	} else {
		if err2 := variant.Store(&state); err2 != nil {
			stationLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store state 'property'")
			return "", err2
		}
	}
	stationLogger.Debugf("State %s", state)
	return state, nil
}

func (s *Station) GetConnectedNetwork() (*Network, error) {
	if variant, err := s.obj.GetProperty(stationPropertyConnectedNetwork); err != nil {
		stationLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'connected network'")
		return nil, nil
	} else {
		var path dbus.ObjectPath
		if err2 := variant.Store(&path); err2 != nil {
			stationLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'connected network'")
			return nil, err2
		}
		if network, err2 := NewNetwork(s.conn, path); err2 != nil {
			stationLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to create new Network from property 'connected network'")
			return nil, err2
		} else {
			stationLogger.Debugf("Created new Network from property 'connected network': %s", network)
			return network, nil
		}
	}
}

func (s *Station) GetScanning() (bool, error) {
	var scanning bool
	if variant, err := s.obj.GetProperty(stationPropertyScanning); err != nil {
		stationLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get property 'scanning'")
		return false, err
	} else {
		if err2 := variant.Store(&scanning); err2 != nil {
			stationLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'scanning'")
			return false, err2
		}
	}
	stationLogger.Debugf("Scanning = %t", scanning)
	return scanning, nil
}

func (s *Station) GetConnectedAccessPoint() (*BasicServiceSet, error) {
	if variant, err := s.obj.GetProperty(stationPropertyConnectedAccessPoint); err != nil {
		stationLogger.WithFields(log.Fields{
			"err": err,
		}).Info("Failed to get optional property 'connected access point'")
		return nil, nil
	} else {
		var path dbus.ObjectPath
		if err2 := variant.Store(&path); err2 != nil {
			stationLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to store property 'connected access point'")
			return nil, err2
		}
		if bss, err2 := NewBasicServiceSet(s.conn, path); err2 != nil {
			stationLogger.WithFields(log.Fields{
				"err": err2,
			}).Error("Failed to create new Basic Service Set from property 'connected access point'")
			return nil, err2
		} else {
			stationLogger.Debugf("Created new Basic Service Set from property 'connected access point': %s", bss)
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
	var scanning bool
	var err error
	if scanning, err = s.GetScanning(); err != nil {
		stationLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to get station scanning state")
		return err
	}

	if !scanning {
		if err := s.obj.Call(stationMethodScan, 0).Err; err != nil {
			stationLogger.WithFields(log.Fields{
				"err": err,
			}).Error("Failed to start scan")
			return err
		}
		stationLogger.Debug("Scanning")
		return nil
	} else {
		stationLogger.Info("Station is already scanning")
		return nil
	}
}

func (s *Station) Disconnect() error {
	if err := s.obj.Call(stationMethodDisconnect, 0).Err; err != nil {
		stationLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to disconnect")
		return err
	}
	stationLogger.Debug("Disconnected")
	return nil
}

func (s *Station) GetOrderedNetworks() ([]*StationOrderedNetwork, error) {
	var objects [][]dbus.Variant
	if err := s.obj.Call(stationMethodGetOrderedNetworks, 0).Store(&objects); err != nil {
		stationLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to call GetOrderedNetworks")
		return nil, err
	}
	orderedNetworks := make([]*StationOrderedNetwork, 0, len(objects))
	for i, obj := range objects {
		networkPath := obj[0].Value().(dbus.ObjectPath)
		signalStrength := obj[1].Value().(int16)
		if network, err := NewNetwork(s.conn, networkPath); err != nil {
			stationLogger.WithFields(log.Fields{
				"err": err,
			}).Error("Failed to create new Network from GetOrderedNetworks")
			return orderedNetworks, nil
		} else {
			orderedNetworks = append(orderedNetworks, &StationOrderedNetwork{
				Network:        network,
				SignalStrength: signalStrength,
			})
			stationLogger.Debugf("Created new Network from GetOrderedNetworks: %s", orderedNetworks[i])
		}
		stationLogger.Debugf("Found %d Networks from GetOrderedNetworkds", i)
	}
	return orderedNetworks, nil
}

func (s *Station) GetHiddenAccessPoints() ([]*HiddenAccessPoint, error) {
	var objects [][]dbus.Variant
	if err := s.obj.Call(stationMethodGetHiddenAccessPoints, 0).Store(&objects); err != nil {
		stationLogger.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to call GetHiddenAccessPoints")
		return nil, err
	}
	accessPoints := make([]*HiddenAccessPoint, 0, len(objects))
	for i, obj := range objects {
		address := obj[0].Value().(string)
		signalStrength := obj[1].Value().(int16)
		accessPoints = append(accessPoints, &HiddenAccessPoint{
			Address:        address,
			SignalStrength: signalStrength,
		})
		stationLogger.Debugf("Created new Hidden Access point from GetHiddenAccessPoints: %v", *accessPoints[i])
	}
	return accessPoints, nil
}

func (s *Station) ConnectHiddenNetwork(ssid string) error {
	if err := s.obj.Call(stationMethodConnectHiddenNetwork, 0, ssid).Err; err != nil {
		stationLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to call ConnectHiddenNetwork for SSID %s", ssid)
		return err
	}
	stationLogger.Debugf("Connected to Hidden Network SSID %s", ssid)
	return nil
}

func (s *Station) RegisterSignalLevelAgent(client SignalLevelAgentClient, levels []int16) error {
	clientPath := client.GetPath()
	if err := s.obj.Call(stationMethodRegisterSignalLevelAgent, 0, clientPath, levels).Err; err != nil {
		stationLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to register SignalLevelAgent %s with level(s) %v", clientPath, levels)
		return err
	}
	stationLogger.Debugf("Registered SignalLevelAgent %s with level(s) %v", clientPath, levels)
	return nil
}

func (s *Station) UnregisterSignalLevelAgent(client SignalLevelAgentClient) error {
	clientPath := client.GetPath()
	if err := s.obj.Call(stationMethodUnregisterSignalLevelAgent, 0, clientPath).Err; err != nil {
		stationLogger.WithFields(log.Fields{
			"err": err,
		}).Errorf("failed to unregister SignalLevelAgent %s", clientPath)
		return err
	}
	stationLogger.Debugf("unregistered SignalLevelAgent %s", clientPath)
	return nil
}
