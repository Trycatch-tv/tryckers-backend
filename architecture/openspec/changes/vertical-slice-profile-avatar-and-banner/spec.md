# Spec: vertical-slice-profile-avatar-and-banner

## Capability

Authenticated users can upload and replace their own profile avatar and profile banner.

## API Contract

| Method | Path | Auth | Request | Response |
| --- | --- | --- | --- | --- |
| POST | `/api/v1/users/me/avatar` | Bearer token | multipart form field `file` | Updated user |
| POST | `/api/v1/users/me/banner` | Bearer token | multipart form field `file` | Updated user |
| DELETE | `/api/v1/users/me/avatar` | Bearer token | none | Updated user |
| DELETE | `/api/v1/users/me/banner` | Bearer token | none | Updated user |

## Acceptance Criteria

- Given an authenticated user uploads a valid JPG/PNG/WEBP avatar, when the request succeeds, then `avatar_url` persists and the profile renders the new image after refresh.
- Given an authenticated user uploads a valid JPG/PNG/WEBP banner, when the request succeeds, then `banner_url` persists and the profile renders the new banner after refresh.
- Given a user uploads an unsupported file type, when the backend validates the file, then the request is rejected.
- Given a user uploads a file over the size limit, when the backend validates the file, then the request is rejected.
- Given a user has no avatar, when the profile renders, then a default avatar or initials fallback appears.
- Given a user has no banner, when the profile renders, then a generated gradient banner appears.
- Given a mobile viewport, when the profile renders, then avatar and banner remain visually acceptable.
- Given a user removes avatar or banner, when the request succeeds, then fallback visuals are rendered.

## Validation Rules

- Only authenticated users can upload media.
- Only the authenticated user's media can be changed.
- Avatar max size: 2 MB.
- Banner max size: 5 MB.
- Supported MIME types: JPG, PNG, WEBP.
- Server generates safe filenames.
- Executable and unknown file types are rejected.

## Persistence Rules

- `users.avatar_url` stores the public avatar URL.
- `users.banner_url` stores the public banner URL.
- Existing `profile_picture` remains supported for backwards compatibility and is updated when avatar changes.

## Frontend Impact

- Profile page renders `avatar_url` before `profile_picture`.
- Profile page renders `banner_url` as cover image when present.
- User upload controls are shown only on the authenticated user's own profile.
- Upload errors are shown through existing notification service.

## Compatibility

This is backwards-compatible for clients using `profile_picture`; new clients should prefer `avatar_url`.
