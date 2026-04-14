<template>
  <aside class="w-56 bg-white border-r border-stone-200 flex flex-col h-screen flex-shrink-0">

    <!-- HEADER -->
    <div class="px-6 py-7 border-b border-stone-200 flex-shrink-0">
      <div class="font-serif text-2xl font-semibold">attic player</div>
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
    <RecentPlays v-if="config.lastfmUser && config.lastfmKey" />

    <!-- FOOTER -->
    <div class="px-6 py-4 border-t border-stone-200 flex-shrink-0">
      <button class="text-xs text-stone-400 hover:text-amber-700 transition-colors" @click="logout">
        Sign out
      </button>
    </div>

  </aside>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useConfigStore } from '../stores/config'
import RecentPlays from './RecentPlays.vue'

const config = useConfigStore()
const route  = useRoute()
const router = useRouter()

const navItems = [
  { to: '/artists',   icon: '🎤', label: 'Artists'   },
  { to: '/albums',    icon: '💿', label: 'Albums'    },
  { to: '/playlists', icon: '📋', label: 'Playlists' },
]

function isActive(path) {
  return route.path.startsWith(path)
}

function logout() {
  config.logout()
  router.push('/login')
}

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