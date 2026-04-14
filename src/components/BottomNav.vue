<template>
  <nav class="fixed bottom-0 left-0 right-0 h-16 bg-white border-t border-stone-200 flex items-center justify-around z-50 md:hidden">
    <RouterLink
      v-for="item in navItems" :key="item.to"
      :to="item.to"
      class="flex flex-col items-center gap-0.5 px-6 py-2 text-stone-400 transition-colors"
      :class="{ 'text-amber-700': isActive(item.to) }"
    >
      <span class="text-2xl leading-none">{{ item.icon }}</span>
      <span class="text-xs font-medium">{{ item.label }}</span>
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
]

function isActive(path) {
  return route.path.startsWith(path)
}
</script>