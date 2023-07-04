package database

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

const (
	sessionOptionOther       = 0
	sessionOptionLimitation  = 1
	sessionOptionCountSelect = 2
	sessionOptionSelect      = 3
	sessionOptionGroupBy     = 4
	sessionOptionOrderBy     = 5
	sessionOptionUpdateCols  = 6
	sessionOptionTable       = 7
	sessionOptionJoin        = 8
	sessionOptionNoCount     = 9
	sessionOptionMatch       = 10
)

type SessionOption struct {
	Type    int
	Process func(builder *gorm.DB) *gorm.DB
}

type SessionOptionList []SessionOption

func ParseSessionOption(data interface{}) []SessionOption {
	result := make([]SessionOption, 0)

	elem := GetElem(reflect.ValueOf(data))

	sortBy := ""
	groupBy := ""
	page := 0
	pageSize := 10
	hasPage := true
	hasCount := true

	for i := 0; i != elem.NumField(); i++ {
		field := elem.Type().Field(i)
		fieldValue := GetElem(elem.Field(i))
		tag := field.Tag.Get("session")

		opt := parseSessionOptionTag(tag)

		if opt.name == "-" {
			continue
		}

		if omitEmpty(fieldValue, opt) {
			continue
		}

		if opt.name == "sort_by" {
			sortBy = loadString(fieldValue, opt)
			continue
		}

		if opt.name == "group_by" {
			groupBy = loadString(fieldValue, opt)
			continue
		}

		if opt.name == "no_count" {
			hasCount = false
			continue
		}

		if opt.name == "page" {
			hasPage = true
			page = loadInt(fieldValue, opt)
			continue
		}

		if opt.name == "page_size" {
			hasPage = true
			if page == 0 {
				page = 1
			}
			pageSize = loadInt(fieldValue, opt)
			continue
		}

		if opt.name == "select" {
			switch fieldValue.Kind() {
			case reflect.String:
				result = append(result, SessionCondition.WithSelect(fieldValue.String()))
			case reflect.Slice:
				fields := make([]interface{}, 0)
				for _, v := range fieldValue.Interface().([]string) {
					fields = append(fields, v)
				}
				if len(fields) == 1 {
					result = append(result, SessionCondition.WithSelect(fields[0]))
				} else if len(fields) > 1 {
					result = append(result, SessionCondition.WithSelect(fields[0], fields[1:]...))
				}
			}
			continue
		}

		switch opt.op {
		case "equal":
			result = append(result, SessionCondition.WithEqual(opt.name, fieldValue.Interface()))
		case "not_equal":
			result = append(result, SessionCondition.WithNotEqual(opt.name, fieldValue.Interface()))
		case "like":
			result = append(result, SessionCondition.WithLike(opt.name, loadString(fieldValue, opt)))
		case "like_or":
			result = append(result, SessionCondition.WithLikeOr(loadString(fieldValue, opt), strings.Split(opt.name, "&")...))
		case "in":
			in := make([]interface{}, fieldValue.Len())
			for i := 0; i != fieldValue.Len(); i++ {
				in[i] = fieldValue.Index(i).Interface()
			}
			result = append(result, SessionCondition.WithIn(opt.name, in...))
		case "not_in":
			in := make([]interface{}, fieldValue.Len())
			for i := 0; i != fieldValue.Len(); i++ {
				in[i] = fieldValue.Index(i).Interface()
			}
			result = append(result, SessionCondition.WithNotIn(opt.name, in...))
		case "lt":
			result = append(result, SessionCondition.WithLt(opt.name, fieldValue.Interface()))
		case "lte":
			result = append(result, SessionCondition.WithLte(opt.name, fieldValue.Interface()))
		case "gt":
			result = append(result, SessionCondition.WithGt(opt.name, fieldValue.Interface()))
		case "gte":
			result = append(result, SessionCondition.WithGte(opt.name, fieldValue.Interface()))
		case "json_contains":
			result = append(
				result,
				SessionCondition.WithWhere(
					fmt.Sprintf("JSON_CONTAINS(%s, ?)", opt.name),
					fieldValue.Interface(),
				),
			)
		case "match":
			result = append(result, SessionCondition.WithMatch(strings.Split(opt.name, "&"), fieldValue.Interface()))
		default:
			result = append(result, SessionCondition.WithWhere(opt.op, fieldValue.Interface()))
		}
	}

	if len(sortBy) != 0 {
		result = append(result, SessionCondition.WithSort(sortBy))
	}

	if len(groupBy) != 0 {
		result = append(result, SessionCondition.WithGroupBy(groupBy))
	}

	if hasPage && hasCount {
		result = append(result, SessionCondition.WithPage(page, pageSize))
	}

	if !hasCount {
		result = append(result, SessionCondition.WithNoCount())
	}

	return result
}

func processSessionOptionForCount(searchOpt SessionOption, session *gorm.DB) *gorm.DB {
	if searchOpt.Type == sessionOptionUpdateCols {
		return session
	}
	if searchOpt.Type != sessionOptionLimitation && searchOpt.Type != sessionOptionSelect &&
		searchOpt.Type != sessionOptionOrderBy && searchOpt.Type != sessionOptionGroupBy {
		session = searchOpt.Process(session)
	}
	return session
}

func processSessionOption(searchOpt SessionOption, session *gorm.DB) *gorm.DB {
	if searchOpt.Type == sessionOptionUpdateCols {
		return session
	}
	if searchOpt.Type != sessionOptionCountSelect {
		return searchOpt.Process(session)
	}
	return session
}

//func processCondSessionOptionForCount(session *gorm.DB, option SessionOption) builder.Cond {
//	switch option.Type {
//	case sessionOptionUpdateCols, sessionOptionLimitation, sessionOptionSelect, sessionOptionOrderBy, sessionOptionGroupBy:
//		return nil
//	case sessionOptionOther:
//		return option.Cond()
//	default:
//		logrus.Debug(token.Serialize.Encode(option))
//		option.Process(session)
//		return nil
//	}
//}
//
//func processCondSessionOption(session *gorm.DB, option SessionOption) builder.Cond {
//	switch option.Type {
//	case sessionOptionUpdateCols, sessionOptionCountSelect:
//		return nil
//	case sessionOptionOther:
//		return option.Cond()
//	default:
//		option.Process(session)
//		return nil
//	}
//}
//
