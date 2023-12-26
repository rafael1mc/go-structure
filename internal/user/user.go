package user

import (
	"database/sql"
	"fmt"
	"gomodel/internal/shared/database"
	"gomodel/internal/shared/env"
	"gomodel/internal/shared/util/password"
	"gomodel/internal/shared/util/rnd"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	database *database.Database
	env      *env.Env
}

func NewUser(
	database *database.Database,
	env *env.Env,
) *User {
	return &User{
		database: database,
		env:      env,
	}
}

func (u User) Create(createReq CreateUserRequest, tx *sqlx.Tx) (dbUser UserModel, err error) {
	myTx, err := u.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	if createReq.Name == "" {
		err = InvalidNameError
		return
	}
	if createReq.Email == "" {
		err = InvalidEmailError
		return
	}
	if createReq.PlainPass == "" || len(createReq.PlainPass) < 8 {
		err = InvalidPasswordError
		return
	}

	passSalt := uuid.NewString()
	passHash, err := password.HashAndSalt(createReq.PlainPass, passSalt)
	if err != nil {
		return
	}

	dbUser = UserModel{
		ID:        uuid.NewString(),
		Name:      createReq.Name,
		Email:     createReq.Email,
		PassHash:  passHash,
		PassSalt:  passSalt,
		IsEnabled: true,
		CreatedAt: nil,
		UpdatedAt: nil,
	}

	_, err = tx.NamedExec(
		`INSERT INTO user_ (id, name_, email, pass_hash, pass_salt, is_enabled)
			VALUES(:id, :name_, :email, :pass_hash, :pass_salt, :is_enabled);`,
		dbUser,
	)

	if myTx.IsExternal() {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	// TODO SEND EMAIl
	fmt.Println("Would be sending email now from user.Create")

	return
}

func (u User) Invite(
	inviteReq InviteUserRequest,
	tx *sqlx.Tx,
) (err error) {
	myTx, err := u.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	if inviteReq.Email == "" {
		err = InvalidEmailError
		return
	}
	if inviteReq.InviterID == "" {
		err = InvalidInviterIDError
		return
	}

	user, err := u.GetByEmail(inviteReq.Email, tx)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	invitation, err := u.GetInviteByEmail(inviteReq.Email, tx)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	invitationCode := rnd.RandStringBytesMaskImprSrcUnsafe(u.env.Api.InviteCodeLength)

	if err == sql.ErrNoRows { // no previous existing invite for this email
		invitationID := uuid.NewString()
		_, err = tx.Exec(
			`INSERT INTO user_invite (id,  user_id, email, invitation_status, code)
				VALUES ($1, $2, $3, $4, $5)`,
			invitationID,
			inviteReq.InviterID,
			inviteReq.Email,
			PendingInvitationStatus,
			invitationCode,
		)
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

		// TODO SEND EMAIL

		return
	}

	if invitation.Status == CompletedInvitationStatus && user.IsEnabled {
		// silently return nothing, because invited email is already an enabled user
		err = nil
		return
	}

	_, err = tx.Exec(
		`UPDATE user_invite
			SET
				updated_at = NOW(),
				invitation_status = $1,
				code = $2
			WHERE email = $3`,
		PendingInvitationStatus,
		invitationCode,
		inviteReq.Email,
	)
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

	// TODO SEND EMAIL

	return
}

func (u User) GetInviteByEmail(
	email string,
	tx *sqlx.Tx,
) (dbModel UserInviteModel, err error) {
	myTx, err := u.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	err = tx.Get(
		&dbModel,
		`SELECT * FROM user_invite WHERE email = $1`,
		email,
	)

	return
}

func (u User) GetByID(
	id string,
	tx *sqlx.Tx,
) (dbUser UserModel, err error) {
	myTx, err := u.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	err = tx.Get(
		&dbUser,
		`SELECT * FROM user_ WHERE id = $1`,
		id,
	)

	return
}

func (u User) GetByEmail(
	email string,
	tx *sqlx.Tx) (dbUser UserModel, err error) {
	myTx, err := u.database.WithTx(&tx)
	if err != nil {
		return
	}
	defer myTx.HandleRollback(&err)

	err = tx.Get(
		&dbUser,
		`SELECT * FROM user_ WHERE email = $1`,
		email,
	)

	return
}
