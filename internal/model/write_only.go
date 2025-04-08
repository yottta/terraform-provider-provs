package model

type WriteOnly struct {
	ID   string `json:"id,omitempty"`
	Attr string `json:"attr,omitempty"`
}

func (o *WriteOnly) GetID() string {
	return o.ID
}

func (o *WriteOnly) SetID(id string) {
	o.ID = id
}

func (o *WriteOnly) GetAttr() string {
	return o.Attr
}

func (o *WriteOnly) SetAttr(attr string) {
	o.Attr = attr
}
