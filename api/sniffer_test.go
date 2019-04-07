package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/wirect/wirect-server/api/mocks"
	"gitlab.com/wirect/wirect-server/model"
)

func createTestContext(req *http.Request) (echo.Context, *httptest.ResponseRecorder) {
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e := echo.New()
	c := e.NewContext(req, rec)

	return c, rec
}
func TestCreateSniffer(t *testing.T) {
	mockSnifferDB := &mocks.SnifferDatabase{}
	mockSnifferDB.On("CreateSniffer", mock.AnythingOfType("*model.Sniffer")).Return(nil)
	snifferAPI := &SnifferAPI{mockSnifferDB}

	snifferJSON := `{"MAC":"11:22:33:44:55:66","name":"lib_sniffer","location":"library"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(snifferJSON))

	c, rec := createTestContext(req)

	snifferAPI.CreateSniffer(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, snifferJSON, strings.TrimRight(rec.Body.String(), "\n"))

	var sniffer model.Sniffer
	json.Unmarshal([]byte(snifferJSON), &sniffer)
	mockSnifferDB.AssertCalled(t, "CreateSniffer", &sniffer)
}

func TestCreateSnifferWithEmptyJSON(t *testing.T) {
	mockSnifferDB := &mocks.SnifferDatabase{}
	mockSnifferDB.On("CreateSniffer", mock.AnythingOfType("*model.Sniffer")).Return(nil)
	snifferAPI := &SnifferAPI{mockSnifferDB}

	snifferJSON := `{}`

	req := httptest.NewRequest(http.MethodPost, "/sniffer", strings.NewReader(snifferJSON))
	c, rec := createTestContext(req)

	snifferAPI.CreateSniffer(c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockSnifferDB.AssertNotCalled(t, "CreateSniffer", &model.Sniffer{})
}

func TestGetSniffers(t *testing.T) {
	sniffers := []model.Sniffer{
		{MAC: "AA:BB:CC:DD:EE:FF", Name: "lib_sniffer", Location: "library"},
		{MAC: "11:22:33:44:55:66", Name: "copy_sniffer", Location: "copy_center"},
	}

	mockSnifferDB := &mocks.SnifferDatabase{}
	mockSnifferDB.On("GetSniffers").Return(sniffers)
	snifferAPI := &SnifferAPI{mockSnifferDB}

	req := httptest.NewRequest(http.MethodGet, "/sniffer", nil)
	c, rec := createTestContext(req)

	snifferAPI.GetSniffers(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var actualSniffers []model.Sniffer
	json.Unmarshal(rec.Body.Bytes(), &actualSniffers)

	assert.Equal(t, sniffers, actualSniffers)
}

func TestUpdateSniffer(t *testing.T) {
	snifferUpdate := model.Sniffer{MAC: "11:22:33:44:55:66", Name: "room_sniffer", Location: "room"}
	updateSnifferJSON, _ := json.Marshal(snifferUpdate)

	mockSnifferDB := &mocks.SnifferDatabase{}
	mockSnifferDB.On("UpdateSniffer", mock.AnythingOfType("*model.Sniffer")).Return(nil)

	snifferAPI := &SnifferAPI{mockSnifferDB}

	req := httptest.NewRequest(http.MethodPut, "/sniffer", bytes.NewReader(updateSnifferJSON))
	c, rec := createTestContext(req)

	snifferAPI.UpdateSniffer(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockSnifferDB.AssertCalled(t, "UpdateSniffer", &snifferUpdate)
}

func TestUpdateSnifferWithEmptyJSON(t *testing.T) {
	updateSnifferJSON := `{}`
	mockSnifferDB := &mocks.SnifferDatabase{}

	snifferAPI := &SnifferAPI{mockSnifferDB}

	req := httptest.NewRequest(http.MethodPut, "/sniffer", strings.NewReader(updateSnifferJSON))
	c, rec := createTestContext(req)

	snifferAPI.UpdateSniffer(c)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockSnifferDB.AssertNotCalled(t, "UpdateSniffer", &model.Sniffer{})
}

func TestUpdateWithFailingDBUpdate(t *testing.T) {
	snifferUpdate := model.Sniffer{MAC: "11:22:33:44:55:66", Name: "room_sniffer", Location: "room"}
	updateSnifferJSON, _ := json.Marshal(snifferUpdate)

	mockSnifferDB := &mocks.SnifferDatabase{}
	mockSnifferDB.On("UpdateSniffer", mock.AnythingOfType("*model.Sniffer")).Return(errors.New(""))
	snifferAPI := &SnifferAPI{mockSnifferDB}

	req := httptest.NewRequest(http.MethodPut, "/sniffer", bytes.NewReader(updateSnifferJSON))
	c, rec := createTestContext(req)

	snifferAPI.UpdateSniffer(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
