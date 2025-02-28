package database

import (
	"errors"
	"fmt"

	"github.com/goda6565/nexus-user-auth/infrastructure/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	InstanceSQLite int = iota
	InstancePostgres
)

var (
	errInvalidDBInstance = errors.New("invalid db instance")
)

func NewDBInstance(instance int) (db *gorm.DB, err error) {
	// DBのインスタンスを生成
	switch instance {
	case InstancePostgres:
		configs := NewConfigPostgres()
		dsn := fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			configs.User,
			configs.Password,
			configs.Host,
			configs.Port,
			configs.Database,
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case InstanceSQLite:
		configs := NewConfigSQLite()
		db, err = gorm.Open(sqlite.Open(configs.Database), &gorm.Config{})
	default:
		return nil, errInvalidDBInstance
	}
	return db, err
}

func NewGormModel() []interface{} {
	// マイグレーション対象のモデルを返す
	return []interface{}{
		&models.User{},
	}
}
