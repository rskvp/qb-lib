package qbl_commons

import (
	qbc "github.com/rskvp/qb-core"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	ModeProduction = qbc.ModeProduction
	ModeDebug      = qbc.ModeDebug
)

func GormConfig(mode string) (config *gorm.Config) {
	if mode == ModeDebug {
		config = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}
	} else {
		config = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}
	return
}
