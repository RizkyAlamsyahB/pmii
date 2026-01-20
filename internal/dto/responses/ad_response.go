package responses

import "github.com/garuda-labs-1/pmii-be/internal/domain"

// AdResponse untuk response single ad (admin) - hanya slotName dan imageUrl
type AdResponse struct {
	SlotName string  `json:"slotName"`
	ImageURL *string `json:"imageUrl"`
}

// AdsByPageResponse untuk response ads grouped by page (admin)
type AdsByPageResponse struct {
	Page     string       `json:"page"`
	PageName string       `json:"pageName"`
	Ads      []AdResponse `json:"ads"`
}

// AllAdsResponse untuk response semua ads (untuk admin)
type AllAdsResponse struct {
	Pages []AdsByPageResponse `json:"pages"`
}

// PublicAdResponse untuk response public - hanya slot dan imageUrl
type PublicAdResponse struct {
	Slot     int     `json:"slot"`
	ImageURL *string `json:"imageUrl"`
}

// PublicAdsByPageResponse untuk response public ads by page
type PublicAdsByPageResponse struct {
	Page string             `json:"page"`
	Ads  []PublicAdResponse `json:"ads"`
}

// ToAdResponse converts domain.Ad to AdResponse
func ToAdResponse(ad *domain.Ad) AdResponse {
	return AdResponse{
		SlotName: ad.GetSlotName(),
		ImageURL: ad.ImageURL,
	}
}

// ToAdResponseList converts slice of domain.Ad to slice of AdResponse
func ToAdResponseList(ads []domain.Ad) []AdResponse {
	result := make([]AdResponse, len(ads))
	for i, ad := range ads {
		result[i] = ToAdResponse(&ad)
	}
	return result
}

// ToPublicAdResponse converts domain.Ad to PublicAdResponse
func ToPublicAdResponse(ad *domain.Ad) PublicAdResponse {
	return PublicAdResponse{
		Slot:     ad.Slot,
		ImageURL: ad.ImageURL,
	}
}

// ToPublicAdResponseList converts slice of domain.Ad to slice of PublicAdResponse
func ToPublicAdResponseList(ads []domain.Ad) []PublicAdResponse {
	result := make([]PublicAdResponse, len(ads))
	for i, ad := range ads {
		result[i] = ToPublicAdResponse(&ad)
	}
	return result
}

// GroupAdsByPage groups ads by page for admin response
func GroupAdsByPage(ads []domain.Ad) []AdsByPageResponse {
	// Define page order
	pageOrder := []domain.AdPage{
		domain.AdPageLanding,
		domain.AdPageNews,
		domain.AdPageOpini,
		domain.AdPageLifeAtPMII,
		domain.AdPageIslamic,
		domain.AdPageDetailArticle,
	}

	// Group ads by page
	pageMap := make(map[domain.AdPage][]AdResponse)
	pageNameMap := make(map[domain.AdPage]string)
	for _, ad := range ads {
		pageMap[ad.Page] = append(pageMap[ad.Page], ToAdResponse(&ad))
		if _, exists := pageNameMap[ad.Page]; !exists {
			pageNameMap[ad.Page] = ad.GetPageDisplayName()
		}
	}

	// Build ordered result
	result := make([]AdsByPageResponse, 0, len(pageOrder))
	for _, page := range pageOrder {
		if adsInPage, exists := pageMap[page]; exists {
			result = append(result, AdsByPageResponse{
				Page:     string(page),
				PageName: pageNameMap[page],
				Ads:      adsInPage,
			})
		}
	}

	return result
}
