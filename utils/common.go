package utils

import (
	"bytes"
	"os/exec"

	"golang.org/x/exp/constraints"
)

// FixRange 矫正数值, 必须大于等于下限，小于等于上限
func FixRange[T constraints.Integer | constraints.Float](arr ...T) T {
	var empty T
	if len(arr) == 0 {
		return empty
	} else if len(arr) == 1 {
		return arr[0]
	} else if len(arr) >= 2 && arr[0] < arr[1] {
		return arr[1]
	} else if len(arr) >= 3 && arr[0] > arr[2] {
		return arr[2]
	}

	return arr[0]
}

func ExecCmd(dir, cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	command.Dir = dir
	// 给标准输入以及标准错误初始化一个buffer，每条命令的输出位置可能是不一样的，
	// 比如有的命令会将输出放到stdout，有的放到stderr
	command.Stdout = &bytes.Buffer{}
	command.Stderr = &bytes.Buffer{}

	err := command.Run()
	if err != nil {
		// 打印程序中的错误以及命令行标准错误中的输出
		return command.Stderr.(*bytes.Buffer).String(), err
	}
	// 打印命令行的标准输出
	return command.Stdout.(*bytes.Buffer).String(), nil
}
