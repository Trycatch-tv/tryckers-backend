---
type: decision-candidates
project_state: ai-assisted
refined_by: architecture-agent
---

> Idioma del proyecto: **español**. Escribe este conocimiento en español. Mantén en inglés el código, los nombres de archivo, los comandos y las claves de configuración.

# Architecture Decision Candidates (ADR Candidates)

Los siguientes puntos representan decisiones arquitectónicas implícitas observadas en la base de código o decisiones técnicas requeridas a corto plazo para ser formalizadas mediante `decision-agent`:

## 1. Candidate: S3-Compatible Cloud Storage Provider
- **Contexto:** Actualmente los avatares y banners se persisten en el disco local (`./uploads`), lo cual no escala en entornos serverless ni en contenedores efímeros.
- **Opciones:** AWS S3, Cloudflare R2, MinIO (desarrollo local / self-hosted).
- **Impacto:** Modificación del constructor en `src/internal/services/storage/` sin alterar la interfaz de consumo.

## 2. Candidate: Normalización de Habilidades y Tags
- **Contexto:** Las habilidades técnicas y tags de posts se almacenan actualmente como cadenas de texto plano separadas por comas.
- **Opciones:**
  - A) Mantener campos de texto plano indexados con PostgreSQL Full Text Search / `pg_trgm`.
  - B) Crear tablas relacionales dedicadas (`skills`, `tags`, `user_skills`, `post_tags`).
  - C) Array types nativos de PostgreSQL (`text[]`).
- **Impacto:** Facilita filtrado estricto, búsqueda facetada y estadísticas de popularidad tecnológica.

## 3. Candidate: Estrategia de Refresh Tokens y Manejo de Sesión
- **Contexto:** El endpoint de refresh token está expuesto pero requiere persistencia de hash de tokens para permitir rotación e invalidación de sesiones comprometidas.
- **Opciones:** Almacenamiento en tabla `refresh_tokens` en PostgreSQL con hash SHA-256 vs Redis con TTL automático.
- **Impacto:** Mayor seguridad y control de sesiones activas concurrentes.

## 4. Candidate: Motor de Migraciones en Base de Datos
- **Contexto:** GORM `AutoMigrate` se ejecuta al inicio de la aplicación en desarrollo.
- **Opciones:** Adoptar `golang-migrate` con archivos `.sql` versionados para control estricto de esquema en producción.
- **Impacto:** Trazabilidad determinista de cambios de base de datos en CI/CD.

────────────────────────
Agent: architecture-agent

Produced:
knowledge/tech/discovery/decision-candidates.md

Next:
decision-agent
roadmap-agent
────────────────────────
