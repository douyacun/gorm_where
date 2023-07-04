package database

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type _sessionCondition struct{}

var SessionCondition = &_sessionCondition{}

func (*_sessionCondition) WithTable(table string) SessionOption {
	return SessionOption{
		Type: sessionOptionTable,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Table(table)
		},
	}
}

func (*_sessionCondition) WithModel(model interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionTable,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Model(model)
		},
	}
}

func (*_sessionCondition) WithJoin(joinOperator string, table interface{}, condition string, args ...interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionTable,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Joins(fmt.Sprintf("%s `%s` on %s ", joinOperator, table, condition))
		},
	}
}

func (*_sessionCondition) WithJoinAs(joinOperator string, table, as interface{}, condition string, args ...interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionTable,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Joins(fmt.Sprintf("%s `%s` on %s ", joinOperator, table, condition))
		},
	}
}

func (*_sessionCondition) WithPage(page int, pageSize int) SessionOption {
	if page == 0 {
		page = 1
	} else if page < 0 {
		page = -1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return SessionOption{
		Type: sessionOptionLimitation,
		Process: func(session *gorm.DB) *gorm.DB {
			if page >= 1 {
				return session.Limit(pageSize).Offset((page - 1) * pageSize)
			}
			return session
		},
	}
}

func (*_sessionCondition) WithSort(sortBy string) SessionOption {
	return SessionOption{
		Type: sessionOptionOrderBy,
		Process: func(session *gorm.DB) *gorm.DB {
			if len(sortBy) == 0 {
				return session
			}
			return session.Order(sortBy)
		},
	}
}

func (*_sessionCondition) WithGroupBy(groupBy string) SessionOption {
	return SessionOption{
		Type: sessionOptionGroupBy,
		Process: func(session *gorm.DB) *gorm.DB {
			if len(groupBy) == 0 {
				return session
			}
			return session.Group(groupBy)
		},
	}
}

func (*_sessionCondition) WithHaving(having string) SessionOption {
	return SessionOption{
		Type: sessionOptionGroupBy,
		Process: func(session *gorm.DB) *gorm.DB {
			if len(having) == 0 {
				return session
			}
			return session.Having(having)
		},
	}
}

func (*_sessionCondition) WithGt(field string, value interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(fmt.Sprintf("%s > ?", field), value)
		},
	}
}

func (*_sessionCondition) WithGte(field string, value interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(fmt.Sprintf("%s >= ?", field), value)
		},
	}
}

func (*_sessionCondition) WithLt(field string, value interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(fmt.Sprintf("%s < ?", field), value)
		},
	}
}

func (*_sessionCondition) WithLte(field string, value interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(fmt.Sprintf("%s <= ?", field), value)
		},
	}
}

func (*_sessionCondition) WithEqual(field string, value interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(fmt.Sprintf("%s = ?", field), value)
		},
	}
}

func (*_sessionCondition) WithMod(field string, value int) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(field, value)
		},
	}
}

func (*_sessionCondition) WithNotEqual(field string, value interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(fmt.Sprintf("%s != ?", field), value)
		},
	}
}

func (*_sessionCondition) WithLike(field string, value string) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(fmt.Sprintf("%s like ?", field), "%"+value+"%")
		},
	}
}

func (*_sessionCondition) WithLikeOr(value string, fieldList ...string) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			where := make([]string, 0)
			for _, field := range fieldList {
				where = append(where, fmt.Sprintf("(%s like '%s')", field, "%"+value+"%"))
			}

			return session.Where(strings.Join(where, " or "))
		},
	}
}

func (*_sessionCondition) With(process func(session *gorm.DB) *gorm.DB) SessionOption {
	return SessionOption{
		Type:    sessionOptionOther,
		Process: process,
	}
}

func (*_sessionCondition) WithIn(field string, data ...interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(fmt.Sprintf("%s in ?", field), data)
		},
	}
}

func (*_sessionCondition) WithNotIn(field string, data ...interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(fmt.Sprintf("%s not in ?", field), data)
		},
	}
}

func (*_sessionCondition) WithBetween(field string, begin interface{}, end interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(fmt.Sprintf("%s between ? and ?", field), begin, end)
		},
	}
}

func (*_sessionCondition) WithNoCount() SessionOption {
	return SessionOption{
		Type: sessionOptionNoCount,
		Process: func(session *gorm.DB) *gorm.DB {
			return session
		},
	}
}

func (*_sessionCondition) WithWhere(condition string, param ...interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(condition, param...)
		},
	}
}

func (*_sessionCondition) WithOr(condition string, param ...interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionOther,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Or(condition, param...)
		},
	}
}

func (*_sessionCondition) WithDistinct(condition string) SessionOption {
	return SessionOption{
		Type: sessionOptionSelect,
		Process: func(session *gorm.DB) *gorm.DB {
			//session.Distinct(condition)
			return session
		},
	}
}

func (*_sessionCondition) WithCountSelect(selectStr string) SessionOption {
	return SessionOption{
		Type: sessionOptionCountSelect,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Select(selectStr)
		},
	}
}

func (*_sessionCondition) WithSelect(field interface{}, fields ...interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionSelect,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Select(field, fields...)
		},
	}
}

func (*_sessionCondition) WithUpdateCols(cols []interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionUpdateCols,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Select(cols)
		},
	}
}

func (*_sessionCondition) ParseSessionOptions(source interface{}) []SessionOption {
	return ParseSessionOption(source)
}

func (*_sessionCondition) WithConflict(fields ...string) SessionOption {
	return SessionOption{
		Type: sessionOptionUpdateCols,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Clauses(clause.OnConflict{
				DoUpdates: clause.AssignmentColumns(fields),
			})
		},
	}
}

func (*_sessionCondition) WithIgnore() SessionOption {
	return SessionOption{
		Type: sessionOptionUpdateCols,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Clauses(clause.Insert{Modifier: "IGNORE"})
		},
	}
}

func (*_sessionCondition) WithMatch(fields []string, value interface{}) SessionOption {
	return SessionOption{
		Type: sessionOptionMatch,
		Process: func(session *gorm.DB) *gorm.DB {
			return session.Where(fmt.Sprintf("match(%s) against('%s' in NATURAL LANGUAGE MODE)", strings.Join(fields, ","), value))
		},
	}
}
