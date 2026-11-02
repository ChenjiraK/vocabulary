# Register and current user

`POST /auth/register` accepts `multipart/form-data`, JSON, or URL-encoded form
data. Use multipart when uploading a profile picture. API field names use
snake_case; map frontend `displayName`, `confirmPassword`, and `profilePicture`
to the fields below.

| Field | Rules |
| --- | --- |
| `display_name` | Required, trimmed, 1–100 characters |
| `email` | Required, valid email, trimmed, lowercased, at most 255 characters |
| `password` | Required, at least 8 characters and at most 72 UTF-8 bytes |
| `confirm_password` | Required, exactly matches password (including spaces) |
| `profile_picture` | Optional multipart file; JPEG/PNG, at most 2 MiB, dimensions at most 4096 × 4096 |

Passwords are hashed with bcrypt. Confirmation is validated but never persisted.
User and identity creation use one database transaction. An email already used
by a password or social account returns 409; registration does not attach a
password to an existing social account. Concurrent password registrations are
also protected by the existing database unique constraint.

## Storage configuration

Add these settings to the backend environment (examples also in `.env.example`):

```dotenv
AVATAR_UPLOAD_DIR=./uploads/avatars
PUBLIC_ASSET_BASE_URL=http://localhost:8000
```

The backend validates and re-encodes image content, generates a random filename,
and stores `/uploads/avatars/<generated-name>.png` (or `.jpg`) in `users.avatar_url`.
Only that directory is publicly served at `/uploads/avatars/`. Uploaded files are
ignored by Git. If database creation fails, the newly stored avatar is removed.

`PUBLIC_ASSET_BASE_URL` must be the browser-accessible origin of this backend,
or a CDN serving the same avatar paths. Without it, development defaults to
`http://localhost:<APP_PORT>` (port 8080 if `APP_PORT` is absent).
In production, mount persistent storage at `AVATAR_UPLOAD_DIR` and use an HTTPS
public base URL. Multiple backend instances must share the avatar storage.

## Example request

```bash
curl http://localhost:8000/auth/register \
  -F 'display_name=Jane' \
  -F 'email=jane@example.com' \
  -F 'password=example-password' \
  -F 'confirm_password=example-password' \
  -F 'profile_picture=@/absolute/path/avatar.png'
```

For browser `FormData`, let the browser set `Content-Type` and its boundary.
Omit `profile_picture` when no image was selected. JSON registration supports
the four text fields; files must be sent as multipart.

A successful request returns HTTP 201 and the existing auth response with a
JWT and `user`. `user` now also contains `email` and `profile_picture`:

```json
{
  "id": 1,
  "display_name": "Jane",
  "email": "jane@example.com",
  "profile_picture": "http://localhost:8000/uploads/avatars/<generated-name>.png",
  "avatar_url": "http://localhost:8000/uploads/avatars/<generated-name>.png",
  "user_type": "registered"
}
```

`avatar_url` remains available for existing clients. `profile_picture` is `null`
when no avatar exists. `GET /auth/me` with `Authorization: Bearer <token>` returns
the same profile fields along with the existing token and game statistics.

Errors use `{ "message": "..." }`: 400 for invalid fields/images, 409 for an
existing email, 413 for oversized uploads/bodies, 415 for unsupported content
types, and 500 for persistence failures. The total request limit is 2 MiB +
64 KiB for form overhead. No schema migration is needed for this change.
