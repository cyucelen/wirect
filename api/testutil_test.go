package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/labstack/echo"
)

type handlerFunc func(ctx echo.Context) error

func createTestContext(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	return c, rec
}

func sendTestRequestToHandler(snifferMAC string, payload interface{}, handler handlerFunc) *httptest.ResponseRecorder {
	payloadJSON, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(payloadJSON))
	c, rec := createTestContext(req)
	addSnifferMACParamToContext(c, "/sniffers/:snifferMAC/packets", snifferMAC)

	handler(c)

	return rec
}

func sendTestRequestToHandlerWithRawBody(payload string, handler handlerFunc) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
	c, rec := createTestContext(req)
	handler(c)

	return rec
}

func sendTestRequestToHandlerWithEmptyJSON(handler handlerFunc) int {
	emptyJSON := `{}`
	rec := sendTestRequestToHandlerWithRawBody(emptyJSON, handler)
	return rec.Code
}

func sendTestRequestToHandlerWithCorruptedJSON(handler handlerFunc) int {
	corruptedJSON := `{"MAC":"AA:AA:AA:AA:AA:AA"`
	rec := sendTestRequestToHandlerWithRawBody(corruptedJSON, handler)
	return rec.Code
}

func addSnifferMACParamToContext(ctx echo.Context, path, snifferMAC string) {
	ctx.SetPath(path)
	ctx.SetParamNames("snifferMAC")
	ctx.SetParamValues(url.QueryEscape(snifferMAC))
}
