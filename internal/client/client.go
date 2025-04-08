package client

import (
	"bytes"
	"encoding/json"
	"io"
	"terraform-provider-provs/internal/model"

	"github.com/google/uuid"
)

type Client[T model.IDer] interface {
	GetAll() ([]T, error)
	GetByID(id string) (T, error)
	// Create will use the obj.GetID as identifier is specified. Otherwise, will generate one and will call obj.SetID with it
	Create(obj T) (T, error)
	Update(obj T) error
	Delete(id string) error
}

type client[T model.IDer] struct {
	resType string
	c       BackendClient
}

func NewClient[T model.IDer](backend BackendClient, resType string) Client[T] {
	return &client[T]{
		c:       backend,
		resType: resType,
	}
}

func (c *client[T]) GetAll() ([]T, error) {
	all, err := c.c.ReadAll(c.resType)
	if err != nil {
		return nil, err
	}
	return c.readersToObj(all)
}

func (c *client[T]) GetByID(id string) (T, error) {
	var out T
	dat, err := c.c.Read(c.resType, id)
	if err != nil {
		return out, err
	}
	return c.readerToObj(dat)
}

func (c *client[T]) Create(obj T) (T, error) {
	if obj.GetID() == "" {
		obj.SetID(uuid.NewString())
	}
	read, err := c.objToReader(obj)
	if err != nil {
		return obj, err
	}
	if err := c.c.CreateWithId(c.resType, obj.GetID(), read); err != nil {
		return obj, err
	}
	return obj, nil
}

func (c *client[T]) Update(obj T) error {
	read, err := c.objToReader(obj)
	if err != nil {
		return err
	}
	return c.c.Update(c.resType, obj.GetID(), read)
}

func (c *client[T]) Delete(id string) error {
	return c.c.Destroy(c.resType, id)
}

func (c *client[T]) readersToObj(in []io.Reader) ([]T, error) {
	res := make([]T, len(in))
	for i := 0; i < len(in); i++ {
		d, err := c.readerToObj(in[i])
		if err != nil {
			return nil, err
		}
		res[i] = d
	}
	return res, nil
}

func (c *client[T]) readerToObj(in io.Reader) (T, error) {
	var out T
	b, err := io.ReadAll(in)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(b, &out); err != nil {
		return out, err
	}
	return out, nil
}

func (c *client[T]) objToReader(obj T) (io.Reader, error) {
	d, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(d), nil
}
