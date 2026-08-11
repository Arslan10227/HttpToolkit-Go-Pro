//go:build windows

package system

import (
	"golang.org/x/sys/windows/registry"
	"os"
)

// RegisterProtocol registers the custom URI scheme "httptoolkitpro" to the current executable path.
func RegisterProtocol() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Create registry keys under HKCU (HKEY_CURRENT_USER) to avoid requiring administrator privileges.
	// HKEY_CURRENT_USER\Software\Classes\httptoolkitpro
	k, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Classes\httptoolkitpro`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()

	if err := k.SetStringValue("", "URL:httptoolkitpro Protocol"); err != nil {
		return err
	}
	if err := k.SetStringValue("URL Protocol", ""); err != nil {
		return err
	}

	shellKey, _, err := registry.CreateKey(k, `shell\open\command`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer shellKey.Close()

	// e.g. "path\to\HttpToolkit-Pro.exe" "%1"
	cmdVal := `"` + execPath + `" "%1"`
	return shellKey.SetStringValue("", cmdVal)
}
