package main

import "strings"

type Headers map[string][]string

// Get returns the first value of a key, and an empty string for no values.
func (h Headers) Get(key string) string {
	if len(h[canonical(key)]) == 0 {
		return ""
	}
	return h[canonical(key)][0]
}

// Add adds the key-value pair to the header.
func (h Headers) Add(key string, value string) {
	h[canonical(key)] = append(h[canonical(key)], value)
}

// Values returns all key-value pairs for the specific key.
func (h Headers) Values(key string) []string {
	return h[canonical(key)]
}

func canonical(key string) string {
	return strings.ToLower(key)
}
