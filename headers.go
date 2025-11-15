package main

import (
	"slices"
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

// Set overwrites the current key-value pair with the single value provided.
func (h Headers) Set(key string, value string) {
	h[canonical(key)] = []string{value}
}

// Add adds the key-value pair to the header.
func (h Headers) Add(key string, value string) {
	h[canonical(key)] = append(h[canonical(key)], value)
}

// ContainsKey returns true if a key is defined in the Headers.
func (h Headers) ContainsKey(key string) bool {
	_, ok := h[canonical(key)]
	return ok
}

// ContainsValue returns true if a specific key-value pair is contained by the Headers.
func (h Headers) ContainsValue(key string, value string) bool {
	return slices.Contains(h[canonical(key)], value)
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
