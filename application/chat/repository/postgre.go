package repository

import (
	"context"
	"errors"
	"fmt"
	"kudago/application/chat"
	"kudago/application/models"
	"kudago/pkg/constants"
	"kudago/pkg/logger"
	"net/http"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

type ChatDatabase struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func NewChatDatabase(conn *pgxpool.Pool, logger logger.Logger) chat.Repository {
	return &ChatDatabase{pool: conn, logger: logger}
}

func (cd ChatDatabase) GetAllDialogues(uid uint64, page int) (models.DialogueCardsSQL, error) {
	var dialogues models.DialogueCardsSQL
	err := pgxscan.Select(context.Background(), cd.pool, &dialogues,
		`SELECT DISTINCT ON(d.id) d.id as ID, d.user_1, d.user_2, m.id as IDMes,
		m.mes_from, m.mes_to, m.text, m.date, m.redact, m.read
	FROM dialogues d JOIN messages m on d.id = m.id_dialogue
	WHERE user_1 = $1 OR user_2 = $1
	ORDER BY id, date DESC
	LIMIT $2 OFFSET $3`, uid, constants.ChatPerPage, (page-1)*constants.ChatPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(dialogues) == 0 {
		cd.logger.Debug("no rows in method GetAllEvents")
		return models.DialogueCardsSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return nil, err
	}

	return dialogues, nil
}

//Исправить значения на read
func (cd ChatDatabase) GetMessages(id uint64, page int) (models.MessagesSQL, error) {
	var messages models.MessagesSQL
	err := pgxscan.Select(context.Background(), cd.pool, &messages,
		`SELECT id, mes_from, mes_to, text,
		date, redact, read FROM messages
	WHERE id_dialogue = $1
	LIMIT $2 OFFSET $3`, id, constants.ChatPerPage, (page-1)*constants.ChatPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(messages) == 0 {
		cd.logger.Debug("no rows in method GetMessages")
		return models.MessagesSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return nil, err
	}

	return messages, nil
}

//Здесь as оставляем, с разных таблиц в одну структурку пишем
func (cd ChatDatabase) GetEasyDialogue(id uint64) (models.EasyDialogueMessageSQL, error) {
	var dialogue []models.EasyDialogueMessageSQL
	err := pgxscan.Select(context.Background(), cd.pool, &dialogue,
		`SELECT id as ID, user_1 as User1, user_2 as User2 FROM dialogues
	WHERE id = $1`, id)

	if errors.As(err, &pgx.ErrNoRows) {
		cd.logger.Debug("no rows in method GetEventsByCategory")
		return models.EasyDialogueMessageSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return models.EasyDialogueMessageSQL{}, nil
	}
	return dialogue[0], nil
}

func (cd ChatDatabase) GetEasyMessage(id uint64) (models.EasyDialogueMessageSQL, error) {
	var message []models.EasyDialogueMessageSQL
	err := pgxscan.Select(context.Background(), cd.pool, &message,
		`SELECT id as ID, mes_from as User1, mes_to as User2 FROM messages
	WHERE id = $1`, id)
	if errors.As(err, &pgx.ErrNoRows) {
		cd.logger.Debug("no rows in method GetEventsByCategory")
		return models.EasyDialogueMessageSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return models.EasyDialogueMessageSQL{}, nil
	}
	return message[0], nil
}

func (cd ChatDatabase) DeleteDialogue(id uint64) error {
	resp, err := cd.pool.Exec(context.Background(),
		`DELETE FROM dialogues WHERE id = $1`, id)

	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	if resp.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("Dialogue with id "+fmt.Sprint(id)+" not found"))
	}

	return nil
}

func (cd ChatDatabase) DeleteMessage(id uint64) error {
	resp, err := cd.pool.Exec(context.Background(),
		`DELETE FROM messages WHERE id = $1`, id)

	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	if resp.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("Message with id "+fmt.Sprint(id)+" not found"))
	}

	return nil
}

func (cd ChatDatabase) SendMessage(id uint64, newMessage *models.NewMessage, uid uint64, now time.Time) error {
	// messages (id, id_dialogue, mes_from, mes_to, text, date, redact, read)
	_, err := cd.pool.Exec(context.Background(),
		`INSERT INTO messages 
		VALUES (default, $1, $2, $3, $4, $5, default, default)`,
		id, uid, newMessage.To, newMessage.Text, now)
	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return nil
}

//в вк показывается время отправки и (ред.), если на него навести, то будет время редактирования
func (cd ChatDatabase) EditMessage(id uint64, text string) error {
	_, err := cd.pool.Exec(context.Background(),
		`UPDATE messages SET text = $1, redact = true WHERE id = $2`, text, id)

	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return nil
}

func (cd ChatDatabase) MessagesSearch(uid uint64, str string, page int) (models.MessagesSQL, error) {
	var messages models.MessagesSQL
	err := pgxscan.Select(context.Background(), cd.pool, &messages,
		`SELECT id, mes_from, mes_to, text,
		date, redact, read FROM messages
		WHERE (LOWER(text) LIKE '%' || $1 || '%')
		ORDER BY date DESC
		LIMIT $2 OFFSET $3`, str, constants.ChatPerPage, (page-1)*constants.ChatPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(messages) == 0 {
		cd.logger.Debug("no rows in method CategorySearch with searchstring " + str)
		return models.MessagesSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return models.MessagesSQL{}, err
	}

	return messages, nil
}

func (cd ChatDatabase) DialogueMessagesSearch(uid uint64, id uint64, str string, page int) (models.MessagesSQL, error) {
	var messages models.MessagesSQL
	err := pgxscan.Select(context.Background(), cd.pool, &messages,
		`SELECT id, mes_from, mes_to, text,
		date, redact, read FROM messages
		WHERE (LOWER(text) LIKE '%' || $1 || '%')
		AND id_dialogue = $2
		ORDER BY date DESC
		LIMIT $3 OFFSET $4`, str, id, constants.ChatPerPage, (page-1)*constants.ChatPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(messages) == 0 {
		cd.logger.Debug("no rows in method CategorySearch with searchstring " + str)
		return models.MessagesSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return models.MessagesSQL{}, err
	}

	return messages, nil
}

func (cd ChatDatabase) CheckDialogue(uid1 uint64, uid2 uint64) (bool, uint64, error) {
	var id []uint64
	err := pgxscan.Select(context.Background(), cd.pool, &id,
		`SELECT id FROM dialogues WHERE 
		(user_1 = $1 AND user_2 = $2) OR (user_1 = $2 AND user_2 = $1)`, uid1, uid2)

	if errors.As(err, &pgx.ErrNoRows) || len(id) == 0 {
		return false, 0, nil
	}
	if err != nil {
		cd.logger.Warn(err)
		return false, 0, err
	}
	return true, id[0], nil
}

func (cd ChatDatabase) CheckMessage(uid1 uint64, uid2 uint64) (bool, uint64, error) {
	var id []uint64
	err := pgxscan.Select(context.Background(), cd.pool, &id,
		`SELECT id FROM messages WHERE 
		(mes_from = $1 AND mes_to = $2) OR (mes_from = $2 AND mes_to = $1)`, uid1, uid2)

	if errors.As(err, &pgx.ErrNoRows) || len(id) == 0 {
		return false, 0, nil
	}
	if err != nil {
		cd.logger.Warn(err)
		return false, 0, err
	}
	return true, id[0], nil
}

func (cd ChatDatabase) NewDialogue(uid1 uint64, uid2 uint64) (uint64, error) {
	var id uint64
	err := cd.pool.QueryRow(context.Background(),
		`INSERT INTO dialogues 
		VALUES (default, $1, $2) RETURNING id`,
		uid1, uid2).Scan(&id)
	if err != nil {
		cd.logger.Warn(err)
		return 0, err
	}

	return id, nil
}
