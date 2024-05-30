-- name: GetSchools :many
SELECT *, COUNT(*) OVER () AS total_rows
FROM schools
OFFSET $1
LIMIT $2;

-- name: GetCompanies :many
SELECT *, COUNT(*) OVER () AS total_rows
FROM companies
OFFSET $1
LIMIT $2;

-- name: GetIssuingOrganizations :many
SELECT *, COUNT(*) OVER () AS total_rows
FROM issuing_organizations
OFFSET $1
LIMIT $2;

-- name: GetSkills :many
SELECT *, COUNT(id) OVER () AS total_rows
FROM skills
OFFSET $1
LIMIT $2;