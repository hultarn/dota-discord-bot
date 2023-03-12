package opendota

import "time"

type OpenDotaGameObject struct {
	MatchID               int            `json:"match_id"`
	BarracksStatusDire    int            `json:"barracks_status_dire"`
	BarracksStatusRadiant int            `json:"barracks_status_radiant"`
	Chat                  []Chat         `json:"chat"`
	Cluster               int            `json:"cluster"`
	Cosmetics             Cosmetics      `json:"cosmetics"`
	DireScore             int            `json:"dire_score"`
	DraftTimings          []DraftTimings `json:"draft_timings"`
	Duration              int            `json:"duration"`
	Engine                int            `json:"engine"`
	FirstBloodTime        int            `json:"first_blood_time"`
	GameMode              int            `json:"game_mode"`
	HumanPlayers          int            `json:"human_players"`
	Leagueid              int            `json:"leagueid"`
	LobbyType             int            `json:"lobby_type"`
	MatchSeqNum           int            `json:"match_seq_num"`
	NegativeVotes         int            `json:"negative_votes"`
	Objectives            []Objectives   `json:"objectives"`
	PicksBans             PicksBans      `json:"picks_bans"`
	PositiveVotes         int            `json:"positive_votes"`
	RadiantGoldAdv        RadiantGoldAdv `json:"radiant_gold_adv"`
	RadiantScore          int            `json:"radiant_score"`
	RadiantWin            bool           `json:"radiant_win"`
	RadiantXpAdv          RadiantXpAdv   `json:"radiant_xp_adv"`
	StartTime             int            `json:"start_time"`
	Teamfights            Teamfights     `json:"teamfights"`
	TowerStatusDire       int            `json:"tower_status_dire"`
	TowerStatusRadiant    int            `json:"tower_status_radiant"`
	Version               int            `json:"version"`
	ReplaySalt            int            `json:"replay_salt"`
	SeriesID              int            `json:"series_id"`
	SeriesType            int            `json:"series_type"`
	RadiantTeam           RadiantTeam    `json:"radiant_team"`
	DireTeam              DireTeam       `json:"dire_team"`
	League                League         `json:"league"`
	Skill                 int            `json:"skill"`
	Players               []Players      `json:"players"`
	Patch                 int            `json:"patch"`
	Region                int            `json:"region"`
	AllWordCounts         AllWordCounts  `json:"all_word_counts"`
	MyWordCounts          MyWordCounts   `json:"my_word_counts"`
	Throw                 int            `json:"throw"`
	Comeback              int            `json:"comeback"`
	Loss                  int            `json:"loss"`
	Win                   int            `json:"win"`
	ReplayURL             string         `json:"replay_url"`
}
type Chat struct {
	Time       int    `json:"time"`
	Unit       string `json:"unit"`
	Key        string `json:"key"`
	Slot       int    `json:"slot"`
	PlayerSlot int    `json:"player_slot"`
}
type Cosmetics struct {
}
type DraftTimings struct {
	Order          int  `json:"order"`
	Pick           bool `json:"pick"`
	ActiveTeam     int  `json:"active_team"`
	HeroID         int  `json:"hero_id"`
	PlayerSlot     int  `json:"player_slot"`
	ExtraTime      int  `json:"extra_time"`
	TotalTimeTaken int  `json:"total_time_taken"`
}
type Objectives struct {
	Time       int    `json:"time"`
	Type       string `json:"type"`
	Slot       int    `json:"slot,omitempty"`
	Key        int    `json:"key,omitempty"`
	PlayerSlot int    `json:"player_slot,omitempty"`
	Unit       string `json:"unit,omitempty"`
	Value      int    `json:"value,omitempty"`
	Killer     int    `json:"killer,omitempty"`
	Team       int    `json:"team,omitempty"`
}
type PicksBans struct {
}
type RadiantGoldAdv struct {
}
type RadiantXpAdv struct {
}
type Teamfights struct {
}
type RadiantTeam struct {
}
type DireTeam struct {
}
type League struct {
}
type AbilityUses struct {
}
type AbilityTargets struct {
}
type DamageTargets struct {
}
type Actions struct {
}
type AdditionalUnits struct {
}
type BuybackLog struct {
	Time       int `json:"time"`
	Slot       int `json:"slot"`
	PlayerSlot int `json:"player_slot"`
}
type ConnectionLog struct {
	Time       int    `json:"time"`
	Event      string `json:"event"`
	PlayerSlot int    `json:"player_slot"`
}
type Damage struct {
}
type DamageInflictor struct {
}
type DamageInflictorReceived struct {
}
type DamageTaken struct {
}
type GoldReasons struct {
}
type HeroHits struct {
}
type ItemUses struct {
}
type KillStreaks struct {
}
type Killed struct {
}
type KilledBy struct {
}
type KillsLog struct {
	Time int    `json:"time"`
	Key  string `json:"key"`
}
type LanePos struct {
}
type LifeState struct {
}
type MaxHeroHit struct {
}
type MultiKills struct {
}
type Obs struct {
}
type ObsLeftLog struct {
}
type ObsLog struct {
}
type PermanentBuffs struct {
}
type Purchase struct {
}
type PurchaseLog struct {
	Time    int    `json:"time"`
	Key     string `json:"key"`
	Charges int    `json:"charges"`
}
type Runes struct {
	Property1 int `json:"property1"`
	Property2 int `json:"property2"`
}
type RunesLog struct {
	Time int `json:"time"`
	Key  int `json:"key"`
}
type Sen struct {
}
type SenLeftLog struct {
}
type SenLog struct {
}
type XpReasons struct {
}
type PurchaseTime struct {
}
type FirstPurchaseTime struct {
}
type ItemWin struct {
}
type ItemUsage struct {
}
type PurchaseTpscroll struct {
}
type Benchmarks struct {
}
type Players struct {
	MatchID                 int                     `json:"match_id"`
	PlayerSlot              int                     `json:"player_slot"`
	AbilityUpgradesArr      []int                   `json:"ability_upgrades_arr"`
	AbilityUses             AbilityUses             `json:"ability_uses"`
	AbilityTargets          AbilityTargets          `json:"ability_targets"`
	DamageTargets           DamageTargets           `json:"damage_targets"`
	AccountID               int                     `json:"account_id"`
	Actions                 Actions                 `json:"actions"`
	AdditionalUnits         AdditionalUnits         `json:"additional_units"`
	Assists                 int                     `json:"assists"`
	Backpack0               int                     `json:"backpack_0"`
	Backpack1               int                     `json:"backpack_1"`
	Backpack2               int                     `json:"backpack_2"`
	BuybackLog              []BuybackLog            `json:"buyback_log"`
	CampsStacked            int                     `json:"camps_stacked"`
	ConnectionLog           []ConnectionLog         `json:"connection_log"`
	CreepsStacked           int                     `json:"creeps_stacked"`
	Damage                  Damage                  `json:"damage"`
	DamageInflictor         DamageInflictor         `json:"damage_inflictor"`
	DamageInflictorReceived DamageInflictorReceived `json:"damage_inflictor_received"`
	DamageTaken             DamageTaken             `json:"damage_taken"`
	Deaths                  int                     `json:"deaths"`
	Denies                  int                     `json:"denies"`
	DnT                     []int                   `json:"dn_t"`
	Gold                    int                     `json:"gold"`
	GoldPerMin              int                     `json:"gold_per_min"`
	GoldReasons             GoldReasons             `json:"gold_reasons"`
	GoldSpent               int                     `json:"gold_spent"`
	GoldT                   []int                   `json:"gold_t"`
	HeroDamage              int                     `json:"hero_damage"`
	HeroHealing             int                     `json:"hero_healing"`
	HeroHits                HeroHits                `json:"hero_hits"`
	HeroID                  int                     `json:"hero_id"`
	Item0                   int                     `json:"item_0"`
	Item1                   int                     `json:"item_1"`
	Item2                   int                     `json:"item_2"`
	Item3                   int                     `json:"item_3"`
	Item4                   int                     `json:"item_4"`
	Item5                   int                     `json:"item_5"`
	ItemUses                ItemUses                `json:"item_uses"`
	KillStreaks             KillStreaks             `json:"kill_streaks"`
	Killed                  Killed                  `json:"killed"`
	KilledBy                KilledBy                `json:"killed_by"`
	Kills                   int                     `json:"kills"`
	KillsLog                []KillsLog              `json:"kills_log"`
	LanePos                 LanePos                 `json:"lane_pos"`
	LastHits                int                     `json:"last_hits"`
	LeaverStatus            int                     `json:"leaver_status"`
	Level                   int                     `json:"level"`
	LhT                     []int                   `json:"lh_t"`
	LifeState               LifeState               `json:"life_state"`
	MaxHeroHit              MaxHeroHit              `json:"max_hero_hit"`
	MultiKills              MultiKills              `json:"multi_kills"`
	Obs                     Obs                     `json:"obs"`
	ObsLeftLog              []ObsLeftLog            `json:"obs_left_log"`
	ObsLog                  []ObsLog                `json:"obs_log"`
	ObsPlaced               int                     `json:"obs_placed"`
	PartyID                 int                     `json:"party_id"`
	PermanentBuffs          []PermanentBuffs        `json:"permanent_buffs"`
	Pings                   int                     `json:"pings"`
	Purchase                Purchase                `json:"purchase"`
	PurchaseLog             []PurchaseLog           `json:"purchase_log"`
	RunePickups             int                     `json:"rune_pickups"`
	Runes                   Runes                   `json:"runes"`
	RunesLog                []RunesLog              `json:"runes_log"`
	Sen                     Sen                     `json:"sen"`
	SenLeftLog              []SenLeftLog            `json:"sen_left_log"`
	SenLog                  []SenLog                `json:"sen_log"`
	SenPlaced               int                     `json:"sen_placed"`
	Stuns                   int                     `json:"stuns"`
	Times                   []int                   `json:"times"`
	TowerDamage             int                     `json:"tower_damage"`
	XpPerMin                int                     `json:"xp_per_min"`
	XpReasons               XpReasons               `json:"xp_reasons"`
	XpT                     []int                   `json:"xp_t"`
	Personaname             string                  `json:"personaname"`
	Name                    string                  `json:"name"`
	LastLogin               time.Time               `json:"last_login"`
	RadiantWin              bool                    `json:"radiant_win"`
	StartTime               int                     `json:"start_time"`
	Duration                int                     `json:"duration"`
	Cluster                 int                     `json:"cluster"`
	LobbyType               int                     `json:"lobby_type"`
	GameMode                int                     `json:"game_mode"`
	Patch                   int                     `json:"patch"`
	Region                  int                     `json:"region"`
	IsRadiant               bool                    `json:"isRadiant"`
	Win                     int                     `json:"win"`
	Lose                    int                     `json:"lose"`
	TotalGold               int                     `json:"total_gold"`
	TotalXp                 int                     `json:"total_xp"`
	KillsPerMin             int                     `json:"kills_per_min"`
	Kda                     int                     `json:"kda"`
	Abandons                int                     `json:"abandons"`
	NeutralKills            int                     `json:"neutral_kills"`
	TowerKills              int                     `json:"tower_kills"`
	CourierKills            int                     `json:"courier_kills"`
	LaneKills               int                     `json:"lane_kills"`
	HeroKills               int                     `json:"hero_kills"`
	ObserverKills           int                     `json:"observer_kills"`
	SentryKills             int                     `json:"sentry_kills"`
	RoshanKills             int                     `json:"roshan_kills"`
	NecronomiconKills       int                     `json:"necronomicon_kills"`
	AncientKills            int                     `json:"ancient_kills"`
	BuybackCount            int                     `json:"buyback_count"`
	ObserverUses            int                     `json:"observer_uses"`
	SentryUses              int                     `json:"sentry_uses"`
	LaneEfficiency          int                     `json:"lane_efficiency"`
	LaneEfficiencyPct       int                     `json:"lane_efficiency_pct"`
	Lane                    int                     `json:"lane"`
	LaneRole                int                     `json:"lane_role"`
	IsRoaming               bool                    `json:"is_roaming"`
	PurchaseTime            PurchaseTime            `json:"purchase_time"`
	FirstPurchaseTime       FirstPurchaseTime       `json:"first_purchase_time"`
	ItemWin                 ItemWin                 `json:"item_win"`
	ItemUsage               ItemUsage               `json:"item_usage"`
	PurchaseTpscroll        PurchaseTpscroll        `json:"purchase_tpscroll"`
	ActionsPerMin           int                     `json:"actions_per_min"`
	LifeStateDead           int                     `json:"life_state_dead"`
	RankTier                int                     `json:"rank_tier"`
	Cosmetics               []int                   `json:"cosmetics"`
	Benchmarks              Benchmarks              `json:"benchmarks"`
}
type AllWordCounts struct {
}
type MyWordCounts struct {
}
