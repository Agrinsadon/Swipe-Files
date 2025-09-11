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
   go run .
   ```
2. Backend will run at `http://localhost:8787`.

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
- Files from the chosen directory (default `~/Desktop`) will appear.
- Click the trash button (or later swipe) to move files to the system Trash/Recycle Bin.
- Files you keep remain untouched.

## License
MIT
