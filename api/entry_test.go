package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/LeandroEstevez/budgetAppAPI/db/mock"
	db "github.com/LeandroEstevez/budgetAppAPI/db/sqlc"
	"github.com/LeandroEstevez/budgetAppAPI/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	entry := createRandomEntry()

	type CreateEntryParamsTest struct {
		Owner   string    `json:"owner"`
		Name    string    `json:"name"`
		DueDate string `json:"due_date"`
		Amount  int64     `json:"amount"`
	}

	reqArg := CreateEntryParamsTest {
		Owner: entry.Owner,
		Name: entry.Name,
		DueDate: "2022-12-11",
		Amount: entry.Amount,
	}

	arg := db.CreateEntryParams {
		Owner: entry.Owner,
		Name: entry.Name,
		DueDate: entry.DueDate,
		Amount: entry.Amount,
	}

	testCases := []struct {
		name string
		reqArg CreateEntryParamsTest
		arg db.CreateEntryParams
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			reqArg: reqArg,
			arg: arg,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					CreateEntry(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(entry, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchEntry(t, recorder.Body, entry)
			},
		},
		{
			name: "InvalidOwner",
			reqArg: reqArg,
			arg: arg,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					CreateEntry(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidDate",
			reqArg: reqArg,
			arg: arg,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					CreateEntry(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InternalError",
			reqArg: reqArg,
			arg: arg,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					CreateEntry(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			if tc.name == "InvalidOwner" {
				tc.reqArg.Owner = "xyz"
			} else if tc.name == "InvalidDate" {
				tc.reqArg.DueDate = "2008-14-14"
			}

			url := fmt.Sprintf("/entry")
			body, err := json.Marshal(tc.reqArg)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	var id int32 = 10

	testCases := []struct {
		name string
		id int32
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id: id,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					DeleteEntry(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadId",
			id: -1,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					DeleteEntry(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id: id,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					DeleteEntry(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/deleteEntry/%d", tc.id)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchEntry(t *testing.T, body *bytes.Buffer, entry db.Entry) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotEntry db.Entry
	err = json.Unmarshal(data, &gotEntry)
	require.NoError(t, err)
	require.Equal(t, entry.Name, gotEntry.Name)
	require.Equal(t, entry.Owner, gotEntry.Owner)
	require.Equal(t, entry.DueDate, gotEntry.DueDate)
	require.Equal(t, entry.Amount, gotEntry.Amount)
}

func createRandomEntry() db.Entry {
	date, _ :=  time.Parse(YYYYMMDD, "2022-12-11")

	return db.Entry {
		ID: 95,
		Owner: util.RandomString(6),
		Name: util.RandomString(6),
		DueDate: date,
		Amount: 5,
	}
}