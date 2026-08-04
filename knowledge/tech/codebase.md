---
type: codebase
project_state: ai-assisted
generated_by: kaddo-bootstrap
template_version: 1
refined_by: architecture-agent
---

> Idioma del proyecto: **español**. Escribe este conocimiento en español. Mantén en inglés el código, los nombres de archivo, los comandos y las claves de configuración.

# Codebase Map

## Repository Structure

```
tryckers-backend/
├── .kaddo/                  # Metadatos, configuración y empaquetado de conocimiento Kaddo
├── docs/                    # Especificaciones OpenAPI/Swagger generadas (swagger.json, swagger.yaml)
├── knowledge/               # Artefactos de conocimiento (business, product, tech, delivery, agents, skills)
├── src/
│   ├── cmd/
│   │   └── main.go          # Punto de entrada principal de la aplicación backend
│   └── internal/
│       ├── api/
│       │   ├── handlers/    # Controladores HTTP (User, Post, Comment)
│       │   ├── middlewares/ # Autenticación JWT, verificación de roles, CORS
│       │   └── routes/      # Enrutamiento centralizado de la API REST (/api/v1)
│       ├── config/          # Carga de variables de entorno y conexión con base de datos PostgreSQL
│       ├── dtos/            # Data Transfer Objects para requests y responses
│       ├── enums/           # Enumeraciones de dominio (UserRole, PostType, PostStatus, Country)
│       ├── errors/          # Manejo estructurado de errores y respuestas HTTP
│       ├── models/          # Entidades GORM y esquemas de persistencia
│       ├── repository/      # Capa de persistencia y consultas SQL/GORM
│       ├── services/        # Lógica de negocio y casos de uso
│       │   └── storage/     # Abstracción e implementación de almacenamiento de archivos
│       └── utils/           # Utilitarios criptográficos, hashing y generación de tokens JWT
├── uploads/                 # Directorio local de persistencia de medios (avatars, banners)
├── go.mod / go.sum          # Gestión de dependencias de Go
└── .env.example             # Plantilla de variables de entorno requeridas
```

## Entry Points

- **Backend Application Entry Point:** `src/cmd/main.go`
  - Inicializa configuración (`config.Load()`).
  - Conecta a PostgreSQL con GORM (`config.InitGormDB()`).
  - Habilita extensión `uuid-ossp` y ejecuta AutoMigrate (`User`, `Post`, `Comment`, `PostVote`).
  - Configura middleware CORS, rutas estáticas (`/docs`, `/uploads`), Swagger (`/swagger/*any`) y rutas de la API (`routes.SetupV1()`).
  - Inicia el servidor HTTP en el puerto configurado (`:8080`).

## Important Modules

- **`src/internal/api/routes/router.go`:** Define el catálogo de rutas públicas (`/register`, `/login`, `/refresh-token`) y rutas protegidas (`/users`, `/perfil/:username`, `/users/me/avatar`, `/users/me/banner`, `/posts`, `/cartelera`, `/comments`).
- **`src/internal/services/storage/local_storage.go`:** Interfaz y servicio de almacenamiento local de archivos en disco.
- **`src/internal/api/middlewares/auth.go`:** Intercepta requests protegidas, valida la firma del JWT y extrae los claims del usuario.

## How to Run

### Requisitos Previos
- Go 1.24+
- PostgreSQL 14+ con extensión `uuid-ossp`
- Archivo `.env` configurado con variables de conexión (`DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `PORT`, `JWT_SECRET`, `FRONTEND_URL`).

### Comandos de Ejecución
```bash
# Descargar dependencias
go mod download

# Ejecutar en modo desarrollo
go run src/cmd/main.go

# Generar / Actualizar documentación Swagger
swag init -g src/cmd/main.go -o docs
```

## How to Test

```bash
# Ejecutar suite de pruebas unitarias
go test ./... -v
```

## Open Questions

- [open] ¿Cuál será la configuración estándar para tests de integración automatizados con base de datos de prueba en CI/CD?
- [open] ¿Se agregará un `Makefile` o scripts Taskfile para estandarizar los comandos de desarrollo, formateo y testing?

────────────────────────
Agent: architecture-agent

Produced:
knowledge/tech/codebase.md

Next:
decision-agent
roadmap-agent
────────────────────────
