package layout

import (
	"bytes"
	"embed"
	"os"
)

//go:embed app
var appStatic embed.FS

func CopyFile(fileName, dstFile string, args ...string) error {
	buff, err := appStatic.ReadFile(fileName)
	if err != nil {
		return err
	}
	n := len(args) / 2
	for i := 0; i < n; i++ {
		buff = bytes.ReplaceAll(buff, []byte(args[i*2]), []byte(args[i*2+1]))
	}

	return os.WriteFile(dstFile, buff, 0644)
}
