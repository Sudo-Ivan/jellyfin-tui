#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
linux_dir="$root/internal/mpv/bundle/linux"
windows_dir="$root/internal/mpv/bundle/windows"

do_linux=false
do_windows=false

usage() {
	cat <<'EOF'
Usage: scripts/fetch-mpv.sh [--linux] [--windows] [--all]

Downloads portable mpv builds into internal/mpv/bundle/ for embedding.

  --linux    fetch Linux anylinux AppImage and write bundle/linux/mpv.zip
  --windows  fetch Windows git-release zip and write bundle/windows/mpv.zip
  --all      fetch both (default when no flag is given)

Environment overrides:
  MPV_LINUX_URL    direct download URL for the Linux AppImage
  MPV_WINDOWS_URL  direct download URL for the Windows zip
EOF
}

while [[ $# -gt 0 ]]; do
	case "$1" in
	--linux)
		do_linux=true
		shift
		;;
	--windows)
		do_windows=true
		shift
		;;
	--all)
		do_linux=true
		do_windows=true
		shift
		;;
	-h | --help)
		usage
		exit 0
		;;
	*)
		echo "unknown argument: $1" >&2
		usage >&2
		exit 1
		;;
	esac
done

if ! $do_linux && ! $do_windows; then
	do_linux=true
	do_windows=true
fi

need() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "missing required command: $1" >&2
		exit 1
	fi
}

gh_asset() {
	local repo="$1"
	local filter="$2"
	local url="${3:-}"
	if [[ -z "$url" ]]; then
		need curl
		need python3
		url="$(curl -fsSL "https://api.github.com/repos/$repo/releases/latest" |
			python3 -c "import json,sys,re; d=json.load(sys.stdin); pat=re.compile(sys.argv[1]);
for a in d.get('assets',[]):
 n=a.get('name','')
 if pat.search(n):
  print(a['browser_download_url']); break
else:
 raise SystemExit('no asset matched: '+sys.argv[1])" "$filter")"
	fi
	printf '%s' "$url"
}

fetch_windows() {
	need curl
	need unzip
	need zip

	mkdir -p "$windows_dir"
	local url
	url="$(gh_asset "mpv-player/mpv" 'x86_64-pc-windows-msvc\.zip$' "${MPV_WINDOWS_URL:-}")"
	if [[ "$url" == *pdb* ]]; then
		echo "windows url looks like a pdb package: $url" >&2
		exit 1
	fi

	local work
	work="$(mktemp -d "$root/.cache/mpv-fetch-windows.XXXXXX")"
	trap 'rm -rf "$work"' RETURN

	echo "fetching windows mpv from $url"
	curl -fsSL -o "$work/mpv-src.zip" "$url"
	unzip -q "$work/mpv-src.zip" -d "$work/extract"
	if [[ ! -f "$work/extract/mpv.exe" ]]; then
		echo "windows archive has no mpv.exe" >&2
		exit 1
	fi
	rm -f "$windows_dir/mpv.zip"
	(
		cd "$work/extract"
		zip -q -r "$windows_dir/mpv.zip" . \
			-i '*.exe' '*.dll' '*.com' \
			-x '*.pdb'
	)
	echo "wrote $windows_dir/mpv.zip"
}

fetch_linux() {
	need curl
	need zip
	if [[ "$(uname -s)" != "Linux" ]]; then
		echo "linux mpv fetch requires Linux (AppImage extraction)" >&2
		exit 1
	fi

	mkdir -p "$linux_dir"
	local url
	url="$(gh_asset "pkgforge-dev/mpv-AppImage" 'anylinux-x86_64\.AppImage$' "${MPV_LINUX_URL:-}")"

	local work
	work="$(mktemp -d "$root/.cache/mpv-fetch-linux.XXXXXX")"
	trap 'rm -rf "$work"' RETURN

	echo "fetching linux mpv from $url"
	curl -fsSL -o "$work/mpv.AppImage" "$url"
	chmod +x "$work/mpv.AppImage"
	(
		cd "$work"
		./mpv.AppImage --appimage-extract >/dev/null
	)
	local extract="$work/squashfs-root"
	if [[ ! -d "$extract" ]]; then
		extract="$work/AppDir"
	fi
	if [[ ! -d "$extract" ]]; then
		echo "appimage extract produced no output directory" >&2
		exit 1
	fi
	for part in bin lib share; do
		if [[ ! -e "$extract/$part" ]]; then
			echo "appimage extract missing $part" >&2
			exit 1
		fi
	done
	rm -f "$linux_dir/mpv.zip"
	(
		cd "$extract"
		zip -q -r "$linux_dir/mpv.zip" bin lib share
	)
	echo "wrote $linux_dir/mpv.zip"
}

mkdir -p "$root/.cache"

if $do_windows; then
	fetch_windows
fi
if $do_linux; then
	fetch_linux
fi
