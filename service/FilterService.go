package service

/*
 过滤处理类

*/
import (
	"gochat/model"
	"gochat/utils"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var (
	FilterService filterService //单例对象
	az, AZ        []rune        //边界
)

//过滤类
type filterService struct {
	top      *FilterNode
	starWord rune

	msgQueue   chan string
	msgRecords *utils.List
}

//词库节点
type FilterNode struct {
	Next  map[rune]*FilterNode
	IsEnd bool
}

//初始化
func init() {
	FilterService.init()
	az = []rune("az")
	AZ = []rune("AZ")
}

//执行过滤，输入待过滤字符串，返回已过滤的字符串
func (this *filterService) Filter(msg string) string {
	//这里压入消息队列，是为了异步处理记录消息历史
	this.msgQueue <- msg
	words := []rune(msg)
	for i := 0; i < len(words); i++ {
		i, _ = this.replaceWords(words, this.top, i)
	}
	return string(words)
}

//用于统计指令秒数内最常用的单词，参数为秒，这里需要注意，这里只能传入不大于60
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
					if v >= az[0] && v <= az[1] || v >= AZ[0] && v <= AZ[1] { //只认为英文字母为单词，其它均不认为是单词
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

//私有方法，进行过滤词替换
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

//对象初始化，由于是单例对象，只会被调用一次
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

//这是一个记录历史消息的消息循环，会随系统启动时而开启
func (this *filterService) recordWordsForPopular() {
	for {
		words := <-this.msgQueue
		var msgRecord *model.MsgRecord = nil
		t := time.Now().Unix()
		pos := int(t % 60) //每一秒我们存放一个历史消息切片，由于是循环队列，故尔只会记录最近60秒内的历史消息，超过的会被覆盖。
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
