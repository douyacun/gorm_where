package database

import "strings"

type sessionOptionTag struct {
	name         string
	op           string
	defaultValue string
	empty        bool
	noUpdate     bool
	sortBy       string
	page         int
	pageSize     int
	ignoreCopy   bool
}

func parseSessionOptionTag(tag string) sessionOptionTag {
	tags := strings.Split(tag, ",")
	option := sessionOptionTag{
		name: tags[0],
		op:   "equal",
	}

	option.name = tags[0]
	for _, tag := range tags[1:] {
		tag = strings.Trim(tag, " ")

		if tag == "empty" {
			option.empty = true
		}

		if tag == "no-update" {
			option.noUpdate = true
		}

		if tag == "ignore-copy" {
			option.ignoreCopy = true
		}

		if strings.HasPrefix(tag, "default:") {
			option.defaultValue = strings.TrimPrefix(tag, "default:")
		}

		if strings.HasPrefix(tag, "op:") {
			option.op = strings.TrimPrefix(tag, "op:")
		}
	}

	return option
}
