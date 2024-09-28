package store

import (
	"bufio"
	"fmt"
	"os"
)

type Message struct {
	Url      string `json:"url"`
	Req_body []byte `json:"req_body,omitempty"`
}

var MessageRepository messageRepository

var sendMsgQueryTempl string

func init() {

	fd, err := os.Open("./store/insert_message.sql")
	if err != nil {
		panic("cant open file ./store/insert_message.sql during store init() function")
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)

	for scanner.Scan() {
		sendMsgQueryTempl += scanner.Text() + "\n"
	}
}

type messageRepository struct {}

func (m *messageRepository) SaveMessage(msg *Message) error {

	//to_delete
	fmt.Println("save message:", *msg)

	_, err := db.Exec(sendMsgQueryTempl, msg.Url, msg.Req_body)

	if err != nil {
		return fmt.Errorf("message saving in db error: %w", err)
	}

	return nil
}
