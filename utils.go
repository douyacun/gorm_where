package database

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

const (
	DbRecordExistsError = "Error 1062: Duplicate entry"
)

func omitEmpty(value reflect.Value, opt sessionOptionTag) bool {
	if len(opt.defaultValue) != 0 {
		return false
	}

	switch value.Kind() {
	case reflect.Slice, reflect.Map:
		if value.IsNil() || value.Len() == 0 {
			return true
		}
	default:
		if !opt.empty && value.Interface() == reflect.Zero(value.Type()).Interface() {
			return true
		}
	}
	return false
}

func loadString(value reflect.Value, opt sessionOptionTag) string {
	result := ""
	if value.Kind() == reflect.String {
		result = value.String()
	}

	if len(result) == 0 {
		result = opt.defaultValue
	}

	return result
}

func loadInt(value reflect.Value, opt sessionOptionTag) int {
	result := int(0)
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = int(value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = int(value.Uint())
	}

	if result == 0 && len(opt.defaultValue) != 0 {
		result, _ = strconv.Atoi(opt.defaultValue)
	}

	return result
}

func isDbRecordExistsError(err error) bool {
	return strings.Contains(err.Error(), DbRecordExistsError)
}

func GetElem(elem reflect.Value) reflect.Value {
	for elem.Kind() == reflect.Ptr || elem.Kind() == reflect.Interface {
		elem = elem.Elem()
	}
	return elem
}

func CamelCaseToUnderscore(str string) string {
	var output []rune
	var segment []rune
	for _, r := range str {
		if !unicode.IsLower(r) && string(r) != "_" {
			output = addSegment(output, segment)
			segment = nil
		}
		segment = append(segment, unicode.ToLower(r))
	}
	output = addSegment(output, segment)
	return string(output)
}

func addSegment(inrune, segment []rune) []rune {
	if len(segment) == 0 {
		return inrune
	}
	if len(inrune) != 0 {
		inrune = append(inrune, '_')
	}
	inrune = append(inrune, segment...)
	return inrune
}

func CamelToSnake(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}
