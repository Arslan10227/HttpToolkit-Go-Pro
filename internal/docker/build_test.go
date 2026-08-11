package docker

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"strings"
	"testing"
)

func TestInjectIntoDockerfile(t *testing.T) {
	df := "FROM alpine:3.18\nRUN echo hi\nFROM ubuntu:22.04\nCMD bash\n"
	out, added := InjectIntoDockerfile(df, 8000, buildInjectionEnv(8000))
	if added == 0 {
		t.Fatal("expected injection commands")
	}
	if !strings.Contains(out, BuildLabel+"=started-8000") {
		t.Fatalf("missing start label: %s", out)
	}
	if !strings.Contains(out, "COPY "+contextInjectPath) {
		t.Fatalf("missing COPY injection: %s", out)
	}
	if strings.Count(strings.ToUpper(out), "FROM ") != 2 {
		t.Fatalf("expected two FROM lines preserved")
	}
}

func TestInjectBuildContext(t *testing.T) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	_ = tw.WriteHeader(&tar.Header{Name: "Dockerfile", Mode: 0o644, Size: 12})
	_, _ = tw.Write([]byte("FROM alpine\n"))
	_ = tw.Close()

	out, added, err := InjectBuildContext(buf.Bytes(), "Dockerfile", BuildSettings{
		ProxyPort:   8000,
		CertContent: []byte("test-cert"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if added == 0 {
		t.Fatal("expected dockerfile commands added")
	}

	tr := tar.NewReader(bytes.NewReader(out))
	foundDockerfile := false
	foundCert := false
	for {
		hdr, err := tr.Next()
		if err != nil {
			break
		}
		if hdr.Name == "Dockerfile" || hdr.Name == "./Dockerfile" {
			foundDockerfile = true
		}
		if strings.HasSuffix(hdr.Name, "ca.pem") {
			foundCert = true
		}
	}
	if !foundDockerfile || !foundCert {
		t.Fatalf("missing injected tar entries")
	}
}

func TestInjectBuildContextGzip(t *testing.T) {
	var raw bytes.Buffer
	gw := gzip.NewWriter(&raw)
	tw := tar.NewWriter(gw)
	_ = tw.WriteHeader(&tar.Header{Name: "Dockerfile", Mode: 0o644, Size: 12})
	_, _ = tw.Write([]byte("FROM alpine\n"))
	_ = tw.Close()
	_ = gw.Close()

	out, _, err := InjectBuildContext(raw.Bytes(), "Dockerfile", BuildSettings{
		ProxyPort:   8000,
		CertContent: []byte("cert"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !isGzip(out) {
		t.Fatal("expected gzipped output")
	}
}

func TestTransformBuildOutput(t *testing.T) {
	in := []byte(`{"stream":"Step 1/4 : FROM alpine\n"}` + "\n" +
		`{"stream":"LABEL ` + BuildLabel + `=started-8000\n"}` + "\n" +
		`{"stream":"Step 2/4 : COPY files\n"}` + "\n" +
		`{"stream":"LABEL ` + BuildLabel + `=8000\n"}` + "\n" +
		`{"stream":" ---\> abcdef\n"}` + "\n" +
		`{"stream":"Step 3/4 : RUN echo ok\n"}`)

	out := TransformBuildOutput(in, 2)
	if !strings.Contains(string(out), "Enabling HTTP Toolkit interception") {
		t.Fatalf("missing interception banner: %s", out)
	}
	if strings.Contains(string(out), "started-8000") {
		t.Fatalf("should hide start label noise: %s", out)
	}
}
