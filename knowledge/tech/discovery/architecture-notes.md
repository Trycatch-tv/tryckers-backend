---
type: architecture-notes
project_state: ai-assisted
refined_by: architecture-agent
---

> Idioma del proyecto: **español**. Escribe este conocimiento en español. Mantén en inglés el código, los nombres de archivo, los comandos y las claves de configuración.

# Architecture Discovery Notes

## Observations from Codebase & Runtime

1. **Patrón de Capas en Backend:**
   - La API sigue la secuencia estándar: `Route -> Handler -> Service -> Repository -> GORM Database`.
   - La inyección de dependencias se realiza de forma manual en `SetupV1(r, db)` instanciando repositorios y servicios en el momento del arranque.

2. **Mecanismo de Almacenamiento de Archivos:**
   - Actualmente existe una interfaz de almacenamiento en `src/internal/services/storage/` cuya implementación activa es `LocalStorage`, guardando archivos en `./uploads`.
   - Se requiere diseñar una implementación `S3Storage` que satisfaga la misma interfaz para facilitar el cambio transparente a almacenamiento de objetos.

3. **Arquitectura de Estado en Frontend:**
   - Angular 20 utiliza Standalone Components y `AuthStore` basado en `@ngrx/signals`.
   - Las rutas privadas están protegidas por guards funcionales `canMatch: [AuthenticatedGuard]`.
   - La comunicación con la API se realiza a través de servicios inyectables tipados con TypeScript.

4. **Persistencia y Modelado de Datos:**
   - Todas las claves primarias utilizan UUID v4 generado por base de datos o por hooks de GORM (`BeforeCreate`).
   - Las migraciones se gestionan mediante `AutoMigrate` al inicio de la aplicación en desarrollo. Para producción, se evaluará una herramienta de migraciones deterministas (como `golang-migrate` o Goose).

────────────────────────
Agent: architecture-agent

Produced:
knowledge/tech/discovery/architecture-notes.md

Next:
decision-agent
roadmap-agent
────────────────────────
