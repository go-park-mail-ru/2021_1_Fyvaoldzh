package http

/*
import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	mock_user "kudago/application/user/mocks"
	"testing"
)

func beforeTest(t *testing.T) (*UserHandler, *mock_user.MockUseCase, *server.MockContext, *server.MockResponseWriter) {
	ctrl := gomock.NewController(t)
	w := server.NewMockResponseWriter(ctrl)
	ctx := server.NewMockContext(ctrl)
	usecase := mock_user.NewMockUseCase(ctrl)
	handler := UserHandler{
		UseCase: usecase,
		Logger: zap.NewExample(), sanitizer: bluemonday.UGCPolicy()}

	response := echo.NewResponse(w, echo.New())
	ctx.EXPECT().Response().Return(response).AnyTimes()
	w.EXPECT().Write(gomock.Any()).AnyTimes()

	return &handler, usecase, ctx, w
}

 */