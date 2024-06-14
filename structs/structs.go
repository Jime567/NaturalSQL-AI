package structs

type Cyclist struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	AddressID int    `json:"address_id"`
	BikeID    int    `json:"bike_id"`
}
