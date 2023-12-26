CREATE TABLE IF NOT EXISTS member (
    id UUID PRIMARY KEY,
    category VARCHAR NOT NULL,
    institution_id UUID REFERENCES institution(id),
    user_id UUID REFERENCES user_(id),
    is_enabled BOOL NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX unique_member_enabled_institution_user ON member (institution_id, user_id) WHERE is_enabled = true;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON member
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();