package main

import "context"
import "dagger/olympusinfrastructure/internal/dagger"

type OlympusInfrastructure struct{}

func (m *OlympusInfrastructure) HelloWorld(ctx context.Context) string { return "Hello from OlympusInfrastructure!" }

func main() { dagger.Serve() }
