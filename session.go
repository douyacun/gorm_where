package database

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/db"
	"gorm.io/gorm"
)

type _session struct{}

var Session = &_session{}

const (
	DbSessionAppointTag        = "__session_appoint__"
	DbSessionRuntimeTag        = "__session_runtime__"
	DbSessionRuntimeDefaultTag = "default"
)

func (*_session) StartToContext(ctx context.Context) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(DbSessionRuntimeTag, db.DB.Session(&gorm.Session{}))
		return ginCtx
	} else {
		return context.WithValue(ctx, DbSessionRuntimeTag, db.DB.Session(&gorm.Session{}))
	}
}

func (s *_session) Close(ctx context.Context) {
	sess, _ := s.getRuntimeFromContext(ctx)
	if sess != nil {
		sess.Commit()
	}
}

func (*_session) getRuntimeFromContext(ctx context.Context) (*gorm.DB, error) {
	ri := ctx.Value(DbSessionRuntimeTag)
	if ri == nil {
		return nil, errors.New("getRuntimeFromContext fail: session runtime not found in context")
	}
	r, ok := ri.(*gorm.DB)
	if !ok {
		return nil, errors.New("getRuntimeFromContext fail: not is sessionRuntime")
	}

	return r, nil
}

func (*_session) getConnection(ctx context.Context) string {
	conn := DbSessionRuntimeDefaultTag

	connFromContext := ctx.Value(DbSessionAppointTag)
	if connFromContext != nil {
		appointConn := connFromContext.(string)
		if appointConn != "" {
			conn = appointConn
		}
	}

	return conn
}

func (*_session) getFromContext(ctx context.Context) (*gorm.DB, error) {
	r, err := Session.getRuntimeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *_session) GetFromContext(ctx context.Context) (*gorm.DB, error) {
	return s.getFromContext(ctx)
}
