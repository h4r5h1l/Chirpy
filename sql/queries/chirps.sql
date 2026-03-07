-- name: CreateChirp :one
INSERT INTO chirps (
        id,
        created_at,
        updated_at,
        body,
        user_id,
        is_chirpy_red
    )
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2, false)
RETURNING *;
-- name: GetAllChirps :many
SELECT *
FROM chirps
ORDER BY created_at ASC;
-- name: GetChirpByChirpID :one
SELECT *
FROM chirps
WHERE id = $1;
-- name: DeleteChirpByChirpID :exec
DELETE FROM chirps
WHERE id = $1;