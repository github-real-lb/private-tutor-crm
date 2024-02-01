-- name: CreateStudent :one
INSERT INTO students (
  first_name, last_name, phone_number, email_address, college_id, funnel_id, notes
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetStudent :one
SELECT * FROM students
WHERE student_id = $1 LIMIT 1;

-- name: ListStudents :many
SELECT * FROM students
ORDER BY last_name, first_name
LIMIT $1
OFFSET $2;