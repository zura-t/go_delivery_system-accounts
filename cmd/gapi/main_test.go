package gapi

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zura-t/go_delivery_system-accounts/config"
	db "github.com/zura-t/go_delivery_system-accounts/pkg/db/sqlc"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := config.Config{}

	server, err := NewServer(store, config)
	require.NoError(t, err)

	return server
}
