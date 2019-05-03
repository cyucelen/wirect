package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cyucelen/wirect/api/mocks"
	"github.com/cyucelen/wirect/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockSnifferDB(sniffers []model.Sniffer) *mocks.SnifferDatabase {
	mockSnifferDB := &mocks.SnifferDatabase{}
	mockSnifferDB.On("CreateSniffer", mock.AnythingOfType("*model.Sniffer")).Return(nil)
	mockSnifferDB.On("GetSniffers").Return(sniffers)
	mockSnifferDB.On("UpdateSniffer", mock.AnythingOfType("*model.Sniffer")).Return(nil)
	return mockSnifferDB
}

func createFailingMockSnifferDB() *mocks.SnifferDatabase {
	mockSnifferDB := &mocks.SnifferDatabase{}
	mockSnifferDB.On("CreateSniffer", mock.AnythingOfType("*model.Sniffer")).Return(errors.New(""))
	mockSnifferDB.On("UpdateSniffer", mock.AnythingOfType("*model.Sniffer")).Return(errors.New(""))

	return mockSnifferDB
}

func TestCreateSniffer(t *testing.T) {
	sniffer := model.Sniffer{MAC: "11:22:33:44:55:66", Name: "lib_sniffer", Location: "library"}
	snifferJSON, _ := json.Marshal(sniffer)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(snifferJSON))
	c, rec := createTestContext(req)

	mockSnifferDB := createMockSnifferDB([]model.Sniffer{})
	snifferAPI := &SnifferAPI{mockSnifferDB}

	err := snifferAPI.CreateSniffer(c)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, string(snifferJSON), strings.TrimRight(rec.Body.String(), "\n"))

	var actualSniffer model.Sniffer
	json.Unmarshal(snifferJSON, &actualSniffer)
	mockSnifferDB.AssertCalled(t, "CreateSniffer", &actualSniffer)
}

func TestCreateSnifferWithEmptyJSON(t *testing.T) {
	snifferJSON := `{}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(snifferJSON))
	c, rec := createTestContext(req)

	mockSnifferDB := createMockSnifferDB([]model.Sniffer{})
	snifferAPI := &SnifferAPI{mockSnifferDB}

	snifferAPI.CreateSniffer(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockSnifferDB.AssertNotCalled(t, "CreateSniffer", &model.Sniffer{})
}

func TestCreateSnifferWithCorruptedJSON(t *testing.T) {
	snifferJSON := `{"Tim}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(snifferJSON))
	c, rec := createTestContext(req)

	mockSnifferDB := createMockSnifferDB([]model.Sniffer{})
	snifferAPI := &SnifferAPI{mockSnifferDB}

	snifferAPI.CreateSniffer(c)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateSnifferWithFailingDB(t *testing.T) {
	sniffer := model.Sniffer{MAC: "11:22:33:44:55:66", Name: "lib_sniffer", Location: "library"}
	snifferJSON, _ := json.Marshal(sniffer)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(snifferJSON))
	c, rec := createTestContext(req)

	mockSnifferDB := createFailingMockSnifferDB()
	snifferAPI := &SnifferAPI{mockSnifferDB}

	snifferAPI.CreateSniffer(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetSniffers(t *testing.T) {
	sniffers := []model.Sniffer{
		{MAC: "AA:BB:CC:DD:EE:FF", Name: "lib_sniffer", Location: "library"},
		{MAC: "11:22:33:44:55:66", Name: "copy_sniffer", Location: "copy_center"},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c, rec := createTestContext(req)

	mockSnifferDB := createMockSnifferDB(sniffers)
	snifferAPI := &SnifferAPI{mockSnifferDB}

	err := snifferAPI.GetSniffers(c)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	var actualSniffers []model.Sniffer
	json.Unmarshal(rec.Body.Bytes(), &actualSniffers)

	assert.Equal(t, sniffers, actualSniffers)
}

func TestUpdateSniffer(t *testing.T) {
	snifferUpdate := model.Sniffer{Name: "room_sniffer", Location: "room"}
	updateSnifferJSON, _ := json.Marshal(snifferUpdate)
	snifferMAC := "11:22:33:44:55:66"

	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(updateSnifferJSON))
	c, rec := createTestContext(req)
	addSnifferMACParamToContext(c, "/sniffers/:snifferMAC", snifferMAC)

	mockSnifferDB := createMockSnifferDB([]model.Sniffer{})
	snifferAPI := &SnifferAPI{mockSnifferDB}

	err := snifferAPI.UpdateSniffer(c)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	expectedSnifferUpdate := snifferUpdate
	expectedSnifferUpdate.MAC = snifferMAC
	mockSnifferDB.AssertCalled(t, "UpdateSniffer", &expectedSnifferUpdate)
}

func TestUpdateSnifferWithEmptyJSON(t *testing.T) {
	updateSnifferJSON := `{}`

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(updateSnifferJSON))
	c, rec := createTestContext(req)

	mockSnifferDB := &mocks.SnifferDatabase{}
	snifferAPI := &SnifferAPI{mockSnifferDB}

	err := snifferAPI.UpdateSniffer(c)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockSnifferDB.AssertNotCalled(t, "UpdateSniffer", &model.Sniffer{})
}

func TestUpdateSnifferWithCorruptedJSON(t *testing.T) {
	snifferJSON := `{"Tim}`

	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(snifferJSON))
	c, rec := createTestContext(req)

	mockSnifferDB := createMockSnifferDB([]model.Sniffer{})
	snifferAPI := &SnifferAPI{mockSnifferDB}

	err := snifferAPI.UpdateSniffer(c)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateWithFailingDBUpdate(t *testing.T) {
	snifferUpdate := model.Sniffer{Name: "room_sniffer", Location: "room"}
	updateSnifferJSON, _ := json.Marshal(snifferUpdate)
	snifferMAC := "11:22:33:44:55:66"

	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(updateSnifferJSON))
	c, rec := createTestContext(req)
	addSnifferMACParamToContext(c, "/sniffers/:snifferMAC", snifferMAC)

	mockSnifferDB := createFailingMockSnifferDB()
	snifferAPI := &SnifferAPI{mockSnifferDB}

	err := snifferAPI.UpdateSniffer(c)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
