package member

import (
	"database/sql"
	"fmt"
	"gomodel/internal/shared/database"
	"gomodel/internal/shared/env"
	"gomodel/internal/shared/timeprovider"
	"gomodel/internal/shared/util/rnd"
	"gomodel/internal/user"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Member struct {
	database     *database.Database
	logger       *slog.Logger
	user         *user.User
	timeProvider timeprovider.TimeProvider
	env          *env.Env
}

func NewMember(
	database *database.Database,
	logger *slog.Logger,
	user *user.User,
	timeProvider timeprovider.TimeProvider,
	env *env.Env,
) *Member {
	return &Member{
		database:     database,
		logger:       logger,
		user:         user,
		timeProvider: timeProvider,
		env:          env,
	}
}

func (m *Member) GetMembersByInstitutionID(
	institutionID string,
	tx *sqlx.Tx,
) (res []MemberResponse, err error) {
	myTx, err := m.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	err = tx.Select(
		&res,
		`SELECT c.*, u.name_ FROM member m
		INNER JOIN user_ u ON u.id = m.user_id
		WHERE m.institution_id = $1
		AND m.is_enabled = TRUE
		AND u.is_enabled = TRUE`,
		institutionID,
	)

	return
}

func (m *Member) InviteMember(
	request InviteMemberRequest,
	tx *sqlx.Tx,
) (err error) {
	myTx, err := m.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	if request.InstitutionID == "" {
		err = InvalidInstitutionError
		return
	}
	if request.Email == "" {
		err = InvalidEmailError
		return
	}
	if request.InviterID == "" {
		err = InvalidInviterIDError
		return
	}
	if request.Category == "" {
		err = InvalidCategoryError
		return
	}

	user, err := m.user.GetByEmail(request.Email, tx)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	invitation, err := m.GetInviteByEmail(request.InstitutionID, request.Email, tx)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	invitationCode := rnd.RandStringBytesMaskImprSrcUnsafe(m.env.Api.InviteCodeLength)

	if err == sql.ErrNoRows { // no previous existing invite for this email
		invitationID := uuid.NewString()
		_, err = tx.Exec(
			`INSERT INTO member_invite (id, email, invitation_status, category, code, institution_id, inviter_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			invitationID,
			request.Email,
			PendingInvitationStatus,
			request.Category,
			invitationCode,
			request.InstitutionID,
			request.InviterID,
		)

		if myTx.IsExternal() {
			return
		}

		err = tx.Commit()
		if err != nil {
			return
		}

		// TODO send email

		return
	}

	if invitation.Status == CompletedInvitationStatus && user.IsEnabled { // already accepted
		err = nil
		return
	}

	// TODO resent email

	return nil
}

func (m Member) ValidateInvite(
	request ValidateInviteRequest,
	tx *sqlx.Tx,
) (invite MemberInviteModel, err error) {
	myTx, err := m.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	invite, err = m.GetInviteByCode(request.InstitutionID, request.Code, tx)
	if err != nil {
		return
	}

	if invite.Status != PendingInvitationStatus {
		err = InvalidInviteStatusError
		return
	}

	currentInviteDuration := m.timeProvider.ProvideUTCMilli() - invite.CreatedAt.UnixMilli() + (5 * time.Minute).Milliseconds()
	if time.Duration(currentInviteDuration).Seconds() > float64(m.env.Api.InviteCodeTTLSeconds) {
		err = ExpiredInviteCodeError
		return
	}

	err = nil
	return
}

func (m Member) AcceptInvite(request AcceptInviteRequest, tx *sqlx.Tx) (err error) {
	myTx, err := m.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	invite, err := m.ValidateInvite(request.ValidateInviteRequest, tx)
	if err != nil {
		return
	}

	if request.Name == "" {
		err = InvalidNameError
		return
	}
	if request.Password == "" || len(request.Password) < 8 {
		err = InvalidPasswordError
		return
	}

	userModel, err := m.user.GetByEmail(invite.Email, tx)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		userReq := user.CreateUserRequest{
			Email:     invite.Email,
			PlainPass: request.Password,
			Name:      request.Name,
		}
		userModel, err = m.user.Create(userReq, tx)
		if err != nil {
			return
		}
	}

	memberReq := CreateMemberRequest{
		Category:      invite.Category,
		InstitutionID: invite.InstitutionID,
		UserID:        userModel.ID,
	}
	_, err = m.CreateMember(memberReq, tx)
	if err != nil {
		return
	}

	err = m.CompleteInvite(invite.ID, tx)
	if err != nil {
		return
	}

	if myTx.IsExternal() {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	// TODO send welcome email
	fmt.Println("Would be sending email from member accept invite")

	return
}

func (m Member) CreateMember(request CreateMemberRequest, tx *sqlx.Tx) (dbModel MemberModel, err error) {
	myTx, err := m.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	if request.Category == "" {
		err = InvalidCategoryError
		return
	}
	if request.InstitutionID == "" {
		err = InvalidInstitutionError
		return
	}
	if request.UserID == "" {
		err = InvalidUserError
		return
	}

	dbModel = MemberModel{
		ID:            uuid.NewString(),
		Category:      request.Category,
		InstitutionID: request.InstitutionID,
		UserID:        request.UserID,
		IsEnabled:     true,
		CreatedAt:     nil,
		UpdatedAt:     nil,
	}

	_, err = tx.NamedExec(
		`INSERT INTO member (id, category, institution_id, user_id, is_enabled)
		VALUES (:id, :category, :institution_id, :user_id, :is_enabled)`,
		dbModel,
	)

	if myTx.IsExternal() {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

func (m Member) CompleteInvite(inviteID string, tx *sqlx.Tx) (err error) {
	myTx, err := m.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	if inviteID == "" {
		err = InvalidInviterIDError
		return
	}

	_, err = tx.Exec(
		`UPDATE member_invite
		SET invitation_status = $1
		WHERE id = $2`,
		CompletedInvitationStatus,
		inviteID,
	)

	if myTx.IsExternal() {
		return
	}

	err = tx.Commit()
	return
}

func (m Member) GetInviteByEmail(
	institutionID string,
	email string,
	tx *sqlx.Tx,
) (dbModel MemberInviteModel, err error) {
	myTx, err := m.database.WithTx(&tx)
	if err != nil {
		return
	}
	myTx.HandleRollback(&err)

	err = tx.Get(
		&dbModel,
		`SELECT * FROM member_invite WHERE institution_id = $1 AND email = $2;`,
		institutionID,
		email,
	)

	return
}

func (m Member) GetInviteByCode(
	institutionID,
	code string,
	tx *sqlx.Tx,
) (dbModel MemberInviteModel, err error) {
	myTx, err := m.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	err = tx.Get(
		&dbModel,
		`SELECT * FROM member_invite
		WHERE institution_id = $1
		AND code = $2`,
		institutionID,
		code,
	)

	return
}
