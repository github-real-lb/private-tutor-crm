-- name: CreateStudent :one
INSERT INTO students (
  first_name, last_name, email, phone_number, address, college_id, funnel_id, hourly_fee, notes 
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
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

-- name: UpdateStudent :exec
UPDATE students
  set   first_name = $2,
        last_name = $3, 
        email = $4,
        phone_number = $5, 
        address =  $6,
        college_id = $7,
        funnel_id = $8, 
        hourly_fee = $9, 
        notes = $10
WHERE student_id = $1;

-- name: DeleteStudent :exec
DELETE FROM students
WHERE student_id = $1;