package utils

import "strings"

func URLWithQueryParams(baseURL string, params map[string]string) string {
	if len(params) == 0 {
		return baseURL
	}

	var url strings.Builder
	url.WriteString(baseURL + "?")
	first := true
	for key, value := range params {
		if !first {
			url.WriteString("&")
		}
		if value == "" {
			continue
		}
		url.WriteString(key + "=" + value)
		first = false
	}
	return url.String()
}
