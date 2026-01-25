package commands

import (
	"regexp"
	"strings"
)

// Version parsers - shared between status and install

// stripAnsi removes ANSI escape codes (colors, formatting) from a string
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripAnsi(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

func trimVersion(s string) string {
	return strings.TrimSpace(strings.TrimPrefix(stripAnsi(s), "v"))
}

func parseBrewVersion(s string) string {
	s = stripAnsi(s)
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		parts := strings.Split(lines[0], " ")
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

func parseGitVersion(s string) string {
	s = stripAnsi(s)
	v := strings.TrimPrefix(strings.TrimSpace(s), "git version ")
	// Truncate Apple Git suffix for cleaner display
	if idx := strings.Index(v, " ("); idx != -1 {
		v = v[:idx]
	}
	return v
}

func parseGoVersion(s string) string {
	s = stripAnsi(s)
	// "go version go1.23.4 darwin/arm64" -> "1.23.4"
	parts := strings.Fields(s)
	if len(parts) >= 3 {
		return strings.TrimPrefix(parts[2], "go")
	}
	return ""
}

func parseJavaVersion(s string) string {
	s = stripAnsi(s)
	// java -version outputs to stderr: openjdk version "21.0.1" ...
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(line, "version") {
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start != -1 && end != -1 && end > start {
				return line[start+1 : end]
			}
		}
	}
	return ""
}

func parsePythonVersion(s string) string {
	s = stripAnsi(s)
	return strings.TrimPrefix(strings.TrimSpace(s), "Python ")
}

func parseTerraformVersion(s string) string {
	s = stripAnsi(s)
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		parts := strings.Split(lines[0], " ")
		if len(parts) >= 2 {
			return strings.TrimPrefix(parts[1], "v")
		}
	}
	return ""
}

func parseAnsibleVersion(s string) string {
	s = stripAnsi(s)
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		line := lines[0]
		if strings.Contains(line, "[core") {
			start := strings.Index(line, "[core ")
			end := strings.Index(line, "]")
			if start != -1 && end != -1 {
				return line[start+6 : end]
			}
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

func parseMultipassVersion(s string) string {
	s = stripAnsi(s)
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

func parseCodexVersion(s string) string {
	s = stripAnsi(s)
	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) >= 2 {
		return parts[1]
	}
	return strings.TrimSpace(s)
}

func parseMoleVersion(s string) string {
	s = stripAnsi(s)
	// "\nMole version 1.14.5\nmacOS: ..." -> "1.14.5"
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, "Mole version") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				return parts[2]
			}
		}
	}
	return ""
}

func parseClaudeVersion(s string) string {
	s = stripAnsi(s)
	// "2.0.76 (Claude Code)" -> "2.0.76"
	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) >= 1 {
		return parts[0]
	}
	return ""
}

func parseAnsibleLintVersion(s string) string {
	s = stripAnsi(s)
	// "ansible-lint 25.12.2 using ansible-core:2.20.1 ..." -> "25.12.2"
	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

func parsePulumiVersion(s string) string {
	s = stripAnsi(s)
	// "v3.100.0" -> "3.100.0"
	return strings.TrimPrefix(strings.TrimSpace(s), "v")
}

func parseHappyCoderVersion(s string) string {
	s = stripAnsi(s)
	// "happy version: 0.13.0\nUsing Claude Code..." -> "0.13.0"
	lines := strings.Split(s, "\n")
	if len(lines) > 0 {
		line := lines[0]
		if strings.HasPrefix(line, "happy version:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "happy version:"))
		}
	}
	return ""
}
