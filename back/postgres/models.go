package postgres

type (
	User struct {
		ID         uint64 `json:"id"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Admin      bool   `json:"admin"`
		Subscribed bool   `json:"subscribed"`
	}
)
