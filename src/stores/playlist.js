import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getPlaylists, updatePlaylist } from '../api/subsonic'

export const usePlaylistStore = defineStore('playlist', () => {
  const playlists = ref([])
  let loaded = false

  async function load() {
    if (loaded) return
    playlists.value = await getPlaylists()
    loaded = true
  }

  async function addTrack(playlistId, trackId) {
    await updatePlaylist(playlistId, { songIdsToAdd: [trackId] })
  }

  return { playlists, load, addTrack }
})
