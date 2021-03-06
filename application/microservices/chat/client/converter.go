package client

import (
	"kudago/application/microservices/chat/proto"
	"kudago/application/models"
)

func ConvertDialogueCards (cards *proto.DialogueCards) models.DialogueCards {
	var newCards models.DialogueCards
	for _, elem := range cards.List {
		var oneCard models.DialogueCard
		oneCard.ID = elem.ID
		oneCard.Interlocutor = ConvertUserOnEvent(elem.Interlocutor)
		oneCard.LastMessage = ConvertMessage(elem.LastMessage)
		newCards = append(newCards, oneCard)
	}
	return newCards
}

func ConvertUserOnEvent(usr *proto.UserOnEvent) models.UserOnEvent {
	var newUser models.UserOnEvent
	newUser.Id = usr.Id
	newUser.Name = usr.Name
	newUser.Avatar = usr.Avatar
	return newUser
}

func ConvertMessage(msg *proto.Message) models.Message {
	var newMsg models.Message
	newMsg.ID = msg.ID
	newMsg.Date = msg.Date
	newMsg.Text = msg.Text
	newMsg.FromMe = msg.FromMe
	newMsg.Read = msg.Read
	newMsg.Redact = msg.Redact
	return newMsg
}

func ConvertMessages(msg *proto.Messages) models.Messages {
	var newMsg models.Messages
	for _, elem := range msg.List {
		newMsg = append(newMsg, ConvertMessage(elem))
	}
	return newMsg
}

func ConvertMessagesToProto(msg models.Messages) *proto.Messages {
	var newMsg proto.Messages
	for _, elem := range msg {
		newMsg.List = append(newMsg.List, ConvertMessageToProto(elem))
	}
	return &newMsg
}

func ConvertDialogue(d *proto.Dialogue) models.Dialogue {
	var newDialogue models.Dialogue
	newDialogue.ID = d.ID
	newDialogue.Interlocutor = ConvertUserOnEvent(d.Interlocutor)
	newDialogue.DialogMessages = ConvertMessages(d.DialogMessages)
	return newDialogue
}

func ConvertIdsToProto(massiv []uint64) *proto.Ids {
	var n proto.Ids
	for _, elem := range massiv {
		n.List = append(n.List, elem)
	}
	return &n
}

func ConvertDialogueCardsToProto(cards models.DialogueCards) *proto.DialogueCards {
	var newCards proto.DialogueCards
	for _, elem := range cards {
		var oneCard proto.DialogueCard
		oneCard.ID = elem.ID
		oneCard.Interlocutor = ConvertUserOnEventToProto(elem.Interlocutor)
		oneCard.LastMessage = ConvertMessageToProto(elem.LastMessage)
		newCards.List = append(newCards.List, &oneCard)
	}
	return &newCards
}

func ConvertUserOnEventToProto (usr models.UserOnEvent) *proto.UserOnEvent{
	var newUser proto.UserOnEvent
	newUser.Id = usr.Id
	newUser.Name = usr.Name
	newUser.Avatar = usr.Avatar
	return &newUser
}

func ConvertMessageToProto(msg models.Message) *proto.Message{
	var newMsg proto.Message
	newMsg.ID = msg.ID
	newMsg.Date = msg.Date
	newMsg.Text = msg.Text
	newMsg.FromMe = msg.FromMe
	newMsg.Read = msg.Read
	newMsg.Redact = msg.Redact
	return &newMsg
}

func ConvertDialogueToProto(d models.Dialogue) *proto.Dialogue {
	var newDialogue proto.Dialogue
	newDialogue.ID = d.ID
	newDialogue.Interlocutor = ConvertUserOnEventToProto(d.Interlocutor)
	newDialogue.DialogMessages = ConvertMessagesToProto(d.DialogMessages)
	return &newDialogue
}

func ConvertIds(massiv *proto.Ids) []uint64 {
	var n []uint64
	for _, elem := range massiv.List {
		n = append(n, elem)
	}
	return n
}