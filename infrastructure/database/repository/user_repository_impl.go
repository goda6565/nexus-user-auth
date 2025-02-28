package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/goda6565/nexus-user-auth/domain/user/entity"
	"github.com/goda6565/nexus-user-auth/domain/user/repository"
	"github.com/goda6565/nexus-user-auth/errs"
	"github.com/goda6565/nexus-user-auth/infrastructure/database/adapter"
	"github.com/goda6565/nexus-user-auth/infrastructure/database/models"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(user *entity.User) (*entity.User, error) {
	tx := r.db.Create(adapter.NewUserAdapter().Convert(user))
	if tx.Error != nil {
		return nil, errs.NewInfraError(fmt.Errorf("ユーザー作成に失敗しました: %w", tx.Error).Error())
	}
	return user, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*entity.User, error) {
	var modelUser models.User
	tx := r.db.Where("email = ?", email).First(&modelUser)
	if tx.Error != nil {
		return nil, errs.NewInfraError(fmt.Errorf("メールアドレス(%s)でのユーザー取得に失敗しました: %w", email, tx.Error).Error())
	}
	user, err := adapter.NewUserAdapter().ReBuild(&modelUser)
	if err != nil {
		return nil, errs.NewInfraError(fmt.Errorf("ユーザーエンティティの再構築に失敗しました: %w", err).Error())
	}
	return user, nil
}

func (r *UserRepositoryImpl) GetUserByObjID(objID string) (*entity.User, error) {
	var modelUser models.User
	tx := r.db.Where("obj_id = ?", objID).First(&modelUser)
	if tx.Error != nil {
		return nil, errs.NewInfraError(fmt.Errorf("オブジェクトID(%s)でのユーザー取得に失敗しました: %w", objID, tx.Error).Error())
	}
	user, err := adapter.NewUserAdapter().ReBuild(&modelUser)
	if err != nil {
		return nil, errs.NewInfraError(fmt.Errorf("ユーザーエンティティの再構築に失敗しました: %w", err).Error())
	}
	return user, nil
}

func (r *UserRepositoryImpl) UpdateUser(user *entity.User) (*entity.User, error) {
	var modelUser models.User
	tx := r.db.Where("obj_id = ?", user.ObjID().Value()).First(&modelUser)
	if tx.Error != nil {
		return nil, errs.NewInfraError(fmt.Errorf("ユーザー更新のための既存レコード取得に失敗しました: %w", tx.Error).Error())
	}

	// ドメインエンティティを永続化用モデルに変換
	converted, ok := adapter.NewUserAdapter().Convert(user).(*models.User)
	if !ok {
		return nil, errs.NewInfraError("変換されたモデルが *models.User ではありません。")
	}

	// 既存の modelUser のフィールドを更新
	modelUser.Email = converted.Email
	modelUser.Password = converted.Password
	modelUser.Username = converted.Username
	modelUser.AvatarURL = converted.AvatarURL
	modelUser.EmailVerifiedAt = converted.EmailVerifiedAt
	modelUser.LastLoginAt = converted.LastLoginAt
	modelUser.Role = converted.Role

	// 更新処理を実行
	tx = r.db.Save(&modelUser)
	if tx.Error != nil {
		return nil, errs.NewInfraError(fmt.Errorf("ユーザー更新に失敗しました: %w", tx.Error).Error())
	}

	// 更新後の永続化用モデルからドメインエンティティに再構築
	updatedUser, err := adapter.NewUserAdapter().ReBuild(&modelUser)
	if err != nil {
		return nil, errs.NewInfraError(fmt.Errorf("ユーザーエンティティの再構築に失敗しました: %w", err).Error())
	}
	return updatedUser, nil
}

func (r *UserRepositoryImpl) DeleteUser(objID string) error {
	tx := r.db.Where("obj_id = ?", objID).Delete(&entity.User{})
	if tx.Error != nil {
		return errs.NewInfraError(fmt.Errorf("オブジェクトID(%s)のユーザー削除に失敗しました: %w", objID, tx.Error).Error())
	}
	return nil
}
