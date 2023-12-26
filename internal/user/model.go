package user

import "time"

type InviteUserRequest struct {
	Email     string `json:"email"`
	InviterID string `json:"-"`
}

type CreateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	PlainPass string `json:"password"`
}

type UserModel struct {
	ID        string     `db:"id"`
	Name      string     `db:"name_"`
	Email     string     `db:"email"`
	PassHash  string     `db:"pass_hash"`
	PassSalt  string     `db:"pass_salt"`
	IsEnabled bool       `db:"is_enabled"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

const (
	PendingInvitationStatus   = "PENDING"
	CompletedInvitationStatus = "COMPLETED"
)

type UserInviteModel struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	Email     string     `db:"email"`
	Status    string     `db:"invitation_status"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
