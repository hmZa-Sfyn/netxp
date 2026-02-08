#!/usr/bin/env python3
"""
Super dangerous cache-che spam script - for memes / testing only
"""

import random
import subprocess
from pathlib import Path
from datetime import datetime

CACHE_FILE = Path("./cache.che")
REPO_ROOT = Path.cwd()

def git(*args, check=True, capture=False):
    """Small git command wrapper"""
    cmd = ["git", *args]
    kw = {}
    if capture:
        kw["capture_output"] = True
        kw["text"] = True
    if check:
        kw["check"] = True
    
    try:
        return subprocess.run(cmd, **kw, cwd=REPO_ROOT)
    except subprocess.CalledProcessError as e:
        print("Git command failed!")
        print(" ".join(cmd))
        print(e.stderr if e.stderr else "(no error output)")
        raise


def main():
    if not CACHE_FILE.exists():
        print(f"Error: {CACHE_FILE} not found in current directory")
        return 1

    # Generate random number (you can change range)
    number = random.randint(0, 19_000_696_969)

    timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # --- Option 1: Append line (most common for cache/test files) ---
    try:
        with CACHE_FILE.open("a", encoding="utf-8") as f:
            f.write(f"\n{number}  # {timestamp}\n")
        print(f"Appended number: {number}")
    except Exception as e:
        print("Failed to write to cache.che:", e)
        return 1

    # Alternative style - overwrite with only the number
    # with CACHE_FILE.open("w", encoding="utf-8") as f:
    #     f.write(str(number))

    try:
        # Stage the file
        git("add", str(CACHE_FILE))

        # Commit
        commit_msg = f"cache test: {number}"
        git("commit", "-m", commit_msg)
        print(f"Committed → {commit_msg}")

        # Push (to master - old name, many repos now use main)
        ##git("push", "origin", "master")
        ##print("Pushed to origin/master")

    except subprocess.CalledProcessError:
        print("\nGit operations failed - stopping here.")
        return 1

    print("\nDone (probably)")
    return 0


if __name__ == "__main__":
    for x in range(1,999): 
        main()

    exit(0)