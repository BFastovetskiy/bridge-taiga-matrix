package Config

import (
	"encoding/json"
	"os"
)

type TaigaProject struct {
	Name              string `json:"name"`
	MatrixProjectRoomID string `json:"matrixProjectRoomID"`
}

type Config struct {
	TaigaBaseURL  string         `json:"taigaBaseURL"`
	TaigaUsername string         `json:"taigaUsername"`
	TaigaPassword string         `json:"taigaPassword"`
	TaigaProjects []TaigaProject `json:"taigaProjects"`

	MatrixServer            string `json:"matrixServer"`
	MatrixToken             string `json:"matrixToken"`
	DuplicateToGeneralGroup bool   `json:"duplicateToGeneralGroup"`
	GeneralRoomID           string `json:"generalRoomId"`

	InsecureSkipVerify bool `json:"InsecureSkipVerify"`

	Language string `json:"language"`

	DaysUntilDeadline int `json:"daysUntilDeadline"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
