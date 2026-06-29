# alter-avenger themed custom tab bar for Kitty
# Palette synced from nvim/colors/alter-avenger.lua

from kitty.borders import BorderColor
from kitty.fast_data_types import Screen, DrawableCache
from kitty.rgb import color_as_srgb
from kitty.tab_bar import DrawData, TabExtraData
from kitty.utils import color_as_int

# === alter-avenger palette ===
BG = (0x1B, 0x15, 0x25)            # #1B1525
BG_ALT = (0x20, 0x1A, 0x2C)        # #201A2C
FG = (0xC3, 0xBA, 0xD0)            # #C3BAD0
COMMENT = (0x5E, 0x53, 0x74)       # #5E5374
SELECTION = (0x3A, 0x2E, 0x4D)     # #3A2E4D
STRING_GOLD = (0xD4, 0xC4, 0x9A)   # #D4C49A
ERROR_RED = (0xB8, 0x5C, 0x5C)     # #B85C5C
SPECIAL_CRIMSON = (0xA8, 0x57, 0x7A)  # #A8577A
GREY_MID = (0x90, 0x88, 0xA0)      # #9088A0

SEP = "│"


def _draw_text(draw_data, screen, tab, index, max_title_length):
    """Draw tab content: index, title, bell/activity symbols."""
    tab_data = tab

    # Bell symbol (red)
    if tab_data.needs_attention:
        screen.cursor.fg = _rgb_to_int(ERROR_RED)
        screen.draw(tab_data.bell_symbol or "🔔 ")
    else:
        # Activity symbol
        if tab_data.activity_symbol:
            screen.draw(tab_data.activity_symbol)

    # Tab index
    screen.cursor.fg = _rgb_to_int(GREY_MID)
    screen.draw(f"{index}:")

    # Title
    title = tab_data.title
    if max_title_length:
        title = title[:max_title_length]
    screen.draw(f" {title}")


def _rgb_to_int(rgb):
    return (rgb[0] << 16) | (rgb[1] << 8) | rgb[2]


def draw_tab(draw_data, screen, tab, before, max_title_length, index, is_last, is_active):
    """Custom tab renderer matching alter-avenger + lualine aesthetic."""
    # Background
    bg = _rgb_to_int(BG) if is_active else _rgb_to_int(BG_ALT)
    screen.cursor.bg = bg

    # Left separator (not for first tab)
    if index > 1:
        screen.cursor.fg = _rgb_to_int(SELECTION)
        screen.draw(f" {SEP} ")
        screen.cursor.bg = bg

    # Tab content
    if is_active:
        screen.cursor.fg = _rgb_to_int(STRING_GOLD)
        screen.cursor.bold = True
    else:
        screen.cursor.fg = _rgb_to_int(COMMENT)
        screen.cursor.bold = False

    _draw_text(draw_data, screen, tab, index, max_title_length)

    screen.cursor.bold = False

    # Right padding
    screen.draw(" ")

    # End separator for last tab
    if is_last:
        screen.cursor.fg = _rgb_to_int(SELECTION)
        screen.cursor.bg = _rgb_to_int(BG)
        screen.draw(f" {SEP} ")

    return screen.cursor.x
