package gapi

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/zura-t/go_delivery_system-accounts/internal/db/mock"
	db "github.com/zura-t/go_delivery_system-accounts/internal/db/sqlc"
	"github.com/zura-t/go_delivery_system-accounts/pb"
	"github.com/zura-t/go_delivery_system-accounts/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_update_user(t *testing.T) {
	user, _ := randomUser(t)
	newName := pkg.RandomString(6)

	tests := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Id:   user.ID,
				Name: newName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
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
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(userUpdated, nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				userUpdated := res.GetUser()
				require.Equal(t, newName, userUpdated.Name)
				require.Equal(t, user.Email, userUpdated.Email)
			},
		},
		{
			name: "NotFound",
			req: &pb.UpdateUserRequest{
				Id: user.ID,
				Name: newName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.UpdateUserRequest{
				Id:   user.ID,
				Name: newName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(1).Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range tests {
		tc := tests[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			tc.buildStubs(store)
			server := newTestServer(t, store)
			rsp, err := server.UpdateUser(context.Background(), tc.req)
			tc.checkResponse(t, rsp, err)
		})
	}
}
