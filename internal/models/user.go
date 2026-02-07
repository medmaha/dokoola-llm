package models

// AuthUser represents an authenticated user from the backend
type AuthUser struct {
	Name            string `json:"name"`
	PublicID        string `json:"public_id"`
	IsActive        bool   `json:"is_active"`
	CompleteProfile bool   `json:"complete_profile"`
	Avatar          string `json:"avatar"`
	IsStaff         bool   `json:"is_staff"`
	IsTalent        bool   `json:"is_talent"`
	IsClient        bool   `json:"is_client"`
}
