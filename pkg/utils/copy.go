package utils

import (
	"log"
	"reflect"
)

// Copy 使用反射将 src 的字段复制到 dst 中
func Copy(src interface{}, dst interface{}) interface{} {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	// 检查 src 和 dst 是否为指针
	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		log.Println("src and dst must be pointers")
		return nil
	}

	// 解析到它们的值
	srcElem := srcVal.Elem()
	dstElem := dstVal.Elem()

	// 检查 src 和 dst 是否是结构体
	if srcElem.Kind() != reflect.Struct || dstElem.Kind() != reflect.Struct {
		log.Println("src and dst must be struct pointers")
		return nil
	}

	// 检查 dst 是否可寻址和可设置
	if !dstElem.CanAddr() || !dstElem.CanSet() {
		log.Println("dst is not addressable or cannot be set")
		return nil
	}

	// 遍历 src 的字段
	for i := 0; i < srcElem.NumField(); i++ {
		srcField := srcElem.Field(i)
		dstField := dstElem.FieldByName(srcElem.Type().Field(i).Name)

		// 如果 dst 中有同名字段，并且它是可设置的
		if dstField.IsValid() && dstField.CanSet() {
			// 如果类型相同，直接赋值
			if srcField.Type() == dstField.Type() {
				dstField.Set(srcField)
			}
		}
	}

	return dst
}
