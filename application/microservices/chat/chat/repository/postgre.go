package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/tarantool/go-tarantool"
	"kudago/application/microservices/chat/chat"
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
	ttool  *tarantool.Connection
	logger logger.Logger
}

func NewChatDatabase(conn *pgxpool.Pool, ttool *tarantool.Connection, logger logger.Logger) chat.Repository {
	return &ChatDatabase{pool: conn, ttool: ttool, logger: logger}
}

func (cd ChatDatabase) GetAllDialogues(uid uint64, page int) (models.DialogueCardsSQL, error) {
	var dialogues models.DialogueCardsSQL
	err := pgxscan.Select(context.Background(), cd.pool, &dialogues,
		`SELECT * FROM (SELECT DISTINCT ON(d.id) d.id as ID, d.user_1, d.user_2, m.id as IDMes,
		m.mes_from, m.mes_to, m.text, m.date, m.redact, m.read
	FROM dialogues d JOIN messages m on d.id = m.id_dialogue
	WHERE user_1 = $1 OR user_2 = $1
	ORDER BY id, date DESC) as pl
	ORDER BY date DESC
	LIMIT $2 OFFSET $3`, uid, constants.ChatPerPage, (page-1)*constants.ChatPerPage)
	/*err := pgxscan.Select(context.Background(), cd.pool, &dialogues,
		`SELECT * FROM (SELECT DISTINCT ON(d.id) d.id as ID, d.user_1, d.user_2, m.id as IDMes,
		m.mes_from, m.mes_to, m.text, m.date, m.redact, m.read
	FROM dialogues d JOIN messages m on d.id = m.id_dialogue
	WHERE user_1 = $1 OR user_2 = $1
	ORDER BY id, date DESC) as pl
	ORDER BY date DESC`, uid)*/

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

func (cd ChatDatabase) GetAllNotifications(uid uint64, page int, now time.Time) (models.NotificationsSQL, error) {
	var notifications models.NotificationsSQL
	err := pgxscan.Select(context.Background(), cd.pool, &notifications,
		`SELECT id, type, date_create as date, read FROM notifications
		WHERE id_to = $1 AND date < $2
		ORDER BY read, date DESC
		LIMIT $3 OFFSET $4`, uid, now, constants.ChatPerPage, (page-1)*constants.ChatPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(notifications) == 0 {
		cd.logger.Debug("no rows in method GetAllEvents")
		return models.NotificationsSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return nil, err
	}

	return notifications, nil
}

func (cd ChatDatabase) GetMessages(id uint64, page int) (models.MessagesSQL, error) {
	var messages models.MessagesSQL
	err := pgxscan.Select(context.Background(), cd.pool, &messages,
		`SELECT id, mes_from, mes_to, text,
		date, redact, read FROM messages
	WHERE id_dialogue = $1 ORDER BY date DESC
	LIMIT $2 OFFSET $3`, id, constants.ChatPerPage, (page-1)*constants.ChatPerPage)
	/*err := pgxscan.Select(context.Background(), cd.pool, &messages,
		`SELECT id, mes_from, mes_to, text,
		date, redact, read FROM messages
	WHERE id_dialogue = $1 ORDER BY date`, id)*/

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
func (cd ChatDatabase) CheckDialogueID(id uint64) (bool, models.EasyDialogueMessageSQL, error) {
	var dialogue []models.EasyDialogueMessageSQL
	err := pgxscan.Select(context.Background(), cd.pool, &dialogue,
		`SELECT id as ID, user_1 as User1, user_2 as User2 FROM dialogues
	WHERE id = $1`, id)

	if errors.As(err, &pgx.ErrNoRows) || len(dialogue) == 0 {
		cd.logger.Debug("no rows in method GetEventsByCategory")
		return false, models.EasyDialogueMessageSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return false, models.EasyDialogueMessageSQL{}, nil
	}
	return true, dialogue[0], nil
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

func (cd ChatDatabase) EditMessage(id uint64, text string) error {
	_, err := cd.pool.Exec(context.Background(),
		`UPDATE messages SET text = $1, redact = true WHERE id = $2`, text, id)

	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return nil
}

func (cd ChatDatabase) MessagesSearch(uid uint64, str string, page int) (models.DialogueCardsSQL, error) {
	var dialogues models.DialogueCardsSQL
	err := pgxscan.Select(context.Background(), cd.pool, &dialogues,
		`SELECT d.id as ID, d.user_1, d.user_2, m.id as IDMes,
		m.mes_from, m.mes_to, m.text, m.date, m.redact, m.read
	FROM dialogues d JOIN messages m on d.id = m.id_dialogue
	WHERE (user_1 = $1 OR user_2 = $1) AND (LOWER(m.text) LIKE '%' || $2 || '%')
	ORDER BY date DESC
	LIMIT $3 OFFSET $4`, uid, str, constants.ChatPerPage, (page-1)*constants.ChatPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(dialogues) == 0 {
		cd.logger.Debug("no rows in method CategorySearch with searchstring " + str)
		return models.DialogueCardsSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return models.DialogueCardsSQL{}, err
	}

	return dialogues, nil
}

func (cd ChatDatabase) DialogueMessagesSearch(uid uint64, id uint64, str string, page int) (models.DialogueCardsSQL, error) {

	var dialogues models.DialogueCardsSQL
	err := pgxscan.Select(context.Background(), cd.pool, &dialogues,
		`SELECT d.id as ID, d.user_1, d.user_2, m.id as IDMes,
		m.mes_from, m.mes_to, m.text, m.date, m.redact, m.read
	FROM dialogues d JOIN messages m on d.id = m.id_dialogue
	WHERE (user_1 = $1 OR user_2 = $1) AND (LOWER(m.text) LIKE '%' || $2 || '%') AND d.id = $3
	ORDER BY date DESC
	LIMIT $4 OFFSET $5`, uid, str, id, constants.ChatPerPage, (page-1)*constants.ChatPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(dialogues) == 0 {
		cd.logger.Debug("no rows in method CategorySearch with searchstring " + str)
		return models.DialogueCardsSQL{}, nil
	}

	if err != nil {
		cd.logger.Warn(err)
		return models.DialogueCardsSQL{}, err
	}

	return dialogues, nil
}

func (cd ChatDatabase) CheckDialogueUsers(uid1 uint64, uid2 uint64) (bool, models.EasyDialogueMessageSQL, error) {
	var dialogue []models.EasyDialogueMessageSQL
	err := pgxscan.Select(context.Background(), cd.pool, &dialogue,
		`SELECT id as ID, user_1 as User1, user_2 as User2 FROM dialogues WHERE 
		(user_1 = $1 AND user_2 = $2) OR (user_1 = $2 AND user_2 = $1)`, uid1, uid2)

	if errors.As(err, &pgx.ErrNoRows) || len(dialogue) == 0 {
		return false, models.EasyDialogueMessageSQL{}, nil
	}
	if err != nil {
		cd.logger.Warn(err)
		return false, models.EasyDialogueMessageSQL{}, err
	}
	return true, dialogue[0], nil
}

func (cd ChatDatabase) CheckMessage(id uint64) (bool, models.EasyDialogueMessageSQL, error) {
	var messages []models.EasyDialogueMessageSQL
	err := pgxscan.Select(context.Background(), cd.pool, &messages,
		`SELECT id as ID, mes_from as User1, mes_to as User2 FROM messages WHERE 
		id = $1`, id)

	if errors.As(err, &pgx.ErrNoRows) || len(messages) == 0 {
		return false, models.EasyDialogueMessageSQL{}, nil
	}
	if err != nil {
		cd.logger.Warn(err)
		return false, models.EasyDialogueMessageSQL{}, err
	}
	return true, messages[0], nil

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

func (cd ChatDatabase) ReadMessages(id uint64, page int, uid uint64) (int64, error) {
	resp, err := cd.pool.Exec(context.Background(),
		`UPDATE messages SET read = true
		WHERE id in
	(SELECT id from messages where mes_to = $1 AND id_dialogue = $2 AND read = false
		ORDER BY date DESC LIMIT $3 OFFSET $4)`,
		uid, id, constants.ChatPerPage, (page-1)*constants.ChatPerPage)
	/*_, err := cd.pool.Exec(context.Background(),
		`UPDATE messages SET read = true
		WHERE id in
	(SELECT id from messages where mes_to = $1 AND id_dialogue = $2
		ORDER BY date)`, uid, id)*/

	if err != nil {
		cd.logger.Warn(err)
		return 0, err
	}

	return resp.RowsAffected(), nil
}

//По просьбе фронта делаем прочитанными все уведомления
func (cd ChatDatabase) ReadNotifications(uid uint64, page int, now time.Time) error {
	_, err := cd.pool.Exec(context.Background(),
		`UPDATE notifications SET read = true
		WHERE id_to = $1 AND date <= $2`, uid, now)

	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return nil
}

func (cd ChatDatabase) AddMailNotification(id uint64, idTo uint64, now time.Time) error {
	_, err := cd.pool.Exec(context.Background(),
		`INSERT INTO notifications 
		VALUES ($1, $2, $3, $4, $4, default)`,
		id, constants.MailNotif, idTo, now)
	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return nil
}

func (cd ChatDatabase) AddCountNotification(id uint64) error {
	_, err := cd.ttool.Update(constants.TarantoolSpaceName2, "primary",
		[]interface{}{id}, []interface{}{[]interface{}{"+", constants.TarantoolNotifications, 1}})
	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return nil
}

func (cd ChatDatabase) SetZeroCountNotifications(id uint64) error {
	resp, err := cd.ttool.Select(constants.TarantoolSpaceName2,
		"primary", 0, 1, tarantool.IterEq, []interface{}{id})

	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	data := resp.Data[0]
	d, ok := data.([]interface{})
	if !ok {
		return errors.New("cast error")
	}

	count, ok := d[2].(uint64)
	if !ok {
		return errors.New("cast error")
	}

	_, err = cd.ttool.Replace(constants.TarantoolSpaceName2, []interface{}{id, 0, count})
	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return errors.New("got no data in tarantool")
}

func (cd ChatDatabase) AddCountMessages(id uint64) error {
	_, err := cd.ttool.Update(constants.TarantoolSpaceName2, "primary",
		[]interface{}{id}, []interface{}{[]interface{}{"+", constants.TarantoolMessages, 1}})
	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return nil
}

func (cd ChatDatabase) DecrementCountMessages(id uint64, count int64) error {
	_, err := cd.ttool.Update(constants.TarantoolSpaceName2, "primary",
		[]interface{}{id}, []interface{}{[]interface{}{"-", constants.TarantoolMessages, count}})
	if err != nil {
		cd.logger.Warn(err)
		return err
	}

	return nil
}

func (cd ChatDatabase) GetAllCounts(uid uint64) (models.Counts, error) {
	var counts models.Counts
	resp, err := cd.ttool.Select(constants.TarantoolSpaceName2,
		"primary", 0, 1, tarantool.IterEq, []interface{}{uid})

	if err != nil {
		cd.logger.Warn(err)
		return models.Counts{}, err
	}

	data := resp.Data[0]
	d, ok := data.([]interface{})
	if !ok {
		return models.Counts{}, errors.New("cast error")
	}

	counts.Notifications, ok = d[1].(uint64)
	if !ok {
		_ = errors.New("cast error")
	}
	counts.Chat, ok = d[2].(uint64)
	if !ok {
		_ = errors.New("cast error")
	}

	return counts, nil
}

func (cd ChatDatabase) GetNotificationCounts(uid uint64, now time.Time) (uint64, error) {
	var counts []uint64
	err := pgxscan.Select(context.Background(), cd.pool, &counts,
		`SELECT count(id) FROM notifications
		WHERE id_to = $1 AND date < $2 AND read = false`, uid, now)

	if err != nil {
		cd.logger.Warn(err)
		return 0, err
	}

	return counts[0], nil
}
