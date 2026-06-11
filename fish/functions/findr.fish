# ~/.config/fish/functions/findr.fish
function findr --description "在指定文件夹内递归搜索包含字符串的文件 (大小写不敏感)"
    # 检查参数数量
    if test (count $argv) -ne 2
        echo "用法: findr <目标文件夹> <匹配字符串>" >&2
        return 1
    end

    set target_dir $argv[1]
    set search_string $argv[2]

    # 检查目标文件夹是否存在且可读可搜索
    if not test -d "$target_dir"
        echo "findr: 错误: '$target_dir' 不是一个目录。" >&2
        return 1
    end
    # 目录需要读(r)和执行(x - 搜索)权限
    # if not test -r "$target_dir" -o not test -x "$target_dir"
    #     echo "findr: 错误: 无法读取或搜索目录 '$target_dir'." >&2
    #     return 1
    # end
    #
    # 使用 grep 递归搜索
    # -r: 递归搜索子目录
    # -i: 忽略大小写
    # -n: 显示行号
    # --color=always: 强制高亮显示匹配项
    # --exclude-dir: 排除指定的目录 (例如 .git, node_modules 等)
    # --: 标记选项结束
    # grep --color=always -r -i -n \
    #     --exclude-dir={.git,.svn,node_modules,build,dist,__pycache__} \
    #     -- "$search_string" "$target_dir"
    #
    # --- 更好的选择：ripgrep (rg) ---
    # 如果你安装了 ripgrep (rg)，它通常更快并且默认会忽略 .gitignore 文件
    # 可以取消注释下面的行来使用 rg (如果已安装)
    if command -q rg
        rg --color=always -i -n --glob '!.git/' --glob '!node_modules/' -- "$search_string" "$target_dir"
    else
        # Fallback to grep if rg is not installed
        grep --color=always -r -i -n \
            --exclude-dir={.git,.svn,node_modules,build,dist,__pycache__} \
            -- "$search_string" "$target_dir"
    end
end
