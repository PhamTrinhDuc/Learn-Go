package evaluation

// DatasetItem present for one item off dataset evaluation
type DatasetItem struct {
	Question           string   // Câu hỏi đầu vào
	GroundTruthAnswer  string   // Câu trả lời mẫu (nếu có)
	GroundTruthContext string   // Ngữ cảnh mẫu (để tính Retrieval metrics)
	RetrievedContexts  []string // Danh sách các văn bản tìm kiếm được
	GeneratedAnswer    string   // Câu trả lời do hệ thống RAG sinh ra
}

// DatasetResultEval hold result for each item evaluation
type DatasetResultEval struct {
	Item                 DatasetItem
	Hit                  int     // rank of document founded
	Status               string  // HIT or MISS
	TimeSearch           float64 // time search
	ScoreAnswerRelevancy float64 // score for answer relevancy
	ScoreFaithfulness    float64 // score for faithfulness
}

// RetrievalMetrics hold metrics for evaluation retrieval
type RetrievalMetrics struct {
	HitRate                 float64 // percentage of questions that retrieved at least one correct document
	MRR                     float64 // mean reciprocal rank
	PrecisionAt1            float64 // percentage of questions that retrieved the correct document as the first result
	AvgTimeSearch           float64 // average time search in seconds
	AvgScoreAnswerRelevancy float64 // average score for answer relevancy
	AvgScoreFaithfulness    float64 // average score for faithfulness
}

// Result hold all result after evaluation
type Result struct {
	Retrieval struct {
		HitRate       float64
		MRR           float64
		PrecisionAt1  float64
		AvgTimeSearch float64
	}
	Generation struct {
		AvgFaithfulness    float64
		AvgAnswerRelevancy float64
	}
	TotalItems    int
	AvgTimeSearch float64
}
