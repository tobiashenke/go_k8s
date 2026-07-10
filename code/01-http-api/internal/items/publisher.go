package items

import (
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

func (p *ItemPublisher) Publish(item Item) error {
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	err = p.conn.Publish("items.created", data)
	return err
}
