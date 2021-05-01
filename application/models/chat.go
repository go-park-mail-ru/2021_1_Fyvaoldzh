package models

import "time"

type NewMessage struct {
	To   uint64 `json:"to"`
	Text string `json:"text"`
}

type Message struct {
	ID     uint64 `json:"id"`
	FromMe bool   `json:"fromMe"`
	Text   string `json:"text"`
	Date   string `json:"date"`
	Redact bool   `json:"redact"`
	Read   bool   `json:"read"`
}

type MessageSQL struct {
	ID     uint64
	From   uint64
	To     uint64
	Text   string
	Date   time.Time
	Redact bool
	Read   bool
}

type Dialogue struct {
	ID             uint64      `json:"id"`
	Interlocutor   UserOnEvent `json:"interlocutor"`
	DialogMessages Messages    `json:"messages"`
}

type EasyDialogueMessageSQL struct {
	ID    uint64
	User1 uint64
	User2 uint64
}

type DialogueSQL struct {
	ID             uint64
	User1          uint64
	User2          uint64
	DialogMessages MessagesSQL
}

type DialogueCard struct {
	ID           uint64      `json:"id"`
	Interlocutor UserOnEvent `json:"interlocutor"`
	LastMessage  Message     `json:"message"`
}

type DialogueCardSQL struct {
	ID          uint64
	User_1      uint64
	User_2      uint64
	LastMessage MessageSQL
}

//easyjson:json
type Messages []Message

//easyjson:json
type MessagesSQL []MessageSQL

//easyjson:json
type Dialogues []Dialogue

//easyjson:json
type DialogueCards []DialogueCard

//easyjson:json
type DialogueCardsSQL []DialogueCardSQL

//easyjson:json
type DialoguesSQL []DialogueSQL

func ConvertMessage(old MessageSQL, uid uint64) Message {
	var newMessage Message
	newMessage.ID = old.ID
	newMessage.Text = old.Text
	newMessage.Redact = old.Redact
	newMessage.Read = old.Read
	newMessage.Date = old.Date.String()
	if old.From == uid {
		newMessage.FromMe = true
	} else {
		newMessage.FromMe = false
	}
	return newMessage
}
