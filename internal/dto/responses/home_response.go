package responses

import "time"

type HeroCategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type HeroSectionResponse struct {
	ID          int64                `json:"id"`
	Title       string               `json:"title"`
	Slug        string               `json:"slug"`
	Excerpt     string               `json:"excerpt"`
	ImageURL    string               `json:"imageUrl"`
	PublishedAt string               `json:"publishedAt"`
	Category    HeroCategoryResponse `json:"category"`
	AuthorID    int64                `json:"authorId"`
	Tags        []string             `json:"tags"`
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

type WhyItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IconURI     string `json:"icon_uri"`
}

type WhySectionResponse struct {
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Description *string   `json:"description"`
	Data        []WhyItem `json:"data"`
}

type TestimonialSectionResponse struct {
	Testimoni string `json:"testimoni"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Career    string `json:"career"`
	ImageURI  string `json:"image_uri"`
}

type FaqItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type FaqSectionResponse struct {
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Description *string   `json:"description"`
	Data        []FaqItem `json:"data"`
}

type CtaSectionResponse struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}
