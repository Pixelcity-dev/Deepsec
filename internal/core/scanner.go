package core

import (
	"context"
)

type Scanner interface {
	Name() string
	Type() ScanType
	Scan(ctx context.Context, target Target, rules []Rule) ([]Finding, error)
	SupportedTargets() []TargetKind
	IsEnabled() bool
	SetEnabled(enabled bool)
}

type BaseScanner struct {
	Enabled bool
}

func (b *BaseScanner) IsEnabled() bool {
	return b.Enabled
}

func (b *BaseScanner) SetEnabled(enabled bool) {
	b.Enabled = enabled
}

type ScannerRegistry struct {
	scanners map[ScanType]Scanner
}

func NewScannerRegistry() *ScannerRegistry {
	return &ScannerRegistry{
		scanners: make(map[ScanType]Scanner),
	}
}

func (r *ScannerRegistry) Register(s Scanner) {
	r.scanners[s.Type()] = s
}

func (r *ScannerRegistry) Get(scanType ScanType) (Scanner, bool) {
	s, ok := r.scanners[scanType]
	return s, ok
}

func (r *ScannerRegistry) GetAll() []Scanner {
	var scanners []Scanner
	for _, s := range r.scanners {
		scanners = append(scanners, s)
	}
	return scanners
}

func (r *ScannerRegistry) GetByTypes(types []ScanType) []Scanner {
	var scanners []Scanner
	for _, t := range types {
		if s, ok := r.scanners[t]; ok && s.IsEnabled() {
			scanners = append(scanners, s)
		}
	}
	return scanners
}

func (r *ScannerRegistry) SupportedTypes() []ScanType {
	var types []ScanType
	for t := range r.scanners {
		types = append(types, t)
	}
	return types
}
