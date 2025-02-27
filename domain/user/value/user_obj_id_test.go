package value

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUserObjID_Valid(t *testing.T) {
	validUUID := uuid.New().String() // 36文字のUUID形式
	objID, err := NewUserObjID(validUUID)
	assert.NoError(t, err, "正しいUUIDの場合、エラーが返られないこと")
	assert.NotNil(t, objID)
	assert.Equal(t, validUUID, objID.Value())
}

func TestNewUserObjID_InvalidLength(t *testing.T) {
	// 36文字でない値
	invalid := "1234"
	objID, err := NewUserObjID(invalid)
	assert.Error(t, err, "長さが不正な場合、エラーが返ること")
	assert.Nil(t, objID)
}

func TestNewUserObjID_InvalidFormat(t *testing.T) {
	// 36文字でもフォーマットがUUIDでない場合
	invalid := "123456787123471234712347123456789012" // 形式上正しく見えるが数字だけで検証
	objID, err := NewUserObjID(invalid)
	// ここは実際の正規表現に合わせて期待値を決定してください
	assert.Error(t, err, "UUID形式でない場合、エラーが返ること")
	assert.Nil(t, objID)
}
