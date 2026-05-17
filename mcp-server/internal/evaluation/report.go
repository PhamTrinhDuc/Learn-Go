package evaluation

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// SaveDetailedResultsToCSV store result detail to file
func SaveDetailedResultsToCSV(evals []DatasetResultEval, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	writer.Write([]string{"Question", "Expected Context", "Generated Answer", "Ground Truth Answer", "Faithfulness", "AnswerRelevancy", "Rank", "Time Search", "Status"})

	for _, e := range evals {
		err := writer.Write([]string{
			e.Item.Question,
			e.Item.GroundTruthContext,
			strings.Join(e.Item.RetrievedContexts, "\n"),
			e.Item.GeneratedAnswer,
			e.Item.GroundTruthAnswer,
			fmt.Sprintf("%.2f", e.ScoreFaithfulness),
			fmt.Sprintf("%.2f", e.ScoreAnswerRelevancy),
			fmt.Sprintf("%d", e.Hit),
			fmt.Sprintf("%.4f", e.TimeSearch),
			e.Status,
		})
		if err != nil {
			return err
		}
	}

	fmt.Printf("✅ Đã lưu log chi tiết tại: %s\n", filePath)
	return nil
}

// PrintSummary print result on console
func PrintSummary(metrics RetrievalMetrics) {
	fmt.Printf("\n--- KẾT QUẢ ĐÁNH GIÁ TỔNG HỢP ---\n")
	fmt.Printf("Hit Rate@5:      %.2f%%\n", metrics.HitRate)
	fmt.Printf("Precision@1:    %.2f%%\n", metrics.PrecisionAt1)
	fmt.Printf("MRR (Ranking):  %.4f\n", metrics.MRR)
	fmt.Printf("Avg Time Search: %.4f\n", metrics.AvgTimeSearch)
	fmt.Println("---------------------------------")
}
