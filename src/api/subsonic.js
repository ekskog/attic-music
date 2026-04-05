import { useConfigStore } from '../stores/config'

function hexEncode(s) {
  return Array.from(s).map(c => c.charCodeAt(0).toString(16).padStart(2,'0')).join('')
}

function buildUrl(endpoint, params = {}) {
  const config = useConfigStore()
  const base   = config.server.replace(/\/$/, '')
  const auth   = `u=${encodeURIComponent(config.username)}&p=enc:${hexEncode(config.password)}&v=1.16.1&c=atticweb&f=json`
  const extra  = Object.entries(params).map(([k,v]) => `${k}=${encodeURIComponent(v)}`).join('&')
  return `${base}/rest/${endpoint}?${auth}${extra ? '&'+extra : ''}`
}

export function coverUrl(id, size = 200) {
  return buildUrl('getCoverArt', { id, size })
}

export function streamUrl(id) {
  return buildUrl('stream', { id, format: 'raw' })
}

function ensureArray(val) {
  if (!val) return []
  return Array.isArray(val) ? val : [val]
}

async function request(endpoint, params = {}) {
  const res  = await fetch(buildUrl(endpoint, params))
  const json = await res.json()
  const data = json['subsonic-response']
  if (data.status !== 'ok') throw new Error(data.error?.message || 'API error')
  return data
}

export async function ping() {
  return request('ping')
}

export async function getIndexes() {
  const data    = await request('getIndexes')
  const indexes = ensureArray(data.indexes?.index)
  const items   = []
  for (const idx of indexes) {
    for (const a of ensureArray(idx.artist)) {
      items.push({ id: a.id, name: a.name, isDir: true })
    }
  }
  return items
}

export async function getDirectory(id) {
  const data = await request('getMusicDirectory', { id })
  return ensureArray(data.directory?.child)
}

export async function getArtists() {
  const data = await request('getArtists')
  return ensureArray(data.artists?.index).map(g => ({
    name:   g.name,
    artist: ensureArray(g.artist),
  }))
}

export async function getArtist(id) {
  const data = await request('getArtist', { id })
  return {
    info:   data.artist,
    albums: ensureArray(data.artist?.album),
  }
}

export async function getAlbum(id) {
  const data = await request('getAlbum', { id })
  return {
    info:   data.album,
    tracks: ensureArray(data.album?.song),
  }
}

export async function getAlbumList(size = 500) {
  const data = await request('getAlbumList2', { type: 'alphabeticalByName', size })
  return ensureArray(data.albumList2?.album)
}