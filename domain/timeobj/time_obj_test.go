package timeobj

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTimeObj(t *testing.T) {
	// 現在時刻を元に TimeObj を生成する
	now := time.Now()
	timeObj, err := NewTimeObj(now)
	if err != nil {
		t.Fatal(err)
	}

	// 生成した TimeObj の値が正しいことを検証する
	assert.Equal(t, now, timeObj.Value())
}

func TestNewTimeObj_Future(t *testing.T) {
	// 未来の日付を指定して TimeObj を生成しようとするとエラーが返ることを検証する
	future := time.Now().Add(time.Hour)
	_, err := NewTimeObj(future)
	assert.NotNil(t, err)
}
