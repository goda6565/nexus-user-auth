package timeobj

import (
	"fmt"
	"time"
)

// TimeObj: 時刻を表すオブジェクト
type TimeObj struct {
	value time.Time
}

func (t *TimeObj) Value() time.Time {
	return t.value
}

// 未来の日付は許容せず、未来の場合はエラーを返す。
func NewTimeObj(value time.Time) (*TimeObj, error) {
	now := time.Now()
	if value.After(now) {
		return nil, fmt.Errorf("未来の日付は許容されません: %v", value)
	}
	return &TimeObj{value: value}, nil
}
