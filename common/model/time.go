package model

import (
	"context"
	"fmt"
	"gorm.io/gorm/schema"
	"reflect"
	"time"
)

type DateTime struct {
	time.Time
}

func (t *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format("2006-01-02") + `"`), nil
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		return nil
	}
	parsedTime, err := time.Parse(`"2006-01-02"`, str)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}

// Value adapts for GORM to handle the value of a custom field type
func (t DateTime) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Format("2006-01-02"), nil
}

// Scan 方法适用于 GORM 处理自定义字段类型的扫描
func (t *DateTime) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, value interface{}) error {
	if value == nil {
		return nil
	}

	var newTime time.Time
	var err error

	switch v := value.(type) {
	case time.Time:
		newTime = v
	case []byte:
		newTime, err = time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
	case string:
		newTime, err = time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("无法将 %T 转换为 DateTime", value)
	}
	fmt.Println(newTime)
	// 现在将时间设置到 DateTime 字段
	*t = DateTime{newTime}

	return nil
}
