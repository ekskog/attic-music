const cache = new Map()

export async function fetchLyrics(track) {
  const key = track.id
  if (cache.has(key)) return cache.get(key)

  const params = new URLSearchParams({
    artist_name: track.artist || '',
    track_name:  track.title  || '',
    album_name:  track.album  || '',
    duration:    String(Math.round(track.duration || 0)),
  })
  let result = null
  try {
    const res = await fetch(`https://lrclib.net/api/get?${params}`)
    if (res.ok) {
      const data = await res.json()
      result = {
        plain:  data.plainLyrics  || null,
        synced: parseLrc(data.syncedLyrics || ''),
      }
    }
  } catch (_) {}
  cache.set(key, result)
  return result
}

function parseLrc(lrc) {
  if (!lrc) return null
  const lines = lrc.split('\n')
    .map(line => {
      const m = line.match(/^\[(\d+):(\d+\.\d+)\](.*)/)
      if (!m) return null
      return { time: parseInt(m[1]) * 60 + parseFloat(m[2]), text: m[3].trim() }
    })
    .filter(l => l && l.text)
  return lines.length ? lines : null
}
