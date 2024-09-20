package kungdota

const ChatMessageFirstBlood = "CHAT_MESSAGE_FIRSTBLOOD"

type Match struct {
	Teams             [][]string `json:"teams"`
	Score             []int      `json:"score"`
	Winner            string     `json:"winner"`
	DotaMatchId       string     `json:"dotaMatchId"`
	LeagueId          string     `json:"leagueId"`
	DiedFirstBlood    string     `json:"diedFirstBlood"`
	ClaimedFirstBlood string     `json:"claimedFirstBlood"`

	CoolaStats []CoolaStats `json:"coolaStats"`
}
type CoolaStats struct {
	ObserverKills string `json:"observer_kills"`
	SentryKills   string `json:"sentry_kills"`
	ObsPlaced     string `json:"obs_placed"`
	SenPlaced     string `json:"sen_placed"`
	Kills         string `json:"kills"`
	Assists       string `json:"assists"`
	Deaths        string `json:"deaths"`
	FantasyPoints string `json:"fantasyPoints"`
}

type ShuffledTeams struct {
	TeamOne      Players
	TeamTwo      Players
	FirstPicker  Player
	SecondPicker Player
	EloDiff      int
}

// TODO
func (rx Players) Names() []string {
	r := make([]string, 0)

	for _, p := range rx.Players {
		r = append(r, p.Username)
	}

	return r
}

type Players struct {
	Players []Player `json:"players"`
}

type Player struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	EloRating       int    `json:"eloRating"`
	Steam32ID       string `json:"steam32id,omitempty"`
	DiscordID       string `json:"discordId,omitempty"`
	DiscordUsername string `json:"discordUsername,omitempty"`
}
