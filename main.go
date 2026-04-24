package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"bridge-taiga-matrix/Config"
	"bridge-taiga-matrix/Locale"
)

var cfg *Config.Config
var apiClient *http.Client
var locale *Locale.Locale

type TaigaItem struct {
	Subject   string `json:"subject"`
	DueDate   string `json:"due_date"`
	IsClosed  bool   `json:"is_closed"`
	IsBlocked bool   `json:"is_blocked"`
}

type AuthResponse struct {
	AuthToken string `json:"auth_token"`
}

type Project struct {
	ID int `json:"id"`
}

func main() {
	configPath := flag.String("config", "settings.json", "path to configuration JSON file")
	flag.Parse()

	var err error
	cfg, err = Config.Load(*configPath)
	if err != nil {
		slog.Error("Error load configuration", "error", err)
		return
	}

	locale, err = Locale.Load("locales", cfg.Language)
	if err != nil {
		slog.Error("Error load locales", "error", err)
		return
	}

	apiClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify},
		},
	}

	token := login()
	if token == "" {
		slog.Error("Authorization error in Taiga")
		return
	}

	now := time.Now()
	for _, project := range cfg.TaigaProjects {
		projectID := getProjectID(token, project.Name)

		stories := getItems(fmt.Sprintf("%s/api/v1/userstories?project=%d", cfg.TaigaBaseURL, projectID), token)
		checkOverdue(stories, locale.T("userstory"), now, project.MatrixProjectRoomID)

		tasks := getItems(fmt.Sprintf("%s/api/v1/tasks?project=%d", cfg.TaigaBaseURL, projectID), token)
		checkOverdue(tasks, locale.T("task"), now, project.MatrixProjectRoomID)
	}
}

func checkOverdue(items []TaigaItem, label string, now time.Time, roomID string) {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for _, item := range items {
		if !item.IsClosed {
			if item.IsBlocked {
				continue
			}
			if item.DueDate == "" {
				sendToMatrix(locale.T("no_deadline", item.Subject), roomID)
				continue
			}
			due, err := time.Parse("2006-01-02", item.DueDate)
			if err != nil {
				continue
			}

			if due.Before(today) {
				sendToMatrix(locale.T("overdue", label, item.Subject, item.DueDate), roomID)
				continue
			}

			daysLeft := int(due.Sub(today).Hours() / 24)
			if daysLeft <= cfg.DaysUntilDeadline {
				sendToMatrix(locale.T("days_left", label, item.Subject, cfg.DaysUntilDeadline, item.DueDate), roomID)
			}
		}
	}
}

func getItems(url string, token string) []TaigaItem {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := apiClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var items []TaigaItem
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &items)
	return items
}

func login() string {
	url := fmt.Sprintf("%s/api/v1/auth", cfg.TaigaBaseURL)
	payload := map[string]string{"type": "normal", "username": cfg.TaigaUsername, "password": cfg.TaigaPassword}
	body, _ := json.Marshal(payload)

	resp, err := apiClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != 200 {
		return ""
	}
	defer resp.Body.Close()

	var auth AuthResponse
	json.NewDecoder(resp.Body).Decode(&auth)
	return auth.AuthToken
}

func getProjectID(token string, projectName string) int {
	url := fmt.Sprintf("%s/api/v1/projects/by_slug?slug=%s", cfg.TaigaBaseURL, projectName)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := apiClient.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	var p Project
	json.NewDecoder(resp.Body).Decode(&p)
	return p.ID
}

func sendToMatrix(text string, roomID string) {
	send := func(room string) {
		url := fmt.Sprintf("%s/_matrix/client/r0/rooms/%s/send/m.room.message?access_token=%s",
			cfg.MatrixServer, room, cfg.MatrixToken)
		payload := map[string]string{"msgtype": "m.text", "body": text}
		body, _ := json.Marshal(payload)
		apiClient.Post(url, "application/json", bytes.NewBuffer(body))
	}

	send(roomID)
	if cfg.DuplicateToGeneralGroup && roomID != cfg.GeneralRoomID {
		send(cfg.GeneralRoomID)
	}
	slog.Info(locale.T("sent", text))
}
