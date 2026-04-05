<template>
  <aside class="w-56 bg-white border-r border-stone-200 flex flex-col h-screen flex-shrink-0">

    <!-- HEADER -->
    <div class="px-6 py-7 border-b border-stone-200 flex-shrink-0">
      <div class="font-serif text-2xl font-semibold">Attic</div>
      <div class="text-xs text-stone-400 mt-0.5 truncate">{{ config.server }}</div>
    </div>

    <!-- NAV -->
    <nav class="p-4 flex-shrink-0">
      <div class="text-xs font-medium uppercase tracking-widest text-stone-400 px-2 py-2">Library</div>
      <RouterLink
        v-for="item in navItems" :key="item.to"
        :to="item.to"
        class="nav-item"
        :class="{ active: isActive(item.to) }"
      >
        <span class="w-4 text-center">{{ item.icon }}</span>
        {{ item.label }}
      </RouterLink>
    </nav>

    <!-- SCROBBLES -->
    <div v-if="config.lastfmUser && config.lastfmKey" class="flex-1 border-t border-stone-200 p-4 overflow-hidden flex flex-col min-h-0">
      <div class="flex items-center justify-between mb-3 flex-shrink-0">
        <span class="text-xs font-medium uppercase tracking-widest text-stone-400">Recent plays</span>
        <span class="text-xs text-amber-700">last.fm</span>
      </div>
      <div class="overflow-y-auto flex-1 space-y-2">
        <div v-if="!scrobbles.length" class="text-xs text-stone-400">Nothing yet…</div>
        <div v-for="(s, i) in scrobbles" :key="i" class="pb-2 border-b border-stone-100 last:border-0">
          <div class="text-xs font-medium truncate" :class="s.nowplaying ? 'text-amber-700' : 'text-stone-800'">
            <span v-if="s.nowplaying" class="inline-block w-1.5 h-1.5 rounded-full bg-amber-700 mr-1 mb-0.5 animate-pulse"></span>
            {{ s.track }}
          </div>
          <div class="text-xs text-stone-400 truncate">{{ s.artist }}</div>
          <div v-if="!s.nowplaying" class="text-xs text-stone-300 mt-0.5">{{ s.when }}</div>
        </div>
      </div>
    </div>

    <!-- FOOTER -->
    <div class="px-6 py-4 border-t border-stone-200 flex-shrink-0">
      <button class="text-xs text-stone-400 hover:text-amber-700 transition-colors" @click="logout">
        Sign out
      </button>
    </div>

  </aside>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useConfigStore } from '../stores/config'
import { getRecentTracks } from '../api/lastfm'

const config = useConfigStore()
const route  = useRoute()
const router = useRouter()

const scrobbles = ref([])
let scrobbleTimer = null

const navItems = [
  { to: '/folders', icon: '📁', label: 'Folders'  },
  { to: '/artists', icon: '🎤', label: 'Artists'  },
  { to: '/albums',  icon: '💿', label: 'Albums'   },
]

function isActive(path) {
  return route.path.startsWith(path)
}

async function fetchScrobbles() {
  scrobbles.value = await getRecentTracks(5)
}

function logout() {
  config.logout()
  router.push('/login')
}

onMounted(() => {
  if (config.lastfmUser && config.lastfmKey) {
    fetchScrobbles()
    scrobbleTimer = setInterval(fetchScrobbles, 30000)
  }
})

onUnmounted(() => {
  if (scrobbleTimer) clearInterval(scrobbleTimer)
})
</script>

<style scoped>
@reference "../style.css";

.nav-item {
  @apply flex items-center gap-2.5 px-2 py-2 rounded text-sm text-stone-700 transition-all cursor-pointer no-underline;
}
.nav-item:hover {
  @apply bg-stone-50;
}
.nav-item.active {
  @apply bg-amber-50 text-amber-700 font-medium;
}
</style>