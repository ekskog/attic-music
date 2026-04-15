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
import { ref, watchEffect } from 'vue'
import { coverUrl } from '../api/subsonic'

const props = defineProps({ artist: Object })
const emit  = defineEmits(['click'])

const imageUrl = ref(null)

// Watch so imageUrl stays in sync if artist.coverArt is enriched after first render.
watchEffect(() => {
  imageUrl.value = coverUrl(props.artist.coverArt || props.artist.id, 200)
})

function onImgError() {
  console.warn(`[artist image] no cover for "${props.artist.name}" (id: ${props.artist.id}, coverArt: ${props.artist.coverArt ?? 'none'})`)
  imageUrl.value = null
}
</script>
