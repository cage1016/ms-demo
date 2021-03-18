package gormtracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

const (
	gormSpanKey        = "__gorm_span"
	callBackBeforeName = "opentracing:before"
	callBackAfterName  = "opentracing:after"
	operName           = "db.gorm"
)

func before(db *gorm.DB) {
	span, _ := opentracing.StartSpanFromContext(db.Statement.Context, operName)
	db.InstanceSet(gormSpanKey, span)
}

func after(db *gorm.DB) {
	spanIn, y := db.InstanceGet(gormSpanKey)
	if !y {
		return
	}
	span := spanIn.(opentracing.Span)
	defer span.Finish()
	ext.DBType.Set(span, db.Statement.Dialector.Name())
	span.SetTag("db.table", db.Statement.Table)
	span.SetTag("db.SQL", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...))

	if err := db.Error; err != nil {
		span.LogFields(log.String("Error", err.Error()))
	}
}

type GormTracePlugin struct{}

func (g *GormTracePlugin) Name() string {
	return "GormTracePlugin"
}

func (g *GormTracePlugin) Initialize(db *gorm.DB) error {
	// before
	db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// after
	db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return nil
}

func NewGormTracePlugin() *GormTracePlugin {
	return &GormTracePlugin{}
}
