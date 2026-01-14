package responses

// PublicSiteSettingResponse adalah response site settings untuk public
// Tidak menyertakan updatedAt karena tidak diperlukan FE
type PublicSiteSettingResponse struct {
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
