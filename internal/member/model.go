package member

import "time"

type MemberResponse struct {
	ID            string `json:"id" db:"id"`
	Name          string `json:"name" db:"name_"`
	InstitutionID string `json:"institution_id" db:"institution_id"`
	UserID        string `json:"user_id" db:"user_id"`
}

const (
	PendingInvitationStatus   = "PENDING"
	CompletedInvitationStatus = "COMPLETED"
)

type InviteMemberRequest struct {
	InstitutionID string `json:"-"`
	Email         string `json:"email"`
	Category      string `json:"category"`
	InviterID     string `json:"-"`
}

type MemberInviteModel struct {
	ID            string     `db:"id"`
	Email         string     `db:"email"`
	Status        string     `db:"invitation_status"`
	Category      string     `db:"category"`
	Code          string     `db:"code"`
	InviterID     string     `db:"inviter_id"`
	InstitutionID string     `db:"institution_id"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

type ValidateInviteRequest struct {
	InstitutionID string `json:"institution_id"`
	Code          string `json:"code"`
}

type AcceptInviteRequest struct {
	ValidateInviteRequest
	Password string `json:"password"`
	Name     string `json:"name"`
}

type CreateMemberRequest struct {
	Category      string
	InstitutionID string
	UserID        string
}

type MemberModel struct {
	ID            string     `db:"id"`
	Category      string     `db:"category"`
	InstitutionID string     `db:"institution_id"`
	UserID        string     `db:"user_id"`
	IsEnabled     bool       `db:"is_enabled"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

// type UpdateMemberRequest struct {
// 	ID        string `json:"id" db:"id"`
// 	Category  string `json:"category" db:"category"`
// 	IsEnabled bool   `json:"is_enabled" db:"is_enabled"`
// }
