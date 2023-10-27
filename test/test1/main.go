package main

import (
	"github.com/jkkkls/hjing/utils"
	"log"
)

func main() {
	str, err := utils.ExecCmd("./aaa", "ls")
	log.Println(str, err)
}
