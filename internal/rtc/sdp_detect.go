package rtc

import (
	"strings"
)

// LooksLikeSDP reports whether body content appears to be an SDP offer/answer.
func LooksLikeSDP(body []byte, contentType string) bool {
	if len(body) == 0 {
		return false
	}
	ct := strings.ToLower(contentType)
	if strings.Contains(ct, "application/sdp") || strings.Contains(ct, "application/sdp+json") {
		return true
	}
	text := string(body)
	if !strings.Contains(text, "v=0") {
		return false
	}
	// WebRTC SDP uses m=application for SCTP data channels, plus audio/video tracks.
	return strings.Contains(text, "m=application") ||
		strings.Contains(text, "m=audio") ||
		strings.Contains(text, "m=video")
}
