# SwipeFiles

SwipeFiles is a Tinder-style file manager. The Go backend lists files and moves unwanted ones to the system Trash/Recycle Bin. The Next.js (React + TypeScript) frontend displays files and lets users quickly keep or delete them with a simple swipe or click.

## Features
- Go backend with cross-platform Trash support (macOS, Linux, Windows).
- REST API for listing files and sending them to Trash.
- Next.js frontend with React + TypeScript.
- Simple interface to keep or delete files.

## Requirements
- Go 1.22+
- Node.js 18+
- npm or yarn

## Getting Started

### Backend
1. Open a terminal:
   ```bash
   cd backend
   go mod tidy
   # Start HTTP API
   go run ./cmd/server
   ```
2. Backend listens at `http://localhost:8787`.

### Frontend
1. Open another terminal:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```
2. Frontend will run at `http://localhost:3000`.

### Configuration
Create `frontend/.env.local`:
```
NEXT_PUBLIC_BACKEND_URL=http://localhost:8787
```

### Usage
- Open the frontend in browser.
- Files from the chosen directory (default `~/Downloads`) will appear.
- Use ←/→ arrow keys or drag left/right to trash/keep a file.
- Files you keep remain untouched; trashed files go to the OS Trash/Recycle Bin.

## Project Structure
```
backend/
  cmd/server/            # entrypoint: starts HTTP server
  internal/
    server/              # router construction
    http/
      handlers/          # /api/files, /api/trash
      middleware/        # CORS
      respond/           # JSON helpers
    dto/                 # transport types (JSON)
    util/                # path resolution, etc.
  platform/              # OS-specific trash implementations

frontend/
  src/
    app/                 # Next.js App Router pages
    components/          # Card, FilePreview, ActionsBar
    lib/                 # API client
    utils/               # formatting helpers
```

## License
MIT
