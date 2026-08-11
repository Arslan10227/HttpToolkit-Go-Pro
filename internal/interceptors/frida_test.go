package interceptors

import (
	"testing"
)

func TestParseFridaPsTabular(t *testing.T) {
	rawOutput := `  PID  Name           Identifier
-----  -------------  -------------------------
 3121  Settings       com.android.settings
    -  Chrome         com.android.chrome
  123  My App Name    org.example.app.name`

	targets := parseFridaPsTabular(rawOutput)
	if len(targets) != 3 {
		t.Fatalf("Expected 3 targets, got %d", len(targets))
	}

	if targets[0]["id"] != "com.android.settings" || targets[0]["name"] != "Settings" {
		t.Fatalf("Expected Settings app target, got %+v", targets[0])
	}
	if targets[1]["id"] != "com.android.chrome" || targets[1]["name"] != "Chrome" {
		t.Fatalf("Expected Chrome app target, got %+v", targets[1])
	}
	if targets[2]["id"] != "org.example.app.name" || targets[2]["name"] != "My App Name" {
		t.Fatalf("Expected My App Name target, got %+v", targets[2])
	}
}

func TestParseFridaPsTabularFallback(t *testing.T) {
	// Tests column fallback parsing if headers are missing or formatted oddly
	rawOutput := `Some random header text here
 3121  Settings       com.android.settings
    -  Chrome         com.android.chrome`

	targets := parseFridaPsTabular(rawOutput)
	if len(targets) != 2 {
		t.Fatalf("Expected 2 targets, got %d", len(targets))
	}
	if targets[0]["id"] != "com.android.settings" || targets[0]["name"] != "Settings" {
		t.Fatalf("Expected Settings app, got %+v", targets[0])
	}
	if targets[1]["id"] != "com.android.chrome" || targets[1]["name"] != "Chrome" {
		t.Fatalf("Expected Chrome app, got %+v", targets[1])
	}
}
