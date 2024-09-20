package service

import (
	"context"
	"dota-discord-bot/src/internal/kungdota"
	"dota-discord-bot/src/internal/kungdota/repository"
	"dota-discord-bot/src/internal/opendota"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"unsafe"

	"go.uber.org/zap"
	"gonum.org/v1/gonum/stat/combin"
)

type KungdotaService interface {
	GetPlayersByNames(ctx context.Context, ids []string) (kungdota.Players, error)
	GetPlayersByDiscordIDs(ctx context.Context, ids []string) ([]string, error)
	ShufflePlayers(ctx context.Context, p kungdota.Players) (kungdota.ShuffledTeams, error)
	PostMatch(m opendota.OpenDotaGameObject) error
	// SignUp(ctx context.Context, username string, i int) (map[string][]string, error)
	Update(ctx context.Context, username string) (map[string][]string, error)
}

type kungdotaService struct {
	logger             *zap.Logger
	kungdotaRepository repository.KungdotaRepository
	leagueId           string
}

func NewKungdotaService(logger *zap.Logger, kungdotaRepository repository.KungdotaRepository, leagueId string) KungdotaService {
	return &kungdotaService{
		logger:             logger,
		kungdotaRepository: kungdotaRepository,
		leagueId:           leagueId,
	}
}

func (rx kungdotaService) convert(m opendota.OpenDotaGameObject) (kungdota.Match, error) {
	pAll, _ := rx.kungdotaRepository.GetAllPlayers()

	fbClaimedPlayerSlot := -1
	fbDiedPlayerSlot := -1

	fbClaimedPlayerKungID := -1
	fbDiedPlayerKungID := -1

	var tmpID string

	for _, o := range m.Objectives {
		if o.Type == kungdota.ChatMessageFirstBlood {
			fbClaimedPlayerSlot = o.PlayerSlot

			nbr, ok := o.Key.(float64)
			if ok {
				fbDiedPlayerSlot = int(nbr)
				break
			} else {
				return kungdota.Match{}, errors.New("mayday")
			}

		}
	}

	teams := make([]string, 0)
	coolaStats := make([]kungdota.CoolaStats, 0)
	for _, pOpen := range m.Players {
		for _, p := range pAll.Players {
			tmpID = p.Steam32ID
			if fmt.Sprintf("%d", pOpen.AccountID) == tmpID {
				if pOpen.PlayerSlot == fbClaimedPlayerSlot {
					fbClaimedPlayerKungID = p.ID
				}

				// You might think this looks weird but it actually makes total sense... I think
				if pOpen.PlayerSlot == fbDiedPlayerSlot || pOpen.PlayerSlot-123 == fbDiedPlayerSlot {
					fbDiedPlayerKungID = p.ID
				}

				teams = append(teams, strconv.Itoa(p.ID))
				coolaStats = append(coolaStats, kungdota.CoolaStats{
					ObserverKills: strconv.Itoa(pOpen.ObserverKills),
					SentryKills:   strconv.Itoa(pOpen.SentryKills),
					ObsPlaced:     strconv.Itoa(pOpen.ObserverUses),
					SenPlaced:     strconv.Itoa(pOpen.SentryUses),
					Kills:         strconv.Itoa(pOpen.Kills),
					Assists:       strconv.Itoa(pOpen.Assists),
					Deaths:        strconv.Itoa(pOpen.Deaths),
					// TODO:
					FantasyPoints: "0",
				})
				goto found
			}
		}
		//TODO:
		return kungdota.Match{}, fmt.Errorf("player not found, %d with name %s", pOpen.AccountID, pOpen.Personaname)
	found:
	}

	tmp, ok := m.MatchID.(string)
	if !ok {
		tmp2 := strconv.FormatFloat(m.MatchID.(float64), 'f', -1, 64)
		// tmp2, _ := m.MatchID.(float64)
		tmp = fmt.Sprint(tmp2)
	}

	r := kungdota.Match{
		Teams:             [][]string{teams[:5], teams[5:10]},
		Score:             []int{m.RadiantScore, m.DireScore},
		Winner:            fmt.Sprintf("%d", ^*(*uint64)(unsafe.Pointer(&m.RadiantWin))&1), //Don't arrest me pls
		DotaMatchId:       tmp,
		LeagueId:          rx.leagueId, //,
		DiedFirstBlood:    strconv.Itoa(fbDiedPlayerKungID),
		ClaimedFirstBlood: strconv.Itoa(fbClaimedPlayerKungID),
		CoolaStats:        coolaStats,
	}

	return r, nil
}

func (rx kungdotaService) PostMatch(m opendota.OpenDotaGameObject) error {
	t, err := rx.convert(m)
	if err != nil {
		return err
	}

	return rx.kungdotaRepository.PostMatch(t)
}

func (rx kungdotaService) GetPlayersByNames(ctx context.Context, ids []string) (kungdota.Players, error) {
	return rx.kungdotaRepository.GetByNames(ctx, ids)
}

func (rx kungdotaService) GetPlayersByDiscordIDs(ctx context.Context, ids []string) ([]string, error) {
	return rx.kungdotaRepository.GetByDiscordIDs(ctx, ids)
}

func (rx *kungdotaService) ShufflePlayers(ctx context.Context, p kungdota.Players) (kungdota.ShuffledTeams, error) {
	var teamPairs [][][]int

	// Generate all possible 5v5 team combinations
	allCombinations := combin.Combinations(10, 5)
	for i := 0; i < len(allCombinations)/2; i++ {
		teamPairs = append(teamPairs, [][]int{allCombinations[i], allCombinations[len(allCombinations)-(i+1)]})
	}

	type teamBalance struct {
		elo   int
		teams [][]int
	}

	var balancedTeams []teamBalance
	for _, teamPair := range teamPairs {
		teamOneElo, teamTwoElo := 0.0, 0.0

		for _, idx := range teamPair[0] {
			teamOneElo += float64(p.Players[idx].EloRating)
		}

		for _, idx := range teamPair[1] {
			teamTwoElo += float64(p.Players[idx].EloRating)
		}

		balancedTeams = append(balancedTeams, teamBalance{
			elo:   int(math.Abs(teamOneElo - teamTwoElo)),
			teams: teamPair,
		})
	}

	sort.Slice(balancedTeams, func(i, j int) bool {
		return balancedTeams[i].elo < balancedTeams[j].elo
	})

	// Select one of the top 5 most balanced team pairs randomly
	selectedTeam := balancedTeams[rand.Intn(5)]

	var teamOne, teamTwo []kungdota.Player
	for _, i := range selectedTeam.teams[0] {
		teamOne = append(teamOne, p.Players[i])
	}
	for _, i := range selectedTeam.teams[1] {
		teamTwo = append(teamTwo, p.Players[i])
	}

	var shuffledTeamOne, shuffledTeamTwo kungdota.Players
	if rand.Intn(2) == 1 {
		shuffledTeamOne.Players = teamOne
		shuffledTeamTwo.Players = teamTwo
	} else {
		shuffledTeamOne.Players = teamTwo
		shuffledTeamTwo.Players = teamOne
	}

	return kungdota.ShuffledTeams{
		TeamOne:      shuffledTeamOne,
		TeamTwo:      shuffledTeamTwo,
		FirstPicker:  shuffledTeamOne.Players[rand.Intn(5)],
		SecondPicker: shuffledTeamTwo.Players[rand.Intn(5)],
		EloDiff:      selectedTeam.elo,
	}, nil
}

func (rx *kungdotaService) SignUp(ctx context.Context, username string, i int) (map[string][]string, error) {
	return rx.kungdotaRepository.SignUp(ctx, username, i)
}

func (rx *kungdotaService) Update(ctx context.Context, username string) (map[string][]string, error) {
	return rx.kungdotaRepository.Update(ctx, username)
}
