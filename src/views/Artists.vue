<template>
  <div class="flex flex-col h-full overflow-hidden">

    <!-- ARTISTS GRID -->
    <template v-if="view === 'grid'">
      <div class="border-b border-stone-200 bg-white flex-shrink-0">
        <div class="px-4 md:px-8 pt-5 md:pt-7 pb-3">
          <h1 class="font-serif text-2xl md:text-4xl font-semibold">Artists</h1>
        </div>
        <div class="px-4 md:px-6 pb-2 flex flex-wrap gap-0.5">
          <button
            v-for="letter in LETTERS"
            :key="letter"
            type="button"
            tabindex="-1"
            class="min-w-[1.75rem] h-7 px-1 text-xs font-medium rounded transition-colors select-none"
            :class="expandedGroups[letter]
              ? 'bg-amber-700 text-white'
              : 'text-stone-500 hover:text-amber-700 hover:bg-amber-50'"
            @click="toggleAndScroll(letter)"
          >
            {{ letter }}
          </button>
        </div>
      </div>
      <div ref="scrollContainer" class="flex-1 overflow-y-auto px-6 py-4 pb-40 md:pb-24">
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <template v-else>
          <template v-for="group in artistIndex" :key="group.name">
          <div
            v-if="expandedGroups[group.name]"
            :ref="el => setGroupRef(group.name, el)"
            class="mb-6"
          >
            <div class="font-serif text-xl font-semibold text-amber-700 border-b border-stone-200 pb-1 mb-3">
              {{ group.name }}
            </div>
            <div class="grid gap-4" style="grid-template-columns: repeat(auto-fill, minmax(120px, 1fr))">
              <ArtistCard
                v-for="artist in group.artist"
                :key="artist.id"
                :artist="artist"
                @click="openArtist(artist)"
              />
            </div>
          </div>
          </template>
        </template>
      </div>
    </template>

    <!-- ARTIST DETAIL -->
    <template v-else-if="view === 'artist'">
      <div class="px-4 md:px-8 py-5 md:py-7 border-b border-stone-200 bg-white flex-shrink-0">
        <div class="flex items-center gap-1.5 text-xs text-stone-400 mb-3">
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="view = 'grid'">Artists</span>
          <span class="opacity-40">›</span>
          <span>{{ currentArtist.name }}</span>
        </div>
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 md:w-14 md:h-14 rounded-xl overflow-hidden flex-shrink-0 bg-stone-100">
            <img
              v-if="artistDetailImageUrl"
              :src="artistDetailImageUrl"
              :alt="currentArtist.name"
              class="w-full h-full object-cover"
              @error="onArtistDetailImgError"
            />
            <div v-else class="w-full h-full flex items-center justify-center font-serif text-2xl font-semibold text-stone-300 select-none">
              {{ currentArtist.name?.[0]?.toUpperCase() }}
            </div>
          </div>
          <h1 class="font-serif text-2xl md:text-4xl font-semibold">{{ currentArtist.name }}</h1>
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
              <div class="w-full h-full flex items-center justify-center text-2xl">💿</div>
              <img :src="coverUrl(album.coverArt || album.id)" :alt="album.name" class="absolute inset-0 w-full h-full object-cover" @error="e => e.target.style.display='none'" />
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
      <div class="px-4 md:px-8 py-5 md:py-7 border-b border-stone-200 bg-white flex-shrink-0">
        <div class="flex items-center gap-1.5 text-xs text-stone-400 mb-2 flex-wrap">
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="view = 'grid'">Artists</span>
          <span class="opacity-40">›</span>
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="view = 'artist'">{{ currentArtist.name }}</span>
          <span class="opacity-40">›</span>
          <span>{{ currentAlbum.name }}</span>
        </div>
        <h1 class="font-serif text-2xl md:text-4xl font-semibold truncate">{{ currentAlbum.name }}</h1>
      </div>
      <div class="flex-1 overflow-y-auto px-4 md:px-8 py-4 md:py-6 pb-40 md:pb-24">
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <div v-else>
          <div class="flex gap-6 mb-8 items-end">
            <div class="w-40 h-40 flex-shrink-0 bg-amber-50 overflow-hidden relative">
              <div class="w-full h-full flex items-center justify-center text-5xl">💿</div>
              <img :src="coverUrl(currentAlbum.coverArt || currentAlbum.id)" :alt="currentAlbum.name" class="absolute inset-0 w-full h-full object-cover" @error="e => e.target.style.display='none'" />
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

          <div class="grid gap-2 text-xs uppercase tracking-widest text-stone-400 px-3 pb-2 border-b border-stone-200 mb-1"
            style="grid-template-columns: 28px 1fr 44px 28px">
            <span class="text-center">#</span><span>Title</span><span>Time</span><span></span>
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
import { ref, reactive, nextTick, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePlayerStore } from '../stores/player'
import { getArtists, getArtist, getAlbum, coverUrl } from '../api/subsonic'
import TrackItem  from '../components/TrackItem.vue'
import ArtistCard from '../components/ArtistCard.vue'

const route  = useRoute()
const router = useRouter()
const player = usePlayerStore()

const LETTERS = ['#', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z']

const view    = ref('grid')
const loading = ref(false)

const artistIndex    = ref([])
const expandedGroups = reactive({})
const scrollContainer = ref(null)
const groupRefs = {}

function setGroupRef(name, el) {
  if (el) groupRefs[name] = el
}

const currentArtist           = ref(null)
const currentArtistAlbums     = ref([])
const loadingArtist           = ref(false)
const artistDetailImageUrl    = ref(null)

const currentAlbum = ref(null)
const albumTracks  = ref([])

async function loadArtists() {
  loading.value = true
  try {
    const index = await getArtists()
    for (const group of index) {
      expandedGroups[group.name] = false
    }
    artistIndex.value = index
  } finally {
    loading.value = false
  }

}

function toggleAndScroll(name) {
  const opening = !expandedGroups[name]
  expandedGroups[name] = opening
  if (opening) {
    nextTick(() => {
      const el = groupRefs[name]
      const container = scrollContainer.value
      if (el && container) {
        container.scrollTo({ top: el.offsetTop - 8, behavior: 'smooth' })
      }
    })
  }
}

async function openArtist(artist) {
  currentArtist.value = { ...artist }
  artistDetailImageUrl.value = null
  view.value = 'artist'
  loadingArtist.value = true
  try {
    const data = await getArtist(artist.id)
    currentArtistAlbums.value = data.albums
    const coverArtId = data.info?.coverArt || artist.coverArt || artist.id
    artistDetailImageUrl.value = coverUrl(coverArtId, 112)
  } finally { loadingArtist.value = false }
}

function onArtistDetailImgError() {
  console.warn(`[artist image] no image for "${currentArtist.value?.name}" (url: ${artistDetailImageUrl.value})`)
  artistDetailImageUrl.value = null
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

async function openArtistById(id) {
  currentArtist.value         = { id, name: '' }
  artistDetailImageUrl.value  = null
  view.value                  = 'artist'
  loadingArtist.value         = true
  try {
    const data = await getArtist(id)
    currentArtist.value        = data.info
    currentArtistAlbums.value  = data.albums
    const coverArtId = data.info?.coverArt || id
    artistDetailImageUrl.value = coverUrl(coverArtId, 112)
  } finally { loadingArtist.value = false }
}

watch(() => route.params.id, (id) => {
  if (id) openArtistById(id)
  else view.value = 'grid'
})

onMounted(async () => {
  await loadArtists()
  if (route.params.id) openArtistById(route.params.id)
})
</script>
