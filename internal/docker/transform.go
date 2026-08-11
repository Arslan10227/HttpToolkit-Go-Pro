package docker

import (
	"encoding/json"
	"fmt"
	"strings"
)

type containerCreateConfig struct {
	Image      string            `json:"Image"`
	Env        []string          `json:"Env"`
	Labels     map[string]string `json:"Labels"`
	HostConfig struct {
		Binds []string `json:"Binds"`
	} `json:"HostConfig"`
}

// TransformContainerCreateConfig injects proxy env, CA volume, and interception labels.
func TransformContainerCreateConfig(raw []byte, settings InterceptionSettings) ([]byte, error) {
	var cfg containerCreateConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	if cfg.Image == "" {
		return raw, nil
	}

	host := DockerHostAddress("")
	proxyURL := fmt.Sprintf("http://%s:%d", host, settings.ProxyPort)
	caPath := injectedCAPath
	if settings.CertPath != "" {
		_ = ensureInjectionVolume(settings.ProxyPort, settings.CertPath)
	}

	injected := injectEnvVars(proxyURL, caPath)
	envKeys := make(map[string]struct{}, len(injected))
	for k := range injected {
		envKeys[k] = struct{}{}
	}
	var env []string
	for _, item := range cfg.Env {
		key := envKey(item)
		if _, skip := envKeys[key]; skip {
			continue
		}
		env = append(env, item)
	}
	for k, v := range injected {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cfg.Env = env

	if cfg.Labels == nil {
		cfg.Labels = map[string]string{}
	}
	cfg.Labels = TransformComposeCreationLabels(settings.ProxyPort, cfg.Labels)
	if cfg.Labels == nil {
		cfg.Labels = map[string]string{}
	}
	cfg.Labels[ContainerLabel] = fmt.Sprintf("%d", settings.ProxyPort)

	vol := InjectionVolumeName(settings.ProxyPort)
	bind := fmt.Sprintf("%s:%s:ro", vol, injectedPath)
	var binds []string
	for _, b := range cfg.HostConfig.Binds {
		if strings.Contains(b, injectedPath) {
			continue
		}
		binds = append(binds, b)
	}
	binds = append(binds, bind)
	cfg.HostConfig.Binds = binds

	return json.Marshal(cfg)
}

type containerInspectResponse struct {
	Config struct {
		Labels map[string]string `json:"Labels"`
	} `json:"Config"`
}

type containerListItem struct {
	Labels map[string]string `json:"Labels"`
}

func remapInspectResponse(proxyPort int, raw []byte) ([]byte, error) {
	var resp containerInspectResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return raw, err
	}
	resp.Config.Labels = TransformComposeResponseLabels(proxyPort, resp.Config.Labels)
	return json.Marshal(resp)
}

func remapListResponse(proxyPort int, raw []byte) ([]byte, error) {
	var list []containerListItem
	if err := json.Unmarshal(raw, &list); err != nil {
		return raw, err
	}
	for i := range list {
		list[i].Labels = TransformComposeResponseLabels(proxyPort, list[i].Labels)
	}
	return json.Marshal(list)
}
