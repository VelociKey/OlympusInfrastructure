package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	mode := os.Getenv("OLYMPUS_MODE")
	if mode == "" {
		mode = "hybrid" // Default: Podman for emulators, Native for Go Bridges
	}

	fmt.Printf("Starting Olympus Workstation Cloud in [%s] mode...\n", mode)

	// High-Fidelity Substrate: VPC Bridge
	bridge := NewVPCBridge()
	fmt.Printf("🌐 VPC Bridge Active (Latency: %v)\n", bridge.Latency)

	if mode == "hybrid" || mode == "sandbox" {
		startPodmanMesh()
	}

	if mode == "hybrid" || mode == "native" {
		startNativeManagers()
		startNativeBridges()
		checkHealth()
	}
	
	fmt.Println("🌌 Olympus Fleet is now OPERATIONAL. Press Ctrl+C to stop.")
	select {} // Block forever to keep subprocesses alive
}

func checkHealth() {
	fmt.Println("🏥 Running Fleet Health Checks...")
	ports := []string{"8092", "8091", "8096", "8098", "8095", "8093", "8094", "8097", "8099", "8090"}
	
	for _, port := range ports {
		online := false
		for i := 0; i < 5; i++ {
			conn, err := net.DialTimeout("tcp", "localhost:"+port, 1*time.Second)
			if err == nil {
				fmt.Printf("✅ Port %s: ONLINE\n", port)
				conn.Close()
				online = true
				break
			}
			time.Sleep(1 * time.Second)
		}
		if !online {
			fmt.Printf("❌ Port %s: OFFLINE\n", port)
		}
	}
}

func startPodmanMesh() {
	fmt.Println("🚀 Initializing Podman GCP Mesh Sandbox...")
	cmd := exec.Command("podman-compose", "-f", "Podman/podman-compose.yaml", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Podman not available or failed: %v\n", err)
	}
}

func startNativeManagers() {
	fmt.Println("🏗️  Scanning for Native Go-GCP Managers...")

	root := ".." // Assuming running from OlympusInfrastructure
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("Error reading root: %v\n", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "OlympusGCP-") {
			continue
		}

		cluster := entry.Name()
		managerDir := filepath.Join(root, cluster, "10000-Autonomous-Actors")
		if _, err := os.Stat(managerDir); err != nil {
			continue
		}

		subEntries, _ := os.ReadDir(managerDir)
		for _, sub := range subEntries {
			if !sub.IsDir() || !strings.HasSuffix(sub.Name(), "Manager") {
				continue
			}

			managerName := sub.Name()
			mainGo := filepath.Join(managerDir, managerName, "main.go")
			if _, err := os.Stat(mainGo); err == nil {
				fmt.Printf("🚀 Starting Manager: %s (%s)\n", managerName, cluster)

				cmd := exec.Command("go", "run", mainGo)
				cmd.Dir = filepath.Join(managerDir, managerName)
				if err := cmd.Start(); err != nil {
					fmt.Printf("Failed to start %s: %v\n", managerName, err)
				}
			}
		}
	}
	fmt.Println("Manager orchestration sequence initiated.")
}

func startNativeBridges() {
	fmt.Println("🏗️  Scanning for Native Go-MCP Bridges...")

	root := ".." // Assuming running from OlympusInfrastructure
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("Error reading root: %v\n", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "OlympusGCP-") {
			continue
		}

		cluster := entry.Name()
		bridgeDir := filepath.Join(root, cluster, "20000-Context-Bridges")
		if _, err := os.Stat(bridgeDir); err != nil {
			continue
		}

		subEntries, _ := os.ReadDir(bridgeDir)
		for _, sub := range subEntries {
			if !sub.IsDir() || !strings.HasSuffix(sub.Name(), "Bridge") {
				continue
			}

			bridgeName := sub.Name()
			mainGo := filepath.Join(bridgeDir, bridgeName, "main.go")
			if _, err := os.Stat(mainGo); err == nil {
				fmt.Printf("🚀 Starting Bridge: %s (%s)\n", bridgeName, cluster)

				cmd := exec.Command("go", "run", mainGo)
				cmd.Dir = filepath.Join(bridgeDir, bridgeName)
				if err := cmd.Start(); err != nil {
					fmt.Printf("Failed to start %s: %v\n", bridgeName, err)
				}
			}
		}
	}
	fmt.Println("Bridge orchestration sequence initiated.")
}
