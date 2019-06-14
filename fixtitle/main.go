package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	dirPath string
)

func init() {
	// 处理的目录 -p 参数 如果不加则默认处理当前目录
	flag.StringVar(&dirPath, "p", ".", "dir path")
}

func main() {
	// 读取json配置文件
	flag.Parse()
	fmt.Println("will to direcotry:", dirPath)

	// 替换目录下面的内容
	work(dirPath)
}

func work(dir string) {
	// 遍历目录 读取每一个文件
	filepath.Walk(dir, walkFunc)
}

func walkFunc(path string, info os.FileInfo, err error) error {
	var (
		ok         bool
		content    []byte
		pat        string
		reg        *regexp.Regexp
		fileName   string
		reWriteFlg bool
	)
	// 如果不是.md结尾的文件则不进行处理
	if ok = strings.HasSuffix(info.Name(), ".md"); !ok {
		return nil
	}

	fileName = strings.SplitN(info.Name(), ".md", -1)[0]
	if fileName == "" {
		return nil
	}

	if content, err = ioutil.ReadFile(path); err != nil {
		return err
	}

	// 处理title:
	// 将文件名称替换为title
	// 如果原先有title则不处理
	pat = `title:\s+(.*)\n`
	reg = regexp.MustCompile(pat)
	if ok = reg.Match(content); !ok {
		// 找不到title 则在slug前面增加title:文件名称
		pat = `slug:.*\n`
		reg = regexp.MustCompile(pat)
		content = reg.ReplaceAll(content, []byte("title: \""+fileName+"\"\n$0"))
		reWriteFlg = true
	}

	// 处理category的格式错误
	// category: "go" 改为：
	// categories:
	// - Development
	// - VIM
	pat = `category:\s+"(.*)"\n`
	reg = regexp.MustCompile(pat)
	if ok = reg.Match(content); ok {
		// 找到category内容
		content = reg.ReplaceAll(content, []byte("categories:\n    - $1\n"))
		reWriteFlg = true
	}

	if reWriteFlg {
		fmt.Println("已处理文件：", info.Name())
		return ioutil.WriteFile(path, content, info.Mode().Perm())
	}

	return nil
}
