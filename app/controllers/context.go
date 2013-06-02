package controllers

import (
	"CrazyTestor/app/models"
	"log"
	"strconv"
)

var (
	cmgr        *ContextMgr = nil
	DEFAULT_LEN             = 100
)

type Message map[string]string

type Context struct {
	ToUserName   string
	FromUserName string
	CreateTime   string
	MsgType      string
	Content      string
	MsgId        int64
	QuesstionID  int64
	testId       int64
}

func NewContext(s *Message) *Context {
	//var err error
	c := Context{ToUserName: (*s)["ToUserName"],
		FromUserName: (*s)["FromUserName"],
		CreateTime:   (*s)["CreateTime"],
		MsgType:      (*s)["MsgType"],
		Content:      (*s)["Content"]}
	/*
		c.MsgId, err = strconv.ParseInt((*s)["MsgId"], 10, 64)
		if err != nil {
			//panic("strconv")
			log.Println("MsgId error")
		}
	*/
	return &c
}

func (c *Context) sameSession(s *Message) bool {
	if c.FromUserName == (*s)["FromUserName"] {
		return true
	}
	return false
}

func (c *Context) handle(s *Message) {
	data := (*s)["Content"]
	log.Println("handling:", data)
	log.Println("handling Id:", (*c).QuesstionID)
	if (*c).QuesstionID == 0 {
		//rsp := mashup question and answers
		rsp := &Message{"ToUserName": (*s)["FromUserName"],
			"FromUserName": (*s)["ToUserName"],
			"CreateTime":   (*s)["CreateTime"],
			//	"MsgId":        (*s)["MsgId"],
			"MsgType": "text"}
		tid, err := strconv.ParseInt(data, 10, 32)
		var t *models.Test
		ts := testService.Find()
		if err == nil {
			log.Println("选择测试序列:", tid)
			t = testService.Get(tid)
			for i, tt := range ts {
				if i == int(tid)-1 {
					t = &tt
					break
				}
			}
		}

		if t == nil {

			(*rsp)["Content"] = "请选择测试列表:"
			for i, t := range ts {
				(*rsp)["Content"] += "\n" + strconv.Itoa(i+1) + " :" + t.Title
			}
		} else {
			list := questionService.Find(t.Id)
			var content = ""
			if len(list) > 0 {
				content = "您选择了 测试\"" + t.Title + "\",问题如下\n"

				question, answers := questionService.Get(list[0].Id)
				content += question.Title
				for i, ans := range answers {
					content += "\n 答案" + strconv.Itoa(i+1) + ":" + ans.Content
				}
				(*c).QuesstionID = question.Id
				(*c).testId = question.TestId
			} else {
				content = "找不到相关的测试问题."
			}

			(*rsp)["Content"] = content

			log.Println("start test", rsp)

		}
		cmgr.out <- rsp

		return
	}
	log.Println("get AnswerByContent")
	aid, err := strconv.ParseInt(data, 10, 32)
	answer := &models.Answer{}
	log.Println("err:", err)
	if err == nil { //序号回答问题
		log.Println("序号:", aid)
		answer = questionService.GetAnswerById((*c).QuesstionID, aid)
	} else {
		log.Println("文本查找:", data)
		answer = questionService.GetAnswer((*c).QuesstionID, data)
	}
	qid := strconv.Itoa(int((*c).QuesstionID))
	log.Println("question id ", qid)
	//amswer := getAnswer(pid, data)

	rsp := &Message{"ToUserName": (*s)["FromUserName"],
		"FromUserName": (*s)["ToUserName"],
		"CreateTime":   (*s)["CreateTime"],
		//	"MsgId":        (*s)["MsgId"],
		"MsgType": "text"}

	if answer == nil {

		(*rsp)["Content"] = "系统错误,找不到响应的问题:" + qid //+ (*c).QuesstionID;
		cmgr.out <- rsp
		return
	}

	log.Println("next id:", answer.NextQuestionId)
	if answer.NextQuestionId == 0 {
		(*rsp)["Content"] = "测试完成\n测试结果:" + answer.Result
		log.Println("Content:", (*rsp)["Content"])
		(*c).QuesstionID = 0
		log.Println("next answer", rsp)
		cmgr.out <- rsp
	} else {
		question, answers := questionService.Get(answer.NextQuestionId)
		(*rsp)["Content"] = question.Title // + "\n" + answers[0].Content
		for i, ans := range answers {
			(*rsp)["Content"] += "\n 答案" + strconv.Itoa(i+1) + ":" + ans.Content
		}
		(*c).QuesstionID = question.Id

		log.Println("next answer", rsp)
		cmgr.out <- rsp
	}

	// input to statistics table
	log.Println("insert statistics...")
	value := &models.Statistics{}
	value.SuitId = (*c).testId
	value.QuestionId = (*c).QuesstionID
	value.UserId = (*c).FromUserName
	value.Content = answer.Content
	statisticsService.Add(value)

	return
}

type ContextMgr struct {
	ctxs []*Context
	in   chan *Message
	out  chan *Message
}

func Invoke(s *map[string]string) *map[string]string {
	(*cmgr).in <- (*Message)(s)
	omsg := <-(*cmgr).out
	return (*map[string]string)(omsg)
}

func RunContextMgr() {
	if cmgr == nil {
		cmgr = new(ContextMgr)
		cmgr.ctxs = make([]*Context, DEFAULT_LEN, DEFAULT_LEN)
		cmgr.in = make(chan *Message)
		cmgr.out = make(chan *Message)
		log.Println("ContextMgr created")
	}

	go func() {
		log.Println("ContextMgr Loop")
		for { //TODO for quite
			select {
			case s := <-cmgr.in:
				var i int
				var k *Context
				for i, k = range cmgr.ctxs {
					if k != nil && k.sameSession(s) {
						log.Println("Message Exists")
						go k.handle(s)
						goto NEXT
					}
				}
				log.Println("Message not Exists", i)
				for i, k = range cmgr.ctxs {
					if k == nil {
						cmgr.ctxs[i] = NewContext(s)
						go cmgr.ctxs[i].handle(s)
						goto NEXT
					} else {
						continue
						//panic("Out of slice len")
					}
				}
			}
		NEXT:
		}
	}()
}
