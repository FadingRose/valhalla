# ~/.config/fish/functions/rgf.fish

function rgf -d "Search file content with ripgrep and fzf"
    # 如果没有提供参数，则不执行
    if test (count $argv) -eq 0
        echo "Usage: rgf <pattern>"
        return 1
    end

    # 使用 rg 搜索内容，并把结果管道给 fzf
    # rg 参数:
    # --line-number: 显示行号
    # --no-heading: 不显示文件名标题
    # --color=always: 始终输出颜色代码，以便 fzf 的 --ansi 能解析
    # fzf 参数:
    # --delimiter ':' : 使用冒号作为分隔符，用于后续解析
    # --preview: 使用 bat 预览文件，并高亮匹配的行
    #   {1} 是文件名 (第一个字段)
    #   {2} 是行号 (第二个字段)
    # --preview-window: 定制预览窗口
    set -l result (rg --line-number --no-heading --color=always "$argv" |
        fzf --ansi \
            --delimiter ':' \
            --preview 'bat --style=numbers --color=always --highlight-line {2} {1}' \
            --preview-window 'up,60%,border-bottom,+{2}+3/3,~3')

    # 如果用户在 fzf 中选择了某一行 (result 不为空)
    if test -n "$result"
        # 解析出文件名和行号
        set -l file (echo "$result" | awk -F: '{print $1}')
        set -l line (echo "$result" | awk -F: '{print $2}')

        # 使用你喜欢的编辑器打开文件并跳转到指定行
        # 将 nvim 修改为你的编辑器，如 vim, code, subl 等
        # 对于 VS Code: code --goto "$file:$line"
        if command -v nvim >/dev/null
            nvim "+$line" "$file"
        else if command -v vim >/dev/null
            vim "+$line" "$file"
        else
            echo "No editor found to open the file."
        end
    end
end
