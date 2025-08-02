package utils

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Converter 接口定义类型转换器
type Converter interface {
	Convert(src reflect.Value) (interface{}, bool)
}

// PbTimestampConverter 将 time.Time 转换为 google.protobuf.Timestamp
type PbTimestampConverter struct{}

func (c PbTimestampConverter) Convert(src reflect.Value) (interface{}, bool) {
	// 处理 time.Time 类型（值类型）
	if src.Type() == reflect.TypeOf(time.Time{}) {
		t := src.Interface().(time.Time)
		if t.IsZero() {
			return &timestamp.Timestamp{}, true // 返回空Timestamp而非nil，避免指针零值
		}
		return &timestamp.Timestamp{
			Seconds: t.Unix(),
			Nanos:   int32(t.Nanosecond()),
		}, true
	}

	// 处理 primitive.DateTime（MongoDB 的时间类型）
	if src.Type() == reflect.TypeOf(primitive.DateTime(0)) {
		dt := src.Interface().(primitive.DateTime)
		t := dt.Time()
		return &timestamp.Timestamp{
			Seconds: t.Unix(),
			Nanos:   int32(t.Nanosecond()),
		}, true
	}

	// 处理 time.Time 指针类型（非 nil 时）
	if src.Type() == reflect.TypeOf(&time.Time{}) && !src.IsNil() {
		t := src.Interface().(*time.Time)
		return &timestamp.Timestamp{
			Seconds: t.Unix(),
			Nanos:   int32(t.Nanosecond()),
		}, true
	}

	return nil, false
}

// ObjectIDConverter 将 primitive.ObjectID 转换为 string
type ObjectIDConverter struct{}

func (c ObjectIDConverter) Convert(src reflect.Value) (interface{}, bool) {
	// 处理值类型
	if src.Type() == reflect.TypeOf(primitive.ObjectID{}) {
		return src.Interface().(primitive.ObjectID).Hex(), true
	}
	// 处理指针类型
	if src.Type() == reflect.TypeOf(&primitive.ObjectID{}) && !src.IsNil() {
		objID := src.Interface().(*primitive.ObjectID)
		return objID.Hex(), true
	}
	return nil, false
}

// TimeConverter 将 time.Time 转换为字符串
type TimeConverter struct{}

func (c TimeConverter) Convert(src reflect.Value) (interface{}, bool) {
	if src.Type() == reflect.TypeOf(time.Time{}) {
		return src.Interface().(time.Time).Format("2006-01-02 15:04:05"), true
	}
	if src.Type() == reflect.TypeOf(primitive.DateTime(0)) {
		t := src.Interface().(primitive.DateTime).Time()
		return t.Format("2006-01-02 15:04:05"), true
	}
	return nil, false
}

// MapOption 映射配置选项
type MapOption struct {
	FieldMappings map[string]string // 目标字段名 -> 源字段名
}

// Map 结构体映射函数
func StructToPB(src, dst interface{}, converters []Converter, opts ...MapOption) error {
	fieldMappings := make(map[string]string)
	for _, opt := range opts {
		if opt.FieldMappings != nil {
			for k, v := range opt.FieldMappings {
				fieldMappings[k] = v
			}
		}
	}

	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	// 检查源和目标是否为指针
	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		return fmt.Errorf("src and dst must be pointers (got src: %T, dst: %T)", src, dst)
	}

	// 检查指针是否为nil
	if srcVal.IsNil() || dstVal.IsNil() {
		return fmt.Errorf("src and dst cannot be nil pointers")
	}

	// 解引用指针
	srcVal = srcVal.Elem()
	dstVal = dstVal.Elem()

	// 检查是否为结构体
	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		return fmt.Errorf("src and dst must point to structs (got src: %T, dst: %T)", srcVal.Interface(), dstVal.Interface())
	}

	// 遍历目标结构体字段
	for i := 0; i < dstVal.NumField(); i++ {
		dstField := dstVal.Field(i)
		if !dstField.CanSet() {
			continue // 跳过不可设置的字段（如未导出字段）
		}

		dstFieldName := dstVal.Type().Field(i).Name
		srcFieldName := dstFieldName
		if mappedName, ok := fieldMappings[dstFieldName]; ok {
			srcFieldName = mappedName // 使用字段映射
		}

		// 获取源结构体对应字段
		srcField := srcVal.FieldByName(srcFieldName)
		if !srcField.IsValid() {
			continue // 源字段不存在则跳过
		}

		// 处理嵌套结构体（值类型）
		if srcField.Kind() == reflect.Struct && dstField.Kind() == reflect.Struct {
			nestedDst := reflect.New(dstField.Type()).Elem() // 初始化嵌套结构体
			if err := StructToPB(srcField.Addr().Interface(), nestedDst.Addr().Interface(), converters, opts...); err == nil {
				dstField.Set(nestedDst)
			}
			continue
		}

		// 处理嵌套结构体指针
		if srcField.Kind() == reflect.Ptr && !srcField.IsNil() && srcField.Elem().Kind() == reflect.Struct {
			if dstField.Kind() == reflect.Ptr {
				// 初始化目标指针（若为nil）
				if dstField.IsNil() {
					dstField.Set(reflect.New(dstField.Type().Elem()))
				}
				// 递归映射嵌套结构体
				if err := StructToPB(srcField.Interface(), dstField.Interface(), converters, opts...); err != nil {
					continue // 映射失败则跳过
				}
				continue
			}
		}

		// 处理切片
		if srcField.Kind() == reflect.Slice && dstField.Kind() == reflect.Slice {
			// 初始化目标切片
			dstSlice := reflect.MakeSlice(dstField.Type(), srcField.Len(), srcField.Cap())

			for j := 0; j < srcField.Len(); j++ {
				srcElem := srcField.Index(j)
				dstElem := dstSlice.Index(j)

				// 处理切片中的指针元素
				realSrcElem := srcElem
				if srcElem.Kind() == reflect.Ptr {
					if srcElem.IsNil() {
						continue // 跳过nil指针
					}
					realSrcElem = srcElem.Elem()
				}

				// 切片元素为结构体时递归映射
				if realSrcElem.Kind() == reflect.Struct {
					// 目标元素为结构体
					if dstElem.Kind() == reflect.Struct {
						nestedDst := reflect.New(dstElem.Type()).Elem()
						if err := StructToPB(realSrcElem.Addr().Interface(), nestedDst.Addr().Interface(), converters, opts...); err == nil {
							dstElem.Set(nestedDst)
						}
					} else if dstElem.Kind() == reflect.Ptr {
						// 目标元素为结构体指针
						nestedDst := reflect.New(dstElem.Type().Elem())
						if err := StructToPB(realSrcElem.Addr().Interface(), nestedDst.Interface(), converters, opts...); err == nil {
							dstElem.Set(nestedDst)
						}
					}
				} else {
					// 非结构体元素直接赋值（类型兼容）
					if realSrcElem.Type().AssignableTo(dstElem.Type()) {
						dstElem.Set(realSrcElem)
					}
				}
			}

			dstField.Set(dstSlice)
			continue
		}

		// 尝试使用转换器转换字段
		converted := false
		if dstField.Kind() == reflect.Ptr {
			// 初始化指针类型目标字段（关键修复：确保指针非nil）
			if dstField.IsNil() {
				dstField.Set(reflect.New(dstField.Type().Elem()))
			}
			elem := dstField.Elem()
			if elem.CanSet() {
				for _, converter := range converters {
					if val, ok := converter.Convert(srcField); ok {
						// 处理转换后的值
						valRef := reflect.ValueOf(val)
						if valRef.Type().AssignableTo(elem.Type()) {
							elem.Set(valRef)
						} else if val != nil {
							// 尝试间接赋值（如转换后为指针，目标为值类型）
							if valRef.Kind() == reflect.Ptr && valRef.Elem().Type().AssignableTo(elem.Type()) {
								elem.Set(valRef.Elem())
							} else {
								return fmt.Errorf("converter result type %T not assignable to target type %T", val, elem.Interface())
							}
						}
						converted = true
						break
					}
				}
			}
		} else {
			// 非指针类型目标字段
			for _, converter := range converters {
				if val, ok := converter.Convert(srcField); ok {
					valRef := reflect.ValueOf(val)
					if valRef.Type().AssignableTo(dstField.Type()) {
						dstField.Set(valRef)
						converted = true
						break
					} else {
						return fmt.Errorf("converter result type %T not assignable to target type %T", val, dstField.Interface())
					}
				}
			}
		}

		if converted {
			continue
		}

		// 类型相同直接赋值
		if srcField.Type().AssignableTo(dstField.Type()) {
			dstField.Set(srcField)
		}
	}

	return nil
}

// ListToPB 结构体列表映射到 Protobuf 消息列表
func ListToPB(srcList, dstList interface{}, converters []Converter, opts ...MapOption) error {
	srcVal := reflect.ValueOf(srcList)
	dstVal := reflect.ValueOf(dstList)

	// 检查源和目标是否为指针
	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		return fmt.Errorf("srcList and dstList must be pointers (got srcList: %T, dstList: %T)", srcList, dstList)
	}

	// 检查指针是否为nil
	if srcVal.IsNil() {
		// 根据错误信息，MapOption 类型没有 DefaultIfNil 字段，此处暂时移除该条件判断
		// 若需要实现此功能，需先在 MapOption 结构体中添加 DefaultIfNil 字段
		if false {
			// 源为nil时初始化目标为零长度切片
			dstVal.Elem().Set(reflect.MakeSlice(dstVal.Type().Elem(), 0, 0))
			return nil
		}
		return fmt.Errorf("srcList cannot be nil pointer")
	}

	if dstVal.IsNil() {
		return fmt.Errorf("dstList cannot be nil pointer")
	}

	// 解引用指针
	srcVal = srcVal.Elem()
	dstVal = dstVal.Elem()

	// 检查是否为切片或数组
	if srcVal.Kind() != reflect.Slice && srcVal.Kind() != reflect.Array {
		return fmt.Errorf("srcList must be a slice or array (got %T)", srcList)
	}

	if dstVal.Kind() != reflect.Slice {
		return fmt.Errorf("dstList must be a slice (got %T)", dstList)
	}

	// 创建目标切片
	dstSlice := reflect.MakeSlice(dstVal.Type(), srcVal.Len(), srcVal.Len())
	dstVal.Set(dstSlice)

	// 获取元素类型

	// 遍历源列表
	for i := 0; i < srcVal.Len(); i++ {
		srcElem := srcVal.Index(i)
		dstElem := dstVal.Index(i)

		// 处理源元素为指针的情况
		if srcElem.Kind() == reflect.Ptr {
			if srcElem.IsNil() {
				// 由于 MapOption 类型没有 DefaultIfNil 字段，移除该条件判断
				if false {
					// 创建目标元素的零值
					if dstElem.Kind() == reflect.Ptr {
						dstElem.Set(reflect.New(dstElem.Type().Elem()))
					}
					continue
				}
				return fmt.Errorf("srcList element at index %d is nil", i)
			}
			srcElem = srcElem.Elem()
		}

		// 确保目标元素是指针
		if dstElem.Kind() != reflect.Ptr {
			return fmt.Errorf("dstList element must be a pointer (got %v)", dstElem.Type())
		}

		// 初始化目标元素
		if dstElem.IsNil() {
			dstElem.Set(reflect.New(dstElem.Type().Elem()))
		}

		// 映射单个元素
		if err := StructToPB(srcElem.Addr().Interface(), dstElem.Interface(), converters, opts...); err != nil {
			return fmt.Errorf("failed to map element at index %d: %v", i, err)
		}
	}

	return nil
}

func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)
	// 获取结构体的类型和值
	if val.Kind() == reflect.Struct {
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			result[typ.Field(i).Name] = field.Interface()
		}
	}
	return result
}

func MapToStruct(m map[string]interface{}, obj interface{}) error {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem() // 获取结构体的值
		for key, value := range m {
			fieldVal := val.FieldByName(key)
			if fieldVal.IsValid() && fieldVal.CanSet() {
				valFieldType := fieldVal.Type()
				valFieldValue := reflect.ValueOf(value)
				if valFieldValue.Type().AssignableTo(valFieldType) {
					fieldVal.Set(valFieldValue)
				}
			}
		}
		return nil
	}
	return fmt.Errorf("provided object is not a pointer to a struct")
}

// PbToMap 将 Protobuf 消息转换为更新用的 map
func PbToMap(pb proto.Message) map[string]interface{} {
	result := make(map[string]interface{})
	ref := pb.ProtoReflect()

	ref.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		tagName := string(fd.Name()) // 使用 JSON 标签名

		// 处理各种字段类型
		switch fd.Kind() {
		case protoreflect.MessageKind:
			if fd.Message().FullName() == "google.protobuf.Timestamp" {
				if fd.IsList() {
					// 处理 repeated google.protobuf.Timestamp 字段
					list := v.List()
					values := make([]string, list.Len())
					for i := 0; i < list.Len(); i++ {
						msg := list.Get(i).Message().Interface()
						if ts, ok := msg.(*timestamppb.Timestamp); ok {
							values[i] = ts.AsTime().Format(time.RFC3339)
						} else {
							values[i] = ""
						}
					}
					result[tagName] = values
				} else {
					msg := v.Message().Interface()
					if ts, ok := msg.(*timestamppb.Timestamp); ok {
						result[tagName] = ts.AsTime().Format(time.RFC3339)
					} else {
						result[tagName] = nil
					}
				}
				return true
			} else if fd.IsList() {
				// 处理 repeated 字段
				list := v.List()
				values := make([]interface{}, list.Len())
				for i := 0; i < list.Len(); i++ {
					item := list.Get(i)
					if fd.Message().FullName() == "google.protobuf.Timestamp" {
						if msg := item.Message().Interface(); msg != nil {
							if ts, ok := msg.(*timestamppb.Timestamp); ok {
								values[i] = ts.AsTime().Format(time.RFC3339)
							} else {
								values[i] = nil
							}
						} else {
							values[i] = nil
						}
					} else if nestedMsg := item.Message().Interface(); nestedMsg != nil {
						values[i] = PbToMap(nestedMsg.(proto.Message))
					} else {
						values[i] = nil
					}
				}
				result[tagName] = values
			} else {
				// 处理嵌套消息
				if nestedMsg := v.Message().Interface(); nestedMsg != nil {
					result[tagName] = PbToMap(nestedMsg.(proto.Message))
				} else {
					result[tagName] = nil
				}
			}
		case protoreflect.StringKind:
			if fd.IsList() {
				// 处理 repeated string 字段
				list := v.List()
				values := make([]string, list.Len())
				for i := 0; i < list.Len(); i++ {
					values[i] = list.Get(i).String()
				}
				result[tagName] = values
			} else {
				result[tagName] = v.String()
			}
		case protoreflect.EnumKind:
			if fd.IsList() {
				// 处理 repeated enum 字段
				list := v.List()
				values := make([]int32, list.Len())
				for i := 0; i < list.Len(); i++ {
					values[i] = int32(list.Get(i).Enum())
				}
				result[tagName] = values
			} else {
				result[tagName] = int32(v.Enum())
			}
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			if fd.IsList() {
				// 处理 repeated int32 类型字段
				list := v.List()
				values := make([]int64, list.Len())
				for i := 0; i < list.Len(); i++ {
					values[i] = int64(list.Get(i).Int())
				}
				result[tagName] = values
			} else {
				result[tagName] = int64(v.Int())
			}
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			if fd.IsList() {
				// 处理 repeated uint32 类型字段
				list := v.List()
				values := make([]int64, list.Len())
				for i := 0; i < list.Len(); i++ {
					values[i] = int64(list.Get(i).Uint())
				}
				result[tagName] = values
			} else {
				result[tagName] = int64(v.Uint())
			}
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			if fd.IsList() {
				// 处理 repeated int64 类型字段
				list := v.List()
				values := make([]int64, list.Len())
				for i := 0; i < list.Len(); i++ {
					values[i] = list.Get(i).Int()
				}
				result[tagName] = values
			} else {
				result[tagName] = v.Int()
			}
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			if fd.IsList() {
				// 处理 repeated uint64 类型字段
				list := v.List()
				values := make([]int64, list.Len())
				for i := 0; i < list.Len(); i++ {
					values[i] = int64(list.Get(i).Uint())
				}
				result[tagName] = values
			} else {
				result[tagName] = int64(v.Uint())
			}
		default:
			result[tagName] = v.Interface()
		}
		return true
	})

	return result
}

func MergerToString(src []string) string {
	var tar string
	for _, v := range src {
		tar += v + ";"
	}
	return tar
}

func SplitToSlice(src string) []string {
	parts := strings.Split(src, ";")
	// 创建一个新的切片来存储处理后的结果
	var result []string
	for _, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart != "" {
			result = append(result, trimmedPart)
		}
	}
	return result
}

// 判断俩个集合是否元素相等
// 时间复杂度是 O(max(n, m))
func AreSetsEqual(set1, set2 []string) bool {
	// 如果两个集合长度不相同，则直接返回 false
	if len(set1) != len(set2) {
		return false
	}

	// 使用 map 来统计每个字符串出现的次数
	count1 := make(map[string]int)
	count2 := make(map[string]int)

	// 统计第一个集合的元素出现次数
	for _, str := range set1 {
		count1[str]++
	}

	// 统计第二个集合的元素出现次数
	for _, str := range set2 {
		count2[str]++
	}

	// 比较两个集合的元素统计结果
	for key, val := range count1 {
		if count2[key] != val {
			return false
		}
	}

	return true
}

// 集合 B - A
// A，B两个集合，只找出B在A中不存在的字符串集合
func DifferenceBInA(A, B []string) []string {
	// 创建一个哈希集合（map）用于快速查找 A 中的元素
	setAMap := make(map[string]struct{})

	// 将 A 中的元素存入 setAMap
	for _, val := range A {
		setAMap[val] = struct{}{}
	}

	// 创建一个结果切片，用于存储 B 中不在 A 中的元素
	var diff []string

	// 遍历 B，找出不在 A 中的元素
	for _, val := range B {
		if _, found := setAMap[val]; !found {
			diff = append(diff, val)
		}
	}

	return diff
}
