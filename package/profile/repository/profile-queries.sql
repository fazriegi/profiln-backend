-- name: InsertUserAvatar :exec
UPDATE users
SET avatar_url = $1,
    updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: GetUserById :one
SELECT u.id, u.email, u.full_name, u.avatar_url, u.bio, u.open_to_work, u.followers_count, u.followings_count
FROM users u
WHERE u.id = $1
LIMIT 1;

-- name: UpdateUserDetailAbout :exec
UPDATE user_details
SET about = @about::text,
    updated_at = NOW()
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
  user_id, job_title, company_id, employment_type, location, location_type, start_date, finish_date, description, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()
)
RETURNING *;

-- name: InsertEducation :one 
INSERT INTO educations (
  user_id, school_id, degree, field_of_study, gpa, start_date, finish_date, description, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
)
RETURNING *;

-- name: InsertCompany :one
INSERT INTO companies (name)
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
  user_id, name, issuing_organization_id, issue_date, expiration_date, credential_id, url, created_at, updated_at
) VALUES (
   $1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
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
SELECT ud.id, ud.user_id, ud.about, ud.updated_at, ud.created_at, u.id, u.email, u.full_name
FROM users u
LEFT JOIN user_details ud
ON u.id = ud.user_id
WHERE u.id = $1
LIMIT 1;

-- name: GetProfile :many
SELECT users.full_name, users.bio, user_social_links.url, user_social_links.platform, user_skills.main_skill, skills.name, users.followers_count, users.followings_count
FROM users
LEFT JOIN user_social_links
ON users.id = user_social_links.user_id
LEFT JOIN user_skills
ON users.id = user_skills.user_id
LEFT JOIN skills
ON user_skills.skill_id = skills.id
WHERE user_skills.main_skill = TRUE AND users.id = $1;

-- name: BatchInsertSkills :exec
INSERT INTO skills (name)
SELECT unnest(@names::text[])
ON CONFLICT (name) DO NOTHING;

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
    avatar_url = $2,
    updated_at = NOW()
WHERE id = $3
RETURNING full_name, avatar_url;

-- name: UpdateUserDetailByUserId :one
UPDATE user_details
SET hide_phone_number = $2,
    phone_number = $3,
    gender = $4,
    updated_at = NOW()
WHERE user_id = $1
RETURNING hide_phone_number, phone_number, gender;

-- name: UpsertUserSocialLink :exec
INSERT INTO user_social_links (user_id, platform, url)
SELECT $1, $2, $3
ON CONFLICT (user_id, platform) DO UPDATE
SET url = EXCLUDED.url;

-- name: UpdateUserCertificate :one
UPDATE certificates 
SET name = @name::text,
    issuing_organization_id = @issuing_organization_id::bigint,
    issue_date = @issue_date::date, 
    expiration_date = @expiration_date::date, 
    credential_id = @credential_id::text, 
    url = @url::text,
    updated_at = NOW()
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
    hide_phone_number = $7,
    updated_at = NOW()
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
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetEducationById :one
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

-- name: BatchInsertEducationFiles :many
INSERT INTO education_files
  (education_id, url)
SELECT @education_id::bigint, UNNEST(@url::text[])
RETURNING *;

-- name: DeleteEducationFilesByEducationId :exec
DELETE FROM education_files
WHERE education_id = @education_id::bigint;

-- name: UpdateUserWorkExperience :one
UPDATE work_experiences
SET job_title = $2,
    company_id = $3,
    employment_type = $4,
    location = $5,
    location_type = $6,
    start_date = $7,
    finish_date = $8,
    description = $9,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetWorkExperienceById :one
SELECT * FROM work_experiences
WHERE id = @id::bigint
LIMIT 1;

-- name: BatchInsertWorkExperienceSkills :exec
INSERT INTO work_experience_skills (work_experience_id, user_skill_id)
SELECT @work_experience_id::bigint, unnest(@user_skill_id::bigint[])
ON CONFLICT (work_experience_id, user_skill_id) DO NOTHING;

-- name: DeleteWorkExperienceSkillsByWorkExperience :many
DELETE FROM work_experience_skills
WHERE work_experience_id = @work_experience_id::bigint
RETURNING user_skill_id;

-- name: BatchInsertWorkExperienceFiles :many
INSERT INTO work_experience_files
  (work_experience_id, url)
SELECT @work_experience_id::bigint, UNNEST(@url::text[])
RETURNING *;

-- name: DeleteWorkExperienceFilesByWorkExperienceId :exec
DELETE FROM work_experience_files
WHERE work_experience_id = @work_experience_id::bigint;

-- name: GetUserEducationFileURLs :many
SELECT url FROM education_files
WHERE education_id = @education_id::bigint;

-- name: GetWorkExperienceFileURLs :many
SELECT url FROM work_experience_files
WHERE work_experience_id = @work_experience_id::bigint;

-- name: GetUserSkillIDsByName :many
SELECT us.id FROM user_skills us
JOIN skills s ON us.skill_id = s.id
WHERE s.name = ANY(@name::text[]);

-- name: GetUserProfile :one
SELECT 
    u.id, u.full_name, u.avatar_url, u.bio, u.open_to_work, u.followers_count, u.followings_count,
    ud.phone_number, ud.gender, ud.location, ud.portfolio_url, ud.about
FROM users u
LEFT JOIN user_details ud ON u.id = ud.user_id 
WHERE u.id = $1
LIMIT 1;

-- name: GetUserSocialLinks :many
SELECT platform, url FROM user_social_links
WHERE user_id = @user_id::bigint;

-- name: GetUserSkills :many
SELECT us.main_skill, s.name FROM user_skills us
LEFT JOIN skills s ON us.skill_id = s.id
WHERE us.user_id = @user_id::bigint;

-- name: GetWorkExperiencesByUserId :many
SELECT 
  we.*, c.name AS company_name, 
  COALESCE(array_agg(DISTINCT s.name) FILTER (WHERE s.name IS NOT NULL), '{}') AS skills, 
  COALESCE(array_agg(DISTINCT wef.url) FILTER (WHERE wef.url IS NOT NULL), '{}') AS file_urls,
  COUNT(*) OVER () AS total_rows
FROM work_experiences we 
LEFT JOIN companies c ON we.company_id = c.id 
LEFT JOIN work_experience_files wef ON we.id = wef.work_experience_id 
LEFT JOIN work_experience_skills wes ON we.id = wes.work_experience_id 
LEFT JOIN user_skills us ON wes.user_skill_id = us.id 
LEFT JOIN skills s ON us.skill_id  = s.id
WHERE we.user_id = @user_id::bigint
GROUP BY we.id, c.name
ORDER BY we.finish_date DESC, we.start_date DESC
OFFSET $1
LIMIT $2;

-- name: GetEducationsByUserId :many
SELECT 
  e.*, 
  schools.name AS school_name, 
  COALESCE(array_agg(DISTINCT skills.name) FILTER (WHERE skills.name IS NOT NULL), '{}') AS skills, 
  COALESCE(array_agg(DISTINCT ef.url) FILTER (WHERE ef.url IS NOT NULL), '{}') AS file_urls,
  COUNT(*) OVER () AS total_rows
FROM educations e 
LEFT JOIN schools ON e.school_id = schools.id 
LEFT JOIN education_files ef ON e.id = ef.education_id 
LEFT JOIN education_skills es ON e.id = es.education_id 
LEFT JOIN user_skills us ON es.user_skill_id = us.id 
LEFT JOIN skills ON us.skill_id = skills.id
WHERE e.user_id = @user_id::bigint
GROUP BY e.id, schools.name
ORDER BY e.finish_date DESC, e.start_date DESC
OFFSET $1
LIMIT $2;


-- name: GetCertificatesByUserId :many
SELECT 
  c.*, 
  io.name AS issuing_organization_name,
  COUNT(*) OVER () AS total_rows
FROM certificates c 
LEFT JOIN issuing_organizations io ON c.issuing_organization_id = io.id 
WHERE user_id = @user_id::bigint
ORDER BY c.issue_date desc, c.expiration_date DESC
OFFSET $1
LIMIT $2;

-- name: GetFollowedUsersByUserId :many
SELECT 
  u.id, u.full_name ,u.avatar_url ,u.bio ,u.open_to_work,
  COUNT(*) OVER () AS total_rows
FROM followings f 
LEFT JOIN users u ON f.follow_user_id = u.id 
WHERE f.user_id = @user_id::bigint
OFFSET $1
LIMIT $2;

-- name: UpdateUserOpenToWork :one
UPDATE users
SET open_to_work = @open_to_work::boolean,
    updated_at = NOW()
WHERE id = @user_id::bigint
RETURNING id;

-- name: InsertJobPosition :one
INSERT INTO job_positions (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING *;

-- name: BatchInsertUserJobInterests :exec
INSERT INTO user_job_interests (user_id, job_position_id)
SELECT @user_id::bigint, UNNEST(@job_position_id::bigint[]);

-- name: BatchInsertUserLocationTypeInterests :exec
INSERT INTO user_location_type_interests (user_id, location_type)
SELECT @user_id::bigint, UNNEST(@location_type::varchar(10)[]);

-- name: BatchInsertUserEmploymentTypeInterests :exec
INSERT INTO user_employment_type_interests (user_id, employment_type)
SELECT @user_id::bigint, UNNEST(@employment_type::varchar(10)[]);

-- name: BatchDeleteUserJobInterests :exec
DELETE FROM user_job_interests
WHERE user_id = @user_id::bigint;

-- name: BatchDeleteUserLocationTypeInterests :exec
DELETE FROM user_location_type_interests
WHERE user_id = @user_id::bigint;

-- name: BatchDeleteUserEmploymentTypeInterests :exec
DELETE FROM user_employment_type_interests
WHERE user_id = @user_id::bigint;

-- name: DeleteWorkExperienceById :exec
DELETE FROM work_experiences
WHERE id = @id::bigint AND user_id = @user_id::bigint;

-- name: DeleteEducationById :exec
DELETE FROM educations
WHERE id = @id::bigint AND user_id = @user_id::bigint;

-- name: DeleteCertificateById :exec
DELETE FROM certificates
WHERE id = @id::bigint AND user_id = @user_id::bigint;

-- name: LockUserForUpdate :one
SELECT 1
FROM users
WHERE id = $1
FOR UPDATE;

-- name: UpdateUserFollowingsCount :one
UPDATE users
SET followings_count = GREATEST(followings_count + @value::smallint, 0),
    updated_at = NOW()
WHERE id = @user_id::bigint
RETURNING followings_count;

-- name: UpdateUserFollowersCount :one
UPDATE users
SET followers_count = GREATEST(followers_count + @value::smallint, 0),
    updated_at = NOW()
WHERE id = @user_id::bigint
RETURNING followers_count;

-- name: InsertFollowings :one
INSERT INTO followings (user_id, follow_user_id)
VALUES (@user_id::bigint, @follow_user_id::bigint)
ON CONFLICT (user_id, follow_user_id) DO NOTHING
RETURNING *;

-- name: DeleteFollowings :one
DELETE FROM followings
WHERE user_id = @user_id::bigint AND follow_user_id = @follow_user_id::bigint
RETURNING id;