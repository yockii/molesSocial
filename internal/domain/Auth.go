package domain

type SignUpRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" form:"email"`
	Reason   string `json:"reason" form:"reason"`
}

type SignInRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
