//go:build !windows

package system

// RegisterProtocol registers the custom URI scheme (stub for non-Windows).
func RegisterProtocol() error {
	return nil
}
