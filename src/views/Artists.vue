<template>
  <div class="flex flex-col h-full overflow-hidden">

    <!-- ARTIST LIST -->
    <template v-if="!currentArtist && !currentAlbum">
      <div class="px-8 py-7 border-b border-stone-200 bg-white flex-shrink-0">
        <h1 class="font-serif text-4xl font-semibold">Artists</h1>
      </div>
      <div class="flex-1 overflow-y-auto px-8 py-6 pb-40 md:pb-24">
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <div v-else>
          <div v-for="group in artistIndex" :key="group.name" class="mb-6">
            <div class="font-serif text-xl font-semibold text-amber-700 border-b border-stone-200 pb-1 mb-2">
              {{ group.name }}
            </div>
            <div class="grid gap-5" style="grid-template-columns: repeat(auto-fill, minmax(168px, 1fr))">
              <div
                v-for="artist in group.artist" :key="artist.id"
                class="cursor-pointer group"
                @click="openArtist(artist)"
              >
                <div class="aspect-square bg-amber-50 mb-2.5 overflow-hidden relative">
                  <img
                    v-if="artist.coverArt"
                    :src="coverUrl(artist.coverArt)"
                    :alt="artist.name"
                    class="w-full h-full object-cover"
                    @error="onImgError"
                  />
                  <div v-else class="w-full h-full flex items-center justify-center text-3xl">🎤</div>
                </div>
                <div class="text-sm font-medium truncate text-center">{{ artist.name }}</div>
                <div v-if="artist.albumCount" class="text-xs text-stone-400 mt-0.5 text-center">
                  {{ artist.albumCount }} album{{ artist.albumCount !== 1 ? 's' : '' }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- ARTIST DETAIL -->
    <template v-else-if="currentArtist && !currentAlbum">
      <div class="px-8 py-7 border-b border-stone-200 bg-white flex-shrink-0">
        <div class="flex items-center gap-1.5 text-xs text-stone-400 mb-2">
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="closeArtist">Artists</span>
          <span class="opacity-40">›</span>
          <span>{{ currentArtist.name }}</span>
        </div>
        <h1 class="font-serif text-4xl font-semibold">{{ currentArtist.name }}</h1>
      </div>
      <div class="flex-1 overflow-y-auto px-8 py-6 pb-40 md:pb-24">
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <div v-else class="grid gap-5" style="grid-template-columns: repeat(auto-fill, minmax(168px, 1fr))">
          <div
            v-for="album in artistAlbums" :key="album.id"
            class="cursor-pointer group"
            @click="openAlbum(album)"
          >
            <div class="aspect-square bg-amber-50 mb-2.5 overflow-hidden relative">
              <img v-if="album.coverArt" :src="coverUrl(album.coverArt)" :alt="album.name" class="w-full h-full object-cover" @error="onImgError" />
              <div v-else class="w-full h-full flex items-center justify-center text-4xl">💿</div>
              <div class="absolute inset-0 bg-black/25 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                <button
                  class="w-11 h-11 rounded-full bg-white flex items-center justify-center text-base pl-0.5"
                  @click.stop="playAlbum(album)"
                >▶</button>
              </div>
            </div>
            <div class="text-sm font-medium truncate">{{ album.name }}</div>
            <div class="text-xs text-stone-400 mt-0.5">{{ album.year || '' }}</div>
          </div>
        </div>
      </div>
    </template>

    <!-- ALBUM DETAIL -->
    <template v-else-if="currentAlbum">
      <div class="px-8 py-7 border-b border-stone-200 bg-white flex-shrink-0">
        <div class="flex items-center gap-1.5 text-xs text-stone-400 mb-2 flex-wrap">
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="closeArtist">Artists</span>
          <span class="opacity-40">›</span>
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="closeAlbum">{{ currentArtist.name }}</span>
          <span class="opacity-40">›</span>
          <span>{{ currentAlbum.name }}</span>
        </div>
        <h1 class="font-serif text-4xl font-semibold">{{ currentAlbum.name }}</h1>
      </div>
      <div class="flex-1 overflow-y-auto px-8 py-6 pb-40 md:pb-24">
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <div v-else>
          <!-- ALBUM HEADER -->
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
              <div class="text-sm text-stone-400 mb-4 cursor-pointer hover:text-amber-700 transition-colors" @click="closeAlbum">
                {{ currentAlbum.artist }}
              </div>
              <div class="flex gap-2.5">
                <button class="bg-stone-900 text-white text-sm font-medium px-5 py-2 hover:bg-amber-700 transition-colors" @click="playAlbumTracks">▶ Play</button>
                <button class="border border-stone-200 text-sm px-4 py-2 hover:border-amber-700 hover:text-amber-700 transition-colors" @click="queueAlbumTracks">+ Queue</button>
              </div>
            </div>
          </div>

          <!-- TRACK LIST -->
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
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { usePlayerStore } from '../stores/player'
import { getArtists, getArtist, getAlbum, coverUrl } from '../api/subsonic'
import TrackItem from '../components/TrackItem.vue'

const player = usePlayerStore()
const route = useRoute()

const loading      = ref(false)
const artistIndex  = ref([])
const currentArtist = ref(null)
const artistAlbums  = ref([])
const currentAlbum  = ref(null)
const albumTracks   = ref([])

async function loadArtists() {
  loading.value = true
  try { artistIndex.value = await getArtists() }
  finally { loading.value = false }
}

async function openArtist(artist) {
  currentArtist.value = artist
  currentAlbum.value  = null
  loading.value = true
  try {
    const data = await getArtist(artist.id)
    artistAlbums.value = data.albums
  } finally { loading.value = false }
}

async function openArtistById(id) {
  if (!id) return
  loading.value = true
  try {
    const data = await getArtist(id)
    currentArtist.value = { id, name: data.info?.name || '' }
    artistAlbums.value = data.albums || []
    currentAlbum.value = null
  } finally { loading.value = false }
}

async function openAlbum(album) {
  currentAlbum.value = album
  loading.value = true
  try {
    const data = await getAlbum(album.id)
    albumTracks.value  = data.tracks
    currentAlbum.value = { ...album, ...data.info }
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

function closeArtist() {
  currentArtist.value = null
  currentAlbum.value  = null
}

function closeAlbum() {
  currentAlbum.value = null
}

function onImgError(e) {
  try { if (e?.target) e.target.style.display = 'none' } catch(_) {}
}

onMounted(async () => {
  await loadArtists()
  if (route.params.id) await openArtistById(route.params.id)
})

watch(() => route.params.id, (id) => {
  if (id) openArtistById(id)
  else {
    currentArtist.value = null
    currentAlbum.value = null
  }
})
</script>