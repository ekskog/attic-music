import { ref, computed } from 'vue'
import { usePlayerStore } from '../stores/player'
import { fetchLyrics } from '../api/lrclib'

// Module-level singleton — shared across Player.vue and FullPlayer.vue
const lyrics  = ref(null)   // null | { instrumental, synced: [...] | null, plain: string | null }
const loading = ref(false)
const error   = ref(false)
let   lastFetchedId = null

export function useLyrics() {
  const player = usePlayerStore()

  async function loadLyrics() {
    const track = player.currentTrack
    if (!track) return
    if (track.id === lastFetchedId) return   // already fetched for this track
    lastFetchedId = track.id
    lyrics.value  = null
    error.value   = false
    loading.value = true
    try {
      lyrics.value = await fetchLyrics({
        artist:   track.artist,
        title:    track.title,
        album:    track.album,
        duration: track.duration,
      })
    } catch {
      error.value = true
    } finally {
      loading.value = false
    }
  }

  // Index of the currently-playing line (for synced lyrics)
  const currentLineIndex = computed(() => {
    if (!lyrics.value?.synced) return -1
    const t = player.currentTime
    let idx = -1
    for (let i = 0; i < lyrics.value.synced.length; i++) {
      if (lyrics.value.synced[i].time <= t) idx = i
      else break
    }
    return idx
  })

  return { lyrics, loading, error, currentLineIndex, loadLyrics }
}
