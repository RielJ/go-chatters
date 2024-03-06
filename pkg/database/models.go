package database

// User struct
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Message struct
type Message struct {
	ID             int    `json:"id"`
	UserID         int    `json:"user_id"`
	Message        string `json:"message"`
	CreatedAt      string `json:"created_at"`
	ConversationID int    `json:"conversation_id"`
}

// Conversation struct
type Conversation struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ConversationUser struct
type ConversationUser struct {
	ID             int `json:"id"`
	ConversationID int `json:"conversation_id"`
	UserID         int `json:"user_id"`
}
