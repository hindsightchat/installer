package theme

import (
	"image/color"

	"github.com/hindsightchat/installer/src/colours"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Dark struct{}

func (d *Dark) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return colours.Background
	case theme.ColorNameButton:
		return colours.Accent
	case theme.ColorNameDisabledButton:
		return colours.Offline
	case theme.ColorNameDisabled:
		return colours.TextMuted
	case theme.ColorNameForeground:
		return colours.TextPrimary
	case theme.ColorNamePrimary:
		return colours.Accent
	case theme.ColorNameFocus:
		return colours.Accent
	case theme.ColorNameHover:
		return colours.AccentHover
	case theme.ColorNameInputBackground:
		return colours.InputBg
	case theme.ColorNameInputBorder:
		return colours.Border
	case theme.ColorNamePlaceHolder:
		return colours.TextMuted
	case theme.ColorNameScrollBar:
		return colours.Scrollbar
	case theme.ColorNameShadow:
		return colours.Shadow
	case theme.ColorNameSuccess:
		return colours.Success
	case theme.ColorNameWarning:
		return colours.Warning
	case theme.ColorNameError:
		return colours.Error
	case theme.ColorNameHeaderBackground:
		return colours.HeaderBg
	case theme.ColorNameSeparator:
		return colours.Border
	case theme.ColorNameMenuBackground:
		return colours.BackgroundCard
	case theme.ColorNameOverlayBackground:
		return colours.BackgroundCard
	default:
		return theme.DarkTheme().Color(name, theme.VariantDark)
	}
}

func (d *Dark) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (d *Dark) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (d *Dark) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInnerPadding:
		return 10
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 22
	case theme.SizeNameSubHeadingText:
		return 16
	case theme.SizeNameScrollBar:
		return 10
	case theme.SizeNameScrollBarSmall:
		return 4
	default:
		return theme.DefaultTheme().Size(name)
	}
}
