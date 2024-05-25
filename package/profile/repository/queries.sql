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
SET about = @about::text
WHERE user_id = @user_id::bigint;

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

-- name: BatchInsertSkills :exec
INSERT INTO skills (name)
SELECT unnest(@names::text[])
ON CONFLICT (name) DO NOTHING;

-- name: GetSkills :many
SELECT *, COUNT(id) OVER () AS total_rows
FROM skills
OFFSET $1
LIMIT $2;

-- name: UpdateUserMainSkillToFalse :exec
UPDATE user_skills 
SET main_skill = false 
WHERE user_id = @user_id::bigint
AND main_skill = true;

-- name: BatchInsertUserMainSkills :many
-- start get skills id 
WITH exist_skills AS (
    SELECT id, name
    FROM skills
    WHERE name = ANY(@names::text[])
)
-- end get skills id 
-- start insert user skills if not exist
INSERT INTO user_skills (user_id, skill_id, main_skill)
SELECT
    @user_id::bigint,
    es.id,
    @is_main_skill::boolean
FROM exist_skills es
WHERE es.name = ANY(@names::text[])
ON CONFLICT (user_id, skill_id) DO UPDATE
SET main_skill = true
RETURNING id;
-- end insert user skills if not exist

-- name: BatchInsertUserSkills :many
-- start get skills id 
WITH exist_skills AS (
    SELECT id, name
    FROM skills
    WHERE name = ANY(@names::text[])
)
-- end get skills id 
-- start insert user skills if not exist
INSERT INTO user_skills (user_id, skill_id, main_skill)
SELECT
    @user_id::bigint,
    es.id,
    @is_main_skill::boolean
FROM exist_skills es
WHERE es.name = ANY(@names::text[])
ON CONFLICT (user_id, skill_id) DO NOTHING
RETURNING id;
-- end insert user skills if not exist

-- name: UpdateUser :one
UPDATE users
SET full_name = $1,
    avatar_url = $2
WHERE id = $3
RETURNING full_name, avatar_url;

-- name: UpdateUserDetailByUserId :one
UPDATE user_details
SET hide_phone_number = $2,
    phone_number = $3,
    gender = $4
WHERE user_id = $1
RETURNING hide_phone_number, phone_number, gender;

-- name: UpsertUserSocialLink :exec
WITH social_link AS (
    SELECT id
    FROM social_links
    WHERE name = $2
    LIMIT 1
)
INSERT INTO user_social_links (user_id, social_link_id, url)
SELECT $1, sl.id, $3
FROM social_link sl
ON CONFLICT (user_id, social_link_id) DO UPDATE
SET url = EXCLUDED.url;

-- name: UpdateUserCertificate :one
UPDATE certificates 
SET name = @name::text,
    issuing_organization_id = @issuing_organization_id::bigint,
    issue_date = @issue_date::date, 
    expiration_date = @expiration_date::date, 
    credential_id = @credential_id::text, 
    url = @url::text
WHERE id = @id::bigint AND user_id = @user_id::bigint
RETURNING id; 

-- name: GetUserAvatarById :one
SELECT avatar_url
FROM users
WHERE users.id = $1
LIMIT 1;

-- name: GetUserDetail :one
SELECT * FROM user_details
WHERE user_id = @user_id::bigint
LIMIT 1;

-- name: UpdateUserDetail :one
UPDATE user_details
SET phone_number = $2,
    gender = $3,
    location = $4,
    portfolio_url = $5,
    about = $6,
    hide_phone_number = $7
WHERE user_id = $1
RETURNING id, phone_number, gender, location, portfolio_url, about, hide_phone_number;

-- name: UpdateUserEducation :one
UPDATE educations
SET school_id = $2,
    degree = $3,
    field_of_study = $4,
    gpa = $5,
    start_date = $6,
    finish_date = $7,
    description = $8,
    document_url = $9
WHERE id = $1
RETURNING *;

-- name: GetUserEducation :one
SELECT * FROM educations
WHERE id = @id::bigint
LIMIT 1;

-- name: BatchInsertEducationSkills :exec
INSERT INTO education_skills (education_id, user_skill_id)
SELECT @education_id::bigint, unnest(@user_skill_id::bigint[])
ON CONFLICT (education_id, user_skill_id) DO NOTHING;

-- name: DeleteEducationSkillsByEducation :many
DELETE FROM education_skills
WHERE education_id = @education_id::bigint
RETURNING user_skill_id;