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
		j.notifier.NotifyAll()
	}
}

func CheckWinConditions(fixture simpleFixture) bool {
	fixtureId := strconv.Itoa(fixture.fixtureId)
	api, err := url.Parse(BASE_API_URL)
	if err != nil {
		fmt.Println("Error parsing baseURL, abandoning checking win conditions")
		return false
	}
	api = api.JoinPath("fixtures")
	params := api.Query()
	params.Set("id", fixtureId)
	api.RawQuery = params.Encode()
	req, err := http.NewRequest("GET", api.String(), nil)
	if err != nil {
		fmt.Println("Bad request, abandoning checking win conditions for", err)
		return false
	}
	req.Header.Add("X-RapidAPI-Key", os.Getenv("FOOTBALL_RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("FOOTBALL_RAPIS_API_HOST"))
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("error performing request abandoning checking win conditions", err)
		return false
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var fixtureResponse models.FixtureResponse
	err = json.Unmarshal(body, &fixtureResponse)

	fmt.Printf("Evaluating %s v %s: Time: %s\n",
		fixtureResponse.Response[0].Teams.Home.Name,
		fixtureResponse.Response[0].Teams.Away.Name,
		fixtureResponse.Response[0].Fixture.Date,
	)

	return fixtureResponse.Response[0].Teams.Home.Winner == true
}
