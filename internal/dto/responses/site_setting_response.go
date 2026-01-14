package responses

// SiteSettingResponse adalah DTO untuk response site settings
type SiteSettingResponse struct {
	SiteName        *string `json:"siteName,omitempty"`
	SiteTitle       *string `json:"siteTitle,omitempty"`
	SiteDescription *string `json:"siteDescription,omitempty"`
	Favicon         string  `json:"favicon,omitempty"`
	LogoHeader      string  `json:"logoHeader,omitempty"`
	LogoBig         string  `json:"logoBig,omitempty"`
	FacebookURL     *string `json:"facebookUrl,omitempty"`
	TwitterURL      *string `json:"twitterUrl,omitempty"`
	LinkedinURL     *string `json:"linkedinUrl,omitempty"`
	InstagramURL    *string `json:"instagramUrl,omitempty"`
	YoutubeURL      *string `json:"youtubeUrl,omitempty"`
	GithubURL       *string `json:"githubUrl,omitempty"`
}
