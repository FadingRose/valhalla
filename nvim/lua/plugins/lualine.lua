return {
  "nvim-lualine/lualine.nvim",
  config = function()
    local ok, theme = pcall(require, "alter-avenger")
    local c = ok and theme.palette or {}
    local char = ok and theme.character or {}

    local accent = c.special_crimson or "#ff9e64"

    local function center_msg()
      if theme and theme.random_voice then
        return theme.random_voice(0)
      end
      return ""
    end

    require("lualine").setup({
      options = {
        theme = "auto",
        component_separators = "",
        section_separators = { left = "::", right = "::" },
        disabled_filetypes = { "snacks_picker_list", "dashboard" },
      },
      sections = {
        lualine_a = { { "mode", separator = { left = "::" }, right_padding = 2 } },
        lualine_b = { "filename", "branch" },
        lualine_c = {
          "%=",
          {
            center_msg,
            color = { gui = "italic" },
          },
        },
        lualine_x = {
          {
            require("noice").api.status.message.get_hl,
            cond = require("noice").api.status.message.has,
          },
          {
            require("noice").api.status.command.get,
            cond = require("noice").api.status.command.has,
            color = { fg = accent },
          },
          {
            require("noice").api.status.mode.get,
            cond = require("noice").api.status.mode.has,
            color = { fg = accent },
          },
          {
            require("noice").api.status.search.get,
            cond = require("noice").api.status.search.has,
            color = { fg = accent },
          },
        },
        lualine_y = { "filetype", "progress" },
        lualine_z = {
          { "location", separator = { right = "::" }, left_padding = 2 },
        },
      },
      inactive_sections = {
        lualine_a = { "filename" },
        lualine_b = {},
        lualine_c = {},
        lualine_x = {},
        lualine_y = {},
        lualine_z = { "location" },
      },
      tabline = {
        lualine_a = {
          {
            "buffers",
            show_filename_only = true,
            hide_filename_extension = false,
            show_modified_status = true,
            mode = 2,
            max_length = vim.o.columns,
            symbols = {
              modified = " ●",
              alternate_file = "",
              directory = "",
            },
            buffers_color = {
              active = { fg = c.bg or "#1B1525", bg = c.ui_crimson or "#8B4F5E", gui = "bold" },
              inactive = { fg = c.comment or "#5E5374", bg = c.bg_alt or "#201A2C" },
            },
          },
        },
        lualine_z = {
          {
            "tabs",
            mode = 2,
            max_length = function()
              return vim.o.columns
            end,
            tabs_color = {
              active = { fg = c.bg or "#1B1525", bg = c.special_crimson or "#A8577A", gui = "bold" },
              inactive = { fg = c.comment or "#5E5374", bg = c.bg_alt or "#201A2C" },
            },
          },
        },
      },
      extensions = { "lazy", "nvim-dap-ui", "fugitive", "mason", "trouble" },
    })
  end,
}
