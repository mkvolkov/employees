package model

type Employee struct {
	ID      int    `db:"id"`
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Gender  string `json:"gender" db:"gender"`
	Age     int    `json:"age" db:"age"`
	Email   string `json:"email" db:"email"`
	Address string `json:"address" db:"address"`
	Vdays   int    `db:"vdays"`
}

type DeleteID struct {
	ID int `json:"id"`
}
