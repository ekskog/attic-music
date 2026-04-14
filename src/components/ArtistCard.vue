<template>
  <div ref="cardEl" class="cursor-pointer group" @click="emit('click')">
    <div class="aspect-square bg-stone-100 rounded-xl overflow-hidden mb-2 transition-transform duration-200 group-hover:scale-[1.03]">
      <img
        v-if="imageUrl"
        :src="imageUrl"
        :alt="artist.name"
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
import { ref, onMounted, onUnmounted } from 'vue'
import { coverUrl } from '../api/subsonic'
import { getArtistImage } from '../api/deezer'

const props = defineProps({ artist: Object })
const emit  = defineEmits(['click'])

const cardEl  = ref(null)
const imageUrl = ref(
  props.artist.coverArt ? coverUrl(props.artist.coverArt, 200) : null
)

let observer = null

function onImgError() {
  imageUrl.value = null
}

onMounted(() => {
  observer = new IntersectionObserver(async ([entry]) => {
    if (!entry.isIntersecting) return
    observer.disconnect()
    const url = await getArtistImage(props.artist.name)
    if (url) imageUrl.value = url
  }, { rootMargin: '100px' })

  if (cardEl.value) observer.observe(cardEl.value)
})

onUnmounted(() => {
  if (observer) observer.disconnect()
})
</script>
