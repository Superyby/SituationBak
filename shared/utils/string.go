package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"unicode/utf8"
)

// IsEmpty 检查字符串是否为空（去除空白后）
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty 检查字符串是否非空
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// DefaultIfEmpty 如果字符串为空则返回默认值
func DefaultIfEmpty(s, defaultVal string) string {
	if IsEmpty(s) {
		return defaultVal
	}
	return s
}

// TruncateString 截断字符串到指定长度
// 支持 UTF-8 字符串，按字符数截断
func TruncateString(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxLen])
}

// TruncateWithEllipsis 截断字符串并添加省略号
func TruncateWithEllipsis(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return "..."
	}
	runes := []rune(s)
	return string(runes[:maxLen-3]) + "..."
}

// ContainsAny 检查字符串是否包含任意一个子串
func ContainsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// ContainsAll 检查字符串是否包含所有子串
func ContainsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

// RemoveWhitespace 移除所有空白字符
func RemoveWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			return -1
		}
		return r
	}, s)
}

// NormalizeWhitespace 将连续空白字符替换为单个空格
func NormalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// MaskString 掩码字符串（用于敏感信息）
// 例如: MaskString("13812345678", 3, 4, '*') => "138****5678"
func MaskString(s string, prefixLen, suffixLen int, maskChar rune) string {
	runes := []rune(s)
	length := len(runes)

	if length <= prefixLen+suffixLen {
		return s
	}

	masked := make([]rune, length)
	copy(masked[:prefixLen], runes[:prefixLen])
	copy(masked[length-suffixLen:], runes[length-suffixLen:])

	for i := prefixLen; i < length-suffixLen; i++ {
		masked[i] = maskChar
	}

	return string(masked)
}

// MaskEmail 掩码邮箱地址
// 例如: MaskEmail("user@example.com") => "u***@example.com"
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	name := parts[0]
	domain := parts[1]

	if len(name) <= 1 {
		return email
	}

	return string(name[0]) + "***@" + domain
}

// MaskPhone 掩码手机号
// 例如: MaskPhone("13812345678") => "138****5678"
func MaskPhone(phone string) string {
	return MaskString(phone, 3, 4, '*')
}

// RandomString 生成随机字符串
func RandomString(length int) string {
	bytes := make([]byte, length/2+1)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// SplitAndTrim 分割字符串并去除每个元素的空白
func SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// JoinNonEmpty 连接非空字符串
func JoinNonEmpty(sep string, strs ...string) string {
	nonEmpty := make([]string, 0, len(strs))
	for _, s := range strs {
		if s != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	return strings.Join(nonEmpty, sep)
}

// FirstNonEmpty 返回第一个非空字符串
func FirstNonEmpty(strs ...string) string {
	for _, s := range strs {
		if s != "" {
			return s
		}
	}
	return ""
}

// PadLeft 左侧填充字符串
func PadLeft(s string, length int, pad rune) string {
	runeCount := utf8.RuneCountInString(s)
	if runeCount >= length {
		return s
	}
	return strings.Repeat(string(pad), length-runeCount) + s
}

// PadRight 右侧填充字符串
func PadRight(s string, length int, pad rune) string {
	runeCount := utf8.RuneCountInString(s)
	if runeCount >= length {
		return s
	}
	return s + strings.Repeat(string(pad), length-runeCount)
}

// Reverse 反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
