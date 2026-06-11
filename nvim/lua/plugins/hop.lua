return {
  "phaazon/hop.nvim",
  branch = "v2",
  config = function()
    require("hop").setup({
      keys = "qwertasdfgzxcv",
      quit_key = "<Esc>",
    })
  end,
}
