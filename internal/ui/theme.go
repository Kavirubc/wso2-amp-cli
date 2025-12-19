package ui

import "github.com/charmbracelet/lipgloss"

// WSO2 Color Palette - Based on design system tokens
var (
	// Neutrals (Note: inverted in original schema for dark theme)
	Black = lipgloss.Color("#000000")
	White = lipgloss.Color("#ffffff")

	// Gray scale
	Gray100 = lipgloss.Color("#1a1a1a")
	Gray200 = lipgloss.Color("#333333")
	Gray300 = lipgloss.Color("#4d4d4d")
	Gray400 = lipgloss.Color("#666666")
	Gray500 = lipgloss.Color("#808080")
	Gray600 = lipgloss.Color("#999999")
	Gray700 = lipgloss.Color("#b3b3b3")
	Gray800 = lipgloss.Color("#cccccc")
	Gray900 = lipgloss.Color("#e6e6e6")

	// Red
	Red100 = lipgloss.Color("#330900")
	Red200 = lipgloss.Color("#661300")
	Red300 = lipgloss.Color("#991c00")
	Red400 = lipgloss.Color("#cc2500")
	Red500 = lipgloss.Color("#ff2f00")
	Red600 = lipgloss.Color("#ff5833")
	Red700 = lipgloss.Color("#ff8266")
	Red800 = lipgloss.Color("#ffac99")
	Red900 = lipgloss.Color("#ffd5cc")

	// Orange (WSO2 Primary)
	Orange100 = lipgloss.Color("#331700")
	Orange200 = lipgloss.Color("#662e00")
	Orange300 = lipgloss.Color("#994500")
	Orange400 = lipgloss.Color("#cc5c00")
	Orange500 = lipgloss.Color("#ff7300")
	Orange600 = lipgloss.Color("#ff8f33")
	Orange700 = lipgloss.Color("#ffab66")
	Orange800 = lipgloss.Color("#ffc799")
	Orange900 = lipgloss.Color("#ffe3cc")

	// Green
	Green100 = lipgloss.Color("#0c271a")
	Green200 = lipgloss.Color("#174f33")
	Green300 = lipgloss.Color("#23764d")
	Green400 = lipgloss.Color("#2f9d66")
	Green500 = lipgloss.Color("#3bc480")
	Green600 = lipgloss.Color("#62d099")
	Green700 = lipgloss.Color("#89dcb3")
	Green800 = lipgloss.Color("#b0e8cc")
	Green900 = lipgloss.Color("#d8f3e6")

	// Indigo
	Indigo100 = lipgloss.Color("#0a0e29")
	Indigo200 = lipgloss.Color("#141d52")
	Indigo300 = lipgloss.Color("#1f2b7a")
	Indigo400 = lipgloss.Color("#2939a3")
	Indigo500 = lipgloss.Color("#3347cc")
	Indigo600 = lipgloss.Color("#5c6cd6")
	Indigo700 = lipgloss.Color("#8591e0")
	Indigo800 = lipgloss.Color("#adb6eb")
	Indigo900 = lipgloss.Color("#d6daf5")

	// Blue (reversed scale - lighter to darker)
	Blue100 = lipgloss.Color("#cceaff")
	Blue200 = lipgloss.Color("#99d5ff")
	Blue300 = lipgloss.Color("#66bfff")
	Blue400 = lipgloss.Color("#33aaff")
	Blue500 = lipgloss.Color("#0095ff")
	Blue600 = lipgloss.Color("#0077cc")
	Blue700 = lipgloss.Color("#005999")
	Blue800 = lipgloss.Color("#003c66")
	Blue900 = lipgloss.Color("#001e33")

	// Yellow
	Yellow100 = lipgloss.Color("#332200")
	Yellow200 = lipgloss.Color("#664400")
	Yellow300 = lipgloss.Color("#996600")
	Yellow400 = lipgloss.Color("#cc8800")
	Yellow500 = lipgloss.Color("#ffaa00")
	Yellow600 = lipgloss.Color("#ffbb33")
	Yellow700 = lipgloss.Color("#ffcc66")
	Yellow800 = lipgloss.Color("#ffdd99")
	Yellow900 = lipgloss.Color("#ffeecc")

	// Teal (reversed scale - lighter to darker)
	Teal100 = lipgloss.Color("#d1faf5")
	Teal200 = lipgloss.Color("#a3f5eb")
	Teal300 = lipgloss.Color("#75f0e1")
	Teal400 = lipgloss.Color("#47ebd8")
	Teal500 = lipgloss.Color("#19e6ce")
	Teal600 = lipgloss.Color("#14b8a5")
	Teal700 = lipgloss.Color("#0f8a7b")
	Teal800 = lipgloss.Color("#0a5c52")
	Teal900 = lipgloss.Color("#052e29")
)

// Semantic color aliases for easy usage
var (
	Primary   = Orange500
	Secondary = Indigo500
	Success   = Green500
	Error     = Red500
	Warning   = Yellow500
	Info      = Blue500
	Accent    = Teal500
)
