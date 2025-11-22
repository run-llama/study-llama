package auth

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestAuthorizeGetFail(t *testing.T) {
	if _, ok := os.LookupEnv("POSTGRES_CONNECTION_STRING"); !ok {
		t.Skip()
	} else {
		app := fiber.New()
		fReqCtx := fasthttp.RequestCtx{Request: *fasthttp.AcquireRequest()}
		defer fasthttp.ReleaseRequest(&fReqCtx.Request)
		c := app.AcquireCtx(&fReqCtx)
		defer app.ReleaseCtx(c)
		c.Request().SetRequestURI("/")
		c.Request().Header.SetMethod("GET")
		c.Request().Header.SetCookie("session_token", "noSession")
		_, err := AuthorizeGet(c)
		if err == nil {
			t.Error("Expected an error, got none")
		}
	}
}

func TestAuthorizePostFail(t *testing.T) {
	if _, ok := os.LookupEnv("POSTGRES_CONNECTION_STRING"); !ok {
		t.Skip()
	} else {
		app := fiber.New()
		fReqCtx := fasthttp.RequestCtx{Request: *fasthttp.AcquireRequest()}
		defer fasthttp.ReleaseRequest(&fReqCtx.Request)
		c := app.AcquireCtx(&fReqCtx)
		defer app.ReleaseCtx(c)
		c.Request().SetRequestURI("/search/gateway")
		c.Request().Header.SetMethod("POST")
		c.Request().Header.SetCookie("session_token", "noSession")
		c.Request().Header.SetCookie("session_token", "noCSRF")
		_, err := AuthorizePost(c)
		if err == nil {
			t.Error("Expected an error, got none")
		}
	}
}
