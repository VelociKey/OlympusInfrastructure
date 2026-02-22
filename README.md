# OlympusInfrastructure

Central orchestration and sandbox management for **provisioning** the Olympus fleet.

## Sovereign Pod Strategy
All GCP emulators are deployed within a unified Pod in Podman Desktop. This provides:
1.  **Isolation**: Network and filesystem sandboxing for the customer workstation.
2.  **Consistency**: Identical network topology between local dev and customer delivery.
3.  **Hybrid Switching**: Developers can run Go bridges natively while talking to the containerized pod.
