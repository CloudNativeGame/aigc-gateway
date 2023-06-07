package utils

import "strings"

func GetDomainFromEndpoint(endpoint string) string {
	var domain string
	tail := 0
	index := strings.Index(endpoint, ".")
	if endpoint[len(endpoint)-1] == '/' {
		tail = 1
	}
	if index != -1 {
		domain = endpoint[index : len(endpoint)-tail]
	}
	return domain
}
