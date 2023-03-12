package steamdota

type Steamdota struct {
	Result Result `json:"result,omitempty"`
}
type Players struct {
	AccountID  int `json:"account_id,omitempty"`
	PlayerSlot int `json:"player_slot,omitempty"`
	TeamNumber int `json:"team_number,omitempty"`
	TeamSlot   int `json:"team_slot,omitempty"`
	HeroID     int `json:"hero_id,omitempty"`
}
type Matches struct {
	SeriesID      int       `json:"series_id,omitempty"`
	SeriesType    int       `json:"series_type,omitempty"`
	MatchID       int64     `json:"match_id,omitempty"`
	MatchSeqNum   int64     `json:"match_seq_num,omitempty"`
	StartTime     int       `json:"start_time,omitempty"`
	LobbyType     int       `json:"lobby_type,omitempty"`
	RadiantTeamID int       `json:"radiant_team_id,omitempty"`
	DireTeamID    int       `json:"dire_team_id,omitempty"`
	Players       []Players `json:"players,omitempty"`
}
type Result struct {
	Status           int       `json:"status,omitempty"`
	NumResults       int       `json:"num_results,omitempty"`
	TotalResults     int       `json:"total_results,omitempty"`
	ResultsRemaining int       `json:"results_remaining,omitempty"`
	Matches          []Matches `json:"matches,omitempty"`
}

type Hero []struct {
	ID              int         `json:"id"`
	Name            string      `json:"name"`
	LocalizedName   string      `json:"localized_name"`
	PrimaryAttr     string      `json:"primary_attr"`
	AttackType      string      `json:"attack_type"`
	Roles           []string    `json:"roles"`
	Img             string      `json:"img"`
	Icon            string      `json:"icon"`
	BaseHealth      int         `json:"base_health"`
	BaseHealthRegen float64     `json:"base_health_regen"`
	BaseMana        int         `json:"base_mana"`
	BaseManaRegen   int         `json:"base_mana_regen"`
	BaseArmor       int         `json:"base_armor"`
	BaseMr          int         `json:"base_mr"`
	BaseAttackMin   int         `json:"base_attack_min"`
	BaseAttackMax   int         `json:"base_attack_max"`
	BaseStr         int         `json:"base_str"`
	BaseAgi         int         `json:"base_agi"`
	BaseInt         int         `json:"base_int"`
	StrGain         float64     `json:"str_gain"`
	AgiGain         float64     `json:"agi_gain"`
	IntGain         float64     `json:"int_gain"`
	AttackRange     int         `json:"attack_range"`
	ProjectileSpeed int         `json:"projectile_speed"`
	AttackRate      float64     `json:"attack_rate"`
	BaseAttackTime  int         `json:"base_attack_time"`
	AttackPoint     float64     `json:"attack_point"`
	MoveSpeed       int         `json:"move_speed"`
	TurnRate        interface{} `json:"turn_rate"`
	CmEnabled       bool        `json:"cm_enabled"`
	Legs            int         `json:"legs"`
	DayVision       int         `json:"day_vision"`
	NightVision     int         `json:"night_vision"`
	HeroID          int         `json:"hero_id"`
	TurboPicks      int         `json:"turbo_picks"`
	TurboWins       int         `json:"turbo_wins"`
	ProBan          int         `json:"pro_ban"`
	ProWin          int         `json:"pro_win"`
	ProPick         int         `json:"pro_pick"`
	OnePick         int         `json:"1_pick"`
	OneWin          int         `json:"1_win"`
	TwoPick         int         `json:"2_pick"`
	TwoWin          int         `json:"2_win"`
	ThreePick       int         `json:"3_pick"`
	ThreeWin        int         `json:"3_win"`
	FourPick        int         `json:"4_pick"`
	FourWin         int         `json:"4_win"`
	FivePick        int         `json:"5_pick"`
	FiveWin         int         `json:"5_win"`
	SixPick         int         `json:"6_pick"`
	SixWin          int         `json:"6_win"`
	SevenPick       int         `json:"7_pick"`
	SevenWin        int         `json:"7_win"`
	EightPick       int         `json:"8_pick"`
	EightWin        int         `json:"8_win"`
	NullPick        int         `json:"null_pick"`
	NullWin         int         `json:"null_win"`
}
