#!/usr/bin/env bash
# From https://github.com/oneclickvirt/UnlockTests

set -euo pipefail

repo_url="https://github.com/oneclickvirt/UnlockTests/releases/download/output"
cdn_success_url=""
cdn_urls=(
  "https://cdn0.spiritlhl.top/"
  "http://cdn3.spiritlhl.net/"
  "http://cdn1.spiritlhl.net/"
  "http://cdn2.spiritlhl.net/"
)

asset_os=""
asset_arch=""
asset=""
bin_name="ut"
install_dir="${INSTALL_DIR:-}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Required command not found: $1" >&2
    exit 1
  fi
}

is_interactive() {
  [ -t 0 ] && [ -z "${CI:-}" ] && [ -z "${NO_INTERACTIVE:-}" ] && [ -z "${NONINTERACTIVE:-}" ] && [ -z "${UT_INSTALL_NONINTERACTIVE:-}" ]
}

confirm_yes() {
  [ "${YES:-}" = "1" ] || [ "${YES:-}" = "true" ] || [ "${UT_INSTALL_ASSUME_YES:-}" = "1" ] || [ "${UT_INSTALL_ASSUME_YES:-}" = "true" ]
}

normalize_path() {
  local path="$1"
  if command -v cygpath >/dev/null 2>&1; then
    cygpath -u "$path" 2>/dev/null || printf '%s\n' "$path"
  else
    printf '%s\n' "$path"
  fi
}

windows_path() {
  local path="$1"
  if command -v cygpath >/dev/null 2>&1; then
    cygpath -w "$path" 2>/dev/null || printf '%s\n' "$path"
  else
    printf '%s\n' "$path"
  fi
}

detect_os_arch() {
  local os arch
  os="$(uname -s)"
  arch="$(uname -m)"

  case "$os" in
    Linux*) asset_os="linux" ;;
    Darwin*) asset_os="darwin" ;;
    FreeBSD*) asset_os="freebsd" ;;
    OpenBSD*) asset_os="openbsd" ;;
    MINGW* | MSYS* | CYGWIN*) asset_os="windows"; bin_name="ut.exe" ;;
    *)
      echo "Unsupported operating system: $os" >&2
      exit 1
      ;;
  esac

  case "$arch" in
    x86_64 | amd64 | x64) asset_arch="amd64" ;;
    i386 | i486 | i586 | i686 | x86) asset_arch="386" ;;
    aarch64 | aarch64_be | arm64 | arm64v8) asset_arch="arm64" ;;
    arm | armel | armhf | armv5 | armv5l | armv6 | armv6l | armv7 | armv7l | armv8l) asset_arch="arm" ;;
    loongarch64 | loong64) asset_arch="loong64" ;;
    mips64el | mips64le) asset_arch="mips64le" ;;
    mips64) asset_arch="mips64" ;;
    mipsel | mipsle | mipsisa32r6el) asset_arch="mipsle" ;;
    mips | mipsisa32r6) asset_arch="mips" ;;
    ppc64el | ppc64le) asset_arch="ppc64le" ;;
    ppc64) asset_arch="ppc64" ;;
    riscv64 | riscv64gc) asset_arch="riscv64" ;;
    s390x) asset_arch="s390x" ;;
    *)
      echo "Unsupported architecture: $arch" >&2
      exit 1
      ;;
  esac

  case "$asset_os:$asset_arch" in
    darwin:amd64 | darwin:arm64) ;;
    windows:amd64 | windows:386 | windows:arm64) ;;
    freebsd:amd64 | freebsd:386 | freebsd:arm64 | freebsd:arm) ;;
    openbsd:amd64 | openbsd:386 | openbsd:arm64 | openbsd:arm) ;;
    linux:amd64 | linux:386 | linux:arm64 | linux:arm | linux:riscv64 | linux:mips64 | linux:mips64le | linux:mipsle | linux:mips | linux:ppc64 | linux:ppc64le | linux:s390x | linux:loong64) ;;
    windows:arm)
      echo "Unsupported release: Windows 32-bit ARM is not provided." >&2
      exit 1
      ;;
    *)
      echo "Unsupported release target: ${asset_os}/${asset_arch} (detected ${os}/${arch})" >&2
      exit 1
      ;;
  esac

  asset="ut-${asset_os}-${asset_arch}"
}

default_system_dir() {
  if [ "$asset_os" = "windows" ]; then
    if [ -n "${ProgramFiles:-}" ]; then
      normalize_path "${ProgramFiles}\\UnlockTests\\bin"
    else
      printf '%s\n' "/c/Program Files/UnlockTests/bin"
    fi
  else
    printf '%s\n' "/usr/bin"
  fi
}

default_user_dir() {
  if [ "$asset_os" = "windows" ]; then
    if [ -n "${LOCALAPPDATA:-}" ]; then
      normalize_path "${LOCALAPPDATA}\\Microsoft\\WindowsApps"
    else
      normalize_path "${HOME:-.}\\AppData\\Local\\Microsoft\\WindowsApps"
    fi
  else
    printf '%s\n' "${XDG_BIN_HOME:-${HOME:-.}/.local/bin}"
  fi
}

is_admin() {
  if [ "$asset_os" = "windows" ]; then
    if command -v powershell.exe >/dev/null 2>&1; then
      powershell.exe -NoProfile -NonInteractive -Command "[Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent() | ForEach-Object { if (\$_.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) { exit 0 } else { exit 1 } }" >/dev/null 2>&1
    else
      return 1
    fi
  else
    [ "$(id -u)" -eq 0 ]
  fi
}

can_write_dir() {
  local dir="$1"
  mkdir -p -- "$dir" 2>/dev/null && [ -w "$dir" ]
}

install_binary() {
  local source="$1" dir="$2" destination="${dir}/${bin_name}"
  mkdir -p -- "$dir"
  if command -v install >/dev/null 2>&1; then
    install -m 0755 -- "$source" "$destination"
  else
    cp -- "$source" "$destination"
    chmod 0755 -- "$destination"
  fi
  echo "Installed ${bin_name} to ${destination}"
}

install_binary_elevated() {
  local source="$1" dir="$2" destination="${dir}/${bin_name}"
  if [ "$asset_os" = "windows" ]; then
    if ! command -v powershell.exe >/dev/null 2>&1; then
      return 1
    fi
    UT_INSTALL_SRC="$(windows_path "$source")" \
    UT_INSTALL_DIR="$(windows_path "$dir")" \
    UT_INSTALL_DEST="$(windows_path "$destination")" \
      powershell.exe -NoProfile -NonInteractive -Command 'Start-Process -FilePath powershell.exe -Verb RunAs -Wait -ArgumentList @("-NoProfile","-NonInteractive","-Command","New-Item -ItemType Directory -Force -LiteralPath $env:UT_INSTALL_DIR | Out-Null; Copy-Item -LiteralPath $env:UT_INSTALL_SRC -Destination $env:UT_INSTALL_DEST -Force")' >/dev/null 2>&1
    if [ -f "$destination" ]; then
      echo "Installed ${bin_name} to ${destination}"
      return 0
    fi
    return 1
  fi
  if command -v sudo >/dev/null 2>&1; then
    sudo mkdir -p -- "$dir"
    if command -v install >/dev/null 2>&1; then
      sudo install -m 0755 -- "$source" "$destination"
    else
      sudo cp -- "$source" "$destination"
      sudo chmod 0755 -- "$destination"
    fi
    echo "Installed ${bin_name} to ${destination}"
    return 0
  fi
  if command -v doas >/dev/null 2>&1; then
    doas mkdir -p -- "$dir"
    if command -v install >/dev/null 2>&1; then
      doas install -m 0755 -- "$source" "$destination"
    else
      doas cp -- "$source" "$destination"
      doas chmod 0755 -- "$destination"
    fi
    echo "Installed ${bin_name} to ${destination}"
    return 0
  fi
  return 1
}

choose_install_dir() {
  local system_dir user_dir choice
  system_dir="$(default_system_dir)"
  user_dir="$(default_user_dir)"

  if [ -n "$install_dir" ]; then
    install_dir="$(normalize_path "$install_dir")"
    return
  fi

  if is_admin || can_write_dir "$system_dir"; then
    install_dir="$system_dir"
    return
  fi

  if ! is_interactive || confirm_yes; then
    install_dir="$user_dir"
    echo "No administrator/root permission detected; falling back to user install: ${install_dir}"
    return
  fi

  echo "No administrator/root permission detected."
  echo "1) Try elevated install to ${system_dir}"
  echo "2) Install for current user to ${user_dir}"
  printf 'Choose [1/2] (default 2): '
  read -r choice || choice="2"
  case "$choice" in
    1) install_dir="$system_dir" ;;
    *) install_dir="$user_dir" ;;
  esac
}

check_cdn() {
  local origin_url="$1"
  for cdn_url in "${cdn_urls[@]}"; do
    if curl -fsSL -k --max-time 6 -- "${cdn_url}${origin_url}" | grep -q "success"; then
      cdn_success_url="$cdn_url"
      return 0
    fi
    sleep 0.5
  done
  cdn_success_url=""
}

download_asset() {
  local asset_name="$1" target="$2" url="${repo_url}/${asset_name}" download_url
  download_url="${cdn_success_url}${url}"
  echo "Downloading ${asset_name}..."
  if command -v wget >/dev/null 2>&1; then
    if wget -q -O "$target" -- "$download_url"; then
      return 0
    fi
    echo "wget download failed, falling back to curl" >&2
  fi
  curl -fsSL -o "$target" -- "$download_url"
}

need_cmd curl
detect_os_arch
choose_install_dir

tmp_file="$(mktemp)"
trap 'rm -f -- "$tmp_file"' EXIT

check_cdn "https://raw.githubusercontent.com/spiritLHLS/ecs/main/back/test"
if [ -n "$cdn_success_url" ]; then
  echo "CDN available, using CDN"
else
  echo "No CDN available, using GitHub release URL"
fi

download_asset "$asset" "$tmp_file"
chmod 0755 -- "$tmp_file"

if install_binary "$tmp_file" "$install_dir" 2>/dev/null; then
  exit 0
fi

if [ -n "${INSTALL_DIR:-}" ]; then
  echo "Failed to install to INSTALL_DIR=${install_dir}" >&2
  exit 1
fi

if [ "$install_dir" = "$(default_system_dir)" ]; then
  if is_interactive && ! confirm_yes; then
    if install_binary_elevated "$tmp_file" "$install_dir"; then
      exit 0
    fi
  fi
  install_dir="$(default_user_dir)"
  echo "Falling back to user install: ${install_dir}"
  install_binary "$tmp_file" "$install_dir"
else
  echo "Failed to install to ${install_dir}" >&2
  exit 1
fi
