-- name: GetTaskByID :one
SELECT id, owner_id,
       group_id, title,
       description, ST_AsBinary(location)::geometry AS location,
       is_favorite, priority,
       icon, status,
       deadline, notify_at,
       created_at, updated_at,
       completed_at, purge_at
FROM pgtasks.task WHERE id = $1
LIMIT 1;

-- name: SaveTask :exec
UPDATE pgtasks.task
SET owner_id = $2, group_id = $3,
    title = $4, description = $5,
    location = $6, is_favorite = $7,
    priority = $8, icon = $9,
    status = $10, deadline = $11,
    notify_at = $12, completed_at = $13,
    created_at = $14, updated_at = $15,
    purge_at = $16, title_lang = $17,
    description_lang = $18
WHERE id = $1;

-- name: CreateTask :exec
INSERT INTO pgtasks.task (
                  id, owner_id,
                  group_id, title,
                  description, location,
                  is_favorite, priority,
                  icon, status,
                  deadline, notify_at,
                  title_lang, description_lang
)
VALUES (
        $1, $2,
        $3, $4,
        $5, $6::geometry,
        $7, $8,
        $9, $10,
        $11, $12,
        $13, $14
);

-- name: CompleteTask :exec
UPDATE pgtasks.task
SET status = $2,
    completed_at = $3,
    updated_at = $4
WHERE id = $1;

-- name: RestoreTask :exec
UPDATE pgtasks.task
SET updated_at = $2, purge_at = NULL
WHERE id = $1;

-- name: SoftDeleteTaskByID :exec
UPDATE pgtasks.task
SET purge_at = $2, updated_at = $3
WHERE id = $1;

-- name: DeleteTasksByOwnerID :exec
DELETE FROM pgtasks.task WHERE owner_id = $1;
