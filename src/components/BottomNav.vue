<template>
  <nav class="fixed bottom-0 left-0 right-0 h-14 bg-white border-t border-stone-200 flex items-center justify-around z-50 md:hidden">
    <RouterLink
      v-for="item in navItems" :key="item.to"
      :to="item.to"
      class="flex-1 flex items-center justify-center py-2 text-xs font-medium uppercase tracking-widest transition-colors text-stone-400"
      :class="{ 'text-amber-700': isActive(item.to) }"
    >
      {{ item.label }}
    </RouterLink>
  </nav>
</template>

<script setup>
import { useRoute } from 'vue-router'

const route = useRoute()

const navItems = [
  { to: '/artists',   icon: '🎤', label: 'Artists'   },
  { to: '/albums',    icon: '💿', label: 'Albums'    },
  { to: '/playlists', icon: '📋', label: 'Playlists' },
  { to: '/search',    icon: '🔍', label: 'Search'    },
]

function isActive(path) {
  return route.path.startsWith(path)
}
</script>