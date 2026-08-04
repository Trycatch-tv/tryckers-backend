---
type: business
project_state: ai-assisted
generated_by: kaddo-bootstrap
template_version: 1
refined_by: business-agent
---

> Idioma del proyecto: **español**. Escribe este conocimiento en español. Mantén en inglés el código, los nombres de archivo, los comandos y las claves de configuración.

# Business Context

## Problem Statement

El ecosistema de talento tecnológico en Latinoamérica enfrenta una desconexión crítica entre la formación de los desarrolladores y los procesos reales de descubrimiento y contratación:

1. **Falta de reputación técnica verificable:** Los currículums tradicionales y perfiles estáticos en redes profesionales no reflejan la habilidad real de resolución de problemas, la calidad de código ni la capacidad de colaboración efectiva.
2. **Barrera de entrada para talento emergente y juniors:** Los desarrolladores formados de manera autodidacta o en bootcamps carecen de canales creíbles para demostrar experiencia práctica en entornos de equipo o proyectos de código abierto del mundo real.
3. **Invisibilidad del trabajo colaborativo:** Las contribuciones en comunidades, revisiones de código, mentorías y discusiones técnicas de alto valor suelen quedar dispersas y no se consolidan en una identidad profesional auditable.
4. **Fricción para reclutadores y líderes técnicos:** Los reclutadores invierten un tiempo excesivo filtrando postulaciones sin poder evaluar señales directas de habilidad técnica, dominio de inglés y actividad comprobada en comunidad.

## Value Proposition

Tryckers es una plataforma comunitaria y red de talento técnico para Latinoamérica que une el aprendizaje práctico, la colaboración open source y la visibilidad profesional.

- **Para los desarrolladores:** Un espacio para construir una identidad profesional viva mediante publicaciones técnicas, debates, proyectos y contribuciones reales que demuestran su competencia y seniority.
- **Para la comunidad:** Un ecosistema de aprendizaje colaborativo y mentoría entre pares impulsado por los miembros de TryCatch.tv.
- **Para reclutadores y empresas:** Un canal directo para descubrir talento validado a través de evidencia tangible (código, proyectos, video pitches, nivel de inglés y actividad comunitaria) en lugar de CVs autodeclarados.

## Target Users and Personas

### Member (Desarrollador / Miembro Comunitario)
- **Perfil:** Desarrollador de software, estudiante o profesional tecnológico en LATAM.
- **Objetivos principales:**
  - Construir un perfil profesional con enlaces a proyectos, nivel de inglés, video de presentación y disponibilidad laboral.
  - Publicar artículos técnicos, tutoriales y reflexiones para construir marca personal.
  - Interactuar con pares mediante comentarios y votos en publicaciones.
  - Ser descubierto por reclutadores y potenciales colaboradores.
  - Descubrir y conectar con otros profesionales del ecosistema.

### Recruiter (Reclutador / Headhunter)
- **Perfil:** Profesional de atracción de talento o líder técnico que busca contratar desarrolladores en LATAM.
- **Objetivos principales:**
  - Filtrar y descubrir candidatos según habilidades, país, seniority y estado de disponibilidad (*Open to work*, *Freelance*, etc.).
  - Analizar perfiles enriquecidos con señales objetivas (video pitch, inglés, intereses y actividad).
  - Contactar directamente a miembros disponibles.

### Contributor (Contribuidor Open Source)
- **Perfil:** Miembro activo interesado en colaborar en el desarrollo de la plataforma Tryckers y herramientas comunitarias.
- **Objetivos principales:**
  - Participar en el ciclo de vida del software mediante issues, PRs, revisiones y testing.
  - Validar habilidades de trabajo en equipo y arquitectura en un proyecto de producción.
  - Obtener reconocimiento y reputación por sus aportes de código y documentación.

### Mentor / Peer Guide
- **Perfil:** Desarrollador senior o especialista técnico dispuesto a guiar a otros miembros.
- **Objetivos principales:**
  - Compartir conocimiento a través de retroalimentación en posts y código.
  - Acompañar a miembros en retos comunitarios y rutas de aprendizaje.

### Company Representative (Representante de Empresa)
- **Perfil:** Organización que busca posicionar su marca empleadora o publicar oportunidades laborales y retos técnicos.
- **Objetivos principales:**
  - Publicar vacantes orientadas a la comunidad.
  - Evaluar candidatos con señales directas de participación.

### Administrator (Moderador / Operador de Plataforma)
- **Perfil:** Equipo responsable de la gobernanza, seguridad y operación de la comunidad.
- **Objetivos principales:**
  - Moderar contenido y reportes para asegurar un espacio constructivo y libre de spam.
  - Gestionar roles de usuario y permisos administrativos.
  - Supervisar la salud y métricas de interacción en la plataforma.

## Core Business Rules

### 1. Identidad y Registro
- **BR-IAM-001:** El registro en la plataforma requiere un correo electrónico único y una contraseña segura procesada con algoritmo criptográfico de hash.
- **BR-IAM-002:** Cada usuario debe poseer un `username` único en el sistema que actuará como identificador en su URL pública de perfil.
- **BR-IAM-003:** El acceso a operaciones privadas requiere autenticación mediante tokens de sesión válidos.
- **BR-IAM-004:** La plataforma debe soportar al menos los roles `Member` y `Admin`, condicionando el acceso a endpoints privilegiados.

### 2. Identidad Profesional y Perfil
- **BR-PRF-001:** Todo miembro registrado posee un perfil profesional editable.
- **BR-PRF-002:** El perfil puede contener información profesional: titular (*headline*), biografía, enlaces a GitHub/LinkedIn, enlace a video de presentación (*pitch video*), seniority, nivel de inglés y disponibilidad laboral.
- **BR-PRF-003:** Los usuarios pueden personalizar su perfil cargando y eliminando imágenes para su avatar y banner decorativo.
- **BR-PRF-004:** La información pública del perfil debe omitir cualquier dato sensible o credencial de seguridad.

### 3. Publicaciones y Contenido Técnico
- **BR-PST-001:** Únicamente los usuarios autenticados pueden redactar y publicar posts.
- **BR-PST-002:** Cada publicación debe contener título y cuerpo de contenido, pudiendo incluir tags temáticos y enlaces a medios (imágenes o videos).
- **BR-PST-003:** Un post solo puede ser editado por su autor original.
- **BR-PST-004:** La eliminación de publicaciones debe realizarse mediante borrado lógico, garantizando que el contenido eliminado deje de ser visible en listados públicos.
- **BR-PST-005:** Los administradores pueden eliminar publicaciones que infrinjan las normas de la comunidad.

### 4. Interacción, Votos y Cartelera
- **BR-INT-001:** Un usuario autenticado puede emitir un único voto de respaldo (*upvote*) por publicación y puede retirar su voto en cualquier momento.
- **BR-INT-002:** Los votos acumulados determinan la visibilidad de las publicaciones en la cartelera semanal de contenido destacado.
- **BR-INT-003:** Los usuarios autenticados pueden registrar comentarios en publicaciones públicas para formular preguntas o aportar al debate técnico.
- **BR-INT-004:** Los comentarios solo pueden ser modificados o eliminados por su respectivo autor o por administradores.

### 5. Privacidad y Seguridad
- **BR-SEC-001:** Ningún endpoint público ni payload de API debe exponer contraseñas, hashes, tokens de sesión o datos internos de auditoría no autorizados.
- **BR-SEC-002:** La visualización de datos de contacto debe respetar la configuración de privacidad establecida por el usuario.

## Constraints

- **Restricción de Recursos:** Desarrollo impulsado principalmente por comunidad open source, requiriendo arquitecturas modulares y bajo costo operativo inicial.
- **Modelo de Despliegue:** Arquitectura desacoplada en backend (Go/Gin) y frontend (Angular/SPA) con API REST versionada (`/api/v1`).
- **Persistencia de Archivos:** En fases iniciales el almacenamiento multimedia utiliza el sistema de archivos local (`/uploads`), requiriendo diseño compatible para migración futura a almacenamiento de objetos (S3/Cloud Storage).
- **Idioma Principal:** La interfaz y contenidos comunitarios están dirigidos prioritariamente a la comunidad hispanohablante de LATAM.

## Glossary

- **Trycker:** Miembro de la comunidad TryCatch.tv registrado en la plataforma.
- **Cartelera:** Sección semanal que agrupa y destaca las publicaciones con mayor engagement y respaldo de la comunidad.
- **Pitch Video:** Video corto (1-2 minutos) donde el profesional presenta su perfil, experiencia e intereses laborales.
- **Open to Work:** Estado de disponibilidad que indica que el miembro busca activamente empleo formal.
- **Vertical Slice:** Metodología de entrega de software en la que cada incremento contiene UI, API, lógica de negocio y persistencia completa.

## Assumptions

- **A-01:** Los miembros valoran una alternativa a LinkedIn enfocada en validación técnica y comunidad abierta.
- **A-02:** Los reclutadores están dispuestos a usar filtros de habilidades e inglés para encontrar candidatos con perfiles activos.
- **A-03:** El sistema de votos y cartelera semanal incentiva la publicación de contenido técnico de calidad sin requerir algoritmos complejos de curación inicial.

## Open Questions

- [open] ¿Cuáles serán las reglas y fórmulas exactas para el algoritmo de trending de la cartelera con decaimiento por tiempo?
- [open] ¿Qué verificación adicional requerirán las cuentas de tipo `Recruiter` antes de contactar candidatos masivamente?
- [open] ¿Cómo se integrará la validación automática de pull requests y badges de GitHub en el perfil de usuario?
- [open] ¿Qué límites de tamaño, formatos y cuotas por usuario aplicarán para avatares y banners?
- [open] ¿Qué mecanismo de moderación comunitaria (reporte de abusos, flags) se implementará para el contenido generado por usuarios?

────────────────────────
Agent: business-agent

Produced:
knowledge/business/business.md

Next:
product-agent
────────────────────────