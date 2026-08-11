package mitm

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Arslan10227/HttpToolkit-Go-Pro/internal/proxy/rules"
)

// serveAmiusingFallback intercepts requests to amiusing.httptoolkit.tech and
// serves the same self-hosted "am I being intercepted" page from the proxy
// itself. This is a safety net for any cached Chrome state, older builds, or
// hardcoded UI references that still request the old external domain.
func (s *Server) serveAmiusingFallback(w io.Writer, req *http.Request, id string) {
	start := time.Now()
	path := req.URL.Path
	if path == "" || path == "/" || path == "/amiusing" {
		body := amiusingFallbackPage("https://amiusing.httptoolkit.tech/amiusing/intercepted")
		writeFixedResponse(w, id, s.events, rules.StepAction{
			Kind:       "fixed-response",
			StatusCode: http.StatusOK,
			Body:       body,
			Headers:    map[string]string{"Content-Type": "text/html; charset=utf-8"},
		}, start)
		return
	}
	if path == "/amiusing/intercepted" {
		writeFixedResponse(w, id, s.events, rules.StepAction{
			Kind:       "fixed-response",
			StatusCode: http.StatusOK,
			Body:       "intercepted",
			Headers:    map[string]string{"Content-Type": "text/plain"},
		}, start)
		return
	}
	writeFixedResponse(w, id, s.events, rules.StepAction{
		Kind:       "fixed-response",
		StatusCode: http.StatusNotFound,
		Body:       "not found",
		Headers:    map[string]string{"Content-Type": "text/plain"},
	}, start)
}

func amiusingFallbackPage(testURL string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>HTTP Toolkit Interception Check</title>
<style>
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #1e2028; color: #fff; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; }
.card { background: #2a2d37; padding: 40px; border-radius: 12px; text-align: center; max-width: 420px; }
h1 { margin-top: 0; font-size: 1.2em; }
p { color: #a0a5b5; line-height: 1.5; }
.success { color: #00e676; }
.error { color: #ff5252; }
</style>
</head>
<body>
<div class="card" id="card">
  <h1>Checking interception...</h1>
  <p id="msg">Verifying that HTTP Toolkit is intercepting this browser's traffic.</p>
</div>
<script>
const testUrl = %q;
fetch(testUrl, { cache: "no-store" })
  .then(r => {
    if (r.ok) {
      document.getElementById("card").innerHTML = '<h1 class="success">Interception active</h1><p>You can close this tab.</p>';
      setTimeout(() => window.close(), 1000);
    } else {
      fail();
    }
  })
  .catch(fail);

function fail() {
  document.getElementById("card").innerHTML = '<h1 class="error">Interception not active</h1><p>The HTTP Toolkit certificate is not trusted in this browser. Install and trust the CA, then restart the browser.</p>';
}
</script>
</body>
</html>`, testURL)
}

func isAmiusingHost(host string) bool {
	host = strings.ToLower(host)
	return host == "amiusing.httptoolkit.tech" || strings.HasSuffix(host, ".amiusing.httptoolkit.tech")
}
