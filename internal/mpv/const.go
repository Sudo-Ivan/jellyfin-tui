package mpv

import "time"

const (
	ipcTimeout    = 4 * time.Second
	ipcRetry      = 40 * time.Millisecond
	ipcDialTry    = 200 * time.Millisecond
	defaultVolume = 100
	hashBytes     = 8
	maxZipMember  = 64 << 20

	dirPerm     = 0o700
	binPerm     = 0o700
	zipFilePerm = 0o644
	zipDirPerm  = 0o755

	cacheApp  = "jellyfin-tui"
	cacheMpv  = "mpv-"
	titleName = "jellyfin-tui"
	zipMagic  = "PK"
	exeUnix   = "mpv"
	exeWin    = "mpv.exe"
	ytdlNo    = "--ytdl=no"

	sharedDir        = "shared"
	libDir           = "lib"
	ldLinuxPrefix    = "ld-linux"
	ldMuslPrefix     = "ld-musl"
	maxSymlinkTarget = 4 << 10

	ipcSetProp  = "set_property"
	flagReplace = "replace"
	optStart    = "start"
	optTitle    = "force-media-title"

	propTime     = "time-pos"
	propDuration = "duration"
	propPause    = "pause"
	propVolume   = "volume"
	propEOF      = "eof-reached"
	propIdle     = "idle-active"
	propPath     = "path"
	propTrack    = "track-list"
	propAid      = "aid"
	propSid      = "sid"

	trackAudio = "audio"
	trackSub   = "sub"

	alangEng = "eng"
	sidNo    = "no"

	evEndFile  = "end-file"
	evProp     = "property-change"
	reasonEOF  = "eof"
	reasonStop = "stop"
)
