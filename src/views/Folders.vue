<template>
  <div class="flex flex-col h-full overflow-hidden">

    <!-- HEADER -->
    <div class="px-8 py-7 border-b border-stone-200 bg-white flex-shrink-0">
      <div class="flex items-center gap-1.5 text-xs text-stone-400 mb-2 flex-wrap">
        <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="goToRoot">Library</span>
        <template v-for="(crumb, i) in breadcrumbs" :key="crumb.id">
          <span class="opacity-40">›</span>
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="goToBreadcrumb(i)">{{ crumb.name }}</span>
        </template>
      </div>
      <h1 class="font-serif text-4xl font-semibold">
        {{ breadcrumbs.length ? breadcrumbs[breadcrumbs.length - 1].name : 'Library' }}
      </h1>
    </div>

    <!-- CONTENT -->
    <div class="flex-1 overflow-y-auto px-8 py-6 pb-40 md:pb-40 md:pb-24">
      <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
      <div v-else-if="!items.length" class="flex flex-col items-center justify-center py-24 text-stone-400 gap-2">
        <span class="text-4xl">📂</span>
        <span class="font-serif text-lg">Empty folder</span>
      </div>
      <div v-else class="flex flex-col gap-px">
        <div
          v-for="item in items" :key="item.id"
          class="flex items-center gap-3 px-3 py-2.5 rounded cursor-pointer hover:bg-white transition-colors group text-sm"
          :class="{ 'text-amber-700 font-medium': player.currentTrack?.id === item.id }"
          @click="handleItem(item)"
        >
          <span class="text-stone-400 flex-shrink-0">{{ item.isDir ? '📁' : '🎵' }}</span>
          <span class="flex-1 truncate">{{ item.title || item.name }}</span>
          <span v-if="!item.isDir" class="text-xs text-stone-400 flex-shrink-0">{{ player.fmt(item.duration) }}</span>
          <button
            v-if="!item.isDir"
            class="opacity-0 group-hover:opacity-100 text-stone-400 hover:text-amber-700 border border-stone-200 hover:border-amber-700 hover:bg-amber-50 px-2 py-0.5 rounded text-base leading-snug transition-all flex-shrink-0"
            @click.stop="player.addToQueue(item)"
            title="Add to queue"
          >+</button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePlayerStore } from '../stores/player'
import { getIndexes, getDirectory } from '../api/subsonic'

const route  = useRoute()
const router = useRouter()
const player = usePlayerStore()

const items       = ref([])
const breadcrumbs = ref([])
const loading     = ref(false)

async function loadRoot() {
  loading.value = true
  try {
    items.value       = await getIndexes()
    breadcrumbs.value = []
  } finally {
    loading.value = false
  }
}

async function loadDir(id) {
  loading.value = true
  try {
    items.value = await getDirectory(id)
  } finally {
    loading.value = false
  }
}

async function handleItem(item) {
  if (item.isDir) {
    breadcrumbs.value.push({ id: item.id, name: item.title || item.name })
    router.push({ name: 'folders', params: { id: item.id } })
  } else {
    const songs = items.value.filter(i => !i.isDir)
    const idx   = songs.findIndex(s => s.id === item.id)
    player.playTrack(item, songs, idx >= 0 ? idx : 0)
  }
}

async function goToRoot() {
  breadcrumbs.value = []
  router.push({ name: 'folders' })
}

async function goToBreadcrumb(index) {
  const crumb = breadcrumbs.value[index]
  breadcrumbs.value = breadcrumbs.value.slice(0, index + 1)
  router.push({ name: 'folders', params: { id: crumb.id } })
}

onMounted(() => {
  if (route.params.id) loadDir(route.params.id)
  else loadRoot()
})

watch(() => route.params.id, (id) => {
  if (id) loadDir(id)
  else loadRoot()
})
</script>