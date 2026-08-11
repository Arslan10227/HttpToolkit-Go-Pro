package origins

import "regexp"

var DevOrigins = []*regexp.Regexp{
	regexp.MustCompile(`^https?://localhost(:\d+)?$`),
	regexp.MustCompile(`^https?://127\.0\.0\.\d+(:\d+)?$`),
	regexp.MustCompile(`^https?://\[::1\](:\d+)?$`),
	regexp.MustCompile(`^http://local\.httptoolkit\.tech(:\d+)?$`),
	regexp.MustCompile(`^https://app\.httptoolkit\.tech$`),
	regexp.MustCompile(`^https://[a-z0-9-]+\.vercel\.app$`),
	regexp.MustCompile(`^https?://wails\.localhost(:\d+)?$`),
	regexp.MustCompile(`^https?://tauri\.localhost(:\d+)?$`),
	regexp.MustCompile(`^https://asset\.localhost$`),
	regexp.MustCompile(`^https?://192\.168\.\d{1,3}\.\d{1,3}(:\d+)?$`),
	regexp.MustCompile(`^https?://10\.\d{1,3}\.\d{1,3}\.\d{1,3}(:\d+)?$`),
	regexp.MustCompile(`^https?://172\.(1[6-9]|2\d|3[01])\.\d{1,3}\.\d{1,3}(:\d+)?$`),
}

var ProdOrigins = []*regexp.Regexp{
	regexp.MustCompile(`^https://app\.httptoolkit\.tech$`),
	regexp.MustCompile(`^https://[a-z0-9-]+\.vercel\.app$`),
	regexp.MustCompile(`^https?://wails\.localhost(:\d+)?$`),
	regexp.MustCompile(`^https?://tauri\.localhost(:\d+)?$`),
	regexp.MustCompile(`^https://asset\.localhost$`),
	regexp.MustCompile(`^https?://localhost(:\d+)?$`),
	regexp.MustCompile(`^https?://127\.0\.0\.\d+(:\d+)?$`),
	regexp.MustCompile(`^https?://\[::1\](:\d+)?$`),
}

func IsAllowed(origin string, devMode bool) bool {
	if origin == "" {
		return devMode
	}
	list := ProdOrigins
	if devMode {
		list = DevOrigins
	}
	for _, p := range list {
		if p.MatchString(origin) {
			return true
		}
	}
	return false
}
