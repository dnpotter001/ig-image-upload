# Instagram Drip-Feed Poster — Planning

A personal project to drip-feed photos from a Google Drive inbox to an Instagram photography page at prime times, via the Instagram Graph API.

## Goal

Upload photos from phone or laptop into a Google Drive folder. A hosted service picks them up, does light processing, queues them, and publishes to Instagram at defined prime-time slots. Low volume, cheap, custom-built as a learning project.

## Architecture

```
[phone/laptop]
     │ upload JPEG or folder of JPEGs
     ▼
[Google Drive inbox folder]
     │
     │ polled every ~5 min
     ▼
┌──────────────────── Hetzner VM ────────────────────┐
│  Go binary (igposter)                              │
│  ├─ poller     → diff Drive vs SQLite              │
│  ├─ processor  → resize, sRGB, EXIF strip, q85     │
│  ├─ queue      → inbox/ → posted/ | failed/        │
│  ├─ scheduler  → publish at prime-time windows     │
│  └─ http       → serves /media/<id>.jpg over HTTPS │
│                  (autocert / Let's Encrypt)        │
│                                                    │
│  /var/lib/igposter  (Hetzner Volume)               │
│  ├─ inbox/  posted/  failed/                       │
│  ├─ state.db  (SQLite)                             │
│  └─ tokens.json                                    │
└────────────────────────────────────────────────────┘
     │ Graph API publish (single image or carousel)
     ▼
[Instagram]
```

## Stack Decisions

- **Language:** Go — new to me, well-suited to small network services, single static binary deploys
- **Hosting:** Hetzner Cloud VM (~€4/month) + small Hetzner Volume for state — maximum learning value (Linux, systemd, ufw)
- **Provisioning:** Hetzner web console for VM + Volume, cloud-init for first-boot config, `deploy/bootstrap.sh` for anything imperative. No Terraform — overkill for one VM and not what this project is teaching.
- **CI/CD:** GitHub Actions — build Go binary, scp to VM, restart systemd service
- **TLS:** Go `golang.org/x/crypto/acme/autocert` — Let's Encrypt directly from the binary, no reverse proxy
- **Queue:** files-as-queue pattern using atomic renames between directories (all on the mounted volume, same filesystem)
- **Metadata/state:** SQLite on the mounted volume
- **Image processing:** Go `image` stdlib + `disintegration/imaging` or similar for resize
- **Inbox:** Google Drive folder, **polled** every ~5 min (no webhooks)
- **Image hosting for IG:** VM serves processed JPEGs at `https://<domain>/media/<id>.jpg` so IG's ingestor can fetch them; files deleted after successful publish
- **Repo visibility:** public

## Deferred / Out of Scope for v1

- **Automatic RAW editing + LUT application.** Creative editing (LUTs, grading) stays on laptop with Lightroom/Capture One/darktable. Processed JPEGs land in Drive already-edited. Possible follow-up project.
- Web dashboard for queue management / reordering
- Video / reels support (would need FFmpeg, larger VM)
- Multi-account support

## Secrets Strategy

Three categories, three homes:

### Deploy-time (GitHub Actions Secrets)
- `SSH_DEPLOY_KEY`
- `VM_HOST`

### Runtime (on VM at `/etc/igposter/env`, chmod 600)
Loaded via systemd `EnvironmentFile=`:
- `IG_ACCESS_TOKEN` (initial)
- `IG_APP_SECRET`
- `GOOGLE_CLIENT_ID`
- `GOOGLE_CLIENT_SECRET`
- `GOOGLE_REFRESH_TOKEN`
- `PUBLIC_DOMAIN` (used by autocert and for `/media/` URLs)

### Rotating tokens (on VM at `/var/lib/igposter/tokens.json`)
- Current IG long-lived token (refreshed every ~50 days, expires at 60)
- Current Google access token

### Master copies
- Password manager (1Password / Bitwarden)

### Hygiene rules
- `.gitignore`: `.env`, `*.env`, `tokens.json`, `*.pem`, `*.key`, `secrets/`
- Commit `.env.example` with keys but no values
- Pre-commit `gitleaks` hook
- GitHub secret scanning enabled (automatic for public repos)
- Settings → Actions → "Require approval for fork PRs"
- If a secret ever leaks: rotate immediately, don't rely on history rewrite

## Prerequisites to Set Up

- Instagram Business or Creator account (switch in IG app)
- Meta Developer app with Instagram content publishing capability
  - **Decide which API path:** "Instagram API with Instagram Login" (newer, no FB Page required) vs the older Page-linked Graph API. Permission set and token refresh endpoints differ — pick before doing app setup.
  - App stays in **Development Mode** indefinitely (no App Review needed) since the only IG account it acts on is the owner's, added as a tester/admin.
- Long-lived IG access token (~60 day expiry, must auto-refresh)
- Google Cloud project with Drive API enabled
- OAuth credentials for Drive access
- A domain pointed at the Hetzner VM (required for HTTPS via Let's Encrypt; IG and your browser both need it)
- Hetzner Cloud account + API token

## Build Order

### 1. Instagram publisher (validate riskiest piece first)
- Decide API path (Instagram Login vs Page-linked Graph API), Meta app setup
- Local Go script that publishes a hardcoded JPEG via Graph API
- Extend to carousel publish (2–10 images, 3-step container flow)
- Implement long-lived token refresh flow
- Prove publish end-to-end

### 2. Local queue + scheduler
- Directory structure on mounted volume: `inbox/`, `posted/`, `failed/`
- A **post** is the queue unit. Loose file in `inbox/` = single-image post. Subfolder in `inbox/` = carousel post (up to 10 children).
- Atomic file moves for state transitions (all dirs on same filesystem)
- SQLite for post history and tokens
- Ticker-based scheduler checking prime-time windows (e.g., 8am, 12pm, 7pm `Pacific/Auckland` — use `time.LoadLocation` so DST is handled)
- Run locally, drop files/folders into `inbox/`, verify posting at correct times

### 3. Drive integration (polling)
- Google OAuth flow (one-time, refresh token stored)
- Poller runs every ~5 min: list Drive inbox folder, diff against SQLite (`seen_files` table tracking Drive file IDs + modifiedTime), enqueue new items
- **Folder = carousel.** Folder children are downloaded together as one post.
- **Stability window:** only enqueue a Drive folder whose contents have been unchanged for ≥5 min (avoids grabbing mid-upload). Single files can enqueue immediately.
- **Ordering within a carousel:** if all filenames start with a number, sort numerically ascending; otherwise use Drive `createdTime` (arrival order).
- Download new files, process (resize 1080px, sRGB conversion, EXIF strip, JPEG q85), move to queue
- Archive originals to a Drive "archive" folder

### 4. Deploy
- Create VM + Volume in the Hetzner console, point DNS at the IP
- cloud-init for first-boot: user creation, package install, ufw, volume mount at `/var/lib/igposter`
- `deploy/bootstrap.sh` for anything cloud-init can't easily express (idempotent, re-runnable)
- systemd unit for the service, with `EnvironmentFile=/etc/igposter/env`
- Open ports 80 + 443 in ufw (autocert needs both)
- GitHub Actions workflow: build → scp → restart
- Manual first deploy to learn the ropes, then codify in cloud-init + Actions

## Repo Layout (planned)

igposter/
├── cmd/igposter/main.go
├── internal/
│   ├── drive/        # OAuth, polling, download
│   ├── instagram/    # Graph API, token refresh, carousel publish
│   ├── queue/        # filesystem queue, atomic moves
│   ├── processor/    # resize, sRGB, EXIF strip
│   ├── scheduler/    # prime-time windows
│   └── httpd/        # autocert TLS, /media/ static serving
├── deploy/
│   ├── cloud-init.yaml
│   ├── bootstrap.sh
│   └── igposter.service
├── .github/workflows/
│   └── deploy.yml
├── .env.example
├── .gitignore
└── README.md
docs/
└── planning.md

## Prime-Time Logic

- Timer/ticker runs every 15 minutes
- Check: (a) inside a posting window, (b) no post yet in this window, (c) queue has items
- If all three → publish oldest post in queue (single image or carousel) → mark in SQLite
- Windows configurable; start with 3 slots/day in `Pacific/Auckland`

## Carousel Logic

- IG Graph API carousel = 3-step flow: create per-child `IMAGE` containers (`is_carousel_item=true`), create a `CAROUSEL` container with `children=[...]`, then publish.
- Limit: 2–10 children, each ≤8MB JPEG, aspect ratio within IG's allowed range.
- If any child fails validation, whole post → `failed/`. Don't silently drop children.
- Children passed in the order described in **Build Order step 3** (numeric prefix → numeric ascending; otherwise Drive `createdTime`).

## Known Gotchas

- IG long-lived tokens expire every 60 days → must auto-refresh (refresh endpoint requires the token be ≥24h old)
- iPhone Display P3 → must convert to sRGB for Instagram or colors wash out
- Rate limit: 100 IG posts per 24h per account (not a concern at low volume)
- Media for IG must be at a publicly fetchable HTTPS URL with correct `Content-Type` header — served from our VM at `/media/<id>.jpg`
- Drive folder polled mid-upload could yield an incomplete carousel → stability window required
- Drive list order is not user-meaningful → rely on numeric prefix or `createdTime`, not API order
- VM disk is ephemeral if rebuilt → all mutable state must live on the mounted Hetzner Volume
- `autocert` needs ports 80 + 443 reachable from the public internet for the ACME HTTP-01 challenge

## Future Ideas

- Automatic RAW + LUT editing pipeline (separate project)
- Web dashboard for queue preview, reorder, skip
- Analytics ingestion (insights per post)
- Healthcheck endpoint + uptime monitoring
- Staging VM on a separate branch
- Rollback-on-failure in deploy workflow
- Caption generation from filename or sidecar `.txt`
- Multi-account support



