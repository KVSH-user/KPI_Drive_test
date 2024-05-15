package stan

import (
	"KPI_Drive_test/internal/entity"
	"KPI_Drive_test/internal/factSender"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log/slog"
)

type Client struct {
	Sc stan.Conn
}

func NewClient(clusterID, clientID, natsURL string) (*Client, error) {
	const op = "handlers.stan.NewClient"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Client{Sc: sc}, nil
}

func (c *Client) Subscribe(subject string, cb stan.MsgHandler) (stan.Subscription, error) {
	return c.Sc.Subscribe(subject, cb, stan.DurableName("my-durable"))
}

func (c *Client) Close() error {
	if c.Sc != nil {
		c.Sc.Close()
	}
	return nil
}

func FactMessage(log *slog.Logger, m *stan.Msg) error {
	const op = "handlers.stan.FactMessage"

	var fact entity.Fact
	if err := json.Unmarshal(m.Data, &fact); err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	// отправляем факт
	if err := factSender.SendFact(fact); err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	log.Info("fact send successfully")

	return nil
}
