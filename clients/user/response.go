package clients

import "github.com/google/uuid"

type UserResponse struct {
	Code    string   `json:"code"`
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

type UserData struct {
	UUID        uuid.UUID `json:"uuid"`
	RegNumber   string    `json:"regNumber"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Photo       string    `json:"photo"`
	Role        string    `json:"role"`
	Library     string    `json:"library"`
}
