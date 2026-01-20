package domain

import "time"

// ActivityActionType represents the type of action performed
type ActivityActionType string

const (
	ActionCreate ActivityActionType = "create"
	ActionUpdate ActivityActionType = "update"
	ActionDelete ActivityActionType = "delete"
	ActionLogin  ActivityActionType = "login"
	ActionLogout ActivityActionType = "logout"
)

// ActivityModuleType represents the module where the action was performed
type ActivityModuleType string

const (
	ModuleUser      ActivityModuleType = "user"
	ModulePost      ActivityModuleType = "post"
	ModuleCategory  ActivityModuleType = "category"
	ModuleTags      ActivityModuleType = "tags"
	ModuleTestimoni ActivityModuleType = "testimoni"
	ModuleMembers   ActivityModuleType = "members"
	ModuleTeams     ActivityModuleType = "teams"
	ModuleDokumen   ActivityModuleType = "dokumen"
	ModuleSettings  ActivityModuleType = "settings"
	ModuleAuth      ActivityModuleType = "auth"
	ModuleAds       ActivityModuleType = "ads"
)

// ActivityLog represents a log entry for user activities
type ActivityLog struct {
	ID          int                `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int                `gorm:"not null" json:"user_id"`
	ActionType  ActivityActionType `gorm:"type:activity_action_type;not null" json:"action_type"`
	Module      ActivityModuleType `gorm:"type:activity_module_type;not null" json:"module"`
	Description *string            `gorm:"type:text" json:"description,omitempty"`
	TargetID    *int               `json:"target_id,omitempty"`
	OldValue    map[string]any     `gorm:"type:jsonb;serializer:json" json:"old_value,omitempty"`
	NewValue    map[string]any     `gorm:"type:jsonb;serializer:json" json:"new_value,omitempty"`
	IPAddress   *string            `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent   *string            `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt   time.Time          `gorm:"default:now()" json:"created_at"`

	// Relationship
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for ActivityLog
func (ActivityLog) TableName() string {
	return "activity_logs"
}
