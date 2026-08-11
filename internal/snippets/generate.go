package snippets

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type HarHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HarQuery struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type HarPostData struct {
	MimeType string `json:"mimeType"`
	Text     string `json:"text"`
}

// HarRequest matches the HAR entry.request shape sent by the UI.
type HarRequest struct {
	Method      string       `json:"method"`
	URL         string       `json:"url"`
	Headers     []HarHeader  `json:"headers"`
	QueryString []HarQuery   `json:"queryString"`
	PostData    *HarPostData `json:"postData"`
}

func ParseHarRequest(raw json.RawMessage) (HarRequest, error) {
	var req HarRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return HarRequest{}, err
	}
	if req.Method == "" || req.URL == "" {
		return HarRequest{}, fmt.Errorf("request method and url required")
	}
	return req, nil
}

func Generate(req HarRequest, target, client string) (string, error) {
	target = strings.ToLower(strings.TrimSpace(target))
	client = strings.ToLower(strings.TrimSpace(client))
	if target == "" || client == "" {
		return "", fmt.Errorf("target and client required")
	}

	switch target {
	case "shell":
		return generateShell(req, client)
	case "http":
		return generateHTTP(req, client)
	case "javascript", "node":
		return generateJS(req, target, client)
	case "python":
		return generatePython(req, client)
	case "go":
		return generateGo(req, client)
	case "java":
		return generateJava(req, client)
	case "php":
		return generatePHP(req, client)
	case "ruby":
		return generateRuby(req, client)
	case "csharp":
		return generateCSharp(req, client)
	case "swift":
		return generateSwift(req, client)
	case "kotlin":
		return generateKotlin(req, client)
	case "c":
		return generateC(req, client)
	case "objc":
		return generateObjC(req, client)
	case "r":
		return generateR(req, client)
	case "clojure":
		return generateClojure(req, client)
	case "ocaml":
		return generateOCaml(req, client)
	case "crystal":
		return generateCrystal(req, client)
	case "rust":
		return generateRust(req, client)
	case "powershell":
		return generatePowerShell(req, client)
	default:
		return generateShell(req, "curl")
	}
}

func headerMap(req HarRequest) map[string]string {
	out := make(map[string]string, len(req.Headers))
	for _, h := range req.Headers {
		if h.Name == "" {
			continue
		}
		out[h.Name] = h.Value
	}
	return out
}

func bodyText(req HarRequest) string {
	if req.PostData == nil {
		return ""
	}
	return req.PostData.Text
}

func hasBody(req HarRequest) bool {
	return req.PostData != nil && req.PostData.Text != ""
}

func shellEscape(s string) string {
	return strings.ReplaceAll(s, "'", "'\\''")
}

func generateShell(req HarRequest, client string) (string, error) {
	switch client {
	case "httpie":
		var b strings.Builder
		fmt.Fprintf(&b, "http %s '%s'", req.Method, shellEscape(req.URL))
		for _, h := range req.Headers {
			fmt.Fprintf(&b, " '%s:%s'", h.Name, h.Value)
		}
		if hasBody(req) {
			fmt.Fprintf(&b, " <<'EOF'\n%s\nEOF", req.PostData.Text)
		}
		return b.String(), nil
	case "wget":
		var b strings.Builder
		fmt.Fprintf(&b, "wget --method=%s", req.Method)
		for _, h := range req.Headers {
			fmt.Fprintf(&b, " --header='%s: %s'", h.Name, shellEscape(h.Value))
		}
		if hasBody(req) {
			fmt.Fprintf(&b, " --body-data='%s'", shellEscape(bodyText(req)))
		}
		fmt.Fprintf(&b, " '%s'", shellEscape(req.URL))
		return b.String(), nil
	default:
		return buildCurl(req), nil
	}
}

func buildCurl(req HarRequest) string {
	var b strings.Builder
	fmt.Fprintf(&b, "curl -X %s", req.Method)
	for _, h := range req.Headers {
		fmt.Fprintf(&b, " -H '%s: %s'", h.Name, shellEscape(h.Value))
	}
	if hasBody(req) {
		fmt.Fprintf(&b, " -d '%s'", shellEscape(bodyText(req)))
	}
	fmt.Fprintf(&b, " '%s'", shellEscape(req.URL))
	return b.String()
}

func generateHTTP(req HarRequest, client string) (string, error) {
	_ = client
	u, err := url.Parse(req.URL)
	if err != nil {
		return "", err
	}
	path := u.RequestURI()
	if path == "" {
		path = "/"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s HTTP/1.1\r\n", req.Method, path)
	fmt.Fprintf(&b, "Host: %s\r\n", u.Host)
	for _, h := range req.Headers {
		if strings.EqualFold(h.Name, "host") {
			continue
		}
		fmt.Fprintf(&b, "%s: %s\r\n", h.Name, h.Value)
	}
	if hasBody(req) {
		if !headerHas(req, "content-length") {
			fmt.Fprintf(&b, "Content-Length: %d\r\n", len(bodyText(req)))
		}
	}
	b.WriteString("\r\n")
	if hasBody(req) {
		b.WriteString(bodyText(req))
	}
	return b.String(), nil
}

func headerHas(req HarRequest, name string) bool {
	for _, h := range req.Headers {
		if strings.EqualFold(h.Name, name) {
			return true
		}
	}
	return false
}

func generateJS(req HarRequest, target, client string) (string, error) {
	switch client {
	case "jquery":
		var b strings.Builder
		b.WriteString("$.ajax({\n")
		fmt.Fprintf(&b, "  url: '%s',\n", shellEscape(req.URL))
		fmt.Fprintf(&b, "  method: '%s',\n", req.Method)
		if len(req.Headers) > 0 {
			b.WriteString("  headers: {\n")
			for _, h := range req.Headers {
				fmt.Fprintf(&b, "    '%s': '%s',\n", h.Name, shellEscape(h.Value))
			}
			b.WriteString("  },\n")
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "  data: '%s',\n", shellEscape(bodyText(req)))
		}
		b.WriteString("  success: (data) => console.log(data)\n});")
		return b.String(), nil
	case "xhr":
		var b strings.Builder
		b.WriteString("const xhr = new XMLHttpRequest();\n")
		fmt.Fprintf(&b, "xhr.open('%s', '%s');\n", req.Method, shellEscape(req.URL))
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "xhr.setRequestHeader('%s', '%s');\n", h.Name, shellEscape(h.Value))
		}
		b.WriteString("xhr.onload = () => console.log(xhr.responseText);\n")
		if hasBody(req) {
			fmt.Fprintf(&b, "xhr.send('%s');", shellEscape(bodyText(req)))
		} else {
			b.WriteString("xhr.send();")
		}
		return b.String(), nil
	case "axios":
		var b strings.Builder
		b.WriteString("const axios = require('axios');\n\n")
		b.WriteString("axios({\n")
		fmt.Fprintf(&b, "  method: '%s',\n", req.Method)
		fmt.Fprintf(&b, "  url: '%s',\n", shellEscape(req.URL))
		if len(req.Headers) > 0 {
			b.WriteString("  headers: {\n")
			for _, h := range req.Headers {
				fmt.Fprintf(&b, "    '%s': '%s',\n", h.Name, shellEscape(h.Value))
			}
			b.WriteString("  },\n")
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "  data: '%s',\n", shellEscape(bodyText(req)))
		}
		b.WriteString("}).then(res => console.log(res.data));")
		return b.String(), nil
	case "request":
		var b strings.Builder
		b.WriteString("const request = require('request');\n\n")
		b.WriteString("request({\n")
		fmt.Fprintf(&b, "  method: '%s',\n", req.Method)
		fmt.Fprintf(&b, "  url: '%s',\n", shellEscape(req.URL))
		if len(req.Headers) > 0 {
			b.WriteString("  headers: {\n")
			for _, h := range req.Headers {
				fmt.Fprintf(&b, "    '%s': '%s',\n", h.Name, shellEscape(h.Value))
			}
			b.WriteString("  },\n")
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "  body: '%s',\n", shellEscape(bodyText(req)))
		}
		b.WriteString("}, (err, res, body) => console.log(body));")
		return b.String(), nil
	case "fetch":
		if target == "node" {
			var b strings.Builder
			b.WriteString("const fetch = require('node-fetch');\n\n")
			b.WriteString("fetch('")
			b.WriteString(shellEscape(req.URL))
			b.WriteString("', {\n")
			fmt.Fprintf(&b, "  method: '%s',\n", req.Method)
			if len(req.Headers) > 0 {
				b.WriteString("  headers: {\n")
				for _, h := range req.Headers {
					fmt.Fprintf(&b, "    '%s': '%s',\n", h.Name, shellEscape(h.Value))
				}
				b.WriteString("  },\n")
			}
			if hasBody(req) {
				fmt.Fprintf(&b, "  body: '%s',\n", shellEscape(bodyText(req)))
			}
			b.WriteString("}).then(res => res.text()).then(console.log);")
			return b.String(), nil
		}
		fallthrough
	case "unirest":
		if target == "node" && client == "unirest" {
			var b strings.Builder
			b.WriteString("const unirest = require('unirest');\n\n")
			b.WriteString("unirest('")
			b.WriteString(req.Method)
			b.WriteString("', '")
			b.WriteString(shellEscape(req.URL))
			b.WriteString("')")
			for _, h := range req.Headers {
				fmt.Fprintf(&b, "\n  .header('%s', '%s')", h.Name, shellEscape(h.Value))
			}
			if hasBody(req) {
				fmt.Fprintf(&b, "\n  .send('%s')", shellEscape(bodyText(req)))
			}
			b.WriteString("\n  .end(res => console.log(res.body));")
			return b.String(), nil
		}
		fallthrough
	case "native":
		if target == "node" {
			return generateNodeNative(req)
		}
		fallthrough
	default:
		var b strings.Builder
		b.WriteString("const options = {\n")
		fmt.Fprintf(&b, "  method: '%s',\n", req.Method)
		if len(req.Headers) > 0 {
			b.WriteString("  headers: {\n")
			for _, h := range req.Headers {
				fmt.Fprintf(&b, "    '%s': '%s',\n", h.Name, shellEscape(h.Value))
			}
			b.WriteString("  },\n")
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "  body: '%s',\n", shellEscape(bodyText(req)))
		}
		b.WriteString("};\n\n")
		fmt.Fprintf(&b, "fetch('%s', options).then(res => res.text()).then(console.log);", shellEscape(req.URL))
		return b.String(), nil
	}
}

func generateNodeNative(req HarRequest) (string, error) {
	u, err := url.Parse(req.URL)
	if err != nil {
		return "", err
	}
	scheme := "http"
	if u.Scheme == "https" {
		scheme = "https"
	}
	var b strings.Builder
	b.WriteString("const http = require('" + scheme + "');\n\n")
	b.WriteString("const options = {\n")
	fmt.Fprintf(&b, "  hostname: '%s',\n", u.Hostname())
	if u.Port() != "" {
		fmt.Fprintf(&b, "  port: %s,\n", u.Port())
	}
	fmt.Fprintf(&b, "  path: '%s',\n", shellEscape(u.RequestURI()))
	fmt.Fprintf(&b, "  method: '%s',\n", req.Method)
	if len(req.Headers) > 0 {
		b.WriteString("  headers: {\n")
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "    '%s': '%s',\n", h.Name, shellEscape(h.Value))
		}
		b.WriteString("  },\n")
	}
	b.WriteString("};\n\n")
	b.WriteString("const req = http.request(options, (res) => {\n")
	b.WriteString("  res.on('data', (d) => process.stdout.write(d));\n")
	b.WriteString("});\n")
	if hasBody(req) {
		fmt.Fprintf(&b, "req.write('%s');\n", shellEscape(bodyText(req)))
	}
	b.WriteString("req.end();")
	return b.String(), nil
}

func generatePython(req HarRequest, client string) (string, error) {
	switch client {
	case "python3":
		var b strings.Builder
		b.WriteString("import urllib.request\n\n")
		fmt.Fprintf(&b, "req = urllib.request.Request('%s', method='%s'", shellEscape(req.URL), req.Method)
		if hasBody(req) {
			fmt.Fprintf(&b, ", data=b'%s'", bodyText(req))
		}
		b.WriteString(")\n")
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "req.add_header('%s', '%s')\n", h.Name, shellEscape(h.Value))
		}
		b.WriteString("with urllib.request.urlopen(req) as resp:\n    print(resp.read().decode())")
		return b.String(), nil
	default:
		var b strings.Builder
		b.WriteString("import requests\n\n")
		b.WriteString("response = requests.request(\n")
		fmt.Fprintf(&b, "    method='%s',\n", req.Method)
		fmt.Fprintf(&b, "    url='%s',\n", shellEscape(req.URL))
		if len(req.Headers) > 0 {
			b.WriteString("    headers={\n")
			for _, h := range req.Headers {
				fmt.Fprintf(&b, "        '%s': '%s',\n", h.Name, shellEscape(h.Value))
			}
			b.WriteString("    },\n")
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "    data='''%s''',\n", bodyText(req))
		}
		b.WriteString(")\nprint(response.text)")
		return b.String(), nil
	}
}

func generateGo(req HarRequest, client string) (string, error) {
	_ = client
	var b strings.Builder
	b.WriteString("package main\n\nimport (\n\t\"fmt\"\n\t\"io\"\n\t\"net/http\"\n\t\"strings\"\n)\n\n")
	b.WriteString("func main() {\n")
	if hasBody(req) {
		fmt.Fprintf(&b, "\tbody := strings.NewReader(`%s`)\n", bodyText(req))
		fmt.Fprintf(&b, "\treq, err := http.NewRequest(\"%s\", \"%s\", body)\n", req.Method, req.URL)
	} else {
		fmt.Fprintf(&b, "\treq, err := http.NewRequest(\"%s\", \"%s\", nil)\n", req.Method, req.URL)
	}
	b.WriteString("\tif err != nil {\n\t\tpanic(err)\n\t}\n")
	for _, h := range req.Headers {
		fmt.Fprintf(&b, "\treq.Header.Set(\"%s\", \"%s\")\n", h.Name, h.Value)
	}
	b.WriteString("\tresp, err := http.DefaultClient.Do(req)\n")
	b.WriteString("\tif err != nil {\n\t\tpanic(err)\n\t}\n")
	b.WriteString("\tdefer resp.Body.Close()\n")
	b.WriteString("\tb, _ := io.ReadAll(resp.Body)\n")
	b.WriteString("\tfmt.Println(string(b))\n}")
	return b.String(), nil
}

func generateJava(req HarRequest, client string) (string, error) {
	switch client {
	case "okhttp":
		var b strings.Builder
		b.WriteString("OkHttpClient client = new OkHttpClient();\n\n")
		if hasBody(req) {
			fmt.Fprintf(&b, "RequestBody body = RequestBody.create(\"%s\", MediaType.parse(\"%s\"));\n\n",
				shellEscape(bodyText(req)), mimeType(req))
			fmt.Fprintf(&b, "Request request = new Request.Builder()\n  .url(\"%s\")\n  .method(\"%s\", body)\n",
				req.URL, req.Method)
		} else {
			fmt.Fprintf(&b, "Request request = new Request.Builder()\n  .url(\"%s\")\n  .method(\"%s\", null)\n",
				req.URL, req.Method)
		}
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "  .addHeader(\"%s\", \"%s\")\n", h.Name, shellEscape(h.Value))
		}
		b.WriteString("  .build();\n\n")
		b.WriteString("try (Response response = client.newCall(request).execute()) {\n")
		b.WriteString("  System.out.println(response.body().string());\n}")
		return b.String(), nil
	case "unirest":
		var b strings.Builder
		b.WriteString("HttpResponse<String> response = Unirest.")
		b.WriteString(strings.ToLower(req.Method))
		b.WriteString("(\"")
		b.WriteString(req.URL)
		b.WriteString("\")")
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "\n  .header(\"%s\", \"%s\")", h.Name, shellEscape(h.Value))
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "\n  .body(\"%s\")", shellEscape(bodyText(req)))
		}
		b.WriteString("\n  .asString();\nSystem.out.println(response.getBody());")
		return b.String(), nil
	case "asynchttp":
		var b strings.Builder
		b.WriteString("AsyncHttpClient client = new DefaultAsyncHttpClient();\n")
		b.WriteString("client.prepare(\"")
		b.WriteString(req.Method)
		b.WriteString("\", \"")
		b.WriteString(req.URL)
		b.WriteString("\")")
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "\n  .setHeader(\"%s\", \"%s\")", h.Name, shellEscape(h.Value))
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "\n  .setBody(\"%s\")", shellEscape(bodyText(req)))
		}
		b.WriteString("\n  .execute()\n  .toCompletableFuture()\n  .thenAccept(r -> System.out.println(r.getResponseBody()));\n")
		b.WriteString("client.close();")
		return b.String(), nil
	case "restclient":
		var b strings.Builder
		b.WriteString("RestClient restClient = RestClient.create();\n")
		b.WriteString("String body = restClient.")
		b.WriteString(strings.ToLower(req.Method))
		b.WriteString("()\n  .uri(\"")
		b.WriteString(req.URL)
		b.WriteString("\")")
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "\n  .header(\"%s\", \"%s\")", h.Name, shellEscape(h.Value))
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "\n  .body(\"%s\")", shellEscape(bodyText(req)))
		}
		b.WriteString("\n  .retrieve()\n  .body(String.class);\nSystem.out.println(body);")
		return b.String(), nil
	case "nethttp":
		fallthrough
	default:
		var b strings.Builder
		fmt.Fprintf(&b, "URI uri = URI.create(\"%s\");\n", req.URL)
		fmt.Fprintf(&b, "HttpRequest.Builder builder = HttpRequest.newBuilder(uri).method(\"%s\", ", req.Method)
		if hasBody(req) {
			fmt.Fprintf(&b, "HttpRequest.BodyPublishers.ofString(\"%s\"))", shellEscape(bodyText(req)))
		} else {
			b.WriteString("HttpRequest.BodyPublishers.noBody())")
		}
		b.WriteString(";\n")
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "builder.header(\"%s\", \"%s\");\n", h.Name, shellEscape(h.Value))
		}
		b.WriteString("HttpResponse<String> response = HttpClient.newHttpClient()\n")
		b.WriteString("  .send(builder.build(), HttpResponse.BodyHandlers.ofString());\n")
		b.WriteString("System.out.println(response.body());")
		return b.String(), nil
	}
}

func mimeType(req HarRequest) string {
	if req.PostData != nil && req.PostData.MimeType != "" {
		return req.PostData.MimeType
	}
	return "text/plain"
}

func generatePHP(req HarRequest, client string) (string, error) {
	switch client {
	case "http1", "http2":
		var b strings.Builder
		b.WriteString("$client = new GuzzleHttp\\Client();\n")
		b.WriteString("$response = $client->request('")
		b.WriteString(req.Method)
		b.WriteString("', '")
		b.WriteString(shellEscape(req.URL))
		b.WriteString("', [\n")
		if len(req.Headers) > 0 {
			b.WriteString("    'headers' => [\n")
			for _, h := range req.Headers {
				fmt.Fprintf(&b, "        '%s' => '%s',\n", h.Name, shellEscape(h.Value))
			}
			b.WriteString("    ],\n")
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "    'body' => '%s',\n", shellEscape(bodyText(req)))
		}
		b.WriteString("]);\necho $response->getBody();")
		return b.String(), nil
	default:
		var b strings.Builder
		b.WriteString("$ch = curl_init();\n")
		fmt.Fprintf(&b, "curl_setopt($ch, CURLOPT_URL, '%s');\n", shellEscape(req.URL))
		fmt.Fprintf(&b, "curl_setopt($ch, CURLOPT_CUSTOMREQUEST, '%s');\n", req.Method)
		if len(req.Headers) > 0 {
			b.WriteString("curl_setopt($ch, CURLOPT_HTTPHEADER, [\n")
			for _, h := range req.Headers {
				fmt.Fprintf(&b, "    '%s: %s',\n", h.Name, shellEscape(h.Value))
			}
			b.WriteString("]);\n")
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "curl_setopt($ch, CURLOPT_POSTFIELDS, '%s');\n", shellEscape(bodyText(req)))
		}
		b.WriteString("curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);\n")
		b.WriteString("$response = curl_exec($ch);\ncurl_close($ch);\necho $response;")
		return b.String(), nil
	}
}

func generateRuby(req HarRequest, client string) (string, error) {
	if client == "faraday" {
		var b strings.Builder
		b.WriteString("require 'faraday'\n\n")
		b.WriteString("conn = Faraday.new\n")
		b.WriteString("response = conn.run_request(:")
		b.WriteString(strings.ToLower(req.Method))
		b.WriteString(", '")
		b.WriteString(shellEscape(req.URL))
		b.WriteString("', ")
		if hasBody(req) {
			fmt.Fprintf(&b, "'%s', ", shellEscape(bodyText(req)))
		} else {
			b.WriteString("nil, ")
		}
		b.WriteString("{\n")
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "  '%s' => '%s',\n", h.Name, shellEscape(h.Value))
		}
		b.WriteString("})\nputs response.body")
		return b.String(), nil
	}
	var b strings.Builder
	b.WriteString("require 'net/http'\nrequire 'uri'\n\n")
	fmt.Fprintf(&b, "uri = URI('%s')\n", shellEscape(req.URL))
	b.WriteString("http = Net::HTTP.new(uri.host, uri.port)\n")
	b.WriteString("http.use_ssl = (uri.scheme == 'https')\n")
	fmt.Fprintf(&b, "request = Net::HTTP::%s.new(uri)\n", strings.ToUpper(req.Method))
	for _, h := range req.Headers {
		fmt.Fprintf(&b, "request['%s'] = '%s'\n", h.Name, shellEscape(h.Value))
	}
	if hasBody(req) {
		fmt.Fprintf(&b, "request.body = '%s'\n", shellEscape(bodyText(req)))
	}
	b.WriteString("response = http.request(request)\nputs response.body")
	return b.String(), nil
}

func generateCSharp(req HarRequest, client string) (string, error) {
	if client == "restsharp" {
		var b strings.Builder
		b.WriteString("var client = new RestClient();\n")
		b.WriteString("var request = new RestRequest(\"")
		b.WriteString(req.URL)
		b.WriteString("\", Method.")
		b.WriteString(strings.ToUpper(req.Method))
		b.WriteString(");\n")
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "request.AddHeader(\"%s\", \"%s\");\n", h.Name, shellEscape(h.Value))
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "request.AddStringBody(\"%s\", DataFormat.Json);\n", shellEscape(bodyText(req)))
		}
		b.WriteString("var response = await client.ExecuteAsync(request);\nConsole.WriteLine(response.Content);")
		return b.String(), nil
	}
	var b strings.Builder
	b.WriteString("using var client = new HttpClient();\n")
	if len(req.Headers) > 0 {
		b.WriteString("var request = new HttpRequestMessage\n{\n")
		fmt.Fprintf(&b, "    Method = HttpMethod.%s,\n", strings.ToUpper(req.Method))
		fmt.Fprintf(&b, "    RequestUri = new Uri(\"%s\"),\n", req.URL)
		b.WriteString("};\n")
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "request.Headers.Add(\"%s\", \"%s\");\n", h.Name, shellEscape(h.Value))
		}
		if hasBody(req) {
			fmt.Fprintf(&b, "request.Content = new StringContent(\"%s\");\n", shellEscape(bodyText(req)))
		}
		b.WriteString("var response = await client.SendAsync(request);\n")
	} else {
		if hasBody(req) {
			fmt.Fprintf(&b, "var content = new StringContent(\"%s\");\n", shellEscape(bodyText(req)))
			fmt.Fprintf(&b, "var response = await client.%sAsync(\"%s\", content);\n",
				titleMethod(req.Method), req.URL)
		} else {
			fmt.Fprintf(&b, "var response = await client.%sAsync(\"%s\");\n",
				titleMethod(req.Method), req.URL)
		}
	}
	b.WriteString("Console.WriteLine(await response.Content.ReadAsStringAsync());")
	return b.String(), nil
}

func titleMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS":
		return strings.ToUpper(method[:1]) + strings.ToLower(method[1:])
	default:
		return "Send"
	}
}

func generateSwift(req HarRequest, client string) (string, error) {
	_ = client
	var b strings.Builder
	fmt.Fprintf(&b, "let url = URL(string: \"%s\")!\n", req.URL)
	b.WriteString("var request = URLRequest(url: url)\n")
	fmt.Fprintf(&b, "request.httpMethod = \"%s\"\n", req.Method)
	for _, h := range req.Headers {
		fmt.Fprintf(&b, "request.setValue(\"%s\", forHTTPHeaderField: \"%s\")\n", shellEscape(h.Value), h.Name)
	}
	if hasBody(req) {
		fmt.Fprintf(&b, "request.httpBody = \"%s\".data(using: .utf8)\n", shellEscape(bodyText(req)))
	}
	b.WriteString("let task = URLSession.shared.dataTask(with: request) { data, _, _ in\n")
	b.WriteString("    print(String(data: data!, encoding: .utf8)!)\n}\ntask.resume()")
	return b.String(), nil
}

func generateKotlin(req HarRequest, client string) (string, error) {
	_ = client
	var b strings.Builder
	b.WriteString("val client = OkHttpClient()\n\n")
	if hasBody(req) {
		fmt.Fprintf(&b, "val body = \"%s\".toRequestBody(\"%s\".toMediaType())\n\n",
			shellEscape(bodyText(req)), mimeType(req))
		fmt.Fprintf(&b, "val request = Request.Builder().url(\"%s\").method(\"%s\", body)\n", req.URL, req.Method)
	} else {
		fmt.Fprintf(&b, "val request = Request.Builder().url(\"%s\").method(\"%s\", null)\n", req.URL, req.Method)
	}
	for _, h := range req.Headers {
		fmt.Fprintf(&b, "    .addHeader(\"%s\", \"%s\")\n", h.Name, shellEscape(h.Value))
	}
	b.WriteString("    .build()\n\n")
	b.WriteString("client.newCall(request).execute().use { response ->\n")
	b.WriteString("    println(response.body?.string())\n}")
	return b.String(), nil
}

func generateC(req HarRequest, client string) (string, error) {
	switch client {
	case "libcurl":
		var b strings.Builder
		b.WriteString("#include <stdio.h>\n#include <curl/curl.h>\n\n")
		b.WriteString("int main(void) {\n")
		b.WriteString("  CURL *curl = curl_easy_init();\n")
		fmt.Fprintf(&b, "  curl_easy_setopt(curl, CURLOPT_URL, \"%s\");\n", req.URL)
		fmt.Fprintf(&b, "  curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, \"%s\");\n", req.Method)
		if hasBody(req) {
			fmt.Fprintf(&b, "  curl_easy_setopt(curl, CURLOPT_POSTFIELDS, \"%s\");\n", shellEscape(bodyText(req)))
		}
		b.WriteString("  curl_easy_perform(curl);\n")
		b.WriteString("  curl_easy_cleanup(curl);\n")
		b.WriteString("  return 0;\n}")
		return b.String(), nil
	default:
		return buildCurl(req), nil
	}
}

func generateObjC(req HarRequest, client string) (string, error) {
	_ = client
	var b strings.Builder
	fmt.Fprintf(&b, "NSURL *url = [NSURL URLWithString:@\"%s\"];\n", req.URL)
	b.WriteString("NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:url];\n")
	fmt.Fprintf(&b, "[request setHTTPMethod:@\"%s\"];\n", req.Method)
	for _, h := range req.Headers {
		fmt.Fprintf(&b, "[request setValue:@\"%s\" forHTTPHeaderField:@\"%s\"];\n", shellEscape(h.Value), h.Name)
	}
	if hasBody(req) {
		fmt.Fprintf(&b, "[request setHTTPBody:[@\"%s\" dataUsingEncoding:NSUTF8StringEncoding]];\n", shellEscape(bodyText(req)))
	}
	b.WriteString("NSURLSessionDataTask *task = [[NSURLSession sharedSession] dataTaskWithRequest:request completionHandler:^(NSData *data, NSURLResponse *response, NSError *error) {\n")
	b.WriteString("    NSLog(@\"%%@\", [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding]);\n}];\n[task resume];")
	return b.String(), nil
}

func generateR(req HarRequest, client string) (string, error) {
	_ = client
	var b strings.Builder
	b.WriteString("library(httr)\n\n")
	fmt.Fprintf(&b, "response <- verb(\"%s\", \"%s\"", strings.ToUpper(req.Method), req.URL)
	if len(req.Headers) > 0 {
		b.WriteString(", add_headers = c(")
		parts := make([]string, 0, len(req.Headers))
		for _, h := range req.Headers {
			parts = append(parts, fmt.Sprintf("\"%s\" = \"%s\"", h.Name, shellEscape(h.Value)))
		}
		b.WriteString(strings.Join(parts, ", "))
		b.WriteString(")")
	}
	if hasBody(req) {
		fmt.Fprintf(&b, ", body = '%s'", shellEscape(bodyText(req)))
	}
	b.WriteString(")\ncontent(response, \"text\")")
	return b.String(), nil
}

func generateClojure(req HarRequest, client string) (string, error) {
	_ = client
	var b strings.Builder
	b.WriteString("(require '[clj-http.client :as client])\n\n")
	b.WriteString("(println (:body (client/request {")
	fmt.Fprintf(&b, ":method :%s", strings.ToLower(req.Method))
	fmt.Fprintf(&b, " :url \"%s\"", req.URL)
	if len(req.Headers) > 0 {
		b.WriteString(" :headers {")
		for i, h := range req.Headers {
			if i > 0 {
				b.WriteString(" ")
			}
			fmt.Fprintf(&b, "\"%s\" \"%s\"", h.Name, shellEscape(h.Value))
		}
		b.WriteString("}")
	}
	if hasBody(req) {
		fmt.Fprintf(&b, " :body \"%s\"", shellEscape(bodyText(req)))
	}
	b.WriteString("})))")
	return b.String(), nil
}

func generateOCaml(req HarRequest, client string) (string, error) {
	_ = client
	var b strings.Builder
	b.WriteString("open Cohttp_lwt_unix\nopen Cohttp\nopen Lwt\n\n")
	b.WriteString("let uri = Uri.of_string \"")
	b.WriteString(req.URL)
	b.WriteString("\" in\n")
	fmt.Fprintf(&b, "let body = Cohttp_lwt.Body.of_string \"%s\" in\n", shellEscape(bodyText(req)))
	fmt.Fprintf(&b, "Client.call %s uri body >>= fun (_, body) ->\n", strings.ToUpper(req.Method))
	b.WriteString("  body |> Cohttp_lwt.Body.to_string >|= print_endline")
	return b.String(), nil
}

func generateCrystal(req HarRequest, client string) (string, error) {
	_ = client
	var b strings.Builder
	b.WriteString("require \"http/client\"\n\n")
	b.WriteString("uri = URI.parse(\"")
	b.WriteString(req.URL)
	b.WriteString("\")\n")
	b.WriteString("client = HTTP::Client.new(uri.host, uri.port, tls: uri.scheme == \"https\")\n")
	b.WriteString("response = client.")
	b.WriteString(strings.ToLower(req.Method))
	b.WriteString("(uri.request_target")
	if hasBody(req) {
		fmt.Fprintf(&b, ", headers: HTTP::Headers{\n")
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "  \"%s\" => \"%s\",\n", h.Name, shellEscape(h.Value))
		}
		b.WriteString("}, body: \"")
		b.WriteString(shellEscape(bodyText(req)))
		b.WriteString("\")")
	} else if len(req.Headers) > 0 {
		b.WriteString(", headers: HTTP::Headers{\n")
		for _, h := range req.Headers {
			fmt.Fprintf(&b, "  \"%s\" => \"%s\",\n", h.Name, shellEscape(h.Value))
		}
		b.WriteString("})")
	} else {
		b.WriteString(")")
	}
	b.WriteString("\nputs response.body")
	return b.String(), nil
}

func generateRust(req HarRequest, client string) (string, error) {
	_ = client
	var b strings.Builder
	b.WriteString("#[tokio::main]\nasync fn main() -> Result<(), Box<dyn std::error::Error>> {\n")
	b.WriteString("    let client = reqwest::Client::new();\n")
	b.WriteString("    let response = client\n        .")
	b.WriteString(strings.ToLower(req.Method))
	b.WriteString("( \"")
	b.WriteString(req.URL)
	b.WriteString("\" )")
	for _, h := range req.Headers {
		fmt.Fprintf(&b, "\n        .header(\"%s\", \"%s\")", h.Name, shellEscape(h.Value))
	}
	if hasBody(req) {
		fmt.Fprintf(&b, "\n        .body(\"%s\")", shellEscape(bodyText(req)))
	}
	b.WriteString("\n        .send()\n        .await?;\n")
	b.WriteString("    println!(\"{}\", response.text().await?);\n    Ok(())\n}")
	return b.String(), nil
}

func generatePowerShell(req HarRequest, client string) (string, error) {
	cmd := "Invoke-WebRequest"
	if client == "restmethod" {
		cmd = "Invoke-RestMethod"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "$response = %s -Uri '%s' -Method %s", cmd, shellEscape(req.URL), req.Method)
	if len(req.Headers) > 0 {
		b.WriteString(" -Headers @{")
		for i, h := range req.Headers {
			if i > 0 {
				b.WriteString("; ")
			}
			fmt.Fprintf(&b, "'%s' = '%s'", h.Name, shellEscape(h.Value))
		}
		b.WriteString("}")
	}
	if hasBody(req) {
		fmt.Fprintf(&b, " -Body '%s'", shellEscape(bodyText(req)))
	}
	if client == "restmethod" {
		b.WriteString("\n$response")
	} else {
		b.WriteString("\n$response.Content")
	}
	return b.String(), nil
}
