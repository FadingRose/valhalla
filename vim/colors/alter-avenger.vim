" ─── alter-avenger.vim ─────────────────────────────────────────────────────
" Vim port of the alter-avenger colorscheme (originally Lua for nvim).
" Maintains the same dark-purple palette; uses pure :hi commands so it runs
" in plain Vim 9 without any Neovim-only APIs.

" Reset
if exists("g:colors_name")
  hi clear
endif
if exists("syntax_on")
  syntax reset
endif
let g:colors_name = "alter-avenger"

set background=dark
if has('termguicolors')
  set termguicolors
endif

" ─── Palette (mirror of nvim/colors/alter-avenger.lua) ─────────────────────
let s:bg            = "#1B1525"
let s:bg_alt        = "#201A2C"
let s:bg_darkest    = "#14101C"
let s:selection     = "#3A2E4D"
let s:cursorline    = "#251E34"
let s:fg            = "#C3BAD0"
let s:comment       = "#5E5374"
let s:line_nr       = "#6E6286"
let s:cur_line_nr   = "#C3BAD0"
let s:fg_brightest  = "#DDD2EC"
let s:type          = "#C4A3D8"
let s:keyword       = "#9A88B2"
let s:error_red     = "#B85C5C"
let s:string_gold   = "#D4C49A"
let s:number_amber  = "#C9A87C"
let s:crimson       = "#A8577A"
let s:rose          = "#C47A8A"
let s:ui_purple     = "#8B6FB5"
let s:ui_crimson    = "#8B4F5E"
let s:grey_light    = "#C8C0D4"
let s:grey_mid      = "#9088A0"

" ─── Core UI ───────────────────────────────────────────────────────────────
hi Normal         guifg=#C3BAD0 guibg=#1B1525 gui=NONE
hi NormalNC       guifg=#C3BAD0 guibg=#201A2C gui=NONE
hi Cursor         guifg=#1B1525 guibg=#DDD2EC gui=NONE
hi CursorLine                    guibg=#251E34 gui=NONE
hi CursorLineNr   guifg=#C3BAD0 guibg=#251E34 gui=bold
hi CursorColumn                  guibg=#251E34 gui=NONE
hi LineNr         guifg=#6E6286 guibg=NONE     gui=NONE
hi SignColumn                    guibg=NONE     gui=NONE
hi VertSplit      guifg=#3A2E4D guibg=NONE     gui=NONE
hi Folded         guifg=#5E5374 guibg=#201A2C gui=NONE
hi FoldColumn     guifg=#6E6286 guibg=NONE     gui=NONE

" ─── Status / tabs / popup ─────────────────────────────────────────────────
hi StatusLine     guifg=#DDD2EC guibg=#201A2C gui=NONE
hi StatusLineNC   guifg=#5E5374 guibg=#14101C gui=NONE
hi TabLine        guifg=#5E5374 guibg=#14101C gui=NONE
hi TabLineFill                  guibg=#14101C gui=NONE
hi TabLineSel     guifg=#DDD2EC guibg=#201A2C gui=bold
hi WildMenu       guifg=#1B1525 guibg=#8B6FB5 gui=bold
hi Pmenu          guifg=#C3BAD0 guibg=#201A2C gui=NONE
hi PmenuSel       guifg=#1B1525 guibg=#8B6FB5 gui=bold
hi PmenuSbar                    guibg=#14101C gui=NONE
hi PmenuThumb                   guibg=#5E5374 gui=NONE

" ─── Search / visual / errors ──────────────────────────────────────────────
hi Search         guifg=#1B1525 guibg=#D4C49A gui=NONE
hi IncSearch      guifg=#1B1525 guibg=#C47A8A gui=bold
hi Visual         guifg=#DDD2EC guibg=#3A2E4D gui=NONE
hi VisualNOS                    guibg=#3A2E4D gui=NONE
hi MatchParen                   guibg=#8B4F5E gui=bold
hi Error          guifg=#B85C5C guibg=#1B1525 gui=bold
hi ErrorMsg       guifg=#B85C5C guibg=NONE     gui=bold
hi WarningMsg     guifg=#D4C49A guibg=NONE     gui=NONE
hi Todo           guifg=#D4C49A guibg=#1B1525 gui=bold,italic
hi NonText        guifg=#5E5374 guibg=NONE     gui=NONE
hi SpecialKey     guifg=#5E5374 guibg=NONE     gui=NONE
hi Conceal        guifg=#9088A0 guibg=NONE     gui=NONE
hi EndOfBuffer    guifg=#1B1525 guibg=NONE     gui=NONE
hi ColorColumn                  guibg=#251E34 gui=NONE
hi SpellBad       guifg=#B85C5C guibg=NONE     gui=undercurl guisp=#B85C5C
hi SpellCap       guifg=#8B6FB5 guibg=NONE     gui=undercurl guisp=#8B6FB5
hi SpellRare      guifg=#C47A8A guibg=NONE     gui=undercurl guisp=#C47A8A
hi SpellLocal     guifg=#C9A87C guibg=NONE     gui=undercurl guisp=#C9A87C
hi Question       guifg=#D4C49A guibg=NONE     gui=bold
hi ModeMsg        guifg=#9088A0 guibg=NONE     gui=NONE
hi MoreMsg        guifg=#D4C49A guibg=NONE     gui=bold
hi Title          guifg=#C4A3D8 guibg=NONE     gui=bold
hi Directory      guifg=#8B6FB5 guibg=NONE     gui=bold
hi diffAdded      guifg=#D4C49A guibg=NONE     gui=NONE
hi diffRemoved    guifg=#B85C5C guibg=NONE     gui=NONE
hi diffChanged    guifg=#C9A87C guibg=NONE     gui=NONE

" ─── Diff-mode highlight (vimdiff / :diffsplit) ────────────────────────────
" Tuned to read clearly against the dark-purple bg without washing out.
hi DiffAdd        guifg=#D4C49A guibg=#251E34 gui=NONE      " added line
hi DiffChange     guifg=#C9A87C guibg=#201A2C gui=NONE      " changed line
hi DiffDelete     guifg=#B85C5C guibg=#1B1525 gui=bold      " deleted line
hi DiffText       guifg=#DDD2EC guibg=#3A2E4D gui=bold      " specific char change

" ─── Syntax groups (port of hi group -> palette mapping) ───────────────────
hi Comment        guifg=#5E5374 guibg=NONE     gui=italic
hi Constant       guifg=#D4C49A guibg=NONE     gui=NONE
hi String         guifg=#D4C49A guibg=NONE     gui=NONE
hi Character      guifg=#C9A87C guibg=NONE     gui=NONE
hi Number         guifg=#C9A87C guibg=NONE     gui=NONE
hi Boolean        guifg=#C9A87C guibg=NONE     gui=NONE
hi Float          guifg=#C9A87C guibg=NONE     gui=NONE
hi Identifier     guifg=#C3BAD0 guibg=NONE     gui=NONE
hi Function       guifg=#C4A3D8 guibg=NONE     gui=NONE
hi Statement      guifg=#9A88B2 guibg=NONE     gui=NONE
hi Conditional    guifg=#9A88B2 guibg=NONE     gui=NONE
hi Repeat         guifg=#9A88B2 guibg=NONE     gui=NONE
hi Label          guifg=#9A88B2 guibg=NONE     gui=NONE
hi Operator       guifg=#9088A0 guibg=NONE     gui=NONE
hi Keyword        guifg=#9A88B2 guibg=NONE     gui=NONE
hi Exception      guifg=#A8577A guibg=NONE     gui=NONE
hi PreProc        guifg=#C47A8A guibg=NONE     gui=NONE
hi Include        guifg=#C47A8A guibg=NONE     gui=NONE
hi Define         guifg=#C47A8A guibg=NONE     gui=NONE
hi Macro          guifg=#C47A8A guibg=NONE     gui=NONE
hi Type           guifg=#C4A3D8 guibg=NONE     gui=NONE
hi StorageClass   guifg=#8B6FB5 guibg=NONE     gui=NONE
hi Structure      guifg=#8B6FB5 guibg=NONE     gui=NONE
hi Typedef        guifg=#8B6FB5 guibg=NONE     gui=NONE
hi Special        guifg=#A8577A guibg=NONE     gui=NONE
hi SpecialChar    guifg=#A8577A guibg=NONE     gui=NONE
hi Tag            guifg=#C47A8A guibg=NONE     gui=bold
hi Delimiter      guifg=#9088A0 guibg=NONE     gui=NONE
hi SpecialComment guifg=#5E5374 guibg=NONE     gui=italic
hi Underlined     guifg=#8B6FB5 guibg=NONE     gui=underline
hi Ignore         guifg=#5E5374 guibg=NONE     gui=NONE

" ─── Markdown ──────────────────────────────────────────────────────────────
hi markdownH1       guifg=#C4A3D8 guibg=NONE gui=bold
hi markdownH2       guifg=#C47A8A guibg=NONE gui=bold
hi markdownH3       guifg=#D4C49A guibg=NONE gui=bold
hi markdownCode     guifg=#D4C49A guibg=#201A2C gui=NONE
hi markdownCodeBlock guifg=#D4C49A guibg=NONE gui=NONE
hi markdownLinkText guifg=#8B6FB5 guibg=NONE gui=underline

" ─── GitSigns-style (placeholder; no plugin required) ──────────────────────
hi GitGutterAdd    guifg=#D4C49A guibg=NONE gui=NONE
hi GitGutterChange guifg=#C9A87C guibg=NONE gui=NONE
hi GitGutterDelete guifg=#B85C5C guibg=NONE gui=NONE

" ─── End ───────────────────────────────────────────────────────────────────
