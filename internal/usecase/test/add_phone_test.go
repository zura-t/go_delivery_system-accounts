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

type AddPhoneRequest struct {
	Id   int64
	data *entity.UserAddPhone
}

func Test_add_phone(t *testing.T) {
	user, _ := randomUser(t)
	newPhone := pkg.RandomPhone()

	tests := []struct {
		name          string
		req           *AddPhoneRequest
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res string, st int, err error)
	}{
		{
			name: "OK",
			req: &AddPhoneRequest{
				Id: user.ID,
				data: &entity.UserAddPhone{
					Phone: newPhone,
				},
			},
			buildStub: func(store *mockdb.MockStore) {
				req := db.AddPhoneParams{
					ID: user.ID,
					Phone: sql.NullString{
						String: newPhone,
						Valid:  true,
					},
				}
				store.EXPECT().AddPhone(gomock.Any(), req).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res string, st int, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.Equal(t, http.StatusOK, st)
			},
		},
		{
			name: "NotFound",
			req: &AddPhoneRequest{
				Id: user.ID,
				data: &entity.UserAddPhone{
					Phone: newPhone,
				},
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().AddPhone(gomock.Any(), gomock.Any()).Times(1).Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res string, st int, err error) {
				require.Error(t, err)
				require.Empty(t, res)
				require.Equal(t, http.StatusNotFound, st)
			},
		},
		{
			name: "InternalError",
			req: &AddPhoneRequest{
				Id: user.ID,
				data: &entity.UserAddPhone{
					Phone: newPhone,
				},
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().AddPhone(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res string, st int, err error) {
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
			res, st, err := u.AddPhone(tc.req.Id, tc.req.data)

			tc.checkResponse(t, res, st, err)
		})
	}
}
