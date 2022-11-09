package account

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}
type DeleteRequest struct {
	AccountID int64
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	AccessExp    int64  `json:"access_exp"`
	RefreshToken string `json:"refresh_token"`
	RefreshExp   int64  `json:"refresh_exp"`
}

type GetRequest struct {
	AccountID int64
}

type GetResponse struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type UpdateRequest struct {
	AccountID int64  `json:"-"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Password  string `json:"password"`
}

type UpdateResponse struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
