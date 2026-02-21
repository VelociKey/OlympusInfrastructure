package main

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

// VPCBridge simulates internal GCP networking constraints.
type VPCBridge struct {
	DNSMap  map[string]string
	Latency time.Duration
	Chaos   map[string]float64 // Port -> Loss Probability (0.0 to 1.0)
}

func NewVPCBridge() *VPCBridge {
	return &VPCBridge{
		DNSMap: map[string]string{
			"metadata.google.internal": "127.0.0.1",
			"storage.googleapis.com":   "127.0.0.1",
			"pubsub.googleapis.com":    "127.0.0.1",
		},
		Latency: 50 * time.Millisecond,
		Chaos:   make(map[string]float64),
	}
}

// SetLoss sets the probability of a connection being dropped for a specific port.
func (v *VPCBridge) SetLoss(port string, probability float64) {
	v.Chaos[port] = probability
}

// Resolve simulates internal Cloud DNS.
func (v *VPCBridge) Resolve(hostname string) string {
	ip, ok := v.DNSMap[hostname]
	if !ok {
		return hostname
	}
	slog.Info("VPCBridge: Resolved internal DNS", "hostname", hostname, "ip", ip)
	return ip
}

// SimulatePath adds artificial latency and chaos.
func (v *VPCBridge) SimulatePath(port string) error {
	if loss, ok := v.Chaos[port]; ok {
		// Simple random loss simulation
		if time.Now().UnixNano()%100 < int64(loss*100) {
			slog.Warn("💥 VPCBridge: Simulated Packet Loss", "port", port)
			return fmt.Errorf("network connection dropped by vpc chaos engine")
		}
	}

	if v.Latency > 0 {
		time.Sleep(v.Latency)
	}
	return nil
}

func (v *VPCBridge) Dial(network, addr string) (net.Conn, error) {
	_, port, _ := net.SplitHostPort(addr)
	if err := v.SimulatePath(port); err != nil {
		return nil, err
	}
	
	host, _, _ := net.SplitHostPort(addr)
	resolved := v.Resolve(host)
	return net.Dial(network, fmt.Sprintf("%s:%s", resolved, port))
}
