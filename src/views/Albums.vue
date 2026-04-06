<template>
  <div class="flex flex-col h-full overflow-hidden">

    <!-- ALBUM LIST -->
    <template v-if="!currentAlbum">
      <div class="px-6 py-6 border-b border-stone-200 bg-white flex-shrink-0">
        <h1 class="font-serif text-3xl font-semibold">Albums</h1>
      </div>
      <div class="flex-1 overflow-y-auto px-4 py-4 pb-40 md:pb-24">
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <div v-else-if="!albums.length" class="flex flex-col items-center justify-center py-24 text-stone-400 gap-2">
          <span class="text-4xl">💿</span>
          <span class="font-serif text-lg">No albums found</span>
        </div>
        <div v-else class="grid gap-4" style="grid-template-columns: repeat(auto-fill, minmax(100px, 1fr))">
          <div
            v-for="album in albums" :key="album.id"
            class="cursor-pointer group"
            @click="openAlbum(album)"
          >
            <div class="aspect-square bg-amber-50 mb-1.5 overflow-hidden relative rounded-lg">
              <img v-if="album.coverArt" :src="coverUrl(album.coverArt)" :alt="album.name" class="w-full h-full object-cover" @error="onImgError" />
              <div v-else class="w-full h-full flex items-center justify-center text-3xl">💿</div>
              <div class="absolute inset-0 bg-black/25 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity rounded-lg">
                <button
                  class="w-9 h-9 rounded-full bg-white flex items-center justify-center text-sm pl-0.5"
                  @click.stop="playAlbum(album)"
                >▶</button>
              </div>
            </div>
            <div class="text-xs font-medium truncate leading-tight">{{ album.name }}</div>
            <div class="text-xs text-stone-400 truncate mt-0.5">{{ album.artist }}</div>
          </div>
        </div>
      </div>
    </template>

    <!-- ALBUM DETAIL -->
    <template v-else>
      <div class="px-6 py-6 border-b border-stone-200 bg-white flex-shrink-0">
        <div class="flex items-center gap-1.5 text-xs text-stone-400 mb-1">
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="closeAlbum">Albums</span>
          <span class="opacity-40">›</span>
          <span class="truncate">{{ currentAlbum.name }}</span>
        </div>
        <h1 class="font-serif text-3xl font-semibold">{{ currentAlbum.name }}</h1>
      </div>
      <div class="flex-1 overflow-y-auto px-4 py-4 pb-40 md:pb-24">
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <div v-else>

          <!-- ALBUM HEADER -->
          <div class="flex gap-4 mb-6 items-end">
            <div class="w-28 h-28 flex-shrink-0 bg-amber-50 overflow-hidden rounded-lg shadow-md">
              <img v-if="currentAlbum.coverArt" :src="coverUrl(currentAlbum.coverArt)" :alt="currentAlbum.name" class="w-full h-full object-cover" @error="onImgError" />
              <div v-else class="w-full h-full flex items-center justify-center text-4xl">💿</div>
            </div>
            <div class="flex-1 overflow-hidden">
              <div class="text-xs uppercase tracking-widest text-stone-400 mb-1">
                Album{{ currentAlbum.year ? ' · ' + currentAlbum.year : '' }}
              </div>
              <div class="font-serif text-xl font-semibold leading-tight mb-0.5">{{ currentAlbum.name }}</div>
              <div class="text-sm text-stone-400 mb-3 truncate">{{ currentAlbum.artist }}</div>
              <div class="flex gap-2">
                <button class="bg-stone-900 text-white text-xs font-medium px-4 py-2 rounded-full hover:bg-amber-700 transition-colors" @click="playAlbumTracks">▶ Play</button>
                <button class="border border-stone-200 text-xs px-4 py-2 rounded-full hover:border-amber-700 hover:text-amber-700 transition-colors" @click="queueAlbumTracks">+ Queue</button>
              </div>
            </div>
          </div>

          <!-- TRACK LIST HEADER -->
          <div class="grid gap-2 text-xs uppercase tracking-widest text-stone-400 px-3 pb-2 border-b border-stone-200 mb-1"
            style="grid-template-columns: 28px 1fr 44px 28px">
            <span class="text-center">#</span><span>Title</span><span>Time</span><span></span>
          </div>

          <!-- TRACKS -->
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
import { useRoute, useRouter } from 'vue-router'
import { usePlayerStore } from '../stores/player'
import { getAlbumList, getAlbum, coverUrl } from '../api/subsonic'
import TrackItem from '../components/TrackItem.vue'

const route  = useRoute()
const router = useRouter()
const player = usePlayerStore()

const loading      = ref(false)
const albums       = ref([])
const currentAlbum = ref(null)
const albumTracks  = ref([])

async function loadAlbums() {
  loading.value = true
  try { albums.value = await getAlbumList() }
  finally { loading.value = false }
}

async function openAlbum(album) {
  loading.value = true
  router.push({ name: 'album-detail', params: { id: album.id } })
  try {
    const data = await getAlbum(album.id)
    currentAlbum.value = { ...album, ...data.info }
    albumTracks.value  = data.tracks
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

function closeAlbum() {
  currentAlbum.value = null
  router.push({ name: 'albums' })
}

function onImgError(e) {
  try { if (e?.target) e.target.style.display = 'none' } catch(_) {}
}

onMounted(async () => {
  await loadAlbums()
  if (route.params.id) {
    const album = albums.value.find(a => a.id === route.params.id)
    if (album) openAlbum(album)
  }
})

watch(() => route.params.id, (id) => {
  if (!id) currentAlbum.value = null
})
</script>