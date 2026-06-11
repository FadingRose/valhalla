return {
  "OXY2DEV/markview.nvim",
  lazy = false,
  dependencies = {
    "nvim-treesitter/nvim-treesitter",
    "nvim-tree/nvim-web-devicons",
  },
  opts = {
    preview = {
      filetypes = { "markdown", "quarto", "rmd" },
      ignore_buftypes = { "nofile" },
      hybrid_modes = { "n" },
      linewise_hybrid_mode = true,
    },
    markdown = {
      headings = {
        heading_1 = {
          style = "label",
          sign = "󰌕 ",
          sign_hl = "MarkviewHeading1Sign",
        },
        heading_2 = {
          style = "label",
          sign = "󰌖 ",
          sign_hl = "MarkviewHeading2Sign",
        },
        heading_3 = {
          style = "label",
          sign = "󰌗 ",
          sign_hl = "MarkviewHeading3Sign",
        },
        heading_4 = {
          style = "label",
          sign = "󰌘 ",
          sign_hl = "MarkviewHeading4Sign",
        },
      },
      code_blocks = {
        style = "language",
        pad_amount = 3,
        language_names = {
          ["solidity"] = "Solidity",
          ["sol"] = "Solidity",
          ["rust"] = "Rust",
          ["go"] = "Go",
          ["lua"] = "Lua",
          ["python"] = "Python",
          ["bash"] = "Bash",
          ["json"] = "JSON",
          ["toml"] = "TOML",
          ["yaml"] = "YAML",
          ["javascript"] = "JS",
          ["typescript"] = "TS",
        },
      },
      block_quotes = {
        enable = true,
        callouts = {
          ["NOTE"] = { title = " Note", hl = "MarkviewNote" },
          ["WARNING"] = { title = " Warning", hl = "MarkviewWarning" },
          ["TODO"] = { title = " Todo", hl = "MarkviewTodo" },
          ["IMPORTANT"] = { title = " Important", hl = "MarkviewImportant" },
          ["CAUTION"] = { title = " Caution", hl = "MarkviewCaution" },
          ["TIP"] = { title = " Tip", hl = "MarkviewTip" },
          ["QUESTION"] = { title = " Question", hl = "MarkviewQuestion" },
        },
      },
      list_items = {
        enable = true,
        indent_size = 2,
        marker_dot = "●",
        marker_minus = "—",
        marker_plus = "◆",
      },
      checkboxes = {
        enable = true,
        checked = {
          text = "✔",
          hl = "MarkviewCheckboxChecked",
        },
        unchecked = {
          text = "✘",
          hl = "MarkviewCheckboxUnchecked",
        },
        custom = {
          ["~"] = { text = "◯", hl = "MarkviewCheckboxPending" },
          ["?"] = { text = "？", hl = "MarkviewCheckboxCancelled" },
        },
      },
      horizontal_rules = {
        style = "dashed",
        text = " ── ",
      },
      tables = {
        enable = true,
        use_virt_lines = true,
      },
    },
  },
}
