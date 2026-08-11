//go:build !windows

package system

func IsHarAssociated() bool {
	return false
}

func RegisterHarAssociation() bool {
	return false
}

func UnregisterHarAssociation() bool {
	return false
}
