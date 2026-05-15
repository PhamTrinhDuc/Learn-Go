package evaluation

import "log"

// RetrievalResult hold metrics for evaluation
type RetrievalResult struct {
	HitRate      float64
	MRR          float64
	PrecisionAt1 float64
}

// EvaluateRetrieval log detail process evaluation and metrics
func EvaluateRetrieval(items []DatasetItem) ([]DatasetResultEval, RetrievalResult) {
	var hits, p1 int
	var mrrSum float64
	total := len(items)

	if total == 0 {
		log.Println("no found items to evaluation!")
		return []DatasetResultEval{}, RetrievalResult{}
	}

	datasetResultEvals := []DatasetResultEval{}

	for _, item := range items {
		rank := 0
		for i, ctx := range item.RetrievedContexts {
			if isMatch(ctx, item.GroundTruthContext) {
				rank = i + 1
				break
			}
		}

		if rank > 0 {
			hits++
			mrrSum += 1.0 / float64(rank)
			if rank == 1 {
				p1++
			}
		}

		datasetResultEval := DatasetResultEval{
			Item:   item,
			Hit:    rank,
			Status: GetStatus(rank),
		}
		datasetResultEvals = append(datasetResultEvals, datasetResultEval)
	}

	resultsMetrics := RetrievalResult{
		HitRate:      float64(hits) / float64(total) * 100,
		MRR:          mrrSum / float64(total),
		PrecisionAt1: float64(p1) / float64(total) * 100,
	}

	return datasetResultEvals, resultsMetrics
}

// isMatch check len(string) equal
func isMatch(retrieved, expected string) bool {
	if len(retrieved) == 0 || len(expected) == 0 {
		return false
	}
	return (len(retrieved) >= len(expected) && retrieved[:len(expected)] == expected) ||
		(len(expected) >= len(retrieved) && expected[:len(retrieved)] == retrieved)
}

// GetStatus return HIT or MISS if document retrieved found in documents groundtruth
func GetStatus(rank int) string {
	if rank > 0 {
		return "HIT"
	}
	return "MISS"
}
