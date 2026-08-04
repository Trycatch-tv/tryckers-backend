---
type: current-state
project_state: ai-assisted
generated_by: kaddo-bootstrap
template_version: 1
refined_by: architecture-agent
---

> Idioma del proyecto: **español**. Escribe este conocimiento en español. Mantén en inglés el código, los nombres de archivo, los comandos y las claves de configuración.

# Current State

## System Overview

Tryckers es una plataforma web desacoplada compuesta por dos componentes principales en arquitectura multirepo:
1. **tryckers-backend (Core):** API REST monolítica modular construida en Go (1.24.4) utilizando el framework Gin, persistencia relacional con PostgreSQL mediante GORM y almacenamiento local de archivos estáticos.
2. **tryckers-frontend:** Single Page Application (SPA) desarrollada en Angular 20 con TypeScript 5.8, gestión reactiva de estado mediante `@ngrx/signals` (SignalStore), componentes de interfaz con PrimeNG 20 y estilos utilitarios con TailwindCSS 4.

## Modules

### Backend Modules (`tryckers-backend/src/internal/`)
- **`api/routes`:** Configuración centralizada de endpoints versión 1 (`SetupV1`) separando grupos públicos y protegidos.
- **`api/handlers`:** Controladores HTTP para usuarios (`UserHandler`), publicaciones (`PostHandler`) y comentarios (`CommentHandler`).
- **`api/middlewares`:** Middleware de autenticación JWT (`AuthMiddleware`), autorización por roles (`RoleMiddleware`) y CORS.
- **`models`:** Estructuras de datos GORM para `User`, `Post`, `Comment` y `PostVote`.
- **`repository`:** Capa de acceso a datos (`UserRepository`, `PostRepository`, `CommentRepository`).
- **`services`:** Lógica de negocio de dominio (`UserService`, `PostService`, `CommentService`).
- **`services/storage`:** Abstracción del sistema de archivos (`LocalStorage` implementado para `/uploads`).
- **`dtos` / `enums` / `errors` / `utils`:** Objetos de transferencia, tipos enumerados de dominio (`UserRole`, `PostStatus`, `Country`), utilitarios criptográficos (bcrypt) y generadores de tokens JWT.

### Frontend Modules (`tryckers-frontend/src/app/`)
- **`auth`:** Rutas de autenticación (`/auth/login`, `/auth/register`), guards (`AuthenticatedGuard`, `NotAuthenticatedGuard`) y gestión de sesión con `AuthStore`.
- **`tryckers`:** Vistas de directorio de miembros (`DashboardPage`) y perfil público (`ProfilePage`).
- **`post`:** Vistas y componentes de creación/detalle de posts y la cartelera semanal (`CarteleraPage`).
- **`shared` / `core` / `components`:** Layouts compartidos (`MainLayoutComponent`), interceptores HTTP para Bearer Token y servicios de notificación.

## Dependencies and Integrations

- **Framework Web:** Gin v1.10.1 (`github.com/gin-gonic/gin`).
- **ORM & Driver DB:** GORM v1.30.0 (`gorm.io/gorm`) con driver `gorm.io/driver/postgres` (pgx v5).
- **Seguridad y Criptografía:** `golang-jwt/jwt/v5` y `golang.org/x/crypto` (bcrypt).
- **Documentación API:** Swagger / OpenAPI mediante `swaggo/swag` y `swaggo/gin-swagger`.
- **Frontend Core:** Angular v20, RxJS v7.8, PrimeNG v20, TailwindCSS v4.
- **Integraciones Externas:** URLs embebidas/referenciadas de GitHub, LinkedIn y plataformas de video (YouTube/Vimeo).

## Data Stores

- **Base de Datos Principal:** PostgreSQL con extensión `uuid-ossp` habilitada para identificadores UUID v4.
- **Esquema Relacional:**
  - `users`: ID (UUID PK), Name, Username (Unique), Email (Unique), Password (hash), AvatarURL, BannerURL, ProfilePicture, GithubURL, LinkedinURL, PitchVideo, Headline, Bio, Seniority, EnglishLevel, EFSetScore, Points, Role, Country, Availability, Interests, Status, Timestamps.
  - `posts`: ID (UUID PK), Title, Content, Image, Type, Tags, Status, MediaURL, UserID (FK), VotesCount, Timestamps.
  - `comments`: ID (UUID PK), UserID (FK), PostID (FK), Content, Status, Timestamps.
  - `posts_votes`: Composite (UserID, PostID), CreatedAt.
- **Almacenamiento de Archivos:** Directorio local `./uploads` servido como ruta estática en `/uploads`.

## Infrastructure and Deployment Signals

- **Contenedores y Orquestación:** Configuración para ejecución local mediante variables de entorno (`.env` con `godotenv`).
- **CORS:** Configurado para admitir el origen definido en `FRONTEND_URL`.
- **Documentación Swagger:** Expuesta en `GET /swagger/*any` consumiendo `docs/swagger.json`.

## Implicit Decisions (candidates)

- **Arquitectura Backend:** Monolito modular por capas (Handler -> Service -> Repository -> GORM) en lugar de microservicios o arquitectura hexagonal estricta.
- **Identificadores:** Adopción uniforme de UUID v4 en todas las entidades primarias.
- **Borrado Lógico:** Manejo de estados (`enums.DELETED`) en lugar de `gorm.DeletedAt` con timestamp de base de datos.
- **Tags de Contenido:** Almacenamiento desnormalizado como string separado por comas (`Tags string`).
- **Gestión de Estado Frontend:** Adopción de `@ngrx/signals` (SignalStore) en lugar de NgRx Store clásico basado en Actions/Reducers.

## Open Questions

- [open] ¿Cuándo se migrará la capa de almacenamiento de `LocalStorage` a un driver de `S3Storage` compatible con AWS S3 / MinIO / Cloudflare R2?
- [open] ¿Se incorporará una tabla relacional normalizada de `tags` y `skills` para soportar indexación y agregaciones complejas?
- [open] ¿Se adoptará `deleted_at` para borrado lógico estándar con soporte nativo de GORM?
- [open] ¿Qué estrategia de caché (Redis / In-memory) se implementará para optimizar las consultas de la Cartelera Semanal y el Directorio?

## Areas Requiring Human Validation

1. Política de retención y limpieza de imágenes huérfanas en el almacenamiento local.
2. Definición formal del algoritmo de refresco de tokens (HttpOnly cookies vs LocalStorage en cliente).
3. Planes de infraestructura para soporte de Server-Side Rendering (SSR) o prerendering en Angular para SEO de perfiles públicos.

────────────────────────
Agent: architecture-agent

Produced:
knowledge/tech/current-state.md

Next:
decision-agent
roadmap-agent
────────────────────────
