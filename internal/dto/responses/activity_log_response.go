package responses

import (
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
)

// ActivityLogUserInfo menampilkan info user yang melakukan aktivitas
type ActivityLogUserInfo struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type ActivityLogResponse struct {
	ID          int                       `json:"id"`
	UserID      int                       `json:"user_id"`
	User        *ActivityLogUserInfo      `json:"user,omitempty"`
	ActionType  domain.ActivityActionType `json:"action_type"`
	Module      domain.ActivityModuleType `json:"module"`
	Description *string                   `json:"description,omitempty"`
	TargetID    *int                      `json:"target_id,omitempty"`
	OldValue    map[string]any            `json:"old_value,omitempty"`
	NewValue    map[string]any            `json:"new_value,omitempty"`
	IPAddress   *string                   `json:"ip_address,omitempty"`
	UserAgent   *string                   `json:"user_agent,omitempty"`
	CreatedAt   time.Time                 `json:"created_at"`
}
