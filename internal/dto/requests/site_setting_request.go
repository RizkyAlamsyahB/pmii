package requests

// UpdateSiteSettingRequest adalah DTO untuk update site settings
type UpdateSiteSettingRequest struct {
	SiteName        *string `form:"site_name"`
	SiteTitle       *string `form:"site_title"`
	SiteDescription *string `form:"site_description"`
	FacebookURL     *string `form:"facebook_url"`
	TwitterURL      *string `form:"twitter_url"`
	LinkedinURL     *string `form:"linkedin_url"`
	InstagramURL    *string `form:"instagram_url"`
	YoutubeURL      *string `form:"youtube_url"`
	GithubURL       *string `form:"github_url"`
	// Images handled separately: favicon, logo_header, logo_big
}
