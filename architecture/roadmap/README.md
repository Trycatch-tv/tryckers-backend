# Tryckers Platform Roadmap

Functionality -> Usability -> Experience -> Scalability

This roadmap tracks the platform vision, implementation status, vertical slices, and release tags.

Status legend:

- Complete: implemented in the current backend/frontend codebase.
- Partial: implemented at a basic level, but needs hardening or UX/product completion.
- Pending: not implemented yet.

## Tagging Rules

Completed vertical slices are tagged after merge to `main`.

Format:

```text
vs/<vs-name>/v<major>.<minor>.<patch>
```

Version guidance:

- Patch: implementation detail or bug fix within an existing capability.
- Minor: new vertical slice or backwards-compatible API capability.
- Major: breaking API, persistence, or architecture contract change.

## Vertical Slice Log

| Date | VS | Branch | Tag | OpenSpec | Summary | Status |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-05-26 | `architecture-sdd-foundation` | `dev` -> `main` | TBD | `architecture/openspec/changes/architecture-sdd-foundation` | Establish SDD, OpenSpec workspace, roadmap, and project agents. | Complete |

## Current Implemented Capabilities

| Area | Capability | Status |
| --- | --- | --- |
| Auth | User registration | Complete |
| Auth | Login with JWT | Complete |
| Auth | Refresh token endpoint | Partial |
| Auth | Frontend session persistence | Partial |
| Auth | Auth guards for public/private routes | Complete |
| Auth | Bearer token interceptor | Complete |
| Auth | Automatic refresh on 401 | Partial |
| Users | User listing / directory base | Complete |
| Users | Profile by username | Complete |
| Users | Profile page UI | Partial |
| Posts | Create post | Complete |
| Posts | List posts | Complete |
| Posts | Post detail | Complete |
| Posts | Update post | Complete |
| Posts | Logical delete by status | Partial |
| Posts | Posts by user | Complete |
| Posts | Tags as comma-separated data | Partial |
| Posts | Video post media URL validation | Partial |
| Votes | Vote posts | Complete |
| Discovery | Weekly popular posts / cartelera | Complete |
| Comments | Create comments | Complete |
| Comments | List comments by post | Complete |
| Comments | Update comments | Complete |
| Comments | Logical delete comments | Partial |
| Frontend UX | Skeleton loading in key pages | Partial |
| Frontend UX | Toast notifications | Partial |
| Frontend UX | Basic UX metrics service | Partial |
| Documentation | SDD architecture workspace | Complete |
| Documentation | OpenSpec change workflow | Complete |
| Documentation | VS and branch agent runbooks | Complete |
| Documentation | Swagger docs | Partial |

## Stage 1 - MVP Consolidation

Objective: close the MVP with stability, UX coherence, functional onboarding, and a real usage flow.

### Authentication & Security Hardening

| Capability | Backend | Frontend | Status |
| --- | --- | --- | --- |
| Refresh token rotation | Required | Token update flow exists | Partial |
| Token revocation | Required | Logout clears local state | Partial |
| Rate limiting | Required | N/A | Pending |
| Anti brute-force protection | Required | N/A | Pending |
| Email verification | Required | Verification UX required | Pending |
| Forgot/reset password | Required | Reset flow required | Pending |
| HTML/Markdown content sanitization | Required | Preview/editor sanitization required | Pending |
| Centralized DTO validation | Basic binding exists | Form validation exists | Partial |
| Basic action audit | Required | N/A | Pending |
| Global session expiration handling | N/A | Interceptor handles 401 | Partial |
| Silent refresh UX | N/A | Interceptor retry exists | Partial |
| Robust persistent auth state | N/A | localStorage state exists | Partial |
| Global loader states | N/A | Required | Pending |
| Skeleton loading | N/A | Present in some screens | Partial |

### User Experience Improvements

| Area | Capability | Status |
| --- | --- | --- |
| Profile | Upload avatar | Pending |
| Profile | Cover/banner profile | Pending |
| Profile | Skills tags | Pending |
| Profile | Tech stack visual | Pending |
| Profile | Social badges | Pending |
| Profile | GitHub preview | Pending |
| Profile | LinkedIn preview | Pending |
| Profile | Video pitch embed | Partial |
| Profile | Visual work availability | Pending |
| Directory | Advanced search | Pending |
| Directory | Persistent URL filters | Pending |
| Directory | Sorting | Pending |
| Directory | Infinite scroll or pagination | Pending |
| Directory | Quick profile cards | Partial |

### Posts System Improvements

| Capability | Backend | Frontend | Status |
| --- | --- | --- | --- |
| Markdown support | Required | Required | Pending |
| Rich content sanitization | Required | Required | Pending |
| Trending score algorithm | Cartelera by weekly votes exists | Cartelera UI exists | Partial |
| Real soft delete with timestamps | Status-based delete exists | Delete action exists | Partial |
| Friendly slugs | Required | Required | Pending |
| Improved editor | N/A | Basic modal/form exists | Partial |
| Draft autosave | Required | Required | Pending |
| Post preview | N/A | Required | Pending |
| Empty states | N/A | Some screens have empty states | Partial |
| Relative dates | N/A | Required | Pending |
| Share actions | N/A | Required | Pending |
| Bookmarking | Required | Required | Pending |

## Stage 2 - Networking Platform

Objective: transform the directory into a platform for real connections.

### Networking Features

| Capability | Status |
| --- | --- |
| Follow users | Pending |
| Connection requests | Pending |
| Basic direct messaging | Pending |
| Matchmaking by skills/interests | Pending |
| Open to work | Pending |
| Looking for co-founder | Pending |
| Available for freelance | Pending |

### Intelligent Discovery

| Capability | Status |
| --- | --- |
| Basic recommendation engine | Pending |
| Similar profiles | Pending |
| Suggested connections | Pending |
| Related posts | Pending |
| Recommended Tryckers section | Pending |
| Suggested collaborators | Pending |
| Smart onboarding questions | Pending |

### Opportunities / Recruiters

| Capability | Status |
| --- | --- |
| Job/opportunity publishing | Pending |
| Recruiter profiles | Pending |
| Simple candidate pipeline | Pending |
| Talent search filters | Pending |
| Bookmark candidates | Pending |
| Contact requests | Pending |
| One-click apply | Pending |
| Public profile share page | Pending |
| Resume generator | Pending |
| Recruiter visibility toggle | Pending |

## Stage 3 - Community Ecosystem

Objective: turn Tryckers into a living community ecosystem.

### Gamification

| Capability | Status |
| --- | --- |
| XP system | Pending |
| Reputation | Pending |
| Badges | Pending |
| Contributor levels | Pending |
| Streaks | Pending |
| Leaderboards | Pending |
| Community challenges | Pending |

### Open Source Collaboration

| Capability | Status |
| --- | --- |
| GitHub integration | Pending |
| Show PRs and contributions | Pending |
| OSS contributor badge | Pending |
| Project showcase | Pending |
| Community teams | Pending |

### Learning Layer

| Capability | Status |
| --- | --- |
| Workshops | Pending |
| Mentorships | Pending |
| Pair programming | Pending |
| Learning paths | Pending |
| Skill progression | Pending |
| Community certifications | Pending |

## Stage 4 - Platform & Scale

Objective: scale the platform technically and socially.

### Technical Evolution

| Capability | Backend | Frontend | Status |
| --- | --- | --- | --- |
| Modular monolith optimization | Architecture documented | N/A | Partial |
| Caching layer | Required | N/A | Pending |
| Queue system | Required | N/A | Pending |
| Search engine | Required | Search UI required | Pending |
| Realtime layer | Required | Realtime updates required | Pending |
| Event-driven architecture | Future architecture | N/A | Pending |
| SSR / SEO optimization | N/A | Required | Pending |
| PWA support | N/A | Required | Pending |
| Offline mode | N/A | Required | Pending |
| Microfrontend exploration | N/A | Future architecture | Pending |

### Analytics Platform

| Capability | Status |
| --- | --- |
| Profile views | Pending |
| Recruiter impressions | Pending |
| Post engagement | Partial |
| Opportunity CTR | Pending |
| Network growth | Pending |
| Community growth | Pending |
| Most active skills | Pending |
| Hiring trends | Pending |
| Country distribution | Pending |
| Contribution analytics | Pending |

## Stage 5 - Tryckers Network Vision

Objective: become a LATAM tech network.

### Ecosystem Vision

| Capability | Status |
| --- | --- |
| AI talent matching | Pending |
| AI profile optimization | Pending |
| AI recruiter assistant | Pending |
| AI career recommendations | Pending |
| AI-generated portfolios | Pending |

### Internationalization

| Capability | Status |
| --- | --- |
| Multi-language support | Pending |
| Country communities | Pending |
| Regional events | Pending |
| Localized feeds | Pending |

### Enterprise Layer

| Capability | Status |
| --- | --- |
| Company pages | Pending |
| Team hiring dashboards | Pending |
| Employer branding | Pending |
| Sponsorships | Pending |
| Verified recruiters | Pending |

## Remaining Vertical Slice Backlog

Each vertical slice should deliver real functional value, cross backend + frontend + UX when needed, remain deployable, and be documented through SDD/OpenSpec.

### Stage 1 - MVP Consolidation VS

#### Authentication & Security

| VS | Includes | Status |
| --- | --- | --- |
| `vertical-slice-refresh-token-hardening` | Refresh token rotation, token invalidation, session refresh strategy, frontend retry stabilization. | Pending |
| `vertical-slice-auth-session-persistence` | Robust auth hydration, restore session on refresh, session expiration UX, global auth store stabilization. | Pending |
| `vertical-slice-email-verification` | Verification tokens, verification endpoint, verification email template, frontend verification screen. | Pending |
| `vertical-slice-forgot-password-flow` | Forgot password request, reset token generation, reset form, reset validation. | Pending |
| `vertical-slice-rate-limit-and-anti-bruteforce` | Rate limiting middleware, login attempt throttling, IP/user tracking. | Pending |
| `vertical-slice-centralized-validation-and-errors` | DTO validation layer, unified error responses, validation middleware, frontend form error mapping. | Pending |
| `vertical-slice-audit-log-foundation` | Basic action logs, user activity tracking, sensitive action registry. | Pending |

#### User Experience

| VS | Includes | Status |
| --- | --- | --- |
| `vertical-slice-profile-avatar-and-banner` | Avatar upload, cover upload, storage integration, profile media rendering. | Pending |
| `vertical-slice-profile-tech-identity` | Skills tags, tech stack visualization, availability badges, social badges. | Pending |
| `vertical-slice-social-profile-preview` | GitHub preview, LinkedIn preview, video pitch embed stabilization. | Pending |
| `vertical-slice-directory-search-and-filtering` | Search engine, filter by skills/country/seniority, query params persistence, sorting. | Pending |
| `vertical-slice-directory-pagination-and-cards` | Infinite scroll or pagination, quick cards, lightweight previews. | Pending |

#### Posts System

| VS | Includes | Status |
| --- | --- | --- |
| `vertical-slice-markdown-post-system` | Markdown rendering, Markdown sanitization, rich content support. | Pending |
| `vertical-slice-post-editor-improvements` | Better editor UX, preview mode, empty states, relative dates. | Pending |
| `vertical-slice-post-drafts-and-autosave` | Draft persistence, autosave strategy, recover draft flow. | Pending |
| `vertical-slice-post-sharing-and-bookmarks` | Share links, bookmark system, saved posts UI. | Pending |
| `vertical-slice-post-slugs-and-seo-foundation` | Friendly slugs, SEO-ready URLs, share metadata. | Pending |
| `vertical-slice-trending-algorithm-v2` | Better ranking formula, time decay scoring, engagement weighting. | Pending |

### Stage 2 - Networking Platform VS

#### Connections & Networking

| VS | Includes | Status |
| --- | --- | --- |
| `vertical-slice-follow-system` | Follow/unfollow, followers/following, feed relationships. | Pending |
| `vertical-slice-connection-requests` | Mutual connections, accept/reject flows. | Pending |
| `vertical-slice-open-to-work-and-freelance` | Availability flags, co-founder status, freelance visibility. | Pending |
| `vertical-slice-direct-messaging-foundation` | Conversations, basic messaging, notification hooks. | Pending |

#### Intelligent Discovery

| VS | Includes | Status |
| --- | --- | --- |
| `vertical-slice-profile-recommendation-engine` | Similar profiles, suggested collaborators, related members. | Pending |
| `vertical-slice-related-posts-and-content` | Related posts, skill similarity scoring. | Pending |
| `vertical-slice-smart-onboarding` | Guided onboarding, interest capture, personalized recommendations. | Pending |

#### Recruiter Platform

| VS | Includes | Status |
| --- | --- | --- |
| `vertical-slice-opportunity-posting-system` | Job posts, opportunity forms, recruiter publishing flow. | Pending |
| `vertical-slice-recruiter-profiles` | Recruiter identity, company association, recruiter visibility. | Pending |
| `vertical-slice-candidate-search-engine` | Talent filtering, search by skills, candidate cards. | Pending |
| `vertical-slice-one-click-apply` | Apply action, candidate submission, recruiter notifications. | Pending |
| `vertical-slice-recruiter-pipeline-foundation` | Candidate stages, recruiter workflow, bookmark candidates. | Pending |
| `vertical-slice-public-profile-share-pages` | Public sharable URLs, recruiter-ready pages, SEO foundation. | Pending |

### Stage 3 - Community Ecosystem VS

#### Gamification

| VS | Includes | Status |
| --- | --- | --- |
| `vertical-slice-reputation-system` | XP, reputation scoring, activity weighting. | Pending |
| `vertical-slice-badges-and-levels` | Contributor badges, levels, milestones. | Pending |
| `vertical-slice-community-challenges` | Weekly challenges, community participation loops. | Pending |

#### Open Source Layer

| VS | Includes | Status |
| --- | --- | --- |
| `vertical-slice-github-integration` | GitHub account linking, public contributions sync. | Pending |
| `vertical-slice-project-showcase` | OSS projects, featured community projects. | Pending |

#### Learning Layer

| VS | Includes | Status |
| --- | --- | --- |
| `vertical-slice-mentorship-system` | Mentor/mentee flows, session scheduling, mentorship profiles. | Pending |
| `vertical-slice-learning-paths` | Skill progression, structured learning tracks. | Pending |

### Stage 4 - Platform & Scale VS

| VS | Includes | Status |
| --- | --- | --- |
| `vertical-slice-search-engine-foundation` | Full text search, search indexing. | Pending |
| `vertical-slice-realtime-foundation` | Realtime notifications, live updates. | Pending |
| `vertical-slice-analytics-foundation` | Profile views, engagement analytics, recruiter impressions. | Pending |
| `vertical-slice-caching-layer` | Redis layer, hot endpoints optimization. | Pending |
| `vertical-slice-queue-processing-foundation` | Async jobs, email queues, notification queues. | Pending |
| `vertical-slice-ssr-and-seo-optimization` | Angular SSR, meta rendering, SEO optimization. | Pending |

### Stage 5 - Tryckers Network Vision VS

| VS | Includes | Status |
| --- | --- | --- |
| `vertical-slice-ai-talent-matching` | AI recommendations, recruiter matching. | Pending |
| `vertical-slice-ai-profile-optimization` | AI-generated suggestions, profile improvements. | Pending |
| `vertical-slice-multi-language-foundation` | i18n setup, language switching. | Pending |
| `vertical-slice-company-pages` | Company profiles, employer branding. | Pending |

## Recommended Next VS Order

| Order | VS | Why |
| --- | --- | --- |
| 1 | `vertical-slice-profile-avatar-and-banner` | Improves profile trust and visible product quality quickly. |
| 2 | `vertical-slice-directory-search-and-filtering` | Makes the directory useful as the member base grows. |
| 3 | `vertical-slice-opportunity-posting-system` | Starts recruiter value and marketplace utility. |
| 4 | `vertical-slice-centralized-validation-and-errors` | Reduces product friction and stabilizes API/frontend feedback. |
| 5 | `vertical-slice-email-verification` | Strengthens account trust and security. |
| 6 | `vertical-slice-post-editor-improvements` | Improves creation quality and community publishing. |
| 7 | `vertical-slice-post-drafts-and-autosave` | Protects user work and makes posting feel reliable. |
| 8 | `vertical-slice-follow-system` | Creates the first real networking graph. |
| 9 | `vertical-slice-profile-recommendation-engine` | Improves discovery and organic connection loops. |
| 10 | `vertical-slice-analytics-foundation` | Gives the platform visibility into engagement and recruiter value. |

This order unlocks better UX, real networking, organic growth, recruiter value, and early differentiation.

## Recommended Immediate Priority

### Priority 1

| Capability | Status |
| --- | --- |
| Upload avatar/banner | Pending |
| Search/filter users | Pending |
| Recruiter opportunities module | Pending |
| Better profile UX | Partial |
| Skeleton/loading states | Partial |
| Email verification | Pending |
| Forgot password | Pending |

### Priority 2

| Capability | Status |
| --- | --- |
| Recommendation engine | Pending |
| Follow system | Pending |
| Notifications | Pending |
| Bookmarks | Pending |
| Drafts | Pending |
| Gamification | Pending |

### Priority 3

| Capability | Status |
| --- | --- |
| Realtime | Pending |
| Messaging | Pending |
| Analytics | Pending |
| AI layer | Pending |

## Strategic Recommendation

Tryckers should not initially position itself as another LinkedIn.

The strongest opportunity is to become the platform where developers build real reputation while learning and collaborating.

That positioning improves engagement, community growth, organic network effects, and recruiter value.

## North Star Vision

Tryckers can evolve toward:

- Talent network.
- OSS collaboration hub.
- Learning ecosystem.
- Technical reputation platform.
- LATAM developer graph.
