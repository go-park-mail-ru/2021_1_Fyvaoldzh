package http

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	mock_user "kudago/application/user/mocks"
	"kudago/pkg/custom_sanitizer"
	"kudago/pkg/logger"
	"log"
	"net/http"
	"testing"
)

func setUp(t *testing.T) (*UserHandler, *mock_user.MockUseCase) {
	ctrl := gomock.NewController(t)

	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	sugar := l.Sugar()
	zap.NewAtomicLevelAt(zapcore.DebugLevel)
	usecase := mock_user.NewMockUseCase(ctrl)
	cs := custom_sanitizer.NewCustomSanitizer(bluemonday.UGCPolicy())
	handler := UserHandler{
		UseCase: usecase,
		Logger: logger.NewLogger(sugar),
		sanitizer: cs,
	}


	return &handler, usecase
}


func TestUserHandler_CheckUser_GetById(t *testing.T) {
	handler, usecase := setUp(t)
	request, err := http.NewRequest("", "", bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Param(param).Return(id)
	usecase.EXPECT().GetById(testPerson.Id).Return(&testPerson, nil)
	w.EXPECT().WriteHeader(ok)

	err = handler.GetById(ctx)

	assert.Nil(t, err)
}
