package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	re, re2, re3, re4 *regexp.Regexp
)

type Question struct {
	Id       int               `json:"Id"`
	Head     string            `json:"Head"`
	Answer   string            `json:"Answer"`
	Class    string            `json:"Class"`
	Option   map[string]string `json:"Option"`
	OkTimes  int               `json:"OkTimes"`
	ErrTimee int               `json:"ErrTime"`
	Location bool              `json:"Location"`
}

func NewQuestion(id int, head string, answer string, class string) *Question {
	return &Question{
		Id:     id,
		Head:   head,
		Answer: answer,
		Class:  class,
		Option: make(map[string]string),
	}
}

/*
func (q *Question) ToString() string {
	if q.Class == "判断题" {
		return fmt.Sprintf("%s|%s|对|错|||||||%s|||\n", q.Head, q.Class, q.Answer)
	}

	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|||\n", q.Head, q.Class, q.Option["A"], q.Option["B"], q.Option["C"], q.Option["D"], q.Option["E"], q.Option["F"], q.Option["G"], q.Option["H"], q.Answer)
}
*/

func (q *Question) ToString() string {
	if q.Class == "判断题" {
		return fmt.Sprintf("%d.%s(%s)\n答案: %s\n解析:\n", q.Id, q.Head, q.Answer, q.Answer)
	}
	o := ""
	keys := sortMapK(q.Option)
	for _, k := range *keys {
		o = o + fmt.Sprintf("%s.%s\n", k, q.Option[k])
	}
	return fmt.Sprintf("%d.%s(%s)\n%s 答案: %s\n解析:\n", q.Id, q.Head, q.Answer, o, q.Answer)
}

func getFirstCharRange(s string) string {
	for _, r := range s {
		return string(r)
	}
	return ""
}
func StringToMap(q *Question, it string, re *regexp.Regexp) {
	//strs := strings.Split(it, " ")
	strs := re.FindAllString(it, -1)
	for _, str := range strs {
		if re.MatchString(str) {
			//s := re.FindString(str)
			sl := strings.Split(str, ".")
			q.Option[getFirstCharRange(sl[0])] = sl[1]

		}
	}

}
func extractWithIndex(s string, st rune, et rune) string {
	//r := re.FindAllString(s, -1)
	start := strings.Split(s, string(st))
	for _, it := range start {
		if i := strings.Index(it, string(et)); i == 0 || i == -1 {
			continue
		} else {
			return it[0:i]
		}
	}
	return ""

}

func getAnswerAndClass(s string, st rune, et rune) (string, string) {
	c := strings.Count(s, ")")
	for i := 0; i < c; i++ {
		a := extractWithIndex(s, st, et)
		a = strings.Trim(a, " ")
		if strings.Contains(a, "对") || strings.Contains(a, "错") {
			return a, "判断题"
		} else if re4.MatchString(a) && len(a) == 1 {
			return a, "单选题"
		} else if re4.MatchString(a) && len(a) > 1 {
			return a, "多选题"
		} else if a == "" {
			s = strings.Replace(s, ")", "", 1)
		}
		s = strings.ReplaceAll(s, a, "")
	}
	return "", "未知"
}
func deleteredundant(s string, st string, et string) string {
	if strings.Count(s, et) > 1 {
		a, c := getAnswerAndClass(s, '(', ')')
		if c != "未知" {
			lstrs := strings.Split(s, ")")
			s = ""
			for _, lstr := range lstrs {
				if strings.Contains(lstr, a) {
					lstr = lstr + ")"
				}
				s = s + lstr
			}
			//line = strings.Replace(line, ")", "", strings.Count(line, ")")-1)

		}

	}
	return s
}

// SearchFileByLine 逐行搜索文件，返回匹配的行号和內容
func SearchFileByLine(ls *list.List, filename string) error {

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
				line = deleteredundant(line, "(", ")")
				//indices := re.FindAllStringIndex(, -1)
				s := re.FindString(line)
				is := re.FindStringIndex(line)
				//h := re.FindStringSubmatch(line)
				//indices := re.FindAllStringIndex(s, -1)
				//s := re.FindAllStringSubmatch(line, -1)
				questionID++
				//log.Printf("Line %d: questionID : %d . %s ; Answer: %s class: %s\n", lineNum, questionID, s[strings.Index(s, ".")+1:], extractWithIndex(s, '（', '）'), extractWithIndex(s, '[', ']'))
				a, c := getAnswerAndClass(s, '(', ')')
				e = ls.PushBack(NewQuestion(questionID, s[strings.Index(s, ".")+1:strings.Index(s, "[")], a, c))
				//e = ls.PushBack(NewQuestion(questionID, s[strings.Index(s, ".")+1:strings.Index(s, "[")], extractWithIndex(s, '(', ')'), extractWithIndex(s, '[', ']')))
				//log.Printf("%v ", e.Value.(*Question).ToString())

				if len(s) < len(line) {
					if ev, ok := e.Prev().Value.(*Question); ok {
						StringToMap(ev, line[0:is[0]], re2)
						//StringToMap(ev, line[is[1]:], re2)
					}
				}

			} else {

				strs := strings.Split(line, "]")
				for i := 0; i < c; i++ {
					strs[i] = strs[i] + "]"
					strs[i] = deleteredundant(strs[i], "(", ")")
					s := re.FindString(strs[i])
					is := re.FindStringIndex(strs[i])
					//s := re.FindStringSubmatch(strs[i])
					questionID++
					a, c := getAnswerAndClass(s, '(', ')')
					e = ls.PushBack(NewQuestion(questionID, s[strings.Index(s, ".")+1:strings.Index(s, "[")], a, c))
					//log.Printf("Line %d: questionID : %d . %s ; Answer: %s class: %s\n", lineNum, questionID, s[strings.Index(s, ".")+1:], extractWithIndex(s, '(', ')'), extractWithIndex(s, '[', ']'))
					//e = ls.PushBack(NewQuestion(questionID, s[strings.Index(s, ".")+1:strings.Index(s, "[")], extractWithIndex(s, '(', ')'), extractWithIndex(s, '[', ']')))
					//log.Printf("%v ", e.Value.(*Question).ToString())
					if len(s) < len(strs[i]) {
						if ev, ok := e.Prev().Value.(*Question); ok {
							StringToMap(ev, strs[i][0:is[0]], re2)
						}
					}
				}
				if len(strs) > c {
					StringToMap(e.Value.(*Question), strs[len(strs)-1], re2)
				}
			}

			// 若需提取具体子匹配内容：
			// matches := re.FindStringSubmatch(line)
			// if len(matches) > 1 { fmt.Println("Captured:", matches) }
		}

		if re3.MatchString(line) && !strings.Contains(line, "]") {
			line = deleteredundant(line, "(", ")")
			strs := re3.FindAllString(line, -1)
			for j := 0; j < len(strs); j++ {

				/*
					a := extractWithIndex(s, '(', ')')
					if strings.Index(a, "对") != -1 || strings.Index(a, "错") != -1 {
						questionID++
						ls.PushBack(NewQuestion(questionID, s[strings.Index(s, ".")+1:strings.LastIndex(s, "(")], a, "判断题"))
					} else if re4.MatchString(a) && len(a) == 1 {
						questionID++
						ls.PushBack(NewQuestion(questionID, s[strings.Index(s, ".")+1:strings.LastIndex(s, "(")], a, "单选题"))
					} else if re4.MatchString(a) && len(a) > 1 {
						questionID++
						ls.PushBack(NewQuestion(questionID, s[strings.Index(s, ".")+1:strings.LastIndex(s, "(")], a, "多选题"))
					} else {
						ls.PushBack(NewQuestion(questionID, s[strings.Index(s, ".")+1:strings.LastIndex(s, "(")], a, "未知"))
					}
				*/
				a, c := getAnswerAndClass(strs[j], '(', ')')
				questionID++
				e = ls.PushBack(NewQuestion(questionID, strs[j][strings.Index(strs[j], ".")+1:strings.LastIndex(strs[j], "(")], a, c))
				line = strings.ReplaceAll(line, strs[j], "")
				if j < len(strs)-1 {
					//line[0:strings.Index(line,strs[j+1])]
					StringToMap(e.Value.(*Question), line[0:strings.Index(line, strs[j+1])], re2)
				}
				StringToMap(e.Value.(*Question), line, re2)

			}
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

type KV struct {
	Key   string
	Value int
}

func sortMapV(m map[string]int) (*[]KV, int) {

	s := make([]KV, 0, len(m))
	sum := 0
	for k, v := range m {
		s = append(s, KV{k, v})
		sum = sum + v
	}

	sort.Slice(s, func(i, j int) bool {
		return s[i].Value > s[j].Value
	})
	return &s, sum
}
func sortMapK(m map[string]string) *[]string {

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys) // 字典序升序
	return &keys
}

func main() {

	ls := list.New()
	// 1. 预编译正则
	var err error
	re, err = regexp.Compile(`\d+\..*\]`)
	if err != nil {
		fmt.Print("invalid regex1: %w", err)
	}

	re2, err = regexp.Compile(`[A-Z]\.[^ABCDEF]*`)
	if err != nil {
		fmt.Print("invalid regex2: %w", err)
	}

	re3, err = regexp.Compile(`\d+\.[^\]\)]*\)`)
	if err != nil {
		fmt.Print("invalid regex3: %w", err)
	}

	re4, err = regexp.Compile(`^[A-Z]`)
	if err != nil {
		fmt.Print("invalid regex4: %w", err)
	}

	err = SearchFileByLine(ls, "F:/go/dataclean/main/Source.txt")
	if err != nil {
		fmt.Println("Error:", err)
	}

	sum := ls.Len()
	questionclass := make(map[string]int)
	questionanswer := make(map[string]int)
	oneanswer := make(map[string]int)

	yesorno := make(map[string]int)

	//var q *Question
	for e := ls.Front(); e != nil; e = e.Next() {
		if q, ok := e.Value.(*Question); ok {
			q.Head = strings.ReplaceAll(q.Head, q.Answer, "")
			//log.Printf("%v ", q.ToString())

			//if q.Class == "未知" && q.Answer == "" {
			//	log.Printf("%v ", e.Value.(*Question).ToString())
			//}
			questionanswer[q.Answer] = questionanswer[q.Answer] + 1
			questionclass[q.Class] = questionclass[q.Class] + 1
		}
	}
	oneanswer["A"] = questionanswer["A"]
	oneanswer["B"] = questionanswer["B"]
	oneanswer["C"] = questionanswer["C"]
	oneanswer["D"] = questionanswer["D"]

	yesorno["对"] = questionanswer["对"]
	yesorno["错"] = questionanswer["错"]

	delete(questionanswer, "A")
	delete(questionanswer, "B")
	delete(questionanswer, "C")
	delete(questionanswer, "D")
	delete(questionanswer, "对")
	delete(questionanswer, "错")

	s1, sum1 := sortMapV(oneanswer)
	s2, sum2 := sortMapV(questionanswer)
	s3, sum3 := sortMapV(yesorno)

	fmt.Printf("sum:%d;\n", sum)

	fmt.Printf("题型:%v;\n", questionclass)

	//fmt.Printf("单选答案:%v;\n", oneanswer)

	for _, item := range *s1 {
		fmt.Printf("单选题答案其中:%s: %d 个; 概率 : %.2f\n", item.Key, item.Value, float32(item.Value)/float32(sum1))
	}

	//fmt.Printf("多选答案:%v;\n", questionanswer)

	for _, item := range *s2 {
		fmt.Printf("多选题答案其中:%s: %d 个; 概率 : %.2f\n", item.Key, item.Value, float32(item.Value)/float32(sum2))
	}

	//fmt.Printf("判断答案:%v;\n", yesorno)

	for _, item := range *s3 {
		fmt.Printf("判断题答案其中:%s: %d 个; 概率 : %.2f\n", item.Key, item.Value, float32(item.Value)/float32(sum3))
	}

	file, err := os.Create("Simple2.txt")
	if err != nil {
		fmt.Errorf("err:%v", err)
		return
	}
	defer file.Close() // 确保文件描述符释放

	// 创建缓冲写入器
	writer := bufio.NewWriter(file)
	sum = 0

	for e := ls.Front(); e != nil; e = e.Next() {
		if q, ok := e.Value.(*Question); ok {
			if q.Answer != "ABCD" && q.Class == "多选题" {
				_, err = writer.WriteString(fmt.Sprintf("%s", q.ToString()))
				sum++
			}
			if q.Answer != "C" && q.Class == "单选题" {
				_, err = writer.WriteString(fmt.Sprintf("%s", q.ToString()))
				sum++
			}

			if q.Answer == "错" && q.Class == "判断题" {
				_, err = writer.WriteString(fmt.Sprintf("%s", q.ToString()))
				sum++
			}
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	/*
		for e := ls.Front(); e != nil; e = e.Next() {
			if q, ok := e.Value.(*Question); ok {
				if q.Answer != "ABCD" && q.Class == "多选题" {
					b, _ := json.Marshal(q)
					_, err = writer.WriteString(string(b))
					sum++
				}
				if q.Answer != "C" && q.Class == "单选题" {
					b, _ := json.Marshal(q)
					_, err = writer.WriteString(string(b))
					sum++
				}

				if q.Answer == "错" && q.Class == "判断题" {
					b, _ := json.Marshal(q)
					_, err = writer.WriteString(string(b))
					sum++
				}
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	*/
	// 【关键步骤】必须调用 Flush，否则缓冲区剩余数据会丢失
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("写入精简版题库:%d;条\n", sum)
}
