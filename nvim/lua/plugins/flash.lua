return {
  "folke/flash.nvim",
  event = "VeryLazy",
  opts = {
    modes = {
      char = {
        jump_labels = true,
        multi_line = true,
      },
    },
    label = {
      style = "overlay",
      min_pattern_length = 0,
    },
    jump = {
      autoindent = true,
    },
    prompt = {
      enabled = false,
    },
  },
  keys = {
    {
      "<leader>jj",
      mode = { "n", "x", "o" },
      function()
        require("flash").jump()
      end,
      desc = "Flash: jump to pattern",
    },
    {
      "<leader>js",
      mode = { "n", "x", "o" },
      function()
        require("flash").treesitter()
      end,
      desc = "Flash: treesitter select",
    },
    {
      "<leader>jr",
      mode = { "n", "x", "o" },
      function()
        require("flash").remote()
      end,
      desc = "Flash: remote operation",
    },
    {
      "<leader>jR",
      mode = { "n", "x", "o" },
      function()
        require("flash").treesitter_search()
      end,
      desc = "Flash: treesitter search",
    },
    {
      "<c-s>",
      mode = { "c" },
      function()
        require("flash").toggle()
      end,
      desc = "Flash: toggle search",
    },
  },
}
