package services

import (
	"bytes"
	"context"
	"testing"

	"vocabulary/models"

	"github.com/xuri/excelize/v2"
)

func TestUploadExcelSkipsRowsWithEmptyCells(t *testing.T) {
	var file bytes.Buffer
	workbook := excelize.NewFile()
	defer func() {
		_ = workbook.Close()
	}()

	rows := [][]string{
		{"eng", "part_of_speech", "thai", "meaning", "synonyms", "level"},
		{"apple", "noun", "แอปเปิล", "fruit", "pome", "A1"},
		{"banana", "noun", "กล้วย", "", "plantain", "A1"},
		{"cat", "noun", "แมว", "animal", "feline", ""},
		{"", "", "", "", "", ""},
	}

	for rowIndex, row := range rows {
		for columnIndex, value := range row {
			cell, err := excelize.CoordinatesToCellName(columnIndex+1, rowIndex+1)
			if err != nil {
				t.Fatalf("build cell name: %v", err)
			}
			if err := workbook.SetCellValue("Sheet1", cell, value); err != nil {
				t.Fatalf("set cell value: %v", err)
			}
		}
	}

	if err := workbook.Write(&file); err != nil {
		t.Fatalf("write workbook: %v", err)
	}

	repository := &fakeVocabularyRepository{}
	service := NewVocabularyService(repository)

	response, err := service.UploadExcel(context.Background(), bytes.NewReader(file.Bytes()))
	if err != nil {
		t.Fatalf("upload excel: %v", err)
	}

	if response.Inserted != 1 {
		t.Fatalf("expected inserted 1, got %d", response.Inserted)
	}
	if response.Skipped != 2 {
		t.Fatalf("expected skipped 2, got %d", response.Skipped)
	}
	if response.Invalid != 0 {
		t.Fatalf("expected invalid 0, got %d", response.Invalid)
	}
	if len(repository.created) != 1 {
		t.Fatalf("expected created records 1, got %d", len(repository.created))
	}
	if repository.created[0].Level != "A1" {
		t.Fatalf("expected created level A1, got %q", repository.created[0].Level)
	}
}

type fakeVocabularyRepository struct {
	created []models.Vocabulary
	updated []models.Vocabulary
}

func (r *fakeVocabularyRepository) Create(ctx context.Context, vocabulary *models.Vocabulary) error {
	r.created = append(r.created, *vocabulary)
	return nil
}

func (r *fakeVocabularyRepository) FindByEng(ctx context.Context, eng string) (*models.Vocabulary, error) {
	return nil, nil
}

func (r *fakeVocabularyRepository) FindByEngAndPartsOfSpeech(ctx context.Context, eng string, partsOfSpeech string) (*models.Vocabulary, error) {
	return nil, nil
}

func (r *fakeVocabularyRepository) Update(ctx context.Context, vocabulary *models.Vocabulary) error {
	r.updated = append(r.updated, *vocabulary)
	return nil
}
