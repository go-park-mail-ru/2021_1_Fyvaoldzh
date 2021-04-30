package repository

import (
	"context"
	"errors"
	"fmt"
	"kudago/application/chat"
	"kudago/application/models"
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

func NewEventDatabase(conn *pgxpool.Pool, logger logger.Logger) chat.Repository {
	return &ChatDatabase{pool: conn, logger: logger}
}

func (cd ChatDatabase) GetAllDialogues(uid uint64, page int) (models.DialogueCardsSQL, error) {
	var dialogues models.DialogueCardsSQL
	err := pgxscan.Select(context.Background(), cd.pool, &dialogues,
		`SELECT DISTINCT ON(d.id) d.id, d.user_1, d.user_2, m.id, m.mes_from, m.mes_to, m.text, m.date, m.redact, m.read
	FROM dialogues d JOIN messages m on d.id = m.dialogue_id
	WHERE user_1 = $1 OR user_2 = $1
	ORDER BY date DESC
	LIMIT 6 OFFSET $2`, uid, (page-1)*6)

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

func (cd ChatDatabase) GetMessages(id uint64) (models.MessagesSQL, error) {
	var messages models.MessagesSQL
	err := pgxscan.Select(context.Background(), cd.pool, &messages,
		`SELECT id, mes_from, mes_to, text, date, redact, read FROM messages
	WHERE dialogue_id = $1`, id)
	if errors.As(err, &pgx.ErrNoRows) || len(messages) == 0 {
		cd.logger.Debug("no rows in method GetAllEvents")
		return models.MessagesSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return nil, err
	}

	return messages, nil
}

//TODO убрать все это, разделить запросы на изи диалог, сообщение и уровень выше будет все собирать
//Исправить значения на read!!!!!!!!!!!!!
func (cd ChatDatabase) GetOneDialogue(id uint64, page int) (models.DialogueSQL, error) {
	//мб здесь будет ошибка, надо будет создавать массив, проверять первый
	var dialogue models.DialogueSQL
	err := pgxscan.Select(context.Background(), cd.pool, &dialogue,
		`SELECT id, user_1, user_2 FROM dialogues
	WHERE id = $1`, id)
	if errors.As(err, &pgx.ErrNoRows) {
		cd.logger.Debug("no rows in method GetEventsByCategory")
		return models.DialogueSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return models.DialogueSQL{}, nil
	}
	//надо бы это вынести на уровень выше
	dialogue.DialogMessages, err = cd.GetMessages(id)
	if err != nil {
		cd.logger.Warn(err)
		return models.DialogueSQL{}, nil
	}
	return dialogue, nil
}

func (cd ChatDatabase) GetEasyDialogue(id uint64) (models.EasyDialogueMessageSQL, error) {
	var dialogue models.EasyDialogueMessageSQL
	err := pgxscan.Select(context.Background(), cd.pool, &dialogue,
		`SELECT id, user_1, user_2 FROM dialogues
	WHERE id = $1`, id)
	if errors.As(err, &pgx.ErrNoRows) {
		cd.logger.Debug("no rows in method GetEventsByCategory")
		return models.EasyDialogueMessageSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return models.EasyDialogueMessageSQL{}, nil
	}
	return dialogue, nil
}

func (cd ChatDatabase) GetEasyMessage(id uint64) (models.EasyDialogueMessageSQL, error) {
	var message models.EasyDialogueMessageSQL
	err := pgxscan.Select(context.Background(), cd.pool, &message,
		`SELECT id, mes_from, mes_to FROM messages
	WHERE id = $1`, id)
	if errors.As(err, &pgx.ErrNoRows) {
		cd.logger.Debug("no rows in method GetEventsByCategory")
		return models.EasyDialogueMessageSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return models.EasyDialogueMessageSQL{}, nil
	}
	return message, nil
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

//На создание нового диалога!
//Подумать насчет проверки валидности переданных значений(чтоб все значения были номральные), ПЛЮС Проверить существует ли диалог с таким id
func (cd ChatDatabase) SendMessage(newMessage *models.NewMessage, uid uint64, now time.Time) error {
	// messages (id, id_dialogue, mes_from, mes_to, text, date, redact, read)
	_, err := cd.pool.Exec(context.Background(),
		`INSERT INTO messages 
		VALUES (default, $1, $2, $3, $4, $5, default, default)`,
		newMessage.DialogueID, uid, newMessage.To, newMessage.Text, now)
	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return nil
}

func (cd ChatDatabase) EditMessage(id uint64, text string, now time.Time) error {
	_, err := cd.pool.Exec(context.Background(),
		`UPDATE messages SET (text = $1, date = $2) WHERE id = $3`, text, now, id)

	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return nil
}

func (cd ChatDatabase) MessagesSearch(uid uint64, str string, page int) (models.MessagesSQL, error) {
	var messages models.MessagesSQL
	err := pgxscan.Select(context.Background(), cd.pool, &messages,
		`SELECT DISTINCT ON(id) id, mes_from, mes_to, text,
		date, redact, read FROM messages
		WHERE (LOWER(text) LIKE '%' || $1 || '%'
		ORDER BY date DESC
		LIMIT 6 OFFSET $2`, str, (page-1)*6)

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
		`SELECT DISTINCT ON(id) id, mes_from, mes_to, text,
		date, redact, read FROM messages
		WHERE (LOWER(text) LIKE '%' || $1 || '%'
		AND dialogue_id = $2
		ORDER BY date DESC
		LIMIT 6 OFFSET $3`, str, id, (page-1)*6)

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
