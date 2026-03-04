package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"olympus.fleet/00SDLC/OlympusForge/70000-Environmental-Harness/dagger/olympusinfrastructure/internal/dagger"
)

type Olympusinfrastructure struct{}

// ExportGCPContainerClusters builds and exports production container images to a local directory.
func (m *Olympusinfrastructure) ExportGCPContainerClusters(ctx context.Context, src *dagger.Directory) *dagger.Directory {
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

	outDir := dag.Directory()

	for _, cluster := range clusters {
		serviceName := strings.ToLower(filepath.Base(cluster))
		
		binary := dag.Container().
			From("golang:1.25-alpine").
			WithDirectory("/src", src).
			WithWorkdir(fmt.Sprintf("/src/%s/10000-Autonomous-Actors", cluster)).
			WithExec([]string{"sh", "-c", "for d in *Manager; do if [ -d \"$d\" ]; then go build -o /out/manager ./$d; fi; done"}).
			File("/out/manager")

		container := dag.Container().
			From("alpine:latest").
			WithFile("/usr/local/bin/manager", binary).
			WithEntrypoint([]string{"/usr/local/bin/manager"}).
			WithLabel("org.velocikey.service", serviceName)

		tarFile := container.AsTarball()
		outDir = outDir.WithFile(fmt.Sprintf("%s.tar", serviceName), tarFile)
	}

	return outDir
}

// BuildNativeClusters compiles all Go-based GCP Managers and Bridges for the host environment.
func (m *Olympusinfrastructure) BuildNativeClusters(ctx context.Context, src *dagger.Directory) *dagger.Directory {
	builder := dag.
		Container().
		From("golang:1.25-alpine").
		WithDirectory("/src", src).
		WithWorkdir("/src")

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

	for _, cluster := range clusters {
		managerPath := fmt.Sprintf("%s/10000-Autonomous-Actors", cluster)
		builder = builder.WithExec([]string{"sh", "-c", fmt.Sprintf(`
			cd %s && for d in *Manager; do
				if [ -d "$d" ]; then
					echo "Building Manager: $d"
					GOOS=windows GOARCH=amd64 go build -o /out/$d.exe ./$d
				fi
			done
		`, managerPath)})

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

func (m *Olympusinfrastructure) HelloWorld(ctx context.Context) string { return "Hello from OlympusInfrastructure!" }
