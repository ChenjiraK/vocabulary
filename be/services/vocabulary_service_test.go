package services

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"testing"

	"vocabulary/models"

	"github.com/xuri/excelize/v2"
)

func TestLevels(t *testing.T) {
	for _, tc := range []struct {
		name     string
		progress []models.VocabularyLevelProgress
		want     models.VocabularyLevelsResponse
	}{
		{
			name: "no vocabulary",
			want: models.VocabularyLevelsResponse{"A1": "inactive", "A2": "inactive", "B1": "inactive", "B2": "inactive", "C1": "inactive", "C2": "inactive"},
		},
		{
			name: "requested example",
			progress: []models.VocabularyLevelProgress{
				{Level: "B1", Total: 2, Completed: 2},
				{Level: "B2", Total: 3, Completed: 1},
			},
			want: models.VocabularyLevelsResponse{"A1": "inactive", "A2": "inactive", "B1": "complete", "B2": "active", "C1": "inactive", "C2": "inactive"},
		},
		{
			name: "unplayed vocabulary and all supported levels",
			progress: []models.VocabularyLevelProgress{
				{Level: "A1", Total: 1},
				{Level: "A2", Total: 1, Completed: 1},
				{Level: "B1", Total: 1, Completed: 1},
				{Level: "B2", Total: 1, Completed: 1},
				{Level: "C1", Total: 1, Completed: 1},
				{Level: "C2", Total: 1, Completed: 1},
				{Level: "unknown", Total: 1},
			},
			want: models.VocabularyLevelsResponse{"A1": "active", "A2": "complete", "B1": "complete", "B2": "complete", "C1": "complete", "C2": "complete"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repository := &fakeVocabularyRepository{progress: tc.progress}
			got, err := NewVocabularyService(repository).Levels(context.Background(), 42)
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, %v; want %v", got, err, tc.want)
			}
			if repository.progressUserID != 42 {
				t.Fatalf("expected progress for user 42, got %d", repository.progressUserID)
			}
		})
	}
}

func TestLevelsRepositoryError(t *testing.T) {
	wantErr := errors.New("database unavailable")
	repository := &fakeVocabularyRepository{progressErr: wantErr}
	response, err := NewVocabularyService(repository).Levels(context.Background(), 42)
	if response != nil || !errors.Is(err, wantErr) {
		t.Fatalf("expected database error and no response, got %v, %v", response, err)
	}
}

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
	progress       []models.VocabularyLevelProgress
	progressErr    error
	progressUserID uint
	created        []models.Vocabulary
	updated        []models.Vocabulary
}

func (r *fakeVocabularyRepository) ListLevelProgress(ctx context.Context, userID uint) ([]models.VocabularyLevelProgress, error) {
	r.progressUserID = userID
	return r.progress, r.progressErr
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
