package telegram

// Session exist for managing the response of sended text by the status
type Session struct {
	ChatID   int `json:"id"`
	Username string
	Status   int
	IsLogged bool
}

// Status Type
const (
	STATUS_NONE              = 0
	STATUS_REGISTER_PASSWORD = 1
	STATUS_LOGIN_PASSWORD    = 2

	STATUS_ADD_IDEA    = 3
	STATUS_REMOVE_IDEA = 4
)
