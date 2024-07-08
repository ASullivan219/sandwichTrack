package lafc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ASullivan219/freeSandwich/internal/lafc/models"
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
		fmt.Println("ERROR: parsing date")
	}
	return simpleFixture{fixtureId: f.Fixture.Id, Time: fixtureDate, Description: description, ExecutionTime: executionTime}
}

func BuildCronnables() []LafcCronJob {
	fixtures := getAllFixtures()
	eligible := filterFixtures(fixtures)
	fmt.Println(len(eligible))
	return convertToCronnables(eligible)
}

func convertToCronnables(fixtures []simpleFixture) []LafcCronJob {
	jobs := make([]LafcCronJob, 0, len(fixtures))
	for _, fixture := range fixtures {
		jobs = append(jobs, *newjob(fixture))
	}
	return jobs
}

func filterFixtures(fixtures []models.FixtureEntry) []simpleFixture {
	eligibleFixtures := make([]simpleFixture, 0, len(fixtures))
	currentTime := time.Now()
	for _, fixture := range fixtures {
		fixtureDate, err := time.Parse(TIME_LAYOUT, fixture.Fixture.Date)
		if err != nil {
			fmt.Printf("error parsing date: %s\n", fixture.Fixture.Date)
			fmt.Println(err)
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

	if fixture.Teams.Home.Id != LAFC_TEAM_ID_INT {
		return false
	}

	if fixtureDate.UTC().Before(currentTime.UTC()) {
		return false
	}

	return true
}

func getAllFixtures() []models.FixtureEntry {

	u, err := url.Parse(BASE_API_URL)
	if err != nil {
		fmt.Println("error parsing base url")
	}
	u = u.JoinPath("fixtures")
	queryParams := u.Query()
	queryParams.Set("league", LEAGUE_ID)
	queryParams.Set("season", SEASON)
	queryParams.Set("team", LAFC_TEAM_ID)
	u.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Println("Bad request", err)
	}
	req.Header.Add("X-RapidAPI-Key", os.Getenv("FOOTBALL_RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("FOOTBALL_RAPIS_API_HOST"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error performing request", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var fixtureResponse models.FixtureResponse
	err = json.Unmarshal(body, &fixtureResponse)
	if err != nil {
		fmt.Println("Error unmarshaling fixture response", err)
	}

	return fixtureResponse.Response

}