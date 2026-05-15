package database

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"mcp-server/internal/evaluation"
	"mcp-server/internal/llm"
	"os"
	"strings"
	"time"
)

type EvalRow struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Context  string `json:"context"`
}

// generateEvalRow create pair question and answer from context using llm
func generateEvalRow(ctx context.Context, client llm.LLMModel, chunk string) (*EvalRow, error) {
	prompt := fmt.Sprintf(`
Bạn là một chuyên gia tạo tập dữ liệu đánh giá cho hệ thống RAG của Salon tóc. 
Dựa trên nội dung (Context) dưới đây, hãy tạo ra MỘT cặp Câu hỏi và Câu trả lời tương ứng.

Yêu cầu:
1. Câu hỏi (question) phải thực tế, giống cách khách hàng hỏi về dịch vụ, giá cả, hoặc thông tin chi nhánh.
2. Câu trả lời (answer) phải chính xác, ngắn gọn và dựa hoàn toàn vào Context.
3. CHỈ trả về duy nhất JSON theo format: {"question": "...", "answer": "..."}

Context:
%s
`, chunk)

	resp, err := client.Chat(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var row EvalRow
	if err := json.Unmarshal([]byte(resp), &row); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from LLM: %v", err)
	}

	row.Context = chunk
	return &row, nil
}

// GenDataset create dataset for benchmark test [question, answer, context]
func GenDataset(filePath string) error {
	ctx := context.Background()
	client, err := llm.NewLLM(llm.NewOpenAIEmbeddingConfig())

	if err != nil {
		return fmt.Errorf("failed to init client for llm and embedding: %w", err)
	}

	db, err := NewDB(ctx, NewDBConfig())
	if err != nil {
		return fmt.Errorf("failed init database: %w", err)
	}

	documents, err := db.ListDocuments(ctx, 5, 0)
	if err != nil {
		return err
	}

	var dataset []EvalRow
	for _, doc := range documents {
		fmt.Printf("Generating Q&A for chunk ID: %s...\n", doc.ID)
		row, err := generateEvalRow(ctx, client, doc.Content)
		if err != nil {
			fmt.Printf("Error generating for doc %s: %v\n", doc.ID, err)
			continue
		}
		dataset = append(dataset, *row)
	}

	// Lưu kết quả ra file CSV
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create csv file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Viết tiêu đề cột
	writer.Write([]string{"question", "answer", "context"})

	for _, row := range dataset {
		err := writer.Write([]string{row.Question, row.Answer, row.Context})
		if err != nil {
			return fmt.Errorf("failed to write record to csv: %w", err)
		}
	}

	fmt.Printf("\nĐã lưu %d bộ dữ liệu vào file: %s\n", len(dataset), filePath)
	return nil
}

// Evaluation test retrieval performance [Hit Rate, MRR, P@1, NDCG, Time Search]
func Evaluation(filePath string, verbose bool) error {
	ctx := context.Background()
	// 1. Khởi tạo DB
	db, err := NewDB(ctx, NewDBConfig())
	if err != nil {
		return err
	}

	// 2. Đọc file dataset
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file dataset: %w", err)
	}
	defer f.Close()

	// create reader to read file
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// 3. Mở file để lưu kết quả chi tiết
	resultPath := strings.Replace(filePath, ".csv", "_results.csv", 1)

	// 4. Evaluation
	fmt.Println("\nBắt đầu đánh giá hiệu năng Retrieval (Hit Rate, MRR, P@1)...")
	var dataEval []evaluation.DatasetItem
	var timeSearch map[int]float64

	for i, record := range records {
		if i == 0 {
			continue
		}
		question := record[0]
		expectedContext := record[2]

		// perform search
		startTime := time.Now()
		results, err := db.HybridSearch(ctx, HybridSearchParams{
			Query:        question,
			Limit:        5,
			BM25Weight:   0.5,
			VectorWeight: 0.5,
		})
		endTime := time.Now()
		timeSearch[i] = endTime.Sub(startTime).Seconds()
		// format to []string
		var retrievedContext []string
		for _, res := range results {
			retrievedContext = append(retrievedContext, res.Document.Content)
		}

		if err != nil {
			fmt.Printf("Lỗi search câu %d: %v\n", i, err)
			continue
		}

		dataEval = append(dataEval, evaluation.DatasetItem{
			Question:           question,
			GroundTruthContext: expectedContext,
			RetrievedContexts:  retrievedContext,
		})
	}
	resultsDetail, resultMetrics := evaluation.EvaluateRetrieval(dataEval)
	// inject time search to result metrics
	for i, res := range resultsDetail {
		res.TimeSearch = timeSearch[i]
	}
	resultMetrics.AvgTimeSearch = calculateAverageTime(timeSearch)
	// save result to csv  and print metrics summary if verbose is true
	if verbose {
		evaluation.PrintSummary(resultMetrics)
		evaluation.SaveDetailedResultsToCSV(resultsDetail, resultPath)
	}
	return nil
}

func calculateAverageTime(timeSearch map[int]float64) float64 {
	var total float64
	for _, time := range timeSearch {
		total += time
	}
	return total / float64(len(timeSearch))
}
