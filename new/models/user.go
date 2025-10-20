package models

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Contact   string `json:"contact"`
	Address   string `json:"address"`
	CreatedAt string `json:"created_at"`
}
