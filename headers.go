package main

import (
	"strings"
)

type Headers map[string][]string

// Get returns the first value of a key and an `ok` bool if it exists in the map.
func (h Headers) Get(key string) (string, bool) {
	if len(h[canonical(key)]) == 0 {
		return "", false
	}
	return h[canonical(key)][0], true
}

// GetAll returns all key-value pairs for the specific key.
func (h Headers) GetAll(key string) []string {
	return h[canonical(key)]
}

// Add adds the key-value pair to the header.
func (h Headers) Add(key string, value string) {
	h[canonical(key)] = append(h[canonical(key)], value)
}

// ForEach iterates over all (key, value) pairs, including duplicates.
func (h Headers) ForEach(f func(key string, value string)) {
	for key, values := range h {
		for _, value := range values {
			f(key, value)
		}
	}
}

// canonical returns a canonical form of a key string.
// Header keys can be case insensitive.
func canonical(key string) string {
	return strings.ToLower(key)
}
