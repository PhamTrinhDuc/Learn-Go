package evaluation

import (
	"context"
	"log"
	"mcp-server/internal/llm"
)

// EvaluateRetrieval log detail process evaluation and metrics
func EvaluateRetrieval(items []DatasetItem, client llm.LLMModel) ([]DatasetResultEval, RetrievalMetrics) {
	// Init metrics variables
	var hits, p1 int
	var mrrSum float64
	total := len(items)

	if total == 0 {
		log.Println("no found items to evaluation!")
		return []DatasetResultEval{}, RetrievalMetrics{}
	}

	// 2. Format datasets and int Judge
	datasetResultEvals := []DatasetResultEval{}
	judge := NewJudge(client)

	for _, item := range items {
		// Caculate rank, hit for metrics
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
		var scoreFaithFulness, scoreAnswerRelevancy float64
		var err error
		if judge != nil {
			// Score faithfulness
			scoreFaithFulness, err = judge.ScoreFaithfulness(context.Background(), item.GeneratedAnswer, item.RetrievedContexts)
			if err != nil {
				log.Printf("Failed to score faithfulness for item: %v\n", err)
			}

			// Score answer relevancy
			scoreAnswerRelevancy, err = judge.ScoreAnswerRelevancy(context.Background(), item.Question, item.GeneratedAnswer)
			if err != nil {
				log.Printf("Failed to score answer relevancy for item: %v\n", err)
			}
		}

		datasetResultEval := DatasetResultEval{
			Item:                 item,
			Hit:                  rank,
			Status:               GetStatus(rank),
			ScoreAnswerRelevancy: scoreAnswerRelevancy,
			ScoreFaithfulness:    scoreFaithFulness,
		}
		datasetResultEvals = append(datasetResultEvals, datasetResultEval)
	}

	resultsMetrics := calculateAverageMetrics(datasetResultEvals)

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

func calculateAverageMetrics(results []DatasetResultEval) RetrievalMetrics {
	var totalHitRate, totalMRR, totalPrecisionAt1, totalScoreAnswerRelevancy, totalScoreFaithfulness float64
	total := len(results)

	if total == 0 {
		return RetrievalMetrics{}
	}

	// Tính tổng từng chỉ số từ tất cả các item
	for _, result := range results {
		totalHitRate += float64(result.Hit)
		totalMRR += float64(result.Hit)
		totalPrecisionAt1 += float64(result.Hit)
		totalScoreAnswerRelevancy += result.ScoreAnswerRelevancy
		totalScoreFaithfulness += result.ScoreFaithfulness
	}

	// Tính trung bình cho từng chỉ số
	avgHitRate := totalHitRate / float64(total)
	avgMRR := totalMRR / float64(total)
	avgPrecisionAt1 := totalPrecisionAt1 / float64(total)
	avgScoreAnswerRelevancy := totalScoreAnswerRelevancy / float64(total)
	avgScoreFaithfulness := totalScoreFaithfulness / float64(total)

	return RetrievalMetrics{
		HitRate:                 avgHitRate,
		MRR:                     avgMRR,
		PrecisionAt1:            avgPrecisionAt1,
		AvgScoreAnswerRelevancy: avgScoreAnswerRelevancy,
		AvgScoreFaithfulness:    avgScoreFaithfulness,
	}
}
