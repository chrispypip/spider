package spider

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	log "github.com/sirupsen/logrus"
)

const (
	stationDiagnosticInterface  = "net.connman.iwd.StationDiagnostic"
	stationMethodGetDiagnostics = stationInterface + "GetDiagnostics"
)

var (
	stationDiagnosticLogger *log.Entry
)

type StationDiagnosticer interface {
	GetPath() dbus.ObjectPath
	GetInterface() string
	GetDiagnostics() error
}

type StationDiagnostic struct {
	conn *dbus.Conn
	obj  dbus.BusObject
	path dbus.ObjectPath
}

func NewStationDiagnostic(conn *dbus.Conn, path dbus.ObjectPath) (*StationDiagnostic, error) {
	log.SetReportCaller(true)
	stationDiagnosticLogger = log.WithFields(log.Fields{
		"type": "StationDiagnostic",
		"path": path,
	})
	obj := conn.Object(IwdService, path)
	sd := &StationDiagnostic{
		conn: conn,
		obj:  obj,
		path: path,
	}
	return sd, nil
}

func (sd *StationDiagnostic) String() string {
	return fmt.Sprintf("{Path: %s, Interface: %s}", sd.path, sd.GetInterface())
}

func (sd *StationDiagnostic) GetPath() dbus.ObjectPath {
	return sd.path
}

func (sd *StationDiagnostic) GetInterface() string {
	return stationDiagnosticInterface
}

func (sd *StationDiagnostic) GetDiagnostics() error {
	stationDiagnosticLogger.Debug("Getting diagnsotics")
	return nil
}
