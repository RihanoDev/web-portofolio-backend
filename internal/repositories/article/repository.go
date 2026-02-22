package article

import (
	"web-porto-backend/internal/domain/models"

	"gorm.io/gorm"
)

type Repository interface {
	Create(article *models.Article) error
	GetByID(id string) (*models.Article, error)
	GetAll(limit, offset int) ([]*models.Article, int64, error)
	Update(article *models.Article) error
	Delete(id string) error
	GetBySlug(slug string) (*models.Article, error)
	GetByAuthorID(authorID int, limit, offset int) ([]*models.Article, int64, error)
	GetPublished(limit, offset int) ([]*models.Article, int64, error)
	GetByCategory(categoryID int, limit, offset int) ([]*models.Article, int64, error)
	GetByTag(tagID int, limit, offset int) ([]*models.Article, int64, error)
	UpdateArticleCategories(articleID string, categoryIDs []int) error
	UpdateArticleTags(articleID string, tagIDs []int) error
	UpdateArticleImages(articleID string, images []models.ArticleImage) error
	UpdateArticleVideos(articleID string, videos []models.ArticleVideo) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(article *models.Article) error {
	return r.db.Create(article).Error
}

func (r *repository) GetByID(id string) (*models.Article, error) {
	// Check if ID is a temporary ID from frontend
	if len(id) > 0 && id[:5] == "temp-" {
		return nil, gorm.ErrRecordNotFound
	}

	var article models.Article
	err := r.db.Preload("Author").
		Preload("Categories").
		Preload("Tags").
		Preload("Images").
		Preload("Videos").
		Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *repository) GetAll(limit, offset int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	// Count total records
	r.db.Model(&models.Article{}).Count(&total)

	// Get paginated results with preloaded associations
	err := r.db.Preload("Author").
		Preload("Categories").
		Preload("Tags").
		Preload("Images").
		Preload("Videos").
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&articles).Error

	return articles, total, err
}

func (r *repository) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

func (r *repository) Delete(id string) error {
	// Check if ID is a temporary ID from frontend
	if len(id) > 0 && id[:5] == "temp-" {
		// Nothing to delete for temporary IDs
		return nil
	}

	return r.db.Delete(&models.Article{}, "id = ?", id).Error
}

func (r *repository) GetBySlug(slug string) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Author").
		Preload("Categories").
		Preload("Tags").
		Preload("Images").
		Preload("Videos").
		Where("slug = ?", slug).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *repository) GetByAuthorID(authorID int, limit, offset int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	query := r.db.Model(&models.Article{}).Where("author_id = ?", authorID)
	query.Count(&total)

	err := query.Preload("Author").
		Preload("Categories").
		Preload("Tags").
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&articles).Error

	return articles, total, err
}

func (r *repository) GetPublished(limit, offset int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	query := r.db.Model(&models.Article{}).Where("status = ?", "published")
	query.Count(&total)

	err := query.Preload("Author").
		Preload("Categories").
		Preload("Tags").
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&articles).Error

	return articles, total, err
}

func (r *repository) GetByCategory(categoryID int, limit, offset int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	query := r.db.Model(&models.Article{}).
		Joins("JOIN article_categories ac ON ac.article_id = articles.id").
		Where("ac.category_id = ?", categoryID)

	query.Count(&total)

	err := query.Preload("Author").
		Preload("Categories").
		Preload("Tags").
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&articles).Error

	return articles, total, err
}

func (r *repository) GetByTag(tagID int, limit, offset int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	query := r.db.Model(&models.Article{}).
		Joins("JOIN article_tags at ON at.article_id = articles.id").
		Where("at.tag_id = ?", tagID)

	query.Count(&total)

	err := query.Preload("Author").
		Preload("Categories").
		Preload("Tags").
		Preload("Images").
		Preload("Videos").
		Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&articles).Error

	return articles, total, err
}

func (r *repository) UpdateArticleCategories(articleID string, categoryIDs []int) error {
	var article models.Article
	if err := r.db.First(&article, "id = ?", articleID).Error; err != nil {
		return err
	}
	var categories []models.Category
	if len(categoryIDs) > 0 {
		r.db.Where("id IN ?", categoryIDs).Find(&categories)
	}
	return r.db.Model(&article).Association("Categories").Replace(categories)
}

func (r *repository) UpdateArticleTags(articleID string, tagIDs []int) error {
	var article models.Article
	if err := r.db.First(&article, "id = ?", articleID).Error; err != nil {
		return err
	}
	var tags []models.Tag
	if len(tagIDs) > 0 {
		r.db.Where("id IN ?", tagIDs).Find(&tags)
	}
	return r.db.Model(&article).Association("Tags").Replace(tags)
}

func (r *repository) UpdateArticleImages(articleID string, images []models.ArticleImage) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ?", articleID).Delete(&models.ArticleImage{}).Error; err != nil {
			return err
		}
		if len(images) > 0 {
			if err := tx.Create(&images).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) UpdateArticleVideos(articleID string, videos []models.ArticleVideo) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ?", articleID).Delete(&models.ArticleVideo{}).Error; err != nil {
			return err
		}
		if len(videos) > 0 {
			if err := tx.Create(&videos).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
