function git_compare_commits
    if test (count $argv) -ne 2
        echo "Usage: git_compare_commits <commitA> <commitB>"
        return 1
    end

    set commitA $argv[1]
    set commitB $argv[2]

    echo "Comparing commits: $commitA and $commitB"
    echo ""

    # Get the list of changed files between commitA and commitB
    set changed_files (git diff --name-only "$commitA" "$commitB")

    if test -z "$changed_files"
        echo "No files changed between $commitA and $commitB."
        return 0
    end

    echo "Files changed:"
    for file in $changed_files
        echo "- $file"
    end
    echo ""

    # Iterate over each changed file
    for file in $changed_files
        echo "--- File: $file ---"
        echo ""

        echo "--- Content at commit $commitA ---"
        if git cat-file -e "$commitA":"$file" >/dev/null 2>&1
            git show "$commitA":"$file"
        else
            echo "[File did not exist at commit $commitA]"
        end
        echo ""

        echo "--- Content at commit $commitB ---"
        if git cat-file -e "$commitB":"$file" >/dev/null 2>&1
            git show "$commitB":"$file"
        else
            echo "[File did not exist at commit $commitB]"
        end
        echo ""
        echo --------------------
        echo ""
    end
end
