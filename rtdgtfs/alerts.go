package rtdgtfs

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

func GetAlerts(url string) ([]*gtfs.Alert, error) {
	var allAlerts []*gtfs.Alert

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching GTFS Alerts feed failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RTD Alerts endpoint returned HTTP status: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading RTD Alerts response body failed: %w", err)
	}

	feed := &gtfs.FeedMessage{}
	err = proto.Unmarshal(body, feed)
	if err != nil {
		return nil, fmt.Errorf("parsing GTFS-RT Alerts protobuf failed: %w", err)
	}

	for _, entity := range feed.GetEntity() {
		if alert := entity.GetAlert(); alert != nil {
			allAlerts = append(allAlerts, alert)
		}
	}

	return allAlerts, nil
}
