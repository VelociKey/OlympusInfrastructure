# Functional Requirements: MCP Service Testing

This document defines the high-level functional tests for the Olympus MCP bridges, driven by real-world developer complaints and usability requirements.

## 1. Security (Vault)
*   **Requirement**: "I want to test my app's secret loading without actually creating service accounts in GCP."
*   **Test Case (SC-01)**: Call `vault_read("api-key")`. Verify it returns the local jeBNF value and NOT a 401 Unauthorized from Google.
*   **Test Case (SC-02)**: Call `vault_write("new-secret", "val")`. Verify it persists locally and is available for subsequent reads.

## 2. Events (Async Debugging)
*   **Requirement**: "I need to know if my event was actually published before I spend hours debugging the consumer."
*   **Test Case (EV-01)**: Call `event_publish`. Verify the response contains a local `message_id`.
*   **Test Case (EV-02)**: Use `task_pause_queue` and then `task_create`. Verify the task is "Staged" and not executed until unpaused.

## 3. Storage (UI Velocity)
*   **Requirement**: "I need a local URL for my Flutter image assets that matches the GCS pattern."
*   **Test Case (ST-01)**: Call `storage_get_url`. Verify it returns a `file://` or `localhost://` URL that is reachable by the InteractionSurface.

## 4. Unified Data
*   **Requirement**: "I'm tired of rewriting my query logic when moving from Firestore to Datastore."
*   **Test Case (DT-01)**: Call `db_query` with a generic filter. Verify it returns valid JSON results regardless of which emulator is active.
