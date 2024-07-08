package lafc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/ASullivan219/freeSandwich/internal/lafc/models"
)

type LafcCronJob struct {
	Sf simpleFixture
}

func newjob(sf simpleFixture) *LafcCronJob {
	return &LafcCronJob{Sf: sf}
}

func (j *LafcCronJob) Run() {
	CheckWinConditions(j.Sf)
}

func CheckWinConditions(fixture simpleFixture) {
	fixtureId := strconv.Itoa(fixture.fixtureId)
	api, err := url.Parse(BASE_API_URL)
	if err != nil {
		fmt.Println("Error parsing baseURL, abandoning checking win conditions")
		return
	}
	api = api.JoinPath("fixtures")
	params := api.Query()
	params.Set("id", fixtureId)
	api.RawQuery = params.Encode()
	req, err := http.NewRequest("GET", api.String(), nil)
	if err != nil {
		fmt.Println("Bad request, abandoning checking win conditions for", err)
	}
	req.Header.Add("X-RapidAPI-Key", os.Getenv("FOOTBALL_RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("FOOTBALL_RAPIS_API_HOST"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error performing request abandoning checking win conditions", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var fixtureResponse models.FixtureResponse
	err = json.Unmarshal(body, &fixtureResponse)
	fmt.Printf("Evaluating %s v %s: Time: %s\n", fixtureResponse.Response[0].Teams.Home.Name, fixtureResponse.Response[0].Teams.Away.Name, fixtureResponse.Response[0].Fixture.Date)

	if fixtureResponse.Response[0].Teams.Home.Winner == true {
		notify()
	}
}

func notify() {
	fmt.Println("Notification of a win")
}
