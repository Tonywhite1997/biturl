package helper

import "net/url"

func IsValidURL(urlString string) bool {
	parsedURL, err := url.ParseRequestURI(urlString)
	if err != nil {
		return false
	}
	// Ensure the URL has a scheme and host
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}
