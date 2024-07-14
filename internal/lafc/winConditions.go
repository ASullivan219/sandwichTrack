package lafc

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/ASullivan219/freeSandwich/internal/lafc/models"
	"github.com/ASullivan219/freeSandwich/internal/notifier"
)

type LafcCronJob struct {
	Sf       simpleFixture
	notifier notifier.I_Notifier
}

func newjob(sf simpleFixture, n notifier.I_Notifier) *LafcCronJob {
	return &LafcCronJob{Sf: sf, notifier: n}
}

func (j *LafcCronJob) Run() {
	notify := CheckWinConditions(j.Sf)
	if notify {
		j.notifier.NotifyAll(
			fmt.Sprintf(
				"LAFC won at home, check the chick-fil-a app for a free sandwich\n%s",
				j.Sf.Description,
			),
		)
	}
}

func CheckWinConditions(fixture simpleFixture) bool {
	fixtureId := strconv.Itoa(fixture.fixtureId)
	api, err := url.Parse(BASE_API_URL)
	if err != nil {
		slog.Error("Error parsing base url", slog.String("error", err.Error()))
		return false
	}
	api = api.JoinPath("fixtures")
	params := api.Query()
	params.Set("id", fixtureId)
	api.RawQuery = params.Encode()
	req, err := http.NewRequest("GET", api.String(), nil)
	if err != nil {
		slog.Error("bad request", slog.String("error", err.Error()))
		return false
	}
	req.Header.Add("X-RapidAPI-Key", os.Getenv("FOOTBALL_RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("FOOTBALL_RAPIS_API_HOST"))
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		slog.Error("error performing request", slog.String("error", err.Error()))
		return false
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var fixtureResponse models.FixtureResponse
	err = json.Unmarshal(body, &fixtureResponse)

	slog.Info(
		"evaluating game",
		slog.String("home", fixtureResponse.Response[0].Teams.Home.Name),
		slog.String("away", fixtureResponse.Response[0].Teams.Away.Name),
		slog.String("date", fixtureResponse.Response[0].Fixture.Date),
	)

	return fixtureResponse.Response[0].Teams.Home.Winner == true
}
