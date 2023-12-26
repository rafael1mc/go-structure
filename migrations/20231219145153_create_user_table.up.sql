CREATE TABLE IF NOT EXISTS user_ (
    id UUID PRIMARY KEY,
    email VARCHAR NOT NULL,
    pass_hash VARCHAR NOT NULL,
    pass_salt VARCHAR NOT NULL,
    name_ VARCHAR NOT NULL,
    is_enabled BOOL NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX unique_enabled_email ON user_ (email) WHERE is_enabled = true;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON user_
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
