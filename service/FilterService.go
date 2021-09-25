package service

import (
	"gochat/model"
	"gochat/utils"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var (
	FilterService filterService
	az, AZ        []rune
)

type filterService struct {
	top      *FilterNode
	starWord rune

	msgQueue   chan string
	msgRecords *utils.List
}

type FilterNode struct {
	Next  map[rune]*FilterNode
	IsEnd bool
}

func init() {
	FilterService.init()
	az = []rune("az")
	AZ = []rune("AZ")
}

func (this *filterService) Filter(msg string) string {
	//这里压入消息队列，是为了异步处理记录消息历史
	this.msgQueue <- msg
	words := []rune(msg)
	for i := 0; i < len(words); i++ {
		i, _ = this.replaceWords(words, this.top, i)
	}
	return string(words)
}

func (this *filterService) PopularWords(n int) string {
	records := this.msgRecords.GetAll()
	t := time.Now().Unix() - int64(n)
	wordsMap := map[string]int{}

	for _, v := range records {
		if v == nil {
			continue
		}
		mr := v.(*model.MsgRecord)
		if mr.Time >= t {
			//分词并加入统计中
			for _, v := range mr.Words {

				rs := []rune(v)
				for k, v := range rs {
					if v >= az[0] && v <= az[1] || v >= AZ[0] && v <= AZ[1] {
						continue
					}
					rs[k] = ' '
				}

				tmp := strings.Split(string(rs), " ")
				for _, v := range tmp {
					if v == "" {
						continue
					}
					if n, ok := wordsMap[v]; ok {
						wordsMap[v] = n + 1
					} else {
						wordsMap[v] = 1
					}
				}
			}
		}
	}

	//查询最高频率的词
	mx := 0
	mxWords := ""
	for w, n := range wordsMap {
		if n > mx {
			mx = n
			mxWords = w
		}
	}

	return mxWords
}

func (this *filterService) replaceWords(words []rune, node *FilterNode, i int) (int, bool) {
	if i >= len(words) {
		if node.IsEnd { //全词匹配
			return i - 1, true
		} else { //没有匹配成功，返回false，回溯看看有没有可以匹配的词
			return i, false
		}
	}
	v, ok := node.Next[words[i]]
	if ok {
		if idx, ok := this.replaceWords(words, v, i+1); ok {
			words[i] = this.starWord
			return idx, true
		} else if node.IsEnd { //贪心匹配失败后，回朔回来时，若存在单词，则匹配之
			words[i-1] = this.starWord
			return i - 1, true
		}
	} else {
		if node.IsEnd {
			//是末位，进行替换，并返回true
			words[i-1] = this.starWord
			return i - 1, true
		}
	}
	return i, false
}

func (this *filterService) init() {
	var (
		words          []rune
		node, lastNode *FilterNode
		ok             bool
	)

	this.msgQueue = make(chan string, 10000)
	this.starWord = []rune("*")[0]
	this.msgRecords = utils.NewList(60) //这里初始化为60个节点，是为了存储60秒的记录

	bs, err := ioutil.ReadFile("config/filter.txt")

	if err != nil {
		bs, err = ioutil.ReadFile("../config/filter.txt")
		if err != nil {
			panic("Load Filter.txt Failed!\n" + err.Error())
		}
	}
	i := 0

	rs := []rune(string(bs))

	this.top = &FilterNode{
		Next:  map[rune]*FilterNode{},
		IsEnd: false,
	}
	//将过滤词典按每字为一个节点转换为层次结构进行存储，以便于快速对比
	for {
		if i >= len(rs) || rs[i] == '\n' {
			words = rs[0:i]
			if i < len(rs) {
				rs = rs[i+1:]
			} else {
				rs = []rune{}
			}
			i = 0
			node = this.top
			for _, v := range words {
				lastNode = node
				node, ok = node.Next[v]
				if !ok {
					node = &FilterNode{
						Next:  map[rune]*FilterNode{},
						IsEnd: false,
					}
					lastNode.Next[v] = node
				}
			}
			node.IsEnd = true //标记为真表示这里存在一个单词结尾
		}
		i++
		if len(rs) == 0 {
			break
		}
	}

	log.Println("fitler inited")

	go this.recordWordsForPopular()
}

func (this *filterService) recordWordsForPopular() {
	for {
		words := <-this.msgQueue
		var msgRecord *model.MsgRecord = nil
		t := time.Now().Unix()
		pos := int(t % 60)
		node := this.msgRecords.GetAt(pos)
		if node != nil {
			msgRecord = node.(*model.MsgRecord)
			if msgRecord.Time != t {
				msgRecord = nil
			}
		}

		if msgRecord == nil {
			msgRecord = &model.MsgRecord{
				Words: []string{},
				Time:  t,
			}
			this.msgRecords.SetAt(pos, msgRecord)
		}
		msgRecord.Words = append(msgRecord.Words, words)
	}
}
