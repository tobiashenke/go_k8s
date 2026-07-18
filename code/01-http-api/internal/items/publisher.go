package items

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type ItemPublisher struct {
	conn *nats.Conn
}

func NewItemPublisher(natsUrl string) (*ItemPublisher, error) {
	conn, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}
	p := &ItemPublisher{
		conn: conn,
	}
	return p, nil
}

func (p *ItemPublisher) publish(data []byte, action string) error {
	err := p.conn.Publish(action, data)
	return err
}

func (p *ItemPublisher) PublishCreate(ctx context.Context, item Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	err = p.publish(data, "items.created")
	if err != nil {
		return err
	}
	return nil
}

func (p *ItemPublisher) PublishDelete(ctx context.Context, id string) error {
	data, err := json.Marshal(id)
	if err != nil {
		return err
	}
	err = p.publish(data, "items.deleted")
	if err != nil {
		return err
	}
	return nil
}
