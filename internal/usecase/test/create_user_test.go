package usecase_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/zura-t/go_delivery_system-accounts/internal/entity"
	"github.com/zura-t/go_delivery_system-accounts/pkg"

	mockdb "github.com/zura-t/go_delivery_system-accounts/pkg/db/mock"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := pkg.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = pkg.RandomString(6)
	hashedPassword, err := pkg.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Email:          pkg.RandomEmail(),
		HashedPassword: hashedPassword,
		Name:           pkg.RandomString(6),
	}
	return user, password
}

func Test_create_user(t *testing.T) {
	user, password := randomUser(t)

	tests := []struct {
		name          string
		req           *entity.UserRegister
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *entity.User, st int, err error)
	}{
		{
			name: "OK",
			req: &entity.UserRegister{
				Email:    user.Email,
				Password: password,
				Name:     user.Name,
			},
			buildStub: func(store *mockdb.MockStore) {
				request := db.CreateUserParams{
					Email:          user.Email,
					HashedPassword: user.HashedPassword,
					Name:           user.Name,
				}
				userCreated := db.User{
					ID:        user.ID,
					Email:     user.Email,
					Name:      user.Name,
					Phone:     user.Phone,
					CreatedAt: user.CreatedAt,
				}
				store.EXPECT().GetUserByEmail(gomock.Any(), request.Email).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)

				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(request, password)).
					Times(1).
					Return(userCreated, nil)
			},
			checkResponse: func(t *testing.T, res *entity.User, st int, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUser := res
				require.Equal(t, user.Email, createdUser.Email)
				require.Equal(t, user.Name, createdUser.Name)
				require.Equal(t, http.StatusOK, st)
			},
		},
		{
			name: "InternalError",
			req: &entity.UserRegister{
				Email:    user.Email,
				Password: password,
				Name:     user.Name,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *entity.User, st int, err error) {
				require.Error(t, err)
				require.Equal(t, http.StatusInternalServerError, st)
			},
		},
		{
			name: "DuplicateEmail",
			req: &entity.UserRegister{
				Email:    user.Email,
				Password: password,
				Name:     user.Name,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(user, nil)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *entity.User, st int, err error) {
				require.Error(t, err)
				require.Equal(t, http.StatusBadRequest, st)
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
			res, st, err := u.CreateUser(tc.req)

			tc.checkResponse(t, res, st, err)
		})
	}
}
