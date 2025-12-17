package responses

type HeroSectionResponse struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	FeaturedImage string `json:"featured_image"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	TotalViews    int64  `json:"total_views"`
}
