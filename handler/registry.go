package handler

import (
	"fmt"
	"sort"
	"strings"
)

// registry stores registered handlers indexed by subject.
var registry = make(map[string]Handler)

// Register adds a handler to the registry.
// It is typically called from init() functions in handler packages.
func Register(h Handler) {
	registry[h.Subject()] = h
}

// Lookup returns the handler for the given subject, or an error if not found.
func Lookup(subject string) (Handler, error) {
	h, ok := registry[subject]
	if !ok {
		return nil, fmt.Errorf("unsupported subject %q; supported subjects: %s", subject, supportedSubjects())
	}
	return h, nil
}

// SupportedSubjects returns a sorted list of all registered subject strings.
func SupportedSubjects() []string {
	subjects := make([]string, 0, len(registry))
	for s := range registry {
		subjects = append(subjects, s)
	}
	sort.Strings(subjects)
	return subjects
}

func supportedSubjects() string {
	return strings.Join(SupportedSubjects(), ", ")
}
