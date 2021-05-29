package models

import "time"

type Counts struct {
	Notifications uint64 `json:"notifications"`
	Chat          uint64 `json:"chat"`
}

type NotificationSQL struct {
	ID   uint64
	Type string
	Date time.Time
	Read bool
}

type Notification struct {
	ID        uint64 `json:"id"`
	IDToImage uint64 `json:"id_to_image"`
	Type      string `json:"type"`
	Date      string `json:"date"`
	Text      string `json:"text"`
	Read      bool   `json:"read"`
}

type NewMessage struct {
	To   uint64
	Text string
}

type NewMessageJSON struct {
	To   string `json:"to"`
	Text string `json:"text"`
}

type RedactMessage struct {
	ID   uint64
	Text string
}

type RedactMessageJSON struct {
	ID   string `json:"id"`
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
	From   uint64 `db:"mes_from"`
	To     uint64 `db:"mes_to"`
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
	ID     uint64 `db:"id"`
	User1  uint64 `db:"user_1"`
	User2  uint64 `db:"user_2"`
	IDMes  uint64 `db:"idmes"`
	From   uint64 `db:"mes_from"`
	To     uint64 `db:"mes_to"`
	Text   string
	Date   time.Time
	Redact bool
	Read   bool
}

type Mailing struct {
	EventID uint64   `json:"event"`
	To      []uint64 `json:"to"`
}

type MailingJSON struct {
	EventID uint64   `json:"event"`
	To      []string `json:"to"`
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

//easyjson:json
type Notifications []Notification

//easyjson:json
type NotificationsSQL []NotificationSQL

func ConvertNotification(old NotificationSQL, imageId uint64) Notification {
	var newNotification Notification
	newNotification.ID = old.ID
	newNotification.IDToImage = imageId
	newNotification.Type = old.Type
	newNotification.Date = old.Date.String()
	newNotification.Read = old.Read
	return newNotification
}

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

func ConvertMessageFromCard(old DialogueCardSQL, uid uint64) Message {
	var newMessage Message
	newMessage.ID = old.IDMes
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

func ConvertDialogueCard(old DialogueCardSQL, uid uint64, interlocutor UserOnEvent) DialogueCard {
	var newDialogueCard DialogueCard
	newDialogueCard.ID = old.ID
	newDialogueCard.LastMessage = ConvertMessageFromCard(old, uid)
	newDialogueCard.Interlocutor = interlocutor
	return newDialogueCard
}

func ConvertDialogue(dialogue EasyDialogueMessageSQL, messages MessagesSQL, uid uint64, interlocutor UserOnEvent) Dialogue {
	var newDialogue Dialogue
	newDialogue.ID = dialogue.ID
	for i := range messages {
		newDialogue.DialogMessages = append(newDialogue.DialogMessages, ConvertMessage(messages[i], uid))
	}
	newDialogue.Interlocutor = interlocutor

	return newDialogue
}
