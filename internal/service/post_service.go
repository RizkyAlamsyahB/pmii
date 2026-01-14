package service

import (
	"context"
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
	CreatePost(ctx context.Context, req requests.PostCreateRequest) (responses.PostResponse, error)
	GetPostDetail(id string) (responses.PostResponse, error)
	UpdatePost(ctx context.Context, id string, req requests.PostUpdateRequest) (responses.PostResponse, error)
	DeletePost(ctx context.Context, id string) error
}

type postService struct {
	repo            repository.PostRepository
	activityLogRepo repository.ActivityLogRepository
}

func NewPostService(repo repository.PostRepository, activityLogRepo repository.ActivityLogRepository) PostService {
	return &postService{
		repo:            repo,
		activityLogRepo: activityLogRepo,
	}
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
func (s *postService) CreatePost(ctx context.Context, req requests.PostCreateRequest) (responses.PostResponse, error) {
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

	// Get user ID from context
	userID, _ := utils.GetUserID(ctx)
	if userID == 0 {
		userID = 1 // Default Admin if not in context
	}

	publishedTime := time.Now()
	post := domain.Post{
		Title:         req.Title,
		Content:       req.Content,
		Slug:          req.GetSlug(),
		CategoryID:    req.CategoryID,
		UserID:        userID,
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

	// Log activity - Create Post
	s.logActivity(ctx, domain.ActionCreate, domain.ModulePost, "Membuat post baru: "+post.Title, nil, map[string]any{
		"id":          post.ID,
		"title":       post.Title,
		"slug":        post.Slug,
		"category_id": post.CategoryID,
	}, &post.ID)

	return responses.FromDomainToPostResponse(updatedPost), nil
}

// 3. GET DETAIL POST
func (s *postService) GetPostDetail(id string) (responses.PostResponse, error) {
	post, err := s.repo.FindBySlugOrID(id)
	if err != nil {
		return responses.PostResponse{}, err
	}
	return responses.FromDomainToPostResponse(post), nil
}

// 4. UPDATE POST
func (s *postService) UpdatePost(ctx context.Context, id string, req requests.PostUpdateRequest) (responses.PostResponse, error) {
	// Cari data eksisting
	post, err := s.repo.FindBySlugOrID(id)
	if err != nil {
		return responses.PostResponse{}, err
	}

	// Store old values for audit
	oldValues := map[string]any{
		"title":       post.Title,
		"slug":        post.Slug,
		"category_id": post.CategoryID,
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

	// Log activity - Update Post
	s.logActivity(ctx, domain.ActionUpdate, domain.ModulePost, "Mengupdate post: "+post.Title, oldValues, map[string]any{
		"id":          post.ID,
		"title":       post.Title,
		"slug":        post.Slug,
		"category_id": post.CategoryID,
	}, &post.ID)

	return responses.FromDomainToPostResponse(updatedPost), nil
}

// 5. DELETE POST
func (s *postService) DeletePost(ctx context.Context, id string) error {
	post, err := s.repo.FindBySlugOrID(id)
	if err != nil {
		return err
	}

	// Log activity sebelum delete
	s.logActivity(ctx, domain.ActionDelete, domain.ModulePost, "Menghapus post: "+post.Title, map[string]any{
		"id":          post.ID,
		"title":       post.Title,
		"slug":        post.Slug,
		"category_id": post.CategoryID,
	}, nil, &post.ID)

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

// logActivity helper untuk mencatat activity log
func (s *postService) logActivity(ctx context.Context, actionType domain.ActivityActionType, module domain.ActivityModuleType, description string, oldValue, newValue map[string]any, targetID *int) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		return // Skip if no user in context
	}

	ipAddress := utils.GetIPAddress(ctx)
	userAgent := utils.GetUserAgent(ctx)

	var ipPtr, uaPtr *string
	if ipAddress != "" {
		ipPtr = &ipAddress
	}
	if userAgent != "" {
		uaPtr = &userAgent
	}

	log := &domain.ActivityLog{
		UserID:      userID,
		ActionType:  actionType,
		Module:      module,
		Description: &description,
		TargetID:    targetID,
		OldValue:    oldValue,
		NewValue:    newValue,
		IPAddress:   ipPtr,
		UserAgent:   uaPtr,
	}

	// Ignore error - logging should not affect main operation
	_ = s.activityLogRepo.Create(log)
}
