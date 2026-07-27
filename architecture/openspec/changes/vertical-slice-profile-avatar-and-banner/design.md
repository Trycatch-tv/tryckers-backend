# Design: vertical-slice-profile-avatar-and-banner

## Architecture Fit

This VS follows the layered modular monolith:

```text
Client -> Route -> Handler -> UserService -> UserRepository -> PostgreSQL
                       |
                       `-> StorageService -> local filesystem
```

Storage is isolated behind `internal/services/storage` so local storage can later be replaced by S3-compatible storage.

## Layer Changes

| Layer | Planned Change |
| --- | --- |
| Route | Add `POST` and `DELETE` endpoints for `/api/v1/users/me/avatar` and `/api/v1/users/me/banner`. |
| Handler | Accept authenticated multipart uploads. |
| DTO | Return updated user data through existing model shape. |
| Service | Validate ownership from auth context, store media, update user URLs. |
| Repository | Add lookup by id and media URL update methods. |
| Model | Add `AvatarURL` and `BannerURL`. |
| Migration | Add `avatar_url TEXT` and `banner_url TEXT`. |
| Frontend | Add profile upload UX and fallback rendering. |

## Data Flow

```text
User selects image
Frontend previews image
Frontend POSTs multipart file
Auth middleware identifies current user
Storage service validates MIME/size and stores file
User repository updates media URL
Frontend refreshes profile and renders updated media
```

## Storage Design

MVP local folders:

```text
uploads/
|-- avatars/
`-- banners/
```

Public URL base:

```text
/uploads/<kind>/<safe-file-name>
```

Production recommendation:

```text
tryckers/
|-- avatars/
`-- banners/
```

on an S3-compatible bucket.

## Validation

Allowed MIME types:

- `image/jpeg`
- `image/png`
- `image/webp`

Size limits:

- Avatar: 2 MB.
- Banner: 5 MB.

Filenames are generated server-side using media kind, user id, timestamp, random suffix, and a safe extension.

## Replacement Strategy

After a successful database update, the old local media file is removed when it belongs to the local uploads path.

If the database update fails after file save, the new file is removed to avoid orphaned files.

## Frontend Design

Profile page adds:

- Avatar upload button.
- Banner upload button.
- Remove actions for existing avatar/banner.
- Client-side preview before upload.
- Loading state per media type.
- Error feedback through notifications.
- Fallback avatar and generated gradient banner.

## Tradeoffs

Local filesystem storage is enough for MVP development and keeps this slice deployable without provisioning cloud infrastructure. The storage service boundary keeps future S3 migration contained.
