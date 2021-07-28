package telegram

import "fmt"

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
	STATUS_ADD_IDEA          = 3
	STATUS_REMOVE_IDEA       = 4
	STATUS_ADD_MEDIA         = 5
)

//Obtain Slice Index by ChatID
func GetSessionID(chatID int) int {
	for i, v := range Sessions {
		if v.ChatID == chatID {
			return i
		}
	}
	return -1
}

// Destroy Session Slice by username
func DestroyOtherSession(username string) bool {
	for i, v := range Sessions {
		if v.Username == username {
			fmt.Printf("Founded Session in Slice, ID %d", i)
			removeSession(Sessions, i)
			return true
		}
	}
	return false
}

// Function to remove slice element
func removeSession(slice []Session, s int) []Session {
	return append(slice[:s], slice[s+1:]...)
}
