package utils

import "regexp"

func IsValidURL(urlPath string) bool{
	var apiValidPath = regexp.MustCompile("^/api/(posts|users|auth)/[a-z]*$")

	matches := apiValidPath.FindStringSubmatch(urlPath)

	return matches != nil
}