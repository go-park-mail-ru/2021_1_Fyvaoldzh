package middleware

/*
import (
	"kudago/application/microservices/auth/client"
	"kudago/pkg/constants"
	"kudago/pkg/generator"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var (
	test_uid = uint64(1)
)

func setUp(t *testing.T) (echo.Context, *client.MockIAuthClient, Auth) {
	e := echo.New()

	ctrl := gomock.NewController(t)
	rpcAuth := client.NewMockIAuthClient(ctrl)

	auth := NewAuth(rpcAuth)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/", nil)
	c := e.NewContext(req, rec)

	return c, rpcAuth, auth
}

func TestUserId_GetSession(t *testing.T) {
	c, rpcAuth, auth := setUp(t)
	key := generator.RandStringRunes(constants.CookieLength)
	newCookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    key,
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
	}
	c.Request().AddCookie(newCookie)
	rpcAuth.EXPECT().Check(key).Return(true, test_uid, nil, 200)

	auth.GetSession(c.Handler())
	uid := c.Get(constants.UserIdKey).(uint64)

	assert.Equal(t, uid, test_uid)
}


 */