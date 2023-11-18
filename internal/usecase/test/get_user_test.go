package usecase_test

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/zura-t/go_delivery_system-accounts/internal/entity"
	mockdb "github.com/zura-t/go_delivery_system-accounts/pkg/db/mock"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
)

func Test_get_user(t *testing.T) {
	user, _ := randomUser(t)

	tests := []struct {
		name          string
		req           int64
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *entity.User, st int, err error)
	}{
		{
			name: "OK",
			req:  user.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), user.ID).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, res *entity.User, st int, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, user.ID, res.Id)
				require.Equal(t, user.Email, res.Email)
				require.Equal(t, http.StatusOK, st)
			},
		},
		{
			name: "NotFound",
			req:  user.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *entity.User, st int, err error) {
				require.Error(t, err)
				require.Empty(t, res)
				require.Equal(t, http.StatusNotFound, st)
			},
		},
		{
			name: "InternalError",
			req:  user.ID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *entity.User, st int, err error) {
				require.Error(t, err)
				require.Empty(t, res)
				require.Equal(t, http.StatusInternalServerError, st)
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			u := newTestServer(t, store)

			tc.buildStub(store)
			res, st, err := u.GetUser(tc.req)

			tc.checkResponse(t, res, st, err)
		})
	}
}
