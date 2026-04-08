# plata.fyi

[![Deploy](https://github.com/platafyi/plata.fyi/actions/workflows/deploy.yml/badge.svg)](https://github.com/platafyi/plata.fyi/actions/workflows/deploy.yml)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Анонимно споделување на плати за македонски работници. Без login, без лозинки, без e-mail адреса, 
само пополни формулар, помини Cloudflare [Turnstile](https://www.cloudflare.com/application-services/products/turnstile/) 
проверка и добиваш session token кој се чува во твојот прелистувач. Сите податоци се јавни и анонимни.

На секој deployment автоматски се објавува целосен CSV snapshot на сите плати како [GitHub Release](../../releases). Слободно преземи и користи.

Инспирирано од [levels.fyi](https://levels.fyi).

---

## Stack

- **Backend** - Go 1.26+
- **Frontend** - Next.js 14 (App Router), TypeScript, Tailwind CSS
- **База на податоци** - PostgreSQL 16
- **Bot заштита** - Cloudflare Turnstile

---

## Локален development

```bash
make dev-up       # стартувај Postgres во Docker
make dev-migrate  # изврши ги сите migrations

# во посебни терминали:
cd backend && go run ./cmd/api
cd frontend && npm run dev
```

Копирај `.env.example` во `.env` и пополни ги потребните вредности.

Turnstile верификацијата се прескокнува кога `TURNSTILE_SECRET` е празно (dev mode).

---

## Kako работи анонимноста

Приватноста е основен дизајн принцип што не треба да се наруши

### Без e-mail, без сметки

Никогаш не се собираат или чуваат e-mail адреси. Наместо класична автентикација, серверот генерира session token по успешна 
Turnstile верификација. Токенот живее во `localStorage` 30 дена. Ако некој го изгуби (нов уред, избришани податоци во прелистувачот), 
почнува нова сесија. Старите submissions се неповратливи по дизајн.

### Session tokens

При прво поднесување, backend-от генерира `randHex(32)` session token, го зачувува во `tokens` табелата и го враќа на прелистувачот. 
Следните requests се автентицираат со `Authorization: Bearer <token>`. Токените истекуваат по 30 дена. Засега нема начин за refresh на сесија.

---

## Тестови

```bash
cd backend
go test -race ./internal/...
```

За интеграциските тестови потребна е база (автоматски се прескокнуваат ако нема база):

```bash
TEST_DB_URL=postgres://platafyi:platafyi@localhost:5433/platafyi?sslmode=disable \
  go test -v -run TestSearch ./internal/database/
```

---

## Миграции

Миграциите се обични `.sql` фајлови во `backend/migrations/`, нумерирани по ред (`001_`, `002_`, ...).
Custom runner-от во `backend/cmd/migrate/` ги применува по ред и ги следи применетите фајлови во `schema_migrations` табела, 
така секој фајл се извршува точно еднаш.

Миграциите се само нанапред (up). Backwards (down) migrations не се дозволени. 
Ако треба да вратиш нешто, само направи нова forward миграција.

За да додадеш миграција:

1. Кеирај `backend/migrations/NNN_description.sql` (зголеми го префиксот).
2. Изврши `make dev-migrate`, ги прескокнува веќе извшените миграции и ја извршува новата миграција во трансакција
3. Ако миграцијата не успее, трансакцијата се враќа (rollback) и извршувањето запира.

Во продукција, migrations се извршуваат автоматски од deploy pipeline-от пред да се изврши deploy. 
Провери [deploy.yml](.github/workflows/deploy.yml).

---

## Контрибуирање

1. Клонирај го repo-то и создај branch од `main`.
2. **Постави локално.** Следи ги чекорите во [Локален development](#локален-development) погоре.
3. **Направи ги промените.** Држи ги pull request-ите фокусирани, една работа по PR.
4. **Изврши ги тестовите** пред да отвориш PR:
   ```bash
   cd backend && go test -race ./internal/...
   cd frontend && npm run lint
   ```
5. **Отвори pull request** врз `main` со јасен опис на што и зошто.

Користи [Conventional Commits](https://www.conventionalcommits.org/) за commit пораки, на пример
`feat: add company search`, `fix: handle missing token`, `docs: update contributing guide`.

Неколку работи да имаш на ум при контрибуција:

- Вклучи тестови за backend промени. Постојните тестови треба да продолжат да поминуваат.
- API-то и базата треба да бидат backwards compatible. Додавај нови опционални полиња наместо да преименуваш или бришеш постоечки.
- Никогаш не треба да се чуваат или логираат лично идентификувачки информации. Анонимноста е на прво место.
- Нови database колони бараат нумериран миграциски фајл, никогаш не менувај постоечки миграциски фајлови, креирај нови.-
- Ако додадеш пакет (dependency), изврши `go mod tidy && go mod vendor`. **A little copying is better than a little dependency** ([why](https://www.youtube.com/watch?v=PAAkCSZUG1c&t=9m28s)).
- Проектот работи на Go 1.26+ и Next.js 14.

---

## Deployment

Deployment-от е целосно автоматизиран преку [GitHub Actions](.github/workflows/deploy.yml). На секој пуш на таг `v*`:

1. **Се прави** Docker image за backend и frontend и ги поставува на GitHub Container Registry (`ghcr.io/platafyi/api`, `ghcr.io/platafyi/frontend`) тагирани со commit SHA.
2. **Се прави backup од salary_submission табелата** како CSV snapshot и се објавува како GitHub Release.
3. **Се deploy-а** со примена на Kubernetes manifests преку `kubectl apply -k k8s/` и чека rollout-от да заврши.

Манифестите се во `k8s/` (управувани со Kustomize). Deploy job-от го patch-ува image tag-от на тековниот commit SHA пред примена, 
така секој deployment е поврзан со точен build преку GIT_SHA.

## Аналитика

Проектот користи [Umami](https://umami.is/) за аналитика на посетеност. Umami е приватна, open-source алтернатива на 
Google Analytics која не собира лични податоци, не поставува cookies и не ги следи корисниците. 
Сите податоци се анонимни по дизајн. Не е потребен cookie banner.

Повеќе информации: [Umami Privacy Policy](https://umami.is/privacy)

---

Потребни GitHub secrets: `KUBECONFIG` (base64-encoded kubeconfig), `TURNSTILE_SECRET`, `DB_URL`, `IP_HMAC_SECRET`. <-- Веќе постојат како Secrets