package gapi

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/zura-t/go_delivery_system-accounts/internal/db/mock"
	"github.com/zura-t/go_delivery_system-accounts/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_delete_user(t *testing.T) {
	user, _ := randomUser(t)

	tests := []struct {
		name          string
		userID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *emptypb.Empty, err error)
	}{
		{
			name:   "OK",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *emptypb.Empty, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:   "InternalError",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *emptypb.Empty, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name:   "NotFound",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *emptypb.Empty, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
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
			rsp, err := server.DeleteUser(context.Background(), &pb.UserId{Id: tc.userID})
			tc.checkResponse(t, rsp, err)
		})
	}
}
