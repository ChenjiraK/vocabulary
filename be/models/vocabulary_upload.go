package models

import "mime/multipart"

type UploadVocabularyExcelRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type UploadVocabularyExcelResponse struct {
	Message  string                     `json:"message"`
	Inserted int                        `json:"inserted"`
	Updated  int                        `json:"updated"`
	Skipped  int                        `json:"skipped"`
	Invalid  int                        `json:"invalid"`
	Errors   []UploadVocabularyRowError `json:"errors,omitempty"`
}

type UploadVocabularyRowError struct {
	Row     int    `json:"row"`
	Field   string `json:"field"`
	Message string `json:"message"`
}
