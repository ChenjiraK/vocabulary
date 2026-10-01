package handlers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"vocabulary/models"
	"vocabulary/services"

	"github.com/gin-gonic/gin"
)

// handler มีการเรียกใช้ service vocabularyService เพื่อทำงานกับข้อมูล vocabulary และมี method UploadExcel สำหรับอัปโหลดไฟล์ Excel ของคำศัพท์ โดยตรวจสอบว่าไฟล์ที่อัปโหลดเป็นไฟล์ Excel (.xlsx) และเรียกใช้ service เพื่อประมวลผลไฟล์นั้น
type VocabularyHandler struct {
	vocabularyService services.VocabularyService
}

func NewVocabularyHandler(vocabularyService services.VocabularyService) *VocabularyHandler {
	return &VocabularyHandler{vocabularyService: vocabularyService}
}

func (h *VocabularyHandler) UploadExcel(c *gin.Context) {
	var request models.UploadVocabularyExcelRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "file is required",
		})
		return
	}

	if !isExcelFile(request.File.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "file must be .xlsx",
		})
		return
	}

	file, err := request.File.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot open uploaded file",
		})
		return
	}
	defer func() {
		_ = file.Close()
	}()

	response, err := h.vocabularyService.UploadExcel(c.Request.Context(), file)
	if err != nil {
		if errors.Is(err, services.ErrInvalidVocabularyExcelHeader) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "excel header must be eng, part_of_speech, thai, meaning, synonyms, level",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "cannot upload vocabulary excel",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func isExcelFile(filename string) bool {
	//strings.EqualFold() เป็นฟังก์ชันใน Go ที่ใช้เปรียบเทียบสตริงสองตัวโดยไม่สนใจตัวพิมพ์ใหญ่หรือตัวพิมพ์เล็ก (case-insensitive comparison) ซึ่งจะคืนค่า true หากสตริงทั้งสองเหมือนกัน และ false หากไม่เหมือนกัน
	return strings.EqualFold(filepath.Ext(filename), ".xlsx")
}
