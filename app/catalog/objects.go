package catalog

type Response struct {
	Products []Product       `json:"products"`
	Filter   *ResponseFilter `json:"page,omitempty"`
}

type ResponseFilter struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type Product struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category string    `json:"category"`
	Variants []Variant `json:"variants"`
}

type Variant struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
