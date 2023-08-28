package gapi

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zura-t/go_delivery_system-accounts/internal"
	db "github.com/zura-t/go_delivery_system-accounts/internal/db/sqlc"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := internal.Config{}

	server, err := NewServer(store, config)
	require.NoError(t, err)

	return server
}
