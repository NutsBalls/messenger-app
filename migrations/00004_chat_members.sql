-- +goose Up
CREATE TABLE chat_members (
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    PRIMARY KEY (chat_id, user_id)
);

-- +goose Down
DROP TABLE chat_members;
