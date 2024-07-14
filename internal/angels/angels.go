package angels

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ASullivan219/freeSandwich/internal/cronmanager"
	"github.com/ASullivan219/freeSandwich/internal/notifier"
)

const (
	BASE_API_URL       = "https://api-baseball.p.rapidapi.com/"
	MLB_LEAGUE_ID      = "1"
	ANGELS_TEAM_ID     = "17"
	ANGELS_TEAM_ID_INT = 17
	SEASON             = "2024"
	TIME_LAYOUT        = "2006-01-02T15:04:05+00:00"
)

func BuildCronnables(n notifier.I_Notifier, cm cronmanager.CronManager) []*AngelsCronJob {
	allGames := ListAngelsGames()
	eligibleGames := filterGames(allGames)
	simplifiedGames := simplifyAllGames(eligibleGames)
	var angelsJobs []*AngelsCronJob
	for _, sg := range simplifiedGames {
		angelsJobs = append(angelsJobs, newAngelsJob(sg, n, cm))
	}
	return angelsJobs
}

func ListAngelsGames() []game {
	u, err := url.Parse(BASE_API_URL)
	if err != nil {
		slog.Error(
			"error parsing base URL",
			"err", err.Error(),
		)
	}
	u = u.JoinPath("games")
	queryParams := u.Query()
	queryParams.Set("team", "17")
	queryParams.Set("league", "1")
	queryParams.Set("season", "2024")
	u.RawQuery = queryParams.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		slog.Error(
			"bad request",
			"error", err.Error(),
		)
	}
	req.Header.Add("X-RapidAPI-Key", os.Getenv("BASEBALL_RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("BASEBALL_RAPID_API_HOST"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error(
			"error performing request",
			"error", err.Error(),
		)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var apiResponse apiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		slog.Error(
			"bad json",
			"error", err.Error(),
		)
	}
	return apiResponse.Response
}

func simplifyAllGames(games []game) []simplifiedGame {
	simplified := make([]simplifiedGame, 0)
	for _, game := range games {
		simpleGame, err := game.toSimplifiedGame()
		if err != nil {
			slog.Error("error simplifiying game",
				"error", err.Error(),
			)
			continue
		}
		simplified = append(simplified, simpleGame)
	}
	return simplified
}

func filterGames(games []game) []game {
	eligible := make([]game, 0)
	currentTime := time.Now().UTC()
	for _, game := range games {
		if !angelsHome(game) {
			continue
		}
		if !afterToday(game, currentTime) {
			continue
		}
		eligible = append(eligible, game)
	}
	return eligible
}

func angelsHome(game game) bool {
	return game.Teams.Home.Id == ANGELS_TEAM_ID_INT
}

func afterToday(game game, currentTime time.Time) bool {
	delay, _ := time.ParseDuration("3h")
	gameDate, err := time.Parse(TIME_LAYOUT, game.Date)
	if err != nil {
		slog.Error(
			"error parsing date",
			slog.String("error", err.Error()),
		)
		return false
	}

	return gameDate.Add(delay).After(currentTime)
}
