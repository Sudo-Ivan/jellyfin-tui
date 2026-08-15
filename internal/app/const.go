package app

import "time"

const (
	appTitle = "jellyfin-tui"

	schemeHTTP  = "http://"
	schemeHTTPS = "https://"

	loginTries    = 3
	reconnectWait = 400 * time.Millisecond
	uiQueueSize   = 32
	frameTick     = 250 * time.Millisecond
	progressEvery = 10 * time.Second

	seekShort = 10 * time.Second
	seekLong  = 60 * time.Second
	volStep   = 5
	pageJump  = 10

	railWidth     = 22
	seasonWidth   = 24
	loginBoxW     = 56
	loginBoxH     = 12
	fieldLabelW   = 12
	helpBoxW      = 64
	helpBoxH      = 22
	headerHostX   = 16
	loginFieldGap = 2

	homeResumeN = 16
	homeNextUpN = 16
	homeLatestN = 20
	searchLimit = 80
	newlyLimit  = 80

	fnvOffset = 2166136261
	fnvPrime  = 16777619

	secondsPerHour = 3600
	secondsPerMin  = 60
	printableMin   = 32
	loginFields    = 3

	idHome        = "home"
	nameHome      = "Home"
	nameAdded     = "newly added"
	idAdded       = "added"
	idDownloads   = "downloads"
	nameDownloads = "downloads"

	statusLoading   = "loading"
	statusPlaying   = "playing"
	statusSigningIn = "signing in"
	statusSearching = "searching"
	statusReconnect = "reconnecting"
	statusRandom    = "random movie"
	statusNoNext    = "no next episode"
	errExpired      = "session expired"
	errNeedServer   = "need a server url"
	errNeedUser     = "need a username"

	keysMain   = "enter open   / search   m random   a added   t tracks   c cast   d download   space pause   n next   ? help   q quit"
	keysSearch = "genre:action  actor:name  year:1999  year:2010-2012"
	keysLogin  = "tab field   enter connect   r reconnect   q quit"

	markPlay     = ">"
	markPause    = "#"
	tickDone     = "* "
	tickMid      = ": "
	tracksBoxW   = 64
	tracksBoxH   = 18
	castBoxW     = 56
	castBoxH     = 20
	fallbackWait = 3 * time.Second

	statusTranscode = "transcoding fallback"
	statusOffline   = "offline mode"
	statusDownload  = "downloading"
	swatchCh        = '▍'
	caretCh         = '▌'
	secretCh        = '*'
)
