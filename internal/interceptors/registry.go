package interceptors

import (
	"fmt"
	"sync"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/cert"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/config"
	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/logger"
)

type Interceptor interface {
	ID() string
	Version() string
	IsActivable() (bool, error)
	IsActive(proxyPort int) (bool, error)
	Activate(proxyPort int, options map[string]any) (map[string]any, error)
	Deactivate(proxyPort int, options map[string]any) error
	Metadata(kind string) (any, error)
}

type Registry struct {
	cfg          *config.Config
	spki         string
	certs        *cert.Manager
	interceptors map[string]Interceptor
	order        []string
	mu           sync.Mutex
}

func NewRegistry(cfg *config.Config, spki string, certs *cert.Manager) *Registry {
	all := buildAll(cfg, spki, certs)
	r := &Registry{
		cfg:          cfg,
		spki:         spki,
		certs:        certs,
		interceptors: make(map[string]Interceptor),
		order:        make([]string, 0, len(all)),
	}
	for _, i := range all {
		r.interceptors[i.ID()] = i
		r.order = append(r.order, i.ID())
	}
	return r
}

func (r *Registry) List(proxyPort int) []map[string]any {
	var wg sync.WaitGroup
	var mu sync.Mutex

	ids := r.order
	if len(ids) == 0 {
		ids = make([]string, 0, len(r.interceptors))
		for id := range r.interceptors {
			ids = append(ids, id)
		}
	}

	results := make([]map[string]any, len(ids))

	for idx, id := range ids {
		wg.Add(1)
		go func(index int, interceptorID string) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error(fmt.Errorf("panic in interceptor list: %v", rec), map[string]any{
						"interceptorId": interceptorID,
					})
				}
				wg.Done()
			}()
			i := r.interceptors[interceptorID]
			activable, _ := i.IsActivable()

			active := false
			if activable && proxyPort > 0 {
				active, _ = i.IsActive(proxyPort)
			}

			meta, _ := i.Metadata("summary")
			item := map[string]any{
				"id": i.ID(), "version": i.Version(), "isActivable": activable,
			}
			if meta != nil {
				item["metadata"] = meta
			}
			if proxyPort > 0 {
				item["isActive"] = active
			}

			mu.Lock()
			results[index] = item
			mu.Unlock()
		}(idx, id)
	}

	wg.Wait()

	out := make([]map[string]any, 0, len(results))
	for _, res := range results {
		if res != nil {
			out = append(out, res)
		}
	}
	return out
}

type SubMetadataProvider interface {
	SubMetadata(subID string) (any, error)
}

func (r *Registry) SubMetadata(id, subID string) (any, error) {
	i, ok := r.interceptors[id]
	if !ok {
		return nil, fmt.Errorf("unknown interceptor")
	}
	if sp, ok := i.(SubMetadataProvider); ok {
		return sp.SubMetadata(subID)
	}
	return i.Metadata("detailed")
}

func (r *Registry) Metadata(id string) (any, error) {
	i, ok := r.interceptors[id]
	if !ok {
		return nil, fmt.Errorf("unknown interceptor")
	}
	return i.Metadata("detailed")
}

func (r *Registry) Activate(id string, proxyPort int, options map[string]any) (map[string]any, error) {
	i, ok := r.interceptors[id]
	if !ok {
		return nil, fmt.Errorf("unknown interceptor %s", id)
	}
	meta, err := i.Activate(proxyPort, options)
	if err != nil {
		return map[string]any{"success": false, "metadata": map[string]any{"error": err.Error()}}, nil
	}
	return map[string]any{"success": true, "metadata": meta}, nil
}

func (r *Registry) Deactivate(id string, proxyPort int, options map[string]any) (map[string]any, error) {
	i, ok := r.interceptors[id]
	if !ok {
		return nil, fmt.Errorf("unknown interceptor %s", id)
	}
	if err := i.Deactivate(proxyPort, options); err != nil {
		return nil, err
	}
	active, _ := i.IsActive(proxyPort)
	return map[string]any{"success": !active}, nil
}

func (r *Registry) DeactivateAll() {
	for _, i := range r.interceptors {
		_ = i.Deactivate(0, nil)
	}
}

// AndroidAdbCertStatus checks if the CA certificate is installed in the
// Android system trust store for the given ADB device ID.
func (r *Registry) AndroidAdbCertStatus(deviceId, certPEM string) (map[string]any, error) {
	i, ok := r.interceptors["android-adb"]
	if !ok {
		return nil, fmt.Errorf("android-adb interceptor unavailable")
	}
	type adbCertChecker interface {
		getDeviceCertificateStatus(deviceId, certPEM string) (map[string]any, error)
	}
	a, ok := i.(adbCertChecker)
	if !ok {
		return nil, fmt.Errorf("android-adb interceptor does not support certificate status")
	}
	return a.getDeviceCertificateStatus(deviceId, certPEM)
}

// AndroidAdbCertInstall installs the CA certificate into the Android system
// trust store for the given ADB device ID (requires root).
func (r *Registry) AndroidAdbCertInstall(deviceId, certPEM, certPath string) (map[string]any, error) {
	i, ok := r.interceptors["android-adb"]
	if !ok {
		return nil, fmt.Errorf("android-adb interceptor unavailable")
	}
	type adbCertInstaller interface {
		installDeviceCertificate(deviceId, certPEM, certPath string) (map[string]any, error)
	}
	a, ok := i.(adbCertInstaller)
	if !ok {
		return nil, fmt.Errorf("android-adb interceptor does not support certificate install")
	}
	return a.installDeviceCertificate(deviceId, certPEM, certPath)
}
