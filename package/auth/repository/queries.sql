-- name: InsertUser :one
INSERT INTO users (
  email, password, full_name, verified_email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2
WHERE id = $1
RETURNING *;

-- name: UpdateVerifiedEmail :one
UPDATE users
SET verified_email = TRUE
FROM user_otps 
WHERE users.id = user_otps.user_id AND user_otps.otp = $1 AND users.email = $2
RETURNING users.id;

-- name: InsertOtp :one
INSERT INTO user_otps (
  user_id, otp
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetUserOtpByOtp :one
SELECT * 
FROM user_otps 
WHERE otp = $1 
LIMIT 1;

-- name: DeleteOtp :exec
DELETE FROM user_otps WHERE otp = $1;

-- name: InsertUserDetail :one
INSERT INTO user_details (
  user_id, phone_number, gender
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: InsertUserAvatar :exec
UPDATE users
SET avatar_url = $1
WHERE email = $2
RETURNING *;

-- name: InsertUserDetailAbout :exec
UPDATE user_details
SET about = $1
WHERE user_id = $2
RETURNING *;

-- name: InsertWorkExperience :one
INSERT INTO work_experiences (
  user_id, job_title, company_id, employment_type_id, location, location_type_id, start_date, finish_date, description
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: InsertEducation :one 
INSERT INTO educations (
  user_id, school_id, degree, field_of_study, gpa, start_date, finish_date
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: InsertCompany :one
INSERT INTO companies (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: InsertEmploymentType :one
INSERT INTO employment_types (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: InsertLocationType :one
INSERT INTO location_types (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: InsertSchool :one 
INSERT INTO schools (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: InsertCertificate :one
INSERT INTO certificates (
  user_id, name, issuing_organization_id, issue_date, expiration_date, credential_id, url
) VALUES (
   $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: InsertIssuingOrganization :one
INSERT INTO issuing_organizations (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: InsertUserSkill :one
INSERT INTO user_skills (
  user_id, skill_id, main_skill
) VALUES (
   $1, $2, $3
)
RETURNING *;

-- name: InsertSkill :one
INSERT INTO skills (
  name
) VALUES (
  $1
)
RETURNING *;