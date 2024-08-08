package models

type UserDatabase struct {
	Username string `json:"username"`
	Hash     string `json:"hash"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserChangePassword struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type CreateUserResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
