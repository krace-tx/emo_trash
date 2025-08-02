package sensitive

import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/Lofanmi/pinyin-golang/pinyin"
	"github.com/zeromicro/go-zero/core/jsonx"
	"strings"
)

// 定义常量和全局变量
const SIGN = '*' // 敏感词过滤替换

var (
	stopwdSet   = make(map[int]struct{})  // 停顿词
	softwdSet   = make(map[int]*WordNode) // 软文本敏感词
	strongwdSet = make(map[int]*WordNode) // 硬文本敏感词
)

const (
	Soft   = "soft"
	Strong = "strong"
)

//go:embed softwd.txt
var softWords []byte

//go:embed strongwd.txt
var strongWords []byte

//go:embed stopwd.txt
var stopWords []byte

// 初始化，加载敏感词和停顿词
func init() {
	// 加载软文本敏感词
	err := loadSensitiveWords(softWords, Soft)
	if err != nil {
		fmt.Println("Failed to load soft words:", err)
	}

	// 加载硬文本敏感词
	err = loadSensitiveWords(strongWords, Strong)
	if err != nil {
		fmt.Println("Failed to load strong words:", err)
	}

	// 加载停顿词
	err = loadStopWords(stopWords)
	if err != nil {
		fmt.Println("Failed to load stop words:", err)
	}
}

// 加载敏感词字节数据，支持不同类型的敏感词
func loadSensitiveWords(data []byte, wordType string) error {
	words := readWordsFromData(data)
	if wordType == Soft {
		addSensitiveWordToSet(words, softwdSet)
	} else if wordType == Strong {
		addSensitiveWordToSet(words, strongwdSet)
	} else {
		return fmt.Errorf("unknown word type: %s", wordType)
	}
	return nil
}

// 加载停顿词字节数据
func loadStopWords(data []byte) error {
	words := readWordsFromData(data)
	addStopWord(words)
	return nil
}

// 读取字节数据
func readWordsFromData(data []byte) []string {
	var words []string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	// 使用 FieldsFunc 来自定义分隔符处理逻辑
	for scanner.Scan() {
		line := scanner.Text()
		// 使用自定义分隔符，处理半角空格和全角空格
		words = append(words, replaceSpaces(line))
	}
	return words
}

func replaceSpaces(src string) string {
	// 去除全角空格（Unicode: \u3000）
	src = strings.ReplaceAll(src, "　", "！")

	// 去除半角空格（Unicode: \u0020）
	src = strings.ReplaceAll(src, " ", "_")

	return src
}

// 增加停顿词
func addStopWord(words []string) {
	for _, word := range words {
		for _, char := range word {
			stopwdSet[charConvert(char)] = struct{}{}
		}
	}
}

// 增加敏感词到指定的集合
func addSensitiveWordToSet(words []string, targetSet map[int]*WordNode) {
	for _, word := range words {
		chs := []rune(word)
		if len(chs) == 0 {
			continue
		}
		fchar := charConvert(chs[0])
		var fnode *WordNode
		if node, exists := targetSet[fchar]; !exists {
			fnode = NewWordNode(fchar, len(chs) == 1)
			targetSet[fchar] = fnode
		} else {
			fnode = node
		}
		lastIndex := len(chs) - 1
		for i := 1; i < len(chs); i++ {
			fnode = fnode.AddIfNoExist(charConvert(chs[i]), i == lastIndex)
		}
	}
}

// 过滤敏感词，将敏感词替换为特定字符
func DoFilter(src string, wordType string) string {
	src = replaceSpaces(src)
	var targetSet map[int]*WordNode
	if wordType == Soft {
		targetSet = softwdSet
	} else if wordType == Strong {
		targetSet = strongwdSet
	} else {
		return src // 如果类型不对，直接返回原始文本
	}

	chs := []rune(src)
	length := len(chs)

	for i := 0; i < length; i++ {
		currc := charConvert(chs[i])
		if _, exists := targetSet[currc]; !exists {
			continue
		}

		node := targetSet[currc]
		if node == nil {
			continue
		}

		couldMark := false
		markNum := -1

		if node.IsLast() {
			couldMark = true
			markNum = 0
		}

		k := i
		cpcurrc := currc

		for k++; k < length; k++ {
			temp := charConvert(chs[k])
			if temp == cpcurrc || isStopWord(temp) {
				continue
			}
			node = node.QuerySub(temp)
			if node == nil {
				break
			}
			if node.IsLast() {
				couldMark = true
				markNum = k - i
			}
			cpcurrc = temp
		}

		if couldMark {
			for k := 0; k <= markNum; k++ {
				chs[i+k] = SIGN
			}
			i += markNum
		}
	}

	return string(chs)
}

// 审核文本是否通过
func Verify(src string) bool {
	dict := pinyin.NewDict()
	fmt.Println(dict.Sentence(src).Unicode())
	fmt.Println(verify(dict.Sentence(src).Unicode()))

	fmt.Println(dict.Name(src, "-").Unicode())
	fmt.Println(verify(dict.Name(src, "-").Unicode()))

	fmt.Println(src)
	fmt.Println(verify(src))

	return verify(dict.Sentence(src).Unicode()) &&
		verify(dict.Name(src, "-").Unicode()) &&
		verify(src)
}

func verify(src string) bool {
	src = replaceSpaces(src)
	// 软过滤
	softFilteredCount := countReplacedChars(src)
	softFilteredRatio := float64(softFilteredCount) / float64(len(src))

	// 如果软过滤后超过50%被替换，则直接返回false
	if softFilteredRatio >= 0.5 {
		return false
	}

	// 硬过滤，直接判断是否有被过滤的内容
	if strongFilter(src) {
		return false // 硬过滤有被过滤的内容，返回false
	}

	// 如果没有超过50%被替换，且硬过滤没有内容被替换，则返回true
	return true
}

// 计算滤替换的字符数量
func countReplacedChars(src string) int {
	soft, _ := jsonx.Marshal(DoFilter(src, Soft))
	marshal, _ := jsonx.Marshal(src)

	cnt := 0
	for i, v := range soft {
		if marshal[i] != v {
			cnt++
		}
	}

	return cnt
}

// 判断是否是停顿词
func isStopWord(r int) bool {
	_, exists := stopwdSet[r]
	return exists
}

// 大写转小写，全角转半角
func charConvert(src rune) int {
	r := Qj2bjChar(src)
	if r >= 'A' && r <= 'Z' {
		return int(r + 32)
	}
	return int(r)
}

// 判断文本是否包含敏感词（可以通过 DoFilter 实现）
func strongFilter(src string) bool {
	return len(DoFilter(src, Strong)) < len(src)
}
