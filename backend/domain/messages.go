package domain

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

type Message struct {
	ID             string `json:"id" db:"id"`
	ConversationID string `json:"conversation_id" db:"conversation_id"`
	Role           Role   `json:"role" db:"role"`
	Content        string `json:"content" db:"content"`
	CreatedAt      string `json:"created_at" db:"created_at"`
	UpdatedAt      string `json:"updated_at" db:"updated_at"`
}

type MessageRepo interface {
	GetByID(id string) (*Message, error)
	Create(message *Message) error
	Update(message *Message) error
	Delete(id string) error
	ListByConversationID(conversationID string) ([]*Message, error)
}
