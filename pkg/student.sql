-- name: GetStudent :one
SELECT * FROM students
WHERE student_id = $1 LIMIT 1;

-- name: ListStudents :many
SELECT * FROM students
ORDER BY last_name, first_name;

-- name: CreateStudent :one
INSERT INTO students (
  first_name, last_name, phone_number, email_address, college_id, funnel_id, notes
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateStudent :exec
UPDATE students
  set   first_name = $2,
        last_name = $2, 
        phone_number = $2, 
        email_address = $2, 
        college_id = $2, 
        funnel_id = $2, 
        notes = $2
WHERE student_id = $1;

-- name: DeleteStudent :exec
DELETE FROM students
WHERE student_id = $1;