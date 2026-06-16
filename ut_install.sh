#!/usr/bin/env bash
# From https://github.com/oneclickvirt/UnlockTests

set -euo pipefail

repo_url="https://github.com/oneclickvirt/UnlockTests/releases/download/output"
install_dir="${INSTALL_DIR:-/usr/bin}"
cdn_success_url=""
cdn_urls=(
  "https://cdn0.spiritlhl.top/"
  "http://cdn3.spiritlhl.net/"
  "http://cdn1.spiritlhl.net/"
  "http://cdn2.spiritlhl.net/"
)

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Required command not found: $1" >&2
    exit 1
  fi
}

check_cdn() {
  local origin_url="$1"
  for cdn_url in "${cdn_urls[@]}"; do
    if curl -fsSL -k --max-time 6 "${cdn_url}${origin_url}" | grep -q "success"; then
      cdn_success_url="$cdn_url"
      return 0
    fi
    sleep 0.5
  done
  cdn_success_url=""
}

detect_asset() {
  local os arch asset_os asset_arch
  os="$(uname -s)"
  arch="$(uname -m)"
  case "$os" in
    Linux) asset_os="linux" ;;
    Darwin) asset_os="darwin" ;;
    FreeBSD) asset_os="freebsd" ;;
    OpenBSD) asset_os="openbsd" ;;
    *)
      echo "Unsupported operating system: $os" >&2
      exit 1
      ;;
  esac

  case "$arch" in
    x86_64 | x86 | amd64 | x64) asset_arch="amd64" ;;
    i386 | i686) asset_arch="386" ;;
    aarch64 | arm64) asset_arch="arm64" ;;
    armv7l | armv7 | armv8 | armv8l) asset_arch="arm" ;;
    s390x | riscv64 | mips64 | mips64le | mips | mipsle | ppc64 | ppc64le) asset_arch="$arch" ;;
    *)
      echo "Unsupported architecture: $arch" >&2
      exit 1
      ;;
  esac

  if [ "$asset_os" = "darwin" ] && [ "$asset_arch" = "386" ]; then
    echo "Unsupported architecture for Darwin release: $arch" >&2
    exit 1
  fi
  if [ "$asset_os" != "linux" ]; then
    case "$asset_arch" in
      amd64 | 386 | arm64 | arm) ;;
      *)
        echo "Unsupported architecture for ${asset_os} release: $arch" >&2
        exit 1
        ;;
    esac
  fi
  printf 'ut-%s-%s\n' "$asset_os" "$asset_arch"
}

download_asset() {
  local asset="$1" target="$2" url="${repo_url}/${asset}"
  local download_url="${cdn_success_url}${url}"
  echo "Downloading ${asset}..."
  if command -v wget >/dev/null 2>&1; then
    if wget -q -O "$target" "$download_url"; then
      return 0
    fi
    echo "wget download failed, falling back to curl" >&2
  fi
  curl -fsSL -o "$target" "$download_url"
}

need_cmd curl
asset="$(detect_asset)"
tmp_file="$(mktemp)"
trap 'rm -f "$tmp_file"' EXIT

check_cdn "https://raw.githubusercontent.com/spiritLHLS/ecs/main/back/test"
if [ -n "$cdn_success_url" ]; then
  echo "CDN available, using CDN"
else
  echo "No CDN available, using GitHub release URL"
fi

download_asset "$asset" "$tmp_file"
chmod 0755 "$tmp_file"
mkdir -p "$install_dir"
install -m 0755 "$tmp_file" "${install_dir}/ut"
echo "Installed ut to ${install_dir}/ut"
