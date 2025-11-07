-- name: CreateMessage :one
INSERT INTO messages (
    chat_id, sender_id, content
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetMessages :many
SELECT * FROM messages
WHERE chat_id = $1
ORDER BY created_at ASC;

-- name: EditMessage :exec
UPDATE messages
SET
    content = $2,
    is_edited = TRUE,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1
RETURNING id;;

-- name: CreateChat :one
INSERT INTO chats (
    is_group
) VALUES (
    $1
)
RETURNING *;

-- name: DeleteChat :exec
DELETE FROM chats
WHERE id = $1;

-- name: DeleteMessages :exec
DELETE FROM messages
WHERE chat_id = $1;

-- name: CreateGroupChat :one
INSERT INTO chats (
    name, is_group
) VALUES (
    $1, TRUE
)
RETURNING *;

-- name: AddUserToChat :exec
INSERT INTO chat_members (
    chat_id, user_id
) VALUES (
    $1, $2
)
ON CONFLICT DO NOTHING;

-- name: RemoveUserFromChat :exec
DELETE FROM chat_members
WHERE chat_id = $1 AND user_id = $2;

-- name: GetChatMembers :many
SELECT user_id FROM chat_members
WHERE chat_id = $1;

-- name: GetUserChats :many
SELECT c.*
FROM chats c
JOIN chat_members cm ON c.id = cm.chat_id
WHERE cm.user_id = $1;

-- name: ChatExists :one
SELECT EXISTS (
    SELECT 1 FROM chats WHERE id = $1
);

-- name: MessageExists :one
SELECT EXISTS (
    SELECT 1 FROM messages WHERE id = $1
);

-- name: IsUserInChat :one
SELECT EXISTS (
    SELECT 1 FROM chat_members WHERE chat_id = $1 AND user_id = $2
);