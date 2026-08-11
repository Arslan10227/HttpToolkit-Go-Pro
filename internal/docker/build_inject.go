package docker

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type memoryTarEntry struct {
	header tar.Header
	body   []byte
}

// InjectBuildContext rewrites a docker build tarball to inject interception assets.
func InjectBuildContext(raw []byte, dockerfileName string, settings BuildSettings) ([]byte, int, error) {
	if dockerfileName == "" {
		dockerfileName = "Dockerfile"
	}
	gzipped := isGzip(raw)
	entries, err := readTarEntries(raw, gzipped)
	if err != nil {
		return nil, 0, err
	}

	normalizedName := normalizeTarName(dockerfileName)
	commandsAdded := 0
	for i, e := range entries {
		if normalizeTarName(e.header.Name) != normalizedName {
			continue
		}
		updated, added := InjectIntoDockerfile(string(e.body), settings.ProxyPort, buildInjectionEnv(settings.ProxyPort))
		commandsAdded = added
		e.body = []byte(updated)
		e.header.Size = int64(len(e.body))
		entries[i] = e
		break
	}

	cert := settings.CertContent
	if len(cert) == 0 && settings.CertPath != "" {
		cert, _ = os.ReadFile(settings.CertPath)
	}
	if len(cert) > 0 {
		entries = append(entries, memoryTarEntry{
			header: tar.Header{Name: normalizeTarName(contextCAPath), Mode: 0o644, Size: int64(len(cert))},
			body:   cert,
		})
	}

	if settings.AssetsDir != "" {
		overrideRoot := filepath.Join(settings.AssetsDir, "overrides")
		overrideEntries, err := packOverridesDir(overrideRoot)
		if err != nil {
			return nil, 0, err
		}
		entries = append(entries, overrideEntries...)
	}

	entries = append(entries, memoryTarEntry{
		header: tar.Header{
			Name:     normalizeTarName(contextInjectPath),
			Mode:     0o755,
			Typeflag: tar.TypeDir,
		},
	})

	out, err := writeTarEntries(entries, gzipped)
	return out, commandsAdded, err
}

func normalizeTarName(name string) string {
	name = strings.ReplaceAll(name, "\\", "/")
	return strings.TrimPrefix(name, "./")
}

func isGzip(b []byte) bool {
	return len(b) >= 2 && b[0] == 0x1f && b[1] == 0x8b
}

func readTarEntries(raw []byte, gzipped bool) ([]memoryTarEntry, error) {
	reader := bytes.NewReader(raw)
	var tr *tar.Reader
	if gzipped {
		gr, err := gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
		defer gr.Close()
		tr = tar.NewReader(gr)
	} else {
		tr = tar.NewReader(reader)
	}

	var entries []memoryTarEntry
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		body, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
		copyHdr := *hdr
		entries = append(entries, memoryTarEntry{header: copyHdr, body: body})
	}
	return entries, nil
}

func writeTarEntries(entries []memoryTarEntry, gzipped bool) ([]byte, error) {
	var buf bytes.Buffer
	var tw *tar.Writer
	var gw *gzip.Writer
	if gzipped {
		gw = gzip.NewWriter(&buf)
		tw = tar.NewWriter(gw)
	} else {
		tw = tar.NewWriter(&buf)
	}
	for _, e := range entries {
		hdr := e.header
		if hdr.Typeflag == 0 && hdr.Mode == 0 {
			hdr.Mode = 0o644
		}
		if hdr.Typeflag != tar.TypeDir {
			hdr.Size = int64(len(e.body))
		}
		if err := tw.WriteHeader(&hdr); err != nil {
			return nil, err
		}
		if len(e.body) > 0 {
			if _, err := tw.Write(e.body); err != nil {
				return nil, err
			}
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	if gw != nil {
		if err := gw.Close(); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func packOverridesDir(root string) ([]memoryTarEntry, error) {
	if _, err := os.Stat(root); err != nil {
		return nil, nil
	}
	var entries []memoryTarEntry
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			return nil
		}
		name := normalizeTarName(contextOverridesPath + "/" + rel)
		if info.IsDir() {
			entries = append(entries, memoryTarEntry{
				header: tar.Header{Name: name, Mode: 0o555, Typeflag: tar.TypeDir},
			})
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		entries = append(entries, memoryTarEntry{
			header: tar.Header{Name: name, Mode: 0o555, Size: int64(len(body))},
			body:   body,
		})
		return nil
	})
	return entries, err
}
