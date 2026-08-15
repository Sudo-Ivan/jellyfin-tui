package jellyfin

import (
	"net/url"
	"strings"
)

// PlayTarget is a direct HTTP stream mpv can open.
type PlayTarget struct {
	URL           string
	ItemID        string
	MediaSourceID string
	PlaySessionID string
	Transcoding   bool
}

type playbackInfo struct {
	PlaySessionID string        `json:"PlaySessionId"`
	MediaSources  []mediaSource `json:"MediaSources"`
}

type mediaSource struct {
	ID                   string `json:"Id"`
	Container            string `json:"Container"`
	DirectStreamURL      string `json:"DirectStreamUrl"`
	TranscodingURL       string `json:"TranscodingUrl"`
	SupportsDirectPlay   bool   `json:"SupportsDirectPlay"`
	SupportsDirectStream bool   `json:"SupportsDirectStream"`
	SupportsTranscoding  bool   `json:"SupportsTranscoding"`
	ETag                 string `json:"ETag"`
}

// OpenStream resolves a static direct-play URL for itemID.
func (c *Client) OpenStream(itemID string) PlayTarget {
	return c.openPlayback(itemID, false)
}

// OpenTranscodeStream resolves a transcoded stream URL for itemID.
func (c *Client) OpenTranscodeStream(itemID string) PlayTarget {
	t := c.openPlayback(itemID, true)
	t.Transcoding = true
	return t
}

func (c *Client) openPlayback(itemID string, transcode bool) PlayTarget {
	info, err := c.playbackInfo(itemID, transcode)
	if err != nil || len(info.MediaSources) == 0 {
		if transcode {
			return PlayTarget{ItemID: itemID, MediaSourceID: itemID, Transcoding: true, URL: c.buildTranscode(itemID, itemID, "")}
		}
		return c.fallbackTarget(itemID)
	}
	src := pickSource(info.MediaSources, transcode)
	t := PlayTarget{
		ItemID:        itemID,
		MediaSourceID: src.ID,
		PlaySessionID: info.PlaySessionID,
		Transcoding:   transcode,
	}
	if t.MediaSourceID == "" {
		t.MediaSourceID = itemID
	}
	if transcode {
		if src.TranscodingURL != "" {
			t.URL = c.withToken(c.absURL(src.TranscodingURL))
			return t
		}
		t.URL = c.buildTranscode(itemID, t.MediaSourceID, t.PlaySessionID)
		return t
	}
	if src.DirectStreamURL != "" {
		t.URL = c.withToken(c.absURL(src.DirectStreamURL))
		return t
	}
	t.URL = c.buildStream(itemID, t.MediaSourceID, src.Container, t.PlaySessionID, src.ETag)
	return t
}

func (c *Client) playbackInfo(itemID string, transcode bool) (playbackInfo, error) {
	q := url.Values{qUserID: {c.UserID}}
	body := map[string]any{
		"UserId":              c.UserID,
		"IsPlayback":          true,
		"AutoOpenLiveStream":  false,
		"EnableDirectPlay":    !transcode,
		"EnableDirectStream":  !transcode,
		"EnableTranscoding":   transcode,
		"MaxStreamingBitrate": maxBitrate,
	}
	var info playbackInfo
	err := c.post(pathItems+"/"+itemID+pathPlaybackInfo, q, body, &info)
	return info, err
}

func pickSource(list []mediaSource, transcode bool) mediaSource {
	if transcode {
		for _, s := range list {
			if s.SupportsTranscoding || s.TranscodingURL != "" {
				return s
			}
		}
		return list[0]
	}
	for _, s := range list {
		if s.SupportsDirectPlay || s.SupportsDirectStream {
			return s
		}
	}
	return list[0]
}

func (c *Client) fallbackTarget(itemID string) PlayTarget {
	return PlayTarget{
		URL:           c.StreamURL(itemID),
		ItemID:        itemID,
		MediaSourceID: itemID,
	}
}

func (c *Client) absURL(raw string) string {
	if strings.Contains(raw, "://") {
		return raw
	}
	if !strings.HasPrefix(raw, "/") {
		raw = "/" + raw
	}
	return c.Base + raw
}

func (c *Client) withToken(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	q := u.Query()
	if q.Get(qAPIKey) == "" && c.Token != "" {
		q.Set(qAPIKey, c.Token)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func firstContainer(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, ','); i >= 0 {
		s = s[:i]
	}
	return strings.ToLower(strings.TrimSpace(s))
}

func (c *Client) buildStream(itemID, sourceID, container, playSess, tag string) string {
	path := pathVideos + itemID + pathStream
	if ext := firstContainer(container); ext != "" {
		path += "." + ext
	}
	q := url.Values{
		qStatic:        {boolTrue},
		qMediaSourceID: {sourceID},
		qAPIKey:        {c.Token},
		qDeviceID:      {c.DeviceID},
	}
	if playSess != "" {
		q.Set(qPlaySession, playSess)
	}
	if tag != "" {
		q.Set(qTag, tag)
	}
	return c.Base + path + "?" + q.Encode()
}

func (c *Client) buildTranscode(itemID, sourceID, playSess string) string {
	path := pathVideos + itemID + pathMaster
	q := url.Values{
		qMediaSourceID: {sourceID},
		qAPIKey:        {c.Token},
		qDeviceID:      {c.DeviceID},
		qVideoCodec:    {codecH264},
		qAudioCodec:    {codecAAC},
	}
	if playSess != "" {
		q.Set(qPlaySession, playSess)
	}
	return c.Base + path + "?" + q.Encode()
}

// StreamHeaders are extra HTTP headers mpv should send with the stream.
func (c *Client) StreamHeaders() []string {
	ua := clientName + "/" + clientVersion
	h := []string{"User-Agent: " + ua}
	if c.Token != "" {
		h = append(h, headerToken+": "+c.Token)
		h = append(h, headerAuth+": "+c.authHeader())
	}
	return h
}
