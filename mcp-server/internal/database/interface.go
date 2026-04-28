package database

import (
	"context"

	"github.com/google/uuid"
)

// Store defines the interface for database operations
// This interface enables testing with mocks
type Store interface {
	InsertDocument(ctx context.Context, doc *KnowledgeBase) error
	InsertDocuments(ctx context.Context, docs []*KnowledgeBase) error
	SearchDocuments(ctx context.Context, query string, limit int) ([]*KnowledgeBase, error)
	GetDocument(ctx context.Context, docID uuid.UUID) (*KnowledgeBase, error)
	ListDocuments(ctx context.Context, limit int, offset int) ([]*KnowledgeBase, error)
	UpdateDocument(ctx context.Context, doc *KnowledgeBase) error
	DeleteDocumentByID(ctx context.Context, docID uuid.UUID) error
}

// Ensure DB implements Store interface
var _ Store = (*DB)(nil)
