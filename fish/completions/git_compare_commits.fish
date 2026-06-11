# 禁用默认文件补全
complete -c git_compare_commits -f

# 1. 补全 Git 分支和 Tag
complete -c git_compare_commits -a "(git for-each-ref --format='%(refname:short)' refs/heads/ refs/tags/ refs/remotes/)"

# 2. 补全最近 20 个 Commit (格式为: Hash <Tab> Message)
# 使用 %x09 代表 Tab 键，Fish 会自动将其识别为 选项+描述 的分隔符
complete -c git_compare_commits -a "(git log --pretty=format:'%h%x09%s' -n 20)"
