<template>
  <div class="flex flex-col h-full overflow-hidden">

    <!-- ARTISTS GRID -->
    <template v-if="view === 'grid'">

      <!-- ── DESKTOP HEADER (letter nav + filters) ── -->
      <div class="hidden md:block border-b border-stone-200 bg-white flex-shrink-0">
        <div class="px-8 pt-7 pb-3">
          <h1 class="font-serif text-4xl font-semibold mb-3">Artists</h1>
          <div v-if="!loading" class="flex items-center gap-2 flex-wrap">
            <!-- Genre -->
            <div v-if="distinctGenres.length" class="relative">
              <div v-if="showGenreDropdown" class="fixed inset-0 z-10" @click="showGenreDropdown = false"></div>
              <button
                class="relative z-20 text-xs border px-2.5 py-1.5 rounded transition-all"
                :class="filterGenre ? 'border-amber-700 text-amber-700 bg-amber-50' : 'border-stone-200 text-stone-400 hover:border-amber-700 hover:text-amber-700 hover:bg-amber-50'"
                @click="showGenreDropdown = !showGenreDropdown; showYearDropdown = false"
              >{{ filterGenre || 'Genre' }} ▾</button>
              <div v-if="showGenreDropdown" class="absolute top-full left-0 mt-1 bg-white border border-stone-200 shadow-md rounded z-20 min-w-[130px] max-h-60 overflow-y-auto">
                <button class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                  :class="!filterGenre ? 'text-amber-700 font-medium' : 'text-stone-600'"
                  @click="filterGenre = null; showGenreDropdown = false">All genres</button>
                <button v-for="g in distinctGenres" :key="g"
                  class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                  :class="filterGenre === g ? 'text-amber-700 font-medium' : 'text-stone-600'"
                  @click="filterGenre = g; showGenreDropdown = false">{{ g }}</button>
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
                <button class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                  :class="!filterYear ? 'text-amber-700 font-medium' : 'text-stone-600'"
                  @click="filterYear = null; showYearDropdown = false">All years</button>
                <button v-for="y in distinctYears" :key="y"
                  class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                  :class="filterYear === y ? 'text-amber-700 font-medium' : 'text-stone-600'"
                  @click="filterYear = y; showYearDropdown = false">{{ y }}</button>
              </div>
            </div>
            <button v-if="filterGenre || filterYear"
              class="text-xs text-stone-400 hover:text-amber-700 transition-colors"
              @click="filterGenre = null; filterYear = null">Clear</button>
          </div>
        </div>
        <div class="px-6 pb-2 flex flex-wrap gap-0.5">
          <button
            v-for="letter in LETTERS" :key="letter"
            type="button" tabindex="-1"
            class="min-w-[1.75rem] h-7 px-1 text-xs font-medium rounded transition-colors select-none"
            :class="expandedGroups[letter] ? 'bg-amber-700 text-white' : 'text-stone-500 hover:text-amber-700 hover:bg-amber-50'"
            @click="toggleAndScroll(letter)"
          >{{ letter }}</button>
        </div>
      </div>

      <!-- ── MOBILE HEADER (search + filters in one row) ── -->
      <div class="md:hidden border-b border-stone-200 bg-white flex-shrink-0 px-4 py-3">
        <div class="flex items-center gap-2">
          <!-- Search -->
          <div class="relative flex-1 min-w-0">
            <input
              v-model="mobileQuery"
              type="search"
              placeholder="Search artists…"
              class="w-full text-sm px-3 py-2 rounded-lg border border-stone-200 bg-stone-50 placeholder-stone-400 focus:outline-none focus:border-amber-400 transition-colors pr-6"
            />
            <button v-if="mobileQuery" class="absolute right-2 top-1/2 -translate-y-1/2 text-stone-400 text-xs" @click="mobileQuery = ''">✕</button>
          </div>
          <!-- Genre filter -->
          <div class="relative flex-shrink-0">
            <div v-if="showGenreDropdown" class="fixed inset-0 z-10" @click="showGenreDropdown = false"></div>
            <button
              class="relative z-20 text-xs border px-2 py-2 rounded transition-all whitespace-nowrap"
              :class="filterGenre ? 'border-amber-700 text-amber-700 bg-amber-50' : 'border-stone-200 text-stone-400'"
              @click="showGenreDropdown = !showGenreDropdown; showYearDropdown = false"
            >{{ filterGenre ? (filterGenre.length > 7 ? filterGenre.slice(0, 7) + '…' : filterGenre) : 'Genre' }} ▾</button>
            <div v-if="showGenreDropdown" class="absolute top-full right-0 mt-1 bg-white border border-stone-200 shadow-md rounded z-20 min-w-[130px] max-h-60 overflow-y-auto">
              <button class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                :class="!filterGenre ? 'text-amber-700 font-medium' : 'text-stone-600'"
                @click="filterGenre = null; showGenreDropdown = false">All genres</button>
              <button v-for="g in distinctGenres" :key="g"
                class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                :class="filterGenre === g ? 'text-amber-700 font-medium' : 'text-stone-600'"
                @click="filterGenre = g; showGenreDropdown = false">{{ g }}</button>
            </div>
          </div>
          <!-- Year filter -->
          <div class="relative flex-shrink-0">
            <div v-if="showYearDropdown" class="fixed inset-0 z-10" @click="showYearDropdown = false"></div>
            <button
              class="relative z-20 text-xs border px-2 py-2 rounded transition-all"
              :class="filterYear ? 'border-amber-700 text-amber-700 bg-amber-50' : 'border-stone-200 text-stone-400'"
              @click="showYearDropdown = !showYearDropdown; showGenreDropdown = false"
            >{{ filterYear || 'Year' }} ▾</button>
            <div v-if="showYearDropdown" class="absolute top-full right-0 mt-1 bg-white border border-stone-200 shadow-md rounded z-20 min-w-[80px] max-h-60 overflow-y-auto">
              <button class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                :class="!filterYear ? 'text-amber-700 font-medium' : 'text-stone-600'"
                @click="filterYear = null; showYearDropdown = false">All years</button>
              <button v-for="y in distinctYears" :key="y"
                class="block w-full text-left px-3 py-1.5 text-xs hover:bg-stone-50 transition-colors"
                :class="filterYear === y ? 'text-amber-700 font-medium' : 'text-stone-600'"
                @click="filterYear = y; showYearDropdown = false">{{ y }}</button>
            </div>
          </div>
        </div>
      </div>

      <!-- ── DESKTOP CONTENT ── -->
      <div ref="scrollContainer" class="hidden md:block flex-1 overflow-y-auto px-6 py-4 pb-24">
        <!-- Recently Added -->
        <div v-if="recentArtists.length" class="mb-8">
          <h2 class="font-serif text-xl font-semibold mb-3">Recently Added</h2>
          <div class="flex gap-3 overflow-x-auto pb-2" style="scrollbar-width:none;-ms-overflow-style:none">
            <div v-for="artist in recentArtists" :key="artist.id" class="flex-none w-24 cursor-pointer group" @click="openArtist(artist)">
              <div class="aspect-square bg-stone-100 rounded-xl overflow-hidden mb-2 transition-transform duration-200 group-hover:scale-[1.03] relative">
                <div class="w-full h-full flex items-center justify-center font-serif text-3xl font-semibold text-stone-300 select-none">{{ artist.name[0]?.toUpperCase() }}</div>
                <img :src="`/artist-images/avatar?name=${encodeURIComponent(artist.name)}`" :alt="artist.name" class="absolute inset-0 w-full h-full object-cover" @error="e => e.target.style.display='none'" />
              </div>
              <div class="text-sm font-medium truncate leading-tight">{{ artist.name }}</div>
            </div>
          </div>
        </div>
        <!-- Discover Artists -->
        <div v-if="discoverArtists.length" class="mb-8">
          <h2 class="font-serif text-xl font-semibold mb-3">Discover Artists</h2>
          <div class="flex gap-3 overflow-x-auto pb-2" style="scrollbar-width:none;-ms-overflow-style:none">
            <div v-for="artist in discoverArtists" :key="artist.id" class="flex-none w-24 cursor-pointer group" @click="openArtist(artist)">
              <div class="aspect-square bg-stone-100 rounded-xl overflow-hidden mb-2 transition-transform duration-200 group-hover:scale-[1.03] relative">
                <div class="w-full h-full flex items-center justify-center font-serif text-3xl font-semibold text-stone-300 select-none">{{ artist.name[0]?.toUpperCase() }}</div>
                <img :src="`/artist-images/avatar?name=${encodeURIComponent(artist.name)}`" :alt="artist.name" class="absolute inset-0 w-full h-full object-cover" @error="e => e.target.style.display='none'" />
              </div>
              <div class="text-sm font-medium truncate leading-tight">{{ artist.name }}</div>
            </div>
          </div>
        </div>
        <div v-if="recentArtists.length || discoverArtists.length" class="border-b border-stone-200 mb-6"></div>
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <div v-else-if="(filterGenre || filterYear) && !filteredArtists.length" class="flex items-center justify-center py-24 text-stone-400 text-sm">No artists match the filter.</div>
        <template v-else>
          <template v-for="group in filteredArtistIndex" :key="group.name">
            <div v-if="expandedGroups[group.name]" :ref="el => setGroupRef(group.name, el)" class="mb-6">
              <div class="font-serif text-xl font-semibold text-amber-700 border-b border-stone-200 pb-1 mb-3">{{ group.name }}</div>
              <div class="grid gap-4" style="grid-template-columns: repeat(auto-fill, minmax(120px, 1fr))">
                <ArtistCard v-for="artist in group.artist" :key="artist.id" :artist="artist" @click="openArtist(artist)" />
              </div>
            </div>
          </template>
        </template>
      </div>

      <!-- ── MOBILE CONTENT ── -->
      <div class="md:hidden flex-1 overflow-y-auto pb-40">
        <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
        <template v-else>

          <!-- BROWSE MODE: 2 carousels (no search, no filter active) -->
          <template v-if="!mobileQuery && !filterGenre && !filterYear">
            <div v-if="recentArtists.length" class="pt-5 mb-6">
              <div class="px-4 text-xs font-medium uppercase tracking-widest text-stone-400 mb-3">Recently Added</div>
              <div class="w-full flex gap-3 overflow-x-auto px-4 pb-1" style="scrollbar-width:none;-ms-overflow-style:none">
                <div v-for="artist in recentArtists" :key="artist.id" class="flex-none cursor-pointer" style="width:calc(100% / 3 - 8px)" @click="openArtist(artist)">
                  <div class="aspect-square bg-stone-100 rounded-xl overflow-hidden mb-1.5 relative">
                    <div class="w-full h-full flex items-center justify-center font-serif text-3xl font-semibold text-stone-300 select-none">{{ artist.name[0]?.toUpperCase() }}</div>
                    <img :src="`/artist-images/avatar?name=${encodeURIComponent(artist.name)}`" :alt="artist.name" class="absolute inset-0 w-full h-full object-cover" @error="e => e.target.style.display='none'" />
                  </div>
                  <div class="text-xs font-medium truncate leading-tight text-center">{{ artist.name }}</div>
                </div>
              </div>
            </div>
            <div v-if="discoverArtists.length" class="mb-6">
              <div class="px-4 text-xs font-medium uppercase tracking-widest text-stone-400 mb-3">Discover Artists</div>
              <div class="w-full flex gap-3 overflow-x-auto px-4 pb-1" style="scrollbar-width:none;-ms-overflow-style:none">
                <div v-for="artist in discoverArtists" :key="artist.id" class="flex-none cursor-pointer" style="width:calc(100% / 3 - 8px)" @click="openArtist(artist)">
                  <div class="aspect-square bg-stone-100 rounded-xl overflow-hidden mb-1.5 relative">
                    <div class="w-full h-full flex items-center justify-center font-serif text-3xl font-semibold text-stone-300 select-none">{{ artist.name[0]?.toUpperCase() }}</div>
                    <img :src="`/artist-images/avatar?name=${encodeURIComponent(artist.name)}`" :alt="artist.name" class="absolute inset-0 w-full h-full object-cover" @error="e => e.target.style.display='none'" />
                  </div>
                  <div class="text-xs font-medium truncate leading-tight text-center">{{ artist.name }}</div>
                </div>
              </div>
            </div>
          </template>

          <!-- SEARCH / FILTER MODE: alphabetical contact list -->
          <template v-else>
            <div v-if="mobileFilteredIndex.length">
              <template v-for="group in mobileFilteredIndex" :key="group.name">
                <div class="px-4 py-1.5 text-xs font-semibold uppercase tracking-widest text-amber-700 bg-stone-50 border-b border-stone-100">{{ group.name }}</div>
                <div
                  v-for="artist in group.artist" :key="artist.id"
                  class="flex items-center gap-3 px-4 py-2 border-b border-stone-100 active:bg-stone-50 cursor-pointer"
                  @click="openArtist(artist)"
                >
                  <div class="w-11 h-11 rounded-full overflow-hidden flex-shrink-0 bg-stone-100 relative">
                    <div class="w-full h-full flex items-center justify-center font-semibold text-stone-300 select-none">{{ artist.name[0]?.toUpperCase() }}</div>
                    <img :src="`/artist-images/avatar?name=${encodeURIComponent(artist.name)}`" :alt="artist.name" class="absolute inset-0 w-full h-full object-cover" @error="e => e.target.style.display='none'" />
                  </div>
                  <span class="text-sm font-medium truncate">{{ artist.name }}</span>
                </div>
              </template>
            </div>
            <div v-else class="flex items-center justify-center py-24 text-stone-400 text-sm">No artists found</div>
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
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="goToLetter(currentArtistLetter)">{{ currentArtistLetter }}</span>
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
            <div class="text-xs text-stone-400 mt-0.5 truncate">
              {{ [album.year, album.songCount ? album.songCount + ' tracks' : '', album.genre].filter(Boolean).join(' · ') }}
            </div>
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
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="goToLetter(currentArtistLetter)">{{ currentArtistLetter }}</span>
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
              <div class="text-xs uppercase tracking-widest text-stone-400 mb-2 flex items-center flex-wrap gap-1">
                <span>{{ ['Album', currentAlbum.year].filter(Boolean).join(' · ') }}</span>
                <template v-if="editingGenre">
                  <span>·</span>
                  <div class="relative normal-case tracking-normal">
                    <input
                      ref="genreInputEl"
                      v-model="genreInput"
                      class="text-xs border-b border-stone-400 bg-transparent outline-none w-36"
                      placeholder="Genre…"
                      @keydown.enter="saveGenre"
                      @keydown.escape="editingGenre = false"
                      @focus="showGenreList = true"
                      @input="showGenreList = true"
                    />
                    <div v-if="showGenreList && filteredGenreList.length" class="absolute top-full left-0 mt-1 bg-white border border-stone-200 shadow-lg rounded z-30 w-48 max-h-52 overflow-y-auto">
                      <button
                        v-for="g in filteredGenreList" :key="g"
                        class="block w-full text-left px-3 py-1.5 text-xs text-stone-600 hover:bg-amber-50 hover:text-amber-700 transition-colors"
                        @mousedown.prevent="selectGenre(g)"
                      >{{ g }}</button>
                    </div>
                  </div>
                  <button class="normal-case tracking-normal text-amber-700 hover:text-amber-800 font-medium text-xs" @click="saveGenre">Save</button>
                  <button class="normal-case tracking-normal text-stone-400 hover:text-stone-600 text-xs" @click="editingGenre = false">Cancel</button>
                </template>
                <template v-else>
                  <button
                    class="normal-case tracking-normal flex items-center gap-1 group hover:text-amber-700 transition-colors cursor-pointer"
                    @click="startEditGenre"
                  >
                    <span>· {{ albumGenre || 'Add genre' }}</span>
                    <svg class="w-3.5 h-3.5 opacity-40 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536M9 13l6.586-6.586a2 2 0 012.828 2.828L11.828 15.828a2 2 0 01-1.414.586H9v-2a2 2 0 01.586-1.414z" />
                    </svg>
                  </button>
                </template>
              </div>
              <div class="font-serif text-3xl font-semibold mb-1.5">{{ currentAlbum.name }}</div>
              <div class="text-sm text-stone-400 mb-4">{{ currentAlbum.artist }}</div>
              <div class="flex gap-2.5">
                <button class="bg-stone-900 text-white text-sm font-medium px-5 py-2 hover:bg-amber-700 transition-colors" @click="playAlbumTracks">▶ Play</button>
                <button class="border border-stone-200 text-sm px-4 py-2 hover:border-amber-700 hover:text-amber-700 transition-colors" @click="queueAlbumTracks">+ Queue</button>
              </div>
            </div>
          </div>

          <div class="grid gap-2 text-xs uppercase tracking-widest text-stone-400 px-3 pb-2 border-b border-stone-200 mb-1 [grid-template-columns:28px_1fr_44px_28px] md:[grid-template-columns:28px_1fr_1fr_44px_28px]">
            <span class="text-center">#</span><span>Title</span><span class="hidden md:block">Artist</span><span>Time</span><span></span>
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
import { ref, reactive, computed, nextTick, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePlayerStore } from '../stores/player'
import { getArtists, getArtist, getAlbum, coverUrl, getNewestAlbums, getArtistGenreMap } from '../api/subsonic'
import { GENRES } from '../api/genres'
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
const currentArtistLetter     = ref(null)

function getArtistLetter(name) {
  const stripped = name?.replace(/^(the|a|an)\s+/i, '') ?? name
  const first = stripped?.[0]?.toLowerCase()
  return (first && first >= 'a' && first <= 'z') ? first : '#'
}

function goToLetter(letter) {
  view.value = 'grid'
  nextTick(() => {
    expandedGroups[letter] = true
    nextTick(() => {
      const el = groupRefs[letter]
      const container = scrollContainer.value
      if (el && container) container.scrollTo({ top: el.offsetTop - 8, behavior: 'smooth' })
    })
  })
}

const recentArtists   = ref([])
const discoverArtists = ref([])
const mobileQuery     = ref('')

const mobileFilteredIndex = computed(() => {
  const q = mobileQuery.value.toLowerCase().trim()
  const base = filteredArtistIndex.value
  if (!q) return base
  return base
    .map(g => ({ ...g, artist: g.artist.filter(a => a.name.toLowerCase().includes(q)) }))
    .filter(g => g.artist.length)
})

const artistGenreMap    = ref({})  // artistId -> { genres: Set, years: Set }
const filterGenre       = ref(null)
const filterYear        = ref(null)
const showGenreDropdown = ref(false)
const showYearDropdown  = ref(false)

const distinctGenres = computed(() => {
  const seen = new Set()
  for (const meta of Object.values(artistGenreMap.value)) {
    for (const g of meta.genres) seen.add(g)
  }
  return [...seen].sort()
})

const distinctYears = computed(() => {
  const seen = new Set()
  for (const meta of Object.values(artistGenreMap.value)) {
    for (const y of meta.years) seen.add(y)
  }
  return [...seen].sort((a, b) => Number(b) - Number(a))
})

const filteredArtists = computed(() => {
  if (!filterGenre.value && !filterYear.value) return null
  return artistIndex.value.flatMap(g => g.artist).filter(a => {
    const meta = artistGenreMap.value[a.id]
    if (!meta) return false
    if (filterGenre.value && !meta.genres.includes(filterGenre.value)) return false
    if (filterYear.value  && !meta.years.includes(filterYear.value))   return false
    return true
  })
})

const filteredArtistIndex = computed(() => {
  if (!filteredArtists.value) return artistIndex.value
  const ids = new Set(filteredArtists.value.map(a => a.id))
  return artistIndex.value
    .map(g => ({ ...g, artist: g.artist.filter(a => ids.has(a.id)) }))
    .filter(g => g.artist.length)
})

const currentAlbum  = ref(null)
const albumTracks   = ref([])
const editingGenre  = ref(false)
const genreInput    = ref('')
const genreInputEl  = ref(null)
const genreVersion  = ref(0)
const showGenreList = ref(false)

const filteredGenreList = computed(() => {
  const q = genreInput.value.toLowerCase().trim()
  return q ? GENRES.filter(g => g.toLowerCase().includes(q)) : GENRES
})

const GENRE_PREFIX = 'attic_genre_'

const albumGenre = computed(() => {
  genreVersion.value
  try { return localStorage.getItem(GENRE_PREFIX + currentAlbum.value?.id) || currentAlbum.value?.genre || null } catch { return currentAlbum.value?.genre || null }
})

async function startEditGenre() {
  genreInput.value = albumGenre.value || ''
  editingGenre.value = true
  showGenreList.value = true
  await nextTick()
  genreInputEl.value?.focus()
}

function selectGenre(g) {
  genreInput.value = g
  showGenreList.value = false
  saveGenre()
}

function saveGenre() {
  if (!currentAlbum.value) return
  const val = genreInput.value.trim()
  if (val) localStorage.setItem(GENRE_PREFIX + currentAlbum.value.id, val)
  else     localStorage.removeItem(GENRE_PREFIX + currentAlbum.value.id)
  editingGenre.value = false
  showGenreList.value = false
  genreVersion.value++
}

watch(() => currentAlbum.value?.id, () => { editingGenre.value = false; showGenreList.value = false })

async function loadArtists() {
  loading.value = true
  try {
    const index = await getArtists()
    for (const group of index) {
      expandedGroups[group.name] = false
      for (const artist of group.artist) artist._letter = group.name
    }
    artistIndex.value = index
    const all = index.flatMap(g => g.artist)
    discoverArtists.value = [...all].sort(() => Math.random() - 0.5).slice(0, 20)
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
  currentArtistLetter.value = artist._letter || getArtistLetter(artist.name)
  artistDetailImageUrl.value = null
  view.value = 'artist'
  loadingArtist.value = true
  try {
    const data = await getArtist(artist.id)
    currentArtistAlbums.value = data.albums
    artistDetailImageUrl.value = `/artist-images/avatar?name=${encodeURIComponent(artist.name)}`
  } finally { loadingArtist.value = false }
}

function onArtistDetailImgError() {
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
    currentArtistLetter.value  = getArtistLetter(data.info?.name)
    currentArtistAlbums.value  = data.albums
    artistDetailImageUrl.value = `/artist-images/avatar?name=${encodeURIComponent(data.info?.name)}`
  } finally { loadingArtist.value = false }
}

watch(() => route.params.id, (id) => {
  if (id) openArtistById(id)
  else view.value = 'grid'
})

onMounted(async () => {
  getArtistGenreMap().then(m => { artistGenreMap.value = m })
  const [, newestAlbums] = await Promise.all([loadArtists(), getNewestAlbums(100)])
  const seen = new Set()
  const artists = []
  for (const album of newestAlbums) {
    if (album.artistId && !seen.has(album.artistId)) {
      seen.add(album.artistId)
      artists.push({ id: album.artistId, name: album.albumArtist || album.artist })
      if (artists.length >= 20) break
    }
  }
  recentArtists.value = artists
  if (route.params.id) openArtistById(route.params.id)
})
</script>
