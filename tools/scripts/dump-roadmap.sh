#!/usr/bin/env bash
#
# Dump the GitHub Project roadmap into docs/roadmap.md so the repo carries a
# plain-markdown mirror of the board (useful for docs sites / other tooling).
#
# The GitHub Project is the source of truth; this file is generated — do not
# edit docs/roadmap.md by hand, re-run this script instead.
#
# Usage:  tools/scripts/dump-roadmap.sh [project_number] [owner]
# Requires: gh (authenticated, with the `project` scope), python3.

set -euo pipefail

NUM="${1:-1}"
OWNER="${2:-adewaleo}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
OUT="$REPO_ROOT/docs/roadmap.md"

meta="$(gh project view "$NUM" --owner "$OWNER" --format json)"
items="$(gh project item-list "$NUM" --owner "$OWNER" --limit 1000 --format json)"

META_JSON="$meta" ITEMS_JSON="$items" NUM="$NUM" OWNER="$OWNER" python3 - "$OUT" <<'PY'
import os, sys, json
from collections import defaultdict

out_path = sys.argv[1]
meta  = json.loads(os.environ["META_JSON"])
items = json.loads(os.environ["ITEMS_JSON"])["items"]
num, owner = os.environ["NUM"], os.environ["OWNER"]

title = meta.get("title") or "Roadmap"
desc  = meta.get("shortDescription") or ""
url   = meta.get("url") or f"https://github.com/users/{owner}/projects/{num}"

# Phase ordering; anything else falls to the end under "Unphased".
PHASE_ORDER = ["P0", "P1", "P2", "P3", "P4", "P5"]
def phase_key(p):
    return (PHASE_ORDER.index(p), p) if p in PHASE_ORDER else (len(PHASE_ORDER), p or "zzz")

by_phase = defaultdict(list)
for it in items:
    by_phase[it.get("phase") or "Unphased"].append(it)

DONE = {"Done"}
def checkbox(status):
    return "[x]" if (status or "") in DONE else "[ ]"

lines = []
lines.append(f"# {title}")
lines.append("")
lines.append("> **Generated file — do not edit by hand.** Mirror of the GitHub Project;")
lines.append(f"> regenerate with `tools/scripts/dump-roadmap.sh`. Source of truth: [Project]({url}).")
lines.append("")
if desc:
    lines.append(desc)
    lines.append("")

# Summary counts
total = len(items)
done = sum(1 for it in items if (it.get("status") or "") in DONE)
lines.append(f"**{done}/{total} done.**")
lines.append("")

for phase in sorted(by_phase, key=phase_key):
    group = sorted(by_phase[phase], key=lambda i: i.get("title", ""))
    gdone = sum(1 for i in group if (i.get("status") or "") in DONE)
    lines.append(f"## {phase} ({gdone}/{len(group)})")
    lines.append("")
    for it in group:
        t = (it.get("title") or "").strip()
        status = it.get("status") or "—"
        body = ((it.get("content") or {}).get("body") or "").strip().replace("\n", " ")
        lines.append(f"- {checkbox(status)} **{t}** — _{status}_")
        if body:
            lines.append(f"  - {body}")
    lines.append("")

with open(out_path, "w") as f:
    f.write("\n".join(lines).rstrip() + "\n")

print(f"Wrote {out_path} ({total} items, {done} done).")
PY
