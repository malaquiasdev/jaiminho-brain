---
name: brain-save
description: Saves a fact, decision, idea, or working preference to a PARA-organized, OKF-formatted brain (an Obsidian vault), filing it under the right project, area, or resource and merging into an existing note instead of duplicating, then commits it. Use when the user says "save this to the brain", "note this down", "remember this decision", or runs /brain-save with the content.
compatibility: Designed for Claude Code. Requires git and a PARA-structured markdown vault at $BRAIN_DIR (default ~/Documents/brain).
---

# brain-save

Usage: `/brain-save "<fact, decision, or idea>"`

Brain root: `${BRAIN_DIR:-$HOME/Documents/brain}`. Always write the full expression in shell commands; variables do not persist between calls. The vault follows PARA; its `index.md` has the rules.

Without an argument, take the last relevant fact or decision from the conversation and confirm it with the user before writing.

## Steps

1. `grep -ril` the brain (excluding `.git` and `.obsidian`) for 2–4 key terms, in both English and Portuguese.
2. If a note already covers the subject, edit it and bump its `timestamp`. If the new content contradicts the old, replace it and say so in the reply.
3. Otherwise pick the destination, first match wins:
   - A rule for how the agent should work, valid in any project → `2-Areas/claude/preferencias/<slug>.md`, plus one line in the `## Preferências` list of `index.md`. Add a `**Como aplicar:**` line.
   - Serves an active goal in `1-Projects/` → that project.
   - Belongs to an ongoing responsibility (a repo, a product, a team) → that repo's existing folder under `2-Areas/` (find it with `find ... -type d -name "<repo>"`; it may sit under a parent folder, see `index.md`). None found → follow the layout in `index.md` and ask.
   - General reference or interest → `3-Resources/<topic>/`.

   Ask when two fit equally. Write `<destination>/<kebab-slug>.md`, content in the vault's language:

   ```markdown
   ---
   type: <Decision | Reference | Preference | Note>
   title: <short title>
   description: "<one line: what this note says>"
   timestamp: YYYY-MM-DD
   tags:
   ---
   # <short title>

   <the fact, 1–5 lines>

   **Por quê:** <reason or context>
   ```

   `type`: `Decision` for a choice with a why, `Reference` for a contract, spec, or fact to look up, `Preference` for the agent rules above, `Note` for anything else.
4. One idea per file. Link related notes with `[[file-name]]` (full path when the name repeats, like `index`). No secrets (tokens, passwords, webhook URLs).
5. Add or update the note's line in `<destination>/index.md`, using its `title` and `description`. Keep indexes current (OKF): every folder has an `index.md` (`type: Index`, `title`, quoted `description`) listing one line per subfolder and note as `- [[path/from/root|Title]] — description`. A new folder gets its own `index.md` and a line in its parent's.
6. Commit: `git -C "${BRAIN_DIR:-$HOME/Documents/brain}" add -A && git -C "${BRAIN_DIR:-$HOME/Documents/brain}" commit -m "save: <title>"`.
7. Reply with the path, one line.

## Edge cases

- Fact about one repo's code that belongs in that repo's docs or the agent's project memory: suggest that instead, save only if the user insists.
- Brain root missing or not a git repo: stop and tell the user; do not create it.
