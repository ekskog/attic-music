<template>
  <div class="cursor-pointer group" @click="emit('click')">
    <div class="aspect-square bg-stone-100 rounded-xl overflow-hidden mb-2 transition-transform duration-200 group-hover:scale-[1.03]">
      <img
        v-if="imageUrl"
        :src="imageUrl"
        :alt="artist.name"
        loading="lazy"
        class="w-full h-full object-cover"
        @error="onImgError"
      />
      <div v-else class="w-full h-full flex items-center justify-center font-serif text-3xl font-semibold text-stone-300 select-none">
        {{ artist.name[0]?.toUpperCase() }}
      </div>
    </div>
    <div class="text-sm font-medium truncate leading-tight">{{ artist.name }}</div>
    <div v-if="artist.albumCount" class="text-xs text-stone-400 mt-0.5">
      {{ artist.albumCount }} album{{ artist.albumCount !== 1 ? 's' : '' }}
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { coverUrl } from '../api/subsonic'
import { getArtistImage } from '../api/deezer'

const props = defineProps({ artist: Object })
const emit  = defineEmits(['click'])

// Try local Gonic cover art first (getCoverArt?id=ar-xxx)
const imageUrl   = ref(coverUrl(props.artist.id, 200))
const triedLocal = ref(false)

async function onImgError() {
  if (!triedLocal.value) {
    triedLocal.value = true
    imageUrl.value = null                              // show letter while fetching
    const url = await getArtistImage(props.artist.name)
    imageUrl.value = url                               // null → stays as letter
  } else {
    imageUrl.value = null
  }
}
</script>
