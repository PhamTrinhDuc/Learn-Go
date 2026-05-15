package evaluation

// DatasetItem present for one item off dataset evaluation
type DatasetItem struct {
	Question           string   // Câu hỏi đầu vào
	GroundTruthAnswer  string   // Câu trả lời mẫu (nếu có)
	GroundTruthContext string   // Ngữ cảnh mẫu (để tính Retrieval metrics)
	RetrievedContexts  []string // Danh sách các văn bản tìm kiếm được
	GeneratedAnswer    string   // Câu trả lời do hệ thống RAG sinh ra
}

type DatasetResultEval struct {
	Item   DatasetItem
	Hit    int    // rank of document founded
	Status string // HIT or MISS
}

// Result hold all result after evaluation
type Result struct {
	Retrieval struct {
		HitRate      float64
		MRR          float64
		PrecisionAt1 float64
	}
	Generation struct {
		AvgFaithfulness    float64
		AvgAnswerRelevancy float64
	}
	TotalItems int
}
