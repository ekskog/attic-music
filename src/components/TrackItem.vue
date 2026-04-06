<template>
  <div
    class="grid gap-2 px-3 py-2 rounded cursor-pointer items-center text-sm transition-colors hover:bg-white group"
    :class="{ 'text-amber-700': isPlaying }"
    style="grid-template-columns: 36px 1fr 160px 56px 32px"
    @click="emit('play')"
  >
    <!-- NUMBER / NOW PLAYING -->
    <div class="text-xs text-stone-400 text-center">
      <span v-if="isPlaying">♪</span>
      <span v-else>{{ track.track || index + 1 }}</span>
    </div>

    <!-- TITLE -->
    <div class="truncate font-medium">{{ track.title }}</div>

    <!-- ARTIST -->
    <div class="truncate text-xs text-stone-400">{{ track.artist }}</div>

    <!-- DURATION -->
    <div class="text-xs text-stone-400">{{ fmt(track.duration) }}</div>

    <!-- ADD TO QUEUE -->
    <button
      class="opacity-0 group-hover:opacity-100 text-stone-400 hover:text-amber-700 transition-all bg-transparent border-none text-lg leading-none justify-self-center"
      @click.stop="emit('queue')"
      title="Add to queue"
    >+</button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { usePlayerStore } from '../stores/player'

const props = defineProps({
  track: { type: Object,  required: true },
  index: { type: Number,  required: true },
})

const emit = defineEmits(['play', 'queue'])

const player    = usePlayerStore()
const isPlaying = computed(() => player.currentTrack?.id === props.track.id)

function fmt(secs) {
  if (!secs && secs !== 0) return '—'
  const s = Math.floor(secs)
  return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, '0')}`
}
</script>