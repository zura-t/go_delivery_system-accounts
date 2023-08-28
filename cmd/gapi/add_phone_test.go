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
	"google.golang.org/protobuf/types/known/emptypb"
)

func Test_add_phone(t *testing.T) {
	user, _ := randomUser(t)
	newPhone := pkg.RandomPhone()

	tests := []struct {
		name          string
		req           *pb.AddPhoneRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *emptypb.Empty, err error)
	}{
		{
			name: "OK",
			req: &pb.AddPhoneRequest{
				Id:    user.ID,
				Phone: newPhone,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.AddPhoneParams{
					ID: user.ID,
					Phone: sql.NullString{
						String: newPhone,
						Valid:  true,
					},
				}
				store.EXPECT().
					AddPhone(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *emptypb.Empty, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "InternalError",
			req: &pb.AddPhoneRequest{
				Id:    user.ID,
				Phone: newPhone,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					AddPhone(gomock.Any(), gomock.Any()).
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
			name: "NotFound",
			req: &pb.AddPhoneRequest{
				Id:    user.ID,
				Phone: newPhone,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					AddPhone(gomock.Any(), gomock.Any()).
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
			rsp, err := server.AddPhone(context.Background(), tc.req)
			tc.checkResponse(t, rsp, err)
		})
	}
}
