# Proposal: vertical-slice-profile-avatar-and-banner

## Metadata

| Field | Value |
| --- | --- |
| type | vertical-slice |
| domain | profile |
| stage | stage-1-mvp-consolidation |
| priority | high |
| status | planned |
| suggested tag | `vs/vertical-slice-profile-avatar-and-banner/v1.0.0` |

## Summary

Enable Tryckers users to personalize and strengthen their public identity by uploading and managing profile avatars and profile banners/covers.

This slice improves profile quality, recruiter perception, social credibility, discoverability, and platform engagement by making the profile feel personal and alive.

## Problem

Profiles currently rely on static or URL-based visuals. Users need a first-party way to upload identity assets that persist, render consistently across profile surfaces, and remain safe through validation and authenticated access.

## Goals

- Add backend media upload support for avatars and banners.
- Persist `avatar_url` and `banner_url` on users.
- Validate image type and size.
- Store files in local filesystem for MVP.
- Expose public media URLs.
- Update profile rendering with fallbacks.
- Add upload UX with preview, loading, and errors.

## Non-Goals

- Image cropping.
- S3 production bucket provisioning.
- AI-generated banners.
- Animated avatars.
- Profile themes.

## Impact

| Area | Impact |
| --- | --- |
| API | Adds authenticated avatar and banner upload endpoints. |
| Backend domain | Adds user media persistence and storage service. |
| Database | Adds `avatar_url` and `banner_url` to users. |
| Frontend | Adds upload controls and profile rendering updates. |
| Tests | Manual validation plus compile/build verification for this VS. |

## Review Notes

Supabase DDL should be applied through the Supabase MCP server when a project id is available:

```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS banner_url TEXT;
```
