package main

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRenameContext(t *testing.T) {
	kubeconfig := KubeConfig{
		Contexts: []map[string]any{
			{"name": "old-context", "context": map[string]any{"cluster": "test-cluster"}},
		},
		CurrentContext: "old-context",
	}

	// Rename context
	success := renameContext(&kubeconfig, "old-context", "new-context")
	if !success {
		t.Fatalf("renameContext failed to rename 'old-context' to 'new-context'")
	}

	// Validate result
	if kubeconfig.CurrentContext != "new-context" {
		t.Errorf("Expected current-context to be 'new-context', got '%s'", kubeconfig.CurrentContext)
	}

	found := false
	for _, context := range kubeconfig.Contexts {
		if context["name"] == "new-context" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected context 'new-context' not found in contexts")
	}
}

func TestMergeKubeConfigsWithIPAddressChange(t *testing.T) {
	dest := KubeConfig{
		Clusters: []map[string]any{
			{"name": "cluster1", "cluster": map[string]any{"server": "https://old-server-ip"}},
		},
		Contexts: []map[string]any{
			{"name": "context1", "context": map[string]any{"cluster": "cluster1"}},
		},
		Users: []map[string]any{
			{"name": "user1", "user": map[string]any{"token": "token1"}},
		},
	}

	src := KubeConfig{
		Clusters: []map[string]any{
			{"name": "cluster1", "cluster": map[string]any{"server": "https://new-server-ip"}},
			{"name": "cluster2", "cluster": map[string]any{"server": "https://cluster2.example.com"}},
		},
		Contexts: []map[string]any{
			{"name": "context2", "context": map[string]any{"cluster": "cluster2"}},
		},
		Users: []map[string]any{
			{"name": "user2", "user": map[string]any{"token": "token2"}},
		},
		CurrentContext: "context2",
	}

	// Merge configs with a new server IP for "cluster1"
	newServerIP := "https://updated-server-ip"
	mergeKubeConfigs(&dest, &src, newServerIP)

	// Validate clusters
	if len(dest.Clusters) != 2 {
		t.Errorf("Expected 2 clusters, got %d", len(dest.Clusters))
	}

	// Check that the server IP for "cluster1" was updated
	cluster1Found := false
	for _, cluster := range dest.Clusters {
		if cluster["name"] == "cluster1" {
			server := cluster["cluster"].(map[string]any)["server"].(string)
			if server != newServerIP {
				t.Errorf("Expected server IP for 'cluster1' to be '%s', got '%s'", newServerIP, server)
			}
			cluster1Found = true
		}
	}
	if !cluster1Found {
		t.Errorf("Cluster 'cluster1' not found")
	}

	// Check that "cluster2" was added correctly
	cluster2Found := false
	for _, cluster := range dest.Clusters {
		if cluster["name"] == "cluster2" {
			server := cluster["cluster"].(map[string]any)["server"].(string)
			if server != "https://cluster2.example.com" {
				t.Errorf("Expected server IP for 'cluster2' to be 'https://cluster2.example.com', got '%s'", server)
			}
			cluster2Found = true
		}
	}
	if !cluster2Found {
		t.Errorf("Cluster 'cluster2' not found")
	}

	// Validate contexts
	if len(dest.Contexts) != 2 {
		t.Errorf("Expected 2 contexts, got %d", len(dest.Contexts))
	}

	// Validate users
	if len(dest.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(dest.Users))
	}

	// Validate current context
	if dest.CurrentContext != "context2" {
		t.Errorf("Expected current-context to be 'context2', got '%s'", dest.CurrentContext)
	}
}

func TestSaveAndLoadKubeConfig(t *testing.T) {
	config := KubeConfig{
		Clusters: []map[string]any{
			{"name": "cluster1", "cluster": map[string]any{"server": "https://cluster1.example.com"}},
		},
		Contexts: []map[string]any{
			{"name": "context1", "context": map[string]any{"cluster": "cluster1"}},
		},
		Users: []map[string]any{
			{"name": "user1", "user": map[string]any{"token": "token1"}},
		},
		CurrentContext: "context1",
	}

	// Save config to a buffer
	var buf bytes.Buffer
	if err := yaml.NewEncoder(&buf).Encode(&config); err != nil {
		t.Fatalf("Failed to encode kubeconfig: %v", err)
	}

	// Load config from the buffer
	var loadedConfig KubeConfig
	if err := yaml.NewDecoder(&buf).Decode(&loadedConfig); err != nil {
		t.Fatalf("Failed to decode kubeconfig: %v", err)
	}

	// Validate loaded config
	if loadedConfig.CurrentContext != config.CurrentContext {
		t.Errorf("Expected current-context to be '%s', got '%s'", config.CurrentContext, loadedConfig.CurrentContext)
	}
	if len(loadedConfig.Clusters) != 1 {
		t.Errorf("Expected 1 cluster, got %d", len(loadedConfig.Clusters))
	}
	if len(loadedConfig.Users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(loadedConfig.Users))
	}
}
