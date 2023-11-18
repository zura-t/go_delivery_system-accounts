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

func Test_login_user(t *testing.T) {
	user, password := randomUser(t)

	tests := []struct {
		name          string
		req           *entity.UserLogin
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *entity.UserLoginResponse, st int, err error)
	}{
		{
			name: "OK",
			req: &entity.UserLogin{
				Email:    user.Email,
				Password: password,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, res *entity.UserLoginResponse, st int, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, user.ID, res.User.Id)
				require.Equal(t, http.StatusOK, st)
			},
		},
		{
			name: "NotFound",
			req: &entity.UserLogin{
				Email:    "email",
				Password: password,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *entity.UserLoginResponse, st int, err error) {
				require.Error(t, err)
				require.Empty(t, res)
				require.Equal(t, http.StatusNotFound, st)
			},
		},
		{
			name: "IncorrectPassword",
			req: &entity.UserLogin{
				Email:    user.Email,
				Password: password + "new",
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, res *entity.UserLoginResponse, st int, err error) {
				require.Error(t, err)
				require.Empty(t, res)
				require.Equal(t, http.StatusBadRequest, st)
			},
		},
		{
			name: "InternalError",
			req: &entity.UserLogin{
				Email:    user.Email,
				Password: password,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *entity.UserLoginResponse, st int, err error) {
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
			res, st, err := u.LoginUser(tc.req)

			tc.checkResponse(t, res, st, err)
		})
	}
}
