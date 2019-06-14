package main

import (
	"fmt"
	"regexp"
)

func main() {
	// 判断在 b（s、r）中能否找到 pattern 所匹配的字符串
	// func Match(pattern string, b []byte) (matched bool, err error)
	// func MatchString(pattern string, s string) (matched bool, err error)
	// func MatchReader(pattern string, r io.RuneReader) (matched bool, err error)

	// 将 s 中的正则表达式元字符转义成普通字符。
	// func QuoteMeta(s string) string
	var (
		pat string
		src string
	)

	// . 匹配任意一个字符，如果设置s = ture， 则可以匹配换行符
	// (子表达式）被捕获的组，该组被编号(子匹配)

	pat = `(((abc.)def.)ghi)`
	src = `abc-def-ghi abc+def+ghi`
	fmt.Println(regexp.MatchString(pat, src))
	// ture nil

	fmt.Println(regexp.QuoteMeta(pat))
	// \(\(\(abc\.\)def\.\)ghi\)

	b := []byte("abc1def1")
	pat = `abc1|abc1def1`
	reg1 := regexp.MustCompile(pat)
	reg2 := regexp.MustCompilePOSIX(pat)
	fmt.Printf("%s\n", reg1.Find(b)) // abc1
	fmt.Printf("%s\n", reg2.Find(b)) // abc1def1

	pat = `(((abc.)def.)ghi)`
	reg := regexp.MustCompile(pat)
	srcByte := []byte(`abc-def-ghi abc+def+ghi`)

	// 查找第一个匹配结果
	fmt.Printf("%s\n", reg.Find(srcByte))
	// abc-def-ghi

	// 查看第一个匹配结果及其分组字符串
	first := reg.FindSubmatch(srcByte)
	for i := 0; i < len(first); i++ {
		fmt.Printf("%d:%s\n", i, first[i])
	}

}
