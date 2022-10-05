package crud

type JSONorder struct {
	Name     string
	UserName string
	Country  string
	Cart     []string
}

type Clients struct {
	Username string
	Password string
	Country  string
}

type Products struct {
	Code     int     `json:"code"`
	Name     string  `json:"name"`
	Cat      string  `json:"cat"`
	ImageUrl string  `json:"imageUrl"`
	Price    float64 `json:"price"`
}

type Orders struct {
	ID         int
	ProductCod int
	ClientDNI  int
}
