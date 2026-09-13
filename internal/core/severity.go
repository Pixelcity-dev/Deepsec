package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Severity int

const (
	SeverityInfo Severity = iota
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

var severityNames = map[Severity]string{
	SeverityInfo:     "INFO",
	SeverityLow:      "LOW",
	SeverityMedium:   "MEDIUM",
	SeverityHigh:     "HIGH",
	SeverityCritical: "CRITICAL",
}

var severityColors = map[Severity]string{
	SeverityInfo:     "\033[36m",
	SeverityLow:      "\033[33m",
	SeverityMedium:   "\033[93m",
	SeverityHigh:     "\033[31m",
	SeverityCritical: "\033[1;31m",
}

func (s Severity) String() string {
	if name, ok := severityNames[s]; ok {
		return name
	}
	return "UNKNOWN"
}

func (s Severity) Color() string {
	if color, ok := severityColors[s]; ok {
		return color
	}
	return ""
}

func (s Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Severity) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = ParseSeverity(str)
	return nil
}

func ParseSeverity(s string) Severity {
	switch strings.ToUpper(s) {
	case "INFO", "INFORMATIONAL":
		return SeverityInfo
	case "LOW":
		return SeverityLow
	case "MEDIUM", "MED":
		return SeverityMedium
	case "HIGH":
		return SeverityHigh
	case "CRITICAL", "CRIT":
		return SeverityCritical
	default:
		return SeverityInfo
	}
}

func ParseSeverities(s string) []Severity {
	parts := strings.Split(s, ",")
	var severities []Severity
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			severities = append(severities, ParseSeverity(p))
		}
	}
	return severities
}

func (s Severity) Weight() int {
	return int(s)
}

func FormatSeverity(s Severity) string {
	color := s.Color()
	reset := "\033[0m"
	if color == "" {
		return fmt.Sprintf("%-10s", s.String())
	}
	return fmt.Sprintf("%s%-10s%s", color, s.String(), reset)
}
