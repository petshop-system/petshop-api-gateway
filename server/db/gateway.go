package db

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GatewayDomain struct {
	ID            int64           `gorm:"primaryKey, column:id"`
	Router        string          `gorm:"column:router"`
	Configuration json.RawMessage `gorm:"column:configuration" sql:"type:json"`
}

func (GatewayDomain) TableName() string {
	return "petshop_gateway.gateway"
}

type GatewayDB struct {
	DB          *gorm.DB
	LoggerSugar *zap.SugaredLogger
}

func NewGatewayDB(db *gorm.DB, loggerSugar *zap.SugaredLogger) *GatewayDB {
	return &GatewayDB{
		DB:          db,
		LoggerSugar: loggerSugar,
	}
}

func (gatewayDB *GatewayDB) GetAllRouter() []GatewayDomain {
	ctx := context.Background()
	var routers []GatewayDomain
	gatewayDB.DB.WithContext(ctx).Find(&routers)

	return routers
}
