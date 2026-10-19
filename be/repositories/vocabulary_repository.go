package repositories

import (
	"context"

	"vocabulary/models"

	"gorm.io/gorm"
)

//NOTE: return models.Vocabulary เหมาะกับ ข้อมูลต้องมีแน่นอน, ต้องการ API ที่ใช้ง่ายและปลอดภัยจาก nil
//NOTE: return *models.Vocabulary เหมาะกับ ข้อมูลอาจไม่มี, ต้องการ API ที่ยืดหยุ่นและสามารถตรวจสอบ nil ได้,ต้องการแก้ไขข้อมูลต้นฉบับผ่าน pointer และ object เป็น optional

type VocabularyRepository interface {
	ListLevelProgress(ctx context.Context, userID uint) ([]models.VocabularyLevelProgress, error)
	Create(ctx context.Context, vocabulary *models.Vocabulary) error
	FindByEng(ctx context.Context, eng string) (*models.Vocabulary, error)
	FindByEngAndPartsOfSpeech(ctx context.Context, eng string, partsOfSpeech string) (*models.Vocabulary, error)
	Update(ctx context.Context, vocabulary *models.Vocabulary) error
}

type vocabularyRepository struct {
	db *gorm.DB
}

func NewVocabularyRepository(db *gorm.DB) VocabularyRepository {
	return &vocabularyRepository{db: db}
}

func (r *vocabularyRepository) ListLevelProgress(ctx context.Context, userID uint) ([]models.VocabularyLevelProgress, error) {
	var progress []models.VocabularyLevelProgress
	err := r.db.WithContext(ctx).
		Table("vocabulary AS v").
		Select(`v.level, COUNT(*) AS total,
			SUM(CASE WHEN EXISTS (
				SELECT 1 FROM game_answers AS ga
				WHERE ga.vocabulary_id = v.id AND ga.user_id = ? AND ga.is_correct = ?
			) THEN 1 ELSE 0 END) AS completed`, userID, true).
		Group("v.level").
		Scan(&progress).Error
	return progress, err
}

func (r *vocabularyRepository) Create(ctx context.Context, vocabulary *models.Vocabulary) error {
	return r.db.WithContext(ctx).Create(vocabulary).Error
}

func (r *vocabularyRepository) FindByEng(ctx context.Context, eng string) (*models.Vocabulary, error) {
	var vocabulary models.Vocabulary
	if err := r.db.WithContext(ctx).Where("eng = ?", eng).Order("id ASC").First(&vocabulary).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &vocabulary, nil
}

func (r *vocabularyRepository) FindByEngAndPartsOfSpeech(ctx context.Context, eng string, partsOfSpeech string) (*models.Vocabulary, error) {
	var vocabulary models.Vocabulary
	if err := r.db.WithContext(ctx).
		Where("eng = ? AND parts_of_speech = ?", eng, partsOfSpeech).
		Order("id ASC").
		First(&vocabulary).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &vocabulary, nil
}

func (r *vocabularyRepository) Update(ctx context.Context, vocabulary *models.Vocabulary) error {
	return r.db.WithContext(ctx).Save(vocabulary).Error
}
