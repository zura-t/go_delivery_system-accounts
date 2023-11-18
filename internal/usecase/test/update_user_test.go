package usecase_test

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/zura-t/go_delivery_system-accounts/internal/entity"
	"github.com/zura-t/go_delivery_system-accounts/pkg"
	mockdb "github.com/zura-t/go_delivery_system-accounts/pkg/db/mock"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
)

type UpdateUserRequest struct {
	Id   int64
	data *entity.UserUpdate
}

func Test_update_user(t *testing.T) {
	user, _ := randomUser(t)
	newName := pkg.RandomString(6)

	tests := []struct {
		name          string
		req           *UpdateUserRequest
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *entity.User, st int, err error)
	}{
		{
			name: "OK",
			req: &UpdateUserRequest{
				Id: user.ID,
				data: &entity.UserUpdate{
					Name: newName,
				},
			},
			buildStub: func(store *mockdb.MockStore) {
				req := db.UpdateUserParams{
					ID:   user.ID,
					Name: newName,
				}
				userUpdated := db.User{
					ID:        user.ID,
					Email:     user.Email,
					Name:      newName,
					Phone:     user.Phone,
					CreatedAt: user.CreatedAt,
				}

				store.EXPECT().UpdateUser(gomock.Any(), req).
					Times(1).
					Return(userUpdated, nil)
			},
			checkResponse: func(t *testing.T, res *entity.User, st int, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				updatedUser := res
				require.Equal(t, user.ID, updatedUser.Id)
				require.Equal(t, newName, updatedUser.Name)
				require.Equal(t, http.StatusOK, st)
			},
		},
		{
			name: "NotFound",
			req: &UpdateUserRequest{
				Id: user.ID,
				data: &entity.UserUpdate{
					Name: newName,
				},
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *entity.User, st int, err error) {
				require.Error(t, err)
				require.Empty(t, res)
				require.Equal(t, http.StatusNotFound, st)
			},
		},
		{
			name: "InternalError",
			req: &UpdateUserRequest{
				Id: user.ID,
				data: &entity.UserUpdate{
					Name: newName,
				},
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
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
			res, st, err := u.UpdateUser(tc.req.Id, tc.req.data)

			tc.checkResponse(t, res, st, err)
		})
	}
}
