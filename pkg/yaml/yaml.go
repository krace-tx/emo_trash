package yaml

import (
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"strconv"
	"time"
)

// LoadAndMergeFromFile 从 JSON 文件加载配置并合并到目标结构体中
func LoadYamlFile(filePath string, target interface{}) error {
	// 读取文件内容
	jsonContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 调用 LoadAndMerge 函数处理 JSON 内容
	return LoadYaml(jsonContent, target)
}

// LoadYaml 从 YAML 字符串加载配置并合并到目标结构体中
func LoadYaml(yamlContent []byte, target interface{}) error {
	// 解析 YAML 到 map
	var data map[string]interface{}
	if err := yaml.Unmarshal(yamlContent, &data); err != nil {
		return err
	}

	// 创建自定义的 mapstructure 解码配置
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           target,
		TagName:          "json", // 使用结构体字段的 `json` 标签
		WeaklyTypedInput: true,   // 允许弱类型转换，例如字符串转整数
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToTypeHookFunc(),                      // 自定义字符串类型转换
			mapstructure.StringToTimeDurationHookFunc(), // 字符串到 time.Duration 转换
			mapstructure.StringToSliceHookFunc(","),     // 字符串到切片的转换，使用逗号分隔
		),
	}

	// 使用 mapstructure 解码，忽略未定义的字段
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if err := decoder.Decode(data); err != nil {
		return err
	}

	return nil
}

// stringToTypeHookFunc 自定义的字符串类型转换钩子函数
func stringToTypeHookFunc() mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		// 仅处理字符串类型的输入
		if from.Kind() != reflect.String {
			return data, nil
		}

		str := data.(string)

		// 尝试转换为布尔值
		if boolVal, err := strconv.ParseBool(str); err == nil {
			return boolVal, nil
		}

		// 尝试转换为整数
		if intVal, err := strconv.Atoi(str); err == nil {
			return intVal, nil
		}

		// 尝试转换为浮点数
		if floatVal, err := strconv.ParseFloat(str, 64); err == nil {
			return floatVal, nil
		}

		// 尝试转换为时间类型
		if timeVal, err := time.Parse(time.RFC3339, str); err == nil {
			return timeVal, nil
		}

		// 如果所有转换都失败，则返回原始字符串
		return str, nil
	}
}
