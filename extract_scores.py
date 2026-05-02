#!/usr/bin/env python3
"""Extract repository scores from habit examiner logs.

Usage:
  python3 extract_scores.py logs/101-176-scores.txt
  python3 extract_scores.py logs/101-176-scores.txt repo/101-176.txt
  python3 extract_scores.py logs/101-176-scores.txt repo/101-176.txt -o logs/101-176-extracted-scores.txt

When a repo list is provided, output follows the repo-list order and any repo
without a completed score is emitted with score 0.
"""

from __future__ import annotations

import argparse
import re
from collections import defaultdict, deque
from pathlib import Path


EXAMINING_RE = re.compile(r"^Examining repo:\s*(.+?)\s*$")
SETUP_FAIL_RE = re.compile(r"^\[SETUP FAIL\]\s+repo:(.+?)\s*$")
SCORE_RE = re.compile(r"^Score:\s*([-+]?\d+(?:\.\d+)?)\s*/\s*[-+]?\d+(?:\.\d+)?\s*$")
MARKDOWN_LINK_RE = re.compile(r"^\[[^\]]+\]\(([^)]+)\)$")


def normalize_repo(value: str) -> str:
    value = value.strip()
    markdown_match = MARKDOWN_LINK_RE.match(value)
    if markdown_match:
        value = markdown_match.group(1).strip()
    return value


def read_repo_list(path: Path) -> list[str]:
    repos: list[str] = []
    for line in path.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        repos.append(normalize_repo(line))
    return repos


def read_scores(path: Path) -> list[tuple[str, str]]:
    entries: list[tuple[str, str]] = []
    current_repo: str | None = None

    for line in path.read_text(encoding="utf-8", errors="replace").splitlines():
        examining_match = EXAMINING_RE.match(line)
        if examining_match:
            current_repo = normalize_repo(examining_match.group(1))
            continue

        setup_fail_match = SETUP_FAIL_RE.match(line)
        if setup_fail_match:
            entries.append((normalize_repo(setup_fail_match.group(1)), "0"))
            current_repo = None
            continue

        score_match = SCORE_RE.match(line)
        if score_match and current_repo:
            entries.append((current_repo, score_match.group(1)))
            current_repo = None

    return entries


def format_scores(log_path: Path, repo_list_path: Path | None) -> list[str]:
    scored_entries = read_scores(log_path)
    if repo_list_path is None:
        return [f"{repo}: {score}" for repo, score in scored_entries]

    scores_by_repo: dict[str, deque[str]] = defaultdict(deque)
    for repo, score in scored_entries:
        scores_by_repo[repo].append(score)

    lines: list[str] = []
    for repo in read_repo_list(repo_list_path):
        score = scores_by_repo[repo].popleft() if scores_by_repo[repo] else "0"
        lines.append(f"{repo}: {score}")
    return lines


def default_output_path(log_path: Path) -> Path:
    return log_path.with_name(f"{log_path.stem}-extracted.txt")


def main() -> int:
    parser = argparse.ArgumentParser(description="Extract repo: score lines from examiner logs.")
    parser.add_argument("log_file", type=Path, help="examiner log file, for example logs/101-176-scores.txt")
    parser.add_argument("repo_list", nargs="?", type=Path, help="optional repo list; missing scores become 0")
    parser.add_argument(
        "-o",
        "--output",
        type=Path,
        help="output file path; defaults to <log-file-stem>-extracted.txt beside the log file",
    )
    args = parser.parse_args()

    output_path = args.output or default_output_path(args.log_file)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text("\n".join(format_scores(args.log_file, args.repo_list)) + "\n", encoding="utf-8")
    print(f"Wrote {output_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
