package structs

type Cyclist struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	AddressID   *int   `json:"address_id"`
	BikeID      *int   `json:"bike_id"`
	SkillLevel  string `json:"skill_level"`
}

type Bike struct {
	ID           int    `json:"id"`
	Nickname     string `json:"nickname"`
	SerialNumber string `json:"serial_number"`
	Year         string `json:"year"`
	Model        string `json:"model"`
	Make         string `json:"make"`
	Mileage      int    `json:"mileage"`
}

type Address struct {
	ID     int    `json:"id"`
	Street string `json:"street"`
	Zip    string `json:"zip"`
	State  string `json:"state"`
}
