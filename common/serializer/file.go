package serializer

import (
	"context"
	"gorm.io/gorm/schema"
	"reflect"
)

type FileSerializer struct{}

// Value 将二进制文件保存到指定位置，返回文件路径
func (f *FileSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {

	return nil, nil
}

// Scan 读取的文件路径，拼接成 http 访问地址
func (f *FileSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error {

	return nil
}
