package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	ID          string   `json:"id" yaml:"id"`
	Name        string   `json:"name" yaml:"name"`
	Severity    Severity `json:"severity" yaml:"severity"`
	Category    string   `json:"category" yaml:"category"`
	Language    string   `json:"language,omitempty" yaml:"language,omitempty"`
	ScannerType ScanType `json:"scanner_type" yaml:"scanner_type"`
	Pattern     string   `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Description string   `json:"description" yaml:"description"`
	Fix         string   `json:"fix,omitempty" yaml:"fix,omitempty"`
	References  []string `json:"references,omitempty" yaml:"references,omitempty"`
	Tags        []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Enabled     bool     `json:"enabled" yaml:"enabled"`
}

type RuleSet struct {
	Rules []Rule `json:"rules" yaml:"rules"`
}

type RuleEngine struct {
	rules     map[ScanType][]Rule
	overrides map[string]bool
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		rules:     make(map[ScanType][]Rule),
		overrides: make(map[string]bool),
	}
}

func (e *RuleEngine) LoadRulesFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	ext := filepath.Ext(path)
	var ruleSet RuleSet

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &ruleSet); err != nil {
			return err
		}
	case ".json":
		if err := json.Unmarshal(data, &ruleSet); err != nil {
			return err
		}
	}

	for _, rule := range ruleSet.Rules {
		e.AddRule(rule)
	}

	return nil
}

func (e *RuleEngine) LoadRulesFromDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ".yaml" || ext == ".yml" || ext == ".json" {
			return e.LoadRulesFromFile(path)
		}
		return nil
	})
}

func (e *RuleEngine) AddRule(rule Rule) {
	if rule.ID == "" {
		return
	}
	e.rules[rule.ScannerType] = append(e.rules[rule.ScannerType], rule)
}

func (e *RuleEngine) GetRulesForScanner(scanType ScanType) []Rule {
	rules := e.rules[scanType]
	var enabled []Rule
	for _, r := range rules {
		if r.Enabled && !e.overrides[r.ID] {
			enabled = append(enabled, r)
		}
	}
	return enabled
}

func (e *RuleEngine) GetAllRules() []Rule {
	var all []Rule
	for _, rules := range e.rules {
		all = append(all, rules...)
	}
	return all
}

func (e *RuleEngine) GetRuleByID(id string) *Rule {
	for _, rules := range e.rules {
		for _, r := range rules {
			if r.ID == id {
				return &r
			}
		}
	}
	return nil
}

func (e *RuleEngine) DisableRule(id string) {
	e.overrides[id] = true
}

func (e *RuleEngine) EnableRule(id string) {
	delete(e.overrides, id)
}

func (e *RuleEngine) SearchRules(query string) []Rule {
	query = strings.ToLower(query)
	var matches []Rule
	for _, rules := range e.rules {
		for _, r := range rules {
			if strings.Contains(strings.ToLower(r.Name), query) ||
				strings.Contains(strings.ToLower(r.Description), query) ||
				strings.Contains(strings.ToLower(r.ID), query) ||
				strings.Contains(strings.ToLower(r.Category), query) {
				matches = append(matches, r)
			}
			for _, tag := range r.Tags {
				if strings.Contains(strings.ToLower(tag), query) {
					matches = append(matches, r)
					break
				}
			}
		}
	}
	return matches
}

func (e *RuleEngine) GetRulesByLanguage(lang string) []Rule {
	var matches []Rule
	for _, rules := range e.rules {
		for _, r := range rules {
			if strings.EqualFold(r.Language, lang) || r.Language == "" {
				matches = append(matches, r)
			}
		}
	}
	return matches
}

func (e *RuleEngine) GetRulesByCategory(category string) []Rule {
	var matches []Rule
	for _, rules := range e.rules {
		for _, r := range rules {
			if strings.EqualFold(r.Category, category) {
				matches = append(matches, r)
			}
		}
	}
	return matches
}

func (e *RuleEngine) Stats() map[string]int {
	stats := make(map[string]int)
	for scanType, rules := range e.rules {
		stats[string(scanType)] = len(rules)
	}
	return stats
}
