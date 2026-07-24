package main

import (
	"fmt"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/peterjohnbishop/shiny-carnival/rtdgtfs"
)

var (
	vehiclePositions []*gtfs.VehiclePosition
	tripUpdates      []*gtfs.TripUpdate
	alerts           []*gtfs.Alert
)

const (
	alertsFeedURL           = "https://open-data.rtd-denver.com/files/gtfs-rt/rtd/Alerts.pb"
	tripUpdatesFeedURL      = "https://open-data.rtd-denver.com/files/gtfs-rt/rtd/TripUpdate.pb"
	vehiclePositionsFeedURL = "https://open-data.rtd-denver.com/files/gtfs-rt/rtd/VehiclePosition.pb"
)

func main() {
	var err error
	vehiclePositions, err = rtdgtfs.GetVehiclePositions(vehiclePositionsFeedURL)
	if err != nil {
		fmt.Println(err)
	}
	tripUpdates, err = rtdgtfs.GetTripUpdates(tripUpdatesFeedURL)
	if err != nil {
		fmt.Println(err)
	}
	alerts, err = rtdgtfs.GetAlerts(alertsFeedURL)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(alerts)
}
