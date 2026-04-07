<template>
  <div class="flex flex-col h-full overflow-hidden">

    <!-- HEADER -->
    <div class="px-6 py-6 border-b border-stone-200 bg-white flex-shrink-0">
      <div class="flex items-center gap-1.5 text-xs text-stone-400 mb-1 flex-wrap">
        <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="goToRoot">Library</span>
        <template v-for="(crumb, i) in breadcrumbs" :key="crumb.id">
          <span class="opacity-40">›</span>
          <span class="cursor-pointer hover:text-amber-700 transition-colors" @click="goToBreadcrumb(i)">{{ crumb.name }}</span>
        </template>
      </div>
      <h1 class="font-serif text-3xl font-semibold">
        {{ breadcrumbs.length ? breadcrumbs[breadcrumbs.length - 1].name : 'Library' }}
      </h1>
    </div>

    <!-- CONTENT -->
    <div class="flex-1 overflow-y-auto px-4 py-4 pb-40 md:pb-24">
      <div v-if="loading" class="flex items-center justify-center py-24 text-stone-400 text-sm">Loading…</div>
      <div v-else-if="!items.length" class="flex flex-col items-center justify-center py-24 text-stone-400 gap-2">
        <span class="text-4xl">📂</span>
        <span class="font-serif text-lg">Empty folder</span>
      </div>
      <div v-else>

        <!-- DIRECTORIES AS EXPANDABLE LIST -->
        <div v-if="dirs.length" class="flex flex-col gap-1 mb-6">
          <FolderNode v-for="item in dirs" :key="item.id" :item="item" />
        </div>

        <!-- TRACKS AS LIST -->
        <div v-if="tracks.length" class="flex flex-col gap-px">
          <div
            v-for="item in tracks" :key="item.id"
            class="flex items-center gap-3 px-3 py-3 rounded-lg cursor-pointer active:bg-white transition-colors group text-sm"
            :class="{ 'text-amber-700 font-medium': player.currentTrack?.id === item.id }"
            @click="handleItem(item)"
          >
            <span class="text-stone-400 flex-shrink-0">🎵</span>
            <span class="flex-1 truncate">{{ item.title || item.name }}</span>
            <span class="text-xs text-stone-400 flex-shrink-0">{{ player.fmt(item.duration) }}</span>
            <button
              class="opacity-0 group-hover:opacity-100 text-stone-400 hover:text-amber-700 border border-stone-200 hover:border-amber-700 hover:bg-amber-50 px-2 py-0.5 rounded text-base leading-snug transition-all flex-shrink-0"
              @click.stop="player.addToQueue(item)"
            >+</button>
          </div>
        </div>

      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, watch, onMounted, computed, provide } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePlayerStore } from '../stores/player'
import { getIndexes, getDirectory } from '../api/subsonic'
import FolderNode from '../components/FolderNode.vue'

const route  = useRoute()
const router = useRouter()
const player = usePlayerStore()

const items       = ref([])
const breadcrumbs = ref([])
const loading     = ref(false)

// top-level accordion controller so top-level FolderNode siblings collapse
const topOpenChildId = ref(null)
provide('openChildId', topOpenChildId)
provide('setOpenChild', (id) => { topOpenChildId.value = id })

const dirs   = computed(() => items.value.filter(i => i.isDir))
const tracks = computed(() => items.value.filter(i => !i.isDir))

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
    const songs = tracks.value
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