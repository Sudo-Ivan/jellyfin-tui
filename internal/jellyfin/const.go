package jellyfin

import "time"

const (
	clientName    = "jellyfin-tui"
	clientVersion = "0.1.0"

	httpTimeout      = 30 * time.Second
	retryBackoff     = 200 * time.Millisecond
	cacheTTL         = 45 * time.Second
	idleConnTimeout  = 90 * time.Second
	tlsHandshake     = 10 * time.Second
	dialTimeout      = 10 * time.Second
	keepAlive        = 30 * time.Second
	expectContinue   = 1 * time.Second
	idleConns        = 32
	idleConnsPerHost = 8
	maxConnsPerHost  = 16
	maxBitrate       = 200000000
	maxTries         = 3
	maxBodyBytes     = 8 << 20
	maxErrRunes      = 240
	ticksPerSec      = 10_000_000

	httpOKMin     = 200
	httpOKMax     = 300
	httpTooMany   = 429
	httpBadGate   = 502
	httpUnavail   = 503
	httpTimeoutC  = 504
	httpAuthMin   = 401
	httpForbidden = 403
	nsPerTick     = 100

	qLimit         = "Limit"
	qFields        = "Fields"
	qUserID        = "UserId"
	qRecurse       = "Recursive"
	qSortBy        = "SortBy"
	qSortOrd       = "SortOrder"
	qInclude       = "IncludeItemTypes"
	qParent        = "ParentId"
	qSearch        = "SearchTerm"
	qSeason        = "SeasonId"
	qAdjacent      = "AdjacentTo"
	qGenres        = "Genres"
	qYears         = "Years"
	qPerson        = "Person"
	qPersonTypes   = "PersonTypes"
	qUserData      = "EnableUserData"
	qStatic        = "static"
	qMediaSourceID = "mediaSourceId"
	qAPIKey        = "api_key"
	qPlaySession   = "PlaySessionId"
	qDeviceID      = "DeviceId"
	qTag           = "Tag"
	qVideoCodec    = "VideoCodec"
	qAudioCodec    = "AudioCodec"
	authQuote      = `"`
	headerJSON     = "application/json"
	headerAuth     = "Authorization"
	headerToken    = "X-Emby-Token" // #nosec G101

	pathPublicInfo   = "/System/Info/Public"
	pathAuth         = "/Users/AuthenticateByName"
	pathMe           = "/Users/Me"
	pathNextUp       = "/Shows/NextUp"
	pathPlaying      = "/Sessions/Playing"
	pathProgress     = "/Sessions/Playing/Progress"
	pathStopped      = "/Sessions/Playing/Stopped"
	pathItems        = "/Items"
	pathUserViews    = "/UserViews"
	pathResume       = "/UserItems/Resume"
	pathPlaybackInfo = "/PlaybackInfo"
	pathVideos       = "/Videos/"
	pathStream       = "/stream"
	pathMaster       = "/master.m3u8"

	fieldsHome     = "Overview,UserData,RecursiveItemCount,ChildCount"
	fieldsLatest   = "Overview,UserData,ProductionYear,DateCreated,Genres"
	fieldsChildren = "Overview,UserData,ProductionYear,DateCreated,Genres,People," +
		"RecursiveItemCount,ChildCount,RunTimeTicks,IndexNumber,ParentIndexNumber,SeriesName,SeasonName"
	fieldsSearch = "Overview,UserData,ProductionYear,DateCreated,Genres,People," +
		"SeriesName,IndexNumber,ParentIndexNumber"
	fieldsItem = "Overview,UserData,ProductionYear,DateCreated,Genres,People,MediaSources," +
		"RecursiveItemCount,ChildCount,RunTimeTicks,IndexNumber,ParentIndexNumber," +
		"SeriesName,SeasonName,SeriesId,SeasonId"
	fieldsSeason   = "UserData,IndexNumber,ChildCount"
	fieldsEpisode  = "Overview,UserData,RunTimeTicks,IndexNumber,ParentIndexNumber,SeriesName,SeasonName,SeriesId,SeasonId"
	fieldsAdjacent = "UserData,RunTimeTicks,IndexNumber,ParentIndexNumber,SeriesName,SeasonId,SeriesId"

	sortNameIndex = "SortName,IndexNumber"
	sortCreated   = "DateCreated"
	sortRandom    = "Random"
	orderAsc      = "Ascending"
	orderDesc     = "Descending"

	includeLatest = "Movie,Episode,Series"
	includeSearch = "Movie,Series,Episode"
	includeMovies = "Movie"
	includeSeries = "Series"

	personActor = "Actor"
	genreSep    = "|"
	yearSep     = ","
	boolTrue    = "true"
	boolFalse   = "false"

	yearMin     = 1888
	yearMax     = 2100
	yearSpanCap = 80

	playDirect    = "DirectPlay"
	playTranscode = "Transcode"
	eventTime     = "TimeUpdate"

	codecH264 = "h264"
	codecAAC  = "aac"

	TypeSeries     = "Series"
	TypeSeason     = "Season"
	TypeMovie      = "Movie"
	TypeEpisode    = "Episode"
	TypeVideo      = "Video"
	TypeMusicVideo = "MusicVideo"
	TypeFolder     = "Folder"
	TypeCollection = "CollectionFolder"
	TypePlaylist   = "Playlist"
	TypeBoxSet     = "BoxSet"
	TypeHome       = "Home"

	CollectionMovies = "movies"
	CollectionTV     = "tvshows"
)
