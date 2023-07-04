package database

import (
	"github.com/teablog/tea/internal/db"
	"gorm.io/gorm"
)

func newSessionRuntime() *gorm.DB {
	return db.DB.Session(&gorm.Session{SkipDefaultTransaction: true})
}