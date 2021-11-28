package server

import (
	"context"
	"kudago/application/microservices/chat/chat"
	"kudago/application/microservices/chat/client"
	"kudago/application/microservices/chat/proto"
	"kudago/application/models"
)

type ChatServer struct {
	usecase chat.UseCase
}

func NewChatServer(usecase chat.UseCase) *ChatServer {
	return &ChatServer{usecase: usecase}
}

func (cs *ChatServer) GetAllDialogues(_ context.Context, idPage *proto.IdPage) (*proto.DialogueCards, error) {
	answer, err := cs.usecase.GetAllDialogues(idPage.Id, int(idPage.Page))
	if err != nil {
		return nil, err
	}

	return client.ConvertDialogueCardsToProto(answer), nil
}

func (cs *ChatServer) GetAllNotifications(_ context.Context, idPage *proto.IdPage) (*proto.Notifications, error) {
	answer, err := cs.usecase.GetAllNotifications(idPage.Id, int(idPage.Page))
	if err != nil {
		return nil, err
	}

	return client.ConvertNotificationsToProto(answer), nil
}

func (cs *ChatServer) GetAllCounts(_ context.Context, id *proto.Id) (*proto.Counts, error) {
	answer, err := cs.usecase.GetAllCounts(id.Id)
	if err != nil {
		return nil, err
	}

	return client.ConvertCountsToProto(answer), nil
}

func (cs *ChatServer) GetOneDialogue(_ context.Context, idIdPage *proto.IdIdPage) (*proto.Dialogue, error) {
	answer, err := cs.usecase.GetOneDialogue(idIdPage.Id1, idIdPage.Id2, int(idIdPage.Page))
	if err != nil {
		return nil, err
	}

	return client.ConvertDialogueToProto(answer), nil
}

func (cs *ChatServer) DeleteDialogue(_ context.Context, idId *proto.IdId) (*proto.Answer, error) {
	err := cs.usecase.DeleteDialogue(idId.Id1, idId.Id2)
	if err != nil {
		return nil, err
	}
	/* TODO ATTENTION
	если в ошибках падает не 500, то можно вернуть в Answer.
	поднимаешь флаг, пишешь соо, а на стороне клиента ловишь флаг и выдаешь с 4хх
	*/

	return &proto.Answer{}, nil
}

func (cs *ChatServer) SendMessage(_ context.Context, msg *proto.SendEditMessage) (*proto.Answer, error) {
	newMsg := models.NewMessage{
		To:   msg.Id1,
		Text: msg.Text,
	}
	err := cs.usecase.SendMessage(&newMsg, msg.Id2)
	if err != nil {
		return nil, err
	}

	return &proto.Answer{}, nil
}

func (cs *ChatServer) EditMessage(_ context.Context, msg *proto.SendEditMessage) (*proto.Answer, error) {
	newMsg := models.RedactMessage{
		ID:   msg.Id2,
		Text: msg.Text,
	}
	err := cs.usecase.EditMessage(msg.Id1, &newMsg)
	if err != nil {
		return nil, err
	}

	return &proto.Answer{}, nil
}

func (cs *ChatServer) DeleteMessage(_ context.Context, msg *proto.IdId) (*proto.Answer, error) {
	err := cs.usecase.DeleteMessage(msg.Id1, msg.Id2)
	if err != nil {
		return nil, err
	}

	return &proto.Answer{}, nil
}

func (cs *ChatServer) Mailing(_ context.Context, in *proto.MailingIn) (*proto.Answer, error) {
	err := cs.usecase.Mailing(in.UserId, &models.Mailing{
		EventID: in.EventId,
		To:      client.ConvertIds(in.To),
	})
	if err != nil {
		return nil, err
	}

	return &proto.Answer{}, nil
}

func (cs *ChatServer) Search(c context.Context, in *proto.SearchIn) (*proto.DialogueCards, error) {
	answer, err := cs.usecase.Search(in.Uid, int(in.Id), in.Str, int(in.Page))
	if err != nil {
		return nil, err
	}

	return client.ConvertDialogueCardsToProto(answer), nil
}
