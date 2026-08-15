package jellyfin

import "time"

// StreamURL is a static direct-play URL when PlaybackInfo is unavailable.
func (c *Client) StreamURL(itemID string) string {
	return c.buildStream(itemID, itemID, "", "", "")
}

// Playing reports the start of playback to the server.
func (c *Client) Playing(t PlayTarget, pos, duration time.Duration, paused bool) error {
	body := playBody(t, pos, duration, paused)
	body["IsMuted"] = false
	body["PlaybackStartTimeTicks"] = ticks(0)
	return c.post(pathPlaying, nil, body, nil)
}

// Progress reports a time-update while playing.
func (c *Client) Progress(t PlayTarget, pos, duration time.Duration, paused bool) error {
	body := playBody(t, pos, duration, paused)
	body["EventName"] = eventTime
	return c.post(pathProgress, nil, body, nil)
}

// Stopped reports the end of a playback session.
func (c *Client) Stopped(t PlayTarget, pos, duration time.Duration) error {
	body := map[string]any{
		"ItemId":        t.ItemID,
		"PositionTicks": ticks(pos),
		"MediaSourceId": sourceID(t),
	}
	if t.PlaySessionID != "" {
		body[qPlaySession] = t.PlaySessionID
	}
	if duration > 0 {
		body["RunTimeTicks"] = ticks(duration)
	}
	return c.post(pathStopped, nil, body, nil)
}

func playBody(t PlayTarget, pos, duration time.Duration, paused bool) map[string]any {
	method := playDirect
	if t.Transcoding {
		method = playTranscode
	}
	body := map[string]any{
		"ItemId":        t.ItemID,
		"CanSeek":       true,
		"IsPaused":      paused,
		"PlayMethod":    method,
		"PositionTicks": ticks(pos),
		"MediaSourceId": sourceID(t),
	}
	if t.PlaySessionID != "" {
		body[qPlaySession] = t.PlaySessionID
	}
	if duration > 0 {
		body["RunTimeTicks"] = ticks(duration)
	}
	return body
}

func sourceID(t PlayTarget) string {
	if t.MediaSourceID != "" {
		return t.MediaSourceID
	}
	return t.ItemID
}
