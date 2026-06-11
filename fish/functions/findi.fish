# ~/.config/fish/functions/findi.fish
function findi --description "在指定文件中搜索字符串 (大小写不敏感)"
    # 检查参数数量
    if test (count $argv) -ne 2
        echo "用法: findi <目标文件> <匹配字符串>" >&2
        return 1
    end

    set target_file $argv[1]
    set search_string $argv[2]

    # 检查目标文件是否存在且可读
    if not test -f "$target_file"
        echo "findi: 错误: '$target_file' 不是一个普通文件。" >&2
        return 1
    end
    if not test -r "$target_file"
        echo "findi: 错误: 无法读取文件 '$target_file'." >&2
        return 1
    end

    # 使用 grep 搜索
    # -i: 忽略大小写
    # -n: 显示行号
    # --color=always: 强制高亮显示匹配项
    # --: 标记选项结束，防止搜索字符串以 '-' 开头时被误认为是选项
    grep --color=always -i -n -- "$search_string" "$target_file"
end
