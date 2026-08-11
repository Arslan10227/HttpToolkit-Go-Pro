package docker

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/system"
)

const ContainerLabel = "tech.httptoolkit.docker.proxy"

const injectedPath = "/.http-toolkit-injections"
const injectedCAPath = injectedPath + "/httptoolkit-ca.pem"

type InterceptionSettings struct {
	ProxyPort int
	CertPath  string
}

type containerInspect struct {
	ID     string `json:"Id"`
	Name   string `json:"Name"`
	Config struct {
		Env        []string          `json:"Env"`
		Cmd        []string          `json:"Cmd"`
		Entrypoint []string          `json:"Entrypoint"`
		Image      string            `json:"Image"`
		Labels     map[string]string `json:"Labels"`
		WorkingDir string            `json:"WorkingDir"`
		User       string            `json:"User"`
	} `json:"Config"`
	HostConfig struct {
		Binds []string `json:"Binds"`
	} `json:"HostConfig"`
	NetworkSettings struct {
		Networks map[string]any `json:"Networks"`
		Gateway  string         `json:"Gateway"`
	} `json:"NetworkSettings"`
}

func DockerHostAddress(gateway string) string {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		return "host.docker.internal"
	}
	if gateway != "" {
		return gateway
	}
	return "172.17.0.1"
}

func IsInterceptedContainer(labels map[string]string, proxyPort int) bool {
	if labels == nil {
		return false
	}
	return labels[ContainerLabel] == fmt.Sprintf("%d", proxyPort)
}

func InjectionVolumeName(proxyPort int) string {
	return fmt.Sprintf("httptoolkit-injection-%d", proxyPort)
}

func RestartAndInjectContainer(containerID string, settings InterceptionSettings) error {
	if containerID == "" {
		return fmt.Errorf("containerId required")
	}
	inspect, err := inspectContainer(containerID)
	if err != nil {
		return err
	}

	if err := stopContainer(containerID); err != nil {
		return err
	}
	if err := removeContainer(containerID); err != nil {
		return err
	}
	for i := 0; i < 50; i++ {
		if _, err := inspectContainer(containerID); err != nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	host := DockerHostAddress(inspect.NetworkSettings.Gateway)
	proxyURL := fmt.Sprintf("http://%s:%d", host, settings.ProxyPort)
	vol := InjectionVolumeName(settings.ProxyPort)
	caPath := injectedCAPath
	if settings.CertPath != "" {
		_ = ensureInjectionVolume(settings.ProxyPort, settings.CertPath)
	}

	args := []string{"create"}
	name := strings.TrimPrefix(inspect.Name, "/")
	if name != "" {
		args = append(args, "--name", name)
	}
	args = append(args, "--label", fmt.Sprintf("%s=%d", ContainerLabel, settings.ProxyPort))
	args = append(args, "--network", "none")

	for k, v := range injectEnvVars(proxyURL, caPath) {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}
	for _, env := range inspect.Config.Env {
		key := envKey(env)
		if key == "" || injectEnvVars(proxyURL, caPath)[key] != "" {
			continue
		}
		args = append(args, "-e", env)
	}
	for _, bind := range inspect.HostConfig.Binds {
		if strings.Contains(bind, injectedPath) {
			continue
		}
		args = append(args, "-v", bind)
	}
	args = append(args, "-v", fmt.Sprintf("%s:%s:ro", vol, injectedPath))

	if inspect.Config.WorkingDir != "" {
		args = append(args, "-w", inspect.Config.WorkingDir)
	}
	if inspect.Config.User != "" {
		args = append(args, "-u", inspect.Config.User)
	}

	if len(inspect.Config.Entrypoint) > 0 {
		args = append(args, "--entrypoint", inspect.Config.Entrypoint[0])
	}

	image := inspect.Config.Image
	if image == "" {
		image = inspect.ID
	}
	args = append(args, image)
	if len(inspect.Config.Cmd) > 0 {
		args = append(args, inspect.Config.Cmd...)
	}

	out, err := system.Command("docker", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker create: %s: %w", out, err)
	}
	newID := strings.TrimSpace(string(out))
	networks := networkNames(inspect.NetworkSettings.Networks)
	for _, net := range networks {
		if net == "none" {
			continue
		}
		_ = system.Command("docker", "network", "disconnect", "none", newID).Run()
		connectOut, connectErr := system.Command("docker", "network", "connect", net, newID).CombinedOutput()
		if connectErr != nil {
			return fmt.Errorf("docker network connect %s: %s: %w", net, connectOut, connectErr)
		}
	}

	startOut, startErr := system.Command("docker", "start", newID).CombinedOutput()
	if startErr != nil {
		return fmt.Errorf("docker start: %s: %w", startOut, startErr)
	}
	return nil
}

func injectEnvVars(proxyURL, certPath string) map[string]string {
	return map[string]string{
		"HTTP_PROXY":                 proxyURL,
		"HTTPS_PROXY":                proxyURL,
		"http_proxy":                 proxyURL,
		"https_proxy":                proxyURL,
		"WS_PROXY":                   proxyURL,
		"WSS_PROXY":                  proxyURL,
		"SSL_CERT_FILE":              certPath,
		"NODE_EXTRA_CA_CERTS":        certPath,
		"REQUESTS_CA_BUNDLE":         certPath,
		"CURL_CA_BUNDLE":             certPath,
		"AWS_CA_BUNDLE":              certPath,
		"HTTP_TOOLKIT_ACTIVE":        "true",
		"HTTP_TOOLKIT_OVERRIDE_PATH": injectedPath,
	}
}

func envKey(env string) string {
	if i := strings.Index(env, "="); i > 0 {
		return env[:i]
	}
	return ""
}

func inspectContainer(id string) (*containerInspect, error) {
	out, err := system.Command("docker", "inspect", id).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("docker inspect: %s: %w", out, err)
	}
	var list []containerInspect
	if err := json.Unmarshal(out, &list); err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("container not found")
	}
	return &list[0], nil
}

func stopContainer(id string) error {
	out, err := system.Command("docker", "stop", "-t", "1", id).CombinedOutput()
	if err != nil && !strings.Contains(string(out), "already stopped") {
		return fmt.Errorf("docker stop: %s: %w", out, err)
	}
	return nil
}

func removeContainer(id string) error {
	out, err := system.Command("docker", "remove", id).CombinedOutput()
	if err != nil && !strings.Contains(string(out), "No such container") {
		return fmt.Errorf("docker rm: %s: %w", out, err)
	}
	return nil
}

func networkNames(networks map[string]any) []string {
	if len(networks) == 0 {
		return nil
	}
	out := make([]string, 0, len(networks))
	for name := range networks {
		out = append(out, name)
	}
	return out
}

func ensureInjectionVolume(proxyPort int, certPath string) error {
	vol := InjectionVolumeName(proxyPort)
	_ = system.Command("docker", "volume", "create", vol).Run()
	run := system.Command("docker", "run", "--rm", "-v", vol+":/inject",
		"-v", certPath+":/ca.pem:ro", "alpine", "sh", "-c",
		"cp /ca.pem /inject/httptoolkit-ca.pem && chmod 644 /inject/httptoolkit-ca.pem")
	_, _ = run.CombinedOutput()
	return nil
}

func ListDetailedContainers() ([]map[string]any, error) {
	out, err := system.Command("docker", "ps", "-a", "--format", "{{json .}}").Output()
	if err != nil {
		return nil, err
	}
	var list []map[string]any
	for _, line := range splitLines(string(out)) {
		var row map[string]any
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			continue
		}
		id, _ := row["ID"].(string)
		if id == "" {
			continue
		}
		inspect, err := inspectContainer(id)
		if err != nil {
			continue
		}
		labels := inspect.Config.Labels
		if labels == nil {
			labels = map[string]string{}
		}
		ips := make([]string, 0)
		for _, netRaw := range inspect.NetworkSettings.Networks {
			if netMap, ok := netRaw.(map[string]any); ok {
				if ip, ok := netMap["IPAddress"].(string); ok && ip != "" {
					ips = append(ips, ip)
				}
			}
		}
		list = append(list, map[string]any{
			"id":     id,
			"names":  []string{strings.TrimPrefix(inspect.Name, "/")},
			"labels": labels,
			"state":  row["State"],
			"status": row["Status"],
			"image":  row["Image"],
			"ips":    ips,
		})
	}
	return list, nil
}

func DeleteInterceptedContainers(proxyPort int) error {
	out, err := system.Command("docker", "ps", "-aq", "--filter", fmt.Sprintf("label=%s=%d", ContainerLabel, proxyPort)).Output()
	if err != nil {
		return err
	}
	for _, id := range splitLines(string(out)) {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		_ = stopContainer(id)
		_ = removeContainer(id)
	}
	return deleteInterceptedBuildImages(proxyPort)
}

func deleteInterceptedBuildImages(proxyPort int) error {
	out, err := system.Command("docker", "images", "-q", "--filter", fmt.Sprintf("label=%s=%d", BuildLabel, proxyPort)).Output()
	if err != nil {
		return err
	}
	for _, id := range splitLines(string(out)) {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		_ = system.Command("docker", "rmi", "-f", id).Run()
	}
	return nil
}
