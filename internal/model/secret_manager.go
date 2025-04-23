package model

type SecretManager struct {
	ID      string            `json:"id,omitempty"`
	Secrets map[string]string `json:"secrets,omitempty"`
	Name    string            `json:"name"`
}

func (o *SecretManager) GetID() string {
	return o.ID
}

func (o *SecretManager) SetID(id string) {
	o.ID = id
}
