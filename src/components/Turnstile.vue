<template>
  <div :id="containerId"></div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'

const props = defineProps({
  sitekey: { type: String, required: true }
})

const emit = defineEmits(['verified', 'expired', 'error'])
const containerId = `turnstile-${Math.random().toString(36).slice(2)}`
let widgetId = null

onMounted(() => {
  const render = () => {
    if (!window.turnstile) { setTimeout(render, 50); return }
    widgetId = window.turnstile.render(`#${containerId}`, {
      sitekey:            props.sitekey,
      callback:           (token) => emit('verified', token),
      'expired-callback': ()      => emit('expired'),
      'error-callback':   ()      => emit('error'),
    })
  }
  render()
})

onUnmounted(() => {
  if (window.turnstile && widgetId !== null) {
    window.turnstile.remove(widgetId)
    widgetId = null
  }
})
</script>