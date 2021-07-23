package telegram

import "fmt"

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
