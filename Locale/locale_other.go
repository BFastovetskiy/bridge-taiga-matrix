//go:build !windows

package Locale

func detectWindowsLang() string {
	return "en"
}
