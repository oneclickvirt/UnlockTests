#!/bin/bash
#From https://github.com/oneclickvirt/UnlockTests
#2024.05.21

rm -rf /usr/bin/ut
os=$(uname -s)
arch=$(uname -m)

case $os in
  Linux)
    case $arch in
      "x86_64" | "x86" | "amd64" | "x64")
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-linux-amd64
        ;;
      "i386" | "i686")
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-linux-386
        ;;
      "armv7l" | "armv8" | "armv8l" | "aarch64" | "arm64")
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-linux-arm64
        ;;
      *)
        echo "Unsupported architecture: $arch"
        exit 1
        ;;
    esac
    ;;
  Darwin)
    case $arch in
      "x86_64" | "x86" | "amd64" | "x64")
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-darwin-amd64
        ;;
      "i386" | "i686")
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-darwin-386
        ;;
      "armv7l" | "armv8" | "armv8l" | "aarch64" | "arm64")
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-darwin-arm64
        ;;
      *)
        echo "Unsupported architecture: $arch"
        exit 1
        ;;
    esac
    ;;
  FreeBSD)
    case $arch in
      amd64)
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-freebsd-amd64
        ;;
      "i386" | "i686")
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-freebsd-386
        ;;
      "armv7l" | "armv8" | "armv8l" | "aarch64" | "arm64")
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-freebsd-arm64
        ;;
      *)
        echo "Unsupported architecture: $arch"
        exit 1
        ;;
    esac
    ;;
  OpenBSD)
    case $arch in
      amd64)
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-openbsd-amd64
        ;;
      "i386" | "i686")
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-openbsd-386
        ;;
      "armv7l" | "armv8" | "armv8l" | "aarch64" | "arm64")
        wget -O ut https://github.com/oneclickvirt/UnlockTests/releases/download/output/ut-openbsd-arm64
        ;;
      *)
        echo "Unsupported architecture: $arch"
        exit 1
        ;;
    esac
    ;;
  *)
    echo "Unsupported operating system: $os"
    exit 1
    ;;
esac

chmod 777 ut
cp ut /usr/bin/
