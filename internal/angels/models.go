package angels

import (
	"fmt"
	"time"
)

type score struct {
	Hits  int
	Total int
}

type scores struct {
	Home score
	Away score
}

type team struct {
	Id   int
	Name string
}

type teams struct {
	Home team
	Away team
}

type status struct {
	Long  string
	Short string
}

type game struct {
	Id     int
	Date   string
	Status status
	Teams  teams
	Scores scores
}

type apiResponse struct {
	Get      string
	Results  int
	Response []game
}

func (g *game) toSimplifiedGame() (simplifiedGame, error) {
	gameLength, err := time.ParseDuration("4h")
	gameTime, err := time.Parse(TIME_LAYOUT, g.Date)
	if err != nil {
		return simplifiedGame{}, err

	}

	execTime := gameTime.Add(gameLength)
	description := fmt.Sprintf("Home: %s Away: %s gameTime: %s execTime: %s", g.Teams.Home.Name, g.Teams.Away.Name, g.Date, execTime)

	return simplifiedGame{
		Id:            g.Id,
		Date:          gameTime,
		ExecutionTime: execTime,
		HomeScore:     g.Scores.Home.Total,
		AwayScore:     g.Scores.Away.Total,
		AwayTeam:      g.Teams.Away.Name,
		Description:   description,
	}, nil
}

type simplifiedGame struct {
	Id            int
	Date          time.Time
	ExecutionTime time.Time
	HomeScore     int
	AwayScore     int
	AwayTeam      string
	Description   string
}

func (sg *simplifiedGame) toString() string {
	return fmt.Sprintf("GameId: %d, Date: %s, ExecutionTime: %s, Away team: %s, Home: %d, Away: %d",
		sg.Id,
		sg.Date,
		sg.ExecutionTime,
		sg.AwayTeam,
		sg.HomeScore,
		sg.AwayScore,
	)
}
