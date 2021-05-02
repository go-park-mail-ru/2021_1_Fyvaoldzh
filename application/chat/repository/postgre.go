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
	"github.com/labstack/gommon/log"
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
		`SELECT DISTINCT d.id as ID, d.user_1 as User_1, d.user_2 as User_2, m.id as ID_mes,
		m.mes_from as From, m.mes_to as To, m.text as Text, m.date as Date, m.redact as Redact, m.read as Read
	FROM dialogues d JOIN messages m on d.id = m.id_dialogue
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
		`SELECT id as ID, mes_from as From, mes_to as To, text as Text,
		date as Date, redact as Redact, read as Read FROM messages
	WHERE id_dialogue = $1`, id)
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

//Исправить значения на read
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

//Подумать насчет проверки валидности переданных значений(чтоб все значения были номральные)
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
		`UPDATE messages SET text = $1, redact = true WHERE id = $3`, text, id)

	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return nil
}

func (cd ChatDatabase) MessagesSearch(uid uint64, str string, page int) (models.MessagesSQL, error) {
	var messages models.MessagesSQL
	err := pgxscan.Select(context.Background(), cd.pool, &messages,
		`SELECT id as ID, mes_from as From, mes_to as To, text as Text,
		date as Date, redact as Redact, read as Read FROM messages
		WHERE (LOWER(text) LIKE '%' || $1 || '%')
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
	log.Info(uid, id, str)
	err := pgxscan.Select(context.Background(), cd.pool, &messages,
		`SELECT id as ID, mes_from as From, mes_to as To, text as Text,
		date as Date, redact as Redact, read as Read FROM messages
		WHERE (LOWER(text) LIKE '%' || $1 || '%')
		AND id_dialogue = $2
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

/*func (cd ChatDatabase) CheckDialogue(userId uint64) error {
	_, err := ud.pool.Query(context.Background(),
		`SELECT id FROM users WHERE id = $1`, userId)
	if err == sql.ErrNoRows {
		ud.logger.Debug("no rows in method IsExistingUser")
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("user does not exist"))
	}
	if err != nil {
		ud.logger.Warn(err)
		return err
	}

	return nil
}*/

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

//Как тут сразу вернуть созданный id?
func (cd ChatDatabase) NewDialogue(uid1 uint64, uid2 uint64) (uint64, error) {
	_, err := cd.pool.Exec(context.Background(),
		`INSERT INTO dialogues 
		VALUES (default, $1, $2)`,
		uid1, uid2)
	if err != nil {
		cd.logger.Warn(err)
		return 0, err
	}
	_, id, err := cd.CheckDialogue(uid1, uid2)
	if err != nil {
		cd.logger.Warn(err)
		return 0, err
	}

	return id, nil
}
