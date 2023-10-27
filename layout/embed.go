package layout

import (
	"bytes"
	"embed"
	"github.com/jkkkls/hjing/utils"
	"os"
)

//go:embed app
//go:embed project
var appStatic embed.FS

// CopyFile 复制embed中的指定文件
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

// CopyDir 递归复制embed中的指定目录
func CopyDir(srcDir, dstDir string, args ...string) error {
	dir, err := appStatic.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, v := range dir {
		if v.IsDir() {
			newDir := dstDir + "/" + v.Name()
			if !utils.PathExists(newDir) {
				err := os.MkdirAll(newDir, 0755)
				if err != nil {
					return err
				}
			}
			err := CopyDir(srcDir+"/"+v.Name(), dstDir+"/"+v.Name(), args...)
			if err != nil {
				return err
			}
			continue
		}
		err := CopyFile(srcDir+"/"+v.Name(), dstDir+"/"+v.Name(), args...)
		if err != nil {
			return err
		}
	}
	return nil
}
