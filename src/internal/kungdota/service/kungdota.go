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
	GetPlayersByNames(ctx context.Context, ids []string) (kungdota.Players2, error)
	GetPlayersByDiscordIDs(ctx context.Context, ids []string) ([]string, error)
	ShufflePlayers(ctx context.Context, p kungdota.Players2) error
	GetProperties() Properties
	PostMatch(m opendota.OpenDotaGameObject) error
	SignUp(ctx context.Context, username string, i int) (map[string][]string, error)
	Update(ctx context.Context, username string) (map[string][]string, error)
}

type kungdotaService struct {
	logger             *zap.Logger
	kungdotaRepository repository.KungdotaRepository
	properties         Properties
	leagueId           string
}

// priv?
type Properties struct {
	ShuffledTeams kungdota.ShuffledTeams
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

func (rx kungdotaService) GetPlayersByNames(ctx context.Context, ids []string) (kungdota.Players2, error) {
	return rx.kungdotaRepository.GetByNames(ctx, ids)
}

func (rx kungdotaService) GetPlayersByDiscordIDs(ctx context.Context, ids []string) ([]string, error) {
	return rx.kungdotaRepository.GetByDiscordIDs(ctx, ids)
}

func (rx *kungdotaService) GetProperties() Properties {
	return rx.properties
}

func (rx *kungdotaService) ShufflePlayers(ctx context.Context, p kungdota.Players2) error {
	//Magic :)
	pList := make([][][]int, 0)
	list := combin.Combinations(10, 5)
	for i := 0; i < len(list)/2; i++ {
		pList = append(pList, [][]int{list[i], list[len(list)-(i+1)]})
	}

	type tmp2 struct {
		elo   int
		lista [][]int
	}

	tList := make([]tmp2, 0)
	for _, teams := range pList {
		tot1 := 0.0
		tot2 := 0.0

		for _, i := range teams[0] {
			tot1 += float64(p.Players[i].EloRating)
		}

		for _, i := range teams[1] {
			tot2 += float64(p.Players[i].EloRating)
		}

		tList = append(tList, tmp2{
			elo:   int(math.Abs(tot1 - tot2)),
			lista: teams,
		})
	}

	sort.Slice(tList, func(i, j int) bool {
		return tList[i].elo < tList[j].elo
	})

	l := tList[rand.Intn(5)]

	t1 := make([]kungdota.Players, 0)
	t2 := make([]kungdota.Players, 0)

	for _, i := range l.lista[0] {
		t1 = append(t1, p.Players[i])
	}
	for _, i := range l.lista[1] {
		t2 = append(t2, p.Players[i])
	}

	var a kungdota.Players2
	var b kungdota.Players2

	if rand.Intn(2) == 1 {
		a.Players = t1
		b.Players = t2
	} else {
		a.Players = t2
		b.Players = t1
	}

	rx.properties.ShuffledTeams = kungdota.ShuffledTeams{
		TeamOne:      a,
		TeamTwo:      b,
		FirstPicker:  a.Players[rand.Intn(5)],
		SecondPicker: b.Players[rand.Intn(5)],
		EloDiff:      tasd.elo,
	}

	return nil
}

func (rx *kungdotaService) SignUp(ctx context.Context, username string, i int) (map[string][]string, error) {
	return rx.kungdotaRepository.SignUp(ctx, username, i)
}

func (rx *kungdotaService) Update(ctx context.Context, username string) (map[string][]string, error) {
	return rx.kungdotaRepository.Update(ctx, username)
}
