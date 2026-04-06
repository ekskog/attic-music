<template>
  <div class="h-screen flex items-center justify-center bg-stone-50">
    <div class="w-96 bg-white border border-stone-200 p-12">
      <h1 class="font-serif text-5xl font-semibold mb-1">Attic</h1>
      <p class="text-stone-400 text-sm mb-8">Connect to your Subsonic server</p>

      <div class="space-y-4">
        <div>
          <label class="block text-xs font-medium uppercase tracking-widest text-stone-400 mb-1.5">Server URL</label>
          <input v-model="form.server" class="input" placeholder="https://gonic.example.com" @keyup.enter="login" />
        </div>
        <div>
          <label class="block text-xs font-medium uppercase tracking-widest text-stone-400 mb-1.5">Username</label>
          <input v-model="form.username" class="input" autocomplete="username" @keyup.enter="login" />
        </div>
        <div>
          <label class="block text-xs font-medium uppercase tracking-widest text-stone-400 mb-1.5">Password</label>
          <input v-model="form.password" type="password" class="input" autocomplete="current-password" @keyup.enter="login" />
        </div>
      </div>

      <!-- TURNSTILE WIDGET -->
      <div class="mt-6" id="turnstile-container"></div>

      <button
        class="mt-4 w-full bg-stone-900 text-white text-sm font-medium tracking-wide py-3 hover:bg-amber-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        @click="login"
        :disabled="logging || !turnstileToken"
      >
        {{ logging ? 'Connecting…' : 'Connect' }}
      </button>

      <p v-if="error" class="mt-3 text-sm text-amber-700">{{ error }}</p>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useConfigStore } from '../stores/config'
import { ping } from '../api/subsonic'

const TURNSTILE_SITE_KEY  = '0x4AAAAAAC1Tg9yFuV4XukU0'
const TURNSTILE_WORKER_URL = 'https://attic-turnstile.ekflare.workers.dev'

const config  = useConfigStore()
const router  = useRouter()
const logging = ref(false)
const error   = ref('')
const turnstileToken = ref('')

const form = reactive({
  server:   '',
  username: '',
  password: '',
})

onMounted(() => {
  form.server   = config.server
  form.username = config.username
  form.password = config.password

  // Wait for turnstile to be available
  const render = () => {
    if (window.turnstile) {
      window.turnstile.render('#turnstile-container', {
        sitekey: TURNSTILE_SITE_KEY,
        callback:          (token) => { turnstileToken.value = token },
        'expired-callback': ()      => { turnstileToken.value = '' },
        'error-callback':   ()      => { turnstileToken.value = '' },
      })
    } else {
      setTimeout(render, 100)
    }
  }
  render()
})

import { reactive, ref, onMounted, onUnmounted } from 'vue'

onUnmounted(() => {
  if (window.turnstile) {
    window.turnstile.remove('#turnstile-container')
  }
})

async function login() {
  error.value = ''
  if (!turnstileToken.value) {
    error.value = 'Please complete the security check.'
    return
  }
  logging.value = true
  try {
    // Verify token with Worker
    const res  = await fetch(TURNSTILE_WORKER_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token: turnstileToken.value })
    })
    const data = await res.json()
    if (!data.success) {
      error.value = 'Security check failed. Please try again.'
      window.turnstile.reset('#turnstile-container')
      turnstileToken.value = ''
      return
    }

    // Proceed with login
    config.server   = form.server.replace(/\/$/, '')
    config.username = form.username
    config.password = form.password
    await ping()
    config.saveSession()
    router.push('/')
  } catch(e) {
    error.value = 'Could not connect — check URL and credentials.'
  } finally {
    logging.value = false
  }
}
</script>

<style scoped>
@reference "../style.css";
.input {
  @apply w-full border border-stone-200 bg-stone-50 px-3 py-2.5 text-sm text-stone-900 outline-none focus:border-amber-700 transition-colors;
}
</style>