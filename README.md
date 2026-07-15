# GO-TODO-LIST

A server-rendered, multi-user todo application written in Go. It is **not** a JSON
REST API — the server renders HTML on the server side with
[templ](https://github.com/a-h/templ) templates and returns HTML fragments that
[HTMX](https://htmx.org/) swaps into the page, giving a single-page feel without a
JavaScript front-end framework. Users register and log in (passwords hashed with
bcrypt, sessions carried by a JWT cookie), then add, edit, toggle, and delete todos
that persist to MySQL through [GORM](https://gorm.io/). The HTTP layer is built on the
[Echo](https://echo.labstack.com/) framework.

## Features

- **User registration** with bcrypt-hashed passwords and unique-email enforcement.
- **Login / logout** issuing an HS256 JWT stored in an HTTP cookie (`token`) that
  expires after 24 hours; logout clears the cookie.
- **Authenticated home page** that loads only the logged-in user's todos.
- **Add todo** — submits via HTMX and appends a rendered todo fragment to the list.
- **Toggle complete** — flips completion state and re-renders that todo (strike-through
  styling and a filled check icon when complete).
- **Inline edit** — an Alpine.js-driven modal posts the new content and swaps the todo
  in place.
- **Delete** — with an HTMX confirmation prompt; when the last todo is removed the list
  shows a "No Todos" message.
- **Server-side rendering** of all pages and fragments via templ.
- **Production middleware**: request logging, 2 MB body limit, gzip compression, an
  in-memory rate limiter (20 req/s), and secure headers.

## How it works / Architecture

```
Browser (HTMX + Alpine.js + Tailwind, all via CDN)
    │  form posts / hx-* requests
    ▼
Echo router (routes/routes.go)
    │
    ▼
Handlers (handler/handler.go)
    │  auth: read "token" cookie ─► parse JWT ─► load user
    │  data: GORM queries
    ▼
GORM (db/connection.go) ──► MySQL (jwt_project)
    │
    ▼
templ components/pages ──► rendered HTML (full page or fragment)
    │
    ▼
Response HTML swapped into the DOM by HTMX
```

**Request flow**

1. `cmd/main.go` connects to the database, constructs an Echo app, installs
   middleware, wires routes, and listens on `:4000`.
2. `db.Connect()` opens a GORM MySQL connection and runs `AutoMigrate` for the `User`
   and `Todo` models on startup.
3. `routes.Setup` maps each URL to a handler.
4. Every authenticated handler (`Home`, `AddTodo`, `DeleteTodo`) reads the `token`
   cookie, parses and validates the JWT with the shared secret, and resolves the user
   by the `id` claim; failures redirect to `/login`.
5. Handlers query MySQL through GORM and respond by rendering a templ component. The
   `Render` helper renders a `templ.Component` into a pooled buffer and writes it as
   HTML. Mutating actions return a fragment (a single `Todo` component) that HTMX
   swaps into the existing page rather than a full reload.

**Templates** (compiled by `templ generate` into `*_templ.go`)

- `view/layout/layout.templ` — HTML shell that pulls in Tailwind, HTMX 2.0.1, and
  Alpine.js from CDNs.
- `view/pages/` — `Register`, `Login`, and `Home` full-page views.
- `view/components/TodoList.templ` — the `TodoList` container and the individual `Todo`
  component, which carry the `hx-post` / `hx-delete` attributes and the Alpine
  `x-data` edit-modal state.

## Tech stack

From `go.mod` (module `github.com/Garv2003/TODOLIST`, Go 1.22):

- **[Echo v4](https://echo.labstack.com/)** `v4.12.0` — HTTP router and middleware
- **[templ](https://github.com/a-h/templ)** `v0.2.747` — typed Go HTML templating
- **[GORM](https://gorm.io/)** `v1.25.11` with the **MySQL driver** `v1.5.7` — ORM / persistence
- **[golang-jwt/jwt v5](https://github.com/golang-jwt/jwt)** `v5.2.1` — JWT signing/parsing
- **[google/uuid](https://github.com/google/uuid)** `v1.6.0` — user and todo IDs
- **golang.org/x/crypto** (bcrypt) — password hashing
- **golang.org/x/time/rate** — token-bucket rate limiting
- **Front-end via CDN**: HTMX 2.0.1, Alpine.js 3.x, Tailwind CSS (no bundler / npm)

## Getting started

**Prerequisites**

- Go 1.22+
- MySQL, with a database named `jwt_project`
- The [`templ` CLI](https://templ.guide/) (`go install github.com/a-h/templ/cmd/templ@latest`)

**Database connection**

The DSN is hardcoded in `db/connection.go`:

```go
gorm.Open(mysql.Open("root:Garv@/jwt_project"), &gorm.Config{})
```

This connects as user `root` with password `Garv` to the `jwt_project` database on the
default local socket. Create the database (`CREATE DATABASE jwt_project;`) and adjust
the DSN to match your credentials before running. The `users` and `todos` tables are
created automatically via `AutoMigrate`.

**Run**

```bash
make run          # runs `templ generate` then `go run cmd/main.go`
```

or manually:

```bash
templ generate
go run cmd/main.go
```

The server listens on `http://localhost:4000`. Open `/register` to create an account.

> Note: the JWT signing secret is the literal string `"secret"` in `handler/handler.go`
> and the DB credentials are in source. This is a learning/demo project; both should be
> moved to configuration/environment before any real deployment.

## Usage

Routes registered in `routes/routes.go`:

| Method | Path | Handler | Purpose |
|---|---|---|---|
| `GET` | `/register` | `GetRegister` | Registration page |
| `POST` | `/register` | `PostRegister` | Create user, redirect to `/login` |
| `GET` | `/login` | `GetLogin` | Login page |
| `POST` | `/login` | `PostLogin` | Verify credentials, set JWT cookie, redirect to `/` |
| `POST` | `/logout` | `Logout` | Clear cookie, redirect to `/login` |
| `GET` | `/` | `Home` | Authenticated todo list for the current user |
| `POST` | `/add` | `AddTodo` | Create a todo, return the rendered fragment |
| `POST` | `/toggle/:id` | `IsComplete` | Flip completion, return the fragment |
| `POST` | `/edit/:id` | `EditToDo` | Update content, return the fragment |
| `DELETE` | `/delete/:id` | `DeleteTodo` | Delete a todo |

The add/toggle/edit/delete endpoints are called by HTMX attributes on the rendered
elements and return HTML fragments rather than JSON.

## Project structure

```
GO-TODO-LIST/
├── cmd/
│   └── main.go              # Entry point: DB connect, Echo setup, middleware, :4000
├── db/
│   └── connection.go        # GORM MySQL connection + AutoMigrate
├── handler/
│   └── handler.go           # All HTTP handlers, JWT auth, bcrypt, templ rendering
├── models/
│   ├── todo.go              # Todo model (Id, Content, IsCompleted, UserId)
│   └── user.go              # User model (Id, Name, Email unique, Password)
├── routes/
│   └── routes.go            # Route table
├── view/
│   ├── layout/
│   │   └── layout.templ     # HTML shell (Tailwind, HTMX, Alpine.js via CDN)
│   ├── pages/
│   │   ├── Home.templ       # Todo list page (add form, logout)
│   │   ├── Login.templ      # Login page
│   │   └── Register.templ   # Registration page
│   └── components/
│       └── TodoList.templ   # TodoList container + Todo item (HTMX/Alpine wiring)
├── Makefile                 # `make run` -> templ generate + go run
├── go.mod
└── go.sum
```

> `*_templ.go` files alongside the `.templ` sources are generated by `templ generate`
> and are not edited by hand.
