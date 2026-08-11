package config

import (
	_ "embed"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed version.yaml
var versionYAML []byte

// AppMeta holds the application identity loaded from version.yaml.
// This is the single source of truth for app name, version, title, etc.
// Edit internal/config/version.yaml to change these values.
type AppMeta struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Author      string `yaml:"author"`
	Repository  string `yaml:"repository"`
	License     string `yaml:"license"`
}

// appMeta is parsed once at init time from the embedded version.yaml.
var appMeta AppMeta

func init() {
	if err := yaml.Unmarshal(versionYAML, &appMeta); err != nil {
		// Fallback to safe defaults if the YAML is malformed.
		appMeta = AppMeta{
			Name:        "HttpToolkit Go Pro",
			Version:     "1.0.0-go",
			Title:       "Httptoolkit Go",
			Description: "Native Go MITM proxy + Wails desktop shell for HttpToolkit Pro",
			Author:      "Arslan10227",
			Repository:  "https://github.com/Arslan10227/HttpToolkit-Go-Pro",
			License:     "AGPL-3.0",
		}
	}
}

// AppName returns the application name from version.yaml.
func AppName() string { return appMeta.Name }

// AppVersion returns the application version string from version.yaml.
func AppVersion() string { return appMeta.Version }

// AppTitle returns the window/title bar name from version.yaml.
func AppTitle() string { return appMeta.Title }

// AppDescription returns the short description from version.yaml.
func AppDescription() string { return appMeta.Description }

// AppAuthor returns the author name from version.yaml.
func AppAuthor() string { return appMeta.Author }

// AppRepository returns the repository URL from version.yaml.
func AppRepository() string { return appMeta.Repository }

// AppLicense returns the license identifier from version.yaml.
func AppLicense() string { return appMeta.License }

// AppMetaString returns a human-readable summary of the app metadata.
func AppMetaString() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s v%s\n", appMeta.Name, appMeta.Version)
	fmt.Fprintf(&b, "  %s\n", appMeta.Description)
	fmt.Fprintf(&b, "  Author: %s\n", appMeta.Author)
	fmt.Fprintf(&b, "  Repo:   %s\n", appMeta.Repository)
	fmt.Fprintf(&b, "  License: %s", appMeta.License)
	return b.String()
}
