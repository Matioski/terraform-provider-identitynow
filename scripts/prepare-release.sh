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

RELEASE="$(sed -r -n 's/^## \[([0-9.]+)\] \(Unreleased\)/\1/p' CHANGELOG.md)"
if [[ "${RELEASE}" == "" ]]; then
  echo "Error: could not determine next release in CHANGELOG.md" >&2
  exit 3
fi


# Ensure latest changes are checked out
( set -x; ${debug}git pull --rebase origin "${TRUNK})" )

# Set the date for the latest release
( set -x; ${debug}sed -r -i.bak "s/^(## \[[0-9.]+\]) \(Unreleased\)/\1 (${DATE})/i" CHANGELOG.md )

${debug}rm CHANGELOG.md.bak

echo "Preparing release v${RELEASE}..."
echo "  - Creating release preparation branch 'release-prep/${RELEASE}'"
${debug}git checkout -b "release-prep/${RELEASE}"
echo "  - Commiting CHANGELOG.md"
${debug}git commit CHANGELOG.md -m "Prepare release v${RELEASE}"
echo "  - Pushing to origin"
${debug}git push --set-upstream origin "release-prep/${RELEASE}"
