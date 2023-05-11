package httpgin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/ports/httpgin/mocks"
	"homework10/internal/users"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

type adData struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	AuthorID     int64     `json:"author_id"`
	Published    bool      `json:"published"`
	DateUpdate   time.Time `json:"date_update"`
	DateCreating time.Time `json:"date_creating"`
}

type adDataResponse struct {
	Data adData `json:"data"`
}

type deleteResponse struct {
	Data string `json:"data"`
}

type adsResponse struct {
	Data []adData `json:"data"`
}

var (
	ErrBadRequest = fmt.Errorf("bad request")
	ErrForbidden  = fmt.Errorf("forbidden")
	ErrNotFound   = fmt.Errorf("not found")
)

type testClient struct {
	client  *http.Client
	baseURL string
}

func getTestClient(service app.App) *testClient {
	server := NewHTTPServer(":18080", service)
	testServer := httptest.NewServer(server.Handler)

	return &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
}

func (tc *testClient) getResponse(req *http.Request, out any) error {
	resp, err := tc.client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return ErrBadRequest
		}
		if resp.StatusCode == http.StatusForbidden {
			return ErrForbidden
		}
		if resp.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response: %w", err)
	}

	err = json.Unmarshal(respBody, out)
	if err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}

	return nil
}

func (tc *testClient) createAd(userID int64, title string, text string) (adDataResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adDataResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/ads", bytes.NewReader(data))
	if err != nil {
		return adDataResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adDataResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adDataResponse{}, err
	}

	return response, nil
}

func (tc *testClient) changeAdStatus(userID int64, adID int64, published bool) (adDataResponse, error) {
	body := map[string]any{
		"user_id":   userID,
		"published": published,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adDataResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d/status", adID), bytes.NewReader(data))
	if err != nil {
		return adDataResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adDataResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adDataResponse{}, err
	}

	return response, nil
}

func (tc *testClient) updateAd(userID int64, adID int64, title string, text string) (adDataResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adDataResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adID), bytes.NewReader(data))
	if err != nil {
		return adDataResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adDataResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adDataResponse{}, err
	}

	return response, nil
}

func (tc *testClient) listAdsByTitle(title string) (adsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads/search/"+title, nil)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

func (tc *testClient) getListAdsWithFilter(filters map[string]any) (adsResponse, error) {
	v := url.Values{}
	for str, filter := range filters {
		v.Add(str, fmt.Sprintf("%v", filter))
	}
	queryString := v.Encode()

	req, err := http.NewRequest(http.MethodGet, tc.baseURL+"/api/v1/ads/with_filter"+"?"+queryString, nil)
	if err != nil {
		return adsResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adsResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adsResponse{}, err
	}

	return response, nil
}

type userData struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type userDataResponse struct {
	Data userData `json:"data"`
}

func (tc *testClient) createUser(nickname string, email string) (userDataResponse, error) {
	body := map[string]any{
		"nickname": nickname,
		"email":    email,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userDataResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, tc.baseURL+"/api/v1/users", bytes.NewReader(data))
	if err != nil {
		return userDataResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userDataResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userDataResponse{}, err
	}

	return response, nil
}

func (tc *testClient) updateUser(userId int64, nickname string, email string) (userDataResponse, error) {
	body := map[string]any{
		"nickname": nickname,
		"email":    email,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userDataResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", userId), bytes.NewReader(data))
	if err != nil {
		return userDataResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userDataResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userDataResponse{}, err
	}

	return response, nil
}

func (tc *testClient) getAdById(adId int64) (adDataResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adId), nil)
	if err != nil {
		return adDataResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response adDataResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adDataResponse{}, err
	}

	return response, nil
}

func (tc *testClient) getUserById(userId int64) (userDataResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", userId), nil)
	if err != nil {
		return userDataResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response userDataResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userDataResponse{}, err
	}

	return response, nil
}

func (tc *testClient) deleteAdById(adId int64, userId int64) (deleteResponse, error) {
	body := map[string]any{
		"user_id": userId,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return deleteResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%d", adId), bytes.NewReader(data))
	if err != nil {
		return deleteResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response deleteResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return deleteResponse{}, err
	}

	return response, nil
}

func (tc *testClient) deleteUserById(userId int64) (deleteResponse, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/users/%d", userId), nil)
	if err != nil {
		return deleteResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	var response deleteResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return deleteResponse{}, err
	}

	return response, nil
}

type AdServiceTestSuite struct {
	suite.Suite
	app mocks.App
}

func (s *AdServiceTestSuite) SetupTest() {
	s.app = mocks.App{}
}

func TestRepoRun(t *testing.T) {
	suite.Run(t, new(AdServiceTestSuite))
}

func EqualAds(response *adData, ad *ads.Ad) bool {
	if response.Published != ad.Published {
		return false
	}
	if response.ID != ad.ID {
		return false
	}
	if response.Title != ad.Title {
		return false
	}
	if response.DateUpdate != ad.DateUpdate {
		return false
	}
	if response.Text != ad.Text {
		return false
	}
	if response.AuthorID != ad.AuthorID {
		return false
	}
	if response.DateCreating != ad.DateCreating {
		return false
	}
	return true
}

func EqualAdsLists(data []adData, ads []ads.Ad) bool {
	for idx := range ads {
		EqualAds(&data[idx], &ads[idx])
	}
	return true
}

func EqualUsers(response *userData, ad *users.User) bool {
	if response.ID != ad.ID {
		return false
	}
	if response.Nickname != ad.Nickname {
		return false
	}
	if response.Email != ad.Email {
		return false
	}
	return true
}

func (s *AdServiceTestSuite) TestAdService_CreateAd() {
	expect := &ads.Ad{ID: 1, Title: "title 1", Text: "text 1", AuthorID: 1}
	s.app.On("CreateAd", expect.Title, expect.Text, expect.AuthorID).Return(expect, nil)

	client := getTestClient(&s.app)

	got, err := client.createAd(expect.AuthorID, expect.Title, expect.Text)
	s.NoError(err)
	s.True(EqualAds(&got.Data, expect))
}

func (s *AdServiceTestSuite) TestAdService_CreateAdIncorrectUserId() {
	expect := &ads.Ad{ID: 1, Title: "title 1", Text: "text 1", AuthorID: 1}
	s.app.On("CreateAd", expect.Title, expect.Text, expect.AuthorID).Return(expect, app.IncorrectUserId)

	client := getTestClient(&s.app)

	_, err := client.createAd(expect.AuthorID, expect.Title, expect.Text)
	s.ErrorIs(err, ErrForbidden)
}

func (s *AdServiceTestSuite) TestAdService_CreateAdValidationErr() {
	expect := &ads.Ad{ID: 1, Title: "", Text: "text 1", AuthorID: 1}
	s.app.On("CreateAd", expect.Title, expect.Text, expect.AuthorID).Return(expect, app.ValidateError)

	client := getTestClient(&s.app)

	_, err := client.createAd(expect.AuthorID, expect.Title, expect.Text)
	s.ErrorIs(err, ErrBadRequest)
}

func (s *AdServiceTestSuite) TestAdService_CreateUser() {
	expect := &users.User{ID: 1, Nickname: "nickname", Email: "email"}
	s.app.On("CreateUser", expect.Nickname, expect.Email).Return(expect, nil)

	client := getTestClient(&s.app)

	got, err := client.createUser(expect.Nickname, expect.Email)
	s.NoError(err)
	s.True(EqualUsers(&got.Data, expect))
}

func (s *AdServiceTestSuite) TestAdService_CreateUserValidationErr() {
	expect := &users.User{ID: 1, Nickname: "", Email: "email"}
	s.app.On("CreateUser", expect.Nickname, expect.Email).Return(expect, app.ValidateError)

	client := getTestClient(&s.app)

	_, err := client.createUser(expect.Nickname, expect.Email)
	s.ErrorIs(err, ErrBadRequest)
}

func (s *AdServiceTestSuite) TestAdService_ChangeAdStatus() {
	expect := &ads.Ad{ID: 1, Title: "title 1", Text: "text 1", AuthorID: 1}
	s.app.On("ChangeAdStatus", expect.ID, expect.AuthorID, expect.Published).Return(expect, nil)

	client := getTestClient(&s.app)

	got, err := client.changeAdStatus(expect.AuthorID, expect.ID, expect.Published)
	s.NoError(err)
	s.True(EqualAds(&got.Data, expect))
}

func (s *AdServiceTestSuite) TestAdService_ChangeAdStatusIncorrectUserId() {
	expect := &ads.Ad{ID: 1, Title: "title 1", Text: "text 1", AuthorID: 1}
	s.app.On("ChangeAdStatus", expect.ID, expect.AuthorID, expect.Published).Return(expect, app.IncorrectUserId)

	client := getTestClient(&s.app)

	_, err := client.changeAdStatus(expect.AuthorID, expect.ID, expect.Published)
	s.ErrorIs(err, ErrForbidden)
}

func (s *AdServiceTestSuite) TestAdService_UpdateAd() {
	expect := &ads.Ad{ID: 1, Title: "title 1", Text: "text 1", AuthorID: 1}
	s.app.On("UpdateAd", expect.ID, expect.AuthorID, expect.Title, expect.Text).Return(expect, nil)

	client := getTestClient(&s.app)

	got, err := client.updateAd(expect.AuthorID, expect.ID, expect.Title, expect.Text)
	s.NoError(err)
	s.True(EqualAds(&got.Data, expect))
}

func (s *AdServiceTestSuite) TestAdService_UpdateAdIncorrectUserId() {
	expect := &ads.Ad{ID: 1, Title: "title 1", Text: "text 1", AuthorID: 1}
	s.app.On("UpdateAd", expect.ID, expect.AuthorID, expect.Title, expect.Text).Return(expect, app.IncorrectUserId)

	client := getTestClient(&s.app)

	_, err := client.updateAd(expect.AuthorID, expect.ID, expect.Title, expect.Text)
	s.ErrorIs(err, ErrForbidden)
}

func (s *AdServiceTestSuite) TestAdService_UpdateAdValidationErr() {
	expect := &ads.Ad{ID: 1, Title: "", Text: "text 1", AuthorID: 1}
	s.app.On("UpdateAd", expect.ID, expect.AuthorID, expect.Title, expect.Text).Return(expect, app.ValidateError)

	client := getTestClient(&s.app)

	_, err := client.updateAd(expect.AuthorID, expect.ID, expect.Title, expect.Text)
	s.ErrorIs(err, ErrBadRequest)
}

func (s *AdServiceTestSuite) TestAdService_UpdateUser() {
	expect := &users.User{ID: 1, Nickname: "nickname", Email: "email"}
	s.app.On("UpdateUser", expect.ID, expect.Nickname, expect.Email).Return(expect, nil)

	client := getTestClient(&s.app)

	got, err := client.updateUser(expect.ID, expect.Nickname, expect.Email)
	s.NoError(err)
	s.True(EqualUsers(&got.Data, expect))
}

func (s *AdServiceTestSuite) TestAdService_UpdateUserValidationErr() {
	expect := &users.User{ID: 1, Nickname: "", Email: "email"}
	s.app.On("UpdateUser", expect.ID, expect.Nickname, expect.Email).Return(expect, app.ValidateError)

	client := getTestClient(&s.app)

	_, err := client.updateUser(expect.ID, expect.Nickname, expect.Email)
	s.ErrorIs(err, ErrBadRequest)
}

func (s *AdServiceTestSuite) TestAdService_UpdateUserIncorrectUserId() {
	expect := &users.User{ID: 1, Nickname: "", Email: "email"}
	s.app.On("UpdateUser", expect.ID, expect.Nickname, expect.Email).Return(expect, app.IncorrectUserId)

	client := getTestClient(&s.app)

	_, err := client.updateUser(expect.ID, expect.Nickname, expect.Email)
	s.ErrorIs(err, ErrForbidden)
}

func (s *AdServiceTestSuite) TestAdService_GetListAds() {
	userId := int64(1)
	published := true
	dateCreating := "2023-04-29"

	expect1 := ads.Ad{ID: 1, Title: "title 1", Text: "text 1", AuthorID: 1, Published: true}
	expect2 := ads.Ad{ID: 2, Title: "title 2", Text: "text 2", AuthorID: 1, Published: true}
	adsList := []ads.Ad{expect1, expect2}

	filters := map[string]any{"user_id": userId, "published": published, "date_creating": dateCreating}
	s.app.On("GetListAds", filters).Return(adsList, nil)

	client := getTestClient(&s.app)
	response, err := client.getListAdsWithFilter(filters)
	s.NoError(err)
	s.True(EqualAdsLists(response.Data, adsList))
}

func (s *AdServiceTestSuite) TestAdService_GetAd() {
	expect := &ads.Ad{ID: 1, Title: "", Text: "text 1", AuthorID: 1}
	s.app.On("GetAd", expect.ID).Return(expect, nil)

	client := getTestClient(&s.app)

	response, err := client.getAdById(expect.ID)
	s.NoError(err)
	s.True(EqualAds(&response.Data, expect))
}

func (s *AdServiceTestSuite) TestAdService_GetAdNotFound() {
	expect := &ads.Ad{ID: 1, Title: "", Text: "text 1", AuthorID: 1}
	s.app.On("GetAd", expect.ID).Return(expect, app.IncorrectAdId)

	client := getTestClient(&s.app)

	_, err := client.getAdById(expect.ID)
	s.ErrorIs(err, ErrNotFound)
}

func (s *AdServiceTestSuite) TestAdService_GetListAdsByTitle() {
	expect1 := ads.Ad{ID: 1, Title: "title", Text: "text 1", AuthorID: 1, Published: true}
	expect2 := ads.Ad{ID: 2, Title: "title", Text: "text 2", AuthorID: 1, Published: true}
	adsList := []ads.Ad{expect1, expect2}

	title := expect1.Title
	s.app.On("GetListAdsByTitle", title).Return(adsList, nil)

	client := getTestClient(&s.app)

	response, err := client.listAdsByTitle(title)
	s.NoError(err)
	s.True(EqualAdsLists(response.Data, adsList))
}

func (s *AdServiceTestSuite) TestAdService_GetUser() {
	expect := &users.User{ID: 1, Nickname: "nickname", Email: "email"}
	s.app.On("GetUser", expect.ID).Return(expect, nil)

	client := getTestClient(&s.app)

	got, err := client.getUserById(expect.ID)
	s.NoError(err)
	s.True(EqualUsers(&got.Data, expect))
}

func (s *AdServiceTestSuite) TestAdService_GetUserNotFound() {
	expect := &users.User{ID: 1, Nickname: "nickname", Email: "email"}
	s.app.On("GetUser", expect.ID).Return(expect, app.IncorrectUserId)

	client := getTestClient(&s.app)

	_, err := client.getUserById(expect.ID)
	s.ErrorIs(err, ErrNotFound)
}

func (s *AdServiceTestSuite) TestAdService_DeleteUser() {
	expect := &users.User{ID: 1, Nickname: "nickname", Email: "email"}
	s.app.On("DeleteUser", expect.ID).Return(nil)

	client := getTestClient(&s.app)

	_, err := client.deleteUserById(expect.ID)
	s.NoError(err)
}

func (s *AdServiceTestSuite) TestAdService_DeleteUserIncorrectUserId() {
	expect := &users.User{ID: 1, Nickname: "nickname", Email: "email"}
	s.app.On("DeleteUser", expect.ID).Return(app.IncorrectUserId)

	client := getTestClient(&s.app)

	_, err := client.deleteUserById(expect.ID)
	s.ErrorIs(err, ErrNotFound)
}

func (s *AdServiceTestSuite) TestAdService_DeleteAd() {
	expect := &ads.Ad{ID: 1, Title: "", Text: "text 1", AuthorID: 1}
	s.app.On("DeleteAd", expect.ID, expect.AuthorID).Return(nil)

	client := getTestClient(&s.app)

	_, err := client.deleteAdById(expect.ID, expect.AuthorID)
	s.NoError(err)
}

func (s *AdServiceTestSuite) TestAdService_DeleteAdNotFound() {
	expect := &ads.Ad{ID: 1, Title: "", Text: "text 1", AuthorID: 1}
	s.app.On("DeleteAd", expect.ID, expect.AuthorID).Return(app.IncorrectAdId)

	client := getTestClient(&s.app)

	_, err := client.deleteAdById(expect.ID, expect.AuthorID)
	s.ErrorIs(err, ErrNotFound)
}

func (s *AdServiceTestSuite) TestAdService_DeleteAdIncorrectUserId() {
	expect := &ads.Ad{ID: 1, Title: "", Text: "text 1", AuthorID: 1}
	s.app.On("DeleteAd", expect.ID, expect.AuthorID).Return(app.IncorrectUserId)

	client := getTestClient(&s.app)

	_, err := client.deleteAdById(expect.ID, expect.AuthorID)
	s.ErrorIs(err, ErrForbidden)
}
