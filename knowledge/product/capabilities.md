---
type: capabilities
project_state: ai-assisted
generated_by: kaddo-bootstrap
template_version: 1
refined_by: capability-agent
---

> Idioma del proyecto: **español**. Escribe este conocimiento en español. Mantén en inglés el código, los nombres de archivo, los comandos y las claves de configuración.

# Existing Capabilities

## Capability Domains

### Domain: Identity and Access Management (IAM)

**Purpose:** Gestión integral de registro, autenticación, control de sesiones mediante tokens JWT y autorización basada en roles.

**Evidence summary:**
- `src/internal/api/routes/router.go` (`/api/v1/register`, `/api/v1/login`, `/api/v1/refresh-token`)
- `src/internal/api/handlers/user_handler.go` (`CreateUser`, `Login`, `RefreshToken`)
- `src/internal/api/middlewares/auth.go` (`AuthMiddleware`, `RoleMiddleware`)
- `src/internal/models/user.go` (`User`, `enums.UserRole`)
- `tryckers-frontend/src/app/auth/` (`auth.routes.ts`, `guards/authenticated.guard.ts`, `store/auth-store.ts`)

#### Capability: User Registration and Password Hashing
- Status: implemented
- Capability type: business
- User-facing: yes
- Evidence:
  - `src/internal/api/routes/router.go:28` (`POST /api/v1/register`)
  - `src/internal/api/handlers/user_handler.go` (`CreateUser`)
  - `src/internal/models/user.go` (`User`)
- Related flows: User registration flow
- Related data: `users` table (ID, Name, Username, Email, Password hash)
- Related integrations: Ninguna
- Current behavior: Permite crear cuentas mediante DTO con validación de unicidad de username y email. La contraseña se hashea antes de persistir.
- Known constraints: No requiere verificación de email previa al login.
- Risks or uncertainty: Falta de rate limiting en registro público.
- Open questions:
  - [open] ¿Se agregará verificación obligatoria por correo electrónico?

#### Capability: JWT Authentication and Route Protection
- Status: implemented
- Capability type: technical
- User-facing: yes
- Evidence:
  - `src/internal/api/routes/router.go:29` (`POST /api/v1/login`)
  - `src/internal/api/middlewares/auth.go` (`AuthMiddleware`)
  - `tryckers-frontend/src/app/auth/guards/authenticated.guard.ts`
- Related flows: Login flow
- Related data: Bearer JWT claims
- Related integrations: Ninguna
- Current behavior: Emite token JWT al validar credenciales; el middleware backend valida firma y expiración; el frontend protege rutas mediante Angular Guards.
- Known constraints: JWT stateless sin lista negra centralizada de revocación inmediata.
- Risks or uncertainty: Revocación de tokens antes de expiración.
- Open questions:
  - [open] ¿Cuál será la vigencia estándar del access token (e.g. 15 min vs 24h)?

#### Capability: Refresh Token Rotation and Session Renewal
- Status: partial
- Capability type: technical
- User-facing: no
- Evidence:
  - `src/internal/api/routes/router.go:30` (`POST /api/v1/refresh-token`)
  - `src/internal/api/handlers/user_handler.go` (`RefreshToken`)
- Related flows: Refresh session flow
- Related data: Refresh token payload
- Related integrations: Ninguna
- Current behavior: Endpoint expuesto para refrescar sesión cuando un token expira; el frontend cuenta con interceptor para reintentar peticiones con error 401.
- Known constraints: Mecanismo de persistencia de refresh tokens en base de datos aún no formalizado.
- Risks or uncertainty: Manejo concurrente de múltiples peticiones refrescando simultáneamente en frontend.
- Open questions:
  - [open] ¿El refresh token se enviará en cookie HttpOnly segura o en el body JSON?

#### Capability: Role-Based Access Control (RBAC)
- Status: implemented
- Capability type: technical
- User-facing: internal
- Evidence:
  - `src/internal/api/routes/router.go:37` (`middlewares.RoleMiddleware(enums.Admin, enums.Member)`)
  - `src/internal/enums/user_role.go` (`Admin`, `Member`, `Recruiter`)
- Related flows: Protected administration flows
- Related data: `User.Role`
- Related integrations: Ninguna
- Current behavior: Middleware valida que el rol contenido en el token coincida con los roles autorizados para el endpoint.
- Known constraints: Matriz de permisos estática en middleware.
- Risks or uncertainty: Ausencia de permisos granulares a nivel de recurso.
- Open questions:
  - [open] ¿Qué endpoints estarán reservados exclusivamente para el rol `Admin`?

---

### Domain: Professional Identity and Profiles

**Purpose:** Creación, edición, visualización y descubrimiento de perfiles profesionales técnicos de los miembros de la comunidad.

**Evidence summary:**
- `src/internal/api/routes/router.go` (`/api/v1/users`, `/api/v1/perfil/:username`, `/api/v1/users/:ownerId/posts`)
- `src/internal/api/handlers/user_handler.go` (`GetAll`, `Perfil`)
- `src/internal/models/user.go` (Headline, Bio, PitchVideo, EnglishLevel, Seniority, Availability, Interests, Links)
- `tryckers-frontend/src/app/tryckers/pages/profile-page/profile-page.ts`
- `tryckers-frontend/src/app/tryckers/pages/dashboard-page/dashboard-page.ts`

#### Capability: Public Profile Resolution by Username
- Status: implemented
- Capability type: product
- User-facing: yes
- Evidence:
  - `src/internal/api/routes/router.go:38` (`GET /api/v1/perfil/:username`)
  - `src/internal/api/handlers/user_handler.go` (`Perfil`)
  - `tryckers-frontend/src/app/tryckers/pages/profile-page/profile-page.ts`
- Related flows: Profile flow
- Related data: `users` table, relaciones de posts
- Related integrations: Ninguna
- Current behavior: Retorna información pública y posts asociados a un miembro consultado mediante su `@username`.
- Known constraints: Requiere autenticación activa para consultar perfiles en frontend (`canMatch: [AuthenticatedGuard]`).
- Risks or uncertainty: Dificulta visualización pública o indexación externa (SEO) si el perfil requiere login obligatorio.
- Open questions:
  - [open] ¿El perfil público debería ser accesible sin login para visitantes externos y buscadores?

#### Capability: Professional Attributes Representation
- Status: partial
- Capability type: business
- User-facing: yes
- Evidence:
  - `src/internal/models/user.go:21-34` (`GithubURL`, `LinkedinURL`, `PitchVideo`, `Seniority`, `EnglishLevel`, `Availability`, `Interests`)
- Related flows: Profile editing flow, Directory discovery
- Related data: Columnas profesionales en tabla `users`
- Related integrations: Enlaces externos a GitHub, LinkedIn, YouTube
- Current behavior: El modelo almacena campos profesionales, pero las interfaces de edición granular y validadores específicos están en desarrollo.
- Known constraints: Los intereses y tecnologías se guardan como cadenas de texto plano.
- Risks or uncertainty: Inconsistencias al filtrar si las skills no pertenecen a una taxonomía estandarizada.
- Open questions:
  - [open] ¿Se mantendrán las habilidades como texto plano o se creará una tabla de catálogo de skills?

#### Capability: Member Directory and Discovery
- Status: partial
- Capability type: product
- User-facing: yes
- Evidence:
  - `src/internal/api/routes/router.go:37` (`GET /api/v1/users`)
  - `tryckers-frontend/src/app/tryckers/pages/dashboard-page/dashboard-page.ts` (`DashboardPage`)
- Related flows: User directory flow
- Related data: `users` list
- Related integrations: Ninguna
- Current behavior: Obtiene y muestra cards de miembros registrados con enlaces a sus perfiles.
- Known constraints: No dispone de filtros combinados (país, seniority, skills), ordenamiento ni paginación en backend.
- Risks or uncertainty: Problemas de rendimiento y escalabilidad al crecer el número de usuarios.
- Open questions:
  - [open] ¿Qué parámetros de búsqueda y paginación se incorporarán en `GET /api/v1/users`?

---

### Domain: Media and Storage

**Purpose:** Almacenamiento, gestión y servicio de archivos multimedia de usuario (avatares y banners decorativos).

**Evidence summary:**
- `src/internal/api/routes/router.go:39-42` (`/api/v1/users/me/avatar`, `/api/v1/users/me/banner`)
- `src/internal/services/storage/local_storage.go`
- `src/internal/api/handlers/user_handler.go` (`UploadAvatar`, `UploadBanner`, `DeleteAvatar`, `DeleteBanner`)
- `src/internal/models/user.go` (`AvatarURL`, `BannerURL`)

#### Capability: Avatar and Banner Upload with Local Storage
- Status: implemented
- Capability type: technical
- User-facing: yes
- Evidence:
  - `src/internal/api/routes/router.go:39-42`
  - `src/internal/services/storage/local_storage.go`
- Related flows: Profile media management flow
- Related data: `AvatarURL`, `BannerURL` en `users`
- Related integrations: Sistema de archivos local (`uploads/avatars`, `uploads/banners`)
- Current behavior: Permite subir y eliminar imágenes asociadas al usuario autenticado, almacenándolas en el filesystem y persistiendo la ruta relativa.
- Known constraints: Almacenamiento local ligado al contenedor/servidor; no distribuido.
- Risks or uncertainty: Pérdida de imágenes en despliegues efímeros o multi-instancia si no se monta volumen persistente o S3.
- Open questions:
  - [open] ¿Cuándo se migrará a un almacenamiento compatible con S3 (MinIO / AWS S3 / Cloudflare R2)?

---

### Domain: Community Discussions and Content (Posts and Comments)

**Purpose:** Publicación, edición, categorización, consumo y moderación de contenido técnico y comentarios comunitarios.

**Evidence summary:**
- `src/internal/api/routes/router.go:45-59` (`/api/v1/posts`, `/api/v1/comments`)
- `src/internal/api/handlers/post_handler.go`
- `src/internal/api/handlers/comment_handler.go`
- `src/internal/models/post.go` (`Post`, `enums.PostStatus`, `enums.PostType`)
- `src/internal/models/comments.go` (`Comment`)
- `tryckers-frontend/src/app/post/`

#### Capability: Post Publication and Media Tagging
- Status: implemented
- Capability type: product
- User-facing: yes
- Evidence:
  - `src/internal/api/routes/router.go:45` (`POST /api/v1/posts`)
  - `src/internal/models/post.go:12-28` (`Title`, `Content`, `Image`, `Type`, `Tags`, `MediaURL`)
- Related flows: Create post flow
- Related data: `posts` table
- Related integrations: URLs multimedia externas
- Current behavior: Permite a usuarios autenticados crear posts con título, cuerpo, imagen de portada, tags (separados por coma) y enlaces a video.
- Known constraints: Tags representados como cadena de texto plano (`Tags string`).
- Risks or uncertainty: Dificultad para generar agregaciones y búsquedas exactas por tags.
- Open questions:
  - [open] ¿Se migrarán los tags a una entidad relacional `tags` con tabla intermedia `post_tags`?

#### Capability: Logical Deletion of Content (Soft Delete)
- Status: implemented
- Capability type: operational
- User-facing: yes
- Evidence:
  - `src/internal/models/post.go:88-90` (`IsDeleted() bool { return p.Status == enums.DELETED }`)
  - `src/internal/api/routes/router.go:50,58` (`DELETE /posts/:id`, `DELETE /comments/:id`)
- Related flows: Delete post flow, Delete comment flow
- Related data: `posts.status`, `comments.status`
- Related integrations: Ninguna
- Current behavior: Marca los registros como eliminados mediante el estado sin eliminarlos físicamente de la base de datos, filtrándolos de consultas públicas.
- Known constraints: Uso de enum `Status` en lugar de columna `deleted_at` nativa de GORM.
- Risks or uncertainty: Auditoría incompleta sobre la fecha y autor exacto del borrado.
- Open questions:
  - [open] ¿Se agregará `deleted_at` timestamp para compatibilidad estándar con `gorm.DeletedAt`?

#### Capability: Discussion Comments Threading
- Status: implemented
- Capability type: product
- User-facing: yes
- Evidence:
  - `src/internal/api/routes/router.go:55-58` (`POST /comments`, `GET /posts/:id/comments`, `PUT /comments/:id`, `DELETE /comments/:id`)
  - `src/internal/models/comments.go`
- Related flows: Commenting flow
- Related data: `comments` table (UserID, PostID, Content, Status)
- Related integrations: Ninguna
- Current behavior: Permite crear, listar, actualizar y eliminar comentarios vinculados a un post específico.
- Known constraints: Modelo lineal sin anidamiento jerárquico de respuestas (*replies / threaded comments*).
- Risks or uncertainty: Debates profundos pueden volverse difíciles de leer linealmente.
- Open questions:
  - [open] ¿Se soportarán respuestas anidadas a comentarios en una fase posterior?

---

### Domain: Social Signals and Curation (Votes and Billboard)

**Purpose:** Medición del interés comunitario mediante votos unívocos y exposición de contenido destacado en cartelera semanal.

**Evidence summary:**
- `src/internal/api/routes/router.go:47,51` (`GET /cartelera`, `POST /posts/:id/vote`)
- `src/internal/models/postsvotes.go` (`PostVote`)
- `src/internal/api/handlers/post_handler.go` (`Cartelera`, `PostVote`)
- `tryckers-frontend/src/app/post/pages/cartelera-page/cartelera-page.ts`

#### Capability: Post Upvoting with Toggle
- Status: implemented
- Capability type: product
- User-facing: yes
- Evidence:
  - `src/internal/api/routes/router.go:51` (`POST /api/v1/posts/:id/vote`)
  - `src/internal/models/postsvotes.go` (`PostVote`)
  - `src/internal/api/handlers/post_handler.go` (`PostVote`)
- Related flows: Vote post flow
- Related data: `posts_votes` table (UserID, PostID, CreatedAt)
- Related integrations: Ninguna
- Current behavior: Registra el voto si no existe o lo remueve si el usuario ya había votado previamente (*toggle*), actualizando el contador.
- Known constraints: Operación protegida por autenticación.
- Risks or uncertainty: Posible condición de carrera en incremento de contador si no se maneja transaccionalmente o con `gorm.Expr`.
- Open questions:
  - [open] ¿Se implementará rate limiting para prevenir manipulación coordinada de votos?

#### Capability: Weekly Billboard Discovery (Cartelera Semanal)
- Status: implemented
- Capability type: business
- User-facing: yes
- Evidence:
  - `src/internal/api/routes/router.go:47` (`GET /api/v1/cartelera`)
  - `src/internal/api/handlers/post_handler.go` (`Cartelera`)
  - `tryckers-frontend/src/app/post/pages/cartelera-page/cartelera-page.ts`
- Related flows: Weekly popular posts flow
- Related data: `posts` con filtro de fecha semanal y ordenamiento por votos
- Related integrations: Ninguna
- Current behavior: Consulta y lista los posts publicados o votados en el periodo semanal ordenados por popularidad.
- Known constraints: Algoritmo actual basado en suma simple de votos sin ponderación temporal ni decaimiento de antigüedad.
- Risks or uncertainty: Posts más antiguos con muchos votos pueden monopolizar la cartelera sin rotación fresca.
- Open questions:
  - [open] ¿Qué fórmula matemática se adoptará para el score de trending (e.g. Hacker News / Reddit algorithm)?

---

## Capability Gaps

- [gap] Búsqueda y filtrado avanzado de miembros por skills, seniority y país
  - Domain: Professional Identity and Profiles
  - Related capability: Member Directory and Discovery
  - Impact: high
  - Possible roadmap candidate: yes

- [gap] Almacenamiento de archivos en infraestructura de objetos en la nube (S3 / Cloudflare R2)
  - Domain: Media and Storage
  - Related capability: Avatar and Banner Upload with Local Storage
  - Impact: high
  - Possible roadmap candidate: yes

- [gap] Catálogo normalizado de tecnologías y habilidades técnicas (Skills Catalog)
  - Domain: Professional Identity and Profiles
  - Related capability: Professional Attributes Representation
  - Impact: medium
  - Possible roadmap candidate: yes

- [gap] Rotación y persistencia segura de Refresh Tokens con invalidación centralizada
  - Domain: Identity and Access Management (IAM)
  - Related capability: Refresh Token Rotation and Session Renewal
  - Impact: high
  - Possible roadmap candidate: yes

- [gap] Módulo de oportunidades laborales y conexión con recruiters
  - Domain: Professional Identity and Profiles
  - Related capability: Public Profile Resolution by Username
  - Impact: high
  - Possible roadmap candidate: yes

- [gap] Algoritmo de trending con decaimiento temporal para la Cartelera
  - Domain: Social Signals and Curation (Votes and Billboard)
  - Related capability: Weekly Billboard Discovery (Cartelera Semanal)
  - Impact: medium
  - Possible roadmap candidate: yes

---

## Roadmap Candidate Signals

- [candidate] Motor de búsqueda y filtrado estructurado de talento
  - Domain: Professional Identity and Profiles
  - Related capability: Member Directory and Discovery
  - Based on: gap

- [candidate] Adaptador de almacenamiento en la nube compatible con S3
  - Domain: Media and Storage
  - Related capability: Avatar and Banner Upload with Local Storage
  - Based on: risk

- [candidate] Formalización de Refresh Tokens y revocación de sesión
  - Domain: Identity and Access Management (IAM)
  - Related capability: Refresh Token Rotation and Session Renewal
  - Based on: partial capability

- [candidate] Módulo de publicación y postulación a oportunidades laborales
  - Domain: Professional Identity and Profiles
  - Related capability: Public Profile Resolution by Username
  - Based on: business goal

- [candidate] Algoritmo de Cartelera con Time Decay Score
  - Domain: Social Signals and Curation (Votes and Billboard)
  - Related capability: Weekly Billboard Discovery (Cartelera Semanal)
  - Based on: open question

────────────────────────
Agent: capability-agent

Produced:
knowledge/product/capabilities.md

Next:
roadmap-agent
────────────────────────
