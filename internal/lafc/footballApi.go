package lafc

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ASullivan219/freeSandwich/internal/lafc/models"
	"github.com/ASullivan219/freeSandwich/internal/notifier"
)

const (
	BASE_API_URL     = "https://api-football-v1.p.rapidapi.com/v3"
	LEAGUE_ID        = "253"
	LAFC_TEAM_ID     = "1616"
	LAFC_TEAM_ID_INT = 1616
	SEASON           = "2024"
	TIME_LAYOUT      = "2006-01-02T15:04:05+00:00"
)

type simpleFixture struct {
	fixtureId     int
	Time          time.Time
	ExecutionTime time.Time
	Description   string
}

func newSimpleFixture(f models.FixtureEntry) simpleFixture {
	fixtureDate, err := time.Parse(TIME_LAYOUT, f.Fixture.Date)
	threeHours, _ := time.ParseDuration("3h")
	// Make the execution time 3 hours after game start to make sure that the game has ended
	// Before evaluating for win conditions
	executionTime := fixtureDate.Add(threeHours)
	description := fmt.Sprintf("Home: %s, Away: %s, DateTime: %s, ExecutionTime: %s",
		f.Teams.Home.Name, f.Teams.Away.Name, f.Fixture.Date, executionTime)
	if err != nil {
		slog.Error(
			"bad date",
			slog.String("error", err.Error()),
		)
	}
	return simpleFixture{fixtureId: f.Fixture.Id, Time: fixtureDate, Description: description, ExecutionTime: executionTime}
}

func BuildCronnables(notifier notifier.I_Notifier) []LafcCronJob {
	fixtures := getAllFixtures()
	eligible := filterFixtures(fixtures)
	return convertToCronnables(eligible, notifier)
}

func convertToCronnables(fixtures []simpleFixture, notifier notifier.I_Notifier) []LafcCronJob {
	jobs := make([]LafcCronJob, 0, len(fixtures))
	for _, fixture := range fixtures {
		jobs = append(jobs, *newjob(fixture, notifier))
	}
	return jobs
}

func filterFixtures(fixtures []models.FixtureEntry) []simpleFixture {
	eligibleFixtures := make([]simpleFixture, 0, len(fixtures))
	currentTime := time.Now()
	for _, fixture := range fixtures {
		fixtureDate, err := time.Parse(TIME_LAYOUT, fixture.Fixture.Date)
		if err != nil {
			slog.Error(
				"error parsing date",
				slog.String("date", fixture.Fixture.Date),
				slog.String("error", err.Error()),
			)
		}
		if filter(fixture, fixtureDate, currentTime) {
			eligibleFixtures = append(
				eligibleFixtures,
				newSimpleFixture(fixture))
		}
	}
	return eligibleFixtures
}

// Return True if The fixture should be checked for win conditions
// Only Home games should be checked and only games in the future
func filter(fixture models.FixtureEntry, fixtureDate time.Time, currentTime time.Time) bool {

	threehours, _ := time.ParseDuration("3h")
	if fixture.Teams.Home.Id != LAFC_TEAM_ID_INT {
		return false
	}

	if fixtureDate.UTC().Add(threehours).Before(currentTime.UTC()) {
		return false
	}

	return true
}

func getAllFixtures() []models.FixtureEntry {

	u, err := url.Parse(BASE_API_URL)
	if err != nil {
		slog.Error(
			"error parsing url",
			slog.String("error", err.Error()),
		)
	}
	u = u.JoinPath("fixtures")
	queryParams := u.Query()
	queryParams.Set("league", LEAGUE_ID)
	queryParams.Set("season", SEASON)
	queryParams.Set("team", LAFC_TEAM_ID)
	u.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		slog.Error(
			"error bad request",
			slog.String("error", err.Error()))
	}
	req.Header.Add("X-RapidAPI-Key", os.Getenv("FOOTBALL_RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("FOOTBALL_RAPIS_API_HOST"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error(
			"error performing request",
			slog.String("error", err.Error()))
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var fixtureResponse models.FixtureResponse
	err = json.Unmarshal(body, &fixtureResponse)
	if err != nil {
		slog.Error(
			"json error",
			slog.String("error", err.Error()))
	}

	return fixtureResponse.Response

}
