package database

import (
	"context"
	"github.com/pkg/errors"
	"github.com/teablog/tea/internal/db"
	"gorm.io/gorm"
)

type DAO struct {
}

func (dao *DAO) SelectByPage(ctx context.Context, list interface{}, searchOpt ...SessionOption) (int64, error) {
	var err error

	needCount := true
	for _, opt := range searchOpt {
		if opt.Type == sessionOptionNoCount {
			needCount = false
		}
	}
	count := int64(0)
	if needCount {
		session := db.DB.Session(&gorm.Session{})
		for _, opt := range searchOpt {
			session = processSessionOptionForCount(opt, session)
		}
		res := session.Count(&count)
		if res.Error != nil {
			return 0, errors.Wrapf(res.Error, "{{查询条数失败!}}")
		}
	}

	session, err := Session.getFromContext(ctx)
	if err != nil {
		return 0, err
	}
	for _, opt := range searchOpt {
		session = processSessionOption(opt, session)
	}
	res := session.Find(list)
	if res.Error != nil {
		return 0, errors.Wrapf(res.Error, "{{查询失败}}")
	}

	return count, nil
}

// Select 单个查询尽量用get，因为select返回的不是nil 是id=0的对象
func (dao *DAO) Select(ctx context.Context, list interface{}, searchOpt ...SessionOption) error {
	session, err := Session.getFromContext(ctx)
	if err != nil {
		return err
	}

	for _, opt := range searchOpt {
		session = processSessionOption(opt, session)
	}
	res := session.Find(list)
	if res.Error != nil {
		return errors.Wrapf(res.Error, "{{查询失败}}")
	}

	return nil
}

func (dao *DAO) Get(ctx context.Context, row interface{}, searchOpt ...SessionOption) (bool, error) {
	sess, err := Session.getFromContext(ctx)
	if err != nil {
		return false, err
	}
	for _, opt := range searchOpt {
		sess = processSessionOption(opt, sess)
	}
	res := sess.Take(row)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errors.Wrapf(err, "{{查询失败}}")
	}
	return true, nil
}

func (dao *DAO) Count(ctx context.Context, searchOpt ...SessionOption) (int64, error) {
	sess, err := Session.getFromContext(ctx)
	if err != nil {
		return 0, err
	}
	for _, opt := range searchOpt {
		sess = processSessionOption(opt, sess)
	}
	count := int64(0)
	res := sess.Count(&count)
	if res.Error != nil {
		return 0, errors.Wrapf(res.Error, "{{查询失败}}")
	}
	return count, nil
}

func (dao *DAO) Updates(ctx context.Context, data interface{}, searchOpt ...SessionOption) (int64, error) {
	session, err := Session.getFromContext(ctx)
	if err != nil {
		return 0, err
	}
	session = session.Model(data)
	for _, opt := range searchOpt {
		session = opt.Process(session)
	}
	res := session.Updates(data)
	if res.Error != nil {
		return 0, res.Error
	}

	return res.RowsAffected, res.Error
}

func (dao *DAO) Insert(ctx context.Context, data interface{}, searchOpt ...SessionOption) (int64, error) {
	session, err := Session.getFromContext(ctx)
	if err != nil {
		return 0, err
	}
	for _, opt := range searchOpt {
		session = opt.Process(session)
	}
	res := session.Create(data)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, res.Error
}

func (dao *DAO) Delete(ctx context.Context, data interface{}, searchOpt ...SessionOption) (int64, error) {
	session, err := Session.getFromContext(ctx)
	if err != nil {
		return 0, err
	}
	for _, opt := range searchOpt {
		session = opt.Process(session)
	}
	res := session.Delete(data)
	return res.RowsAffected, res.Error
}

func (dao *DAO) Pluck(ctx context.Context, field string, data []interface{}, searchOpt ...SessionOption) error {
	session, err := Session.getFromContext(ctx)
	if err != nil {
		return err
	}
	for _, opt := range searchOpt {
		session = opt.Process(session)
	}
	res := session.Pluck(field, data)
	return res.Error
}

func (dao *DAO) Save(ctx context.Context, data interface{}, searchOpt ...SessionOption) (int64, error) {
	session, err := Session.getFromContext(ctx)
	if err != nil {
		return 0, err
	}
	for _, opt := range searchOpt {
		session = opt.Process(session)
	}
	res := session.Create(data)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, res.Error
}

func (dao *DAO) CreateInBatches(ctx context.Context, data interface{}, options ...SessionOption) error {
	session, err := Session.getFromContext(ctx)
	if err != nil {
		return err
	}
	for _, opt := range options {
		session = opt.Process(session)
	}
	res := session.CreateInBatches(data, 100)
	return res.Error
}

func (dao *DAO) Transaction(ctx context.Context, actions ...func(ctx context.Context) error) error {
	sess, err := Session.getFromContext(ctx)
	if err != nil {
		return err
	}

	sess.Begin()
	for _, act := range actions {
		if actError := act(ctx); actError != nil {
			sess.Rollback()
			return actError
		}
	}
	sess.Commit()

	return nil
}
