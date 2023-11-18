package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zura-t/go_delivery_system-accounts/config"
	"github.com/zura-t/go_delivery_system-accounts/internal/usecase"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
	"github.com/zura-t/go_delivery_system-accounts/token"
)

func newTestServer(t *testing.T, store db.Store) *usecase.UserUseCase {
	config := &config.Config{}
	tokenMaker, err := token.NewJwtMaker("test")
	require.NoError(t, err)
	
	server := usecase.New(store, config, tokenMaker)

	return server
}
