package service

import (
	"math"
	"strings"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
	"github.com/garuda-labs-1/pmii-be/pkg/utils"
)

type PostService interface {
	GetAllPosts(page, limit int, search string) ([]responses.PostResponse, int, int64, error)
	CreatePost(req requests.PostCreateRequest) (responses.PostResponse, error)
	UpdatePost(id string, req requests.PostUpdateRequest) (responses.PostResponse, error)
	DeletePost(id string) error
	GetPostDetail(id string, ip, ua string) (responses.PostResponse, error)
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

// 1. GET ALL POSTS WITH PAGINATION
func (s *postService) GetAllPosts(page, limit int, search string) ([]responses.PostResponse, int, int64, error) {
	offset := (page - 1) * limit
	posts, total, err := s.repo.FindAll(offset, limit, search)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if lastPage < 1 {
		lastPage = 1
	}

	return responses.FromDomainListToPostResponse(posts), lastPage, total, nil
}

// 2. CREATE POST
func (s *postService) CreatePost(req requests.PostCreateRequest) (responses.PostResponse, error) {
	var featuredImage *string

	// Logika Upload Gambar ke Cloudinary
	if req.Image != nil {
		file, _ := req.Image.Open()
		defer file.Close()
		filename := "post-" + time.Now().Format("20060102-150405")

		url, err := utils.UploadToCloudinary(file, filename)
		if err == nil {
			featuredImage = &url
		}
	}

	// Logika Pembuatan Excerpt otomatis
	excerptText := req.Content
	if len(excerptText) > 150 {
		excerptText = excerptText[:150] + "..."
	}

	publishedTime := time.Now()
	post := domain.Post{
		Title:         req.Title,
		Content:       req.Content,
		Slug:          req.GetSlug(),
		CategoryID:    req.CategoryID,
		UserID:        1, // Default Admin
		Excerpt:       &excerptText,
		FeaturedImage: featuredImage,
		PublishedAt:   &publishedTime,
		Tags:          s.processTags(req.Tags), // Logika Many-to-Many Tags
	}

	if err := s.repo.Create(&post); err != nil {
		return responses.PostResponse{}, err
	}

	// Reload data untuk mendapatkan relasi Category & Tags lengkap
	updatedPost, _ := s.repo.FindByID(post.ID)
	return responses.FromDomainToPostResponse(updatedPost), nil
}

// 3. GET DETAIL POST
func (s *postService) GetPostDetail(id string, ip, ua string) (responses.PostResponse, error) {
	// 1. Ambil detail berita
	post, err := s.repo.FindBySlugOrID(id)
	if err != nil {
		return responses.PostResponse{}, err
	}

	// 2. Logika Anti-Spam: Cek apakah IP sudah melihat dalam 24 jam terakhir
	since := time.Now().Add(-24 * time.Hour)
	hasViewed, _ := s.repo.HasViewed(post.ID, ip, since)

	if !hasViewed {
		newView := domain.PostView{
			PostID:    post.ID,
			IPAddress: &ip,
			UserAgent: &ua,
			ViewedAt:  time.Now(),
		}

		// Simpan view baru ke database
		if errAdd := s.repo.AddView(&newView); errAdd == nil {
			// Ambil ulang data agar ViewsCount terbaru langsung dikirim ke user
			updatedPost, errReload := s.repo.FindBySlugOrID(id)
			if errReload == nil {
				post = updatedPost
			}
		}
	}

	return responses.FromDomainToPostResponse(post), nil
}

// 4. UPDATE POST
func (s *postService) UpdatePost(id string, req requests.PostUpdateRequest) (responses.PostResponse, error) {
	// Cari data eksisting
	post, err := s.repo.FindBySlugOrID(id)
	if err != nil {
		return responses.PostResponse{}, err
	}

	// Update field jika dikirim
	if req.Title != "" {
		post.Title = req.Title
		post.Slug = strings.ToLower(strings.ReplaceAll(req.Title, " ", "-"))
	}

	if req.Content != "" {
		post.Content = req.Content
		excerpt := req.Content
		if len(excerpt) > 150 {
			excerpt = excerpt[:150] + "..."
		}
		post.Excerpt = &excerpt
	}

	if req.CategoryID != 0 {
		post.CategoryID = req.CategoryID
	}

	// Update Tags (Replace Association)
	if req.Tags != "" {
		post.Tags = s.processTags(req.Tags)
	}

	// Update Gambar jika ada file baru
	if req.Image != nil {
		file, _ := req.Image.Open()
		defer file.Close()
		filename := "post-update-" + time.Now().Format("20060102-150405")
		url, err := utils.UploadToCloudinary(file, filename)
		if err == nil {
			post.FeaturedImage = &url
		}
	}

	if err := s.repo.Update(&post); err != nil {
		return responses.PostResponse{}, err
	}

	// Reload untuk response DTO terbaru
	updatedPost, _ := s.repo.FindByID(post.ID)
	return responses.FromDomainToPostResponse(updatedPost), nil
}

// 5. DELETE POST
func (s *postService) DeletePost(id string) error {
	post, err := s.repo.FindBySlugOrID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(&post, false) // Soft delete sesuai domain.go
}

// HELPER: Proses String Tags menjadi Domain Entities (FirstOrCreate)
func (s *postService) processTags(tagsInput string) []domain.Tag {
	var tags []domain.Tag
	if tagsInput == "" {
		return tags
	}

	tagNameList := strings.Split(tagsInput, ",")
	for _, tagName := range tagNameList {
		tagName = strings.TrimSpace(tagName)
		if tagName == "" {
			continue
		}

		tagSlug := strings.ToLower(strings.ReplaceAll(tagName, " ", "-"))
		tag, err := s.repo.GetTagBySlug(tagSlug, tagName)
		if err == nil {
			tags = append(tags, tag)
		}
	}
	return tags
}
