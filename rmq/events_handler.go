package rmq

import (
	"context"
	"encoding/json"
)

type EventPayload struct {
	Name string
	Data any
}

type Response struct {
	Data any
}

const (
	CreateUser string = "create_user"
	UpdateUser string = "update_user"
	DelateUser string = "delete_user"
)

func (consumer *Consumer) HandleRPC(ctx context.Context, event []byte) {
	var payload EventPayload
	err := json.Unmarshal(event, &payload)
	if err != nil {
		// return nil, err
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
		// return nil, nil
	default:
		// return nil, fmt.Errorf("no match events")
	}
}
