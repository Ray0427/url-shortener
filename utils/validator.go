package utils

import "net/url"

func CheckUrl(fullUrl string) bool {
	_, err := url.ParseRequestURI(fullUrl)
	if err != nil {
		return false
	}
	return true
}
