<div align="center">
  <img src="assets/ao-logo.svg" alt="GIT_WIZ" width="160" height="160" />

# GIT_WIZ

An Agentic IDE that supervises parallel AI coding agents in isolated workspaces, with complete control and automatic feedback loops from CI failures, review comments, and merge conflicts.

</div>

---

## What is GIT_WIZ?

GIT_WIZ is a meta-harness agent IDE for running AI coding agents in parallel. It gives terminal-based agents like Claude Code, Codex, Cursor, Kimi Code, opencode, and others a shared workspace where their sessions, terminals, branches, pull requests, and feedback loops can be supervised from one place.

The agents still do the coding. GIT_WIZ provides the harness around them: isolated workspaces, live terminal access, session state, PR awareness, and automatic loops that send CI failures, review comments, and merge conflicts back to the right agent. Instead of manually coordinating a pile of agent terminals, GIT_WIZ turns parallel agent work into a managed workflow.

## Why GIT_WIZ?

AI coding agents become much more useful when they can work in parallel, but parallel work gets messy quickly. Branches overlap, terminals get lost, CI failures need follow-up, review comments need replies, and merge conflicts have to reach the right worker.

GIT_WIZ is built to keep that loop visible and manageable. It helps you:

* **Start multiple agents from the same project without mixing their work**
* **Keep every session in a separate git worktree**
* **See which agents are working, waiting, finished, or blocked**
* **Route CI failures, review comments, and merge conflicts back to the right session**
* **Use different agent CLIs through one common supervisor**

## How it works

At a high level, GIT_WIZ follows a simple loop:

1. Add a project you want agents to work on.
2. Start one or more sessions from the desktop app or CLI.
3. GIT_WIZ creates an isolated git worktree for each session.
4. GIT_WIZ launches the selected coding agent in that session's terminal runtime.
5. The local daemon watches session state, terminal activity, pull requests, CI, and review feedback.
6. The desktop app and CLI show the current state and let you send follow-up instructions to the right session.

The result is a local control layer for agentic coding: agents still do the coding, while GIT_WIZ keeps their workspaces, status, terminals, and feedback loops organized.

## Features

The desktop app is the main control surface: projects on the left, active sessions in the center, and the selected session's terminal, pull request state, review runs, and browser preview in the inspector.

* **Parallel agent sessions**: Start multiple coding agents from the same project without mixing files, branches, terminals, or pull request state.
* **Live terminal control**: Open any session and attach to the worker terminal while keeping session summary, PR state, and follow-up actions in view.
* **Review feedback loop**: Run reviewer agents, inspect review status, and route requested changes back to the right worker session.
* **In-app browser preview**: Preview a session's local app beside the terminal so UI work, browser state, and agent output stay together.

## Quickstart & Development

### Prerequisites

- **Go** (1.22+ recommended)
- **Node.js** (v18+) & **npm**

### Backend (Go Daemon & CLI)

The backend daemon (`ao`) manages parallel worker sessions, terminal muxing, git worktrees, and loopback HTTP services.

```bash
# Navigate to backend directory
cd backend

# Build backend binaries
go build ./...

# Run backend unit tests
go test ./...

# Start the local daemon
go run ./cmd/ao start
```

### Frontend (Electron + React Supervisor)

The desktop app is an Electron supervisor wired to the daemon via typed HTTP API calls.

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Run application in development mode
npm run dev

# Run TypeScript type check
npm run typecheck
```

### Code Generation & Repository Scripts

Run root-level scripts for linting and code generation when modifying database schemas or API specs:

```bash
# Run lint checks across backend and frontend
npm run lint
npm run frontend:typecheck

# Regenerate SQLite queries (sqlc)
npm run sqlc

# Regenerate OpenAPI spec and frontend TypeScript client types
npm run api
```

## Telemetry

GIT_WIZ's Electron renderer sends anonymous usage events to PostHog for reliability and product understanding, and PostHog session recording is enabled with local paths and local URLs redacted before transmission. Set `VITE_AO_POSTHOG_KEY` to an empty string before building to disable transmission.

## License

Apache License 2.0. See [LICENSE](LICENSE).
