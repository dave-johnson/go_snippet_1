package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

const Pub = "P"
const Sub = "S"
const ReqRsp = "RR"

var NextTopic = 1

type Address interface {
}

type CatalogMessage struct {
	code string
	data map[string]string
}

type TopicMessage struct {
	token string
	topic string
}

var wg sync.WaitGroup

func main() {
	str := `{"name": "dave", "age": 53, "hair": "yes"}`
	var data map[string]interface{}
	json.Unmarshal([]byte(str), &data)
	if v, ok := data["hair"]; ok {
		fmt.Println("field value", v)
	} else {
		fmt.Println("field not found")
	}
	fmt.Println(data["name"], "is", data["age"])

	// These are Catalog Service messages sent from HSLib
	catMsg := make(chan CatalogMessage)
	go procCatMsg(catMsg)

	// These are Topic Manager messages sent from HSLib
	topicMsg := make(chan TopicMessage)
	go procTopicMsg(topicMsg)

	fmt.Println("\n")
	wg.Add(1)
	catMsg <- newPubMsg("JobComplete", "tk123")
	wg.Add(1)
	catMsg <- newReqMsg("SCH-Huey", "tk432")
	wg.Add(1)
	catMsg <- newSubMsg("JobComplete", "procJobs")
	wg.Add(1)
	topicMsg <- newTopicMsg("tk432")

	wg.Wait()
}

func newPubMsg(name string, token string) CatalogMessage {
	cm := new(CatalogMessage)
	cm.code = Pub
	cm.data = make(map[string]string)
	cm.data["name"] = name
	cm.data["type"] = cm.code
	cm.data["topic"] = ""
	cm.data["token"] = token
	return *cm
}

func newReqMsg(service string, token string) CatalogMessage {
	cm := new(CatalogMessage)
	cm.code = ReqRsp
	cm.data = make(map[string]string)
	cm.data["service"] = service
	cm.data["type"] = cm.code
	cm.data["topic"] = ""
	cm.data["token"] = token
	return *cm
}

func newSubMsg(name string, cb string) CatalogMessage {
	cm := new(CatalogMessage)
	cm.code = Sub
	cm.data = make(map[string]string)
	cm.data["name"] = name
	cm.data["function"] = cb
	return *cm
}

func procCatMsg(cm chan CatalogMessage) {
	for {
		m := <-cm
		switch m.code {
		case Pub:
			Publish(m)
		case Sub:
			Subscribe(m)
		case ReqRsp:
			RequestRepsonse(m)
		default:
			fmt.Println("unknown type")
		}
		wg.Done()
	}
}

func Subscribe(m CatalogMessage) {
	fmt.Println("Process ", m.data["name"], " with function: ", m.data["function"])
}

func Publish(m CatalogMessage) {
	writeDB("catalog", m.data)
}

func RequestRepsonse(m CatalogMessage) {
	writeDB("catalog", m.data)
}

func newTopicMsg(token string) TopicMessage {
	tm := new(TopicMessage)
	tm.token = token
	return *tm
}

func procTopicMsg(tm chan TopicMessage) {
	for {
		m := <-tm
		// m.topic = tokenMagic(m.token)
		// fmt.Println("Update Publication ", m.token, " with token: ", m.topic)
		d := make(map[string]string)
		d["topic"] = tokenMagic(m.token)
		d["token"] = m.token
		writeDB("catalog", d)
		wg.Done()
	}
}

func tokenMagic(topic string) string {
	NextTopic += 1
	return fmt.Sprintf("Topic-%d", NextTopic)
}

func writeDB(db string, data map[string]string) {
	fmt.Printf("%s: key[%s]\t", db, data["token"])
	for key, value := range data {
		if key != "token" {
			fmt.Print("\t", key, ": ", value)
		}
	}
	fmt.Println("")
}
