package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

var (
	ckDbSql = `CREATE DATABASE IF NOT EXISTS {db}`
)

func CreateClickhouseDBSQL(dbName string) string {
	str := strings.ReplaceAll(ckDbSql, "{db}", dbName)
	return str
}

func CreateClickhouseSQL(v any) string {
	rt := reflect.TypeOf(v).Elem()
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	str := "CREATE TABLE IF NOT EXISTS %v.%v ("
	var arr []string
	var key, partStr, sortStr string
	for i := 0; i < rt.NumField(); i++ {
		key = rt.Field(i).Tag.Get("ck")
		if key == "" || key == "-" {
			continue
		}
		name := rt.Field(i).Name
		name = strings.ToLower(name[:1]) + name[1:]
		arr = append(arr,
			fmt.Sprintf("%v %v", name, getCreateDesc(key)),
		)

		partKey := rt.Field(i).Tag.Get("ck_part")
		if partKey != "" {
			partStr = fmt.Sprintf(" PARTITION BY %v", partKey)
		}

		sortKey := rt.Field(i).Tag.Get("ck_sort")
		if partKey != "" {
			sortStr = fmt.Sprintf(" ORDER BY (%v)", sortKey)
		}
	}

	create := str + strings.Join(arr, ",") + ") ENGINE=MergeTree" + partStr + sortStr

	return create
}

func InsertClickhouseSQLByLog(v any, l string) (fields []any, err error) {
	var (
		e reflect.Type
	)
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Ptr {
		e = rt.Elem()
	} else {
		e = rt
	}

	varr := strings.Split(l, ";")
	index := 0
	for i := 0; i < e.NumField(); i++ {
		key := e.Field(i).Tag.Get("ck")
		if key == "" || key == "-" {
			continue
		}
		desc := getCreateDesc(key)
		value := ""
		if index < len(varr) {
			value = varr[index]
		}
		if strings.Contains(desc, "String") {
			fields = append(fields, template.HTMLEscapeString(value))
		} else if strings.Contains(desc, "UInt64") {
			n, err := strconv.ParseUint(value, 10, 64)
			if value != "" && err != nil {
				return nil, fmt.Errorf("日志字段[%v]格式错误, value: %v", e.Field(i).Name, value)
			}
			fields = append(fields, n)
		} else if strings.Contains(desc, "Int64") {
			n, err := strconv.ParseInt(value, 10, 64)
			if value != "" && err != nil {
				return nil, fmt.Errorf("日志字段[%v]格式错误, value: %v", e.Field(i).Name, value)
			}
			fields = append(fields, n)
		} else if strings.Contains(desc, "UInt32") {
			n, err := strconv.ParseUint(value, 10, 32)
			if value != "" && err != nil {
				return nil, fmt.Errorf("日志字段[%v]格式错误, value: %v", e.Field(i).Name, value)
			}
			fields = append(fields, uint32(n))
		} else {
			n, err := strconv.ParseInt(value, 10, 32)
			if value != "" && err != nil {
				return nil, fmt.Errorf("日志字段[%v]格式错误, value: %v", e.Field(i).Name, value)
			}
			fields = append(fields, int32(n))
		}

		index++
	}

	return fields, nil
}
