package filesystem

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"terraform-provider-provs/internal/client"

	"github.com/google/uuid"
	"github.com/spf13/afero"
	_ "github.com/spf13/afero"
)

// fsClient implements BackendClient to provide a local storage solution for resources management
type fsClient struct {
	fs afero.Fs
}

func NewFsClient(basePath string) (client.BackendClient, error) {
	if !path.IsAbs(basePath) {
		return nil, fmt.Errorf("only absolute paths allowed")
	}
	if err := os.MkdirAll(basePath, 0744); err != nil {
		return nil, err
	}
	return &fsClient{
		fs: afero.NewBasePathFs(
			afero.NewOsFs(),
			basePath,
		),
	}, nil
}

func (c *fsClient) Create(resType string, body io.Reader) (string, error) {
	return c.CreateWithId(resType, uuid.NewString(), body)
}

func (c *fsClient) Read(resType string, resId string) (io.Reader, error) {
	fileName := fmt.Sprintf("%s%s%s", resType, string(os.PathSeparator), resId)
	f, err := c.fs.OpenFile(fileName, os.O_RDONLY, os.FileMode(0644))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(content), nil
}

func (c *fsClient) Destroy(resType string, resId string) error {
	fileName := fmt.Sprintf("%s%s%s", resType, string(os.PathSeparator), resId)
	return c.fs.Remove(fileName)
}

func (c *fsClient) Update(resType string, resId string, newContent io.Reader) error {
	if err := c.Destroy(resType, resId); err != nil {
		return err
	}
	if _, err := c.CreateWithId(resType, resId, newContent); err != nil {
		return err
	}
	return nil
}

func (c *fsClient) CreateWithId(resType string, resId string, body io.Reader) (string, error) {
	if err := c.fs.Mkdir(resType, 0744); err != nil && !errors.Is(err, os.ErrExist) {
		return "", err
	}
	fileName := fmt.Sprintf("%s%s%s", resType, string(os.PathSeparator), resId)
	f, err := c.fs.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	b := make([]byte, 1024)
	for {
		n, err := body.Read(b)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := f.Write(b[:n]); err != nil {
			return "", err
		}

	}
	return resId, nil

}

func (c *fsClient) ReadAll(resType string) ([]io.Reader, error) {
	var res []io.Reader
	dir, err := c.fs.Open(resType)
	if err != nil {
		return nil, err
	}
	defer func() { _ = dir.Close() }()
	fileInfo, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, fi := range fileInfo {
		content, err := c.Read(resType, fi.Name())
		if err != nil {
			return nil, err
		}
		res = append(res, content)
	}
	return res, nil
}
