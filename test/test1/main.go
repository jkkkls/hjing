package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"

	"github.com/gookit/goutil/dump"
	"github.com/jkkkls/hjing/cmds/xlsx2proto/examples"
)

func isValidDiomainName(appName string) bool {
	ok, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9/.-]+$`, appName)
	return ok
}

func main() {
	buff, _ := os.ReadFile("../../cmds/xlsx2proto/examples/config.json")

	x := &examples.Configs{}

	err := json.Unmarshal(buff, x)
	log.Println(err)
	dump.Println(x)
}
