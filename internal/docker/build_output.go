package docker

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
)

var (
	buildStepLineRe = regexp.MustCompile(`(?i)^step\s+(\d+)/(\d+)\s*:`)
	buildEndStepRe  = regexp.MustCompile(`^ ---\> [a-f0-9]+\n$`)
)

type buildOutputState struct {
	hide          string // "none", "all", "until-next"
	totalSteps    int
	step          int
	extraCommands int
}

// TransformBuildOutput simplifies injected build steps in Docker's NDJSON stream.
func TransformBuildOutput(raw []byte, extraCommands int) []byte {
	state := buildOutputState{hide: "none", totalSteps: -1, extraCommands: extraCommands}
	var out bytes.Buffer
	for _, line := range bytes.Split(raw, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		transformed := transformBuildLine(line, &state)
		if len(transformed) > 0 {
			out.Write(transformed)
			out.WriteByte('\n')
		}
	}
	return out.Bytes()
}

func transformBuildLine(line []byte, state *buildOutputState) []byte {
	var data struct {
		Stream string `json:"stream"`
	}
	if err := json.Unmarshal(line, &data); err != nil || data.Stream == "" {
		return line
	}

	if state.totalSteps == -1 {
		if m := buildStepLineRe.FindStringSubmatch(data.Stream); len(m) == 3 {
			total := atoiSafe(m[2]) - state.extraCommands
			if total > 0 {
				state.totalSteps = total
			}
		}
	}

	if strings.Contains(data.Stream, "LABEL "+BuildLabel+"=started") {
		state.hide = "all"
		return []byte(`{"stream":" *** Enabling HTTP Toolkit interception ***\n"}`)
	}

	if state.hide == "all" {
		if strings.Contains(data.Stream, "LABEL "+BuildLabel+"=") {
			state.hide = "until-next"
		}
		return nil
	}

	if state.hide == "until-next" {
		if !buildEndStepRe.MatchString(data.Stream) {
			return nil
		}
		state.hide = "none"
	}

	if m := buildStepLineRe.FindStringSubmatch(data.Stream); len(m) == 3 {
		state.step++
		data.Stream = replaceBuildStepLine(data.Stream, state.step, state.totalSteps)
	}
	out, _ := json.Marshal(data)
	return out
}

func replaceBuildStepLine(line string, step, total int) string {
	if total <= 0 {
		return line
	}
	return buildStepLineRe.ReplaceAllString(line, fmtStep(step, total))
}

func fmtStep(step, total int) string {
	return "Step " + itoa(step) + "/" + itoa(total) + " :"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

func atoiSafe(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}
