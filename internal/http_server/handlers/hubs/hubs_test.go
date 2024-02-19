package hubs

import (
	"HestiaHome/internal/database"
	"HestiaHome/internal/lib/api/response"
	"HestiaHome/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the database.Storage interface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) CreateHub(ctx context.Context, hub *models.Hub) (int, error) {
	args := m.Called(ctx, hub)
	return args.Int(0), args.Error(1)
}

func (m *MockDB) GetHubByID(ctx context.Context, id int) (*models.DBHub, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.DBHub), args.Error(1)
}

func (m *MockDB) UpdateHub(ctx context.Context, hub *models.DBHub) error {
	args := m.Called(ctx, hub)
	return args.Error(0)
}

func (m *MockDB) DeleteHub(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDB) GetUserByID(ctx context.Context, id uuid.UUID) (*models.DBUser, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.DBUser), args.Error(1)
}

func (m *MockDB) GetHubsByUserID(ctx context.Context, id uuid.UUID) ([]*models.DBHub, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]*models.DBHub), args.Error(1)
}

func TestHubHandler_HubCtx(t *testing.T) {
	mockDB := new(MockDB)
	hubHandler := &HubHandler{log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})), db: mockDB}

	expectedHub := &models.DBHub{ID: 123, Name: "test", Description: "test description"}
	mockDB.On("GetHubByID", mock.Anything, 123).Return(expectedHub, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "123")
	req := httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.NotNil(t, r.Context().Value("hub"))
	})

	hubHandler.HubCtx(nextHandler).ServeHTTP(w, req)

	mockDB.AssertExpectations(t)
}

func TestHubHandler_CreateHub(t *testing.T) {
	testCases := []struct {
		name               string
		requestBody        HubRequest
		expectedStatusCode int
		mockCreateHub      struct {
			returnID  int
			returnErr error
		}
		mockGetUserByID struct {
			returnUser *models.DBUser
			returnErr  error
		}
		expectedResponse interface{}
		isError          bool
	}{
		{
			name: "Successful creation without owner_id",
			requestBody: HubRequest{
				Hub: &models.Hub{
					Name:        "hub_1",
					Description: "random description",
				},
			},
			expectedStatusCode: http.StatusCreated,
			mockCreateHub: struct {
				returnID  int
				returnErr error
			}{123, nil},
			mockGetUserByID: struct {
				returnUser *models.DBUser
				returnErr  error
			}{&models.DBUser{}, database.ErrUserNotExist},
			expectedResponse: HubResponse{
				Hub: &models.DBHub{
					ID:          123,
					OwnerID:     uuid.Nil,
					Name:        "hub_1",
					Description: "random description",
				},
			},
		},

		{
			name: "Successful creation with owner_id",
			requestBody: HubRequest{
				Hub: &models.Hub{
					OwnerID:     uuid.NameSpaceDNS,
					Name:        "hub_1",
					Description: "random description",
				},
			},
			expectedStatusCode: http.StatusCreated,
			mockCreateHub: struct {
				returnID  int
				returnErr error
			}{123, nil},
			mockGetUserByID: struct {
				returnUser *models.DBUser
				returnErr  error
			}{&models.DBUser{ID: uuid.NameSpaceDNS, Username: "some-username"}, nil},
			expectedResponse: HubResponse{
				Hub: &models.DBHub{
					ID:          123,
					OwnerID:     uuid.NameSpaceDNS,
					Name:        "hub_1",
					Description: "random description",
				},
				User: &models.DBUser{
					ID:       uuid.NameSpaceDNS,
					Username: "some-username",
				},
			},
		},

		{
			name: "Error creation without name",
			requestBody: HubRequest{
				Hub: &models.Hub{
					Description: "random description",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			mockCreateHub: struct {
				returnID  int
				returnErr error
			}{0, nil},
			mockGetUserByID: struct {
				returnUser *models.DBUser
				returnErr  error
			}{&models.DBUser{}, nil},
			expectedResponse: response.ErrResponse{
				StatusText: "Invalid request.",
				ErrorText:  errors.New("missing required Hub fields.").Error(),
			},
			isError: true,
		},

		{
			name: "Error storage mistake",
			requestBody: HubRequest{
				Hub: &models.Hub{
					Name:        "random name",
					Description: "random description",
				},
			},
			expectedStatusCode: http.StatusInternalServerError,
			mockCreateHub: struct {
				returnID  int
				returnErr error
			}{0, errors.New("")},
			mockGetUserByID: struct {
				returnUser *models.DBUser
				returnErr  error
			}{&models.DBUser{}, nil},
			expectedResponse: response.ErrResponse{
				StatusText: "Storage mistake",
			},
			isError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB := new(MockDB)
			hubHandler := &HubHandler{log: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})), db: mockDB}

			mockDB.On("CreateHub", mock.Anything, mock.Anything).Return(tc.mockCreateHub.returnID, tc.mockCreateHub.returnErr)
			mockDB.On("GetUserByID", mock.Anything, mock.Anything).Return(tc.mockGetUserByID.returnUser, tc.mockGetUserByID.returnErr)

			reqBody, err := json.Marshal(&tc.requestBody)
			if err != nil {
				t.Fatal("can't marshall request body to json")
			}
			req := httptest.NewRequest("POST", "/", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler := http.HandlerFunc(hubHandler.CreateHub)
			handler.ServeHTTP(w, req)

			resp := w.Result()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal("can't read body from response")
			}

			if !tc.isError {
				var r HubResponse
				err = json.Unmarshal(body, &r)
				if err != nil {
					t.Fatal("can't unmarshal body to HubResponse")
				}
				assert.Equal(t, r, tc.expectedResponse)
				assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
				mockDB.AssertExpectations(t)
			} else {
				var r response.ErrResponse
				err = json.Unmarshal(body, &r)
				if err != nil {
					t.Fatal("can't unmarshal body to HubResponse")
				}
				assert.Equal(t, r, tc.expectedResponse)
				assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)
			}
		})
	}

}

//func TestHubHandler_GetHub(t *testing.T) {
//	mockDB := new(MockDB)
//	hubHandler := &HubHandler{log: nil, db: mockDB}
//
//	expectedHub := &models.DBHub{ID: 123, Name: "test", Description: "test description"}
//	mockDB.On("GetHubByID", mock.Anything, 123).Return(expectedHub, nil)
//
//	rctx := chi.NewRouteContext()
//	rctx.URLParams.Add("id", strconv.Itoa(expectedHub.ID))
//	req := httptest.NewRequest("GET", "/", nil)
//	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
//	w := httptest.NewRecorder()
//
//	hubHandler.GetHub(w, req)
//
//	resp := w.Result()
//
//	assert.Equal(t, http.StatusOK, resp.StatusCode)
//	mockDB.AssertExpectations(t)
//}

//
//	func TestHubHandler_UpdateHub(t *testing.T) {
//		mockDB := new(MockDB)
//		hubHandler := &HubHandler{log: nil, db: mockDB}
//
//		expectedHub := &models.DBHub{ID: 123, Name: "test", Description: "test description"}
//		mockDB.On("GetHubByID", mock.Anything, 123).Return(expectedHub, nil)
//		mockDB.On("UpdateHub", mock.Anything, mock.Anything).Return(nil)
//
//		reqBody := strings.NewReader(`{"name":"updated test","description":"updated test description"}`)
//		req := httptest.NewRequest("PUT", "/", reqBody)
//		rctx := chi.NewRouteContext()
//		rctx.URLParams.Add("id", strconv.Itoa(expectedHub.ID))
//		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
//		w := httptest.NewRecorder()
//
//		hubHandler.UpdateHub(w, req)
//
//		resp := w.Result()
//
//		assert.Equal(t, http.StatusOK, resp.StatusCode)
//		mockDB.AssertExpectations(t)
//	}
//
//	func TestHubHandler_DeleteHub(t *testing.T) {
//		mockDB := new(MockDB)
//		hubHandler := &HubHandler{log: nil, db: mockDB}
//
//		expectedHub := &models.DBHub{ID: 123, Name: "test", Description: "test description"}
//		mockDB.On("GetHubByID", mock.Anything, 123).Return(expectedHub, nil)
//		mockDB.On("DeleteHub", mock.Anything, expectedHub.ID).Return(nil)
//
//		rctx := chi.NewRouteContext()
//		rctx.URLParams.Add("id", strconv.Itoa(expectedHub.ID))
//		req := httptest.NewRequest("DELETE", "/", nil)
//		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
//		w := httptest.NewRecorder()
//
//		hubHandler.DeleteHub(w, req)
//
//		resp := w.Result()
//
//		assert.Equal(t, http.StatusOK, resp.StatusCode)
//		mockDB.AssertExpectations(t)
//	}
//
//	func TestHubRequest_Bind(t *testing.T) {
//		validRequest := &HubRequest{Hub: &models.Hub{Name: "Test", Description: "Test Description"}}
//		err := validRequest.Bind(nil)
//		assert.NoError(t, err)
//
//		invalidRequest := &HubRequest{Hub: nil}
//		err = invalidRequest.Bind(nil)
//		assert.Error(t, err)
//		assert.Equal(t, "missing required Hub fields.", err.Error())
//	}
