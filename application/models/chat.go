package models

import "time"

type NewMessage struct {
	From uint64 `json:"from"`
	To   uint64 `json:"to"`
	Text string `json:"text"`
}

type Message struct {
	ID uint64 `json:"id"`
	//Sender bool   `json:"sender"` Вместо from/to
	From   uint64 `json:"from"`
	To     uint64 `json:"to"`
	Text   string `json:"text"`
	Date   string `json:"date"`
	Redact bool   `json:"redact"`
}

type MessageSQL struct {
	ID uint64
	//Sender bool
	From   uint64 `json:"from"`
	To     uint64 `json:"to"`
	Text   string
	Date   time.Time
	Redact bool
}

type Dialogue struct {
	ID             uint64      `json:"id"`
	Interlocutor   UserOnEvent `json:"interlocutor"`
	DialogMessages Messages    `json:"messages"`
}

type DialogueCard struct {
	ID           uint64      `json:"id"`
	Interlocutor UserOnEvent `json:"interlocutor"`
	LastMessage  Message     `json:"message"`
}

//easyjson:json
type Messages []Message

//easyjson:json
type Dialogues []Dialogue

//easyjson:json
type DialogueCards []DialogueCard

func ConvertMessage(old MessageSQL) Message {
	var newMessage Message
	newMessage.ID = old.ID
	newMessage.From = old.From
	newMessage.To = old.To
	newMessage.Text = old.Text
	newMessage.Redact = old.Redact
	newMessage.Date = old.Date.String()
	return newMessage
}
