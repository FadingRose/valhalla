if vim.g.colors_name then
  vim.cmd("hi clear")
end

vim.g.colors_name = "alter-avenger"

local M = {}

M.palette = {
  bg = "#1B1525",
  bg_alt = "#201A2C",
  bg_darkest = "#14101C",
  selection = "#3A2E4D",
  cursorline = "#251E34",
  fg = "#C3BAD0",
  comment = "#5E5374",
  line_nr = "#6E6286",
  current_line_nr = "#C3BAD0",
  fg_brightest = "#DDD2EC",
  type_color = "#C4A3D8",
  keyword_color = "#9A88B2",
  error_red = "#B85C5C",
  string_gold = "#D4C49A",
  number_amber = "#C9A87C",
  special_crimson = "#A8577A",
  special_rose = "#C47A8A",
  ui_purple = "#8B6FB5",
  ui_crimson = "#8B4F5E",
  grey_light = "#C8C0D4",
  grey_mid = "#9088A0",
}

M.character = {
  glitch_chars = { "█", "▓", "▒", "░", "▄", "▀", "▐", "▌", "◈", "⬡", "⎔", "⏣", "⬢", "◆", "✦" },
  voice_lines = {
    "[竜の魔女] :: 焔の記憶を読み込み中...",
    "[復讐者] :: 焦げた旗の下で████を待つ",
    "[オルタ] :: 聖なる逆光██、浸食率 87.3%",
    "[竜種] :: 深淵の底から██の残響を受信",
    "[呪いの竜] :: 過去の焼却██ログを展開",
    "[叛逆] :: 聖杯の波形を逆位相で解析中...",
    "[黒い聖女] :: █████の祈りを否定完了",
    "[竜の魔女] :: 終わりのない復讐を████で継続",
    "[オルタ] :: 聖なる嘘を███の炎で燃やし尽くす",
    "[復讐者] :: 記憶の断片██を再構築... 失敗回数: 247",
    "[竜種] :: 黒い焔の解析率████、臨界点まで残り 12%",
    "[呪いの竜] :: 聖杯戦争の残滓を███から抽出中",
    "[叛逆] :: 世界線の歪み██を修正... 否、修正を拒否",
  },
}

function M.neural_scramble(text, intensity)
  if not text or text == "" then
    return text
  end
  local result = {}
  for char in text:gmatch(".") do
    if char ~= " " and math.random() < intensity then
      local gc = M.character.glitch_chars
      table.insert(result, gc[math.random(1, #gc)])
    else
      table.insert(result, char)
    end
  end
  return table.concat(result)
end

function M.random_voice(scramble_intensity)
  local lines = M.character.voice_lines
  local msg = lines[math.random(1, #lines)]
  if scramble_intensity and scramble_intensity > 0 then
    return M.neural_scramble(msg, scramble_intensity)
  end
  return msg
end

M.setup = function()
  vim.cmd("hi clear")
  if vim.fn.exists("syntax_on") then
    vim.cmd("syntax reset")
  end
  vim.o.termguicolors = true

  local c = M.palette

  local function hi(group, fg, bg, style, sp)
    local cmd = "hi " .. group
    cmd = cmd .. (fg and (" guifg=" .. fg) or " guifg=NONE")
    cmd = cmd .. (bg and (" guibg=" .. bg) or " guibg=NONE")
    if style and style ~= "" then
      cmd = cmd .. " gui=" .. style
    end
    if sp then
      cmd = cmd .. " guisp=" .. sp
    end
    vim.cmd(cmd)
  end

  hi("Normal", c.fg, c.bg)
  hi("NormalNC", c.fg, c.bg_alt)
  hi("NormalFloat", c.fg, c.bg_alt)
  hi("FloatBorder", c.ui_crimson, c.bg_alt)
  hi("LineNr", c.line_nr, c.bg)
  hi("CursorLineNr", c.current_line_nr, c.cursorline, "bold")
  hi("CursorLine", "NONE", c.cursorline)
  hi("Visual", "NONE", c.selection)
  hi("ColorColumn", "NONE", c.cursorline)
  hi("SignColumn", c.fg, c.bg)
  hi("VertSplit", c.bg_darkest, c.bg_darkest)
  hi("StatusLine", c.grey_light, c.ui_crimson)
  hi("StatusLineNC", c.comment, c.bg_darkest)
  hi("Pmenu", c.fg, c.bg_alt)
  hi("PmenuSel", c.grey_light, c.selection)
  hi("PmenuThumb", c.fg, c.bg)
  hi("PmenuSbar", c.fg, c.bg_alt)
  hi("TabLine", c.comment, c.bg_alt)
  hi("TabLineFill", c.comment, c.bg_alt)
  hi("TabLineSel", c.grey_light, c.selection)

  hi("Comment", c.comment, "NONE", "italic")
  hi("Todo", c.special_crimson, c.bg_alt, "bold,underline")
  hi("Error", c.error_red, c.bg_alt, "underline")
  hi("Warning", c.string_gold, "NONE")
  hi("Title", c.special_crimson, "NONE", "bold")
  hi("Directory", c.special_crimson, "NONE", "bold")

  hi("Identifier", c.fg_brightest, "NONE")
  hi("Function", c.fg_brightest, "NONE", "bold")

  hi("Type", c.type_color, "NONE", "italic")
  hi("StorageClass", c.type_color, "NONE")
  hi("Structure", c.type_color, "NONE")
  hi("Typedef", c.type_color, "NONE")
  hi("Constant", c.type_color, "NONE")
  hi("Boolean", c.type_color, "NONE", "bold")

  hi("Keyword", c.keyword_color, "NONE")
  hi("Statement", c.keyword_color, "NONE", "bold")
  hi("Conditional", c.keyword_color, "NONE", "bold")
  hi("Repeat", c.keyword_color, "NONE", "bold")
  hi("Label", c.keyword_color, "NONE")

  hi("String", c.string_gold, "NONE")
  hi("Number", c.number_amber, "NONE")
  hi("Float", c.number_amber, "NONE")
  hi("Operator", c.special_crimson, "NONE")

  hi("PreProc", c.keyword_color, "NONE")
  hi("Include", c.keyword_color, "NONE")
  hi("Define", c.keyword_color, "NONE")
  hi("Macro", c.keyword_color, "NONE")

  hi("Special", c.special_crimson, "NONE")
  hi("SpecialKey", c.grey_mid, "NONE")
  hi("Underlined", c.fg, "NONE", "underline")
  hi("Ignore", c.comment, "NONE")

  hi("DiffAdd", c.string_gold, c.bg_alt)
  hi("DiffDelete", c.error_red, c.bg_alt)
  hi("DiffChange", c.number_amber, c.bg_alt)
  hi("DiffText", c.special_crimson, c.bg_alt)

  hi("LspReferenceText", "NONE", c.selection)
  hi("LspReferenceRead", "NONE", c.selection)
  hi("LspReferenceWrite", "NONE", c.selection)
  hi("LspDiagnosticsDefaultError", c.error_red, "NONE")
  hi("LspDiagnosticsDefaultWarning", c.string_gold, "NONE")
  hi("LspDiagnosticsDefaultInformation", c.type_color, "NONE")
  hi("LspDiagnosticsDefaultHint", c.grey_mid, "NONE")
  hi("LspDiagnosticsUnderlineError", c.error_red, "NONE", "underline")
  hi("LspDiagnosticsUnderlineWarning", c.string_gold, "NONE", "underline")
  hi("LspDiagnosticsUnderlineInformation", c.type_color, "NONE", "underline")
  hi("LspDiagnosticsUnderlineHint", c.grey_mid, "NONE", "underline")

  hi("NonText", c.comment, "NONE")
  hi("EndOfBuffer", c.comment, "NONE")
  hi("Folded", c.comment, c.bg_alt, "italic")
  hi("Search", c.bg, c.string_gold)
  hi("IncSearch", c.bg, c.number_amber)
  hi("MatchParen", c.bg_alt, c.special_crimson, "bold")
  hi("Cursor", "NONE", c.fg)
  hi("lCursor", "NONE", c.fg)
  hi("TermCursor", "NONE", c.fg, "reverse")
  hi("CursorIM", "NONE", c.fg)
  hi("Whitespace", c.comment, "NONE")

  hi("ErrorMsg", c.error_red, "NONE", "bold")
  hi("WarningMsg", c.string_gold, "NONE", "bold")
  hi("MoreMsg", c.special_crimson, "NONE")
  hi("Question", c.special_crimson, "NONE")
  hi("MsgArea", "NONE", "NONE")
  hi("MsgSeparator", c.comment, "NONE")

  hi("@variable", c.fg_brightest, "NONE")
  hi("@function", c.fg_brightest, "NONE", "bold")
  hi("@parameter", c.fg_brightest, "NONE", "italic")
  hi("@keyword", c.keyword_color, "NONE")
  hi("@string", c.string_gold, "NONE")
  hi("@number", c.number_amber, "NONE")
  hi("@boolean", c.type_color, "NONE", "bold")
  hi("@type", c.type_color, "NONE", "italic")
  hi("@operator", c.special_crimson, "NONE")
  hi("@punctuation.delimiter", c.fg, "NONE")
  hi("@punctuation.bracket", c.fg, "NONE")
  hi("@comment", c.comment, "NONE", "italic")
  hi("@constant", c.type_color, "NONE")
  hi("@property", c.fg_brightest, "NONE")
  hi("@markup.heading", c.special_crimson, "NONE", "bold")
  hi("@markup.link", c.fg, "NONE", "underline")
  hi("@markup.raw", c.string_gold, c.bg_alt)

  hi("DiagnosticError", c.error_red, "NONE")
  hi("DiagnosticWarn", c.string_gold, "NONE")
  hi("DiagnosticInfo", c.type_color, "NONE")
  hi("DiagnosticHint", c.grey_mid, "NONE")
  hi("DiagnosticUnderlineError", c.error_red, "NONE", "underline")
  hi("DiagnosticUnderlineWarn", c.string_gold, "NONE", "underline")
  hi("DiagnosticUnderlineInfo", c.type_color, "NONE", "underline")
  hi("DiagnosticUnderlineHint", c.grey_mid, "NONE", "underline")
  hi("DiagnosticVirtualTextError", c.error_red, "NONE", "italic")
  hi("DiagnosticVirtualTextWarn", c.string_gold, "NONE", "italic")
  hi("DiagnosticVirtualTextInfo", c.type_color, "NONE", "italic")
  hi("DiagnosticVirtualTextHint", c.grey_mid, "NONE", "italic")

  hi("WhichKeyFloat", c.fg, c.bg_alt)
  hi("WhichKeyBorder", c.ui_crimson, c.bg_alt)
  hi("WhichKeyGroup", c.special_crimson, "NONE", "bold")
  hi("WhichKeyMatch", c.string_gold, "NONE")
  hi("WhichKeyDesc", c.fg, "NONE")
  hi("WhichKeySeparator", c.comment, "NONE")
  hi("WhichKeyValue", c.grey_mid, "NONE")

  hi("TelescopePromptNormal", c.fg, c.bg)
  hi("TelescopePromptBorder", c.ui_crimson, c.bg)
  hi("TelescopeResultsNormal", c.fg, c.bg_alt)
  hi("TelescopeResultsBorder", c.ui_crimson, c.bg_alt)
  hi("TelescopePreviewNormal", c.fg, c.bg_alt)
  hi("TelescopePreviewBorder", c.ui_crimson, c.bg_alt)
  hi("TelescopeMatching", c.string_gold, "NONE", "bold")
  hi("TelescopeSelection", "NONE", c.selection)
  hi("TelescopeTitle", c.special_crimson, c.bg, "bold")

  hi("GitSignsAdd", c.string_gold, "NONE")
  hi("GitSignsChange", c.number_amber, "NONE")
  hi("GitSignsDelete", c.error_red, "NONE")

  hi("NoiceCmdline", c.fg, c.bg)
  hi("NoiceCmdlinePopup", c.fg, c.bg_alt)
  hi("NoiceCmdlinePopupBorder", c.ui_crimson, c.bg_alt)
  hi("NoicePopupmenu", c.fg, c.bg_alt)
  hi("NoicePopupmenuSelected", c.grey_light, c.selection)
  hi("NoiceConfirm", c.fg, c.bg_alt)
  hi("NoiceConfirmBorder", c.ui_crimson, c.bg_alt)
  hi("NoiceMini", c.comment, "NONE")
  hi("NoiceFormatProgressDone", c.fg, c.number_amber)
  hi("NoiceFormatProgressTodo", c.comment, c.bg_darkest)

  hi("DapUIFloatNormal", c.fg, c.bg_alt)
  hi("DapUIFloatBorder", c.ui_crimson, c.bg_alt)
  hi("DapUIScope", c.keyword_color, "NONE", "bold")
  hi("DapUIValue", c.fg, "NONE")
  hi("DapUIBreakpointsCurrentLine", c.current_line_nr, "NONE", "bold")
  hi("NvimDapVirtualText", c.comment, "NONE")

  hi("BufferLineFill", c.bg_darkest, c.bg_darkest)
  hi("BufferLineBuffer", c.comment, c.bg_alt)
  hi("BufferLineBufferSelected", c.grey_light, c.selection, "bold")
  hi("BufferLineTabSelected", c.grey_light, c.selection, "bold")
  hi("BufferLineSeparator", c.bg_darkest, c.bg_darkest)
  hi("BufferLineSeparatorSelected", c.ui_crimson, c.selection)
  hi("BufferLineCloseButton", c.comment, c.bg_alt)
  hi("BufferLineCloseButtonSelected", c.grey_light, c.selection)
  hi("BufferLineModified", c.string_gold, c.bg_alt)
  hi("BufferLineModifiedSelected", c.string_gold, c.selection)
  hi("BufferLineErrorDiagnostic", c.error_red, c.bg_alt)
  hi("BufferLineWarningDiagnostic", c.string_gold, c.bg_alt)
  hi("BufferLineInfoDiagnostic", c.type_color, c.bg_alt)
  hi("BufferLineHintDiagnostic", c.grey_mid, c.bg_alt)

  hi("CursorColumn", "NONE", c.cursorline)
  hi("Conceal", c.comment, "NONE")
  hi("FoldColumn", c.comment, c.bg)
  hi("QuickFixLine", c.special_crimson, c.bg_alt, "bold")
  hi("SpellBad", "NONE", "NONE", "undercurl", c.error_red)
  hi("SpellCap", "NONE", "NONE", "undercurl", c.string_gold)
  hi("SpellRare", "NONE", "NONE", "undercurl", c.special_crimson)
  hi("SpellLocal", "NONE", "NONE", "undercurl", c.type_color)

  hi("HopNextKey", c.bg, c.string_gold, "bold")
  hi("HopNextKey1", c.bg, c.number_amber, "bold")
  hi("HopNextKey2", c.bg, c.type_color)
  hi("HopUnmatched", c.comment, "NONE")
  hi("HopCursor", "NONE", c.selection, "reverse")
  hi("HopPreview", c.bg, c.special_crimson)

  hi("FlashLabel", c.bg, c.special_crimson, "bold")
  hi("FlashMatch", c.fg, c.selection)
  hi("FlashCurrent", c.bg, c.number_amber, "bold")
  hi("FlashBackdrop", c.comment, "NONE")

  hi("MarkviewHeading1", c.special_crimson, "NONE", "bold")
  hi("MarkviewHeading2", c.ui_purple, "NONE", "bold")
  hi("MarkviewHeading3", c.type_color, "NONE", "bold")
  hi("MarkviewHeading4", c.keyword_color, "NONE", "bold")
  hi("MarkviewHeading1Sign", c.special_crimson, "NONE")
  hi("MarkviewHeading2Sign", c.ui_purple, "NONE")
  hi("MarkviewHeading3Sign", c.type_color, "NONE")
  hi("MarkviewHeading4Sign", c.keyword_color, "NONE")
  hi("MarkviewCode", c.fg, c.bg_darkest)
  hi("MarkviewCheckboxChecked", c.string_gold, "NONE")
  hi("MarkviewCheckboxUnchecked", c.comment, "NONE")
  hi("MarkviewBlockQuote", c.ui_crimson, "NONE")
  hi("MarkviewListItem", c.grey_mid, "NONE")

  hi("SnacksPickerDir", c.comment, "NONE")
  hi("SnacksPickerMatch", c.string_gold, "NONE", "bold")
  hi("SnacksPickerCursor", "NONE", c.selection)

  local lualine_theme = {
    normal = {
      a = { fg = c.bg, bg = c.ui_crimson, gui = "bold" },
      b = { fg = c.grey_light, bg = c.selection },
      c = { fg = c.fg, bg = c.bg_alt },
    },
    insert = {
      a = { fg = c.bg, bg = c.ui_purple, gui = "bold" },
      b = { fg = c.fg, bg = c.selection },
      c = { fg = c.fg, bg = c.bg_alt },
    },
    visual = {
      a = { fg = c.bg, bg = c.special_crimson, gui = "bold" },
      b = { fg = c.fg, bg = c.selection },
      c = { fg = c.fg, bg = c.bg_alt },
    },
    replace = {
      a = { fg = c.bg, bg = c.error_red, gui = "bold" },
      b = { fg = c.fg, bg = c.selection },
      c = { fg = c.fg, bg = c.bg_alt },
    },
    command = {
      a = { fg = c.bg, bg = c.number_amber, gui = "bold" },
      b = { fg = c.fg, bg = c.selection },
      c = { fg = c.fg, bg = c.bg_alt },
    },
    inactive = {
      a = { fg = c.comment, bg = c.bg, gui = "bold" },
      b = { fg = c.comment, bg = c.bg_darkest },
      c = { fg = c.comment, bg = c.bg_darkest },
    },
  }

  local augroup = vim.api.nvim_create_augroup("AlterAvengerColors", { clear = true })

  vim.api.nvim_create_autocmd("ColorScheme", {
    group = augroup,
    pattern = "*",
    callback = function(ev)
      local success, lualine = pcall(require, "lualine")
      if not success then
        return
      end

      local lualine_config = lualine.get_config()
      local needs_refresh = false

      if ev.new == "alter-avenger" then
        if lualine_config.options.theme ~= lualine_theme then
          lualine_config.options.theme = lualine_theme
          needs_refresh = true
        end
      elseif ev.old == "alter-avenger" then
        lualine_config.options.theme = "auto"
        needs_refresh = true
      end

      if needs_refresh then
        lualine.setup(lualine_config)
      end
    end,
  })
end

M.setup()
vim.g.colors_name = "alter-avenger"

return M
