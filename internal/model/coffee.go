package model

type Coffee struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Teaser      string       `json:"teaser"`
	Description string       `json:"description"`
	Price       float64      `json:"price"`
	Image       string       `json:"image"`
	Ingredient  []Ingredient `json:"ingredients"`
}

func (c *Coffee) GetID() string {
	return c.ID
}

func (c *Coffee) SetID(id string) {
	c.ID = id
}

type Ingredient struct {
	ID       string `json:"ingredient_id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}

func (c *Ingredient) GetID() string {
	return c.ID
}

func (c *Ingredient) SetID(id string) {
	c.ID = id
}
