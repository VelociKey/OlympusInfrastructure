package main

import (
	"context"
	"fmt"
	"dagger/olympusinfrastructure/internal/dagger"
)

type OlympusInfrastructure struct{}

// BuildFirebaseCluster constructs the unified GCP/Firebase emulator cluster image for Podman.
func (m *OlympusInfrastructure) BuildFirebaseCluster(ctx context.Context, src *dagger.Directory) *dagger.Container {
	return dagger.
		Container().
		Build(src, dagger.ContainerBuildOpts{
			Dockerfile: "Podman/Firebase.Containerfile",
		})
}

// BuildNativeClusters compiles all Go-based GCP Managers and Bridges for the host environment.
func (m *OlympusInfrastructure) BuildNativeClusters(ctx context.Context, src *dagger.Directory) *dagger.Directory {
	// Start with a standard Go build environment
	builder := dagger.
		Container().
		From("golang:1.25-alpine").
		WithDirectory("/src", src).
		WithWorkdir("/src")

	// Define the clusters we want to build
	clusters := []string{
		"00SDLC/OlympusGCP-Compute",
		"00SDLC/OlympusGCP-Data",
		"00SDLC/OlympusGCP-Events",
		"00SDLC/OlympusGCP-FinOps",
		"00SDLC/OlympusGCP-Firebase",
		"00SDLC/OlympusGCP-Intelligence",
		"00SDLC/OlympusGCP-Messaging",
		"00SDLC/OlympusGCP-Observability",
		"00SDLC/OlympusGCP-Storage",
		"00SDLC/OlympusGCP-Vault",
	}

	// Output directory for binaries
	binDir := dagger.Directory()

	for _, cluster := range clusters {
		// Build Managers
		managerPath := fmt.Sprintf("%s/10000-Autonomous-Actors", cluster)
		// We'll use a shell script in the container to find and build any *Manager directories
		builder = builder.WithExec([]string{"sh", "-c", fmt.Sprintf(`
			cd %s && for d in *Manager; do
				if [ -d "$d" ]; then
					echo "Building Manager: $d"
					GOOS=windows GOARCH=amd64 go build -o /out/$d.exe ./$d
				fi
			done
		`, managerPath)})

		// Build Bridges
		bridgePath := fmt.Sprintf("%s/20000-Context-Bridges", cluster)
		builder = builder.WithExec([]string{"sh", "-c", fmt.Sprintf(`
			cd %s && for d in *Bridge; do
				if [ -d "$d" ]; then
					echo "Building Bridge: $d"
					GOOS=windows GOARCH=amd64 go build -o /out/$d.exe ./$d
				fi
			done
		`, bridgePath)})
	}

	return builder.Directory("/out")
}

func (m *OlympusInfrastructure) HelloWorld(ctx context.Context) string { return "Hello from OlympusInfrastructure!" }

func main() { dagger.Serve() }
