package rmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zura-t/go_delivery_system-accounts/cmd/api"
)

type EventPayload struct {
	Name string
	Data any
}

type CreateUserPayload struct {
	Name string
	Data api.CreateUserRequest
}

type Response struct {
	Data any
}

const (
	CreateUser string = "create_user"
	UpdateUser string = "update_user"
	DelateUser string = "delete_user"
)

func (consumer *Consumer) HandleRPC(ctx context.Context, event []byte) (*Response, error) {
	var payload EventPayload
	err := json.Unmarshal(event, &payload)
	if err != nil {
		return nil, err
	}

	switch payload.Name {
	// case CreateUser:
		// var payload CreateUserPayload
		// err := json.Unmarshal(event, &payload)
		// if err != nil {
		// 	return nil, err
		// }
		// user, err := consumer.server.CreateUser(ctx, payload.Data)
		// if err != nil {
		// 	return nil, err
		// }
		// res := &Response{Data: user}
		// return res, nil
	case UpdateUser:
		return nil, nil
	default:
		return nil, fmt.Errorf("no match events")
	}
}
