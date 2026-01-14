package repository

import (
	"strconv"
	"time"

	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

type PostRepository interface {
	FindAll(offset, limit int, search string) ([]domain.Post, int64, error)
	FindByID(id int) (domain.Post, error)
	FindBySlugOrID(identifier string) (domain.Post, error)
	Create(post *domain.Post) error
	Update(post *domain.Post) error
	Delete(post *domain.Post, unscoped bool) error
	GetTagBySlug(slug string, name string) (domain.Tag, error)
	HasViewed(postID int, ip string, since time.Time) (bool, error)
	AddView(view *domain.PostView) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) FindAll(offset, limit int, search string) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	// Tambahkan subquery Select ini agar setiap item di list memiliki data views_count
	query := r.db.Model(&domain.Post{}).
		Select("posts.*, (SELECT COUNT(*) FROM post_views WHERE post_views.post_id = posts.id) as views_count").
		Preload("Tags").
		Preload("Category")

	if search != "" {
		searchKeyword := "%" + search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchKeyword, searchKeyword)
	}

	query.Count(&total)

	// Pastikan urutan query tetap benar
	err := query.Limit(limit).Offset(offset).Order("published_at DESC").Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) FindByID(id int) (domain.Post, error) {
	var post domain.Post
	err := config.DB.Preload("Category").Preload("Tags").First(&post, id).Error
	return post, err
}

func (r *postRepository) FindBySlugOrID(identifier string) (domain.Post, error) {
	var post domain.Post

	// Siapkan base query dengan subquery views_count
	query := r.db.Preload("Category").Preload("Tags").
		Select("posts.*, (SELECT COUNT(*) FROM post_views WHERE post_views.post_id = posts.id) as views_count")

	// Cek apakah identifier adalah integer (ID)
	id, err := strconv.Atoi(identifier)

	if err == nil {
		// Jika sukses dikonversi ke angka, cari berdasarkan ID
		err = query.Where("id = ?", id).First(&post).Error
	} else {
		// Jika gagal (berarti itu slug string), cari berdasarkan Slug
		err = query.Where("slug = ?", identifier).First(&post).Error
	}

	return post, err
}

func (r *postRepository) Create(post *domain.Post) error {
	return config.DB.Create(post).Error
}

func (r *postRepository) Update(post *domain.Post) error {
	return config.DB.Save(post).Error
}

func (r *postRepository) Delete(post *domain.Post, unscoped bool) error {
	db := config.DB
	if unscoped {
		db = db.Unscoped()
	}
	return db.Delete(post).Error
}

func (r *postRepository) GetTagBySlug(slug string, name string) (domain.Tag, error) {
	var tag domain.Tag
	err := config.DB.Where(domain.Tag{Slug: slug}).Attrs(domain.Tag{Name: name}).FirstOrCreate(&tag).Error
	return tag, err
}

func (r *postRepository) HasViewed(postID int, ip string, since time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&domain.PostView{}).
		Where("post_id = ? AND ip_address = ? AND viewed_at > ?", postID, ip, since).
		Count(&count).Error
	return count > 0, err
}

func (r *postRepository) AddView(view *domain.PostView) error {
	return r.db.Create(view).Error
}
