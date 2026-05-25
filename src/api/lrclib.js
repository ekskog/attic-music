const BASE = 'https://lrclib.net/api'

function parseLrc(lrc) {
  return lrc.split('\n')
    .map(line => {
      const m = line.match(/^\[(\d+):(\d+(?:\.\d+)?)\](.*)/)
      if (!m) return null
      return { time: parseInt(m[1]) * 60 + parseFloat(m[2]), text: m[3].trim() }
    })
    .filter(l => l !== null)
}

export async function fetchLyrics({ artist, title, album, duration }) {
  const params = new URLSearchParams({ artist_name: artist, track_name: title })
  if (album)    params.set('album_name', album)
  if (duration) params.set('duration', String(Math.round(duration)))
  try {
    const res = await fetch(`${BASE}/get?${params}`)
    if (res.status === 404) return null
    if (!res.ok) return null
    const data = await res.json()
    if (data.instrumental) return { instrumental: true, synced: null, plain: null }
    if (!data.syncedLyrics && !data.plainLyrics) return null
    return {
      instrumental: false,
      synced: data.syncedLyrics ? parseLrc(data.syncedLyrics) : null,
      plain:  data.plainLyrics || null,
    }
  } catch {
    return null
  }
}
