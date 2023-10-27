package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"slices"
)

var (
	dbSql = `CREATE DATABASE [{db}] ON  PRIMARY
	( NAME = N'{db}', FILENAME = N'{disk}\\database\\{db}.mdf' , SIZE = 4096KB , FILEGROWTH = 1024KB )
	 LOG ON
	( NAME = N'{db}_log', FILENAME = N'{disk}\\database\\{db}_log.ldf' , SIZE = 1024KB , FILEGROWTH = 10%);`
)

func CreateDBSQL(dbName, disk string) string {
	str := strings.ReplaceAll(dbSql, "{db}", dbName)
	str = strings.ReplaceAll(str, "{disk}", disk)
	return str
}

func getCreateDesc(key string) string {
	arr := strings.Split(key, ";")
	for _, v := range arr {
		arr2 := strings.Split(v, ":")
		if arr2[0] == "create" {
			return arr2[1]
		}
	}
	return ""
}

func getInsertValue(rt reflect.Type, rv reflect.Value) string {
	var fields []string
	for i := 0; i < rt.NumField(); i++ {
		key := rt.Field(i).Tag.Get("orm")
		if key == "" || key == "-" {
			continue
		}
		desc := getCreateDesc(key)

		v := fmt.Sprint(rv.Field(i).Interface())
		if !strings.Contains(desc, "varchar") && v == "" {
			fields = append(fields, "NULL")
			continue
		}
		if strings.Contains(desc, "datetime") || strings.Contains(desc, "varchar") {
			v = "'" + template.HTMLEscapeString(v) + "'"
		}
		fields = append(fields,
			v,
		)
	}

	return fmt.Sprintf("(%v)", strings.Join(fields, ","))
}

func CreateSQL(v any) (string, []string) {
	rt := reflect.TypeOf(v).Elem()
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	str := "CREATE TABLE [%v].[dbo].[%v]("
	var arr []string
	var key string
	var indexes []string
	for i := 0; i < rt.NumField(); i++ {
		key = rt.Field(i).Tag.Get("orm")
		if key == "" || key == "-" {
			continue
		}
		name := rt.Field(i).Name
		name = strings.ToLower(name[:1]) + name[1:]
		arr = append(arr,
			fmt.Sprintf("[%v] %v", name, getCreateDesc(key)),
		)

		indexKey := rt.Field(i).Tag.Get("index")
		if indexKey != "" {
			index := fmt.Sprintf("CREATE CLUSTERED INDEX [idx_%%v_%v] ON [%%v].[dbo].[%%v]([%v] %v)", name, name, indexKey)
			indexes = append(indexes, index)
		}
	}

	create := str + strings.Join(arr, ",") + ") ON [PRIMARY]"

	return create, indexes
}

func InsertSQL(v any) string {
	var (
		e      reflect.Type
		arr    []string
		fields []string
	)
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		rv := reflect.ValueOf(v)
		e = rv.Index(0).Type()

		if e.Kind() == reflect.Ptr {
			e = e.Elem()
		}
		for i := 0; i < rv.Len(); i++ {
			ev := rv.Index(i)
			if ev.Kind() == reflect.Ptr {
				ev = ev.Elem()
			}

			arr = append(arr, getInsertValue(e, ev))
		}
	} else {
		rt := reflect.TypeOf(v)
		if rt.Kind() == reflect.Ptr {
			e = rt.Elem()
			arr = append(arr, getInsertValue(e, reflect.ValueOf(v).Elem()))
		} else {
			e = rt
			arr = append(arr, getInsertValue(e, reflect.ValueOf(v)))
		}

	}

	for i := 0; i < e.NumField(); i++ {
		key := e.Field(i).Tag.Get("orm")
		if key == "" || key == "-" {
			continue
		}
		name := e.Field(i).Name
		name = strings.ToLower(name[:1]) + name[1:]
		fields = append(fields,
			fmt.Sprintf("[%v]", name),
		)
	}

	str := "INSERT INTO [%v].[dbo].[%v] "
	return str +
		"(" + strings.Join(fields, ",") + ") VALUES" +
		strings.Join(arr, ",") +
		";"
}

func InsertSQLByLog(v any, ls []string) (string, error) {
	var (
		e      reflect.Type
		arr    []string
		fields []string
	)

	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Ptr {
		e = rt.Elem()
	} else {
		e = rt
	}

	for i := 0; i < e.NumField(); i++ {
		key := e.Field(i).Tag.Get("orm")
		if key == "" || key == "-" {
			continue
		}
		name := e.Field(i).Name
		name = strings.ToLower(name[:1]) + name[1:]
		fields = append(fields,
			fmt.Sprintf("[%v]", name),
		)
	}

	for _, v := range ls {
		varr := strings.Split(v, ";")
		var values []string
		index := 0
		for i := 0; i < e.NumField(); i++ {
			key := e.Field(i).Tag.Get("orm")
			if key == "" || key == "-" {
				continue
			}
			desc := getCreateDesc(key)
			value := ""
			if index < len(varr) {
				value = varr[index]
			}
			if strings.Contains(desc, "datetime") || strings.Contains(desc, "varchar") {
				values = append(values, fmt.Sprintf("'%v'", template.HTMLEscapeString(value)))

			} else {
				//int
				n, err := strconv.ParseInt(value, 10, 64)
				if value != "" && err != nil {
					return "", fmt.Errorf("日志字段[%v]格式错误, value: %v", e.Field(i).Name, value)
				}
				values = append(values, fmt.Sprint(n))
			}

			index++
		}

		arr = append(arr, fmt.Sprintf("(%v)", strings.Join(values, ",")))
	}

	str := "INSERT INTO [%v].[dbo].[%v] "
	return str +
		"(" + strings.Join(fields, ",") + ") VALUES" +
		strings.Join(arr, ",") +
		";", nil
}

func UpdateSQL(v any, where string, fields ...string) string {
	var (
		e    reflect.Type
		sets []string
	)

	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if rt.Kind() == reflect.Ptr {
		e = rt.Elem()
		rv = rv.Elem()
	} else {
		e = rt
	}

	for i := 0; i < e.NumField(); i++ {
		key := e.Field(i).Tag.Get("orm")
		if key == "" || key == "-" {
			continue
		}
		name := e.Field(i).Name
		name = strings.ToLower(name[:1]) + name[1:]
		if len(fields) > 0 && !slices.Contains(fields, name) {
			continue
		}
		desc := getCreateDesc(key)

		v := fmt.Sprint(rv.Field(i).Interface())
		if strings.Contains(desc, "datetime") || strings.Contains(desc, "varchar") {
			v = "'" + template.HTMLEscapeString(v) + "'"
		}

		sets = append(sets,
			fmt.Sprintf("[%v]=%v", name, v),
		)
	}

	// str := "UPDATE [%v].[dbo].[%v%v] SET "
	// return str +
	// 	strings.Join(sets, ", ") +
	// 	" WHERE " + where
	var b strings.Builder
	b.WriteString("UPDATE [%v].[dbo].[%v%v] SET ")
	b.WriteString(strings.Join(sets, ","))
	if where != "" {
		b.WriteString(" WHERE ")
		b.WriteString(where)
	}
	return b.String()
}
