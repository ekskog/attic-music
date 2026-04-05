<template>
  <Transition name="slide-up">
    <div v-if="show" class="fixed inset-0 bg-stone-50 z-50 flex flex-col md:hidden">

      <!-- HEADER -->
      <div class="flex items-center justify-between px-6 pt-12 pb-4 flex-shrink-0">
        <button class="text-stone-400 active:text-amber-700 p-2 -ml-2" @click="emit('collapse')">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>
        <span class="text-xs font-medium uppercase tracking-widest text-stone-400">Now Playing</span>
        <button class="text-stone-400 active:text-amber-700 p-2 -mr-2" @click="showQueue = !showQueue">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h10" />
          </svg>
        </button>
      </div>

      <!-- MAIN CONTENT -->
      <div class="flex-1 flex flex-col px-8 overflow-hidden">

        <!-- QUEUE VIEW -->
        <div v-if="showQueue" class="flex-1 overflow-y-auto">
          <div class="text-xs font-medium uppercase tracking-widest text-stone-400 mb-4">Queue ({{ player.queue.length }})</div>
          <div
            v-for="(track, i) in player.queue" :key="i"
            class="flex items-center gap-3 py-3 border-b border-stone-100 last:border-0 cursor-pointer active:bg-stone-100 rounded px-2 -mx-2"
            :class="{ 'text-amber-700': i === player.currentIndex }"
            @click="player.jumpToQueue(i)"
          >
            <span class="text-xs text-stone-300 w-5 text-right flex-shrink-0">{{ i + 1 }}</span>
            <div class="flex-1 overflow-hidden">
              <div class="text-sm font-medium truncate">{{ track.title }}</div>
              <div class="text-xs text-stone-400 truncate">{{ track.artist }}</div>
            </div>
            <span class="text-xs text-stone-300 flex-shrink-0">{{ player.fmt(track.duration) }}</span>
          </div>
        </div>

        <!-- PLAYER VIEW -->
        <template v-else>

          <!-- ALBUM ART -->
          <div class="flex-1 flex items-center justify-center py-4">
            <div class="w-full max-w-xs aspect-square bg-amber-50 overflow-hidden shadow-2xl rounded-sm">
              <img
                v-if="player.currentTrack?.coverArt"
                :src="coverUrl(player.currentTrack.coverArt, 600)"
                class="w-full h-full object-cover"
                @error="onImgError"
              />
              <div v-else class="w-full h-full flex items-center justify-center text-8xl">🎵</div>
            </div>
          </div>

          <!-- TRACK INFO -->
          <div class="flex-shrink-0 mb-6">
            <div class="text-xl font-semibold truncate">{{ player.currentTrack?.title }}</div>
            <div class="text-stone-400 truncate mt-0.5">{{ player.currentTrack?.artist }}</div>
          </div>

          <!-- PROGRESS -->
          <div class="flex-shrink-0 mb-6">
            <div
              ref="progressEl"
              class="w-full h-1 bg-stone-200 rounded-full cursor-pointer mb-2"
              @click="handleSeek"
            >
              <div class="h-full bg-stone-900 rounded-full" :style="{ width: player.progressPct + '%' }"></div>
            </div>
            <div class="flex justify-between text-xs text-stone-400 tabular-nums">
              <span>{{ player.fmt(player.currentTime) }}</span>
              <span>{{ player.fmt(player.duration) }}</span>
            </div>
          </div>

          <!-- CONTROLS -->
          <div class="flex-shrink-0 mb-12">
            <div class="flex items-center justify-between mb-6">
              <button
                class="text-2xl p-2 transition-colors"
                :class="player.shuffle ? 'text-amber-700' : 'text-stone-300'"
                @click="player.shuffle = !player.shuffle"
              >⇄</button>
              <button class="text-4xl p-2 text-stone-700 active:text-amber-700" @click="player.prevTrack()">⏮</button>
              <button
                class="w-16 h-16 rounded-full bg-stone-900 text-white flex items-center justify-center text-2xl active:bg-amber-700 transition-colors shadow-lg"
                @click="player.togglePlay()"
              >
                {{ player.isPlaying ? '⏸' : '▶' }}
              </button>
              <button class="text-4xl p-2 text-stone-700 active:text-amber-700" @click="player.nextTrack()">⏭</button>
              <button
                class="text-2xl p-2 transition-colors"
                :class="player.repeat ? 'text-amber-700' : 'text-stone-300'"
                @click="player.repeat = !player.repeat"
              >↻</button>
            </div>
          </div>

        </template>
      </div>

    </div>
  </Transition>
</template>

<script setup>
import { ref } from 'vue'
import { usePlayerStore } from '../stores/player'
import { coverUrl } from '../api/subsonic'

const props = defineProps({
  show: { type: Boolean, required: true }
})

const emit = defineEmits(['collapse'])

const player     = usePlayerStore()
const progressEl = ref(null)
const showQueue  = ref(false)

function handleSeek(event) {
  player.seek(event, progressEl.value)
}

function onImgError(e) {
  try { if (e?.target) e.target.style.display = 'none' } catch(_) {}
}
</script>

<style scoped>
.slide-up-enter-active, .slide-up-leave-active { transition: transform 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { transform: translateY(100%); }
</style>