# SwipeFiles

Tinder-style file cleanup: swipe through files with previews and remove or keep them in one gesture. A Go backend handles file listing and conversion, while a Next.js frontend renders the cards and swipe interactions.

## Features
- Fast swipe UI (mouse, touch, arrow keys)
- Live previews in‑app:
  - Images: png, jpg, jpeg, gif, webp, bmp, svg
  - PDF: embedded viewer
  - Video: mp4, webm, ogv/ogg, mov, m4v
  - Audio: mp3, wav, ogg/oga, m4a, aac
  - Text/logs: txt, md, json, csv/tsv, etc. (first ~64KB)
  - Office: docx, dotx, xlsx, pptx, … converted to PDF server‑side (LibreOffice)
- Go backend with cross‑platform Trash support (macOS, Linux, Windows)
- Simple REST API

## Requirements
- Go 1.22+
- Node.js 18+
- npm or yarn
 - LibreOffice/soffice on the server (only for Office previews)

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
2. Frontend runs at `http://localhost:3000`.

### Configuration
Create `frontend/.env.local`:
```
NEXT_PUBLIC_BACKEND_URL=http://localhost:8787
```

### Usage
- Open the frontend in a browser.
- Files from the chosen directory (default `~/Downloads`) appear as cards.
- Previews open inline for images/PDF/video/audio/text and Office (PDF convert).
- Use ←/→ or drag to Trash/Keep. Trashed files go to the OS Trash/Recycle Bin.

## Project Structure
```
backend/
  main.go                # entrypoint: starts HTTP server
  internal/
    server/              # router construction
    http/
      handlers/          # /api/files, /api/trash, /api/open, /api/convert
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
