package tester

import (
	"fmt"
	"os"

	"github.com/goda6565/nexus-user-auth/infrastructure/database"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type DBSQLiteSuite struct {
	suite.Suite
	DB     *gorm.DB
	DBName string
}

func (suite *DBSQLiteSuite) SetupSuite() {
	suite.DBName = fmt.Sprintf("%s.unittest.sqlite", suite.T().Name()) // test名をDB名にする
	os.Setenv("DB_NAME", suite.DBName)
	db, err := database.NewDBInstance(database.InstanceSQLite)
	suite.Assert().Nil(err)
	suite.DB = db

	for _, model := range database.NewGormModel() {
		err := suite.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

func (suite *DBSQLiteSuite) TearDownSuite() { // テスト終了時にDBファイルを削除
	err := os.Remove(suite.DBName)
	suite.Assert().Nil(err)
	os.Unsetenv(suite.DBName)
}