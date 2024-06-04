#!/bin/bash
#From https://github.com/oneclickvirt/UnlockTests
#2024.05.21

rm -rf /usr/bin/UT
os=$(uname -s)
arch=$(uname -m)

case $os in
  Linux)
    case $arch in
      "x86_64" | "x86" | "amd64" | "x64")
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-linux-amd64
        ;;
      "i386" | "i686")
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-linux-386
        ;;
      "armv7l" | "armv8" | "armv8l" | "aarch64" | "arm64")
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-linux-arm64
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
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-darwin-amd64
        ;;
      "i386" | "i686")
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-darwin-386
        ;;
      "armv7l" | "armv8" | "armv8l" | "aarch64" | "arm64")
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-darwin-arm64
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
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-freebsd-amd64
        ;;
      "i386" | "i686")
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-freebsd-386
        ;;
      "armv7l" | "armv8" | "armv8l" | "aarch64" | "arm64")
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-freebsd-arm64
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
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-openbsd-amd64
        ;;
      "i386" | "i686")
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-openbsd-386
        ;;
      "armv7l" | "armv8" | "armv8l" | "aarch64" | "arm64")
        wget -O UT https://github.com/oneclickvirt/UnlockTests/releases/download/output/UT-openbsd-arm64
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

chmod 777 UT
if [ ! -f /usr/bin/UT ]; then
  mv UT /usr/bin/
  UT
else
  ./UT
fi