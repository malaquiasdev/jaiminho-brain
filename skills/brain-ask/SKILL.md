---
name: brain-ask
description: Answers a question from the notes in a PARA-organized, OKF-formatted brain (an Obsidian vault), citing each note it uses. Read-only. Use when the user asks "why did we decide X", "what did I do on Y", "is there anything in the brain about Z", or runs /brain-ask with a question.
compatibility: Designed for Claude Code. Requires a PARA-structured markdown vault at $BRAIN_DIR (default ~/Documents/brain).
allowed-tools: Bash(grep:*) Bash(ls:*) Bash(find:*) Read Grep Glob
---

# brain-ask

Usage: `/brain-ask "<question>"`

Brain root: `${BRAIN_DIR:-$HOME/Documents/brain}`. Read-only: never edit or commit.

## Steps

1. Navigate first: read the root `index.md`, then the `index.md` of each folder that looks relevant; each lists its subfolders and notes with a one-line description.
   Then extract search terms from the question, plus synonyms in English and Portuguese and any repo or issue names, and `grep -ril` the whole brain, `4-Archives/` included, excluding `.git` and `.obsidian`, for anything the indexes missed. `grep -l '^type: Decision'` and similar narrow by note type.
2. Question about a repo or project → start in the folder named after it: `find` for `-type d -name "<repo-or-slug>"` under `2-Areas/` and `1-Projects/` (it may sit under a parent folder). Question about a time range ("last week") → `find` the `sessoes/` folders and filter by the date in the file name.
3. Read the most relevant notes (up to ~10) and follow `[[links]]` when they help.
4. Answer directly: the first line is the answer. Cite the note behind each claim by path.
5. Nothing found → say it is not in the brain. Knowledge from outside the brain only when clearly marked as such.
