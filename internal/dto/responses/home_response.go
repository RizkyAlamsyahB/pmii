package responses

import "time"

type HeroSectionResponse struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	FeaturedImage string    `json:"featured_image"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	TotalViews    int64     `json:"total_views"`
}

type LatestNewsSectionResponse struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	FeaturedImage string    `json:"featured_image"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	TotalViews    int64     `json:"total_views"`
}
