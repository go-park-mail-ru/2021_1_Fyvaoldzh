package repository

import (
	"context"
	"fmt"
	"kudago/application/microservices/chat/chat"
	"kudago/application/models"
	"kudago/pkg/logger"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgmock"
	"github.com/jackc/pgproto3/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	userId    uint64 = 1
	userId2   uint64 = 2
	elemId    uint64 = 1
	pageNum          = 1
	test_text        = "test text"
)

var testNewMessage = models.NewMessage{
	To:   userId2,
	Text: test_text,
}

func newDb(t *testing.T) chat.Repository {
	pool := setUp(t)
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)

	h := NewChatDatabase(pool, logger.NewLogger(sugar))
	return h
}
func setUp(t *testing.T) *pgxpool.Pool {
	script := &pgmock.Script{
		Steps: pgmock.AcceptUnauthenticatedConnRequestSteps(),
	}
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: ""}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.RowDescription{
		Fields: []pgproto3.FieldDescription{
			{
				Name:                 []byte("id"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         8,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("id_dialogue"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         8,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("mes_from"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         8,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("mes_to"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         8,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("text"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("date"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         60,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("redact"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         1,
				TypeModifier:         -1,
				Format:               0,
			},
			{
				Name:                 []byte("read"),
				TableOID:             0,
				TableAttributeNumber: 0,
				DataTypeOID:          23,
				DataTypeSize:         1,
				TypeModifier:         -1,
				Format:               0,
			},
		},
	}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.DataRow{
		Values: [][]byte{[]byte("1")},
	}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.CommandComplete{CommandTag: []byte("SELECT 1")}))
	script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.ReadyForQuery{TxStatus: 'I'}))
	script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Terminate{}))

	ln, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer ln.Close()

	serverErrChan := make(chan error, 1)
	go func() {
		defer close(serverErrChan)

		conn, err := ln.Accept()
		if err != nil {
			serverErrChan <- err
			return
		}
		defer conn.Close()

		err = conn.SetDeadline(time.Now().Add(time.Second))
		if err != nil {
			serverErrChan <- err
			return
		}

		err = script.Run(pgproto3.NewBackend(pgproto3.NewChunkReader(conn), conn))
		if err != nil {
			serverErrChan <- err
			return
		}
	}()

	parts := strings.Split(ln.Addr().String(), ":")
	host := parts[0]
	port := parts[1]
	connStr := fmt.Sprintf("sslmode=disable host=%s port=%s", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pool, err := pgxpool.Connect(ctx, connStr)
	require.NoError(t, err)
	return pool
}

func TestChatDatabase_GetAllDialogues(t *testing.T) {
	h := newDb(t)
	_, err := h.GetAllDialogues(userId, pageNum)

	assert.Nil(t, err)
}

func TestChatDatabase_GetMessages(t *testing.T) {
	h := newDb(t)
	_, err := h.GetMessages(elemId, pageNum)

	assert.Nil(t, err)
}

func TestChatDatabase_CheckDialogueID(t *testing.T) {
	h := newDb(t)
	is, _, err := h.CheckDialogueID(elemId)

	assert.Equal(t, is, false)
	assert.Nil(t, err)
}

func TestChatDatabase_GetEasyMessage(t *testing.T) {
	h := newDb(t)
	_, err := h.GetEasyMessage(elemId)

	assert.Nil(t, err)
}

func TestChatDatabase_DeleteDialogue(t *testing.T) {
	h := newDb(t)
	err := h.DeleteDialogue(elemId)

	assert.NotNil(t, err)
}

func TestChatDatabase_DeleteMessage(t *testing.T) {
	h := newDb(t)
	err := h.DeleteMessage(elemId)

	assert.NotNil(t, err)
}

func TestChatDatabase_SendMessage(t *testing.T) {
	h := newDb(t)
	err := h.SendMessage(elemId, &testNewMessage, userId, time.Now())

	assert.NotNil(t, err)
}

func TestChatDatabase_EditMessage(t *testing.T) {
	h := newDb(t)
	err := h.EditMessage(elemId, test_text)

	assert.NotNil(t, err)
}

func TestChatDatabase_MessagesSearch(t *testing.T) {
	h := newDb(t)
	_, err := h.MessagesSearch(userId, test_text, pageNum)

	assert.Nil(t, err)
}

func TestChatDatabase_DialogueMessagesSearch(t *testing.T) {
	h := newDb(t)
	_, err := h.DialogueMessagesSearch(userId, elemId, test_text, pageNum)

	assert.Nil(t, err)
}

func TestChatDatabase_CheckDialogueUsers(t *testing.T) {
	h := newDb(t)
	is, _, err := h.CheckDialogueUsers(userId, userId2)

	assert.Equal(t, is, false)
	assert.Nil(t, err)
}

func TestChatDatabase_CheckMessage(t *testing.T) {
	h := newDb(t)
	is, _, err := h.CheckMessage(elemId)

	assert.Equal(t, is, false)
	assert.Nil(t, err)
}

func TestChatDatabase_NewDialogue(t *testing.T) {
	h := newDb(t)
	_, err := h.NewDialogue(userId, userId2)

	assert.NotNil(t, err)
}

func TestChatDatabase_ReadMessages(t *testing.T) {
	h := newDb(t)
	err := h.ReadMessages(elemId, pageNum, userId)

	assert.NotNil(t, err)
}
