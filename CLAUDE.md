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
artist-images/
  main.go         # Go sidecar — serves artist cover.jpg from NFS
  Dockerfile      # multi-stage build → minimal Alpine image
src/
  api/
    subsonic.js   # Subsonic REST API client
    lastfm.js     # Last.fm scrobbling integration
  stores/
    config.js     # user config (server URL, credentials, lastfm)
    player.js     # playback state, queue, shuffle/repeat
    playlist.js   # playlist list cache + addTrack action
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
- In CI, the GHA workflow writes `public/config.json` (including `lastfmKey` from the `LASTFM_API_KEY` GitHub Actions secret) before the Docker image is built, so the key is baked into the image as a static file served by Nginx

### Subsonic API
- All music data comes from `src/api/subsonic.js`, which wraps the Subsonic REST API
- Auth uses hex-encoded password per Subsonic spec; client ID is `atticweb`, API version `1.16.1`
- Dev server proxies `/rest` → `https://gonic.ekskog.net` (see `vite.config.js`)
- Production traffic hits the server the user logs into directly from the browser

### NFS File Structure
- Subsonic indexes media from an NFS Share. The structure of the share is:
./mp3/<first letter in artist name>/YYYY-album_name_with_underscore, as below:
/var/lib/media/music/mp3 $ tree | more
.
├── 1
│   ├── !!!
│   │   ├── 2013-thr_er
│   │   │   ├── 01-01-even_when_the_water's_cold.mp3
│   │   │   ├── 01-02-get_that_rhythm_right.mp3
│   │   │   └── cover.jpg
│   │   └── cover.jpg
│   ├── 10,000_maniacs
│   │   ├── 1983-unplugged_on_mtv_preshow
│   │   │   ├── 01-01-how_you've_grown_take_1.mp3
│   │   │   ├── 01-09-how_you've_grown_take_3.mp3
│   │   │   └── cover.jpg
│   │   └── cover.jpg
├── a
│   ├── abba
│   │   ├── 1993-gold
│   │   │   ├── 01-01-dancing_queen.mp3
│   │   │   ├── 01-06-super_trooper.mp3
│   │   │   └── cover.jpg
│   │   └── cover.jpg

### Artist Images
- All artist folders on the NFS share have a `cover.jpg` (pre-fetched externally)
- Gonic does **not** populate `artist.coverArt` in `getArtists` — the field is always empty regardless of indexing
- A dedicated Go sidecar (`artist-images/`) serves images directly from the NFS volume, bypassing the Subsonic API entirely
- Gonic does **not** serve standalone `cover.jpg` files via `getCoverArt` — it only returns embedded ID3 art; the sidecar is the authoritative source for both artist and album cover art
- `ArtistCard` image fallback chain: sidecar (`/artist-images/avatar?name=<artist>`) → Subsonic `getCoverArt?id=<artistId>` → letter placeholder
- Album cover fallback chain (carousel, grid, detail): Subsonic `getCoverArt?id=<albumId>` → sidecar (`/artist-images/album?artist=<artist>&album=<album>`) → 💿 placeholder
- The artist detail view: Subsonic first, then sidecar fallback
- The sidecar builds two maps at startup and rescans every 5 minutes: artist covers (`normalize(artist)` → path) and album covers (`normalize(artist)|normalize(album)` → path)
- Album folders are named `YYYY-album_name_with_underscores`; the sidecar strips the leading `YYYY-` before normalizing so it matches the API album name
- The artist folder name corresponds to the **Album Artist** mp3 tag (e.g. `bob_seger`), not the track Artist tag (e.g. `bob seger & the silver bullet band`); the frontend uses `album.albumArtist || album.artist` when building the sidecar URL
- Directory structure scanned: `<root>/<letter>/<artist_folder>/cover.jpg` (artist) and `<root>/<letter>/<artist_folder>/YYYY-album_folder/cover.jpg` (album)
- Request logging (HIT/MISS with latency) is controlled by the `LOG_REQUESTS` env var, set via ConfigMap `artist-images-config` in namespace `webapps`
- In dev, Vite proxies `/artist-images` → `http://localhost:8081`; in production, Nginx proxies it to the `artist-images` ClusterIP service
- Manifest: `k8s/artist-images.yaml` (ConfigMap + Deployment + Service); image: `ghcr.io/ekskog/artist-images:latest`

### Player
- HTML5 `<audio>` element handles actual playback; Pinia store (`player.js`) manages reactive state
- Stream URLs are authenticated Subsonic `/stream` endpoints
- Queue, shuffle, repeat, seek, and progress are all store-managed

### Playlists
- Playlists are fully server-side (Gonic); `getPlaylists`, `getPlaylist`, `updatePlaylist`, `deletePlaylist` are all in `subsonic.js`
- `playlist.js` store caches the playlist list (fetched once on first use); `addTrack(playlistId, trackId)` calls `updatePlaylist` to append a song
- `TrackItem`'s `+` button opens an inline dropdown listing playlists plus an "Add to queue" option — no view changes required; the dropdown is self-contained in the component

### TrackItem
- `TrackItem.vue` is used in both album detail (`Artists.vue`) and playlist detail (`Playlists.vue`)
- On mobile: artist shown as a subtitle line below the track title
- On desktop (`md:`): artist shown in a dedicated second column; grid is `28px 1fr 1fr 44px 28px` (# / Title / Artist / Time / action)
- Header rows in `Artists.vue` and `Playlists.vue` use the same responsive grid template to stay aligned

### Artist Detail
- Album cards show year and track count (`album.songCount` from the `getArtist` response) as `1967 · 13 tracks`
- Breadcrumb letter (e.g. `t` for The Beatles) is stored and displayed lowercase; do not apply `uppercase` CSS to it

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
