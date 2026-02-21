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
}

func NewVPCBridge() *VPCBridge {
	return &VPCBridge{
		DNSMap: map[string]string{
			"metadata.google.internal": "127.0.0.1",
			"storage.googleapis.com":   "127.0.0.1",
			"pubsub.googleapis.com":    "127.0.0.1",
		},
		Latency: 50 * time.Millisecond,
	}
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

// SimulatePath adds artificial latency to simulate cross-region or cross-zone traffic.
func (v *VPCBridge) SimulatePath() {
	if v.Latency > 0 {
		time.Sleep(v.Latency)
	}
}

func (v *VPCBridge) Dial(network, addr string) (net.Conn, error) {
	v.SimulatePath()
	host, port, _ := net.SplitHostPort(addr)
	resolved := v.Resolve(host)
	return net.Dial(network, fmt.Sprintf("%s:%s", resolved, port))
}
