package api

import (
	"encoding/json"
	"errors"
	"net/http"
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
	mockSnifferDB := createMockSnifferDB([]model.Sniffer{})
	snifferAPI := &SnifferAPI{mockSnifferDB}

	expectedSniffer := model.Sniffer{MAC: "11:22:33:44:55:66", Name: "lib_sniffer", Description: "library"}

	rec := sendTestRequestToHandler("", expectedSniffer, snifferAPI.CreateSniffer, http.MethodPost)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var actualSniffer model.Sniffer
	json.NewDecoder(rec.Body).Decode(&actualSniffer)
	assert.Equal(t, expectedSniffer, actualSniffer)
	mockSnifferDB.AssertCalled(t, "CreateSniffer", &expectedSniffer)
}

func TestCreateSnifferWithEmptyJSON(t *testing.T) {
	mockSnifferDB := createMockSnifferDB([]model.Sniffer{})
	snifferAPI := &SnifferAPI{mockSnifferDB}

	responseCode := sendTestRequestToHandlerWithEmptyJSON(snifferAPI.CreateSniffer)
	assert.Equal(t, http.StatusBadRequest, responseCode)
	mockSnifferDB.AssertNotCalled(t, "CreateSniffer", &model.Sniffer{})
}

func TestCreateSnifferWithCorruptedJSON(t *testing.T) {
	mockSnifferDB := createMockSnifferDB([]model.Sniffer{})
	snifferAPI := &SnifferAPI{mockSnifferDB}

	responseCode := sendTestRequestToHandlerWithCorruptedJSON(snifferAPI.CreateSniffer)
	assert.Equal(t, http.StatusBadRequest, responseCode)
}

func TestCreateSnifferWithFailingDB(t *testing.T) {
	mockSnifferDB := createFailingMockSnifferDB()
	snifferAPI := &SnifferAPI{mockSnifferDB}

	sniffer := model.Sniffer{MAC: "11:22:33:44:55:66", Name: "lib_sniffer", Description: "library"}

	rec := sendTestRequestToHandler("", sniffer, snifferAPI.CreateSniffer, http.MethodPost)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetSniffers(t *testing.T) {
	expectedSniffers := []model.Sniffer{
		{MAC: "AA:BB:CC:DD:EE:FF", Name: "lib_sniffer", Description: "library"},
		{MAC: "11:22:33:44:55:66", Name: "copy_sniffer", Description: "copy_center"},
	}
	mockSnifferDB := createMockSnifferDB(expectedSniffers)
	snifferAPI := &SnifferAPI{mockSnifferDB}

	rec := sendTestRequestToHandler("", nil, snifferAPI.GetSniffers, http.MethodGet)
	var actualSniffers []model.Sniffer
	json.NewDecoder(rec.Body).Decode(&actualSniffers)
	assert.Equal(t, expectedSniffers, actualSniffers)
}

func TestUpdateSniffer(t *testing.T) {
	mockSnifferDB := createMockSnifferDB([]model.Sniffer{})
	snifferAPI := &SnifferAPI{mockSnifferDB}

	snifferUpdate := model.Sniffer{Name: "room_sniffer", Description: "room"}
	snifferMAC := "11:22:33:44:55:66"

	rec := sendTestRequestToHandler(snifferMAC, snifferUpdate, snifferAPI.UpdateSniffer, http.MethodPut)
	assert.Equal(t, http.StatusOK, rec.Code)

	expectedSnifferUpdate := snifferUpdate
	expectedSnifferUpdate.MAC = snifferMAC
	mockSnifferDB.AssertCalled(t, "UpdateSniffer", &expectedSnifferUpdate)
}

func TestUpdateSnifferWithInvalidSnifferMACParam(t *testing.T) {
	mockSnifferDB := createMockSnifferDB([]model.Sniffer{})
	snifferAPI := &SnifferAPI{mockSnifferDB}

	snifferUpdate := model.Sniffer{Name: "room_sniffer", Description: "room"}
	rec := sendTestRequestToHandlerWithInvalidParam(snifferUpdate, snifferAPI.UpdateSniffer)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateSnifferWithEmptyJSON(t *testing.T) {
	mockSnifferDB := &mocks.SnifferDatabase{}
	snifferAPI := &SnifferAPI{mockSnifferDB}

	responseCode := sendTestRequestToHandlerWithEmptyJSON(snifferAPI.UpdateSniffer)
	assert.Equal(t, http.StatusNotFound, responseCode)
	mockSnifferDB.AssertNotCalled(t, "UpdateSniffer", &model.Sniffer{})
}

func TestUpdateSnifferWithCorruptedJSON(t *testing.T) {
	mockSnifferDB := createMockSnifferDB([]model.Sniffer{})
	snifferAPI := &SnifferAPI{mockSnifferDB}

	responseCode := sendTestRequestToHandlerWithCorruptedJSON(snifferAPI.UpdateSniffer)
	assert.Equal(t, http.StatusBadRequest, responseCode)
	mockSnifferDB.AssertNotCalled(t, "UpdateSniffer", &model.Sniffer{})
}

func TestUpdateWithFailingDBUpdate(t *testing.T) {
	mockSnifferDB := createFailingMockSnifferDB()
	snifferAPI := &SnifferAPI{mockSnifferDB}

	snifferUpdate := model.Sniffer{Name: "room_sniffer", Description: "room"}
	snifferMAC := "11:22:33:44:55:66"

	rec := sendTestRequestToHandler(snifferMAC, snifferUpdate, snifferAPI.UpdateSniffer, http.MethodPut)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
