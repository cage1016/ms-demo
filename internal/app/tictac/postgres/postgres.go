package postgres

import (
	"context"
	"database/sql"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"gorm.io/gorm"

	"github.com/cage1016/ms-sample/internal/app/tictac/model"
)

type tictacRespository struct {
	mu  sync.RWMutex
	log log.Logger
	db  *gorm.DB
}

func New(db *gorm.DB, logger log.Logger) model.TictacRespository {
	return &tictacRespository{
		mu:  sync.RWMutex{},
		log: logger,
		db:  db,
	}
}

func (cr *tictacRespository) Tic(ctx context.Context, value int64) (err error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	err = cr.db.WithContext(ctx).Model(&model.Tictac{}).Where("1=1").Update("value", value).Error
	if err != nil {
		level.Error(cr.log).Log("method", "tictacRespository_tic", "err", err)
	}
	return
}

func (cr *tictacRespository) Tac(ctx context.Context) (res int64, err error) {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	err = cr.db.WithContext(ctx).Model(&model.Tictac{}).First(&res).Error
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		level.Error(cr.log).Log("method", "tictacRespository_tac", "err", err)
	}
	return
}
