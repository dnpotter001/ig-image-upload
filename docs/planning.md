# Instagram Drip-Feed Poster — Planning

A personal project to drip-feed photos from a Google Drive inbox to an Instagram photography page at prime times, via the Instagram Graph API.

## Goal

Upload photos from phone or laptop into a Google Drive folder. A hosted service picks them up, does light processing, queues them, and publishes to Instagram at defined prime-time slots. Low volume, cheap, custom-built as a learning project.

## Architecture



## Stack Decisions

- **Language:** Go — new to me, well-suited to small network services, single static binary deploys
- **Hosting:** Hetzner Cloud VM (~€4/month) — maximum learning value (Linux, systemd, Caddy, ufw)
- **Provisioning:** Terraform (already familiar)
- **CI/CD:** GitHub Actions — build Go binary, scp to VM, restart systemd service
- **Reverse proxy / TLS:** Caddy — auto-provisions Let's Encrypt certs
- **Queue:** files-as-queue pattern using atomic renames between directories
- **Metadata/state:** SQLite on local disk
- **Image processing:** Go `image` stdlib + `disintegration/imaging` or similar for resize
- **Inbox:** Google Drive folder, push notifications via `files.watch`
- **Repo visibility:** public

## Deferred / Out of Scope for v1

- **Automatic RAW editing + LUT application.** Creative editing (LUTs, grading) stays on laptop with Lightroom/Capture One/darktable. Processed JPEGs land in Drive already-edited. Possible follow-up project.
- Web dashboard for queue management / reordering
- Video / reels support (would need FFmpeg, larger VM)
- Multi-account support

## Secrets Strategy

Three categories, three homes:

### Deploy-time (GitHub Actions Secrets)
- `HETZNER_TOKEN`
- `SSH_DEPLOY_KEY`
- `VM_HOST`

### Runtime (on VM at `/etc/igposter/env`, chmod 600)
Loaded via systemd `EnvironmentFile=`:
- `IG_ACCESS_TOKEN` (initial)
- `IG_APP_SECRET`
- `GOOGLE_CLIENT_ID`
- `GOOGLE_CLIENT_SECRET`
- `GOOGLE_REFRESH_TOKEN`
- `WEBHOOK_VERIFY_TOKEN`

### Rotating tokens (on VM at `/var/lib/igposter/tokens.json`)
- Current IG long-lived token (refreshed every ~50 days, expires at 60)
- Current Google access token

### Master copies
- Password manager (1Password / Bitwarden)

### Hygiene rules
- `.gitignore`: `.env`, `*.env`, `tokens.json`, `*.pem`, `*.key`, `terraform.tfstate*`, `secrets/`
- Commit `.env.example` with keys but no values
- Pre-commit `gitleaks` hook
- GitHub secret scanning enabled (automatic for public repos)
- Settings → Actions → "Require approval for fork PRs"
- If a secret ever leaks: rotate immediately, don't rely on history rewrite

## Prerequisites to Set Up

- Instagram Business or Creator account (switch in IG app)
- Facebook Page linked to the IG account
- Meta Developer app with Instagram Graph API product
- Permissions: `instagram_basic`, `instagram_content_publish`, `pages_show_list`, `pages_read_engagement`
- Long-lived IG access token (~60 day expiry, must auto-refresh)
- Google Cloud project with Drive API enabled
- OAuth credentials for Drive access
- Hetzner Cloud account + API token

## Build Order

### 1. Instagram publisher (validate riskiest piece first)
- Meta app setup, IG Business account, Facebook Page link
- Local Go script that publishes a hardcoded JPEG via Graph API
- Implement long-lived token refresh flow
- Prove publish end-to-end

### 2. Local queue + scheduler
- Directory structure: `inbox/`, `posted/`, `failed/`, `metadata/`
- Atomic file moves for state transitions
- SQLite for post history and tokens
- Ticker-based scheduler checking prime-time windows (e.g., 8am, 12pm, 7pm NZST)
- Run locally, drop files into `inbox/`, verify posting at correct times

### 3. Drive integration
- Google OAuth flow (one-time, refresh token stored)
- Webhook receiver with Caddy-terminated HTTPS
- `files.watch` setup and daily renewal (channels expire in 7 days)
- Download new files, process (resize 1080px, sRGB conversion, EXIF strip, JPEG q85), move to queue
- Archive originals to a Drive "archive" folder

### 4. Deploy
- Terraform for Hetzner VM provisioning
- cloud-init for first-boot: user creation, package install, ufw, Caddy
- systemd unit for the service
- GitHub Actions workflow: build → scp → restart
- Manual first deploy to learn the ropes, then codify in Terraform/Actions

## Repo Layout (planned)

igposter/
├── cmd/igposter/main.go
├── internal/
│   ├── drive/
│   ├── instagram/
│   ├── queue/
│   ├── processor/
│   └── scheduler/
├── deploy/
│   ├── terraform/
│   ├── cloud-init.yaml
│   ├── Caddyfile
│   └── igposter.service
├── .github/workflows/
│   ├── deploy.yml
│   └── terraform.yml
├── .env.example
├── .gitignore
├── PLANNING.md
└── README.md

## Prime-Time Logic

- Timer/ticker runs every 15 minutes
- Check: (a) inside a posting window, (b) no post yet in this window, (c) queue has items
- If all three → publish oldest in queue → mark in SQLite
- Windows configurable; start with 3 slots/day in NZST

## Known Gotchas

- IG long-lived tokens expire every 60 days → must auto-refresh
- Drive `files.watch` channels expire max 7 days → daily renewal job
- iPhone Display P3 → must convert to sRGB for Instagram or colors wash out
- Rate limit: 100 IG posts per 24h per account (not a concern at low volume)
- Media for IG must be at a publicly fetchable URL with correct `Content-Type` header
- Terraform state contains sensitive data → don't commit `.tfstate`, use remote backend or keep local + gitignored

## Future Ideas

- Automatic RAW + LUT editing pipeline (separate project)
- Web dashboard for queue preview, reorder, skip
- Analytics ingestion (insights per post)
- Healthcheck endpoint + uptime monitoring
- Staging VM on a separate branch
- Rollback-on-failure in deploy workflow
- Caption generation from filename or sidecar `.txt`
- Multi-account support



