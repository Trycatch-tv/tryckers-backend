---
type: adr
status: accepted
date: 2026-08-04
created_from: knowledge/tech/discovery/decision-candidates.md
governed_code:
  - src/internal/services/storage/**
---

# ADR-001 — Adopción de Interfaz de Almacenamiento con Driver Compatible con S3

## Context
Actualmente la plataforma almacena avatares y banners decorativos en el sistema de archivos local (`./uploads`) mediante la implementación `LocalStorage`. En despliegues de producción basados en contenedores efímeros o arquitecturas multi-instancia, el almacenamiento en disco local genera pérdida de archivos y fallos de consistencia.

## Options Considered
1. **Permanecer en LocalStorage con volúmenes persistentes montados:** Sencillo pero complejo de orquestar en nubes serverless y poco escalable.
2. **Implementar almacenamiento en base de datos (BLOB / bytea):** Aumenta drásticamente el tamaño y coste de las copias de seguridad de PostgreSQL.
3. **Implementar adaptador S3-compatible (AWS S3 / Cloudflare R2 / MinIO):** Permite desacoplar el almacenamiento binario utilizando URLs públicas o prefirmadas, manteniendo compatibilidad con MinIO para desarrollo local sin conexión.

## Decision
Se decide implementar una interfaz desacoplada de almacenamiento (`StorageService`) con soporte para múltiples drivers:
- **`LocalStorage`:** Para entornos de desarrollo local y tests sin conexión.
- **`S3Storage`:** Para entornos de staging y producción utilizando proveedores compatibles con la API de Amazon S3 (e.g. Cloudflare R2, AWS S3 o MinIO).

## Consequences
- **Positivas:**
  - Despliegue sin estado (*stateless*) del contenedor backend.
  - Escalabilidad infinita de almacenamiento multimedia con bajo coste operativo.
  - Entrega eficiente de imágenes mediante CDNs globales.
- **Negativas / Mitigaciones:**
  - Requiere configuración de credenciales de acceso a bucket en variables de entorno (`S3_BUCKET`, `S3_REGION`, `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`).
  - Necesidad de asegurar políticas de limpieza para archivos huérfanos tras eliminaciones o reemplazos.

## Related Capabilities
- Media and Storage / Avatar and Banner Upload with Local Storage

## Related Work Items
- [open] WI-STORAGE-001
