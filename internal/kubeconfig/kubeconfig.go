package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultKubeConfigPath = "~/.kube/config"

type Config struct {
	APIVersion     string                 `json:"apiVersion"`
	Kind           string                 `json:"kind"`
	Preferences    map[string]interface{} `json:"preferences,omitempty"`
	CurrentContext string                 `json:"current-context"`
	Contexts       []Context              `json:"contexts"`
	Clusters       []Cluster              `json:"clusters"`
	Users          []User                 `json:"users"`
}

type Context struct {
	Name    string `json:"name"`
	Context struct {
		Cluster   string `json:"cluster"`
		User      string `json:"user"`
		Namespace string `json:"namespace,omitempty"`
	} `json:"context"`
}

type Cluster struct {
	Name    string `json:"name"`
	Cluster struct {
		Server                   string `json:"server"`
		CertificateAuthority     string `json:"certificate-authority,omitempty"`
		CertificateAuthorityData string `json:"certificate-authority-data,omitempty"`
	} `json:"cluster"`
}

type User struct {
	Name    string `json:"name"`
	User    struct {
		Username                string `json:"username,omitempty"`
		Password                string `json:"password,omitempty"`
		Token                   string `json:"token,omitempty"`
		ClientCertificate       string `json:"client-certificate,omitempty"`
		ClientCertificateData   string `json:"client-certificate-data,omitempty"`
		ClientKey               string `json:"client-key,omitempty"`
		ClientKeyData           string `json:"client-key-data,omitempty"`
	} `json:"user"`
}

// Simple YAML parser for kubeconfig
func parseKubeConfig(data []byte) (*Config, error) {
	var config Config

	// Simple parsing logic for kubeconfig files
	lines := strings.Split(string(data), "\n")

	contextStart := -1
	currentContext := ""

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "current-context:") {
			currentContext = strings.TrimSpace(strings.TrimPrefix(trimmed, "current-context:"))
			config.CurrentContext = currentContext
		} else if strings.HasPrefix(trimmed, "contexts:") {
			contextStart = i + 1
		} else if strings.HasPrefix(trimmed, "- name:") && contextStart != -1 {
			// Parse context
			contextName := strings.TrimSpace(strings.TrimPrefix(trimmed, "- name:"))
			ctx := Context{Name: contextName}

			// Look for context details in following lines
			for j := i + 1; j < len(lines); j++ {
				detailLine := strings.TrimSpace(lines[j])
				if strings.HasPrefix(detailLine, "- name:") {
					// Next context starts
					break
				}
				if strings.HasPrefix(detailLine, "cluster:") {
					ctx.Context.Cluster = strings.TrimSpace(strings.TrimPrefix(detailLine, "cluster:"))
				} else if strings.HasPrefix(detailLine, "user:") {
					ctx.Context.User = strings.TrimSpace(strings.TrimPrefix(detailLine, "user:"))
				} else if strings.HasPrefix(detailLine, "namespace:") {
					ctx.Context.Namespace = strings.TrimSpace(strings.TrimPrefix(detailLine, "namespace:"))
				}
			}

			config.Contexts = append(config.Contexts, ctx)
		}
	}

	// Set defaults
	if config.APIVersion == "" {
		config.APIVersion = "v1"
	}
	if config.Kind == "" {
		config.Kind = "Config"
	}

	return &config, nil
}

func serializeKubeConfig(config *Config) ([]byte, error) {
	var builder strings.Builder

	// Write header
	builder.WriteString("apiVersion: " + config.APIVersion + "\n")
	builder.WriteString("kind: " + config.Kind + "\n")
	builder.WriteString("current-context: " + config.CurrentContext + "\n")

	if len(config.Preferences) > 0 {
		builder.WriteString("preferences: {}\n")
	}

	// Write contexts
	builder.WriteString("contexts:\n")
	for _, ctx := range config.Contexts {
		builder.WriteString("  - name: " + ctx.Name + "\n")
		builder.WriteString("    context:\n")
		builder.WriteString("      cluster: " + ctx.Context.Cluster + "\n")
		builder.WriteString("      user: " + ctx.Context.User + "\n")
		if ctx.Context.Namespace != "" {
			builder.WriteString("      namespace: " + ctx.Context.Namespace + "\n")
		}
	}

	// Write clusters
	if len(config.Clusters) > 0 {
		builder.WriteString("clusters:\n")
		for _, cluster := range config.Clusters {
			builder.WriteString("  - name: " + cluster.Name + "\n")
			builder.WriteString("    cluster:\n")
			builder.WriteString("      server: " + cluster.Cluster.Server + "\n")
		}
	}

	// Write users
	if len(config.Users) > 0 {
		builder.WriteString("users:\n")
		for _, user := range config.Users {
			builder.WriteString("  - name: " + user.Name + "\n")
			builder.WriteString("    user: {}\n")
		}
	}

	return []byte(builder.String()), nil
}

func Load() (*Config, error) {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig file: %w", err)
	}

	config, err := parseKubeConfig(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubeconfig: %w", err)
	}

	return config, nil
}

func (c *Config) Save() error {
	configPath := getConfigPath()

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := serializeKubeConfig(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) GetContexts() []Context {
	return c.Contexts
}

func (c *Config) ContextExists(name string) bool {
	for _, ctx := range c.Contexts {
		if ctx.Name == name {
			return true
		}
	}
	return false
}

func (c *Config) SetCurrentContext(name string, namespace string) error {
	found := false
	for i, ctx := range c.Contexts {
		if ctx.Name == name {
			found = true
			c.Contexts[i].Context.Namespace = namespace
			break
		}
	}

	if !found {
		return fmt.Errorf("context '%s' not found", name)
	}

	c.CurrentContext = name
	return nil
}

func (c *Config) RemoveContext(name string) error {
	for i, ctx := range c.Contexts {
		if ctx.Name == name {
			c.Contexts = append(c.Contexts[:i], c.Contexts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("context '%s' not found", name)
}

func (c *Config) CurrentNamespace() string {
	for _, ctx := range c.Contexts {
		if ctx.Name == c.CurrentContext {
			return ctx.Context.Namespace
		}
	}
	return ""
}

func getConfigPath() string {
	if path := os.Getenv("KUBECONFIG"); path != "" {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return defaultKubeConfigPath
	}

	return filepath.Join(home, ".kube", "config")
}