package sync_conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
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

type Structs []*Struct

// sort.Sort(ByCount(data))
func (a Structs) Len() int      { return len(a) }
func (a Structs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Structs) Less(i, j int) bool {
	if a[i].Name == "Public" {
		return true
	} else if a[j].Name == "Public" {
		return false
	}

	return a[i].Name < a[j].Name
}

var (
	types = map[jsonparser.ValueType]string{
		jsonparser.String:  "string",
		jsonparser.Number:  "int64",
		jsonparser.Boolean: "bool",
	}
	kvp      = "KVP"
	comments map[string]string
	defaults map[string]string
	structs  map[string]*Struct
)

func getDefault(key string) string {
	if strings.HasPrefix(key, "Public") {
		return defaults[key[6:]]
	}
	return defaults[key]
}

func getComment(key string) string {
	if strings.HasPrefix(key, "Public") {
		return comments[key[6:]]
	}
	return comments[key]
}

func extractGormTag(s string) (string, string) {
	r := regexp.MustCompile(`\(.*\)`)
	str := r.FindString(string(s))

	index := strings.Index(s, "(")
	return str[1 : len(str)-1], s[:index]
}

func ExportJson(inDirs []string, includeComment bool) ([]byte, error) {
	comments = map[string]string{}
	j := make(map[string]interface{})
	cnt := 0

	for _, v := range inDirs {
		files, err := os.ReadDir(v)
		if err != nil {
			continue
		}

		// 开始导出

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if !strings.HasSuffix(f.Name(), ".xlsx") {
				continue
			}
			name, comment := extractGormTag(f.Name())
			r, err := exportXLSX2Json(v+"/"+f.Name(), name)
			if err != nil {
				return nil, err
			}

			j[name] = r
			a, _ := checkName(name, true)
			comments[a] = comment
			cnt++
		}
	}

	// buff, _ := json.MarshalIndent(j, "", "\t")
	if includeComment {
		j["comments"] = comments
	}
	return json.Marshal(j)
}

func exportCellJson(t, value string) (interface{}, error) {
	switch t {
	case "str":
		return value, nil
	case "int":
		if value == "" {
			return 0, nil
		}
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("非法的数值")
		}
		return int64(n), nil
	case "[]int":
		var m []int64
		arr := strings.Split(value, ";")
		for _, v := range arr {
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("非法的数值: %v", v)
			}
			m = append(m, int64(n))
		}
		return m, nil
	case "[][]int":
		var m [][]int64
		arr := strings.Split(value, "#")
		for _, v := range arr {
			var m1 []int64

			arr1 := strings.Split(v, ";")
			for _, v1 := range arr1 {
				n, err := strconv.ParseInt(v1, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("非法的数值: %v", v1)
				}
				m1 = append(m1, int64(n))
			}

			m = append(m, m1)
		}
		return m, nil
	case "[]str":
		m := strings.Split(value, ";")
		return m, nil
	case "[][]str":
		var m [][]string
		arr := strings.Split(value, "#")
		for _, v := range arr {
			m = append(m, strings.Split(v, ";"))
		}
		return m, nil
	}

	return nil, fmt.Errorf("不支持的数据类型[%v]", t)
}

type FieldType struct {
	Type string
	Name string
}

func exportSheet2Json(file, sname string, sheetIndex int, f *excelize.File) (interface{}, error) {
	var data []interface{}
	name := f.GetSheetName(sheetIndex)
	index := 1
	var fn []*FieldType
	for {
		x, _ := excelize.CoordinatesToCellName(index, 1)
		comment, _ := f.GetCellValue(name, x)

		x, _ = excelize.CoordinatesToCellName(index, 2)
		fieldType, _ := f.GetCellValue(name, x)
		x, _ = excelize.CoordinatesToCellName(index, 3)
		fieldName, _ := f.GetCellValue(name, x)

		if fieldType == "" {
			break
		}
		if comment != "" {
			fieldName1, _ := checkName(sname+"_"+fieldName, true)
			comments[fieldName1] = comment
		}

		fn = append(fn, &FieldType{
			Name: fieldName,
			Type: fieldType,
		})
		index++
	}

	index = 4
	for {
		x, _ := excelize.CoordinatesToCellName(1, index)
		str, _ := f.GetCellValue(name, x)
		if str == "" {
			break
		}

		row := make(map[string]interface{})
		for i := 1; i <= len(fn); i++ {

			x, _ = excelize.CoordinatesToCellName(i, index)
			str, _ = f.GetCellValue(name, x)

			//
			if str == "" {
				continue
			}

			v, err := exportCellJson(fn[i-1].Type, str)
			if err != nil {
				return data, fmt.Errorf("文件[%v]的页[%v]的%v行%v列的数据%v", file, name, index, i, err.Error())
			}

			row[fn[i-1].Name] = v
		}

		data = append(data, row)
		index++
	}

	return data, nil
}

func exportXLSX2Json(file, name string) (interface{}, error) {
	// fmt.Println("开始处理", file)
	f, err := excelize.OpenFile(file)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("file: %v", file))
	}

	return exportSheet2Json(file, name, 0, f)
}

// Run ...
func checkName(s string, firstUpper bool) (string, error) {
	if len(s) == 0 {
		return "", fmt.Errorf("为空")
	}
	// if s[0] < 'a' || s[0] > 'z' {
	// 	return "", fmt.Errorf("首字母必须为小写字母")
	// }
	var (
		buff bytes.Buffer
		u    int
	)

	for i := 0; i < len(s); i++ {
		b := s[i]
		if (b < '0' || b > '9') && (b < 'A' || b > 'Z') && (b < 'a' || b > 'z') && b != '_' && b != '-' {
			return "", fmt.Errorf("存在不合法字符")
		}

		if u == 1 || (firstUpper && i == 0) {
			if b >= 'a' && b <= 'z' {
				b -= 'a' - 'A'
			}
			u = 0
		} else if b == '_' {
			u = 1
			continue
		}

		buff.WriteByte(b)
	}

	return buff.String(), nil
}

func isFloat(s string) bool {
	n, _ := strconv.ParseInt(s, 10, 64)
	f, _ := strconv.ParseFloat(s, 64)
	return float64(n) != f
}

func parseArray(value []byte) (d int, t jsonparser.ValueType, b [][]byte) {
	d = 1
	for {
		v, t1, _, _ := jsonparser.Get(value, "[0]")
		if t1 != jsonparser.Array {
			t = t1
			jsonparser.ArrayEach(value, func(v1 []byte, dataType jsonparser.ValueType, offset int, err error) {
				b = append(b, v1)
			})

			return
		}

		d++
		value = v
	}
}
