package main

import (
	"fmt"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/peterjohnbishop/shiny-carnival/rtdgtfs"
)

var vehiclePositions []*gtfs.VehiclePosition

func main() {
	var err error
	vehiclePositions, err = rtdgtfs.GetVehiclePositions()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(vehiclePositions)
}
