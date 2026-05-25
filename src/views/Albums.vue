<template>
  <div class="flex flex-col flex-1 min-h-0 overflow-hidden">

    <!-- ALBUM LIST -->
    <template v-if="!currentAlbum">
      <div class="px-6 py-6 border-b border-stone-200 bg-white flex-shrink-0">
        <h1 class="font-serif text-3xl font-semibold mb-3">Albums</h1>
        <div v-if="albums.length" class="flex items-center gap-2 flex-wrap">

          <!-- Genre -->
          <div v-if="distinctGenres.length" class="relative">
            <div v-if="showGenreDropdown" class="fixed inset-0 z-10" @click="showGenreDropdown = false"></div>
            <button
              class="relative z-20 text-xs border px-2.5 py-1.5 rounded transition-all"
              :class="filterGenre ? 'border-amber-700 text-amber-700 bg-amber-50' : 'border-stone-200 text-stone-400 hover:border-amber-700 hover:text-amber-700 hover:bg-amber-50'"
              @click="showGenreDropdown = !showGenreDropdown; showYearDropdown = false"
            >{{ filterGenre || 'Genre' }} ▾</button>
            <div v-if="showGenreDropdown" class="absolute top-full left-0 mt-1 bg-white border border-stone-200 shadow-md rounded z-20 min-w-[130px] max-h-60 overflow-y-auto">
              <button
                class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                :class="!filterGenre ? 'text-amber-700 font-medium' : 'text-stone-600'"
                @click="filterGenre = null; showGenreDropdown = false"
              >All genres</button>
              <button
                v-for="g in distinctGenres" :key="g"
                class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                :class="filterGenre === g ? 'text-amber-700 font-medium' : 'text-stone-600'"
                @click="filterGenre = g; showGenreDropdown = false"
              >{{ g }}</button>
            </div>
          </div>

          <!-- Year -->
          <div v-if="distinctYears.length" class="relative">
            <div v-if="showYearDropdown" class="fixed inset-0 z-10" @click="showYearDropdown = false"></div>
            <button
              class="relative z-20 text-xs border px-2.5 py-1.5 rounded transition-all"
              :class="filterYear ? 'border-amber-700 text-amber-700 bg-amber-50' : 'border-stone-200 text-stone-400 hover:border-amber-700 hover:text-amber-700 hover:bg-amber-50'"
              @click="showYearDropdown = !showYearDropdown; showGenreDropdown = false"
            >{{ filterYear || 'Year' }} ▾</button>
            <div v-if="showYearDropdown" class="absolute top-full left-0 mt-1 bg-white border border-stone-200 shadow-md rounded z-20 min-w-[80px] max-h-60 overflow-y-auto">
              <button
                class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                :class="!filterYear ? 'text-amber-700 font-medium' : 'text-stone-600'"
                @click="filterYear = null; showYearDropdown = false"
              >All years</button>
              <button
                v-for="y in distinctYears" :key="y"
                class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                :class="filterYear === y ? 'text-amber-700 font-medium' : 'text-stone-600'"
                @click="filterYear = y; showYearDropdown = false"
              >{{ y }}</button>
            </div>
          </div>

          <button
            v-if="filterGenre || filterYear"
            class="text-xs text-stone-400 hover:text-amber-700 transition-colors"
            @click="filterGenre = null; filterYear = null"
          >Clear</button>
        </div>
      </div>
      <div class="flex-1 min-h-0 overflow-y-auto pb-40 md:pb-24">

        <!-- RECENTLY ADDED CAROUSEL -->
        <div v-if="recentAlbums.length" class="px-4 pt-5 pb-4 border-b border-stone-100">
          <div class="text-xs font-medium uppercase tracking-widest text-stone-400 mb-3">Recently Added</div>
          <div class="relative">
            <button
              v-if="canScrollLeft"
              class="absolute left-0 top-1/2 -translate-y-1/2 -translate-x-2 z-10 w-7 h-7 rounded-full bg-white shadow-md border border-stone-200 flex items-center justify-center text-stone-500 hover:text-amber-700 transition-colors"
              @click="scrollCarousel(-1)"
            >‹</button>
            <div ref="carousel" class="flex gap-3 overflow-x-auto scroll-smooth" style="scrollbar-width: none;" @scroll="onCarouselScroll">
              <div
                v-for="album in recentAlbums" :key="album.id"
                class="flex-shrink-0 w-28 cursor-pointer group"
                @click="openAlbum(album)"
              >
                <div class="aspect-square bg-amber-50 mb-1.5 overflow-hidden relative rounded-lg">
                  <div class="w-full h-full flex items-center justify-center text-3xl">💿</div>
                  <img :src="coverUrl(album.coverArt || album.id)" :alt="album.name" class="absolute inset-0 w-full h-full object-cover" @error="onAlbumCoverError($event, album)" />
                  <div class="absolute inset-0 bg-black/25 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity rounded-lg">
                    <button class="w-8 h-8 rounded-full bg-white flex items-center justify-center text-sm pl-0.5" @click.stop="playAlbum(album)">▶</button>
                  </div>
                </div>
                <div class="text-xs font-medium truncate leading-tight">{{ album.name }}</div>
                <div class="text-xs text-stone-400 truncate mt-0.5">{{ album.artist }}</div>
              </div>
            </div>
            <button
              v-if="canScrollRight"
              class="absolute right-0 top-1/2 -translate-y-1/2 translate-x-2 z-10 w-7 h-7 rounded-full bg-white shadow-md border border-stone-200 flex items-center justify-center text-stone-500 hover:text-amber-700 transition-colors"
              @click="scrollCarousel(1)"
            >›</button>
          </div>
        </div>

        <!-- DISCOVER -->
        <div v-if="discoverAlbums.length" class="px-4 pt-5 pb-4 border-b border-stone-100">
          <div class="text-xs font-medium uppercase tracking-widest text-stone-400 mb-3">Discover</div>
          <div class="flex gap-3 overflow-x-auto" style="scrollbar-width: none;">
            <div
              v-for="album in discoverAlbums" :key="album.id"
              class="flex-shrink-0 w-28 cursor-pointer group"
              @click="openAlbum(album)"
            >
              <div class="aspect-square bg-amber-50 mb-1.5 overflow-hidden relative rounded-lg">
                <div class="w-full h-full flex items-center justify-center text-3xl">💿</div>
                <img :src="coverUrl(album.coverArt || album.id)" :alt="album.name" class="absolute inset-0 w-full h-full object-cover" @error="onAlbumCoverError($event, album)" />
                <div class="absolute inset-0 bg-black/25 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity rounded-lg">
                  <button class="w-8 h-8 rounded-full bg-white flex items-center justify-center text-sm pl-0.5" @click.stop="playAlbum(album)">▶</button>
                </div>
              </div>
              <div class="text-xs font-medium truncate leading-tight">{{ album.name }}</div>
              <div class="text-xs text-stone-400 truncate mt-0.5">{{ album.albumArtist || album.artist }}</div>
            </div>
          </div>
        </div>

        <!-- ALL ALBUMS GRID -->
        <div class="px-4 py-4">
          <div v-if="loading && !albums.length" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
          <div v-else-if="!albums.length" class="flex flex-col items-center justify-center py-24 text-stone-400 gap-2">
            <span class="text-4xl">💿</span>
            <span class="font-serif text-lg">No albums found</span>
          </div>
          <div v-else-if="!filteredAlbums.length" class="flex items-center justify-center py-24 text-stone-400 text-sm">No albums match the filter.</div>
          <div v-else class="grid gap-4" style="grid-template-columns: repeat(auto-fill, minmax(100px, 1fr))">
            <div
              v-for="album in filteredAlbums" :key="album.id"
              class="cursor-pointer group"
              @click="openAlbum(album)"
            >
              <div class="aspect-square bg-amber-50 mb-1.5 overflow-hidden relative rounded-lg">
                <div class="w-full h-full flex items-center justify-center text-3xl">💿</div>
                <img :src="coverUrl(album.coverArt || album.id)" :alt="album.name" class="absolute inset-0 w-full h-full object-cover" @error="onAlbumCoverError($event, album)" />
                <div class="absolute inset-0 bg-black/25 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity rounded-lg">
                  <button
                    class="w-9 h-9 rounded-full bg-white flex items-center justify-center text-sm pl-0.5"
                    @click.stop="playAlbum(album)"
                  >▶</button>
                </div>
              </div>
              <div class="text-xs font-medium truncate leading-tight">{{ album.name }}</div>
              <div class="text-xs text-stone-400 truncate mt-0.5">
                <button class="text-stone-400 hover:text-amber-700" @click.stop="openArtistByName(album.artist, $event)">{{ album.artist }}</button>
              </div>
            </div>
          </div>
          <div ref="sentinel" class="h-8"></div>
          <div v-if="loading && albums.length" class="flex items-center justify-center py-4 text-stone-400 text-sm">Loading…</div>
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
      <div class="flex-1 min-h-0 overflow-y-auto px-4 py-4 pb-40 md:pb-24">
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <div v-else>
          <div class="flex gap-4 mb-6 items-end">
            <div class="w-28 h-28 flex-shrink-0 bg-amber-50 overflow-hidden rounded-lg shadow-md relative">
              <div class="w-full h-full flex items-center justify-center text-4xl">💿</div>
              <img :src="coverUrl(currentAlbum.coverArt || currentAlbum.id)" :alt="currentAlbum.name" class="absolute inset-0 w-full h-full object-cover" @error="onAlbumCoverError($event, currentAlbum)" />
            </div>
            <div class="flex-1 overflow-hidden">
              <div class="text-xs uppercase tracking-widest text-stone-400 mb-1">
                Album{{ currentAlbum.year ? ' · ' + currentAlbum.year : '' }}
              </div>
              <div class="font-serif text-xl font-semibold leading-tight mb-0.5">{{ currentAlbum.name }}</div>
              <div class="text-sm text-stone-400 mb-3 truncate">
                <button class="hover:text-amber-700 transition-colors" @click="openArtistByName(currentAlbum.artist, $event)">{{ currentAlbum.artist }}</button>
              </div>
              <div class="flex gap-2">
                <button class="bg-stone-900 text-white text-xs font-medium px-4 py-2 rounded-full hover:bg-amber-700 transition-colors" @click="playAlbumTracks">▶ Play</button>
                <button class="border border-stone-200 text-xs px-4 py-2 rounded-full hover:border-amber-700 hover:text-amber-700 transition-colors" @click="queueAlbumTracks">+ Queue</button>
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
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePlayerStore } from '../stores/player'
import { getAlbumPage, getNewestAlbums, getRandomAlbums, getAlbum, coverUrl, getArtists } from '../api/subsonic'
import TrackItem from '../components/TrackItem.vue'

const route  = useRoute()
const router = useRouter()
const player = usePlayerStore()

const loading        = ref(false)
const albums         = ref([])
const recentAlbums   = ref([])
const discoverAlbums = ref([])
const currentAlbum = ref(null)
const albumTracks  = ref([])

const filterGenre       = ref(null)
const filterYear        = ref(null)
const showGenreDropdown = ref(false)
const showYearDropdown  = ref(false)

const distinctGenres = computed(() => {
  const genres = [...new Set(albums.value.map(a => a.genre).filter(Boolean))]
  return genres.sort()
})

const distinctYears = computed(() => {
  const years = [...new Set(albums.value.map(a => a.year).filter(Boolean))]
  return years.sort((a, b) => b - a)
})

const filteredAlbums = computed(() => {
  if (!filterGenre.value && !filterYear.value) return albums.value
  return albums.value.filter(a => {
    if (filterGenre.value && a.genre !== filterGenre.value) return false
    if (filterYear.value  && a.year  !== filterYear.value)  return false
    return true
  })
})

const sentinel   = ref(null)
const allLoaded  = ref(false)
const pageSize   = 100
let   pageOffset = 0
let   observer   = null

const carousel      = ref(null)
const canScrollLeft  = ref(false)
const canScrollRight = ref(false)

function onCarouselScroll() {
  if (!carousel.value) return
  canScrollLeft.value  = carousel.value.scrollLeft > 0
  canScrollRight.value = carousel.value.scrollLeft + carousel.value.clientWidth < carousel.value.scrollWidth - 1
}

function scrollCarousel(dir) {
  if (!carousel.value) return
  carousel.value.scrollBy({ left: dir * 320, behavior: 'smooth' })
}

async function loadRecent() {
  const [recent, random] = await Promise.all([getNewestAlbums(20), getRandomAlbums(20)])
  recentAlbums.value   = recent
  discoverAlbums.value = random
  await nextTick()
  onCarouselScroll()
}

async function loadAlbums() {
  albums.value = []
  allLoaded.value = false
  pageOffset = 0
  await loadMore()
  await nextTick()
  setupObserver()
}

async function loadMore() {
  if (loading.value || allLoaded.value) return
  loading.value = true
  try {
    const page = await getAlbumPage(pageSize, pageOffset)
    albums.value.push(...page)
    if (page.length < pageSize) allLoaded.value = true
    else pageOffset += pageSize
  } finally {
    loading.value = false
  }
}

function setupObserver() {
  if (observer) observer.disconnect()
  if (!sentinel.value) return
  observer = new IntersectionObserver(entries => {
    if (entries[0].isIntersecting) loadMore()
  }, { threshold: 0.1 })
  observer.observe(sentinel.value)
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

async function openArtistByName(name, e) {
  if (e && e.stopPropagation) e.stopPropagation()
  try {
    const idx = await getArtists()
    for (const g of idx) {
      for (const a of g.artist) {
        if (a.name === name) {
          router.push({ name: 'artist-detail', params: { id: a.id } })
          return
        }
      }
    }
  } catch (err) {
    // ignore
  }
  router.push({ name: 'artists' })
}

function closeAlbum() {
  currentAlbum.value = null
  router.push({ name: 'albums' })
}

function onAlbumCoverError(e, album) {
  const img = e.target
  if (!img.dataset.triedSidecar) {
    img.dataset.triedSidecar = '1'
    const artist = album.albumArtist || album.artist
    img.src = `/artist-images/album?artist=${encodeURIComponent(artist)}&album=${encodeURIComponent(album.name)}`
  } else {
    img.style.display = 'none'
  }
}

onMounted(async () => {
  loadRecent()
  await loadAlbums()
  if (route.params.id) {
    const album = albums.value.find(a => a.id === route.params.id)
    if (album) openAlbum(album)
    else openAlbumById(route.params.id)
  }
})

onUnmounted(() => {
  if (observer) observer.disconnect()
})

async function openAlbumById(id) {
  loading.value = true
  router.push({ name: 'album-detail', params: { id } })
  try {
    const data = await getAlbum(id)
    currentAlbum.value = data.info
    albumTracks.value  = data.tracks
  } finally { loading.value = false }
}

watch(() => route.params.id, (id) => {
  if (!id) {
    currentAlbum.value = null
  } else {
    const album = albums.value.find(a => a.id === id)
    if (album) openAlbum(album)
    else openAlbumById(id)
  }
})
</script>