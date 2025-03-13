package model

type UserDetails struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Province string `json:"province"`
	ZipCode  string `json:"zip_code"`

	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func NewUserDetails(userID int, phone, address, city, province, zipCode string) *UserDetails {
	return &UserDetails{
		UserID:   userID,
		Phone:    phone,
		Address:  address,
		City:     city,
		Province: province,
		ZipCode:  zipCode,
	}
}
