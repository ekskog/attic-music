const express = require('express');
const path = require('path');
const crypto = require('crypto');
const fetch = global.fetch || require('node-fetch');

const app = express();
app.use(express.json());

// Serve static frontend
app.use(express.static(path.join(__dirname, '..')));

function apiSig(params, secret) {
  const keys = Object.keys(params).sort();
  const str = keys.map(k => `${k}${params[k]}`).join('') + secret;
  return crypto.createHash('md5').update(str, 'utf8').digest('hex');
}

app.post('/api/scrobble', async (req, res) => {
  try {
    const { artist, track, album, timestamp } = req.body || {};
    const apiKey = process.env.LASTFM_API_KEY;
    const apiSecret = process.env.LASTFM_API_SECRET;
    // sessionKey may be provided per-request (user session) or via env fallback
    const sessionKey = req.body.session_key || process.env.LASTFM_SESSION_KEY;
    if (!apiKey || !apiSecret || !sessionKey) return res.status(400).json({ error: 'Last.fm credentials not configured' });
    if (!artist || !track) return res.status(400).json({ error: 'artist and track required' });

    const method = 'track.scrobble';
    const params = {
      artist: artist,
      track: track,
      api_key: apiKey,
      sk: sessionKey,
      method: method,
      timestamp: timestamp || Math.floor(Date.now()/1000),
    };

    const api_sig = apiSig(params, apiSecret);
    const body = new URLSearchParams({ ...params, api_sig, format: 'json' });

    const resp = await fetch('https://ws.audioscrobbler.com/2.0/', {
      method: 'POST',
      body,
    });
    const json = await resp.json();
    return res.json(json);
  } catch (err) {
    console.error(err);
    return res.status(500).json({ error: err.message });
  }
});

// --- Last.fm auth helpers ---
app.get('/api/auth/token', async (req, res) => {
  try {
    const apiKey = process.env.LASTFM_API_KEY;
    if (!apiKey) return res.status(400).json({ error: 'Missing LASTFM_API_KEY' });
    const url = `https://ws.audioscrobbler.com/2.0/?method=auth.gettoken&api_key=${encodeURIComponent(apiKey)}&format=json`;
    const r = await fetch(url);
    const j = await r.json();
    // return token and api_key so frontend can open auth URL
    return res.json({ token: j?.token || (j?.['token']), api_key: apiKey });
  } catch (err) {
    console.error(err);
    return res.status(500).json({ error: err.message });
  }
});

app.post('/api/auth/session', async (req, res) => {
  try {
    const { token } = req.body || {};
    const apiKey = process.env.LASTFM_API_KEY;
    const apiSecret = process.env.LASTFM_API_SECRET;
    if (!apiKey || !apiSecret) return res.status(400).json({ error: 'Missing API credentials' });
    if (!token) return res.status(400).json({ error: 'Missing token' });

    const params = { api_key: apiKey, method: 'auth.getSession', token };
    const keys = Object.keys(params).sort();
    const sig = crypto.createHash('md5').update(keys.map(k => `${k}${params[k]}`).join('') + apiSecret, 'utf8').digest('hex');
    const body = new URLSearchParams({ ...params, api_sig: sig, format: 'json' });
    const r = await fetch('https://ws.audioscrobbler.com/2.0/', { method: 'POST', body });
    const j = await r.json();
    // j.session.key contains the session key
    return res.json(j);
  } catch (err) {
    console.error(err);
    return res.status(500).json({ error: err.message });
  }
});

const port = process.env.PORT || 80;
app.listen(port, () => console.log('server listening on', port));
