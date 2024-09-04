#!/bin/sh

# Uncomment to print commands instead of executing them
#debug="echo "

TRUNK="main"
DATE="$(date '+%B %d, %Y')"

BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [[ "${BRANCH}" != "${TRUNK}" ]]; then
  echo "Release must be prepped on \`${TRUNK}\` branch." >&2
  exit 1
fi

echo "Preparing changelog for release..."

if [[ ! -f CHANGELOG.md ]]; then
  echo "Error: CHANGELOG.md not found."
  exit 2
fi

RELEASE="$(sed -r -n 's/^## v?([0-9.]+) \(Unreleased\)/\1/p' CHANGELOG.md)"
if [[ "${RELEASE}" == "" ]]; then
  echo "Error: could not determine next release in CHANGELOG.md" >&2
  exit 3
fi


# Ensure latest changes are checked out
( set -x; ${debug}git pull --rebase origin "${TRUNK})" )

# Set the date for the latest release
( set -x; ${debug}sed -r "s/^(## \[[0-9.]+\]) \(Unreleased\)/\1 (${DATE})/i" CHANGELOG.md )

echo "Preparing release v${RELEASE}..."

(
  set -x
  ${debug}git checkout -b "release-prep/${RELEASE}"
  ${debug}git add CHANGELOG.md
  ${debug}git commit -m "Prepare release v${RELEASE}"
  ${debug}git push origin"
)


