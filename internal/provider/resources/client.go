package resources

import (
	"bytes"
	"encoding/json"
	"io"
	"terraform-provider-hashicups/internal/client"
	"terraform-provider-hashicups/internal/model"

	"github.com/google/uuid"
)

type orderClient struct {
	c client.BackendClient
}

func (c *orderClient) get(id string) (model.Order, error) {
	all, err := c.c.Read(orderResourceType, id)
	if err != nil {
		return model.Order{}, err
	}
	return readerToOrder(all)
}

func (c *orderClient) update(id string, items []model.OrderItem) error {
	o := model.Order{
		ID:    id,
		Items: items,
	}
	read, err := orderToReader(o)
	if err != nil {
		return err
	}
	return c.c.Update(orderResourceType, o.ID, read)
}

func (c *orderClient) create(orderItems []model.OrderItem) (model.Order, error) {
	o := model.Order{
		ID:    uuid.NewString(),
		Items: orderItems,
	}
	read, err := orderToReader(o)
	if err != nil {
		return model.Order{}, err
	}
	_, err = c.c.CreateWithId(orderResourceType, o.ID, read)
	if err != nil {
		return model.Order{}, err
	}
	return o, nil
}

func (c *orderClient) delete(id string) error {
	return c.c.Destroy(orderResourceType, id)
}

func readerToOrder(in io.Reader) (model.Order, error) {
	b, err := io.ReadAll(in)
	if err != nil {
		return model.Order{}, err
	}
	var c model.Order
	if err := json.Unmarshal(b, &c); err != nil {
		return model.Order{}, err
	}
	return c, nil
}

func orderToReader(o model.Order) (io.Reader, error) {
	d, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(d), nil
}
