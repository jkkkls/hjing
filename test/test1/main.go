package main

import (
	"log"
	"regexp"
)

func isValidDiomainName(appName string) bool {
	ok, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9/.-]+$`, appName)
	return ok
}

func main() {
	str := "qqqq.co-m/123/aa"
	log.Println(isValidDiomainName(str))
}
