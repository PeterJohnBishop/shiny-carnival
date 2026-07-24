// Package rtdgtfs handles GTFS protobuf data provided by RTD
package rtdgtfs

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

func GetVehiclePositions(url string) ([]*gtfs.VehiclePosition, error) {
	var allVehiclePositions []*gtfs.VehiclePosition

	// fetch raw pb data
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching GTFS feed from RTD failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RTD endpoint returned unexpected HTTP status: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading RTD response body failed: %w", err)
	}

	// unmarshal the raw protobuf data into a FeedMessage struct
	feed := &gtfs.FeedMessage{}
	err = proto.Unmarshal(body, feed)
	if err != nil {
		return nil, fmt.Errorf("parsing GTFS-RT Protobuf payload failed: %w", err)
	}

	// check if feed returned zero entities
	if len(feed.GetEntity()) == 0 {
		return nil, fmt.Errorf("received empty GTFS feed (0 entities returned)")
	}

	fmt.Printf("Feed Timestamp: %v\n", time.Unix(int64(feed.GetHeader().GetTimestamp()), 0))
	fmt.Printf("Total Entities: %d\n\n", len(feed.GetEntity()))

	for _, entity := range feed.GetEntity() {
		if vehicle := entity.GetVehicle(); vehicle != nil {
			allVehiclePositions = append(allVehiclePositions, vehicle)
		}
	}

	return allVehiclePositions, nil
}
