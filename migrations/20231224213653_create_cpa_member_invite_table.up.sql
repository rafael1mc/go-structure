CREATE TABLE IF NOT EXISTS member_invite (
    id UUID PRIMARY KEY,
    email VARCHAR NOT NULL,
    invitation_status VARCHAR NOT NULL,
    category VARCHAR NOT NULL,
    code VARCHAR NOT NULL,
    institution_id UUID NOT NULL REFERENCES institution(id),
    inviter_id UUID NOT NULL REFERENCES user_(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON member_invite
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
