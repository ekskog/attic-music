<template>
  <div
    class="grid gap-2 px-3 py-3 rounded-lg cursor-pointer items-center text-sm transition-colors active:bg-white group"
    :class="{ 'text-amber-700': isPlaying }"
    style="grid-template-columns: 28px 1fr 44px 28px"
    @click="emit('play')"
  >
    <div class="text-xs text-stone-400 text-center">
      <span v-if="isPlaying">♪</span>
      <span v-else>{{ track.track || index + 1 }}</span>
    </div>
    <div class="truncate font-medium text-sm">{{ track.title }}</div>
    <div class="text-xs text-stone-400 text-right">{{ fmt(track.duration) }}</div>

    <div class="relative justify-self-center">
      <button
        class="md:opacity-0 md:group-hover:opacity-100 text-stone-400 hover:text-amber-700 active:text-amber-700 transition-all bg-transparent border-none text-lg leading-none"
        @click.stop="toggleMenu"
        title="Add to playlist"
      >+</button>

      <div
        v-if="menuOpen"
        class="absolute right-0 top-full mt-1 z-50 bg-white border border-stone-200 rounded-lg shadow-lg py-1 w-44"
      >
        <div class="px-3 py-1.5 text-xs text-stone-400 uppercase tracking-wider border-b border-stone-100">Add to playlist</div>
        <div v-if="!plStore.playlists.length" class="px-3 py-2 text-xs text-stone-400 italic">No playlists</div>
        <button
          v-for="pl in plStore.playlists"
          :key="pl.id"
          class="w-full text-left px-3 py-2 text-sm hover:bg-amber-50 hover:text-amber-700 transition-colors"
          @click.stop="addTo(pl)"
        >{{ pl.name }}</button>
        <div class="border-t border-stone-100 mt-1">
          <button
            class="w-full text-left px-3 py-2 text-sm text-stone-500 hover:bg-stone-50 transition-colors"
            @click.stop="addToQueue"
          >Add to queue</button>
          <button
            v-if="removable"
            class="w-full text-left px-3 py-2 text-sm text-red-400 hover:bg-red-50 hover:text-red-600 transition-colors"
            @click.stop="removeFromPlaylist"
          >Remove from playlist</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { usePlayerStore } from '../stores/player'
import { usePlaylistStore } from '../stores/playlist'

const props = defineProps({
  track:    { type: Object,  required: true },
  index:    { type: Number,  required: true },
  removable: { type: Boolean, default: false },
})

const emit = defineEmits(['play', 'queue', 'remove'])

const player  = usePlayerStore()
const plStore = usePlaylistStore()

const isPlaying = computed(() => player.currentTrack?.id === props.track.id)
const menuOpen  = ref(false)

async function toggleMenu() {
  if (!menuOpen.value) {
    await plStore.load()
  }
  menuOpen.value = !menuOpen.value
  if (menuOpen.value) {
    document.addEventListener('click', closeMenu)
  }
}

function closeMenu() {
  menuOpen.value = false
  document.removeEventListener('click', closeMenu)
}

async function addTo(pl) {
  await plStore.addTrack(pl.id, props.track.id)
  closeMenu()
}

function addToQueue() {
  player.addToQueue(props.track)
  closeMenu()
}

function removeFromPlaylist() {
  emit('remove', props.index)
  closeMenu()
}

onUnmounted(() => document.removeEventListener('click', closeMenu))

function fmt(secs) {
  if (!secs && secs !== 0) return '—'
  const s = Math.floor(secs)
  return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, '0')}`
}
</script>
