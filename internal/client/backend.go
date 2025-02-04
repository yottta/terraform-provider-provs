package client

import "io"

type BackendClient interface {
	Create(resType string, body io.Reader) (string, error)
	CreateWithId(resType string, id string, body io.Reader) (string, error)
	Read(resType string, resId string) (io.Reader, error)
	ReadAll(resType string) ([]io.Reader, error)
	Destroy(resType string, resId string) error
	Update(resType string, resId string, newContent io.Reader) error
}
