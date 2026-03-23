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
		startCoreAgents()
		startNativeManagers()
		startNativeBridges()
		checkHealth()
	}

	fmt.Println("🌌 Olympus Fleet is now OPERATIONAL. Press Ctrl+C to stop.")
	select {} // Block forever to keep subprocesses alive
}

func checkHealth() {
	fmt.Println("🏥 Running Fleet Health Checks...")
	// Expanded port list based on genesis.json and Mesh Registry
	// 8080-8089: Core Agents
	// 8090-8099: GCP Emulators / Managers
	ports := []string{
		"8080", "8081", "8083", "8084", "8085", "8086", "8087", "8088", "8089",
		"8090", "8091", "8092", "8093", "8094", "8095", "8096", "8097", "8098", "8099",
	}

	for _, port := range ports {
		online := false
		for i := 0; i < 3; i++ {
			conn, err := net.DialTimeout("tcp", "localhost:"+port, 1*time.Second)
			if err == nil {
				fmt.Printf("✅ Port %s: ONLINE\n", port)
				conn.Close()
				online = true
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		if !online {
			fmt.Printf("❌ Port %s: OFFLINE\n", port)
		}
	}
}

func startCoreAgents() {
	fmt.Println("🏗️  Starting Core Olympus Agents (Orchestration Mesh)...")

	root := ".."
	agents := []struct {
		Name string
		Path string
	}{
		{"MeshHub", "olympus.fleet/00SDLC/OlympusActors-Delegation/10000-Autonomous-Actors/10500-Delegation-Management/10520-MeshHub"},
		{"Orchestrator", "olympus.fleet/00SDLC/OlympusActors-Delegation/10000-Autonomous-Actors/10500-Delegation-Management/10530-Orchestration-Control"},
		{"Coder", "olympus.fleet/00SDLC/OlympusActors-Cognition/10000-Autonomous-Actors/10600-Cognitive-Specialties/10620-Logic-Construction"},
		{"Architect", "olympus.fleet/00SDLC/OlympusActors-Cognition/10000-Autonomous-Actors/10600-Cognitive-Specialties/10610-Architectural-Synthesis"},
		{"MemoryAnchor", "olympus.fleet/00SDLC/Olympus2/60000-Information-Storage/610-Memory-Anchors/900-MemoryAnchor"},
		{"AegisGuardian", "olympus.fleet/00SDLC/Olympus2/01000-Identity-Foundations/010-Vision/900-AegisGuardian"},
		{"SovereignAudit", "olympus.fleet/00SDLC/Olympus2/80000-System-Governance/820-Sovereign-Audit/900-SovereignAudit"},
		{"KnowledgeHub", "olympus.fleet/00SDLC/Olympus2/30000-Federated-Services/310-Core-Registry/900-OlympusRegistry"},
		{"Forge", "olympus.fleet/00SDLC/OlympusForge/90000-Enablement-Labs/900-Forge"},
		{"Preflight", "olympus.fleet/00SDLC/OlympusFabric/70000-Environmental-Harness/730-Campaign-Execution"},
		{"MCPGateway", "olympus.fleet/00SDLC/OlympusMCP/10000-Autonomous-Actors/MCPGateway"},
	}

	for _, agent := range agents {
		mainGo := filepath.Join(root, agent.Path, "main.go")
		if _, err := os.Stat(mainGo); err == nil {
			fmt.Printf("🚀 Starting Agent: %s\n", agent.Name)

			cmd := exec.Command("go", "run", "main.go")
			cmd.Dir = filepath.Join(root, agent.Path)
			// Standardizing log output to Ephemeral Scratch
			logPath := filepath.Join(root, "olympus.fleet/00SDLC/Olympus2/C0500-Agent-Intelligence-Outputs/LPSV", strings.ToLower(agent.Name)+".log")
			logFile, _ := os.Create(logPath)
			cmd.Stdout = logFile
			cmd.Stderr = logFile

			if err := cmd.Start(); err != nil {
				fmt.Printf("Failed to start agent %s: %v\n", agent.Name, err)
			}
		}
	}
	fmt.Println("Core mesh orchestration sequence initiated.")
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
		bridgeDir := filepath.Join(root, cluster, "30000-Context-Bridges")
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
