package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// FindJava resolves the path to the system's Java executable.
// It checks PATH, JAVA_HOME, and standard Windows directories.
// When multiple candidates exist, a JDK installation is preferred over a JRE
// because JVM attach requires the com.sun.tools.attach module / tools.jar.
func findJavaCandidates() []string {
	candidates := []string{}

	// 1. Try PATH
	if p, err := exec.LookPath("java"); err == nil {
		candidates = append(candidates, p)
	}
	// 2. Try JAVA_HOME
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		ext := ""
		if runtime.GOOS == "windows" {
			ext = ".exe"
		}
		p := filepath.Join(jh, "bin", "java"+ext)
		if _, err := os.Stat(p); err == nil {
			candidates = append(candidates, p)
		}
	}
	// 3. Try standard paths on Windows
	if runtime.GOOS == "windows" {
		commonPaths := []string{
			`C:\Program Files\Common Files\Oracle\Java\javapath\java.exe`,
			`C:\Program Files\Java`,
			`C:\Program Files\Eclipse Adoptium`,
			`C:\Program Files\BellSoft`,
			`C:\Program Files\Amazon Corretto`,
			`C:\Program Files\Zulu`,
			`C:\Program Files\Microsoft`,
		}
		for _, cp := range commonPaths {
			if strings.HasSuffix(cp, ".exe") {
				if _, err := os.Stat(cp); err == nil {
					candidates = append(candidates, cp)
				}
			} else {
				// Search folders
				if entries, err := os.ReadDir(cp); err == nil {
					for _, entry := range entries {
						if entry.IsDir() {
							p := filepath.Join(cp, entry.Name(), "bin", "java.exe")
							if _, err := os.Stat(p); err == nil {
								candidates = append(candidates, p)
							}
						}
					}
				}
			}
		}
	}

	return candidates
}

// FindJava resolves the path to the system's Java executable, preferring a JDK.
func FindJava() string {
	candidates := findJavaCandidates()
	for _, p := range candidates {
		if isJDKJava(p) {
			return p
		}
	}
	if len(candidates) > 0 {
		return candidates[0]
	}
	return "java" // Fallback
}

// FindJDK resolves the path to a JDK Java executable that supports attach.
// It returns an empty string if no JDK is found.
func FindJDK() string {
	candidates := findJavaCandidates()
	for _, p := range candidates {
		if isJDKJava(p) {
			return p
		}
	}
	return ""
}

// isJDKJava returns true if the java executable belongs to a JDK by checking
// for sibling tools that are only present in a JDK (jcmd/javac).
func isJDKJava(javaPath string) bool {
	base := strings.TrimSuffix(javaPath, filepath.Ext(javaPath))
	base = strings.TrimSuffix(base, "java")
	siblings := []string{base + "jcmd", base + "javac"}
	if runtime.GOOS == "windows" {
		for i := range siblings {
			siblings[i] += ".exe"
		}
	}
	for _, s := range siblings {
		if _, err := os.Stat(s); err == nil {
			return true
		}
	}
	return false
}
