package Locale

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Locale struct {
	messages map[string]string
}

func (l *Locale) T(key string, args ...any) string {
	tmpl, ok := l.messages[key]
	if !ok {
		return key
	}
	if len(args) == 0 {
		return tmpl
	}
	return fmt.Sprintf(tmpl, args...)
}

func Load(localesDir string, lang string) (*Locale, error) {
	if lang == "" {
		lang = detectSystemLang()
	}
	path := fmt.Sprintf("%s/%s.json", localesDir, lang)
	data, err := os.ReadFile(path)
	if err != nil {
		// fallback to English
		path = fmt.Sprintf("%s/en.json", localesDir)
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, err
		}
	}
	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}
	return &Locale{messages: messages}, nil
}

func detectSystemLang() string {
	for _, env := range []string{"LANG", "LANGUAGE", "LC_ALL", "LC_MESSAGES"} {
		if val := os.Getenv(env); val != "" {
			tag := strings.ToLower(val)
			if strings.HasPrefix(tag, "ru") {
				return "ru"
			}
			return "en"
		}
	}
	// Windows fallback
	return detectWindowsLang()
}
