package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// TrackEntity represents a real-time tactical contact (Aircraft, Satellite, Naval)
type TrackEntity struct {
	ID        string    `json:"id"`
	Callsign  string    `json:"callsign"`
	Type      string    `json:"type"` // "UAV", "SATELLITE", "MARITIME"
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	AltitudeM float64   `json:"altitude_m"`
	SpeedKts  float64   `json:"speed_kts"`
	Heading   float64   `json:"heading"`
	Timestamp time.Time `json:"timestamp"`
}

// GeospatialStreamRelay coordinates concurrent telemetry pipelines via channels
type GeospatialStreamRelay struct {
	mu          sync.RWMutex
	subscribers map[string]chan TrackEntity
	broadcast   chan TrackEntity
}

func NewRelay() *GeospatialStreamRelay {
	return &GeospatialStreamRelay{
		subscribers: make(map[string]chan TrackEntity),
		broadcast:   make(chan TrackEntity, 1024),
	}
}

func (r *GeospatialStreamRelay) Subscribe(clientHost string) <-chan TrackEntity {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch := make(chan TrackEntity, 256)
	r.subscribers[clientHost] = ch
	return ch
}

func (r *GeospatialStreamRelay) Dispatch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case track := <-r.broadcast:
			r.mu.RLock()
			for _, ch := range r.subscribers {
				select {
				case ch <- track:
				default:
					// Drop packet if subscriber buffer is full (backpressure mitigation)
				}
			}
			r.mu.RUnlock()
		}
	}
}

func main() {
	fmt.Println("=== Tactical Geospatial Stream Relay (Go Engine Initialized) ===")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	relay := NewRelay()
	go relay.Dispatch(ctx)

	clientFeed := relay.Subscribe("dashboard_client_01")

	// Sample synthetic drone track stream
	go func() {
		testTarget := TrackEntity{
			ID:        "HEX-482B",
			Callsign:  "REAPER-01",
			Type:      "UAV",
			Lat:       34.0522,
			Lon:       -118.2437,
			AltitudeM: 6500.0,
			SpeedKts:  240.5,
			Heading:   182.4,
			Timestamp: time.Now().UTC(),
		}
		relay.broadcast <- testTarget
	}()

	select {
	case packet := <-clientFeed:
		rawJSON, _ := json.MarshalIndent(packet, "", "  ")
		fmt.Printf("[RELAY OUT] Ingested Track Broadcast:\n%s\n", string(rawJSON))
	case <-time.After(1 * time.Second):
		fmt.Println("[WARN] Pipeline timeout.")
	}
}
