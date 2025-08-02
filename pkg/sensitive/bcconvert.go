package sensitive

// ASCII表中可见字符从!开始，偏移位值为33(Decimal)
const DBC_CHAR_START = 33 // 半角!

// ASCII表中可见字符到~结束，偏移位值为126(Decimal)
const DBC_CHAR_END = 126 // 半角~

// 全角对应于ASCII表的可见字符从！开始，偏移值为65281
const SBC_CHAR_START = 65281 // 全角！

// 全角对应于ASCII表的可见字符到～结束，偏移值为65374
const SBC_CHAR_END = 65374 // 全角～

const CONVERT_STEP = 65248 // 全角半角转换间隔

// 全角空格的值，它没有遵从与ASCII的相对偏移，必须单独处理
const SBC_SPACE = 12288 // 全角空格 12288

// 半角空格的值，在ASCII中为32(Decimal)
const DBC_SPACE = ' ' // 半角空格

// Bj2qj 半角字符->全角字符转换
// 只处理空格，!到~之间的字符，忽略其他
func Bj2qj(src string) string {
	if src == "" {
		return src
	}

	buf := make([]rune, 0, len(src))
	for _, char := range src {
		if char == DBC_SPACE { // 如果是半角空格，直接用全角空格替代
			buf = append(buf, SBC_SPACE)
		} else if char >= DBC_CHAR_START && char <= DBC_CHAR_END { // 字符是!到~之间的可见字符
			buf = append(buf, char+CONVERT_STEP)
		} else { // 不对空格以及ascii表中其他可见字符之外的字符做任何处理
			buf = append(buf, char)
		}
	}

	return string(buf)
}

// Bj2qjChar 半角转换全角（字符级别）
func Bj2qjChar(char rune) rune {
	if char == DBC_SPACE { // 如果是半角空格，直接用全角空格替代
		return SBC_SPACE
	} else if char >= DBC_CHAR_START && char <= DBC_CHAR_END { // 字符是!到~之间的可见字符
		return char + CONVERT_STEP
	}
	return char
}

// Qj2bj 全角字符->半角字符转换
// 只处理全角的空格，全角！到全角～之间的字符，忽略其他
func Qj2bj(src string) string {
	if src == "" {
		return src
	}

	buf := make([]rune, 0, len(src))
	for _, char := range src {
		if char >= SBC_CHAR_START && char <= SBC_CHAR_END { // 如果位于全角！到全角～区间内
			buf = append(buf, char-CONVERT_STEP)
		} else if char == SBC_SPACE { // 如果是全角空格
			buf = append(buf, DBC_SPACE)
		} else { // 不处理全角空格，全角！到全角～区间外的字符
			buf = append(buf, char)
		}
	}

	return string(buf)
}

// Qj2bjChar 全角转换半角（字符级别）
func Qj2bjChar(char rune) rune {
	if char >= SBC_CHAR_START && char <= SBC_CHAR_END { // 如果位于全角！到全角～区间内
		return char - CONVERT_STEP
	} else if char == SBC_SPACE { // 如果是全角空格
		return DBC_SPACE
	}
	return char
}
