package docker

import "strings"

// TransformComposeCreationLabels remaps compose config-hash for intercepted containers.
func TransformComposeCreationLabels(proxyPort int, labels map[string]string) map[string]string {
	if labels == nil {
		return nil
	}
	hash, ok := labels["com.docker.compose.config-hash"]
	if !ok {
		return labels
	}
	out := copyLabels(labels)
	out["com.docker.compose.config-hash"] = hash + interceptionSuffix(proxyPort)
	return out
}

// TransformComposeResponseLabels remaps compose hashes in API responses.
func TransformComposeResponseLabels(proxyPort int, labels map[string]string) map[string]string {
	if labels == nil {
		return nil
	}
	hash, ok := labels["com.docker.compose.config-hash"]
	if !ok {
		return labels
	}
	suffix := interceptionSuffix(proxyPort)
	out := copyLabels(labels)
	if strings.HasSuffix(hash, suffix) {
		out["com.docker.compose.config-hash"] = strings.TrimSuffix(hash, suffix)
	} else {
		out["com.docker.compose.config-hash"] = hash + "+unintercepted"
	}
	return out
}

func copyLabels(labels map[string]string) map[string]string {
	out := make(map[string]string, len(labels))
	for k, v := range labels {
		out[k] = v
	}
	return out
}
