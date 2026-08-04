---
type: product
project_state: ai-assisted
generated_by: kaddo-bootstrap
template_version: 1
refined_by: product-agent
---

> Idioma del proyecto: **español**. Escribe este conocimiento en español. Mantén en inglés el código, los nombres de archivo, los comandos y las claves de configuración.

# Product Context

## Existing Product Behavior

Tryckers cuenta actualmente con una base funcional completa para la gestión de identidad, perfiles de miembros, publicaciones técnicas, debates mediante comentarios, sistema de votos y cartelera comunitaria.

### Authentication and Security
El producto permite:
- Registrar nuevos usuarios con validación de unicidad de email y username.
- Iniciar sesión emitiendo tokens JWT con roles (`Member`, `Admin`).
- Proteger rutas en backend mediante `AuthMiddleware` y `RoleMiddleware`.
- Proteger rutas en frontend mediante `AuthenticatedGuard` y `NotAuthenticatedGuard`.
- Interceptar peticiones HTTP en frontend con Bearer Token y soporte inicial de renovación de token (`/api/v1/refresh-token`).

### User Profiles and Directory
El producto permite:
- Listar miembros en el directorio del dashboard (`GET /api/v1/users`).
- Consultar el perfil público de cualquier miembro mediante su `@username` (`GET /api/v1/perfil/:username`).
- Almacenar datos profesionales en el perfil: titular (*headline*), biografía, enlaces a GitHub y LinkedIn, video de presentación (*pitch video*), seniority, nivel de inglés, disponibilidad e intereses.
- Cargar, actualizar y eliminar avatar y banner decorativo utilizando almacenamiento local (`POST /api/v1/users/me/avatar`, `POST /api/v1/users/me/banner`).

### Content and Discussions (Posts and Comments)
El producto permite:
- Crear publicaciones con título, cuerpo de texto, imagen de portada, tags separados por coma y enlaces multimedia (`POST /api/v1/posts`).
- Listar todas las publicaciones activas y consultar detalles de un post específico.
- Listar los posts pertenecientes a un usuario determinado (`GET /api/v1/users/:ownerId/posts`).
- Editar posts existentes validando la autoría del usuario.
- Eliminar publicaciones de forma lógica mediante el estado (`status: deleted`).
- Comentar en posts existentes (`POST /api/v1/comments`, `GET /api/v1/posts/:id/comments`).
- Editar y eliminar comentarios propios.

### Social Curation and Engagement
El producto permite:
- Votar publicaciones con comportamiento tipo *toggle* (votar / quitar voto) mediante `POST /api/v1/posts/:id/vote`.
- Exponer la **Cartelera Semanal** (`GET /api/v1/cartelera`) que lista las publicaciones más votadas y populares del periodo.

---

## Main User Flows

### 1. User Registration Flow
1. El visitante accede a la ruta `/auth/register`.
2. Completa los campos requeridos (nombre, username, email, contraseña).
3. El frontend valida la estructura del formulario y envía `POST /api/v1/register`.
4. El backend valida el DTO, genera el hash de la contraseña y crea el registro en PostgreSQL.
5. El sistema responde con éxito y redirige al login.

### 2. Login and Session Flow
1. El usuario ingresa sus credenciales en `/auth/login`.
2. El frontend envía `POST /api/v1/login`.
3. El backend verifica las credenciales y genera un token JWT.
4. El frontend almacena el token en `AuthStore` y redirige al feed/dashboard (`/home`).
5. Las peticiones autenticadas adjuntan el token en el header `Authorization: Bearer <token>`.

### 3. Profile Viewing and Media Upload Flow
1. El usuario navega a su perfil (`/profile/:username`) o al de otro miembro.
2. El frontend solicita los datos del perfil (`GET /api/v1/perfil/:username`) y sus posts asociados (`GET /api/v1/users/:ownerId/posts`).
3. Para actualizar su imagen o banner, el usuario selecciona un archivo y envía `POST /api/v1/users/me/avatar` o `POST /api/v1/users/me/banner` (multipart/form-data).
4. El backend almacena el archivo en `/uploads` y actualiza la URL en la tabla `users`.

### 4. Create, Edit and Delete Post Flow
1. El usuario autenticado redacta un post (título, contenido, tags, media).
2. El frontend envía `POST /api/v1/posts`.
3. El backend valida el payload y crea el post en estado activo.
4. Para modificarlo, el autor envía `PUT /api/v1/posts`.
5. Para eliminarlo, el autor o un administrador envía `DELETE /api/v1/posts/:id`, marcando el estado como `deleted`.

### 5. Vote and Billboard Flow
1. Un usuario visualiza un post en el feed o detalle.
2. Hace clic en el botón de votar; se envía `POST /api/v1/posts/:id/vote`.
3. El backend registra o revoca el voto y actualiza el contador.
4. Al visitar `/cartelera`, el backend calcula y retorna los posts con mayor puntuación de la semana.

---

## Planned Product Behavior

### 1. Advanced Profile and Talents
- **Skills Catalog:** Catálogo normalizado de habilidades técnicas con autocompletado y niveles de dominio.
- **GitHub Integration:** Sincronización de contribuciones, repositorios destacados y métricas de PRs.
- **Video Pitch Embedded Player:** Reproductor integrado para visualizar el video de presentación del candidato.

### 2. Recruiter Portal and Opportunities
- **Módulo de Oportunidades:** Publicación de ofertas de empleo (Full-time, Freelance, Pasantías, Co-founder).
- **Talent Discovery Search:** Filtros avanzados en el directorio por país, seniority, nivel de inglés y stack tecnológico.
- **One-Click Apply:** Postulación directa a vacantes utilizando el perfil de Tryckers como CV vivo.

### 3. Content and Community Evolution
- **Markdown Rich Editor:** Editor de texto enriquecido con soporte de código, vista previa en tiempo real y sanitización XSS.
- **Threaded Comments:** Soporte para respuestas anidadas en debates complejos.
- **Trending Algorithm:** Cálculo de popularidad con fórmula de decaimiento temporal tipo Hacker News.

### 4. Infrastructure and Scalability
- **Cloud Object Storage:** Migración del almacenamiento de medios de filesystem local a Amazon S3 / Cloudflare R2.
- **Public Profiles and SSR/SEO:** Habilitar visualización de perfiles sin login obligatorio y meta tags dinámicos para indexación y compartición en redes sociales.

---

## Open Questions

- [open] ¿Cuál será la política de visibilidad pública de los perfiles para visitantes no autenticados?
- [open] ¿Qué límites de almacenamiento y formatos se establecerán para videos de presentación subidos directamente?
- [open] ¿Cómo se estructurará el modelo de datos para las ofertas de empleo y postulaciones?
- [open] ¿Qué eventos alimentarán el sistema de reputación y puntos de los miembros?

────────────────────────
Agent: product-agent

Produced:
knowledge/product/product.md

Next:
architecture-agent
────────────────────────