package models

type Fixture struct {
	Id       int
	Referee  string
	Timezone string
	Date     string
}

type Team struct {
	Id     int
	Name   string
	Logo   string
	Winner any
}

type Teams struct {
	Home Team
	Away Team
}

type FixtureEntry struct {
	Fixture Fixture
	League  any
	Teams   Teams
	Goals   any
}

type FixtureResponse struct {
	Get        string
	Parameters any
	Errors     []any
	Results    int
	Paging     any
	Response   []FixtureEntry
}
