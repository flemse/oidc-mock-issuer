#!/usr/bin/env bash
# Updates all GitHub Actions pins in workflow files to the latest release SHA.
# Compatible with bash 3.2+ (macOS default).
# Usage: ./scripts/pin-actions.sh [workflow files...]
# Defaults to all .github/workflows/*.yml files.
set -euo pipefail

if [[ $# -gt 0 ]]; then
  WORKFLOW_FILES=("$@")
else
  WORKFLOW_FILES=(.github/workflows/*.yml)
fi

resolve_sha() {
  local repo="$1" tag="$2"
  local obj_sha obj_type
  obj_sha=$(gh api "repos/$repo/git/ref/tags/$tag" --jq '.object.sha')
  obj_type=$(gh api "repos/$repo/git/ref/tags/$tag" --jq '.object.type')
  if [[ "$obj_type" == "tag" ]]; then
    # Annotated tag — resolve to the commit it points to
    obj_sha=$(gh api "repos/$repo/git/tags/$obj_sha" --jq '.object.sha')
  fi
  echo "$obj_sha"
}

# Collect all unique "owner/repo" references from the workflow files
actions=$(
  grep -hE 'uses: [A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+@' "${WORKFLOW_FILES[@]}" \
    | sed -E 's/.*uses: ([A-Za-z0-9_.-]+\/[A-Za-z0-9_.-]+)@.*/\1/' \
    | sort -u
)

echo "Found $(echo "$actions" | wc -l | tr -d ' ') unique actions to update:"
echo "$actions" | sed 's/^/  /'
echo

# For each action: fetch latest tag+sha, then update all workflow files in one pass
while IFS= read -r repo; do
  printf '  %s → ' "$repo"
  tag=$(gh api "repos/$repo/releases/latest" --jq '.tag_name' 2>/dev/null || true)
  if [[ -z "$tag" ]]; then
    echo "no release found, skipping"
    continue
  fi
  sha=$(resolve_sha "$repo" "$tag")
  echo "$tag ($sha)"

  for file in "${WORKFLOW_FILES[@]}"; do
    [[ -f "$file" ]] || continue
    sed -i.bak -E \
      "s|uses: ${repo}@[A-Za-z0-9_./+-]+( # .*)?|uses: ${repo}@${sha} # ${tag}|g" \
      "$file"
    rm -f "$file.bak"
  done
done <<< "$actions"

echo
echo "Done. Updated files:"
printf '  %s\n' "${WORKFLOW_FILES[@]}"
