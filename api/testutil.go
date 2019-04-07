package api

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
)

func createTestContext(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	return c, rec
}
