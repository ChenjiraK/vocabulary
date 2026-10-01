package models

import "time"

// NOTE: exported fields ตัวแปรที่ขึ้นต้นด้วยตัวใหญ่หมายถึง exported คือสามารถเข้าถึงได้จาก package อื่นเช่น Handler, Service, Repository สามารถเรียกใช้ struct Vocabulary ได้
type Vocabulary struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Eng           string    `gorm:"column:eng;size:255;not null" json:"eng"`
	PartsOfSpeech string    `gorm:"column:parts_of_speech;size:100;not null" json:"parts_of_speech"`
	Thai          string    `gorm:"column:thai;size:255;not null" json:"thai"`
	Meaning       string    `gorm:"column:meaning;type:text" json:"meaning"`
	Synonyms      string    `gorm:"column:synonyms;type:text" json:"synonyms"`
	Level         string    `gorm:"column:level;size:5;not null" json:"level"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Vocabulary) TableName() string {
	return "vocabulary"
}
