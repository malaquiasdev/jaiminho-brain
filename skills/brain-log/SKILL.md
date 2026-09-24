---
name: brain-log
description: Logs the current coding session as a note in a PARA-organized, OKF-formatted brain (an Obsidian vault), filed under the matching project or area, capturing what was done, decisions, links, and follow-ups, then commits it. Use when the user asks to log, save, or record the session or today's work, or runs /brain-log with an optional topic.
compatibility: Designed for Claude Code. Requires git and a PARA-structured markdown vault at $BRAIN_DIR (default ~/Documents/brain).
---

# brain-log

Usage: `/brain-log [topic]`

Brain root: `${BRAIN_DIR:-$HOME/Documents/brain}`. Always write the full expression in shell commands; variables do not persist between calls. The vault follows PARA; its `index.md` has the rules.

## Steps

1. Repo = `basename -s .git "$(git remote get-url origin)"`, if inside a repo (falls back to `basename "$(git rev-parse --show-toplevel)"` when there is no remote). The remote name stays right inside a worktree. Topic = the argument, or a 3–6 word summary of the session.
2. Pick the destination folder:
   - The session advanced something in `1-Projects/` (same epic, issue, or goal), or an active change in the repo's `openspec/changes/<id>/` → that project (`<id>` is its folder name).
   - Otherwise, inside a repo → the repo's existing folder: `find "${BRAIN_DIR:-$HOME/Documents/brain}/2-Areas" -type d -name "<repo>" -not -path "*/.git/*"`. None found → follow the layout in `index.md` (it may group repos under a parent folder, e.g. per client) and ask which parent before creating it.
   - Otherwise → ask the user which project or area, listing `1-Projects/` and `2-Areas/`. A new project folder only when the user confirms it has a goal and an end.
3. Look in `<destination>/sessoes/` for a note with today's date on the same topic. If one exists, update it instead of creating another.
4. Write `<destination>/sessoes/YYYY-MM-DD - <topic>.md`, content in the vault's language (see `index.md`):

   ```markdown
   ---
   type: Session
   timestamp: YYYY-MM-DD
   repo: <repo>
   branch: <branch>
   base: <short sha of the base commit>
   issues: [<IDs>]
   estado: <commitado / PR aberta / não commitado / em prod>
   tags:
   ---
   # <topic>

   ## Contexto
   ## O que foi feito
   ## Decisões
   ## Links
   ## Pendências
   ```

   - Drop frontmatter fields without a value and empty sections.
   - Every decision states why, and the alternative that was rejected.
   - Links are full URLs (PR, issue, commit), never a bare identifier.
   - No secrets: tokens, passwords, webhook URLs, secret values. Say one exists, never its value.
5. `grep -ril` the brain for the session's key terms and link related notes with `[[file-name]]` (full path when the name repeats, like `index`).
6. Add or update the note's line at the top of `<destination>/sessoes/index.md`: `- [[<path without .md>|YYYY-MM-DD — <topic>]] — <estado>`. Keep indexes current (OKF): every folder has an `index.md` (`type: Index`, `title`, quoted `description`) listing one line per subfolder and note as `- [[path/from/root|Title]] — description`. A new folder gets its own `index.md` and a line in its parent's.
7. Commit: `git -C "${BRAIN_DIR:-$HOME/Documents/brain}" add -A && git -C "${BRAIN_DIR:-$HOME/Documents/brain}" commit -m "log: <destination> - <topic>"`.
8. Reply with the note's path, one line. If the session finished the project, offer to move it to `4-Archives/<year>/`, keeping any parent folder it had (`1-Projects/<parent>/<slug>/` → `4-Archives/<year>/<parent>/<slug>/`) and moving its lines between the indexes.

## Edge cases

- Session with nothing worth keeping (a quick question, no decision): say so and ask before writing.
- Brain root missing or not a git repo: stop and tell the user; do not create it.
