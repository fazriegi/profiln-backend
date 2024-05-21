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
WHERE id = $2
RETURNING *;

-- name: GetUserById :one
SELECT *
FROM users
WHERE users.id = $1
LIMIT 1;

-- name: UpdateUserDetailAbout :exec
UPDATE user_details
SET about = $1
WHERE user_id = $2
RETURNING *;

-- name: InsertUserDetailAbout :one
INSERT INTO user_details (
  user_id, about
) VALUES (
  $1, $2
)
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
INSERT INTO companies (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING *;

-- name: InsertEmploymentType :one
INSERT INTO employment_types (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING *;

-- name: InsertLocationType :one
INSERT INTO location_types (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING *;

-- name: InsertSchool :one 
INSERT INTO schools (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING *;

-- name: InsertCertificate :one
INSERT INTO certificates (
  user_id, name, issuing_organization_id, issue_date, expiration_date, credential_id, url
) VALUES (
   $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: InsertIssuingOrganization :one
INSERT INTO issuing_organizations (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING *;

-- name: InsertUserSkill :one
INSERT INTO user_skills (
  user_id, skill_id, main_skill
) VALUES (
   $1, $2, $3
)
RETURNING *;

-- name: InsertSkill :one
INSERT INTO skills (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING *;

-- name: GetUserAbout :one
SELECT users.id, user_details.about
FROM users
LEFT JOIN user_details
ON users.id = user_details.user_id
WHERE users.id = $1
LIMIT 1;

-- name: GetProfile :many
SELECT users.full_name, users.bio, user_social_links.url, social_links.name, user_skills.main_skill, skills.name, (
    SELECT COUNT(*) 
    FROM users 
    INNER JOIN followings 
    ON users.id = followings.user_id
    GROUP BY users.id
  ) AS count_following
FROM users
LEFT JOIN user_social_links
ON users.id = user_social_links.user_id
LEFT JOIN social_links
ON user_social_links.social_link_id = social_links.id
LEFT JOIN user_skills
ON users.id = user_skills.user_id
LEFT JOIN skills
ON user_skills.skill_id = skills.id
WHERE user_skills.main_skill = TRUE AND users.id = $1;