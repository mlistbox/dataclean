package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Question struct {
	Id     int                    `json:"Id"`
	Head   string                 `json:"Head"`
	Answer string                 `json:"Answer"`
	Class  string                 `json:"Class"`
	Option map[string]interface{} `json:"Option"`
}

func NewQuestion(id int, head string, answer string, class string) *Question {
	return &Question{
		Id:     id,
		Head:   head,
		Answer: answer,
		Class:  class,
		Option: make(map[string]interface{}),
	}
}
func (q *Question) ToString() string {
	return fmt.Sprintf("questionID : %d ; Answer: %s class: %s ;Head: %s ;Option: %v\n", q.Id, q.Answer, q.Class, q.Head, q.Option)
}

func getFirstCharRange(s string) string {
	for _, r := range s {
		return string(r)
	}
	return ""
}
func StringToMap(q *Question, it string, re *regexp.Regexp) {
	strs := strings.Split(it, "；")
	for _, str := range strs {
		if re.MatchString(str) {
			//s := re.FindString(str)
			sl := strings.Split(str, "、")
			q.Option[getFirstCharRange(sl[0])] = sl[1]

		}
	}

}
func extractWithIndex(s string, st rune, et rune) string {
	start := strings.IndexRune(s, st)
	end := strings.IndexRune(s, et)

	// 确保括号存在且顺序正确
	if start == -1 || end == -1 || end <= start {
		return ""
	}

	return s[start+utf8.RuneLen(st) : end]
}

// SearchFileByLine 逐行搜索文件，返回匹配的行号和內容
func SearchFileByLine(ls *list.List, filename, pattern1 string, pattern2 string) error {
	// 1. 预编译正则
	re, err := regexp.Compile(pattern1)
	if err != nil {
		return fmt.Errorf("invalid regex1: %w", err)
	}

	re2, err := regexp.Compile(pattern2)
	if err != nil {
		return fmt.Errorf("invalid regex2: %w", err)
	}

	// 2. 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 3. 创建 Scanner
	scanner := bufio.NewScanner(file)

	// 【可选】如果预期有超长行（如压缩日志或单行JSON），扩大缓冲区
	// const maxCapacity = 1024 * 1024 // 1MB
	// buf := make([]byte, 0, 64*1024)
	// scanner.Buffer(buf, maxCapacity)

	lineNum := 0
	questionID := 0
	var e *list.Element
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// 4. 执行匹配
		if re.MatchString(line) {
			//indices := re.FindAllStringIndex(line, -1)
			c := strings.Count(line, "]")
			if strings.HasSuffix(line, "]") && c == 1 {
				//indices := re.FindAllStringIndex(, -1)
				s := re.FindString(line)
				is := re.FindStringIndex(line)
				//indices := re.FindAllStringIndex(s, -1)
				//s := re.FindAllStringSubmatch(line, -1)
				questionID++
				//log.Printf("Line %d: questionID : %d . %s ; Answer: %s class: %s\n", lineNum, questionID, s[strings.Index(s, ".")+1:], extractWithIndex(s, '（', '）'), extractWithIndex(s, '[', ']'))

				e = ls.PushBack(NewQuestion(questionID, s, extractWithIndex(s, '（', '）'), extractWithIndex(s, '[', ']')))

				if len(s) < len(line) {
					if ev, ok := e.Prev().Value.(*Question); ok {
						StringToMap(ev, line[0:is[0]], re2)
					}
				}

			} else {

				strs := strings.Split(line, "]")
				for i := 0; i < c; i++ {
					strs[i] = strs[i] + "]"
					s := re.FindString(strs[i])
					is := re.FindStringIndex(strs[i])
					//s := re.FindStringSubmatch(strs[i])
					questionID++

					//log.Printf("Line %d: questionID : %d . %s ; Answer: %s class: %s\n", lineNum, questionID, s[strings.Index(s, ".")+1:], extractWithIndex(s, '（', '）'), extractWithIndex(s, '[', ']'))
					e = ls.PushBack(NewQuestion(questionID, s, extractWithIndex(s, '（', '）'), extractWithIndex(s, '[', ']')))
					if len(s) < len(strs[i]) {
						if ev, ok := e.Prev().Value.(*Question); ok {
							StringToMap(ev, strs[i][0:is[0]], re2)
						}
					}
				}
			}

			// 若需提取具体子匹配内容：
			// matches := re.FindStringSubmatch(line)
			// if len(matches) > 1 { fmt.Println("Captured:", matches) }
		}
		if re2.MatchString(line) {
			if ev, ok := e.Value.(*Question); ok {
				StringToMap(ev, line, re2)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading file error: %w", err)
	}

	return nil
}

func main() {

	ls := list.New()
	err := SearchFileByLine(ls, "Source.txt", `(\d+)\.(.*)`, `^[A-Z]、`)
	if err != nil {
		fmt.Println("Error:", err)
	}
	for e := ls.Front(); e != nil; e = e.Next() {
		log.Printf("%v ", e.Value.(*Question).ToString())
	}
}
