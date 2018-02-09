package util

import "strings"

// SliceContains checks if a string exists in a string slice
func SliceContains(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(b, a) {
			return true
		}
	}
	return false
}

// TrimSuffix removes a specified trailing substring from the string
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
