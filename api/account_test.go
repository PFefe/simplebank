package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/PFefe/simplebank/db/mock"
	db "github.com/PFefe/simplebank/db/sqlc"
	"github.com/PFefe/simplebank/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {
	account := RandomAccount()

	testCases := []struct {
		name         string
		accountID    int64
		buildStubs   func(store *mockdb.MockStore)
		checkRespose func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(
						gomock.Any(),
						gomock.Eq(account.ID),
					).
					Times(1).
					Return(
						account,
						nil,
					)
			},
			checkRespose: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(
					t,
					http.StatusOK,
					recorder.Code,
				)
				requireBodyMatchAccount(
					t,
					recorder.Body,
					account,
				)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(
						gomock.Any(),
						gomock.Eq(account.ID),
					).
					Times(1).
					Return(
						db.Account{},
						sql.ErrNoRows,
					)
			},
			checkRespose: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(
					t,
					http.StatusNotFound,
					recorder.Code,
				)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(
						gomock.Any(),
						gomock.Eq(account.ID),
					).
					Times(1).
					Return(
						db.Account{},
						sql.ErrConnDone,
					)
			},
			checkRespose: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(
					t,
					http.StatusInternalServerError,
					recorder.Code,
				)
			},
		},
		{
			name:      "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(
					gomock.Any(),
					gomock.Any(),
				).Times(0)
			},
			checkRespose: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(
					t,
					http.StatusBadRequest,
					recorder.Code,
				)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(
			tc.name,
			func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				server := NewServer(store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf(
					"/accounts/%d",
					tc.accountID,
				)
				request, err := http.NewRequest(
					http.MethodGet,
					url,
					nil,
				)
				require.NoError(
					t,
					err,
				)
				server.router.ServeHTTP(
					recorder,
					request,
				)

				tc.checkRespose(
					t,
					recorder,
				)
			},
		)
	}
}

func RandomAccount() db.Account {
	return db.Account{
		ID: util.RandomInt(
			1,
			1000,
		),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(
		t,
		err,
	)
	var gotAccount db.Account
	err = json.Unmarshal(
		data,
		&gotAccount,
	)
	require.NoError(
		t,
		err,
	)
	require.Equal(
		t,
		account,
		gotAccount,
	)
}
