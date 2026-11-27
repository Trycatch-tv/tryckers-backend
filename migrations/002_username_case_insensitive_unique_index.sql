-- WI-009: Refuerzo de unicidad case-insensitive sobre users.username
-- Previamente auditado en RDS: 0 duplicados lógicos existentes.
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_lower_unique
ON users (LOWER(TRIM(username)));
