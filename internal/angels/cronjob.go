package angels

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ASullivan219/freeSandwich/internal/cronmanager"
	"github.com/ASullivan219/freeSandwich/internal/notifier"
)

var RESCHEDULE_ERROR = errors.New("RESCHEDULE")

func newAngelsJob(sg simplifiedGame, n notifier.I_Notifier, cm cronmanager.CronManager) *AngelsCronJob {
	return &AngelsCronJob{
		SimpleGame:    sg,
		ExecutionTime: sg.ExecutionTime,
		CronMan:       cm,
		notifier:      n,
	}
}

type AngelsCronJob struct {
	SimpleGame    simplifiedGame
	ExecutionTime time.Time
	CronMan       cronmanager.CronManager
	notifier      notifier.I_Notifier
}

func (j *AngelsCronJob) Run() {
	slog.Info("Checking game" + j.SimpleGame.Description)
	notify, err := CheckWinConditions(j.SimpleGame.Id)
	if err != nil {
		tenMinutes, _ := time.ParseDuration("10m")
		rescheduleTime := time.Now().Add(tenMinutes)
		slog.Error(
			"Rescheduling check",
			slog.String("description", j.SimpleGame.Description),
			slog.String("reschedule time", rescheduleTime.String()))

		j.CronMan.AddJob(rescheduleTime, j, j.SimpleGame.toString())
		return
	}

	if notify {
		j.notifier.NotifyAll(
			fmt.Sprintf("Angels Scored 7 or more Runs at home Check the Chick-Fil-A app\n%s", j.SimpleGame.toString()),
		)
		return
	}
	slog.Info(
		"Win conditions not met:",
		slog.String("game", j.SimpleGame.Description))
}

func CheckWinConditions(gameId int) (bool, error) {
	intId := strconv.Itoa(gameId)
	u, err := url.Parse(BASE_API_URL)
	if err != nil {
		return false, err
	}
	u = u.JoinPath("games")
	queryParams := u.Query()
	queryParams.Set("id", intId)
	u.RawQuery = queryParams.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("X-RapidAPI-Key", os.Getenv("BASEBALL_RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("BASEBALL_RAPID_API_HOST"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var apiResponse apiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return false, err
	}

	if apiResponse.Response[0].Status.Short != "FT" {
		return false, RESCHEDULE_ERROR
	}

	if apiResponse.Response[0].Scores.Home.Total < 7 {
		return false, nil
	}

	return true, nil
}
