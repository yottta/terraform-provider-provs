package model

type Order struct {
	ID    string      `json:"id,omitempty"`
	Items []OrderItem `json:"items,omitempty"`
}

func (o *Order) GetID() string {
	return o.ID
}

func (o *Order) SetID(id string) {
	o.ID = id
}

type OrderItem struct {
	Coffee   Coffee `json:"coffee"`
	Quantity int    `json:"quantity"`
}
