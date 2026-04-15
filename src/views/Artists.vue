<template>
  <div class="flex flex-col h-full overflow-hidden">

    <!-- ARTISTS GRID -->
    <template v-if="view === 'grid'">
      <div class="px-8 py-7 border-b border-stone-200 bg-white flex-shrink-0">
        <h1 class="font-serif text-4xl font-semibold">Artists</h1>
      </div>
      <div class="flex-1 overflow-y-auto px-6 py-4 pb-40 md:pb-24">
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <div v-else>
          <div v-for="group in artistIndex" :key="group.name" class="mb-6">
            <div
              class="flex items-center justify-between font-serif text-xl font-semibold text-amber-700 border-b border-stone-200 pb-1 mb-3 cursor-pointer select-none"
              @click="toggleGroup(group.name)"
            >
              <span>{{ group.name }}</span>
              <span class="text-base transition-transform duration-200 inline-block" :class="{ 'rotate-90': expandedGroups[group.name] }">›</span>
            </div>
            <div v-if="expandedGroups[group.name]" class="grid gap-4" style="grid-template-columns: repeat(auto-fill, minmax(120px, 1fr))">
              <ArtistCard
                v-for="artist in group.artist"
                :key="artist.id"
                :artist="artist"
                @click="openArtist(artist)"
              />
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- ARTIST DETAIL -->
    <template v-else-if="view === 'artist'">
      <div class="px-8 py-7 border-b border-stone-200 bg-white flex-shrink-0">
        <div class="flex items-center gap-1.5 text-xs text-stone-400 mb-3">
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="view = 'grid'">Artists</span>
          <span class="opacity-40">›</span>
          <span>{{ currentArtist.name }}</span>
        </div>
        <div class="flex items-center gap-4">
          <div class="w-14 h-14 rounded-xl overflow-hidden flex-shrink-0 bg-stone-100">
            <img
              v-if="artistDetailImageUrl"
              :src="artistDetailImageUrl"
              :alt="currentArtist.name"
              class="w-full h-full object-cover"
              @load="onArtistDetailImgLoad"
              @error="onArtistDetailImgError"
            />
            <div v-else class="w-full h-full flex items-center justify-center font-serif text-2xl font-semibold text-stone-300 select-none">
              {{ currentArtist.name?.[0]?.toUpperCase() }}
            </div>
          </div>
          <h1 class="font-serif text-4xl font-semibold">{{ currentArtist.name }}</h1>
        </div>
      </div>
      <div class="flex-1 overflow-y-auto px-6 py-4 pb-40 md:pb-24">
        <div v-if="loadingArtist" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <div v-else class="grid gap-3" style="grid-template-columns: repeat(auto-fill, minmax(88px, 1fr))">
          <div
            v-for="album in currentArtistAlbums" :key="album.id"
            class="cursor-pointer group"
            @click="openAlbum(album)"
          >
            <div class="aspect-square bg-amber-50 mb-1 overflow-hidden relative rounded">
              <img v-if="album.coverArt" :src="coverUrl(album.coverArt)" :alt="album.name" class="w-full h-full object-cover" @error="onImgError" />
              <div v-else class="w-full h-full flex items-center justify-center text-2xl">💿</div>
              <div class="absolute inset-0 bg-black/25 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity rounded">
                <button class="w-7 h-7 rounded-full bg-white flex items-center justify-center text-xs pl-0.5" @click.stop="playAlbum(album)">▶</button>
              </div>
            </div>
            <div class="text-xs font-medium truncate leading-tight">{{ album.name }}</div>
            <div class="text-xs text-stone-400 mt-0.5">{{ album.year || '' }}</div>
          </div>
        </div>
      </div>
    </template>

    <!-- ALBUM DETAIL -->
    <template v-else-if="view === 'album'">
      <div class="px-8 py-7 border-b border-stone-200 bg-white flex-shrink-0">
        <div class="flex items-center gap-1.5 text-xs text-stone-400 mb-2 flex-wrap">
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="view = 'grid'">Artists</span>
          <span class="opacity-40">›</span>
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="view = 'artist'">{{ currentArtist.name }}</span>
          <span class="opacity-40">›</span>
          <span>{{ currentAlbum.name }}</span>
        </div>
        <h1 class="font-serif text-4xl font-semibold">{{ currentAlbum.name }}</h1>
      </div>
      <div class="flex-1 overflow-y-auto px-8 py-6 pb-40 md:pb-24">
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <div v-else>
          <div class="flex gap-6 mb-8 items-end">
            <div class="w-40 h-40 flex-shrink-0 bg-amber-50 overflow-hidden">
              <img v-if="currentAlbum.coverArt" :src="coverUrl(currentAlbum.coverArt)" :alt="currentAlbum.name" class="w-full h-full object-cover" @error="onImgError" />
              <div v-else class="w-full h-full flex items-center justify-center text-5xl">💿</div>
            </div>
            <div>
              <div class="text-xs uppercase tracking-widest text-stone-400 mb-2">
                Album{{ currentAlbum.year ? ' · ' + currentAlbum.year : '' }}
              </div>
              <div class="font-serif text-3xl font-semibold mb-1.5">{{ currentAlbum.name }}</div>
              <div class="text-sm text-stone-400 mb-4">{{ currentAlbum.artist }}</div>
              <div class="flex gap-2.5">
                <button class="bg-stone-900 text-white text-sm font-medium px-5 py-2 hover:bg-amber-700 transition-colors" @click="playAlbumTracks">▶ Play</button>
                <button class="border border-stone-200 text-sm px-4 py-2 hover:border-amber-700 hover:text-amber-700 transition-colors" @click="queueAlbumTracks">+ Queue</button>
              </div>
            </div>
          </div>

          <div class="grid gap-px text-xs uppercase tracking-widest text-stone-400 px-3 pb-2 border-b border-stone-200 mb-1"
            style="grid-template-columns: 36px 1fr 160px 56px 32px">
            <span class="text-center">#</span><span>Title</span><span>Artist</span><span>Time</span><span></span>
          </div>
          <TrackItem
            v-for="(track, i) in albumTracks" :key="track.id"
            :track="track"
            :index="i"
            @play="player.playTrack(track, albumTracks, i)"
            @queue="player.addToQueue(track)"
          />
        </div>
      </div>
    </template>

  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { usePlayerStore } from '../stores/player'
import { getArtists, getArtist, getAlbum, coverUrl } from '../api/subsonic'
import { getArtistImage } from '../api/deezer'
import TrackItem  from '../components/TrackItem.vue'
import ArtistCard from '../components/ArtistCard.vue'

const player = usePlayerStore()

const view    = ref('grid')
const loading = ref(false)

const artistIndex    = ref([])
const expandedGroups = reactive({})

const currentArtist           = ref(null)
const currentArtistAlbums     = ref([])
const loadingArtist           = ref(false)
const artistDetailImageUrl    = ref(null)

const currentAlbum = ref(null)
const albumTracks  = ref([])

async function loadArtists() {
  loading.value = true
  try {
    artistIndex.value = await getArtists()
    for (const group of artistIndex.value) expandedGroups[group.name] = true
  } finally { loading.value = false }
}

function toggleGroup(name) {
  expandedGroups[name] = !expandedGroups[name]
}

async function openArtist(artist) {
  currentArtist.value = { ...artist }
  artistDetailImageUrl.value = null
  view.value = 'artist'
  loadingArtist.value = true
  try {
    const data = await getArtist(artist.id)
    currentArtistAlbums.value = data.albums
    // Gonic doesn't expose artist-level cover art; use the first album's coverArt as the avatar
    const coverArtId = data.info?.coverArt || data.albums?.[0]?.coverArt || null
    if (coverArtId) {
      console.log(`[artist image] "${artist.name}" — using coverArt id: ${coverArtId}`)
      artistDetailImageUrl.value = coverUrl(coverArtId, 112)
    } else {
      console.log(`[artist image] "${artist.name}" — no coverArt in response, fetching from Deezer…`)
      const deezerUrl = await getArtistImage(artist.name)
      if (deezerUrl) {
        console.log(`[artist image] "${artist.name}" — using Deezer image`)
      } else {
        console.log(`[artist image] "${artist.name}" — no image found anywhere, showing initial`)
      }
      artistDetailImageUrl.value = deezerUrl
    }
  } finally { loadingArtist.value = false }
}

function onArtistDetailImgLoad() {
  console.log(`[artist image] "${currentArtist.value?.name}" — loaded from local library`)
}

async function onArtistDetailImgError() {
  const failedUrl = artistDetailImageUrl.value
  console.log(`[artist image] "${currentArtist.value?.name}" — local URL failed (${failedUrl}), fetching from Deezer…`)
  artistDetailImageUrl.value = null
  const url = await getArtistImage(currentArtist.value?.name)
  if (url) {
    console.log(`[artist image] "${currentArtist.value?.name}" — using Deezer image`)
  } else {
    console.log(`[artist image] "${currentArtist.value?.name}" — no image found anywhere, showing initial`)
  }
  artistDetailImageUrl.value = url
}

async function openAlbum(album) {
  loading.value = true
  try {
    const data = await getAlbum(album.id)
    currentAlbum.value = { ...album, ...data.info }
    albumTracks.value  = data.tracks
    view.value = 'album'
  } finally { loading.value = false }
}

async function playAlbum(album) {
  await openAlbum(album)
  playAlbumTracks()
}

function playAlbumTracks() {
  if (!albumTracks.value.length) return
  player.playTrack(albumTracks.value[0], albumTracks.value, 0)
}

function queueAlbumTracks() {
  player.queue.push(...albumTracks.value)
}

function onImgError(e) {
  try { if (e?.target) e.target.style.display = 'none' } catch (_) {}
}

onMounted(loadArtists)
</script>
