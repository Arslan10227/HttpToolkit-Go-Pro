package interceptors

import (
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/webextension"
)

var installedWebExtension *webextension.Manager

// SetWebExtension configures the temp-installed Chromium extension path.
func SetWebExtension(m *webextension.Manager) {
	installedWebExtension = m
}

func chromiumExtensionPath(cfg *config.Config) string {
	if installedWebExtension != nil {
		if err := installedWebExtension.EnsureInstalled(); err == nil {
			return installedWebExtension.ExtensionPath()
		}
	}
	return assetPath(cfg, "webextension")
}
