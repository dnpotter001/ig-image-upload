# Instagram Drip-Feed Poster — Plan

Drip-feed photos from a Google Drive inbox to an Instagram photography page at prime times, via the Instagram Graph API. A Go binary on a Hetzner VM polls Drive, processes photos, queues them, and publishes on a schedule.

## Setup / Prerequisites

- [ ] Switch Instagram to a Business or Creator account
- [ ] Create a Meta Developer app with content publishing (decide: "Instagram API with Instagram Login" vs Page-linked Graph API — pick before app setup)
- [ ] Add own IG account as tester/admin (app stays in Development Mode, no App Review)
- [ ] Get a long-lived IG access token
- [ ] Create a Google Cloud project with Drive API enabled + OAuth credentials
- [ ] Register a domain and point it at the VM
- [ ] Create a Hetzner Cloud account + API token

## Build Order

### 1. Instagram publisher
- [ ] Local Go script that publishes a hardcoded JPEG via Graph API
- [ ] Extend to carousel publish (2–10 images, 3-step container flow)
- [ ] Prove publish end-to-end
- [ ] Static secrets from env vars; paste a fresh 60-day IG token by hand when it expires (no auto-refresh)

### 2. Local queue + scheduler
- [ ] Directory structure: `inbox/`, `posted/`, `failed/` (loose file = single post, subfolder = carousel)
- [ ] Atomic file moves for state transitions
- [ ] File-based state (last-post timestamp + `posted/` dir) for the prime-time window check
- [ ] Ticker-based scheduler for prime-time windows (`Pacific/Auckland`, DST-aware)
- [ ] Run locally, drop files into `inbox/`, verify posting at correct times

### 3. Image processing (`internal/editor`)
- [x] White-border module via ImageMagick: resize to fit 1020px, pad to a 1080×1080 white square (`AddWhiteBorder`, returns error on magick failure)
- [ ] Wire into the queue pipeline (process on the way from inbox → queue)

### 4. Drive integration (polling) — next step
- [ ] Google OAuth flow (one-time, store refresh token)
- [ ] Poller every ~5 min: list Drive inbox, diff against a local seen-files record, enqueue new items
- [ ] Stability window: only enqueue a folder unchanged for ≥5 min; single files enqueue immediately
- [ ] Carousel ordering: numeric filename prefix ascending, else Drive `createdTime`
- [ ] Download new items, run through image processing, move to queue
- [ ] Archive originals to a Drive "archive" folder

### 5. Deploy
- [ ] Create VM + Volume in Hetzner console, point DNS at the IP
- [ ] cloud-init first-boot: user, packages, ufw, volume mount at `/var/lib/igposter`
- [ ] `deploy/bootstrap.sh` for anything cloud-init can't express (idempotent)
- [ ] systemd unit with `EnvironmentFile=/etc/igposter/env`
- [ ] Open ports 80 + 443 in ufw (autocert needs both)
- [ ] GitHub Actions workflow: build → scp → restart
- [ ] Manual first deploy, then codify in cloud-init + Actions
