package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	// 	"runtime"
)

func main() {
	mode := os.Getenv("OLYMPUS_MODE")
	if mode == "" {
		mode = "hybrid" // Default: Podman for emulators, Native for Go Bridges
	}

	fmt.Printf("Starting Olympus Workstation Cloud in [%s] mode...\n", mode)

	if mode == "hybrid" || mode == "sandbox" {
		startPodmanMesh()
	}

	if mode == "hybrid" || mode == "native" {
		startNativeBridges()
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

				// Run bridge in background
				cmd := exec.Command("go", "run", mainGo)
				cmd.Dir = filepath.Join(bridgeDir, bridgeName)
				// We don't want to block, so we don't call Wait()
				if err := cmd.Start(); err != nil {
					fmt.Printf("Failed to start %s: %v\n", bridgeName, err)
				}
			}
		}
	}
	fmt.Println("Bridge orchestration sequence initiated.")
}
