package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type KubeConfig struct {
	APIVersion     string                 `yaml:"apiVersion"`
	Kind           string                 `yaml:"kind"`
	Clusters       []map[string]any       `yaml:"clusters"`
	Contexts       []map[string]any       `yaml:"contexts"`
	CurrentContext string                 `yaml:"current-context"`
	Users          []map[string]any       `yaml:"users"`
	Extensions     map[string]interface{} `yaml:",inline"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <new_context_name> [new_server_ip]")
		os.Exit(1)
	}

	newContextName := os.Args[1]
	var newServerIP string
	if len(os.Args) == 3 {
		newServerIP = os.Args[2]
	}
	defaultKubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	// Read kubeconfig from stdin
	fmt.Println("Paste your kubeconfig file content below (Ctrl+D to end):")
	input := readFromStdin()

	var newKubeConfig KubeConfig
	if err := yaml.Unmarshal(input, &newKubeConfig); err != nil {
		fmt.Printf("Error parsing kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Rename the default context
	oldContextName := newKubeConfig.CurrentContext
	if oldContextName == "" {
		fmt.Println("The provided kubeconfig does not have a 'current-context'")
		os.Exit(1)
	}

	if !renameContext(&newKubeConfig, oldContextName, newContextName) {
		fmt.Printf("Failed to rename context '%s' to '%s'\n", oldContextName, newContextName)
		os.Exit(1)
	}

	// Merge with the default kubeconfig
	defaultKubeConfig, err := loadKubeConfig(defaultKubeconfigPath)
	if err != nil {
		fmt.Printf("Error loading default kubeconfig: %v\n", err)
		os.Exit(1)
	}

	mergeKubeConfigs(&defaultKubeConfig, &newKubeConfig, newServerIP)

	// Write the merged kubeconfig back to the default file
	if err := saveKubeConfig(defaultKubeconfigPath, &defaultKubeConfig); err != nil {
		fmt.Printf("Error saving kubeconfig: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Context '%s' added successfully to %s\n", newContextName, defaultKubeconfigPath)
}

func readFromStdin() []byte {
	var input bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil && err != io.EOF {
		fmt.Printf("Error reading from stdin: %v\n", err)
		os.Exit(1)
	}
	return input.Bytes()
}

func loadKubeConfig(path string) (KubeConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return KubeConfig{}, err
	}
	var config KubeConfig
	err = yaml.Unmarshal(data, &config)
	return config, err
}

func saveKubeConfig(path string, config *KubeConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func renameContext(config *KubeConfig, oldContext, newContext string) bool {
	for _, context := range config.Contexts {
		if context["name"] == oldContext {
			context["name"] = newContext
			config.CurrentContext = newContext
			return true
		}
	}
	return false
}

func mergeKubeConfigs(dest, src *KubeConfig, newServerIP string) {
	// Merge clusters
	for _, srcCluster := range src.Clusters {
		srcClusterName := srcCluster["name"].(string)
		found := false

		for _, destCluster := range dest.Clusters {
			if destCluster["name"] == srcClusterName {
				if newServerIP != "" {
					// Update server IP if newServerIP is provided
					destCluster["cluster"].(map[string]any)["server"] = newServerIP
				}
				found = true
				break
			}
		}

		if !found {
			dest.Clusters = append(dest.Clusters, srcCluster)
		}
	}

	// Merge users
	dest.Users = append(dest.Users, src.Users...)

	// Merge contexts
	dest.Contexts = append(dest.Contexts, src.Contexts...)

	// Update current context if it exists in the source
	if src.CurrentContext != "" {
		dest.CurrentContext = src.CurrentContext
	}
}
