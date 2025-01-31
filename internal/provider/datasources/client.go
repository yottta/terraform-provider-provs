package datasources

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"terraform-provider-hashicups/internal/client"
	"terraform-provider-hashicups/internal/model"
)

type coffeesClient struct {
	c client.BackendClient
}

func (c *coffeesClient) allCoffees() ([]model.Coffee, error) {
	all, err := c.c.ReadAll(coffeesResourceType)
	if err != nil {
		return nil, err
	}
	return readersToCoffee(all)
}

func (c *coffeesClient) createCoffee(cof model.Coffee) error {
	d, err := json.Marshal(cof)
	if err != nil {
		return err
	}
	_, err = c.c.CreateWithId(coffeesResourceType, strconv.Itoa(cof.ID), bytes.NewReader(d))
	return err
}

func readersToCoffee(in []io.Reader) ([]model.Coffee, error) {
	res := make([]model.Coffee, len(in))
	for i := 0; i < len(in); i++ {
		c, err := readerToCoffee(in[i])
		if err != nil {
			return nil, err
		}
		res[i] = c
	}
	return res, nil
}

func readerToCoffee(in io.Reader) (model.Coffee, error) {
	b, err := io.ReadAll(in)
	if err != nil {
		return model.Coffee{}, err
	}
	var c model.Coffee
	if err := json.Unmarshal(b, &c); err != nil {
		return model.Coffee{}, err
	}
	return c, nil
}
