package domain

import "time"

// MemberDepartment represents department type for members
type MemberDepartment string

const (
	DepartmentPengurusHarian MemberDepartment = "pengurus_harian"
	DepartmentKabid          MemberDepartment = "kabid"
	DepartmentWasekbid       MemberDepartment = "wasekbid"
	DepartmentWakilBendahara MemberDepartment = "wakil_bendahara"
)

// GetDepartmentLabel returns human-readable label for department
func (d MemberDepartment) GetLabel() string {
	switch d {
	case DepartmentPengurusHarian:
		return "Pengurus Harian"
	case DepartmentKabid:
		return "Ketua Bidang (Kabid)"
	case DepartmentWasekbid:
		return "Wakil Sekretaris Bidang (Wasekbid)"
	case DepartmentWakilBendahara:
		return "Wakil Bendahara"
	default:
		return string(d)
	}
}

// ValidDepartments returns all valid department values
func ValidDepartments() []MemberDepartment {
	return []MemberDepartment{
		DepartmentPengurusHarian,
		DepartmentKabid,
		DepartmentWasekbid,
		DepartmentWakilBendahara,
	}
}

// Member represents an organization member (e.g., team member, staff)
type Member struct {
	ID          int              `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName    string           `gorm:"type:varchar(100);not null" json:"full_name"`
	Position    string           `gorm:"type:varchar(100);not null" json:"position"`
	Department  MemberDepartment `gorm:"type:member_department;not null;default:'kabid'" json:"department"`
	PhotoURI    *string          `gorm:"type:varchar(255)" json:"photo_uri,omitempty"`
	SocialLinks map[string]any   `gorm:"type:jsonb;serializer:json" json:"social_links,omitempty"`
	IsActive    bool             `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time        `gorm:"default:now()" json:"created_at"`
}

// TableName specifies the table name for Member
func (Member) TableName() string {
	return "members"
}
