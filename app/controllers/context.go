package controllers

import (
	"log"
	"strconv"
	//"CrazyTestor/app/models"
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
	if (*c).QuesstionID == 0  {
		question, answers := questionService.Get(0)
		//rsp := mashup question and answers
		rsp := &Message{"ToUserName": (*s)["FromUserName"],
			"FromUserName": (*s)["ToUserName"],
			"CreateTime":   (*s)["CreateTime"],
		//	"MsgId":        (*s)["MsgId"],
			"MsgType":      "text"}

		(*rsp)["Content"] = question.Title; // + "\n" + answers[0].Content
		for i,ans := range  answers {
			(*rsp)["Content"] +="\n 答案"+strconv.Itoa(i+1)+":" + ans.Content
		}
		log.Println("start test", rsp)
		(*c).QuesstionID = question.Id
		cmgr.out <- rsp
		

		log.Print(question.Id)
		return
	}
	log.Println("get AnswerByContent")
	//amswer := getAnswer(pid, data)
	answer := questionService.GetAnswer((*c).QuesstionID, data)
	rsp := &Message{"ToUserName": (*s)["FromUserName"],
		"FromUserName": (*s)["ToUserName"],
		"CreateTime":   (*s)["CreateTime"],
	//	"MsgId":        (*s)["MsgId"],
		"MsgType":      "text"}
	log.Println("next id:" ,answer.NextQuestionId)
	if(answer.NextQuestionId == 0) {
		(*rsp)["Content"] = answer.Result
		log.Println("Content:",(*rsp)["Content"] )
		(*c).QuesstionID = 0
		log.Println("next answer", rsp)
		cmgr.out <- rsp
	}  else {
		question, answers := questionService.Get(answer.NextQuestionId)
		(*rsp)["Content"] = question.Title;// + "\n" + answers[0].Content
		for i,ans := range  answers {
			(*rsp)["Content"] +="\n 答案"+strconv.Itoa(i+1)+":" + ans.Content
		}
		(*c).QuesstionID = question.Id
		log.Println("next answer", rsp)
		cmgr.out <- rsp
	}




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
						log.Println("Message Exit")
						go k.handle(s)
						goto NEXT
					}
				}
				log.Println("Message not Exit", i)
				for i, k = range cmgr.ctxs {
					if k == nil {
						cmgr.ctxs[i] = NewContext(s)
						go cmgr.ctxs[i].handle(s)
					} else {
						panic("Out of slice len")
					}
				}
			}
		NEXT:
		}
	}()
}
