#!/bin/bash

latest=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
IFS='.' read -r major minor patch <<< "${latest#v}"

case ${1:-patch} in
    major) new="v$((major+1)).0.0" ;;
    minor) new="v${major}.$((minor+1)).0" ;;
    patch) new="v${major}.${minor}.$((patch+1))" ;;
    *) echo "Usage: bump.sh [major|minor|patch]" >&2; exit 1 ;;
esac

git tag -a "$new"

