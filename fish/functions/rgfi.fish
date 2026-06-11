# ~/.config/fish/functions/rgfi.fish

function rgfi -d "Interactively search file content with rg and fzf"
    # 这可以确保 fzf 启动时有一个空的输入流，从而可以接受用户的动态输入。
    set -l result (echo "" | fzf --ansi \
        --delimiter ':' \
        --prompt 'Rg> ' \
        --header 'Type to search content in git files. Press ENTER to open in editor.' \
        --preview 'bat --style=numbers --color=always --highlight-line {2} {1}' \
        --preview-window 'right:55%:border-left:+{2}+3/3,~3' \
        --bind "change:reload(rg --line-number --no-heading --color=always {q} || true)" \
        --bind "start:reload(rg --line-number --no-heading --color=always {q} || true)")

    # 后续处理逻辑保持不变
    if test -n "$result"
        set -l file (echo "$result" | awk -F: '{print $1}')
        set -l line (echo "$result" | awk -F: '{print $2}')

        # ！！！请在这里修改成你自己的编辑器！！！
        if command -v nvim >/dev/null
            nvim "+$line" "$file"
        else if command -v vim >/dev/null
            vim "+$line" "$file"
        else
            echo "Editor not configured. Openning: $file at line $line"
        end
    end
end
