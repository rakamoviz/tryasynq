package tasks

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TypeExpireOrder = "orderregistration:expire_order"
)

type ExpireOrderEmitter struct {
	myDbConn string
	AsynqEventEmitter
}

func NewExpireOrderTask(sagaId string) *asynq.Task {
	payload := map[string]interface{}{"saga_id": sagaId}
	return asynq.NewTask(TypeExpireOrder, payload)
}

func (p *ExpireOrderEmitter) ProcessTask(ctx context.Context, t *asynq.Task) error {
	sagaId, err := t.Payload.GetString("saga_id")
	if err != nil {
		return err
	}

	fmt.Printf("Send OrderExpired event to saga %s\n", sagaId)
	return nil
}

func NewExpireOrderEmitter() *ExpireOrderEmitter {
	return &ExpireOrderEmitter{
		myDbConn:          "helloDb",
		AsynqEventEmitter: AsynqEventEmitter{eventBus: "myBus"},
	}
}
