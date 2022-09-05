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

func TestAddEntry(t *testing.T) {
	user := CreateRandomUser()
	entry := createRandomEntry(user)

	entryResult := db.AddEntryTxResult {
		Entry: entry,
		User: user,
	}

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

	arg := db.AddEntryTxParams {
		Username: entry.Owner,
		Name: entry.Name,
		DueDate: entry.DueDate,
		Amount: entry.Amount,
	}

	testCases := []struct {
		name string
		reqArg CreateEntryParamsTest
		arg db.AddEntryTxParams
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
					AddEntryTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(entryResult, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchEntryResult(t, recorder.Body, entryResult)
			},
		},
		{
			name: "InvalidOwner",
			reqArg: reqArg,
			arg: arg,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					AddEntryTx(gomock.Any(), gomock.Any()).
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
					AddEntryTx(gomock.Any(), gomock.Any()).
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
					AddEntryTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.AddEntryTxResult{}, sql.ErrConnDone)
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
			server := newTestServer(t, store)
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
	user := CreateRandomUser()
	entry := createRandomEntry(user)

	type DeleteEntryParamsTest struct {
		Owner   string    `json:"owner"`
		ID int32 `json:"id"`
	}

	reqArg := DeleteEntryParamsTest {
		Owner: entry.Owner,
		ID: entry.ID,
	}

	arg := db.DeleteEntryTxParams {
		Username: user.Username,
		ID: entry.ID,
	}

	user.TotalExpenses = user.TotalExpenses - entry.Amount
	result := db.DeleteEntryTxResult {
		User: user,
	}

	testCases := []struct {
		name string
		reqArg DeleteEntryParamsTest
		arg db.DeleteEntryTxParams
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
					DeleteEntryTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(result, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchDeletedEntryResult(t, recorder.Body, result)
			},
		},
		{
			name: "BadId",
			reqArg: reqArg,
			arg: arg,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					DeleteEntryTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			reqArg: reqArg,
			arg: arg,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					DeleteEntryTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.DeleteEntryTxResult{}, sql.ErrConnDone)
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			if tc.name == "BadId" {
				tc.reqArg.ID = -1
			}

			url := fmt.Sprintf("/deleteEntry")
			body, err := json.Marshal(tc.reqArg)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetEntries(t *testing.T) {
	user := CreateRandomUser()

	var entries []db.Entry
	for i := 0; i < 3; i++ {
		entry := createRandomEntry(user)
		entry.Owner = "JhonDoe"
		entries = append(entries, entry)
	}

	owner := entries[0].Owner

	reqArg := struct{
		username string
	}{
		username: owner,
	}

	testCases := []struct {
		name string
		owner string
		reqArg struct{username string}
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			owner: owner,
			reqArg: reqArg,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					GetEntries(gomock.Any(), gomock.Eq(owner)).
					Times(1).
					Return(entries, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchEntries(t, recorder.Body, entries)
			},
		},
		{
			name: "BadRequest",
			owner: owner,
			reqArg: reqArg,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					GetEntries(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			owner: owner,
			reqArg: reqArg,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					GetEntries(gomock.Any(), gomock.Any()).
					Times(1).Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// check http status code
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			owner: owner,
			reqArg: reqArg,
			buildStubs: func(store *mockdb.MockStore) {
				// build stub
				store.EXPECT().
					GetEntries(gomock.Any(), gomock.Any()).
					Times(1).Return(nil, sql.ErrConnDone)
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
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			if  tc.name == "BadRequest" {
				tc.owner = "xyz"
			}

			url := fmt.Sprintf("/entries?username=%s", tc.owner)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func requireBodyMatchEntryResult(t *testing.T, body *bytes.Buffer, entryResult db.AddEntryTxResult) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotEntry db.AddEntryTxResult
	err = json.Unmarshal(data, &gotEntry)
	require.NoError(t, err)

	require.Equal(t, entryResult.Entry.Name, gotEntry.Entry.Name)
	require.Equal(t, entryResult.Entry.Owner, gotEntry.Entry.Owner)
	require.Equal(t, entryResult.Entry.DueDate, gotEntry.Entry.DueDate)
	require.Equal(t, entryResult.Entry.Amount, gotEntry.Entry.Amount)
}

func requireBodyMatchDeletedEntryResult(t *testing.T, body *bytes.Buffer, userResult db.DeleteEntryTxResult) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.DeleteEntryTxResult
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, userResult.User.Username, gotUser.User.Username)
	require.Equal(t, userResult.User.FullName, gotUser.User.FullName)
	require.Equal(t, userResult.User.Email, gotUser.User.Email)
	require.Equal(t, userResult.User.TotalExpenses, gotUser.User.TotalExpenses)
}

func requireBodyMatchEntries(t *testing.T, body *bytes.Buffer, entries []db.Entry) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotEntries []db.Entry
	err = json.Unmarshal(data, &gotEntries)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		require.Equal(t, entries[i].Name, gotEntries[i].Name)
		require.Equal(t, entries[i].Owner, gotEntries[i].Owner)
		require.Equal(t, entries[i].DueDate, gotEntries[i].DueDate)
		require.Equal(t, entries[i].Amount, gotEntries[i].Amount)
	}
}

func createRandomEntry(user db.User) db.Entry {
	date, _ :=  time.Parse(YYYYMMDD, "2022-12-11")

	return db.Entry {
		ID: 95,
		Owner: user.Username,
		Name: util.RandomString(6),
		DueDate: date,
		Amount: 5,
	}
}