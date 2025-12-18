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

type AboutUsSectionResponse struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	ImageURI    string `json:"image_uri"`
}

type WhySectionResponse struct {
	Title       string              `json:"title"`
	Subtitle    string              `json:"subtitle"`
	Description *string             `json:"description"`
	Data        []map[string]string `json:"data"`
}

type TestimonialSectionResponse struct {
	Testimoni string `json:"testimoni"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Career    string `json:"career"`
	ImageURI  string `json:"image_uri"`
}
