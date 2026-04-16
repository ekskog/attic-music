# attic-music

A responsive Vue 3 music player webapp that connects to a Subsonic-compatible server (Gonic) for music library browsing and playback, with optional Last.fm scrobbling.

## Tech Stack

- **Vue 3** — Composition API with `<script setup>` SFCs throughout
- **Pinia** — state management (player queue/state, user config)
- **Vue Router v5** — SPA routing with auth guards
- **Tailwind CSS v4** — via `@tailwindcss/vite`; accent color `--accent: #B85C38`
- **Vite** — dev server and build tool
- **No test framework** — none configured

## Commands

```bash
npm install       # install dependencies
npm run dev       # start dev server (http://localhost:5173)
npm run build     # production build → dist/
npm run preview   # serve production build locally
```

## Project Structure

```
src/
  api/
    subsonic.js   # Subsonic REST API client
    lastfm.js     # Last.fm scrobbling integration
  stores/
    config.js     # user config (server URL, credentials, lastfm)
    player.js     # playback state, queue, shuffle/repeat
  router/
    index.js      # route definitions + auth guards
  views/
    Login.vue     # auth form
    Folders.vue   # folder-based library browsing
    Artists.vue   # artist/album/track browsing
  components/
    Player.vue       # desktop footer player + queue
    MiniPlayer.vue   # mobile mini player
    FullPlayer.vue   # mobile full-screen player
    BottomNav.vue    # mobile bottom nav
    SideBar.vue      # desktop sidebar (nav + last.fm scrobbles)
    FolderNode.vue   # expandable folder tree node
    TrackItem.vue    # track list item
    ArtistCard.vue   # artist grid card with avatar
  App.vue           # root layout + auth check
  main.js           # entry point
  style.css         # Tailwind import + root styles
```

## Architecture Notes

### Auth & Config
- On load, the app fetches `/config.json` for server-provisioned settings (Last.fm API key, optional pre-configured server)
- User credentials are persisted to `localStorage` under key `attic_cfg`
- Route guards block all routes until `config.loggedIn` is true
- In Kubernetes, an init container writes `/config.json` from the `lastfm-secret` K8s secret

### Subsonic API
- All music data comes from `src/api/subsonic.js`, which wraps the Subsonic REST API
- Auth uses hex-encoded password per Subsonic spec; client ID is `atticweb`, API version `1.16.1`
- Dev server proxies `/rest` → `https://gonic.ekskog.net` (see `vite.config.js`)
- Production traffic hits the server the user logs into directly from the browser

### Artist Images
- All artist folders on the server have a `cover.jpg` (pre-fetched externally)
- Gonic does not populate `artist.coverArt` in `getArtists` — the field is always empty regardless of indexing
- `getCoverArt?id=ar-xxx` works for some artists when Gonic has indexed the folder cover, but not reliably for all
- `ArtistCard` tries `getCoverArt?id=artist.id` first, then falls back to Deezer `search/artist` on failure
- Deezer is called lazily in `onImgError` — only for visible artists whose local cover failed, naturally rate-limited by demand
- Deezer returns a placeholder URL containing `/artist//` when no image exists; these are filtered out and fall back to the letter placeholder

### Player
- HTML5 `<audio>` element handles actual playback; Pinia store (`player.js`) manages reactive state
- Stream URLs are authenticated Subsonic `/stream` endpoints
- Queue, shuffle, repeat, seek, and progress are all store-managed

### Responsive Layout
- Desktop: sidebar + main content + footer player
- Mobile: full-width content + mini player + full-screen player modal + bottom nav
- Breakpoint is Tailwind's `md:` (768px)

### Last.fm
- Optional; polls `https://ws.audioscrobbler.com/2.0/` every 30 seconds to show recent scrobbles in the sidebar
- Requires `LASTFM_API_KEY` (injected at runtime, not build time)

## Deployment

### Docker
```bash
docker build -t attic-music .
docker run -p 80:8080 attic-music
```
Multi-stage build: Node 22-Alpine compiles the app, Nginx Alpine serves `dist/`.

### Kubernetes
Manifests live in `k8s/`. The deployment runs in namespace `webapps` and pulls from `ghcr.io/ekskog/attic-music:latest`. An init container generates `/config.json` from the `lastfm-secret` K8s secret before the Nginx container starts.

## Conventions

- All components use `<script setup>` — no Options API
- Composables and stores are imported directly; no global registration
- Tailwind utility classes only — no custom CSS beyond root variables in `style.css`
- Keep API calls inside `src/api/`; views and components call store actions, not API methods directly
