return {
  "karb94/neoscroll.nvim",
  event = "VeryLazy",
  opts = {
    mappings = {
      "<C-u>",
      "<C-d>",
      "<C-b>",
      "<C-f>",
      "zt",
      "zz",
      "zb",
    },
    hide_cursor = true,
    stop_eof = true,
    respect_scrolloff = true,
    cursor_scrolls_alone = true,
    pre_hook = function(info)
      if info == "cursorline" then
        vim.wo.cursorline = false
      end
    end,
    post_hook = function(info)
      if info == "cursorline" then
        vim.wo.cursorline = true
      end
    end,
    performance_mode = false,
  },
}
