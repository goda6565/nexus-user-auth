package repository_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/goda6565/nexus-user-auth/domain/user/entity"
	"github.com/goda6565/nexus-user-auth/domain/user/repository"
	"github.com/goda6565/nexus-user-auth/domain/user/value"
	. "github.com/goda6565/nexus-user-auth/infrastructure/database/repository"
	"github.com/goda6565/nexus-user-auth/pkg/tester"
	"github.com/goda6565/nexus-user-auth/pkg/utils"
)

type UserRepositoryImplTestSuite struct {
	tester.DBSQLiteSuite
	userRepo repository.UserRepository
}

func TestUserRepositoryImplTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryImplTestSuite))
}

func (suite *UserRepositoryImplTestSuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.userRepo = NewUserRepository(suite.DB)
}

func (suite *UserRepositoryImplTestSuite) TestCreateUser() {
	// 値オブジェクトの生成
	email, err := value.NewUserEmail("test@example.com")
	suite.NoError(err, "メールアドレスの生成に失敗してはいけない")
	password, err := value.NewUserPassword("Password123!")
	suite.NoError(err, "パスワードの生成に失敗してはいけない")
	username, err := value.NewUserUsername("testuser")
	suite.NoError(err, "ユーザー名の生成に失敗してはいけない")

	// 新規ユーザーエンティティの生成（ファクトリ関数を利用）
	userEntity, err := entity.NewUser(email, password, username)
	suite.NoError(err, "ユーザーエンティティの生成に失敗してはいけない")

	// リポジトリを通じてユーザー作成
	createdUser, err := suite.userRepo.CreateUser(userEntity)
	suite.NoError(err, "ユーザー作成に失敗してはいけない")
	suite.NotNil(createdUser, "作成されたユーザーがnilであってはならない")
	suite.NotEmpty(createdUser.ObjID().Value(), "自動生成されたユーザーIDが空であってはならない")
	suite.Equal("test@example.com", createdUser.Email().Value(), "メールアドレスが期待通りであること")
	suite.Equal("testuser", createdUser.Username().Value(), "ユーザー名が期待通りであること")
	suite.NoError(utils.CheckPassword(createdUser.Password().Value(), "Password123!"), "パスワードが期待通りであること")
}

func (suite *UserRepositoryImplTestSuite) TestGetUserByEmail() {
	// 作成済みのユーザーを用意
	email, err := value.NewUserEmail("find@example.com")
	suite.NoError(err)
	password, err := value.NewUserPassword("Password123!")
	suite.NoError(err)
	username, err := value.NewUserUsername("finduser")
	suite.NoError(err)

	userEntity, err := entity.NewUser(email, password, username)
	suite.NoError(err)

	createdUser, err := suite.userRepo.CreateUser(userEntity)
	suite.NoError(err)

	// メールアドレスでユーザー取得
	foundUser, err := suite.userRepo.GetUserByEmail("find@example.com")
	suite.NoError(err)
	suite.NotNil(foundUser)
	suite.Equal(createdUser.ObjID().Value(), foundUser.ObjID().Value(), "取得したユーザーIDが一致すること")
	suite.Equal(createdUser.Email().Value(), foundUser.Email().Value(), "取得したユーザーのメールが一致すること")
	suite.Equal(createdUser.Username().Value(), foundUser.Username().Value(), "取得したユーザー名が一致すること")
	suite.NoError(utils.CheckPassword(foundUser.Password().Value(), "Password123!"), "取得したユーザーのパスワードが一致すること")
}

func (suite *UserRepositoryImplTestSuite) TestGetUserByObjID() {
	// 作成済みのユーザーを用意
	email, err := value.NewUserEmail("obj@example.com")
	suite.NoError(err)
	password, err := value.NewUserPassword("Password123!")
	suite.NoError(err)
	username, err := value.NewUserUsername("objuser")
	suite.NoError(err)

	userEntity, err := entity.NewUser(email, password, username)
	suite.NoError(err)

	createdUser, err := suite.userRepo.CreateUser(userEntity)
	suite.NoError(err)

	// オブジェクトIDでユーザー取得
	foundUser, err := suite.userRepo.GetUserByObjID(createdUser.ObjID().Value())
	suite.NoError(err)
	suite.NotNil(foundUser)
	suite.Equal(createdUser.ObjID().Value(), foundUser.ObjID().Value(), "取得したユーザーIDが一致すること")
	suite.Equal(createdUser.Email().Value(), foundUser.Email().Value(), "取得したユーザーのメールが一致すること")
	suite.Equal(createdUser.Username().Value(), foundUser.Username().Value(), "取得したユーザー名が一致すること")
	suite.NoError(utils.CheckPassword(foundUser.Password().Value(), "Password123!"), "取得したユーザーのパスワードが一致すること")
}

func (suite *UserRepositoryImplTestSuite) TestUpdateUser() {
	// ユーザー作成
	email, err := value.NewUserEmail("update@example.com")
	suite.NoError(err)
	password, err := value.NewUserPassword("Password123!")
	suite.NoError(err)
	username, err := value.NewUserUsername("updateuser")
	suite.NoError(err)

	userEntity, err := entity.NewUser(email, password, username)
	suite.NoError(err)

	createdUser, err := suite.userRepo.CreateUser(userEntity)
	suite.NoError(err)

	// ユーザー名を更新
	updatedUsername, err := value.NewUserUsername("updateduser")
	suite.NoError(err)

	// 既存のユーザーエンティティから新たなエンティティを再構築（BuildUserを利用）
	updatedUser, err := entity.BuildUser(
		createdUser.ObjID(),
		createdUser.Email(),
		createdUser.Password(),
		updatedUsername,
		createdUser.AvatarURL(),
		createdUser.EmailVerifiedAt(),
		createdUser.LastLoginAt(),
		createdUser.Role(),
	)
	suite.NoError(err)

	resultUser, err := suite.userRepo.UpdateUser(updatedUser)
	suite.NoError(err)
	suite.NotNil(resultUser)
	suite.Equal("updateduser", resultUser.Username().Value(), "更新後のユーザー名が期待通りであること")
	suite.Equal(createdUser.ObjID().Value(), resultUser.ObjID().Value(), "更新後のユーザーIDが一致すること")
	suite.Equal(createdUser.Email().Value(), resultUser.Email().Value(), "更新後のユーザーのメールが一致すること")
	suite.NoError(utils.CheckPassword(resultUser.Password().Value(), "Password123!"), "更新後のユーザーのパスワードが一致すること")
}

func (suite *UserRepositoryImplTestSuite) TestDeleteUser() {
	// ユーザー作成
	email, err := value.NewUserEmail("delete@example.com")
	suite.NoError(err)
	password, err := value.NewUserPassword("Password123!")
	suite.NoError(err)
	username, err := value.NewUserUsername("deleteuser")
	suite.NoError(err)

	userEntity, err := entity.NewUser(email, password, username)
	suite.NoError(err)

	createdUser, err := suite.userRepo.CreateUser(userEntity)
	suite.NoError(err)

	// ユーザー削除
	err = suite.userRepo.DeleteUser(createdUser.ObjID().Value())
	suite.NoError(err)

	// 削除後にユーザー取得を試みる
	deletedUser, err := suite.userRepo.GetUserByObjID(createdUser.ObjID().Value())
	suite.Error(err, "削除されたユーザーは取得できないはず")
	suite.Nil(deletedUser, "削除されたユーザーはnilであるはず")
}
