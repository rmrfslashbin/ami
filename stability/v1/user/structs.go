package user

type Response struct {
	User    *ResponseUser        `json:"user,omitempty"`
	Error   *ResponseUserError   `json:"error,omitempty"`
	Credits *ResponseUserBalance `json:"credits,omitempty"`
}

type ResponseUserBalance struct {
	Credits float64 `json:"credits"`
}

type ResponseUser struct {
	Id             string            `json:"id"`
	Email          string            `json:"email"`
	ProfilePicture string            `json:"profile_picture"`
	Organizations  []ResponseUserOrg `json:"organizations"`
}

type ResponseUserOrg struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	IsDefault bool   `json:"is_default"`
}

type ResponseUserError struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type BalanceInput struct {
	Organization           *string `json:"organization,omitempty"`
	StabilityClientID      *string `json:"stability_client_id,omitempty"`
	StabilityClientVersion *string `json:"stability_client_version,omitempty"`
}
