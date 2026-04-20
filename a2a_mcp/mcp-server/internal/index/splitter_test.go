package index

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDocument(t *testing.T) {
	file_path := "../../../data/VTI_Quy định_Quy định điền thông tin trên hệ thống VMS_v2.0.pdf"
	docs, err := loadDocs(context.Background(), file_path)

	assert.NoError(t, err)
	assert.Condition(t, func() (success bool) {
		return len(docs) > 0
	}, "The len documents should be > 0")
}

func TestFormatDocument(t *testing.T) {
	tenantID := "11111111-1111-1111-1111-111111111111"
	file_path := "../../../data/VTI_Quy định_Quy định điền thông tin trên hệ thống VMS_v2.0.pdf"
	docs, _ := loadDocs(context.Background(), file_path)
	docsFormatted, err := convertToDocumentFormat(docs, tenantID)
	assert.NoError(t, err)

	assert.Condition(t, func() (success bool) {
		return len(docsFormatted) > 0
	}, "The len documents should be > 0")

	fmt.Println(docsFormatted[0])
}

func TestLoadBatchDocument(t *testing.T) {
	filePaths := []string{
		"../../../data/VTI_Quy định thưởng đề xuất IP Kaizen-2019_V1.0.pdf",
		"../../../data/VTI_Quy định_Quy định điền thông tin trên hệ thống VMS_v2.0.pdf",
		"../../../data/VTI_Quy định_Quy định tạm ứng lương_v3.0.pdf",
		// "../../../data/VTI_Thỏa thuận tạm ứng lương_v1.0.pdf",
	}
	docs, err := loadBatchDocs(context.Background(), filePaths)

	assert.NoError(t, err)
	assert.Condition(t, func() (success bool) {
		return len(docs) > 0
	}, "The len documents should be > 0")
}

func TestFormatBatchDocument(t *testing.T) {
	tenantID := "11111111-1111-1111-1111-111111111111"
	filePaths := []string{
		"../../../data/VTI_Quy định thưởng đề xuất IP Kaizen-2019_V1.0.pdf",
		"../../../data/VTI_Quy định_Quy định điền thông tin trên hệ thống VMS_v2.0.pdf",
		"../../../data/VTI_Quy định_Quy định tạm ứng lương_v3.0.pdf",
		// "../../../data/VTI_Thỏa thuận tạm ứng lương_v1.0.pdf",
	}
	docs, _ := loadBatchDocs(context.Background(), filePaths)
	docsFormatted, err := convertToDocumentFormat(docs, tenantID)
	assert.NoError(t, err)

	assert.Condition(t, func() (success bool) {
		return len(docsFormatted) > 0
	}, "The len documents should be > 0")

	fmt.Println(docsFormatted[0])
}
