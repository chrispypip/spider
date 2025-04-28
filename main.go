package main

import (
	"fmt"
	"os"

	"github.com/chrispypip/spider/pkg"
	"github.com/godbus/dbus/v5"
)

type SignalLevelAgent struct {
	conn *dbus.Conn
	path dbus.ObjectPath
}

func (a *SignalLevelAgent) GetConnection() *dbus.Conn {
	return a.conn
}

func (a *SignalLevelAgent) GetPath() dbus.ObjectPath {
	return a.path
}

func (a *SignalLevelAgent) Release(device dbus.ObjectPath) *dbus.Error {
	fmt.Printf("Released agent %s\n", device)
	return nil
}

func (a *SignalLevelAgent) Changed(device dbus.ObjectPath, level uint8) *dbus.Error {
	fmt.Printf("Signal level %s changed to %d\n", device, level)
	return nil
}

type Agent struct {
	conn *dbus.Conn
	path dbus.ObjectPath
}

func (a *Agent) GetConnection() *dbus.Conn {
	return a.conn
}

func (a *Agent) GetPath() dbus.ObjectPath {
	return a.path
}

func (a *Agent) Release() *dbus.Error {
	fmt.Printf("Releasing agent\n")
	return nil
}

func (a *Agent) RequestPassphrase(network dbus.ObjectPath) (string, *dbus.Error) {
	fmt.Printf("Requesting passhprase for %s\n", network)
	return "test123", nil
}

func (a *Agent) RequestPrivateKeyPassphrase(network dbus.ObjectPath) (string, *dbus.Error) {
	fmt.Printf("Requesting private key passphrase for %s\n", network)
	return "KeyPassword", nil
}

func (a *Agent) RequestUserNameAndPassword(network dbus.ObjectPath) (string, string, *dbus.Error) {
	fmt.Printf("Requesting user name and password for %s\n", network)
	return "User", "password", nil
}

func (a *Agent) RequestUserPassword(network dbus.ObjectPath, user string) (string, *dbus.Error) {
	fmt.Printf("Requesting password for for object %s for user %s\n", network, user)
	return "UserPassword", nil
}

func (a *Agent) Cancel(reason string) *dbus.Error {
	fmt.Printf("Canceling agent because: %s\n", reason)
	return nil
}

func main() {
	spider.SetLogLevel(spider.LogLevelDebug)
	_ = spider.SetLogFormatter(spider.LogJSONFormatter, false, true, true)
	if _, err := spider.AddLogFile("/root/mylog.txt", 0777); err != nil {
		fmt.Printf("Failed to add log file: %s\n", err)
		os.Exit(1)
	}
	spdr := &spider.Spider{}
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Printf("Failed to get system bus\n")
		os.Exit(1)
	}
	defer func() {
		_ = conn.Close()
	}()
	ad, err := spider.GetAdapterByName(conn, "phy0")
	if err != nil {
		fmt.Printf("failed to get adapter\n")
		os.Exit(1)
	}
	fmt.Print("\n\n")
	fmt.Printf("Found adapter phy0: %s\n", ad)
	_ = ad.SetPowered(true)

	var dev *spider.Device
	dev, err = spider.GetDeviceByName(conn, "wlan0")
	if err != nil {
		fmt.Printf("failed to get device wlan0\n")
		os.Exit(1)
	}
	fmt.Print("\n\n")
	fmt.Printf("Found device wlan0: %s\n", dev)
	_ = dev.SetPowered(true)
	_ = dev.SetMode("station")

	var nets []*spider.StationOrderedNetwork
	nets, err = spider.ScanForNetworks(conn, 5)
	if err != nil {
		fmt.Printf("Failed to scan for networks\n")
		os.Exit(1)
	}
	fmt.Print("\n\n")
	fmt.Printf("Found networks: %v\n", nets)
	//	sigAgent := &SignalLevelAgent{
	//		conn: conn,
	//		path:  "/test/sigagent/path",
	//	}
	//	err = conn.Export(sigAgent, sigAgent.GetPath(), "net.connman.iwd.SignalLevelAgent")
	//	if err != nil {
	//		fmt.Printf("Failed to export signal agent: %s\n", err)
	//		os.Exit(1)
	//	}
	//
	//	agent := spider.NewSimpleAgent(conn)
	//	err = conn.Export(agent, agent.GetPath(), "net.connman.iwd.Agent")
	//	if err != nil {
	//		fmt.Printf("Failed to export agent: %s\n", err)
	//		os.Exit(1)
	//	}
	//	agent.SetPassphrase("test123")
	//
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	objectManager := conn.Object(spider.IwdService, "/")
	if err = objectManager.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects); err != nil {
		fmt.Printf("Failed to get managed objects: %s\n", err)
		os.Exit(2)
	}
	for k, v := range objects {
		for resource := range v {
			switch resource {
			case "net.connman.iwd.AccessPoint":
				var accessPoint *spider.AccessPoint
				accessPoint, err = spider.NewAccessPoint(conn, k)
				if err != nil {
					fmt.Printf("Failed to created IWD access point: %s\n", err)
					os.Exit(3)
				}
				fmt.Printf("New access point:\n%s\n", accessPoint)
				spdr.AccessPoints = append(spdr.AccessPoints, accessPoint)
			case "net.connman.iwd.Adapter":
				var adapter *spider.Adapter
				adapter, err = spider.NewAdapter(conn, k)
				if err != nil {
					fmt.Printf("Failed to create IWD adapter: %s\n", err)
					os.Exit(3)
				}
				fmt.Printf("New adapter:\n%s\n", adapter)
				spdr.Adapters = append(spdr.Adapters, adapter)
			case "net.connman.iwd.AdHoc":
				var adhoc *spider.AdHoc
				adhoc, err = spider.NewAdHoc(conn, k)
				if err != nil {
					fmt.Printf("Failed to created IWD ad hoc: %s\n", err)
					os.Exit(3)
				}
				fmt.Printf("New adhoc:\n%s\n", adhoc)
				spdr.AdHocs = append(spdr.AdHocs, adhoc)
			case "net.connman.iwd.AgentManager":
				var am *spider.AgentManager
				am = spider.NewAgentManager(conn, k)
				fmt.Printf("New agent manager:\n%s\n", am)
				spdr.AgentManager = am
			case "net.connman.iwd.BasicServiceSet":
				var bss *spider.BasicServiceSet
				bss, err = spider.NewBasicServiceSet(conn, k)
				if err != nil {
					fmt.Printf("Failed to created IWD basic service set: %s\n", err)
					os.Exit(3)
				}
				fmt.Printf("New basic service set:\n%s\n", bss)
				spdr.BasicServiceSets = append(spdr.BasicServiceSets, bss)
			case "net.connman.iwd.Device":
				var device *spider.Device
				device, err = spider.NewDevice(conn, k)
				if err != nil {
					fmt.Printf("Failed to created IWD device: %s\n", err)
					os.Exit(3)
				}
				fmt.Printf("New device:\n%s\n", device)
				spdr.Devices = append(spdr.Devices, device)
			case "net.connman.iwd.KnownNetwork":
				var knownNetwork *spider.KnownNetwork
				knownNetwork, err = spider.NewKnownNetwork(conn, k)
				if err != nil {
					fmt.Printf("Failed to created IWD known network: %s\n", err)
					os.Exit(3)
				}
				fmt.Printf("New known network:\n%s\n", knownNetwork)
				spdr.KnownNetworks = append(spdr.KnownNetworks, knownNetwork)
			case "net.connman.iwd.Network":
				var network *spider.Network
				network, err = spider.NewNetwork(conn, k)
				if err != nil {
					fmt.Printf("Failed to created IWD network: %s\n", err)
					os.Exit(3)
				}
				fmt.Printf("New network:\n%s\n", network)
				spdr.Networks = append(spdr.Networks, network)
			case "net.connman.iwd.Station":
				var station *spider.Station
				station, err = spider.NewStation(conn, k)
				if err != nil {
					fmt.Printf("Failed to created IWD station: %s\n", err)
					os.Exit(3)
				}
				fmt.Printf("New station:\n%s\n", station)
				spdr.Stations = append(spdr.Stations, station)
			case "net.connman.iwd.StationDiagnostic":
				var stationDiagnostic *spider.StationDiagnostic
				stationDiagnostic, err = spider.NewStationDiagnostic(conn, k)
				if err != nil {
					fmt.Printf("Failed to create IWD station diagnostic: %s\n", err)
					os.Exit(3)
				}
				spdr.StationDiagnostics = append(spdr.StationDiagnostics, stationDiagnostic)
			}
		}
	}
	// fmt.Println("")
	// fmt.Println("")
	// fmt.Println("")
	// fmt.Println("")
	// fmt.Println("")
	// fmt.Println("")
	// orderedNets := make([]*spider.StationOrderedNetwork, 0)
	//
	if spdr.AgentManager == nil {
		fmt.Printf("No agent manager\n")
		os.Exit(10)
	}

	myAgentManager, err := spider.GetAgentManager(conn)
	if err != nil {
		fmt.Printf("No agent manager")
		os.Exit(10)
	}
	fmt.Printf("Found agent manager %s\n", myAgentManager)
	//
	//	if myerr := spdr.AgentManager.RegisterAgent(agent); myerr != nil {
	//		fmt.Printf("Failed to register agent: %s\n", myerr)
	//		os.Exit(127)
	//	}
	//
	// defer spdr.AgentManager.UnregisterAgent(agent)
	//
	//	if len(spdr.Stations) == 0 {
	//		fmt.Printf("No stations found\n")
	//		os.Exit(128)
	//	}
	//
	//	for _, st := range spdr.Stations {
	//		if myerr := st.Scan(); myerr != nil {
	//			fmt.Printf("Failed to scan station %s\n", st.GetPath())
	//		} else {
	//			if stOns, myerr2 := st.GetOrderedNetworks(); myerr2 != nil {
	//				fmt.Printf("Failed to get ordered netoworks for station %s: %s\n", st.GetPath(), myerr2)
	//			} else {
	//				orderedNets = append(orderedNets, stOns...)
	//			}
	//		}
	//	}
	//
	// fmt.Println("")
	// fmt.Println("")
	// fmt.Println("")
	// fmt.Printf("Ordered Networks:\n")
	// fmt.Printf("%v\n", orderedNets)
	// fmt.Println("")
	// fmt.Println("")
	// fmt.Println("")
	//
	// var mynet *spider.Network
	//
	//	for _, on := range orderedNets {
	//		if on.Network.GetName() == "TestNetwork" {
	//			mynet = on.Network
	//			break
	//		}
	//	}
	//
	//	if mynet == nil {
	//		fmt.Printf("Network TestNetwork not found\n")
	//		os.Exit(129)
	//	}
	//
	//	if myerr := mynet.Connect(); myerr != nil {
	//		fmt.Printf("Failed to connect to network: %s\n", myerr)
	//		os.Exit(130)
	//	}
	//
	// select {}
}
