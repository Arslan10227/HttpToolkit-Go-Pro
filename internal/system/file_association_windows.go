//go:build windows

package system

import (
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"strings"
)

// IsHarAssociated checks if the .har extension is associated with this executable.
func IsHarAssociated() bool {
	execPath, err := os.Executable()
	if err != nil {
		return false
	}

	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\.har`, registry.READ)
	if err != nil {
		return false
	}
	defer k.Close()

	progID, _, err := k.GetStringValue("")
	if err != nil || progID != "httptoolkitpro.har" {
		return false
	}

	cmdKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\httptoolkitpro.har\shell\open\command`, registry.READ)
	if err != nil {
		return false
	}
	defer cmdKey.Close()

	cmdVal, _, err := cmdKey.GetStringValue("")
	if err != nil {
		return false
	}

	expectedCmd := `"` + execPath + `" "%1"`
	return strings.EqualFold(cmdVal, expectedCmd)
}

// RegisterHarAssociation registers the .har extension to open with this executable.
func RegisterHarAssociation() bool {
	execPath, err := os.Executable()
	if err != nil {
		return false
	}

	k, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Classes\.har`, registry.ALL_ACCESS)
	if err != nil {
		return false
	}
	defer k.Close()
	if err := k.SetStringValue("", "httptoolkitpro.har"); err != nil {
		return false
	}

	progKey, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Classes\httptoolkitpro.har`, registry.ALL_ACCESS)
	if err != nil {
		return false
	}
	defer progKey.Close()
	if err := progKey.SetStringValue("", "HTTP Archive File"); err != nil {
		return false
	}

	iconKey, _, err := registry.CreateKey(progKey, "DefaultIcon", registry.ALL_ACCESS)
	if err != nil {
		return false
	}
	defer iconKey.Close()
	harIconPath := filepath.Join(filepath.Dir(execPath), "har_logo.ico")
	if _, err := os.Stat(harIconPath); err == nil {
		if err := iconKey.SetStringValue("", `"`+harIconPath+`"`); err != nil {
			return false
		}
	} else {
		if err := iconKey.SetStringValue("", `"`+execPath+`",0`); err != nil {
			return false
		}
	}

	cmdKey, _, err := registry.CreateKey(progKey, `shell\open\command`, registry.ALL_ACCESS)
	if err != nil {
		return false
	}
	defer cmdKey.Close()
	return cmdKey.SetStringValue("", `"`+execPath+`" "%1"`) == nil
}

// UnregisterHarAssociation removes registry entries for .har association.
func UnregisterHarAssociation() bool {
	_ = registry.DeleteKey(registry.CURRENT_USER, `Software\Classes\.har`)
	_ = deleteRegistryKeyRecursively(registry.CURRENT_USER, `Software\Classes\httptoolkitpro.har`)
	return true
}

func deleteRegistryKeyRecursively(root registry.Key, path string) error {
	k, err := registry.OpenKey(root, path, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()

	names, err := k.ReadSubKeyNames(-1)
	if err == nil {
		for _, name := range names {
			_ = deleteRegistryKeyRecursively(k, name)
		}
	}
	return registry.DeleteKey(root, path)
}
