package db

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouterDomain struct {
	ID            int64           `gorm:"primaryKey, column:id"`
	Router        string          `gorm:"column:router"`
	Configuration json.RawMessage `gorm:"column:configuration" sql:"type:json"`
}

func (RouterDomain) TableName() string {
	return "petshop_gateway.router"
}

type GatewayDB struct {
	DB          *gorm.DB
	LoggerSugar *zap.SugaredLogger
}

type GatewayDBOption func(*GatewayDB)

func NewGatewayDB(opts ...GatewayDBOption) *GatewayDB {

	gdb := GatewayDB{}

	for _, opt := range opts {
		opt(&gdb)
	}

	return &gdb
}

func (gatewayDB *GatewayDB) GetAllRouter() []RouterDomain {
	ctx := context.Background()
	var routers []RouterDomain
	gatewayDB.DB.WithContext(ctx).Find(&routers)

	return routers
}

func WithDB(db *gorm.DB) GatewayDBOption {
	return func(gdb *GatewayDB) {
		gdb.DB = db
	}
}

func WithLogger(loggerSugar *zap.SugaredLogger) GatewayDBOption {
	return func(db *GatewayDB) {
		db.LoggerSugar = loggerSugar
	}
}
