package model

type Role string

var (
	Admin    Role = "admin"
	Analyser Role = "analyser"
	NoRole   Role = "no role"
)

type User struct {
	Role     Role   `json:"role"`
	Username string `json:"userName"`
}
