package Locale

import (
	"os/exec"
	"strings"
)

func detectWindowsLang() string {
	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		"(Get-Culture).TwoLetterISOLanguageName").Output()
	if err != nil {
		return "en"
	}
	lang := strings.TrimSpace(strings.ToLower(string(out)))
	if lang == "ru" {
		return "ru"
	}
	return "en"
}
