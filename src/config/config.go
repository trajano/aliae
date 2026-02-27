package config

import (
	"bytes"
	context_ "context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/jandedobbeleer/aliae/src/context"
	"github.com/jandedobbeleer/aliae/src/shell"
	aliaeState "github.com/jandedobbeleer/aliae/src/state"
	yamlv3 "gopkg.in/yaml.v3"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	defaultTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}
	client httpClient = &http.Client{Transport: defaultTransport}

	configPathCache string
)

func LoadConfig(configPath string) (*Aliae, error) {
	configPathCache = resolveConfigPath(configPath)
	setTemplateConfigContext(configPathCache)
	rootProgress, _ := loadRootProgress(configPathCache)

	if strings.HasPrefix(configPathCache, "http://") || strings.HasPrefix(configPathCache, "https://") {
		aliae, err := getRemoteConfig(configPathCache)
		if err != nil {
			return nil, err
		}
		if err := validateScriptStateRuntime(aliae); err != nil {
			return nil, err
		}

		// progress.internal is root-only and must ignore included/extended sources.
		aliae.Progress.Internal = rootProgress.Internal
		return aliae, nil
	}

	aliae, err := loadLocalConfig(configPathCache)
	if err != nil {
		return nil, err
	}
	if err := aliae.computeVars(nil); err != nil {
		return nil, err
	}
	if err := validateScriptStateRuntime(aliae); err != nil {
		return nil, err
	}

	// progress.internal is root-only and must ignore included/extended sources.
	aliae.Progress.Internal = rootProgress.Internal
	shell.SetStatTimeout(aliae.StatTimeout)

	return aliae, nil
}

func setTemplateConfigContext(configPath string) {
	if context.Current == nil {
		return
	}

	context.Current.ConfigPath = configPath
	context.Current.ConfigDir = resolveConfigDir(configPath)
}

func resolveConfigDir(configPath string) string {
	if !strings.HasPrefix(configPath, "http://") && !strings.HasPrefix(configPath, "https://") {
		return filepath.Dir(configPath)
	}

	parsed, err := url.Parse(configPath)
	if err != nil {
		return filepath.Dir(configPath)
	}

	parsed.Path = path.Dir(parsed.Path)
	parsed.RawQuery = ""
	parsed.Fragment = ""

	return parsed.String()
}

func home() string {
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return home
	}

	// fallback to older implemenations on Windows
	home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if len(home) == 0 {
		home = os.Getenv("USERPROFILE")
	}

	return home
}

func resolveConfigPath(configPath string) string {
	if len(configPath) == 0 {
		configPath = os.Getenv("ALIAE_CONFIG")
	}

	if len(configPath) == 0 {
		configPath = path.Join(home(), ".aliae.yaml")
	}

	return replaceTildePrefixWithHomeDir(configPath)
}

func ResolveTemplateContext(configPath string) (string, string) {
	resolvedPath := resolveConfigPath(configPath)
	return resolvedPath, resolveConfigDir(resolvedPath)
}

func replaceTildePrefixWithHomeDir(dir string) string {
	if !strings.HasPrefix(dir, "~") {
		return dir
	}

	rem := dir[1:]
	if len(rem) == 0 || isSeparator(rem[0]) {
		return home() + rem
	}

	return dir
}

func isSeparator(c uint8) bool {
	if c == '/' {
		return true
	}

	if runtime.GOOS == context.WINDOWS && c == '\\' {
		return true
	}

	return false
}

func getRemoteConfig(configURL string) (*Aliae, error) {
	req, err := http.NewRequestWithContext(context_.Background(), "GET", configURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download config file: %s\n→ %s", configURL, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseConfig(data)
}

func parseConfig(data []byte) (*Aliae, error) {
	if err := validateScriptWeightsInYAML(data); err != nil {
		return nil, err
	}
	if err := validateScriptStateInYAML(data); err != nil {
		return nil, err
	}

	var aliae Aliae

	decoder := yaml.NewDecoder(bytes.NewBuffer(data), yaml.CustomUnmarshaler(aliaeUnmarshaler))
	err := decoder.Decode(&aliae)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %s", err)
	}

	shell.SetStatTimeout(aliae.StatTimeout)
	if err := aliae.computeVars(nil); err != nil {
		return nil, err
	}

	return &aliae, nil
}

func validateScriptWeightsInYAML(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var root yamlv3.Node
	if err := yamlv3.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("failed to parse config file: %s", err)
	}

	if len(root.Content) == 0 {
		return nil
	}

	doc := root.Content[0]
	if doc == nil || doc.Kind != yamlv3.MappingNode {
		return nil
	}

	scriptNode := findMappingValue(doc, "script")
	if scriptNode == nil || scriptNode.Kind != yamlv3.SequenceNode {
		return nil
	}

	for i, item := range scriptNode.Content {
		if item == nil || item.Kind != yamlv3.MappingNode {
			continue
		}

		weightNode := findMappingValue(item, "weight")
		if weightNode == nil {
			continue
		}

		weight, err := strconv.ParseFloat(strings.TrimSpace(weightNode.Value), 64)
		if err != nil {
			return fmt.Errorf("invalid script[%d].weight: %q", i, weightNode.Value)
		}

		if weight <= 0 {
			return fmt.Errorf("invalid script[%d].weight: must be greater than 0", i)
		}
	}

	return nil
}

func validateScriptStateInYAML(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	var root yamlv3.Node
	if err := yamlv3.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("failed to parse config file: %s", err)
	}

	if len(root.Content) == 0 {
		return nil
	}

	doc := root.Content[0]
	if doc == nil || doc.Kind != yamlv3.MappingNode {
		return nil
	}

	scriptNode := findMappingValue(doc, "script")
	if scriptNode == nil || scriptNode.Kind != yamlv3.SequenceNode {
		return nil
	}

	seenStateFiles := make(map[string]int)
	for i, item := range scriptNode.Content {
		if item == nil || item.Kind != yamlv3.MappingNode {
			continue
		}

		stateNode := findMappingValue(item, "state")
		if stateNode == nil || stateNode.Kind != yamlv3.MappingNode {
			continue
		}

		fileNode := findMappingValue(stateNode, "file")
		stateFile := ""
		if fileNode != nil {
			stateFile = strings.TrimSpace(fileNode.Value)
		}

		if !aliaeState.IsValidFileName(stateFile) {
			return fmt.Errorf("invalid script[%d].state.file: %q", i, stateFile)
		}

		if previousIndex, exists := seenStateFiles[stateFile]; exists {
			return fmt.Errorf("invalid script[%d].state.file: duplicates script[%d].state.file", i, previousIndex)
		}
		seenStateFiles[stateFile] = i

		runEveryNode := findMappingValue(stateNode, "runEvery")
		if runEveryNode == nil || len(strings.TrimSpace(runEveryNode.Value)) == 0 {
			continue
		}

		runEvery, err := time.ParseDuration(strings.TrimSpace(runEveryNode.Value))
		if err != nil {
			return fmt.Errorf("invalid script[%d].state.runEvery: %q", i, runEveryNode.Value)
		}
		if runEvery <= 0 {
			return fmt.Errorf("invalid script[%d].state.runEvery: must be greater than 0", i)
		}
	}

	return nil
}

func findMappingValue(node *yamlv3.Node, key string) *yamlv3.Node {
	if node == nil || node.Kind != yamlv3.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}

	return nil
}

func validateScriptStateRuntime(aliae *Aliae) error {
	if aliae == nil {
		return nil
	}

	seenStateFiles := make(map[string]int)
	for i, script := range aliae.Scripts {
		if script == nil || len(script.State.File) == 0 {
			continue
		}

		stateFile := strings.TrimSpace(string(script.State.File))
		if !aliaeState.IsValidFileName(stateFile) {
			return fmt.Errorf("invalid script[%d].state.file: %q", i, stateFile)
		}

		if previousIndex, exists := seenStateFiles[stateFile]; exists {
			return fmt.Errorf("invalid script[%d].state.file: duplicates script[%d].state.file", i, previousIndex)
		}
		seenStateFiles[stateFile] = i

		if len(strings.TrimSpace(script.State.RunEvery)) == 0 {
			continue
		}

		runEvery, err := time.ParseDuration(strings.TrimSpace(script.State.RunEvery))
		if err != nil {
			return fmt.Errorf("invalid script[%d].state.runEvery: %q", i, script.State.RunEvery)
		}
		if runEvery <= 0 {
			return fmt.Errorf("invalid script[%d].state.runEvery: must be greater than 0", i)
		}
	}

	return nil
}
