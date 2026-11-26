-- WI-008: Limpieza controlada e idempotente de referencias multimedia antiguas
-- No modifica enlaces externos (github_url, linkedin_url, pitch_video, posts.media_url).
BEGIN;

UPDATE users
SET avatar_url = ''
WHERE avatar_url IS NOT NULL
  AND avatar_url <> ''
  AND avatar_url NOT LIKE 'avatars/%';

UPDATE users
SET profile_picture = ''
WHERE profile_picture IS NOT NULL
  AND profile_picture <> ''
  AND profile_picture NOT LIKE 'avatars/%';

UPDATE users
SET banner_url = ''
WHERE banner_url IS NOT NULL
  AND banner_url <> ''
  AND banner_url NOT LIKE 'banners/%';

UPDATE posts
SET image = ''
WHERE image IS NOT NULL
  AND image <> ''
  AND image NOT LIKE 'posts/%';

UPDATE comments
SET image = ''
WHERE image IS NOT NULL
  AND image <> ''
  AND image NOT LIKE 'posts/%'
  AND image NOT LIKE 'comments/%';

COMMIT;
