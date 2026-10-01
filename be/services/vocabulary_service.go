package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"vocabulary/models"
	"vocabulary/repositories"

	"github.com/xuri/excelize/v2"
)

var ErrInvalidVocabularyExcelHeader = errors.New("invalid vocabulary excel header")

type VocabularyService interface {
	UploadExcel(ctx context.Context, reader io.Reader) (*models.UploadVocabularyExcelResponse, error)
}

type vocabularyService struct {
	vocabularyRepository repositories.VocabularyRepository
}

func NewVocabularyService(vocabularyRepository repositories.VocabularyRepository) VocabularyService {
	return &vocabularyService{vocabularyRepository: vocabularyRepository}
}

func (s *vocabularyService) UploadExcel(ctx context.Context, reader io.Reader) (*models.UploadVocabularyExcelResponse, error) {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("open excel file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	sheetName := firstSheetName(file)
	if sheetName == "" {
		return nil, ErrInvalidVocabularyExcelHeader
	}

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("read excel rows: %w", err)
	}

	if len(rows) == 0 || !hasExpectedVocabularyHeader(rows[0]) {
		return nil, ErrInvalidVocabularyExcelHeader
	}

	response := &models.UploadVocabularyExcelResponse{
		Message: "uploaded vocabulary successfully",
	}
	seenRows := make(map[string]struct{})

	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		rowNumber := rowIndex + 1
		row := normalizeVocabularyRow(rows[rowIndex])
		if isEmptyVocabularyRow(row) {
			continue
		}
		if hasEmptyVocabularyCell(row) {
			response.Skipped++
			continue
		}

		if rowErrors := validateVocabularyRow(rowNumber, row); len(rowErrors) > 0 {
			response.Errors = append(response.Errors, rowErrors...)
			response.Invalid++
			continue
		}

		key := vocabularyDuplicateKey(row.Eng, row.PartsOfSpeech)
		if _, ok := seenRows[key]; ok {
			response.Skipped++
			continue
		}
		seenRows[key] = struct{}{}

		duplicate, err := s.vocabularyRepository.FindByEngAndPartsOfSpeech(ctx, row.Eng, row.PartsOfSpeech)
		if err != nil {
			return nil, fmt.Errorf("find duplicate vocabulary: %w", err)
		}
		if duplicate != nil {
			response.Skipped++
			continue
		}

		existing, err := s.vocabularyRepository.FindByEng(ctx, row.Eng)
		if err != nil {
			return nil, fmt.Errorf("find vocabulary by eng: %w", err)
		}
		if existing != nil {
			existing.PartsOfSpeech = row.PartsOfSpeech
			existing.Thai = row.Thai
			existing.Meaning = row.Meaning
			existing.Synonyms = row.Synonyms
			existing.Level = row.Level

			if err := s.vocabularyRepository.Update(ctx, existing); err != nil {
				return nil, fmt.Errorf("update vocabulary: %w", err)
			}

			response.Updated++
			continue
		}

		vocabulary := models.Vocabulary{
			Eng:           row.Eng,
			PartsOfSpeech: row.PartsOfSpeech,
			Thai:          row.Thai,
			Meaning:       row.Meaning,
			Synonyms:      row.Synonyms,
			Level:         row.Level,
		}
		if err := s.vocabularyRepository.Create(ctx, &vocabulary); err != nil {
			return nil, fmt.Errorf("create vocabulary: %w", err)
		}

		response.Inserted++
	}

	if response.Invalid > 0 {
		response.Message = "uploaded vocabulary with some invalid rows"
	}

	return response, nil
}

type vocabularyExcelRow struct {
	Eng           string
	PartsOfSpeech string
	Thai          string
	Meaning       string
	Synonyms      string
	Level         string
}

func firstSheetName(file *excelize.File) string {
	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return ""
	}

	return sheets[0]
}

func hasExpectedVocabularyHeader(row []string) bool {
	expectedHeaders := []string{"eng", "parts_of_speech", "thai", "meaning", "synonyms", "level"}
	if len(row) < len(expectedHeaders) {
		return false
	}

	for index, expected := range expectedHeaders {
		header := strings.ToLower(strings.TrimSpace(row[index]))
		if index == 1 && header == "part_of_speech" {
			continue
		}
		if header != expected {
			return false
		}
	}

	return true
}

func normalizeVocabularyRow(row []string) vocabularyExcelRow {
	return vocabularyExcelRow{
		Eng:           valueAt(row, 0),
		PartsOfSpeech: valueAt(row, 1),
		Thai:          valueAt(row, 2),
		Meaning:       valueAt(row, 3),
		Synonyms:      valueAt(row, 4),
		Level:         valueAt(row, 5),
	}
}

func valueAt(row []string, index int) string {
	if index >= len(row) {
		return ""
	}

	return strings.TrimSpace(row[index])
}

func isEmptyVocabularyRow(row vocabularyExcelRow) bool {
	return row.Eng == "" &&
		row.PartsOfSpeech == "" &&
		row.Thai == "" &&
		row.Meaning == "" &&
		row.Synonyms == "" &&
		row.Level == ""
}

func hasEmptyVocabularyCell(row vocabularyExcelRow) bool {
	return row.Eng == "" ||
		row.PartsOfSpeech == "" ||
		row.Thai == "" ||
		row.Meaning == "" ||
		row.Synonyms == "" ||
		row.Level == ""
}

func validateVocabularyRow(rowNumber int, row vocabularyExcelRow) []models.UploadVocabularyRowError {
	var rowErrors []models.UploadVocabularyRowError

	if row.Eng == "" {
		rowErrors = append(rowErrors, newUploadVocabularyRowError(rowNumber, "eng", "eng is required"))
	}
	if row.PartsOfSpeech == "" {
		rowErrors = append(rowErrors, newUploadVocabularyRowError(rowNumber, "parts_of_speech", "parts_of_speech is required"))
	}
	if row.Thai == "" {
		rowErrors = append(rowErrors, newUploadVocabularyRowError(rowNumber, "thai", "thai is required"))
	}

	return rowErrors
}

func newUploadVocabularyRowError(row int, field string, message string) models.UploadVocabularyRowError {
	return models.UploadVocabularyRowError{
		Row:     row,
		Field:   field,
		Message: message,
	}
}

func vocabularyDuplicateKey(eng string, partsOfSpeech string) string {
	return eng + "\x00" + partsOfSpeech
}
