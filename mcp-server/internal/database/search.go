// https://www.tigerdata.com/blog/you-dont-need-elasticsearch-bm25-is-now-in-postgres
package database

import (
	"context"
	"fmt"

	"github.com/pgvector/pgvector-go"
)

// VectorSearchResult holds params for vector search
type VectorSearchResult struct {
	Document KnowledgeBase
	Score    float64
}

// FullTextSearchResult holds params for full text search
type FullTextSearchResult struct {
	Document KnowledgeBase
}

// HybridSearchParams holds parameters for hybrid search
type HybridSearchParams struct {
	Query        string
	Embedding    []float32
	Limit        int
	BM25Weight   float64 // Weight for lexical search (0.0 to 1.0)
	VectorWeight float64 // Weight for semantic search (0.0 to 1.0)
	MinBM25Score float64 // Minimum BM25 score threshold
	MinVectorSim float64 // Minimum vector similarity threshold
}

// HybridSearchResult represents a result from hybrid search
type HybridSearchResult struct {
	Document      KnowledgeBase
	BM25Score     float64
	VectorScore   float64
	CombinedScore float64
}

// VectorSearch performs similarity search using pgvector
func (db *DB) VectorSearch(ctx context.Context, embedding []float32, limit int) ([]*VectorSearchResult, error) {
	queryString := `
		SELECT id, branch_id, title, content, metadata, embedding, created_at, updated_at,
		1 - (embedding <=>$1) AS similarity_score
		FROM knowledge_base 
		WHERE embedding IS NOT NULL AND is_active = TRUE
		ORDER BY embedding <=> $1
		LIMIT $2
	`
	vec := pgvector.NewVector(embedding)
	rows, err := db.pool.Query(ctx, queryString, vec, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to perform vector search: %w", err)
	}
	defer rows.Close()

	var results []*VectorSearchResult
	for rows.Next() {
		doc := &KnowledgeBase{}
		var score float64
		var dbEmbedding pgvector.Vector

		err := rows.Scan(
			&doc.ID,
			&doc.BranchId,
			&doc.Title,
			&doc.Content,
			&doc.Metadata,
			&dbEmbedding,
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&score,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan documents: %w", err)
		}
		doc.Embedding = dbEmbedding.Slice()
		results = append(results, &VectorSearchResult{
			Document: *doc,
			Score:    score,
		})
	}
	return results, nil
}

func (db *DB) FullTextSearch(ctx context.Context, query string, limit int) ([]*FullTextSearchResult, error) {
	queryString := `
		SELECT
			id, branch_id, title, content, metadata, created_at, updated_at,
			ts_rank_cd(
				to_tsvector('simple', f_unaccent(title || ' ' || content)),
				plainto_tsquery('simple', f_unaccent($1))
			) AS bm25_score
		FROM knowledge_base
		WHERE is_active = TRUE AND to_tsvector('simple', f_unaccent(title || ' ' || content)) @@ websearch_to_tsquery('simple', f_unaccent($1))
		ORDER BY bm25_score DESC
		LIMIT $2
	`

	rows, err := db.pool.Query(ctx, queryString, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to perform full text search: %w", err)
	}
	defer rows.Close()

	var results []*FullTextSearchResult
	for rows.Next() {
		doc := &KnowledgeBase{}
		var score float64

		err := rows.Scan(
			&doc.ID,
			&doc.BranchId,
			&doc.Title,
			&doc.Content,
			&doc.Metadata,
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&score,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan documents: %w", err)
		}
		results = append(results, &FullTextSearchResult{
			Document: *doc,
		})
	}
	return results, nil
}

// This implements a Reciprocal Rank Fusion (RRF) approach for combining results
func (db *DB) HybridSearch(ctx context.Context, params HybridSearchParams) ([]HybridSearchResult, error) {
	// Normalize weights if they don't sum to 1.0
	totalWeight := params.BM25Weight + params.VectorWeight
	if totalWeight == 0 {
		params.BM25Weight = 0.5
		params.VectorWeight = 0.5
		totalWeight = 1.0
	}
	bm25Weight := params.BM25Weight / totalWeight
	vectorWeight := params.VectorWeight / totalWeight

	if params.Limit <= 0 {
		params.Limit = 10
	}

	query := `
		WITH bm25_results AS (
			SELECT
				id, branch_id, title, content, metadata, embedding, created_at, updated_at,
				ts_rank_cd(
					to_tsvector('simple', f_unaccent(title || ' ' || content)),
					plainto_tsquery('simple', f_unaccent($1))
				) AS bm25_score,
				ROW_NUMBER() OVER (ORDER BY ts_rank_cd(
					to_tsvector('simple', f_unaccent(title || ' ' || content)),
					websearch_to_tsquery('simple', f_unaccent($1))
				) DESC) AS bm25_rank
			FROM knowledge_base
			WHERE is_active = TRUE AND to_tsvector('simple', f_unaccent(title || ' ' || content)) @@ websearch_to_tsquery('simple', f_unaccent($1))
		),
		vector_results AS (
			SELECT
				id, branch_id, title, content, metadata, embedding, created_at, updated_at,
				1 - (embedding <=> $2) AS vector_score,
				ROW_NUMBER() OVER (ORDER BY embedding <=> $2) AS vector_rank
			FROM knowledge_base
			WHERE is_active = TRUE AND embedding IS NOT NULL
		),
		combined AS (
			SELECT
				COALESCE(b.id, v.id) AS id,
				COALESCE(b.branch_id, v.branch_id) AS branch_id,
				COALESCE(b.title, v.title) AS title,
				COALESCE(b.content, v.content) AS content,
				COALESCE(b.metadata, v.metadata) AS metadata,
				COALESCE(b.embedding, v.embedding) AS embedding,
				COALESCE(b.created_at, v.created_at) AS created_at,
				COALESCE(b.updated_at, v.updated_at) AS updated_at,
				COALESCE(b.bm25_score, 0) AS bm25_score,
				COALESCE(v.vector_score, 0) AS vector_score,
				-- Reciprocal Rank Fusion score
				(
					COALESCE(1.0 / (60 + b.bm25_rank), 0) * $3 +
					COALESCE(1.0 / (60 + v.vector_rank), 0) * $4
				) AS combined_score
			FROM bm25_results b
			FULL OUTER JOIN vector_results v ON b.id = v.id
			WHERE
				COALESCE(b.bm25_score, 0) >= $5
				OR COALESCE(v.vector_score, 0) >= $6
		)
		SELECT
			id, branch_id, title, content, metadata, embedding,
			created_at, updated_at,
			bm25_score, vector_score, combined_score
		FROM combined
		ORDER BY combined_score DESC
		LIMIT $7
	`

	var embedding interface{}
	if params.Embedding != nil {
		embedding = pgvector.NewVector(params.Embedding)
	}

	rows, err := db.pool.Query(ctx, query,
		params.Query,
		embedding,
		bm25Weight,
		vectorWeight,
		params.MinBM25Score,
		params.MinVectorSim,
		params.Limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to perform hybrid search: %w", err)
	}
	defer rows.Close()

	var results []HybridSearchResult
	for rows.Next() {
		var doc KnowledgeBase
		var bm25Score, vectorScore, combinedScore float64
		var dbEmbedding *pgvector.Vector // Use pointer to handle NULL

		err := rows.Scan(
			&doc.ID,
			&doc.BranchId,
			&doc.Title,
			&doc.Content,
			&doc.Metadata,
			&dbEmbedding,
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&bm25Score,
			&vectorScore,
			&combinedScore,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hybrid search result: %w", err)
		}

		if dbEmbedding != nil && dbEmbedding.Slice() != nil {
			doc.Embedding = dbEmbedding.Slice()
		}

		results = append(results, HybridSearchResult{
			Document:      doc,
			BM25Score:     bm25Score,
			VectorScore:   vectorScore,
			CombinedScore: combinedScore,
		})
	}
	return results, nil
}
