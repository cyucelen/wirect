package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/cyucelen/wirect/api/mocks"
	"github.com/cyucelen/wirect/model"
	"github.com/stretchr/testify/assert"

	"github.com/cyucelen/wirect/test"
)

type RouterAPISuite struct {
	suite.Suite
	routerAPI *RouterAPI
	routerDB  *test.InMemoryDB
}

func TestRouterAPISuite(t *testing.T) {
	suite.Run(t, new(RouterAPISuite))
}

func (s *RouterAPISuite) BeforeTest(string, string) {
	s.routerDB = &test.InMemoryDB{}
	s.routerAPI = &RouterAPI{DB: s.routerDB}
}

func createMockRouterDB(routers []model.Router) *mocks.RouterDatabase {
	mockRouterDB := &mocks.RouterDatabase{}
	mockRouterDB.On("GetRoutersBySniffer", mock.Anything).Return(routers)
	return mockRouterDB
}

func createFailingMockRouterDB() *mocks.RouterDatabase {
	mockRouterDatabase := &mocks.RouterDatabase{}
	mockRouterDatabase.On("CreateRouter", mock.AnythingOfType("*model.Router")).Return(errors.New(""))
	return mockRouterDatabase
}

func (s *RouterAPISuite) TestCreateRouters() {
	db := &test.InMemoryDB{}
	routerAPI := RouterAPI{DB: db}

	routers := []model.RouterExternal{
		{SSID: "1010", LastSeen: 1000},
		{SSID: "2020", LastSeen: 1200},
	}

	rec := sendTestRequestToHandler(defaultTestSnifferMAC, routers, routerAPI.CreateRouters, http.MethodPost)
	assert.Equal(s.T(), http.StatusCreated, rec.Code)

	expectedRouters := []model.RouterExternal{}

	for _, router := range db.Routers {
		expectedRouters = append(expectedRouters, *toExternal(&router))
	}

	assert.Equal(s.T(), routers, expectedRouters)

	var actualResponse []model.RouterExternal
	json.NewDecoder(rec.Body).Decode(&actualResponse)

	assert.Equal(s.T(), routers, actualResponse)
}

func (s *RouterAPISuite) TestCreateRoutersWithCorruptedJSON() {
	db := &test.InMemoryDB{}
	routerAPI := RouterAPI{DB: db}

	responseStatusCode := sendTestRequestToHandlerWithEmptyJSON(routerAPI.CreateRouters)
	assert.Equal(s.T(), http.StatusBadRequest, responseStatusCode)
}

// func (s *RouterAPISuite) TestCreateRoutersWithEmptyRequiredFileds() {
// 	onlyNotValidRouters := []model.RouterExternal{
// 		{MAC: "", SSID: "1010"},
// 		{MAC: "", SSID: "2020"},
// 	}

// 	rec := sendTestRequestToHandler(defaultTestSnifferMAC, onlyNotValidRouters, s.routerAPI.CreateRouters, http.MethodPost)
// 	assert.Len(s.T(), s.routerDB.Routers, 0)
// 	assert.Equal(s.T(), http.StatusBadRequest, rec.Code)

// 	validAndNotValidRouters := []model.RouterExternal{
// 		{MAC: "", SSID: "1010"},
// 		{MAC: "CC:FF:AA:EE:FF:FF", SSID: "2020"},
// 	}

// 	rec = sendTestRequestToHandler(defaultTestSnifferMAC, validAndNotValidRouters, s.routerAPI.CreateRouters, http.MethodPost)
// 	assert.Len(s.T(), s.routerDB.Routers, 1)
// 	assert.Equal(s.T(), http.StatusCreated, rec.Code)

// 	var actualResponse []model.RouterExternal
// 	json.NewDecoder(rec.Body).Decode(&actualResponse)
// 	assert.Len(s.T(), actualResponse, 1)
// 	assert.Equal(s.T(), validAndNotValidRouters[1], actualResponse[0])
// }

func (s *RouterAPISuite) TestCreateRoutersWithInvalidSnifferMACParam() {
	db := &test.InMemoryDB{}
	routerAPI := RouterAPI{DB: db}

	routers := []model.RouterExternal{
		{SSID: "1010"},
		{SSID: "2020"},
	}

	rec := sendTestRequestToHandlerWithInvalidParam(routers, routerAPI.CreateRouters)
	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
}

func (s *RouterAPISuite) TestCreateRoutersWithFailingDB() {
	db := createFailingMockRouterDB()
	routerAPI := RouterAPI{DB: db}

	routers := []model.RouterExternal{
		{SSID: "1010"},
		{SSID: "2020"},
	}
	rec := sendTestRequestToHandler(defaultTestSnifferMAC, routers, routerAPI.CreateRouters, http.MethodPost)
	assert.Equal(s.T(), http.StatusInternalServerError, rec.Code)
}

func (s *RouterAPISuite) TestGetRouters() {
	routers := []model.Router{
		{SSID: "1010", SnifferMAC: defaultTestSnifferMAC},
		{SSID: "2020", SnifferMAC: defaultTestSnifferMAC},
	}
	db := createMockRouterDB(routers)
	routerAPI := RouterAPI{DB: db}
	rec := sendTestRequestToHandler(defaultTestSnifferMAC, nil, routerAPI.GetRouters, http.MethodGet)

	expectedRouters := []model.RouterExternal{
		{SSID: "1010"},
		{SSID: "2020"},
	}

	assert.Equal(s.T(), http.StatusOK, rec.Code)

	var actualResponse []model.RouterExternal
	json.NewDecoder(rec.Body).Decode(&actualResponse)
	assert.Equal(s.T(), expectedRouters, actualResponse)
}

func (s *RouterAPISuite) TestGetRoutersWithInvalidSnifferMACParam() {
	db := &test.InMemoryDB{}
	routerAPI := RouterAPI{DB: db}

	routers := []model.RouterExternal{
		{SSID: "1010"},
		{SSID: "2020"},
	}

	rec := sendTestRequestToHandlerWithInvalidParam(routers, routerAPI.GetRouters)
	assert.Equal(s.T(), http.StatusNotFound, rec.Code)
}
