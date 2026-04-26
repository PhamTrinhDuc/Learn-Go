package domain

type Conversation struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"user_id" db:"user_id"`
	Title     string `json:"title" db:"title"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type ConversationRepo interface {
	GetByID(id string) (*Conversation, error)
	Create(conversation *Conversation) error
	Update(conversation *Conversation) error
	Delete(id string) error
	ListByUserID(userID string) ([]*Conversation, error)
}
