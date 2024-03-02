package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jkkkls/hjing/cmds/xlsx2proto/sync_conf"
	"github.com/jkkkls/hjing/utils"
)

type Field struct {
	Name  string
	Array int
	Type  string
	Tag   string
}

type Struct struct {
	Name   string
	Fields []*Field
}

var (
	outDir         = flag.String("out", "", "输出目录")
	inDir          = flag.String("in", "", "输入目录")
	data           = flag.String("data", "", "数据格式, json or package name")
	includeComment = flag.Bool("includeComment", false, "是否导出注释")
)

func main() {
	flag.Parse()

	if *data == "" || *outDir == "" || *inDir == "" {
		flag.Usage()
		return
	}

	if *data == "json" {
		exportJson()
		return
	}

	arr := strings.Split(*inDir, ";")
	buff, err := sync_conf.ExportStruct2(arr, *data)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.WriteFile(*outDir, buff, 0o644)

	newBuff, err := utils.ExecCmd("", "gofmt", *outDir)
	if err != nil {
		color.Red(err.Error())
		return
	}
	os.WriteFile(*outDir, []byte(newBuff), 0o644)
}

func exportJson() {
	arr := strings.Split(*inDir, ";")
	buff, err := sync_conf.ExportJson(arr, *includeComment)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.WriteFile(*outDir, buff, 0o644)
}
