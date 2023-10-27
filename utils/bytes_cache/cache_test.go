package bytes_cache

import "testing"

func Test_getIndex(t *testing.T) {
	var data = [][]int{
		{0, 0},
		{1, 0},
		{2, 1},
		{3, 2},
		{4, 2},
	}
	for _, v := range data {
		if i := getIndex(v[0]); i != v[1] {
			t.Errorf("getIndex, v: %v, i: %v", v, i)
		}
	}
}
