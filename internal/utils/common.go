package utils

import "golang.org/x/exp/constraints"

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
