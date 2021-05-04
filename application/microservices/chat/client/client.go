package client

/*
protoc --go_out=plugins=grpc:. *.proto
*/

import (
	"context"
	"google.golang.org/grpc"
	chat_proto "kudago/application/microservices/chat/proto"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
)

type ChatClient struct {
	client chat_proto.ChatClient
	gConn  *grpc.ClientConn
	logger logger.Logger
}

func NewChatClient(port string, logger logger.Logger) (IChatClient, error) {
	gConn, err := grpc.Dial(
		constants.Localhost+port,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &ChatClient{client: chat_proto.NewChatClient(gConn), gConn: gConn, logger: logger}, nil
}

func (c *ChatClient) GetAllDialogues(uid uint64, page int) (models.DialogueCards, error) {
	idPage := &chat_proto.IdPage{
		Id:   uid,
		Page: int32(page),
	}

	cards, err := c.client.GetAllDialogues(context.Background(), idPage)
	if err != nil {
		return models.DialogueCards{}, err
	}

	return ConvertDialogueCards(cards), nil
}

func (c *ChatClient) GetOneDialogue(uid1 uint64, uid2 uint64, page int) (models.Dialogue, error) {
	idIdPage := &chat_proto.IdIdPage{
		Id1:  uid1,
		Id2:  uid2,
		Page: int32(page),
	}

	dialogue, err := c.client.GetOneDialogue(context.Background(), idIdPage)
	if err != nil {
		return models.Dialogue{}, err
	}

	return ConvertDialogue(dialogue), nil
}

func (c *ChatClient) DeleteDialogue(uid uint64, id uint64) error {
	idId := &chat_proto.IdId{
		Id1: uid,
		Id2: id,
	}

	_, err := c.client.DeleteDialogue(context.Background(), idId)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatClient) SendMessage(newMessage *models.NewMessage, uid uint64) error {
	msg := &chat_proto.SendEditMessage{
		Id1:  newMessage.To,
		Text: newMessage.Text,
		Id2:  uid,
	}

	_, err := c.client.SendMessage(context.Background(), msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatClient) EditMessage(uid uint64, newMessage *models.RedactMessage) error {
	msg := &chat_proto.SendEditMessage{
		Id1:  uid,
		Text: newMessage.Text,
		Id2:  newMessage.ID,
	}

	_, err := c.client.EditMessage(context.Background(), msg)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatClient) DeleteMessage(uid uint64, id uint64) error {
	idId := &chat_proto.IdId{
		Id1: uid,
		Id2: id,
	}

	_, err := c.client.DeleteMessage(context.Background(), idId)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatClient) Mailing(uid uint64, mailing *models.Mailing) error {
	in := &chat_proto.MailingIn{
		UserId:  uid,
		EventId: mailing.EventID,
		To:      ConvertIdsToProto(mailing.To),
	}

	_, err := c.client.Mailing(context.Background(), in)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatClient) Search(uid uint64, id int, str string, page int) (models.Messages, error) {
	in := &chat_proto.SearchIn{
		Uid:  uid,
		Id:   int32(id),
		Str:  str,
		Page: int32(page),
	}

	answer, err := c.client.Search(context.Background(), in)
	if err != nil {
		return models.Messages{}, err
	}

	return ConvertMessages(answer), nil
}

func (c *ChatClient) Close() {
	if err := c.gConn.Close(); err != nil {
		c.logger.Warn(err)
	}
}
