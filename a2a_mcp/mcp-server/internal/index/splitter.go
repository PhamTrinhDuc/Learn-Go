package index

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"regexp"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"

	"learn-go/a2a_mcp/mcp-server/internal/database"
)

func loadDocs(ctx context.Context, filePath string) ([]schema.Document, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(filePath))

	// Splitter dùng chung — RecursiveCharacter tốt nhất cho mixed content
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(512),
		textsplitter.WithChunkOverlap(50),
		textsplitter.WithSeparators([]string{"\n\n", "\n", ".", "!", "?", " "}),
	)

	var loader documentloaders.Loader

	switch ext {
	case ".pdf":
		info, err := f.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to get file info: %w", err)
		}
		loader = documentloaders.NewPDF(f, info.Size())
	case ".csv":
		loader = documentloaders.NewCSV(f)
	case ".html", ".htm":
		loader = documentloaders.NewHTML(f)
	default: // .txt, .md, .doc đã convert...
		loader = documentloaders.NewText(f)
	}

	// LoadAndSplit = load + chunk trong 1 bước
	docs, err := loader.LoadAndSplit(ctx, splitter)
	if err != nil {
		return nil, err
	}

	// Inject metadata & Clean content
	for i := range docs {
		docs[i].PageContent = cleanContent(docs[i].PageContent)
		docs[i].Metadata["source"] = filepath.Base(filePath)
		docs[i].Metadata["file_type"] = ext
	}

	return docs, nil
}

func cleanContent(content string) string {
	// 1. Thêm khoảng trắng sau các dấu liệt kê như ●, -, * nếu dính với chữ
	reBullet := regexp.MustCompile(`([●\-\*])([^\s])`)
	content = reBullet.ReplaceAllString(content, "$1 $2")

	// 2. Thêm khoảng trắng sau các dấu câu nếu dính với chữ (tránh dính "đơnvị.")
	rePunct := regexp.MustCompile(`([,;:\.\?!])([^\s\d])`)
	content = rePunct.ReplaceAllString(content, "$1 $2")

	// 3. Xử lý khoảng trắng thừa
	return strings.Join(strings.Fields(content), " ")
}

func loadBatchDocs(ctx context.Context, filePaths []string) ([]schema.Document, error) {
	// 1. Khai báo các công cụ
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allDocs []schema.Document

	// 2. Tạo Semaphore để giới hạn 5 file đồng thời
	// Kênh này chứa 5 "vé", ai có vé mới được chạy
	sem := make(chan struct{}, 5)

	// 3. Kênh để bắt lỗi (vì hàm goroutine không return trực tiếp về hàm cha được)
	errChan := make(chan error, len(filePaths))

	for _, path := range filePaths {
		wg.Add(1) // Đếm thêm 1 việc cần làm
		// Đợi lấy "vé" từ semaphore, nếu đủ 5 người rồi thì sẽ kẹt lại ở đây
		sem <- struct{}{}

		go func(p string) {
			defer wg.Done() // Đảm bảo dù lỗi hay không cũng sẽ trừ máy đếm khi xong

			defer func() {
				<-sem
			}() // Xong việc thì trả lại vé cho người sau. Không để ngoài func vì phải chờ for xong mới thực thi => không ai trả vé

			docs, err := loadDocs(ctx, p)
			if err != nil {
				fmt.Printf("Error loading file %s: %v\n", p, err) // Thêm dòng này
				errChan <- err                                    // Nếu lỗi thì gửi vào kênh lỗi
				return
			}
			// Xong viêc thì xin khóa viết vào sổ chung
			mu.Lock()
			allDocs = append(allDocs, docs...)
			mu.Unlock() // Xong việc thì trả khóa
		}(path)
	}
	// 4. Đợi tất cả hoàn thành
	wg.Wait()
	close(errChan)

	// 5. Kiểm tra xem có lỗi không
	if len(errChan) > 0 {
		return nil, <-errChan
	}
	return allDocs, nil
}

func convertToDocumentFormat(docs []schema.Document, tenantID string) ([]*database.Document, error) {
	documents := make([]*database.Document, 0, len(docs))

	for _, doc := range docs {
		// Khởi tạo struct mới
		formattedDoc := &database.Document{
			ID:       uuid.New().String(),
			TenantID: tenantID,
			Content:  doc.PageContent,
			Metadata: doc.Metadata,
		}

		// Ép kiểu metadata["source"] về string
		if source, ok := doc.Metadata["source"].(string); ok {
			ext := filepath.Ext(source)
			formattedDoc.Title = strings.TrimSuffix(source, ext)
		}

		documents = append(documents, formattedDoc)
	}
	return documents, nil
}
