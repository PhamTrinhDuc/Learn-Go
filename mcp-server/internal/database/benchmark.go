package database

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"mcp-server/internal/evaluation"
	"mcp-server/internal/llm"
	"mcp-server/internal/utils"
	"os"
	"strings"
)

type EvalRow struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Context  string `json:"context"`
}

func getConfig() (llm.Config, DBConfig) {
	return llm.Config{
			LLM: llm.ProviderConfig{
				Provider: llm.ProviderGroq,
				Model:    "llama-3.3-70b-versatile",
				APIKey:   utils.GetEnvString("GROQ_API_KEY", ""),
			},
			Embed: llm.ProviderConfig{
				Provider: llm.ProviderOllama,
				BaseURL:  "http://localhost:11434",
				Model:    "qwen3-embedding:0.6b",
			},
		},
		DBConfig{
			Host:     utils.GetEnvString("DB_HOST", "localhost"),
			Port:     utils.GetEnvInt("DB_PORT", 5433),
			User:     utils.GetEnvString("DB_USER", "mcp_user"),
			Password: utils.GetEnvString("DB_PASSWORD", "mcp_password"),
			DBName:   utils.GetEnvString("DB_NAME", "salon_chain"),
			SSLMode:  utils.GetEnvString("DB_SSLMODE", "disable"),
			MaxConns: int32(utils.GetEnvInt("MAX_CONNS", 10)),
			MinConns: int32(utils.GetEnvInt("MAX_CONNS", 2)),
		}
}

func generateEvalRow(ctx context.Context, client *llm.Client, chunk string) (*EvalRow, error) {
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

func GenDataset(filePath string) error {
	ctx := context.Background()
	clientCfg, dbCfg := getConfig()
	client, err := llm.NewClient(clientCfg)

	if err != nil {
		return fmt.Errorf("failed to init client for llm and embedding: %w", err)
	}

	db, err := NewDB(ctx, dbCfg)
	if err != nil {
		return fmt.Errorf("failed init database: %w", err)
	}

	documents, err := db.ListDocuments(ctx, 5, 0) // Giới hạn 10 bản ghi để test
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

func Evaluation(filePath string, verbose bool) error {
	ctx := context.Background()
	// 1. Khởi tạo DB
	_, dbCfg := getConfig()
	db, err := NewDB(ctx, dbCfg)
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

	for i, record := range records {
		if i == 0 {
			continue
		}
		question := record[0]
		expectedContext := record[2]

		// perform search
		results, err := db.HybridSearch(ctx, HybridSearchParams{
			Query:        question,
			Limit:        5,
			BM25Weight:   0.5,
			VectorWeight: 0.5,
		})
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
	if verbose {
		evaluation.PrintSummary(resultMetrics)
		evaluation.SaveDetailedResultsToCSV(resultsDetail, resultPath)
	}
	return nil
}
